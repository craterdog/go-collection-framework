/*
................................................................................
.    Copyright (c) 2009-2024 Crater Dog Technologies.  All Rights Reserved.    .
................................................................................
.  DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               .
.                                                                              .
.  This code is free software; you can redistribute it and/or modify it under  .
.  the terms of The MIT License (MIT), as published by the Open Source         .
.  Initiative. (See https://opensource.org/license/MIT)                        .
................................................................................
*/

package agent

import (
	fmt "fmt"
	cmp "math/cmplx"
	ref "reflect"
	sts "strings"
	syn "sync"
)

// CLASS ACCESS

// Reference

var collatorClass = map[string]any{}
var collatorMutex syn.Mutex

// Function

func Collator[V any]() CollatorClassLike[V] {
	// Generate the name of the bound class type.
	var class *collatorClass_[V]
	var name = fmt.Sprintf("%T", class)

	// Check for existing bound class type.
	collatorMutex.Lock()
	var value = collatorClass[name]
	switch actual := value.(type) {
	case *collatorClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Add a new bound class type.
		class = &collatorClass_[V]{
			defaultMaximum_: 16,
		}
		collatorClass[name] = class
	}
	collatorMutex.Unlock()

	// Return a reference to the bound class type.
	return class
}

// CLASS METHODS

// Target

type collatorClass_[V any] struct {
	defaultMaximum_ int
}

// Constants

func (c *collatorClass_[V]) DefaultMaximum() int {
	return c.defaultMaximum_
}

// Constructors

func (c *collatorClass_[V]) Make() CollatorLike[V] {
	return &collator_[V]{
		maximum_: c.defaultMaximum_,
	}
}

func (c *collatorClass_[V]) MakeWithMaximum(maximum int) CollatorLike[V] {
	if maximum < 0 {
		maximum = c.defaultMaximum_
	}
	return &collator_[V]{
		maximum_: maximum,
	}
}

// INSTANCE METHODS

// Target

type collator_[V any] struct {
	class_   CollatorClassLike[V]
	depth_   int
	maximum_ int
}

// Attributes

func (v *collator_[V]) GetClass() CollatorClassLike[V] {
	return v.class_
}

func (v *collator_[V]) GetDepth() int {
	return v.depth_
}

func (v *collator_[V]) GetMaximum() int {
	return v.maximum_
}

// Public

func (v *collator_[V]) CompareValues(first V, second V) bool {
	return v.compareValues(ref.ValueOf(first), ref.ValueOf(second))
}

func (v *collator_[V]) RankValues(first V, second V) Rank {
	return v.rankValues(ref.ValueOf(first), ref.ValueOf(second))
}

// Private

func (v *collator_[V]) compareArrays(first ref.Value, second ref.Value) bool {
	// Check for maximum traversal depth.
	if v.depth_ == v.maximum_ {
		panic(fmt.Sprintf("The maximum traversal depth was exceeded: %v", v.depth_))
	}

	// Compare the sizes of the Go arrays.
	var size = first.Len()
	if second.Len() != size {
		// The Go arrays are of different lengths.
		return false
	}

	// Compare the values of the Go arrays.
	for i := 0; i < size; i++ {
		v.depth_++
		if !v.compareValues(first.Index(i), second.Index(i)) {
			// Two of the values in the Go arrays are different.
			v.depth_--
			return false
		}
		v.depth_--
	}
	return true
}

func (v *collator_[V]) compareInterfaces(first ref.Value, second ref.Value) bool {
	var typeRef = first.Type() // We know the structures are the same type.
	var count = typeRef.NumMethod()
	for index := 0; index < count; index++ {
		var name = typeRef.Method(index).Name
		var arguments = first.Method(index).Type().NumIn()
		if sts.HasPrefix(name, "Get") && arguments == 0 {
			var firstValue = first.Method(index).Call([]ref.Value{})[0]
			var secondValue = second.Method(index).Call([]ref.Value{})[0]
			if !v.compareValues(firstValue, secondValue) {
				// Found a difference.
				return false
			}
		}
	}
	// All getter values are equal.
	return true
}

