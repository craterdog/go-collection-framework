/*******************************************************************************
 *   Copyright (c) 2009-2023 Crater Dog Technologies™.  All Rights Reserved.   *
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

// This private type defines the namespace structure associated with the
// constants, constructors and functions for the Collator class namespace.
type collatorClass_ struct {
	defaultDepth int
}

// This private constant defines the singleton reference to the Collator
// class namespace.  It also initializes any class constants as needed.
var collatorClassSingleton = &collatorClass_{
	defaultDepth: 16,
}

// This public function returns the singleton reference to the Collator
// class namespace.
func Collator() *collatorClass_ {
	return collatorClassSingleton
}

// CLASS CONSTANTS

// This public class constant represents the default depth of a collection
// at which the Collator gives up trying to collate.  This handles cycles
// in a sensible and efficient manner.
func (c *collatorClass_) DefaultDepth() int {
	return c.defaultDepth
}

// CLASS CONSTRUCTORS

// This public class constructor creates a new Collator with the default
// maximum traversal depth.
func (c *collatorClass_) WithDefaultDepth() CollatorLike {
	var collator = &collator_{
		maximumDepth: c.defaultDepth,
	}
	return collator
}

// This public class constructor creates a new Collator with the specified
// maximum traversal depth.
func (c *collatorClass_) WithSpecifiedDepth(depth int) CollatorLike {
	if depth < 0 || depth > c.defaultDepth {
		depth = c.defaultDepth
	}
	var collator = &collator_{
		maximumDepth: depth,
	}
	return collator
}

// CLASS FUNCTIONS

// This public class function determines whether or not the specified values are
// equal using their natural comparison criteria.
func (c *collatorClass_) CompareValues(first Value, second Value) bool {
	var v = c.WithDefaultDepth()
	return v.CompareValues(first, second)
}

// This public class function returns the ranking for the specified values based
// on their natural ordering:
//   - -1: first < second
//   - 0: first = second
//   - 1: first > second
func (c *collatorClass_) RankValues(first Value, second Value) int {
	var v = c.WithDefaultDepth()
	return v.RankValues(first, second)
}

// CLASS TYPE

// Encapsulated Type

// This private class type encapsulates a Go structure containing private
// attributes that can only be accessed and manipulated using methods that
// implement the collator-like abstract type.
type collator_ struct {
	depth        int
	maximumDepth int
}

// Discerning Interface

// This public class method determines whether or not the specified values are equal.
func (v *collator_) CompareValues(first Value, second Value) bool {
	return v.compareValues(ref.ValueOf(first), ref.ValueOf(second))
}

// This public class method ranks the specified values using their natural ordering.
func (v *collator_) RankValues(first Value, second Value) int {
	return v.rankValues(ref.ValueOf(first), ref.ValueOf(second))
}

// Private Interface

// This private class method removes the generics from the type string for the
// specified type and converts an empty interface into type "any".
func (v *collator_) baseTypeName(t ref.Type) string {
	var result = t.String()
	var index = sts.Index(result, "[")
	if index > -1 {
		result = result[:index]
	}
	if result == "interface {}" {
		result = "any"
	}
	return result
}

// This private class method determines whether or not the specified arrays have
// the same value.
func (v *collator_) compareArrays(first ref.Value, second ref.Value) bool {
	// Check for maximum recursion depth.
	if v.depth+1 > v.maximumDepth {
		panic(fmt.Sprintf("The maximum recursion depth was exceeded: %v", v.depth))
	}

	// Compare the sizes of the arrays.
	var size = first.Len()
	if second.Len() != size {
		// The arrays are of different lengths.
		return false
	}

	// Compare the values of the arrays.
	for i := 0; i < size; i++ {
		v.depth++
		if !v.compareValues(first.Index(i), second.Index(i)) {
			// Two of the values in the arrays are different.
			v.depth--
			return false
		}
		v.depth--
	}
	return true
}

// This private class method determines whether or not the specified elements
// have the same value.
func (v *collator_) compareElements(first, second ref.Value) bool {
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

// This private class method determines whether or not the specified maps have
// the same value.
func (v *collator_) compareMaps(first ref.Value, second ref.Value) bool {
	// Check for maximum recursion depth.
	if v.depth+1 > v.maximumDepth {
		panic(fmt.Sprintf("The maximum recursion depth was exceeded: %v", v.depth))
	}

	// Compare the sizes of the two maps.
	if first.Len() != second.Len() {
		// The maps are different sizes.
		return false
	}

	// Compare the keys and values for the two maps.
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
	// Compare the arrays for the two sequences.
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
	var firstType = v.baseTypeName(first.Type())
	var secondType = v.baseTypeName(second.Type())
	if firstType != secondType && firstType != "any" && secondType != "any" {
		// The values have different types.
		return false
	}

	// We now know that the types of the values are the same, and neither of
	// the values is invalid.
	switch first.Kind() {

	// Handle all native elemental types.
	case ref.Bool,
		ref.Uint8, ref.Uint16, ref.Uint32, ref.Uint64, ref.Uint,
		ref.Int8, ref.Int16, ref.Int32, ref.Int64, ref.Int,
		ref.Float32, ref.Float64, ref.Complex64, ref.Complex128,
		ref.String:
		return v.compareElements(first, second)

	// Handle all native collection types.
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

// This private class method returns the ranking order of the specified arrays
// using a recursive descent algorithm.
func (v *collator_) rankArrays(first ref.Value, second ref.Value) int {
	// Check for maximum recursion depth.
	if v.depth+1 > v.maximumDepth {
		panic(fmt.Sprintf("The maximum recursion depth was exceeded: %v", v.depth))
	}

	// Determine the smallest array.
	var firstSize = first.Len()
	var secondSize = second.Len()
	if firstSize > secondSize {
		// Swap the order of the arrays and reverse the sign of the result.
		return -1 * v.rankArrays(second, first)
	}

	// Iterate through the smallest array.
	for i := 0; i < firstSize; i++ {
		v.depth++
		var rank = v.rankValues(first.Index(i), second.Index(i))
		if rank < 0 {
			// The value in the first array comes before its matching value.
			v.depth--
			return -1
		}
		if rank > 0 {
			// The value in the first array comes after its matching value.
			v.depth--
			return 1
		}
		// The two values match.
		v.depth--
	}

	// The arrays contain the same initial values.
	if secondSize > firstSize {
		// The shorter array is ranked before the longer array.
		return -1
	}
	// The arrays are the same length and contain the same values.
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
// elements.  NOTE: Go does not provide an easy way to a apply a possible tilde
// type (e.g. ~string) to its elemental (named) type without knowing whether or
// not it is actually a tilde type. So we must convert it to a string and then
// parse it back again...
// This method attempts to hide that ugliness from the rest of the code.
func (v *collator_) rankElements(first, second ref.Value) int {
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
	case ref.Int8, ref.Int16, ref.Int32, ref.Int64, ref.Int:
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

// This private class method returns the ranking order of the specified maps
// using a recursive descent algorithm. NOTE: currently the implementation of Go
// maps is hashtable based. The order of the keys is random, even for two maps
// with the same keys if the associations were entered in different sequences.
// Therefore at this time it is necessary to sort the key arrays for each map.
// This introduces a circular dependency between the implementation of the
// Collator and the sorter (i.e. rankMaps() -> SortValues() -> RankingFunction
// type).
func (v *collator_) rankMaps(first ref.Value, second ref.Value) int {
	// Check for maximum recursion depth.
	if v.depth+1 > v.maximumDepth {
		panic(fmt.Sprintf("The maximum recursion depth was exceeded: %v", v.depth))
	}

	// Extract and sort the keys for the two maps.
	var Sorter = Sorter[ref.Value]()
	var firstKeys = first.MapKeys() // The returned keys are in random order.
	Sorter.SortValues(firstKeys, v.rankReflective)
	var secondKeys = second.MapKeys() // The returned keys are in random order.
	Sorter.SortValues(secondKeys, v.rankReflective)

	// Determine the smallest map.
	var firstSize = len(firstKeys)
	var secondSize = len(secondKeys)
	if firstSize > secondSize {
		// Swap the order of the maps and reverse the sign of the result.
		return -1 * v.rankMaps(second, first)
	}

	// Iterate through the smallest map.
	for i := 0; i < firstSize; i++ {
		v.depth++

		// Rank the two keys.
		var firstKey = firstKeys[i]
		var secondKey = secondKeys[i]
		var keyRank = v.rankValues(firstKey, secondKey)
		if keyRank < 0 {
			// The key in the first map comes before its matching key.
			v.depth--
			return -1
		}
		if keyRank > 0 {
			// The key in the first map comes after its matching key.
			v.depth--
			return 1
		}

		// The two keys match so rank the corresponding values.
		var firstValue = first.MapIndex(firstKey)
		var secondValue = second.MapIndex(secondKey)
		var valueRank = v.rankValues(firstValue, secondValue)
		if valueRank < 0 {
			// The value in the first map comes before its matching value.
			v.depth--
			return -1
		}
		if valueRank > 0 {
			// The value in the first map comes after its matching value.
			v.depth--
			return 1
		}
		v.depth--
	}

	// The maps contain the same initial associations.
	if secondSize > firstSize {
		// The shorter map is ranked before the longer map.
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

// This private class method returns the ranking order of the specified
// sequences using a recursive descent algorithm.
func (v *collator_) rankSequences(first ref.Value, second ref.Value) int {
	// Rank the arrays for the two sequences.
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
	var firstType = v.baseTypeName(first.Type())
	var secondType = v.baseTypeName(second.Type())
	if firstType != secondType && firstType != "any" && secondType != "any" {
		// The values have different types.
		return v.RankValues(firstType, secondType)
	}

	// We now know that the types of the values are the same, and neither of
	// the values is nil.
	switch first.Kind() {

	// Handle all native elemental types.
	case ref.Bool,
		ref.Uint8, ref.Uint16, ref.Uint32, ref.Uint64, ref.Uint,
		ref.Int8, ref.Int16, ref.Int32, ref.Int64, ref.Int,
		ref.Float32, ref.Float64, ref.Complex64, ref.Complex128,
		ref.String:
		return v.rankElements(first, second)

	// Handle all native collection types.
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
