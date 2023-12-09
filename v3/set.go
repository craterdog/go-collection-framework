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

// SET IMPLEMENTATION

// This constructor creates a new empty set that uses the natural rank
// function.
func Set[V Value]() SetLike[V] {
	var rank = RankValues
	var List = List[V]()
	var values = List.FromNothing()
	return &set[V]{values, rank}
}

// This constructor creates a new empty set that uses the specified rank
// function.
func SetWithRanker[V Value](rank RankingFunction) SetLike[V] {
	var List = List[V]()
	var values = List.FromNothing()
	return &set[V]{values, rank}
}

// This constructor creates a new set from the specified array. The set uses the
// natural rank function.
func SetFromArray[V Value](array []V) SetLike[V] {
	var v = Set[V]()
	for _, value := range array {
		v.AddValue(value)
	}
	return v
}

// This constructor creates a new set from the specified sequence. The set uses
// the natural rank function.
func SetFromSequence[V Value](sequence Sequential[V]) SetLike[V] {
	var v = Set[V]()
	var iterator = sequence.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		v.AddValue(value)
	}
	return v
}

// This function returns the logical inverse of the specified set.
func Not[V Value](set SetLike[V]) SetLike[V] {
	panic("Not(set) is meaningless, use Sans(fullSet, set) instead.")
}

// This function returns the logical conjunction of the specified sets.
func And[V Value](first, second SetLike[V]) SetLike[V] {
	var result = Set[V]()
	var iterator = first.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		if second.ContainsValue(value) {
			result.AddValue(value)
		}
	}
	return result
}

// This function returns the logical material non-implication of the
// specified sets.
func Sans[V Value](first, second SetLike[V]) SetLike[V] {
	var result = Set[V]()
	result.AddValues(first)
	result.RemoveValues(second)
	return result
}

// This function returns the logical disjunction of the specified sets.
func Or[V Value](first, second SetLike[V]) SetLike[V] {
	var result = Set[V]()
	result.AddValues(first)
	result.AddValues(second)
	return result
}

// This function returns the logical exclusive disjunction of the
// specified sets.
func Xor[V Value](first, second SetLike[V]) SetLike[V] {
	var result = Or(Sans(first, second), Sans(second, first))
	return result
}

// This type defines the structure and methods associated with a set of values.
// The set uses ORDINAL based indexing rather than ZERO based indexing (see
// the description of what this means in the Sequential interface definition).
// This type is parameterized as follows:
//   - V is any type of value.
type set[V Value] struct {
	values ListLike[V]
	rank   RankingFunction
}

// Sequential Interface

// This public class method determines whether or not this array is empty.
func (v *set[V]) IsEmpty() bool {
	return v.values.IsEmpty()
}

// This public class method returns the number of values contained in this
// array.
func (v *set[V]) GetSize() int {
	return v.values.GetSize()
}

// This public class method returns all the values in this array. The values
// retrieved are in the same order as they are in the array.
func (v *set[V]) AsArray() []V {
	return v.values.AsArray()
}

// Accessible Interface

// This public class method generates for this set an iterator that can be
// used to traverse its values.
func (v *set[V]) GetIterator() Ratcheted[V] {
	var iterator = v.values.GetIterator()
	return iterator
}

// This public class method retrieves from this array the value that is
// associated with the specified index.
func (v *set[V]) GetValue(index int) V {
	return v.values.GetValue(index)
}

// This public class method retrieves from this array all values from the first
// index through the last index (inclusive).
func (v *set[V]) GetValues(first int, last int) Sequential[V] {
	return v.values.GetValues(first, last)
}

// SEARCHABLE INTERFACE

// This method returns the index of the FIRST occurrence of the specified value in
// this list, or zero if this list does not contain the value.
func (v *set[V]) GetIndex(value V) int {
	var index, found = v.search(value)
	if !found {
		return 0
	}
	return index
}

// This method determines whether or not this set contains the specified value.
func (v *set[V]) ContainsValue(value V) bool {
	var _, found = v.search(value)
	return found
}

// This method determines whether or not this set contains ANY of the
// specified values.
func (v *set[V]) ContainsAny(values Sequential[V]) bool {
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		if v.ContainsValue(value) {
			// This set contains at least one of the values.
			return true
		}
	}
	// This set does not contain any of the values.
	return false
}

// This method determines whether or not this set contains ALL of the
// specified values.
func (v *set[V]) ContainsAll(values Sequential[V]) bool {
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		if !v.ContainsValue(value) {
			// This set is missing at least one of the values.
			return false
		}
	}
	// This set does contains all of the values.
	return true
}

// FLEXIBLE INTERFACE

// This method adds the specified value to this set if it is not already a
// member of the set.
func (v *set[V]) AddValue(value V) {
	var slot, found = v.search(value)
	if !found {
		// The value is not already a member, so add it.
		v.values.InsertValue(slot, value)
	}
}

// This method adds the specified values to this set if they are not already
// members of the set.
func (v *set[V]) AddValues(values Sequential[V]) {
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		v.AddValue(value)
	}
}

// This method removes the specified value from this set. It returns true if the
// value was in the set and false otherwise.
func (v *set[V]) RemoveValue(value V) {
	var index, found = v.search(value)
	if found {
		// The value is a member, so remove it.
		v.values.RemoveValue(index)
	}
}

// This method removes the specified values from this set. It returns the number
// of values that were removed.
func (v *set[V]) RemoveValues(values Sequential[V]) {
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		v.RemoveValue(value)
	}
}

// This method removes all values from this set.
func (v *set[V]) RemoveAll() {
	v.values.RemoveAll()
}

// GO INTERFACE

func (v *set[V]) String() string {
	return FormatCollection(v)
}

// PRIVATE INTERFACE

// This private method performs a binary search of the set for the specified
// value. It returns two results:
//   - index: The index of the value, or if not found, the index of the value
//     before which it could be inserted in the underlying list.
//   - found: A boolean stating whether or not the value was found.
//
// The algorithm performs a true O[log(n)] worst case search.
func (v *set[V]) search(value V) (index int, found bool) {
	// We use iteration instead of recursion for better performance.
	//    start        first      middle       last          end
	//    |-------------||----------||----------||-------------|
	//                  |<-- size -------------->|
	//
	var first = 1          // Start at the beginning.
	var last = v.GetSize() // End at the end.
	var size = last        // Initially all values are candidates.
	for size > 0 {
		var middle = first + size/2 // Rounds down to the nearest integer.
		var candidate = v.GetValue(middle)
		switch v.rank(value, candidate) {
		case -1:
			// The index of the value is less than the middle
			// index so the first index stays the same.
			last = middle - 1 // We already tried the middle index.
			size = middle - first
		case 0:
			// The index of the value is the middle index.
			return middle, true
		case 1:
			// The index of the value is greater than the middle
			// index so the last index stays the same.
			first = middle + 1 // We already tried the middle index.
			size = last - middle
		}
	}
	// The value was not found, the last index represents the SLOT where it
	// would be inserted. Note that since the value was not found, the
	// indexes are inverted: last < first (i.e. last = first - 1).
	return last, false
}