func (v *collator_[V]) compareMaps(first ref.Value, second ref.Value) bool {
	// Check for maximum traversal depth.
	if v.depth_ == v.maximum_ {
		panic(fmt.Sprintf("The maximum traversal depth was exceeded: %v", v.depth_))
	}

	// Compare the sizes of the two Go maps.
	if first.Len() != second.Len() {
		// The Go maps are different sizes.
		return false
	}

	// Compare the keys and values for the two Go maps.
	var iterator = first.MapRange()
	for iterator.Next() {
		v.depth_++
		var key = iterator.Key()
		var firstValue = iterator.Value()
		var secondValue = second.MapIndex(key)
		if !v.compareValues(firstValue, secondValue) {
			// The values don't match.
			v.depth_--
			return false
		}
		v.depth_--
	}
	return true
}

func (v *collator_[V]) compareIntrinsics(first, second ref.Value) bool {
	return first.Interface() == second.Interface()
}

func (v *collator_[V]) compareSequences(first ref.Value, second ref.Value) bool {
	// Compare the Go arrays for the two sequences.
	var firstArray = first.MethodByName("AsArray").Call([]ref.Value{})[0]
	var secondArray = second.MethodByName("AsArray").Call([]ref.Value{})[0]
	return v.compareArrays(firstArray, secondArray)
}

func (v *collator_[V]) compareValues(first ref.Value, second ref.Value) bool {
	// Handle any invalid values.
	if !first.IsValid() {
		return !second.IsValid()
	}
	if !second.IsValid() {
		return false
	}

	// At this point, neither of the values are invalid.
	var firstType = v.getType(first.Type())
	var secondType = v.getType(second.Type())
	if firstType != secondType && firstType != "any" && secondType != "any" {
		// The values have different types.
		return false
	}

	// We now know that the types of the values are the same, and neither of
	// the values is invalid.
	switch first.Kind() {

	// Handle all intrinsic primitive types.
	case ref.Bool,
		ref.Uint8, ref.Uint16, ref.Uint32, ref.Uint64, ref.Uint,
		ref.Int8, ref.Int16, ref.Int32, ref.Int64, ref.Int,
		ref.Float32, ref.Float64, ref.Complex64, ref.Complex128,
		ref.String:
		return v.compareIntrinsics(first, second)

	// Handle all intrinsic collection types.
	case ref.Array, ref.Slice:
		switch {
		case first.IsNil():
			return second.IsNil()
		case second.IsNil():
			return false // We know that first isn't nil.
		default:
			return v.compareArrays(first, second)
		}
	case ref.Map:
		switch {
		case first.IsNil():
			return second.IsNil()
		case second.IsNil():
			return false // We know that first isn't nil.
		default:
			return v.compareMaps(first, second)
		}

	// Handle all interfaces and pointers.
	case ref.Interface, ref.Pointer:
		switch {
		case first.IsNil():
			return second.IsNil()
		case second.IsNil():
			return false // We know that first isn't nil.
		case first.MethodByName("AsArray").IsValid():
			// The value is a sequence.
			return v.compareSequences(first, second)
		case first.NumMethod() > 0:
			// The value is an interface or pointer to a structure with methods.
			return v.compareInterfaces(first, second)
		default:
			// The values are pointers to the values to be compared.
			first = first.Elem()
			second = second.Elem()
			return v.compareValues(first, second)
		}

	// Handle all Go structures.
	case ref.Struct:
		// The Go comparison operator performs a deep comparison on structures.
		return first.Interface() == second.Interface()

	default:
		panic(fmt.Sprintf(
			"Attempted to compare:\n    first: %v\n    type: %v\n    kind: %v\nand\n    second: %v\n    type: %v\n    kind: %v\n",
			first.Interface(),
			first.Type(),
			first.Kind(),
			second.Interface(),
			second.Type(),
			second.Kind()))
	}
}

func (v *collator_[V]) getType(type_ ref.Type) string {
	var result = type_.String()
	result = sts.TrimPrefix(result, "*")
	if sts.HasPrefix(result, "[") {
		result = "array"
	}
	if sts.HasPrefix(result, "map[") {
		result = "map"
	}
	var index = sts.Index(result, "[")
	if index > -1 {
		result = result[:index]
	}
	if result == "interface {}" {
		result = "any"
	}
	if sts.HasPrefix(result, "bool") {
		result = "boolean"
	}
	if sts.HasPrefix(result, "uint8") {
		result = "byte"
	}
	if sts.HasPrefix(result, "int32") {
		result = "rune"
	}
	if sts.HasPrefix(result, "uint") {
		result = "unsigned"
	}
	if sts.HasPrefix(result, "int") {
		result = "integer"
	}
	if sts.HasPrefix(result, "float") {
		result = "float"
	}
	if sts.HasPrefix(result, "complex") {
		result = "complex"
	}
	return result
}

