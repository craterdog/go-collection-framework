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
// constants, constructors and functions for the Set class namespace.
type setClass_[V Value] struct {
	// This class defines no constants.
}

// This private constant defines a map to hold all the singleton references to
// the type specific Set class namespaces.
var setClassSingletons = map[string]any{}

// This public function returns the singleton reference to a type specific
// Set class namespace.  It also initializes any class constants as needed.
func Set[V Value]() *setClass_[V] {
	var class *setClass_[V]
	var key = fmt.Sprintf("%T", class) // The name of the bound class type.
	var value = setClassSingletons[key]
	switch actual := value.(type) {
	case *setClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Create a new bound class type.
		class = &setClass_[V]{
			// This class defines no constants.
		}
		setClassSingletons[key] = class
	}
	return class
}

// CLASS CONSTRUCTORS

// This public class constructor creates a new empty Set.
// The Set uses the default ranking function to order its values.
func (c *setClass_[V]) Empty() SetLike[V] {
	var set = c.WithRanker(Collator().RankValues)
	return set
}

// This public class constructor creates a new Set from the specified sequence.
// The Set uses the default ranking function to order its values.
func (c *setClass_[V]) FromSequence(sequence Sequential[V]) SetLike[V] {
	var set = c.FromSequenceWithRanker(sequence, Collator().RankValues)
	return set
}

// This public class constructor creates a new Set from the specified sequence.
// The Set uses the specified ranking function to order its values.
func (c *setClass_[V]) FromSequenceWithRanker(
	sequence Sequential[V],
	ranker RankingFunction,
) SetLike[V] {
	var set = c.WithRanker(ranker)
	var iterator = sequence.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		set.AddValue(value)
	}
	return set
}

// This public class constructor creates a new empty Set.
// The Set uses the specified ranking function to order its values.
func (c *setClass_[V]) WithRanker(ranker RankingFunction) SetLike[V] {
	var values = List[V]().Empty()
	var set = &set_[V]{
		ranker,
		values,
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

// This public class function returns the logical disjunction of the specified
// sets.
func (c *setClass_[V]) Or(first, second SetLike[V]) SetLike[V] {
	var result = c.FromSequenceWithRanker(first, first.GetRanker())
	result.AddValues(second)
	return result
}

// This public class function returns the logical material non-implication of
// the specified sets.
func (c *setClass_[V]) Sans(first, second SetLike[V]) SetLike[V] {
	var result = c.FromSequenceWithRanker(first, first.GetRanker())
	result.RemoveValues(second)
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
// implement the set-like abstract type.  This type maintains a Set of values.
// The Set uses ORDINAL based indexing rather than ZERO based indexing (see
// the description of what this means in the Sequential interface definition).
// This type is parameterized as follows:
//   - V is any type of value.
type set_[V Value] struct {
	rank   RankingFunction
	values ListLike[V]
}

// Accessible Interface

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

// Flexible Interface

// This public class method adds the specified value to this Set if it is not
// already a member of the Set.
func (v *set_[V]) AddValue(value V) {
	var slot, found = v.indexOf(value)
	if !found {
		// The value is not already a member, so add it.
		v.values.InsertValue(slot, value)
	}
}

// This public class method adds the specified values to this Set if they are
// not already members of the Set.
func (v *set_[V]) AddValues(values Sequential[V]) {
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		v.AddValue(value)
	}
}

// This public class method returns the ranker function for this Set.
func (v *set_[V]) GetRanker() RankingFunction {
	return v.rank
}

// This public class method removes all values from this Set.
func (v *set_[V]) RemoveAll() {
	v.values.RemoveAll()
}

// This public class method removes the specified value from this Set. It
// returns true if the value was in the Set and false otherwise.
func (v *set_[V]) RemoveValue(value V) {
	var index, found = v.indexOf(value)
	if found {
		// The value is a member, so remove it.
		v.values.RemoveValue(index)
	}
}

// This public class method removes the specified values from this Set. It
// returns the number of values that were removed.
func (v *set_[V]) RemoveValues(values Sequential[V]) {
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		v.RemoveValue(value)
	}
}

// Searchable Interface

// This public class method determines whether or not this Set contains ALL of
// the specified values.
func (v *set_[V]) ContainsAll(values Sequential[V]) bool {
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		if !v.ContainsValue(value) {
			// This Set is missing at least one of the values.
			return false
		}
	}
	// This Set does contains all of the values.
	return true
}

// This public class method determines whether or not this Set contains ANY of
// the specified values.
func (v *set_[V]) ContainsAny(values Sequential[V]) bool {
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		if v.ContainsValue(value) {
			// This Set contains at least one of the values.
			return true
		}
	}
	// This Set does not contain any of the values.
	return false
}

// This public class method determines whether or not this Set contains the
// specified value.
func (v *set_[V]) ContainsValue(value V) bool {
	var _, found = v.indexOf(value)
	return found
}

// This public class method returns the comparing function for this Set.
func (v *set_[V]) GetComparer() ComparingFunction {
	return v.values.GetComparer()
}

// This public class method returns the index of the FIRST occurrence of the
// specified value in this set, or zero if this Set does not contain the value.
func (v *set_[V]) GetIndex(value V) int {
	var index, found = v.indexOf(value)
	if !found {
		return 0
	}
	return index
}

// Sequential Interface

// This public class method returns all the values in this array. The values
// retrieved are in the same order as they are in the array.
func (v *set_[V]) AsArray() []V {
	return v.values.AsArray()
}

// This public class method generates for this Set an iterator that can be
// used to traverse its values.
func (v *set_[V]) GetIterator() Ratcheted[V] {
	var iterator = v.values.GetIterator()
	return iterator
}

// This public class method returns the number of values contained in this
// array.
func (v *set_[V]) GetSize() int {
	return v.values.GetSize()
}

// This public class method determines whether or not this array is empty.
func (v *set_[V]) IsEmpty() bool {
	return v.values.IsEmpty()
}

// Private Interface

// This private class method performs a binary search of the Set for the
// specified value. It returns two results:
//   - index: The index of the value, or if not found, the slot in which it
//     could be inserted in the underlying list.
//   - found: A boolean stating whether or not the value was found.
//
// The algorithm performs a true O[log(n)] worst case search.
func (v *set_[V]) indexOf(value V) (index int, found bool) {
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

// This public class method is used by Go to generate a canonical string for
// the Set.
func (v *set_[V]) String() string {
	return Formatter().FormatCollection(v)
}
