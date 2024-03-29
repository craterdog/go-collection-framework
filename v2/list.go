/*
................................................................................
.    Copyright (c) 2009-2024 Crater Dog Technologies.  All Rights Reserved.    .
................................................................................
.  DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               .
.                                                                              .
.  This code is free software; you can redistribute it and/or modify it under  .
.  the terms of The MIT License (MIT), as published by the Open Source         .
.  Initiative. (See http://opensource.org/licenses/MIT)                        .
................................................................................
*/

package collections

// LIST IMPLEMENTATION

// This constructor creates a new empty list that uses the natural compare
// function.
func List[V Value]() ListLike[V] {
	var capacity = 4 // The minimum value.
	var values = make([]V, 0, capacity)
	var compare = CompareValues
	return &list[V]{Array[V](values), values, compare}
}

// This constructor creates a new empty list that uses the specified compare
// function.
func ListWithComparer[V Value](compare ComparisonFunction) ListLike[V] {
	var capacity = 4 // The minimum value.
	var values = make([]V, 0, capacity)
	return &list[V]{Array[V](values), values, compare}
}

// This constructor creates a new list from the specified array. The list uses
// the natural compare function.
func ListFromArray[V Value](array []V) ListLike[V] {
	var v = List[V]()
	for _, value := range array {
		v.AddValue(value)
	}
	return v
}

// This constructor creates a new list from the specified sequence. The list
// uses the natural compare function.
func ListFromSequence[V Value](sequence Sequential[V]) ListLike[V] {
	var v = List[V]()
	var iterator = Iterator(sequence)
	for iterator.HasNext() {
		var value = iterator.GetNext()
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
	ArrayLike[V]
	values  Array[V]
	compare ComparisonFunction
}

// SEARCHABLE INTERFACE

// This method returns the index of the FIRST occurrence of the specified value in
// this list, or zero if this list does not contain the value.
func (v *list[V]) GetIndex(value V) int {
	for index, candidate := range v.values {
		if v.compare(candidate, value) {
			// Found the value.
			return index + 1 // Convert to an ORDINAL based index.
		}
	}
	// The value was not found.
	return 0
}

// This method determines whether or not this list contains the specified value.
func (v *list[V]) ContainsValue(value V) bool {
	return v.GetIndex(value) > 0
}

// This method determines whether or not this list contains ANY of the specified
// values.
func (v *list[V]) ContainsAny(values Sequential[V]) bool {
	var iterator = Iterator(values)
	for iterator.HasNext() {
		var candidate = iterator.GetNext()
		if v.GetIndex(candidate) > 0 {
			// Found one of the values.
			return true
		}
	}
	// Did not find any of the values.
	return false
}

// This method determines whether or not this list contains ALL of the specified
// values.
func (v *list[V]) ContainsAll(values Sequential[V]) bool {
	var iterator = Iterator(values)
	for iterator.HasNext() {
		var candidate = iterator.GetNext()
		if v.GetIndex(candidate) == 0 {
			// Didn't find one of the values.
			return false
		}
	}
	// Found all of the values.
	return true
}

// MALLEABLE INTERFACE

// This method appends the specified value to the end of this list.
func (v *list[V]) AddValue(value V) {
	// Add space for the new value.
	var index = len(v.values)
	var length = index + 1
	v.resize(length)

	// Append the new value.
	v.values[index] = value
}

// This method appends the specified values to the end of this list.
func (v *list[V]) AddValues(values Sequential[V]) {
	// Add space for the new values.
	var index = len(v.values)
	var length = index + values.GetSize()
	v.resize(length)

	// Append the new values.
	copy(v.values[index:], values.AsArray())
}

// This method inserts the specified value into this list in the specified
// slot between existing values.
func (v *list[V]) InsertValue(slot int, value V) {
	// Add space for the new value.
	var length = len(v.values) + 1
	v.resize(length)

	// Insert the new value.
	copy(v.values[slot+1:], v.values[slot:])
	v.values[slot] = value
}

// This method inserts the specified values into this list in the specified
// slot between existing values.
func (v *list[V]) InsertValues(slot int, values Sequential[V]) {
	// Add space for the new values.
	var size = values.GetSize()
	var length = len(v.values) + size
	v.resize(length)

	// Insert the new values.
	copy(v.values[slot+size:], v.values[slot:])
	copy(v.values[slot:], values.AsArray())
}

// This method removes the value at the specified index from this list. The
// removed value is returned.
func (v *list[V]) RemoveValue(index int) V {
	// Remove the old value.
	index = v.values.GoIndex(index)
	var old = v.values[index]
	copy(v.values[index:], v.values[index+1:])

	// Remove extra space.
	var length = len(v.values) - 1
	v.resize(length)
	return old
}

// This method removes the values in the specified index range from this list.
// The removed values are returned.
func (v *list[V]) RemoveValues(first int, last int) Sequential[V] {
	// Remove the specified values.
	first = v.values.GoIndex(first)
	last = v.values.GoIndex(last)
	var result = ListFromArray(v.values[first : last+1])
	copy(v.values[first:], v.values[last+1:])

	// Remove the extra space.
	var length = len(v.values)
	var size = last - first + 1
	v.resize(length - size)
	return result
}

// This method removes all values from this list.
func (v *list[V]) RemoveAll() {
	v.values = make([]V, 0, 4)
	v.ArrayLike = Array[V](v.values)
}

// SORTABLE INTERFACE

// This method sorts the values in this list using the natural rank function.
func (v *list[V]) SortValues() {
	v.SortValuesWithRanker(nil)
}

// This method sorts the values in this list using the specified rank function.
func (v *list[V]) SortValuesWithRanker(rank RankingFunction) {
	if rank == nil {
		rank = RankValues
	}
	if len(v.values) > 1 {
		SortArray(v.values, rank)
	}
}

// This method reverses the order of all values in this list.
func (v *list[V]) ReverseValues() {
	ReverseArray(v.values)
}

// This method pseudo-randomly shuffles the values in this list.
func (v *list[V]) ShuffleValues() {
	ShuffleArray(v.values)
}

// GO INTERFACE

func (v *list[V]) String() string {
	return FormatCollection(v)
}

// PRIVATE INTERFACE

// This method resizes this list and adjusts the capacity of this list if
// necessary to make it more efficient. Note: Any additional values that are
// added to the length of the list are NOT zeroed out.
func (v *list[V]) resize(length int) {
	var capacity = cap(v.values)
	for length > capacity {
		capacity *= 2
	}
	for length < capacity/4 {
		capacity /= 2
	}
	if capacity != cap(v.values) {
		// Adjust the capacity accordingly.
		var values = make([]V, length, capacity)
		copy(values, v.values)
		v.values = values
	}
	v.values = v.values[:length] // A change the length of the slice.
	v.ArrayLike = Array[V](v.values)
}