func (v *collator_[V]) rankArrays(first ref.Value, second ref.Value) Rank {
	// Check for maximum traversal depth.
	if v.depth_ == v.maximum_ {
		panic(fmt.Sprintf("The maximum traversal depth was exceeded: %v", v.depth_))
	}

	// Determine the smallest Go array.
	var firstSize = first.Len()
	var secondSize = second.Len()
	if firstSize > secondSize {
		// Swap the order of the Go arrays and reverse the result.
		switch v.rankArrays(second, first) {
		case LesserRank:
			return GreaterRank
		case GreaterRank:
			return LesserRank
		default:
			return EqualRank
		}
	}

	// Iterate through the smallest Go array.
	for i := 0; i < firstSize; i++ {
		v.depth_++
		var rank = v.rankValues(first.Index(i), second.Index(i))
		if rank != EqualRank {
			// The values are different.
			v.depth_--
			return rank
		}
		// The two values match.
		v.depth_--
	}

	// The Go arrays contain the same initial values.
	if secondSize > firstSize {
		// The shorter Go array is ranked before the longer Go array.
		return LesserRank
	}
	// The Go arrays are the same length and contain the same values.
	return EqualRank
}

func (v *collator_[V]) rankBooleans(first, second bool) Rank {
	if !first && second {
		return LesserRank
	}
	if first && !second {
		return GreaterRank
	}
	return EqualRank
}

func (v *collator_[V]) rankBytes(first, second byte) Rank {
	if first < second {
		return LesserRank
	}
	if first > second {
		return GreaterRank
	}
	return EqualRank
}

func (v *collator_[V]) rankComplex(first, second complex128) Rank {
	if first == second {
		return EqualRank
	}
	switch {
	case cmp.Abs(first) < cmp.Abs(second):
		// The magnitude of the first vector is less than the second.
		return LesserRank
	case cmp.Abs(first) > cmp.Abs(second):
		// The magnitude of the first vector is greater than the second.
		return GreaterRank
	default:
		// The magnitudes of the vectors are equal.
		switch {
		case cmp.Phase(first) < cmp.Phase(second):
			// The phase of the first vector is less than the second.
			return LesserRank
		case cmp.Phase(first) > cmp.Phase(second):
			// The phase of the first vector is greater than the second.
			return GreaterRank
		default:
			// The phases of the vectors are also equal.
			return EqualRank
		}
	}
}

func (v *collator_[V]) rankFloats(first, second float64) Rank {
	if first < second {
		return LesserRank
	}
	if first > second {
		return GreaterRank
	}
	return EqualRank
}

func (v *collator_[V]) rankInterfaces(first ref.Value, second ref.Value) Rank {
	var typeRef = first.Type() // We know the structures are the same type.
	var count = first.NumMethod()
	for index := 0; index < count; index++ {
		var method = typeRef.Method(index)
		if sts.HasPrefix(method.Name, "Get") {
			var firstValue = first.Method(index).Call([]ref.Value{})[0]
			var secondValue = second.Method(index).Call([]ref.Value{})[0]
			var rank = v.rankValues(firstValue, secondValue)
			if rank != EqualRank {
				// Found a difference.
				return rank
			}
		}
	}
	// All getter values are equal.
	return EqualRank
}

