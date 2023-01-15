/*******************************************************************************
 *   Copyright (c) 2009-2022 Crater Dog Technologies™.  All Rights Reserved.   *
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

// COLLATOR INTERFACE

// This function determines whether or not the specified values are equal using
// their natural comparison criteria.
func CompareValues(first Value, second Value) bool {
	var v = Collator()
	return v.CompareValues(first, second)
}

// This function ranks the specified values based on their natural ordering.
func RankValues(first Value, second Value) int {
	var v = Collator()
	return v.RankValues(first, second)
}

// This constructor creates a new instance of a collator that can be used to
// compare or rank any two values.
func Collator() CollatorLike {
	return &collator{}
}

// COLLATOR IMPLEMENTATION

// The maximum recursion level handles overly complex data models and cyclical
// references in a clean and efficient way.
const maximumRecursion int = 100

// This type defines the structure and methods for a natural collator agent.
type collator struct {
	depth int // The depth starts out at zero.
}

// DISCERNING INTERFACE

// This method determines whether or not the specified values are equal.
func (v *collator) CompareValues(first Value, second Value) bool {
	return v.compareValues(ref.ValueOf(first), ref.ValueOf(second))
}

// This method ranks the specified values using their natural ordering.
func (v *collator) RankValues(first Value, second Value) int {
	return v.rankValues(ref.ValueOf(first), ref.ValueOf(second))
}

// PRIVATE INTERFACE

// This function determines whether or not the specified reflective values are
// equal using reflection and a recursive descent algorithm.
func (v *collator) compareValues(first ref.Value, second ref.Value) bool {
	// Handle any nil pointers.
	if !first.IsValid() {
		return !second.IsValid()
	}
	if !second.IsValid() {
		return false
	}

	// At this point, neither of the values are nil.
	var firstType = baseTypeName(first.Type())
	var secondType = baseTypeName(second.Type())
	if firstType != secondType && firstType != "any" && secondType != "any" {
		// The values have different types.
		return false
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
		return v.compareElements(first, second)

	// Handle all native collection types.
	case ref.Array, ref.Slice:
		return v.compareArrays(first, second)
	case ref.Map:
		return v.compareMaps(first, second)

	// Handle all interfaces and pointers.
	case ref.Interface, ref.Pointer:
		switch {
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
			"Attempted to compare:\n\tfirst: %v\n\ttype: %v\n\tkind: %v\nand\n\tsecond: %v\n\ttype: %v\n\tkind: %v\n",
			first.Interface(),
			first.Type(),
			first.Kind(),
			second.Interface(),
			second.Type(),
			second.Kind()))
	}
}

// This private method determines whether or not the specified elements have
// the same value.
func (v *collator) compareElements(first, second ref.Value) bool {
	return first.Interface() == second.Interface()
}

// This private method determines whether or not the specified arrays have
// the same value.
func (v *collator) compareArrays(first ref.Value, second ref.Value) bool {
	var size = first.Len()
	if second.Len() != size {
		// The arrays are of different lengths.
		return false
	}
	for i := 0; i < size; i++ {
		v.depth++
		if v.depth > maximumRecursion {
			panic(fmt.Sprintf("The maximum recursion depth was exceeded: %v", maximumRecursion))
		}
		if !v.compareValues(first.Index(i), second.Index(i)) {
			// Two of the values in the arrays are different.
			v.depth--
			return false
		}
		v.depth--
	}
	// All values in the arrays are equal.
	return true
}

// This private method determines whether or not the specified maps have
// the same value.
func (v *collator) compareMaps(first ref.Value, second ref.Value) bool {
	// Compare the sizes of the two maps.
	if first.Len() != second.Len() {
		// The maps are different sizes.
		return false
	}

	// Compare the keys and values for the two maps.
	var iterator = first.MapRange()
	for iterator.Next() {
		v.depth++
		if v.depth > maximumRecursion {
			panic(fmt.Sprintf("The maximum recursion depth was exceeded: %v", maximumRecursion))
		}
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

	// All keys and values match.
	return true
}

// This private method determines whether or not the specified sequences
// have the same value.
func (v *collator) compareSequences(first ref.Value, second ref.Value) bool {
	// Compare the arrays for the two sequences.
	var firstArray = first.MethodByName("AsArray").Call([]ref.Value{})[0]
	var secondArray = second.MethodByName("AsArray").Call([]ref.Value{})[0]
	return v.compareArrays(firstArray, secondArray)
}


// This private method determines whether or not the specified interfaces
// have the same getter values.
func (v *collator) compareInterfaces(first ref.Value, second ref.Value) bool {
	var typeRef = first.Type() // We know the structures are the same type.
	var count = first.NumMethod()
	for index := 0; index < count; index++ {
		var method = typeRef.Method(index)
		if sts.HasPrefix(method.Name, "Get") {
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

// PRIVATE FUNCTIONS

func (v *collator) rankReflective(first Value, second Value) int {
	var firstValue = first.(ref.Value)
	var secondValue = second.(ref.Value)
	return v.rankValues(firstValue, secondValue)
}

// This private method returns the ranking order of the specified values using
// reflection and a recursive descent algorithm.
func (v *collator) rankValues(first ref.Value, second ref.Value) int {
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
	var firstType = baseTypeName(first.Type())
	var secondType = baseTypeName(second.Type())
	if firstType != secondType && firstType != "any" && secondType != "any" {
		// The values have different types.
		return RankValues(firstType, secondType)
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
		return v.rankArrays(first, second)
	case ref.Map:
		return v.rankMaps(first, second)

	// Handle all interfaces and pointers.
	case ref.Interface, ref.Pointer:
		switch {
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
			"Attempted to rank:\n\tfirst: %v\n\ttype: %v\n\tkind: %v\nand\n\tsecond: %v\n\ttype: %v\n\tkind: %v\n",
			first.Interface(),
			first.Type(),
			first.Kind(),
			second.Interface(),
			second.Type(),
			second.Kind()))
	}
}

// This private method returns the ranking order of the specified elements. Note
// that Go does not provide an easy way to a possible tilda type (e.g. ~string)
// to its elemental (named) type without knowing whether or not it is actually a
// tilda type. So we must convert it to a string and then parse it back again...
// This method attempts to hide that ugliness from the rest of the code.
func (v *collator) rankElements(first, second ref.Value) int {
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

// This private function returns the ranking order of the specified boolean
// values.
func (v *collator) rankBooleans(first, second bool) int {
	if !first && second {
		return -1
	}
	if first && !second {
		return 1
	}
	return 0
}

// This private function returns the ranking order of the specified unsigned
// integers.
func (v *collator) rankUnsigned(first, second uint64) int {
	if first < second {
		return -1
	}
	if first > second {
		return 1
	}
	return 0
}

// This private function returns the ranking order of the specified signed
// integers.
func (v *collator) rankSigned(first, second int64) int {
	if first < second {
		return -1
	}
	if first > second {
		return 1
	}
	return 0
}

// This private function returns the ranking order of the specified floating
// point numbers.
func (v *collator) rankFloats(first, second float64) int {
	if first < second {
		return -1
	}
	if first > second {
		return 1
	}
	return 0
}

// This private function returns the ranking order of the specified complex
// number values.
func (v *collator) rankComplex(first, second complex128) int {
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

// This private function returns the ranking order of the specified string
// values.
func (v *collator) rankStrings(first, second string) int {
	if first < second {
		// The first string comes before the second string lexigraphically.
		return -1
	}
	if first > second {
		// The first string comes after the second string lexigraphically.
		return 1
	}
	// The two strings are the same.
	return 0
}

// This private method returns the ranking order of the specified arrays using
// a recursive descent algorithm.
func (v *collator) rankArrays(first ref.Value, second ref.Value) int {
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
		if v.depth > maximumRecursion {
			panic(fmt.Sprintf("The maximum recursion depth was exceeded: %v", maximumRecursion))
		}
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

// This private method returns the ranking order of the specified maps using a
// recursive descent algorithm. Note: currently the implementation of Go maps is
// hashtable based. The order of the keys is random, even for two maps with the
// same keys if the associations were entered in different sequences. Therefore
// at this time it is necessary to sort the key arrays for each map. This
// introduces a circular dependency between the implementation of the collator
// and the sorter (i.e. rankMaps() -> SortArray() -> RankingFunction type).
func (v *collator) rankMaps(first ref.Value, second ref.Value) int {
	// Extract and sort the keys for the two maps.
	var firstKeys = first.MapKeys() // The returned keys are in random order.
	SortArray(firstKeys, v.rankReflective)
	var secondKeys = second.MapKeys() // The returned keys are in random order.
	SortArray(secondKeys, v.rankReflective)

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
		if v.depth > maximumRecursion {
			panic(fmt.Sprintf("The maximum recursion depth was exceeded: %v", maximumRecursion))
		}
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

// This private method returns the ranking order of the specified sequences
// using a recursive descent algorithm.
func (v *collator) rankSequences(first ref.Value, second ref.Value) int {
	// Rank the arrays for the two sequences.
	var firstArray = first.MethodByName("AsArray").Call([]ref.Value{})[0]
	var secondArray = second.MethodByName("AsArray").Call([]ref.Value{})[0]
	return v.rankArrays(firstArray, secondArray)
}

// This private method returns the ranking order of the specified structures
// by ranking the associated fields.
func (v *collator) rankStructures(first ref.Value, second ref.Value) int {
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

// This private method returns the ranking order of the specified interfaces
// by ranking the results of their getter methods.
func (v *collator) rankInterfaces(first ref.Value, second ref.Value) int {
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

// PRIVATE FUNCTIONS

// This function removes the generics from the type string for the specified
// type and converts an empty interface into type "any".
func baseTypeName(t ref.Type) string {
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
