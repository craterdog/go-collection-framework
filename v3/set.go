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
)

// CLASS NAMESPACE

// Private Class Namespace Type

type setClass_[V Value] struct {
	// This class defines no constants.
}

// Private Class Namespace References

var setClass = map[string]any{}

// Public Class Namespace Access

func SetClass[V Value]() SetClassLike[V] {
	var class *setClass_[V]
	var key = fmt.Sprintf("%T", class) // The name of the bound class type.
	var value = setClass[key]
	switch actual := value.(type) {
	case *setClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Create a new bound class type.
		class = &setClass_[V]{
			// This class defines no constants.
		}
		setClass[key] = class
	}
	return class
}

// Public Class Constructors

func (c *setClass_[V]) Empty() SetLike[V] {
	var collator = CollatorClass().Default()
	var ranker = collator.RankValues
	var set = c.WithRanker(ranker)
	return set
}

func (c *setClass_[V]) FromArray(values []V) SetLike[V] {
	var array = ArrayClass[V]().FromArray(values)
	var set = c.FromSequence(array)
	return set
}

func (c *setClass_[V]) FromSequence(values Sequential[V]) SetLike[V] {
	var collator = CollatorClass().Default()
	var ranker = collator.RankValues
	var set = c.FromSequenceWithRanker(values, ranker)
	return set
}

func (c *setClass_[V]) FromSequenceWithRanker(
	values Sequential[V],
	ranker RankingFunction,
) SetLike[V] {
	var set = c.WithRanker(ranker)
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		set.AddValue(value)
	}
	return set
}

func (c *setClass_[V]) FromString(values string) SetLike[V] {
	// First we parse it as a collection of any type value.
	var cdcn = CDCNClass().Default()
	var collection = cdcn.ParseCollection(values).(Sequential[Value])

	// Then we convert it to a set of type V.
	var set = c.Empty()
	var iterator = collection.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext().(V)
		set.AddValue(value)
	}
	return set
}

func (c *setClass_[V]) WithRanker(ranker RankingFunction) SetLike[V] {
	var values = ListClass[V]().Empty()
	var set = &set_[V]{
		ranker,
		values,
	}
	return set
}

// Public Class Functions

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

// CLASS INSTANCES

// Private Class Type Definition

type set_[V Value] struct {
	rank   RankingFunction
	values ListLike[V]
}

// Accessible Interface

func (v *set_[V]) GetValue(index int) V {
	return v.values.GetValue(index)
}

func (v *set_[V]) GetValues(first int, last int) Sequential[V] {
	return v.values.GetValues(first, last)
}

// Flexible Interface

func (v *set_[V]) AddValue(value V) {
	var slot, found = v.findIndex(value)
	if !found {
		// The value is not already a member, so add it.
		v.values.InsertValue(slot, value)
	}
}

func (v *set_[V]) AddValues(values Sequential[V]) {
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		v.AddValue(value)
	}
}

func (v *set_[V]) GetRanker() RankingFunction {
	return v.rank
}

func (v *set_[V]) RemoveAll() {
	v.values.RemoveAll()
}

func (v *set_[V]) RemoveValue(value V) {
	var index, found = v.findIndex(value)
	if found {
		// The value is a member, so remove it.
		v.values.RemoveValue(index)
	}
}

func (v *set_[V]) RemoveValues(values Sequential[V]) {
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		v.RemoveValue(value)
	}
}

// Searchable Interface

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

func (v *set_[V]) ContainsValue(value V) bool {
	var _, found = v.findIndex(value)
	return found
}

func (v *set_[V]) GetComparer() ComparingFunction {
	return v.values.GetComparer()
}

func (v *set_[V]) GetIndex(value V) int {
	var index, found = v.findIndex(value)
	if !found {
		return 0
	}
	return index
}

// Sequential Interface

func (v *set_[V]) AsArray() []V {
	return v.values.AsArray()
}

func (v *set_[V]) GetIterator() Ratcheted[V] {
	var iterator = v.values.GetIterator()
	return iterator
}

func (v *set_[V]) GetSize() int {
	return v.values.GetSize()
}

func (v *set_[V]) IsEmpty() bool {
	return v.values.IsEmpty()
}

// Private Interface

// This private class method performs a binary search of the set for the
// specified value. It returns two results:
//   - index: The index of the value, or if not found, the slot in which it
//     could be inserted in the underlying list.
//   - found: A boolean stating whether or not the value was found.
//
// The algorithm performs a true O[log(n)] worst case search.
func (v *set_[V]) findIndex(value V) (index int, found bool) {
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
// the set.
func (v *set_[V]) String() string {
	var cdcn = CDCNClass().Default()
	return cdcn.FormatCollection(v)
}