// NOTE:
// Currently the implementation of Go maps is hashtable based. The  order of the
// keys is random, even for two Go maps with the same keys if the associations
// were entered in different sequences.  Therefore at this time it is necessary
// to sort the key arrays for each Go map.  This introduces a circular
// dependency between the implementation of the collator and the sorter types:
//
// rankMaps() -> SortValues() -> RankingFunction
func (v *collator_[V]) rankMaps(first ref.Value, second ref.Value) Rank {
	// Check for maximum traversal depth.
	if v.depth_ == v.maximum_ {
		panic(fmt.Sprintf("The maximum traversal depth was exceeded: %v", v.depth_))
	}

	// Extract and sort the keys for the two Go maps.
	var sorter = Sorter[ref.Value]().MakeWithRanker(v.rankValues)
	var firstKeys = first.MapKeys() // The returned keys are in random order.
	sorter.SortValues(firstKeys)
	var secondKeys = second.MapKeys() // The returned keys are in random order.
	sorter.SortValues(secondKeys)

	// Determine the smallest Go map.
	var firstSize = len(firstKeys)
	var secondSize = len(secondKeys)
	if firstSize > secondSize {
		// Swap the order of the Go maps and reverse the result.
		switch v.rankMaps(second, first) {
		case LesserRank:
			return GreaterRank
		case GreaterRank:
			return LesserRank
		default:
			return EqualRank
		}
	}

	// Iterate through the smallest Go map.
	for i := 0; i < firstSize; i++ {
		v.depth_++

		// Rank the two keys.
		var firstKey = firstKeys[i]
		var secondKey = secondKeys[i]
		var keyRank = v.rankValues(firstKey, secondKey)
		if keyRank != EqualRank {
			// The two keys are different.
			v.depth_--
			return keyRank
		}

		// The two keys match so rank the corresponding values.
		var firstValue = first.MapIndex(firstKey)
		var secondValue = second.MapIndex(secondKey)
		var valueRank = v.rankValues(firstValue, secondValue)
		if valueRank != EqualRank {
			// The two values are different.
			v.depth_--
			return valueRank
		}
		v.depth_--
	}

	// The Go maps contain the same initial associations.
	if secondSize > firstSize {
		// The shorter Go map is ranked before the longer Go map.
		return LesserRank
	}

	// All keys and values match.
	return EqualRank
}

func (v *collator_[V]) rankIntrinsics(first, second ref.Value) Rank {
	var firstValue = first.Interface()
	var secondValue = second.Interface()
	switch first.Kind() {
	case ref.Bool:
		var firstBoolean = bool(first.Bool())
		var secondBoolean = bool(second.Bool())
		return v.rankBooleans(firstBoolean, secondBoolean)
	case ref.Uint8: // Byte
		var firstByte = byte(first.Uint())
		var secondByte = byte(second.Uint())
		return v.rankBytes(firstByte, secondByte)
	case ref.Uint16, ref.Uint32, ref.Uint64, ref.Uint:
		var firstUnsigned = uint64(first.Uint())
		var secondUnsigned = uint64(second.Uint())
		return v.rankUnsigned(firstUnsigned, secondUnsigned)
	case ref.Int8, ref.Int16, ref.Int64, ref.Int:
		var firstSigned = int64(first.Int())
		var secondSigned = int64(second.Int())
		return v.rankSigned(firstSigned, secondSigned)
	case ref.Float32, ref.Float64:
		var firstFloat = float64(first.Float())
		var secondFloat = float64(second.Float())
		return v.rankFloats(firstFloat, secondFloat)
	case ref.Complex64, ref.Complex128:
		var firstComplex = complex128(first.Complex())
		var secondComplex = complex128(second.Complex())
		return v.rankComplex(firstComplex, secondComplex)
	case ref.Int32: // Runes
		var firstRune = rune(first.Int())
		var secondRune = rune(second.Int())
		return v.rankRunes(rune(firstRune), rune(secondRune))
	case ref.String:
		var firstString = string(first.String())
		var secondString = string(second.String())
		return v.rankStrings(firstString, secondString)
	default:
		var message = fmt.Sprintf("Attempted to rank %v(%T) and %v(%T)", firstValue, firstValue, secondValue, secondValue)
		panic(message)
	}
}

func (v *collator_[V]) rankRunes(first, second int32) Rank {
	if first < second {
		return LesserRank
	}
	if first > second {
		return GreaterRank
	}
	return EqualRank
}

func (v *collator_[V]) rankSequences(first ref.Value, second ref.Value) Rank {
	// Rank the Go arrays for the two sequences.
	var firstArray = first.MethodByName("AsArray").Call([]ref.Value{})[0]
	var secondArray = second.MethodByName("AsArray").Call([]ref.Value{})[0]
	return v.rankArrays(firstArray, secondArray)
}

