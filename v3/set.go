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
)

// CLASS NAMESPACE

// This private type defines the namespace structure associated with the
// constants, constructors and functions for the set class namespace.
type setClass_[V Value] struct {
	// This class defines no constants.
}

// This private constant defines a map to hold all the singleton references to
// the type specific set namespaces.
var setClassSingletons = map[string]any{}

// This public function returns the singleton reference to a type specific
// set namespace.  It also initializes any class constants as needed.
func Set[V Value]() *setClass_[V] {
	var class *setClass_[V]
	var key = fmt.Sprintf("%T", class)
	var value = setClassSingletons[key]
	switch actual := value.(type) {
	case *setClass_[V]:
		class = actual
	default:
		class = &setClass_[V]{
			// This class defines no constants.
		}
		setClassSingletons[key] = class
	}
	return class
}

// CLASS CONSTRUCTORS

// This public class constructor creates a new empty set.
// The set uses the natural ranking function to order its values.
func (c *setClass_[V]) Empty() SetLike[V] {
	var set = c.WithRanker(Collator().RankValues)
	return set
}

// This public class constructor creates a new empty set.
// The set uses the specified ranking function to order its values.
func (c *setClass_[V]) WithRanker(ranker RankingFunction) SetLike[V] {
	var List = List[V]()
	var values = List.Empty()
	var set = &set_[V]{values, ranker}
	return set
}

// This public class constructor creates a new set from the specified sequence.
// The set uses the natural ranking function to order its values.
func (c *setClass_[V]) FromSequence(sequence Sequential[V]) SetLike[V] {
	var set = c.FromSequenceWithRanker(sequence, Collator().RankValues)
	return set
}

// This public class constructor creates a new set from the specified sequence.
// The set uses the specified ranking function to order its values.
func (c *setClass_[V]) FromSequenceWithRanker(sequence Sequential[V], ranker RankingFunction) SetLike[V] {
	var set = c.WithRanker(ranker)
	var iterator = sequence.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		set.AddValue(value)
	}
	return set
}

// CLASS FUNCTIONS

// This public class function returns the logical conjunction of the specified
// sets.
func (c *setClass_[V]) And(first, second SetLike[V]) SetLike[V] {
	var result = c.WithRanker(first.GetRanker())
	var iterator = first.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		if second.ContainsValue(value) {
			result.AddValue(value)
		}
	}
	return result
}

// This public class function returns the logical material non-implication of
// the specified sets.
func (c *setClass_[V]) Sans(first, second SetLike[V]) SetLike[V] {
	var result = c.FromSequenceWithRanker(first, first.GetRanker())
	result.RemoveValues(second)
	return result
}

// This public class function returns the logical disjunction of the specified
// sets.
func (c *setClass_[V]) Or(first, second SetLike[V]) SetLike[V] {
	var result = c.FromSequenceWithRanker(first, first.GetRanker())
	result.AddValues(second)
	return result
}

// This public class function returns the logical exclusive disjunction of the
// specified sets.
func (c *setClass_[V]) Xor(first, second SetLike[V]) SetLike[V] {
	var result = c.Or(c.Sans(first, second), c.Sans(second, first))
	return result
}

// CLASS TYPE

// Encapsulated Type

// This private class type encapsulates a Go structure containing private
// attributes that can only be accessed and manipulated using methods that
// implement the set-like abstract type.  This type maintains a set of values.
// The set uses ORDINAL based indexing rather than ZERO based indexing (see
// the description of what this means in the Sequential interface definition).
// This type is parameterized as follows:
//   - V is any type of value.
type set_[V Value] struct {
	values ListLike[V]
	rank   RankingFunction
}

// Sequential Interface

// This public class method determines whether or not this array is empty.
func (v *set_[V]) IsEmpty() bool {
	return v.values.IsEmpty()
}

// This public class method returns the number of values contained in this
// array.
func (v *set_[V]) GetSize() int {
	return v.values.GetSize()
}

// This public class method returns all the values in this array. The values
// retrieved are in the same order as they are in the array.
func (v *set_[V]) AsArray() []V {
	return v.values.AsArray()
}

// Accessible Interface

// This public class method generates for this set an iterator that can be
// used to traverse its values.
func (v *set_[V]) GetIterator() Ratcheted[V] {
	var iterator = v.values.GetIterator()
	return iterator
}

// This public class method retrieves from this array the value that is
// associated with the specified index.
func (v *set_[V]) GetValue(index int) V {
	return v.values.GetValue(index)
}

// This public class method retrieves from this array all values from the first
// index through the last index (inclusive).
func (v *set_[V]) GetValues(first int, last int) Sequential[V] {
	return v.values.GetValues(first, last)
}

// Searchable Interface

// This public class method returns the comparing function for this set.
func (v *set_[V]) GetComparer() ComparingFunction {
	return v.values.GetComparer()
}

// This public class method returns the index of the FIRST occurrence of the
// specified value in this set, or zero if this set does not contain the value.
func (v *set_[V]) GetIndex(value V) int {
	var index, found = v.search(value)
	if !found {
		return 0
	}
	return index
}

// This public class method determines whether or not this set contains the
// specified value.
func (v *set_[V]) ContainsValue(value V) bool {
	var _, found = v.search(value)
	return found
}

// This public class method determines whether or not this set contains ANY of
// the specified values.
func (v *set_[V]) ContainsAny(values Sequential[V]) bool {
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

// This public class method determines whether or not this set contains ALL of
// the specified values.
func (v *set_[V]) ContainsAll(values Sequential[V]) bool {
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

// Flexible Interface

// This public class method returns the ranker function for this set.
func (v *set_[V]) GetRanker() RankingFunction {
	return v.rank
}

// This public class method adds the specified value to this set if it is not
// already a member of the set.
func (v *set_[V]) AddValue(value V) {
	var slot, found = v.search(value)
	if !found {
		// The value is not already a member, so add it.
		v.values.InsertValue(slot, value)
	}
}

// This public class method adds the specified values to this set if they are
// not already members of the set.
func (v *set_[V]) AddValues(values Sequential[V]) {
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		v.AddValue(value)
	}
}

// This public class method removes the specified value from this set. It
// returns true if the value was in the set and false otherwise.
func (v *set_[V]) RemoveValue(value V) {
	var index, found = v.search(value)
	if found {
		// The value is a member, so remove it.
		v.values.RemoveValue(index)
	}
}

// This public class method removes the specified values from this set. It
// returns the number of values that were removed.
func (v *set_[V]) RemoveValues(values Sequential[V]) {
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		v.RemoveValue(value)
	}
}

// This public class method removes all values from this set.
func (v *set_[V]) RemoveAll() {
	v.values.RemoveAll()
}

// Private Interface

// This public class method is used by Go to generate a canonical string for
// the set.
func (v *set_[V]) String() string {
	return Formatter().FormatCollection(v)
}

// This private class method performs a binary search of the set for the
// specified value. It returns two results:
//   - index: The index of the value, or if not found, the index of the value
//     before which it could be inserted in the underlying list.
//   - found: A boolean stating whether or not the value was found.
//
// The algorithm performs a true O[log(n)] worst case search.
func (v *set_[V]) search(value V) (index int, found bool) {
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
	// would be inserted. NOTE: Since the value was not found, the indexes are
	// inverted: last < first (i.e. last = first - 1).
	return last, false
}
