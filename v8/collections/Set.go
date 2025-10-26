/*
................................................................................
.    Copyright (c) 2009-2025 Crater Dog Technologies.  All Rights Reserved.    .
................................................................................
.  DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               .
.                                                                              .
.  This code is free software; you can redistribute it and/or modify it under  .
.  the terms of The MIT License (MIT), as published by the Open Source         .
.  Initiative. (See https://opensource.org/license/MIT)                        .
................................................................................
*/

package collections

import (
	fmt "fmt"
	age "github.com/craterdog/go-collection-framework/v8/agents"
	uti "github.com/craterdog/go-missing-utilities/v8"
	syn "sync"
)

// CLASS INTERFACE

// Access Function

func SetClass[V any]() SetClassLike[V] {
	return setClass[V]()
}

// Constructor Methods

func (c *setClass_[V]) Set() SetLike[V] {
	var collatorClass = age.CollatorClass[V]()
	var collator = collatorClass.Collator()
	var instance = c.SetWithCollator(collator)
	return instance
}

func (c *setClass_[V]) SetWithCollator(
	collator age.CollatorLike[V],
) SetLike[V] {
	if uti.IsUndefined(collator) {
		panic("The \"collator\" attribute is required by this class.")
	}
	var listClass = ListClass[V]()
	var values = listClass.List()
	var instance = &set_[V]{
		// Initialize the instance attributes.
		collator_: collator,
		values_:   values,
	}
	return instance
}

func (c *setClass_[V]) SetFromArray(
	values []V,
) SetLike[V] {
	var set = c.Set()
	for _, value := range values {
		set.AddValue(value)
	}
	return set
}

func (c *setClass_[V]) SetFromSequence(
	values Sequential[V],
) SetLike[V] {
	var set = c.Set()
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		set.AddValue(value)
	}
	return set
}

// Constant Methods

// Function Methods

func (c *setClass_[V]) And(
	first SetLike[V],
	second SetLike[V],
) SetLike[V] {
	var collator = first.GetCollator()
	var result = c.SetWithCollator(collator)
	var iterator = first.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		if second.ContainsValue(value) {
			result.AddValue(value)
		}
	}
	return result
}

func (c *setClass_[V]) Ior(
	first SetLike[V],
	second SetLike[V],
) SetLike[V] {
	var collator = first.GetCollator()
	var result = c.SetWithCollator(collator)
	result.AddValues(first)
	result.AddValues(second)
	return result
}

func (c *setClass_[V]) San(
	first SetLike[V],
	second SetLike[V],
) SetLike[V] {
	var collator = first.GetCollator()
	var result = c.SetWithCollator(collator)
	result.AddValues(first)
	result.RemoveValues(second)
	return result
}

func (c *setClass_[V]) Xor(
	first SetLike[V],
	second SetLike[V],
) SetLike[V] {
	return c.Ior(c.San(first, second), c.San(second, first))
}

// INSTANCE INTERFACE

// Principal Methods

func (v *set_[V]) GetClass() SetClassLike[V] {
	return setClass[V]()
}

// Attribute Methods

func (v *set_[V]) GetCollator() age.CollatorLike[V] {
	return v.collator_
}

// Accessible[V] Methods

func (v *set_[V]) GetValue(
	index int,
) V {
	var value = v.values_.GetValue(index)
	return value
}

func (v *set_[V]) GetValues(
	first int,
	last int,
) Sequential[V] {
	var values = v.values_.GetValues(first, last)
	return values
}

func (v *set_[V]) GetIndex(
	value V,
) int {
	var index, found = v.findIndex(value)
	if !found {
		return 0
	}
	return index
}

// Elastic[V] Methods

func (v *set_[V]) AddValue(
	value V,
) {
	var slot, found = v.findIndex(value)
	if !found {
		// The value is not already a member, so add it.
		v.values_.InsertValue(uint(slot), value)
	}
}

func (v *set_[V]) AddValues(
	values Sequential[V],
) {
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		v.AddValue(value)
	}
}

func (v *set_[V]) RemoveValue(
	value V,
) {
	var index, found = v.findIndex(value)
	if found {
		// The value is a member, so remove it.
		v.values_.RemoveValue(index)
	}
}

func (v *set_[V]) RemoveValues(
	values Sequential[V],
) {
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		v.RemoveValue(value)
	}
}

func (v *set_[V]) RemoveAll() {
	v.values_.RemoveAll()
}

// Searchable[V] Methods

func (v *set_[V]) ContainsValue(
	value V,
) bool {
	var _, found = v.findIndex(value)
	return found
}

func (v *set_[V]) ContainsAny(
	values Sequential[V],
) bool {
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

func (v *set_[V]) ContainsAll(
	values Sequential[V],
) bool {
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

// Sequential[V] Methods

func (v *set_[V]) IsEmpty() bool {
	return v.values_.IsEmpty()
}

func (v *set_[V]) GetSize() uint {
	var size = v.values_.GetSize()
	return size
}

func (v *set_[V]) AsArray() []V {
	var array = v.values_.AsArray()
	return array
}

func (v *set_[V]) GetIterator() uti.IteratorLike[V] {
	var iterator = v.values_.GetIterator()
	return iterator
}

// PROTECTED INTERFACE

func (v *set_[V]) String() string {
	return uti.Format(v)
}

// Private Methods

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
	var first = 1               // Start at the beginning.
	var last = int(v.GetSize()) // End at the end.
	var size = last             // Initially all values are candidates.
	for size > 0 {
		var middle = first + size/2 // Rounds down to the nearest integer.
		var candidate = v.GetValue(int(middle))
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

// Instance Structure

type set_[V any] struct {
	// Declare the instance attributes.
	collator_ age.CollatorLike[V]
	values_   ListLike[V]
}

// Class Structure

type setClass_[V any] struct {
	// Declare the class constants.
}

// Class Reference

var setMap_ = map[string]any{}
var setMutex_ syn.Mutex

func setClass[V any]() *setClass_[V] {
	// Generate the name of the bound class type.
	var class *setClass_[V]
	var name = fmt.Sprintf("%T", class)

	// Check for an existing bound class type.
	setMutex_.Lock()
	var value = setMap_[name]
	switch actual := value.(type) {
	case *setClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Add a new bound class type.
		class = &setClass_[V]{
			// Initialize the class constants.
		}
		setMap_[name] = class
	}
	setMutex_.Unlock()

	// Return a reference to the bound class type.
	return class
}
