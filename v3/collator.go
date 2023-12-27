/*******************************************************************************
 *   Copyright (c) 2009-2024 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package collections

import (
	fmt "fmt"
	cmp "math/cmplx"
	ref "reflect"
	stc "strconv"
	sts "strings"
)

// CLASS NAMESPACE

// Private Class Namespace Type

type collatorClass_ struct {
	defaultDepth int
}

// Private Namespace Reference(s)

var collatorClass = &collatorClass_{
	defaultDepth: 16,
}

// Public Namespace Access

func CollatorClass() *collatorClass_ {
	return collatorClass
}

// Public Class Constants

func (c *collatorClass_) GetDefaultDepth() int {
	return c.defaultDepth
}

// Public Class Constructors

func (c *collatorClass_) Default() *collator_ {
	var collator = &collator_{
		maximum: c.defaultDepth,
	}
	return collator
}

func (c *collatorClass_) WithDepth(depth int) *collator_ {
	if depth < 0 || depth > c.defaultDepth {
		depth = c.defaultDepth
	}
	var collator = &collator_{
		maximum: depth,
	}
	return collator
}

// CLASS TYPE

// Private Class Type Definition

type collator_ struct {
	depth   int
	maximum int
}

// Discerning Interface

func (v *collator_) CompareValues(first Value, second Value) bool {
	return v.compareValues(ref.ValueOf(first), ref.ValueOf(second))
}

func (v *collator_) RankValues(first Value, second Value) int {
	return v.rankValues(ref.ValueOf(first), ref.ValueOf(second))
}

// Private Interface

// This private class method determines whether or not the specified Go arrays
// have the same value.
func (v *collator_) compareArrays(first ref.Value, second ref.Value) bool {
	// Check for maximum traversal depth.
	if v.depth == v.maximum {
		panic(fmt.Sprintf("The maximum traversal depth was exceeded: %v", v.depth))
	}

	// Compare the sizes of the Go arrays.
	var size = first.Len()
	if second.Len() != size {
		// The Go arrays are of different lengths.
		return false
	}

	// Compare the values of the Go arrays.
	for i := 0; i < size; i++ {
		v.depth++
		if !v.compareValues(first.Index(i), second.Index(i)) {
			// Two of the values in the Go arrays are different.
			v.depth--
			return false
		}
		v.depth--
	}
	return true
}

// This private class method determines whether or not the specified primitives
// have the same value.
func (v *collator_) comparePrimitives(first, second ref.Value) bool {
	return first.Interface() == second.Interface()
}

// This private class method determines whether or not the specified interfaces
// have the same getter values.
func (v *collator_) compareInterfaces(first ref.Value, second ref.Value) bool {
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

// This private class method determines whether or not the specified Go maps
// have the same value.
func (v *collator_) compareMaps(first ref.Value, second ref.Value) bool {
	// Check for maximum traversal depth.
	if v.depth == v.maximum {
		panic(fmt.Sprintf("The maximum traversal depth was exceeded: %v", v.depth))
	}

	// Compare the sizes of the two Go maps.
	if first.Len() != second.Len() {
		// The Go maps are different sizes.
		return false
	}

	// Compare the keys and values for the two Go maps.
	var iterator = first.MapRange()
	for iterator.Next() {
		v.depth++
		var key = iterator.Key()
		var firstValue = iterator.Value()
		var secondValue = second.MapIndex(key)
		if !v.compareValues(firstValue, secondValue) {
			// The values don't match.
			v.depth--
			return false
		}
		v.depth--
	}
	return true
}

// This private class method determines whether or not the specified sequences
// have the same value.
func (v *collator_) compareSequences(first ref.Value, second ref.Value) bool {
	// Compare the Go arrays for the two sequences.
	var firstArray = first.MethodByName("AsArray").Call([]ref.Value{})[0]
	var secondArray = second.MethodByName("AsArray").Call([]ref.Value{})[0]
	return v.compareArrays(firstArray, secondArray)
}

// This private class method determines whether or not the specified reflective
// values are equal using reflection and a recursive descent algorithm.
func (v *collator_) compareValues(first ref.Value, second ref.Value) bool {
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

	// Handle all primitive types.
	case ref.Bool,
		ref.Uint8, ref.Uint16, ref.Uint32, ref.Uint64, ref.Uint,
		ref.Int8, ref.Int16, ref.Int32, ref.Int64, ref.Int,
		ref.Float32, ref.Float64, ref.Complex64, ref.Complex128,
		ref.String:
		return v.comparePrimitives(first, second)

	// Handle all primitive collection types.
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

// This private class method removes the generics from the type string for the
// specified type and converts an empty interface into type "any".
func (v *collator_) getType(type_ ref.Type) string {
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

// This private class method returns the ranking order of the specified Go
// arrays using a recursive descent algorithm.
func (v *collator_) rankArrays(first ref.Value, second ref.Value) int {
	// Check for maximum traversal depth.
	if v.depth == v.maximum {
		panic(fmt.Sprintf("The maximum traversal depth was exceeded: %v", v.depth))
	}

	// Determine the smallest Go array.
	var firstSize = first.Len()
	var secondSize = second.Len()
	if firstSize > secondSize {
		// Swap the order of the Go arrays and reverse the sign of the result.
		return -1 * v.rankArrays(second, first)
	}

	// Iterate through the smallest Go array.
	for i := 0; i < firstSize; i++ {
		v.depth++
		var rank = v.rankValues(first.Index(i), second.Index(i))
		if rank < 0 {
			// The value in the first Go array comes before its matching value.
			v.depth--
			return -1
		}
		if rank > 0 {
			// The value in the first Go array comes after its matching value.
			v.depth--
			return 1
		}
		// The two values match.
		v.depth--
	}

	// The Go arrays contain the same initial values.
	if secondSize > firstSize {
		// The shorter Go array is ranked before the longer Go array.
		return -1
	}
	// The Go arrays are the same length and contain the same values.
	return 0
}

// This private class method returns the ranking order of the specified boolean
// values.
func (v *collator_) rankBooleans(first, second bool) int {
	if !first && second {
		return -1
	}
	if first && !second {
		return 1
	}
	return 0
}

// This private class method returns the ranking order of the specified complex
// number values.
func (v *collator_) rankComplex(first, second complex128) int {
	if first == second {
		return 0
	}
	switch {
	case cmp.Abs(first) < cmp.Abs(second):
		// The magnitude of the first vector is less than the second.
		return -1
	case cmp.Abs(first) > cmp.Abs(second):
		// The magnitude of the first vector is greater than the second.
		return 1
	default:
		// The magnitudes of the vectors are equal.
		switch {
		case cmp.Phase(first) < cmp.Phase(second):
			// The phase of the first vector is less than the second.
			return -1
		case cmp.Phase(first) > cmp.Phase(second):
			// The phase of the first vector is greater than the second.
			return 1
		default:
			// The phases of the vectors are also equal.
			return 0
		}
	}
}

// This private class method returns the ranking order of the specified
// primitives.  NOTE: Go does not provide an easy way to a apply a possible tilde
// type (e.g. ~string) to its primitive (named) type without knowing whether or
// not it is actually a tilde type. So we must convert it to a string and then
// parse it back again...
// This method attempts to hide that ugliness from the rest of the code.
func (v *collator_) rankPrimitives(first, second ref.Value) int {
	var firstValue = first.Interface()
	var secondValue = second.Interface()
	switch first.Kind() {
	case ref.Bool:
		var firstBoolean, _ = stc.ParseBool(fmt.Sprintf("%v", firstValue))
		var secondBoolean, _ = stc.ParseBool(fmt.Sprintf("%v", secondValue))
		return v.rankBooleans(firstBoolean, secondBoolean)
	case ref.Uint8, ref.Uint16, ref.Uint32, ref.Uint64, ref.Uint:
		var firstUnsigned, _ = stc.ParseUint(fmt.Sprintf("%v", firstValue), 10, 64)
		var secondUnsigned, _ = stc.ParseUint(fmt.Sprintf("%v", secondValue), 10, 64)
		return v.rankUnsigned(firstUnsigned, secondUnsigned)
	case ref.Int8, ref.Int16, ref.Int64, ref.Int:
		var firstSigned, _ = stc.ParseInt(fmt.Sprintf("%v", firstValue), 10, 64)
		var secondSigned, _ = stc.ParseInt(fmt.Sprintf("%v", secondValue), 10, 64)
		return v.rankSigned(firstSigned, secondSigned)
	case ref.Float32, ref.Float64:
		var firstFloat, _ = stc.ParseFloat(fmt.Sprintf("%v", firstValue), 64)
		var secondFloat, _ = stc.ParseFloat(fmt.Sprintf("%v", secondValue), 64)
		return v.rankFloats(firstFloat, secondFloat)
	case ref.Complex64, ref.Complex128:
		var firstComplex, _ = stc.ParseComplex(fmt.Sprintf("%v", firstValue), 128)
		var secondComplex, _ = stc.ParseComplex(fmt.Sprintf("%v", secondValue), 128)
		return v.rankComplex(firstComplex, secondComplex)
	case ref.Int32: // Runes
		var firstRune, _ = stc.ParseInt(fmt.Sprintf("%v", firstValue), 10, 32)
		var secondRune, _ = stc.ParseInt(fmt.Sprintf("%v", secondValue), 10, 32)
		return v.rankRunes(rune(firstRune), rune(secondRune))
	case ref.String:
		var firstString = fmt.Sprintf("%v", firstValue)
		var secondString = fmt.Sprintf("%v", secondValue)
		return v.rankStrings(firstString, secondString)
	default:
		var message = fmt.Sprintf("Attempted to rank %v(%T) and %v(%T)", firstValue, firstValue, secondValue, secondValue)
		panic(message)
	}
}

// This private class method returns the ranking order of the specified floating
// point numbers.
func (v *collator_) rankFloats(first, second float64) int {
	if first < second {
		return -1
	}
	if first > second {
		return 1
	}
	return 0
}

// This private class method returns the ranking order of the specified
// interfaces by ranking the results of their getter methods.
func (v *collator_) rankInterfaces(first ref.Value, second ref.Value) int {
	var typeRef = first.Type() // We know the structures are the same type.
	var count = first.NumMethod()
	for index := 0; index < count; index++ {
		var method = typeRef.Method(index)
		if sts.HasPrefix(method.Name, "Get") {
			var firstValue = first.Method(index).Call([]ref.Value{})[0]
			var secondValue = second.Method(index).Call([]ref.Value{})[0]
			var ranking = v.rankValues(firstValue, secondValue)
			if ranking != 0 {
				// Found a difference.
				return ranking
			}
		}
	}
	// All getter values are equal.
	return 0
}

// This private class method returns the ranking order of the specified Go maps
// using a recursive descent algorithm. NOTE: currently the implementation of Go
// maps is hashtable based. The order of the keys is random, even for two Go
// maps with the same keys if the associations were entered in different
// sequences.  Therefore at this time it is necessary to sort the key arrays for
// each Go map.  This introduces a circular dependency between the implementation
// of the collator and the sorter types:
//
//	rankMaps() -> SortValues() -> RankingFunction
func (v *collator_) rankMaps(first ref.Value, second ref.Value) int {
	// Check for maximum traversal depth.
	if v.depth == v.maximum {
		panic(fmt.Sprintf("The maximum traversal depth was exceeded: %v", v.depth))
	}

	// Extract and sort the keys for the two Go maps.
	var sorter = SorterClass[ref.Value]().WithRanker(v.rankReflective)
	var firstKeys = first.MapKeys() // The returned keys are in random order.
	sorter.SortValues(firstKeys)
	var secondKeys = second.MapKeys() // The returned keys are in random order.
	sorter.SortValues(secondKeys)

	// Determine the smallest Go map.
	var firstSize = len(firstKeys)
	var secondSize = len(secondKeys)
	if firstSize > secondSize {
		// Swap the order of the Go maps and reverse the sign of the result.
		return -1 * v.rankMaps(second, first)
	}

	// Iterate through the smallest Go map.
	for i := 0; i < firstSize; i++ {
		v.depth++

		// Rank the two keys.
		var firstKey = firstKeys[i]
		var secondKey = secondKeys[i]
		var keyRank = v.rankValues(firstKey, secondKey)
		if keyRank < 0 {
			// The key in the first Go map comes before its matching key.
			v.depth--
			return -1
		}
		if keyRank > 0 {
			// The key in the first Go map comes after its matching key.
			v.depth--
			return 1
		}

		// The two keys match so rank the corresponding values.
		var firstValue = first.MapIndex(firstKey)
		var secondValue = second.MapIndex(secondKey)
		var valueRank = v.rankValues(firstValue, secondValue)
		if valueRank < 0 {
			// The value in the first Go map comes before its matching value.
			v.depth--
			return -1
		}
		if valueRank > 0 {
			// The value in the first Go map comes after its matching value.
			v.depth--
			return 1
		}
		v.depth--
	}

	// The Go maps contain the same initial associations.
	if secondSize > firstSize {
		// The shorter Go map is ranked before the longer Go map.
		return -1
	}

	// All keys and values match.
	return 0
}

// This private class method returns the relative ranking of the specified
// values using reflection.
func (v *collator_) rankReflective(first Value, second Value) int {
	var firstValue = first.(ref.Value)
	var secondValue = second.(ref.Value)
	return v.rankValues(firstValue, secondValue)
}

// This private class method returns the ranking order of the specified runes.
func (v *collator_) rankRunes(first, second int32) int {
	if first < second {
		return -1
	}
	if first > second {
		return 1
	}
	return 0
}

// This private class method returns the ranking order of the specified
// sequences using a recursive descent algorithm.
func (v *collator_) rankSequences(first ref.Value, second ref.Value) int {
	// Rank the Go arrays for the two sequences.
	var firstArray = first.MethodByName("AsArray").Call([]ref.Value{})[0]
	var secondArray = second.MethodByName("AsArray").Call([]ref.Value{})[0]
	return v.rankArrays(firstArray, secondArray)
}

// This private class method returns the ranking order of the specified signed
// integers.
func (v *collator_) rankSigned(first, second int64) int {
	if first < second {
		return -1
	}
	if first > second {
		return 1
	}
	return 0
}

// This private class method returns the ranking order of the specified string
// values.
func (v *collator_) rankStrings(first, second string) int {
	if first < second {
		// The first string comes before the second string alphabetically.
		return -1
	}
	if first > second {
		// The first string comes after the second string alphabetically.
		return 1
	}
	// The two strings are the same.
	return 0
}

// This private class method returns the ranking order of the specified
// structures by ranking the associated fields.
func (v *collator_) rankStructures(first ref.Value, second ref.Value) int {
	var count = first.NumField() // The structures are the same type.
	for index := 0; index < count; index++ {
		var firstField = first.Field(index)
		var secondField = second.Field(index)
		if firstField.CanInterface() {
			var ranking = v.rankValues(firstField, secondField)
			if ranking != 0 {
				// Found a difference.
				return ranking
			}
		}
	}
	// All fields have matching values.
	return 0
}

// This private class method returns the ranking order of the specified unsigned
// integers.
func (v *collator_) rankUnsigned(first, second uint64) int {
	if first < second {
		return -1
	}
	if first > second {
		return 1
	}
	return 0
}

// This private method returns the ranking order of the specified values using
// reflection and a recursive descent algorithm.
func (v *collator_) rankValues(first ref.Value, second ref.Value) int {
	// Handle any nil pointers.
	if !first.IsValid() {
		if !second.IsValid() {
			// Both values are nil.
			return 0
		}
		// Only the first value is nil.
		return -1
	} else if !second.IsValid() {
		// Only the second value is nil.
		return 1
	}

	// At this point, neither of the values are nil.
	var firstType = v.getType(first.Type())
	var secondType = v.getType(second.Type())
	if firstType != secondType && firstType != "any" && secondType != "any" {
		// The values have different types.
		return v.RankValues(firstType, secondType)
	}

	// We now know that the types of the values are the same, and neither of
	// the values is nil.
	switch first.Kind() {

	// Handle all primitive types.
	case ref.Bool,
		ref.Uint8, ref.Uint16, ref.Uint32, ref.Uint64, ref.Uint,
		ref.Int8, ref.Int16, ref.Int32, ref.Int64, ref.Int,
		ref.Float32, ref.Float64, ref.Complex64, ref.Complex128,
		ref.String:
		return v.rankPrimitives(first, second)

	// Handle all primitive collection types.
	case ref.Array, ref.Slice:
		switch {
		case first.IsNil():
			if second.IsNil() {
				return 0
			}
			// Only the first value is nil.
			return -1
		case second.IsNil():
			return 1 // We know that first isn't nil.
		default:
			return v.rankArrays(first, second)
		}
	case ref.Map:
		switch {
		case first.IsNil():
			if second.IsNil() {
				return 0
			}
			// Only the first value is nil.
			return -1
		case second.IsNil():
			return 1 // We know that first isn't nil.
		default:
			return v.rankMaps(first, second)
		}

	// Handle all interfaces and pointers.
	case ref.Interface, ref.Pointer:
		switch {
		case first.IsNil():
			if second.IsNil() {
				return 0
			}
			// Only the first value is nil.
			return -1
		case second.IsNil():
			return 1 // We know that first isn't nil.
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
		if ranking != 0 {
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