func (v *collator_[V]) rankSigned(first, second int64) Rank {
	if first < second {
		return LesserRank
	}
	if first > second {
		return GreaterRank
	}
	return EqualRank
}

func (v *collator_[V]) rankStrings(first, second string) Rank {
	if first < second {
		// The first string comes before the second string alphabetically.
		return LesserRank
	}
	if first > second {
		// The first string comes after the second string alphabetically.
		return GreaterRank
	}
	// The two strings are the same.
	return EqualRank
}

func (v *collator_[V]) rankStructures(first ref.Value, second ref.Value) Rank {
	var count = first.NumField() // The structures are the same type.
	for index := 0; index < count; index++ {
		var firstField = first.Field(index)
		var secondField = second.Field(index)
		if firstField.CanInterface() {
			var rank = v.rankValues(firstField, secondField)
			if rank != EqualRank {
				// Found a difference.
				return rank
			}
		}
	}
	// All fields have matching values.
	return EqualRank
}

func (v *collator_[V]) rankUnsigned(first, second uint64) Rank {
	if first < second {
		return LesserRank
	}
	if first > second {
		return GreaterRank
	}
	return EqualRank
}

func (v *collator_[V]) rankValues(first ref.Value, second ref.Value) Rank {
	// Handle any nil pointers.
	if !first.IsValid() {
		if !second.IsValid() {
			// Both values are nil.
			return EqualRank
		}
		// Only the first value is nil.
		return LesserRank
	} else if !second.IsValid() {
		// Only the second value is nil.
		return GreaterRank
	}

	// At this point, neither of the values are nil.
	var firstType = v.getType(first.Type())
	var secondType = v.getType(second.Type())
	if firstType != secondType && firstType != "any" && secondType != "any" {
		// The values have different types.
		return v.rankStrings(firstType, secondType)
	}

	// We now know that the types of the values are the same, and neither of
	// the values is nil.
	switch first.Kind() {

	// Handle all intrinsic primitive types.
	case ref.Bool,
		ref.Uint8, ref.Uint16, ref.Uint32, ref.Uint64, ref.Uint,
		ref.Int8, ref.Int16, ref.Int32, ref.Int64, ref.Int,
		ref.Float32, ref.Float64, ref.Complex64, ref.Complex128,
		ref.String:
		return v.rankIntrinsics(first, second)

	// Handle all intrinsic collection types.
	case ref.Array, ref.Slice:
		switch {
		case first.IsNil():
			if second.IsNil() {
				return EqualRank
			}
			// Only the first value is nil.
			return LesserRank
		case second.IsNil():
			return GreaterRank // We know that first isn't nil.
		default:
			return v.rankArrays(first, second)
		}
	case ref.Map:
		switch {
		case first.IsNil():
			if second.IsNil() {
				return EqualRank
			}
			// Only the first value is nil.
			return LesserRank
		case second.IsNil():
			return GreaterRank // We know that first isn't nil.
		default:
			return v.rankMaps(first, second)
		}

	// Handle all interfaces and pointers.
	case ref.Interface, ref.Pointer:
		switch {
		case first.IsNil():
			if second.IsNil() {
				return EqualRank
			}
			// Only the first value is nil.
			return LesserRank
		case second.IsNil():
			return GreaterRank // We know that first isn't nil.
		case first.MethodByName("AsArray").IsValid():
			// The value is a collection.
			return v.rankSequences(first, second)
		case first.NumMethod() > 0:
			// The value is an interface or pointer to a structure with methods.
			return v.rankInterfaces(first, second)
		default:
			// The values are pointers to the values to be ranked.
			first = first.Elem()
			second = second.Elem()
			return v.rankValues(first, second)
		}

	// Handle all Go structures.
	case ref.Struct:
		// Rank the corresponding fields for each structure.
		var ranking = v.rankStructures(first, second)
		if ranking != EqualRank {
			return ranking
		}
		// Rank the corresponding getter values for each structure.
		return v.rankInterfaces(first, second)

	default:
		panic(fmt.Sprintf(
			"Attempted to rank:\n    first: %v\n    type: %v\n    kind: %v\nand\n    second: %v\n    type: %v\n    kind: %v\n",
			first.Interface(),
			first.Type(),
			first.Kind(),
			second.Interface(),
			second.Type(),
			second.Kind()))
	}
}
