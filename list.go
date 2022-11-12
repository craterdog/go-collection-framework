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

// LIST IMPLEMENTATION

// This constructor creates a new empty list that uses the canonical compare
// function.
func List[V Value]() ListLike[V] {
	var capacity = 4 // The minimum value.
	var values = make([]V, 0, capacity)
	var compare = CompareValues
	return &list[V]{values, compare}
}

// This constructor creates a new list from the specified array that uses the
// canonical compare function.
func ListFromArray[V Value](array []V) ListLike[V] {
	var v = List[V]()
	for _, value := range array {
		v.AddValue(value)
	}
	return v
}

// This function returns the concatenation of the two specified lists.
func Concatenate[V Value](first, second ListLike[V]) ListLike[V] {
	var result = List[V]()
	result.AddValues(first)
	result.AddValues(second)
	return result
}

// This type defines the structure and methods associated with a list of values.
// Each value is associated with an implicit positive integer index. The list
// uses ORDINAL based indexing rather than ZERO based indexing (see the
// description of what this means in the Sequential interface definition).
// This type is parameterized as follows:
//   - V is any type of value.
type list[V Value] struct {
	Array[V]
	compare ComparisonFunction
}

// STRINGER INTERFACE

func (v *list[V]) String() string {
	return FormatValue(v)
}

// INDEXED INTERFACE

// This method sets the comparer function for this list.
func (v *list[V]) SetComparer(compare ComparisonFunction) {
	if compare == nil {
		compare = CompareValues
	}
	v.compare = compare
}

// MALLEABLE INTERFACE

// This method appends the specified value to the end of this list.
func (v *list[V]) AddValue(value V) {
	// Add space for the new value.
	var index = len(v.Array)
	var length = index + 1
	v.resize(length)

	// Append the new value.
	v.Array[index] = value
}

// This method appends the specified values to the end of this list.
func (v *list[V]) AddValues(values Sequential[V]) {
	// Add space for the new values.
	var index = len(v.Array)
	var length = index + values.GetSize()
	v.resize(length)

	// Append the new values.
	copy(v.Array[index:], values.AsArray())
}

// This method inserts the specified value into this list in the specified
// slot between existing values.
func (v *list[V]) InsertValue(slot int, value V) {
	// Add space for the new value.
	var length = len(v.Array) + 1
	v.resize(length)

	// Insert the new value.
	copy(v.Array[slot+1:], v.Array[slot:])
	v.Array[slot] = value
}

// This method inserts the specified values into this list in the specified
// slot between existing values.
func (v *list[V]) InsertValues(slot int, values Sequential[V]) {
	// Add space for the new values.
	var size = values.GetSize()
	var length = len(v.Array) + size
	v.resize(length)

	// Insert the new values.
	copy(v.Array[slot+size:], v.Array[slot:])
	copy(v.Array[slot:], values.AsArray())
}

// This method removes the value at the specified index from this list. The
// removed value is returned.
func (v *list[V]) RemoveValue(index int) V {
	// Remove the old value.
	index = v.normalizedIndex(index)
	var old = v.Array[index]
	copy(v.Array[index:], v.Array[index+1:])

	// Remove extra space.
	var length = len(v.Array) - 1
	v.resize(length)
	return old
}

// This method removes the values in the specified index range from this list.
// The removed values are returned.
func (v *list[V]) RemoveValues(first int, last int) Sequential[V] {
	// Remove the specified values.
	first = v.normalizedIndex(first)
	last = v.normalizedIndex(last)
	var result = ListFromArray[V](v.Array[first : last+1])
	copy(v.Array[first:], v.Array[last+1:])

	// Remove the extra space.
	var length = len(v.Array)
	var size = last - first + 1
	v.resize(length - size)
	return result
}

// This method removes all values from this list.
func (v *list[V]) RemoveAll() {
	v.Array = make([]V, 0, 4)
}

// SORTABLE INTERFACE

// This method sorts the values in this list using the canonical rank function.
func (v *list[V]) SortValues() {
	v.SortValuesWithRanker(nil)
}

// This method sorts the values in this list using the specified rank function.
func (v *list[V]) SortValuesWithRanker(rank RankingFunction) {
	if rank == nil {
		rank = RankValues
	}
	if len(v.Array) > 1 {
		SortArray[V](v.Array, rank)
	}
}

// This method reverses the order of all values in this list.
func (v *list[V]) ReverseValues() {
	// Allocate a new array/slice.
	var length = len(v.Array)
	var capacity = cap(v.Array)
	var reversed = make([]V, length, capacity)

	// Copy the values into the new array in reverse.
	for i, _ := range v.Array {
		reversed[i] = v.Array[length-i-1]
	}
	v.Array = reversed
}

// This method pseudo-randomly shuffles the values in this list.
func (v *list[V]) ShuffleValues() {
	ShuffleArray[V](v.Array)
}

// PRIVATE INTERFACE

// This method resizes this list and adjusts the capacity of this list if
// necessary to make it more efficient. Note: Any additional values that are
// added to the length of the list are NOT zeroed out.
func (v *list[V]) resize(length int) {
	var capacity = cap(v.Array)
	for length > capacity {
		capacity *= 2
	}
	for length < capacity/4 {
		capacity /= 2
	}
	if capacity != cap(v.Array) {
		// Adjust the capacity accordingly.
		var values = make([]V, length, capacity)
		copy(values, v.Array)
		v.Array = values
	}
	v.Array = v.Array[:length] // A change the length of the slice.
}
