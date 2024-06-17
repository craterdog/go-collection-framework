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

package collection

import (
	fmt "fmt"
	age "github.com/craterdog/go-collection-framework/v4/agent"
	syn "sync"
)

// CLASS ACCESS

// Reference

var setClass = map[string]any{}
var setMutex syn.Mutex

// Function

func Set[V any](notation NotationLike) SetClassLike[V] {
	// Generate the name of the bound class type.
	var class *setClass_[V]
	var name = fmt.Sprintf("%T", class)

	// Check for existing bound class type.
	setMutex.Lock()
	var value = setClass[name]
	switch actual := value.(type) {
	case *setClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Add a new bound class type.
		class = &setClass_[V]{
			// Initialize the class constants.
			notation_: notation,
		}
		setClass[name] = class
	}
	setMutex.Unlock()

	// Return a reference to the bound class type.
	return class
}

// CLASS METHODS

// Target

type setClass_[V any] struct {
	// Define the class constants.
	notation_ NotationLike
}

// Constants

func (c *setClass_[V]) Notation() NotationLike {
	return c.notation_
}

// Constructors

func (c *setClass_[V]) Make() SetLike[V] {
	var collator = age.Collator[V]().Make()
	var set = c.MakeWithCollator(collator)
	return set
}

func (c *setClass_[V]) MakeWithCollator(collator age.CollatorLike[V]) SetLike[V] {
	var values = List[V](c.notation_).Make()
	return &set_[V]{
		class_:    c,
		collator_: collator,
		values_:   values,
	}
}

func (c *setClass_[V]) MakeFromArray(values []V) SetLike[V] {
	var array = Array[V](c.notation_).MakeFromArray(values)
	var set = c.MakeFromSequence(array)
	return set
}

func (c *setClass_[V]) MakeFromSequence(values Sequential[V]) SetLike[V] {
	var set = c.Make()
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		set.AddValue(value)
	}
	return set
}

// Functions

func (c *setClass_[V]) And(first, second SetLike[V]) SetLike[V] {
	var result = c.MakeWithCollator(first.GetCollator())
	var iterator = first.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		if second.ContainsValue(value) {
			result.AddValue(value)
		}
	}
	return result
}

func (c *setClass_[V]) Or(first, second SetLike[V]) SetLike[V] {
	var result = c.MakeWithCollator(first.GetCollator())
	result.AddValues(first)
	result.AddValues(second)
	return result
}

func (c *setClass_[V]) Sans(first, second SetLike[V]) SetLike[V] {
	var result = c.MakeWithCollator(first.GetCollator())
	result.AddValues(first)
	result.RemoveValues(second)
	return result
}

func (c *setClass_[V]) Xor(first, second SetLike[V]) SetLike[V] {
	return c.Or(c.Sans(first, second), c.Sans(second, first))
}

// INSTANCE METHODS

// Target

type set_[V any] struct {
	class_    SetClassLike[V]
	collator_ age.CollatorLike[V]
	values_   ListLike[V]
}

// Attributes

func (v *set_[V]) GetClass() SetClassLike[V] {
	return v.class_
}

func (v *set_[V]) GetCollator() age.CollatorLike[V] {
	return v.collator_
}

// Accessible

func (v *set_[V]) GetValue(index int) V {
	return v.values_.GetValue(index)
}

func (v *set_[V]) GetValues(first int, last int) Sequential[V] {
	return v.values_.GetValues(first, last)
}

// Flexible

func (v *set_[V]) AddValue(value V) {
	var slot, found = v.findIndex(value)
	if !found {
		// The value is not already a member, so add it.
		v.values_.InsertValue(uint(slot), value)
	}
}

func (v *set_[V]) AddValues(values Sequential[V]) {
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		v.AddValue(value)
	}
}

func (v *set_[V]) RemoveValue(value V) {
	var index, found = v.findIndex(value)
	if found {
		// The value is a member, so remove it.
		v.values_.RemoveValue(index)
	}
}

func (v *set_[V]) RemoveValues(values Sequential[V]) {
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		v.RemoveValue(value)
	}
}

func (v *set_[V]) RemoveAll() {
	v.values_.RemoveAll()
}

// Searchable

func (v *set_[V]) ContainsValue(value V) bool {
	var _, found = v.findIndex(value)
	return found
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

func (v *set_[V]) GetIndex(value V) int {
	var index, found = v.findIndex(value)
	if !found {
		return 0
	}
	return index
}

// Sequential

func (v *set_[V]) IsEmpty() bool {
	return v.values_.IsEmpty()
}

func (v *set_[V]) GetSize() int {
	return v.values_.GetSize()
}

func (v *set_[V]) AsArray() []V {
	return v.values_.AsArray()
}

func (v *set_[V]) GetIterator() age.IteratorLike[V] {
	var iterator = v.values_.GetIterator()
	return iterator
}

// Stringer

func (v *set_[V]) String() string {
	return v.GetClass().Notation().FormatValue(v)
}

// Private

// This private instance method performs a binary search of the set for the
// specified value. It returns two results:
//   - index: The index of the value, or if not found, the slot in which it could
//     be inserted in the underlying list.
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
		switch v.collator_.RankValues(value, candidate) {
		case age.LesserRank:
			// The index of the value is less than the middle
			// index so the first index stays the same.
			last = middle - 1 // We already tried the middle index.
			size = middle - first
		case age.EqualRank:
			// The index of the value is the middle index.
			return middle, true
		case age.GreaterRank:
			// The index of the value is greater than the middle
			// index so the last index stays the same.
			first = middle + 1 // We already tried the middle index.
			size = last - middle
		}
	}
	// The value was not found, the last index represents the SLOT where it
	// would be inserted.  Since the value was not found, the indexes are
	// inverted: last < first (i.e. last = first - 1).
	return last, false
}
