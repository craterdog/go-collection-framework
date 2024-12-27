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
	age "github.com/craterdog/go-collection-framework/v5/agent"
	uti "github.com/craterdog/go-missing-utilities/v2"
	syn "sync"
)

// CLASS INTERFACE

// Access Function

func ListClass[V any]() ListClassLike[V] {
	return listClass[V]()
}

// Constructor Methods

func (c *listClass_[V]) List() ListLike[V] {
	var arrayClass = ArrayClass[V]()
	var array = arrayClass.Array(0)
	var instance = &list_[V]{
		array_: array,
	}
	return instance
}

func (c *listClass_[V]) ListFromArray(
	values []V,
) ListLike[V] {
	var arrayClass = ArrayClass[V]()
	var array = arrayClass.ArrayFromArray(values)
	var instance = &list_[V]{
		array_: array,
	}
	return instance
}

func (c *listClass_[V]) ListFromSequence(
	values Sequential[V],
) ListLike[V] {
	var arrayClass = ArrayClass[V]()
	var array = arrayClass.ArrayFromSequence(values)
	var instance = &list_[V]{
		array_: array,
	}
	return instance
}

// Constant Methods

// Function Methods

func (c *listClass_[V]) Concatenate(
	first ListLike[V],
	second ListLike[V],
) ListLike[V] {
	var list = c.ListFromSequence(first)
	list.AppendValues(second)
	return list
}

// INSTANCE INTERFACE

// Principal Methods

func (v *list_[V]) GetClass() ListClassLike[V] {
	return listClass[V]()
}

// Attribute Methods

// Accessible[V] Methods

func (v *list_[V]) GetValue(
	index Index,
) V {
	var value = v.array_.GetValue(index)
	return value
}

func (v *list_[V]) GetValues(
	first Index,
	last Index,
) Sequential[V] {
	var values = v.array_.GetValues(first, last)
	return values
}

// Malleable[V] Methods

func (v *list_[V]) InsertValue(
	slot age.Slot,
	value V,
) {
	// Create a new larger array.
	var size = v.GetSize() + 1
	var arrayClass = ArrayClass[V]()
	var array = arrayClass.Array(size)

	// Copy the values into the new array.
	var existingValues = v.GetIterator()
	var index Index
	for existingValues.HasNext() {
		if existingValues.GetSlot() == slot {
			index++
			array.SetValue(index, value)
		}
		index++
		var existingValue = existingValues.GetNext()
		array.SetValue(index, existingValue)
	}
	if existingValues.GetSlot() == slot {
		index++
		array.SetValue(index, value)
	}

	// Update the internal array.
	v.array_ = array
}

func (v *list_[V]) InsertValues(
	slot age.Slot,
	values Sequential[V],
) {
	// Create a new larger array.
	var size = v.GetSize() + values.GetSize()
	var arrayClass = ArrayClass[V]()
	var array = arrayClass.Array(size)

	// Copy the values into the new array.
	var existingValues = v.GetIterator()
	var index Index
	for existingValues.HasNext() {
		if existingValues.GetSlot() == slot {
			var newValues = values.GetIterator()
			for newValues.HasNext() {
				index++
				var newValue = newValues.GetNext()
				array.SetValue(index, newValue)
			}
		}
		index++
		var existingValue = existingValues.GetNext()
		array.SetValue(index, existingValue)
	}
	if existingValues.GetSlot() == slot {
		var newValues = values.GetIterator()
		for newValues.HasNext() {
			index++
			var newValue = newValues.GetNext()
			array.SetValue(index, newValue)
		}
	}

	// Update the internal array.
	v.array_ = array
}

func (v *list_[V]) AppendValue(
	value V,
) {
	// Create a new larger array.
	var size = v.GetSize() + 1
	var arrayClass = ArrayClass[V]()
	var array = arrayClass.Array(size)

	// Copy the existing values into the new array.
	var existingValues = v.GetIterator()
	var index Index
	for existingValues.HasNext() {
		index++
		var existingValue = existingValues.GetNext()
		array.SetValue(index, existingValue)
	}

	// Copy the new value to the end of the new array.
	index++
	array.SetValue(index, value)

	// Update the internal array.
	v.array_ = array
}

func (v *list_[V]) AppendValues(
	values Sequential[V],
) {
	// Create a new larger array.
	var size = v.GetSize() + values.GetSize()
	var arrayClass = ArrayClass[V]()
	var array = arrayClass.Array(size)

	// Copy the existing values into the new array.
	var existingValues = v.GetIterator()
	var index Index
	for existingValues.HasNext() {
		index++
		var existingValue = existingValues.GetNext()
		array.SetValue(index, existingValue)
	}

	// Copy the new values into the new array.
	var newValues = values.GetIterator()
	for newValues.HasNext() {
		index++
		var newValue = newValues.GetNext()
		array.SetValue(index, newValue)
	}

	// Update the internal array.
	v.array_ = array
}

func (v *list_[V]) RemoveValue(
	index Index,
) V {
	// Create a new smaller array.
	var removed = v.GetValue(index)
	var size = v.GetSize() - 1
	var arrayClass = ArrayClass[V]()
	var array = arrayClass.Array(size)

	// Copy the remaining values into the new array.
	var counter = v.toNormalized(index)
	index = 1
	var existingValues = v.GetIterator()
	for existingValues.HasNext() {
		counter--
		var existingValue = existingValues.GetNext()
		if counter == 0 {
			continue // Skip this value.
		}
		array.SetValue(index, existingValue)
		index++
	}

	// Update the internal array.
	v.array_ = array
	return removed
}

func (v *list_[V]) RemoveValues(
	first Index,
	last Index,
) Sequential[V] {
	// Create two smaller arrays.
	first = v.toNormalized(first)
	last = v.toNormalized(last)
	var delta = age.Size(last - first + 1)
	var size = v.GetSize() - delta
	var arrayClass = ArrayClass[V]()
	var removed = arrayClass.Array(delta)
	var array = arrayClass.Array(size)

	// Split the existing values into the two new arrays.
	var counter Index
	var arrayIndex Index
	var removedIndex Index
	var existingValues = v.GetIterator()
	for existingValues.HasNext() {
		counter++
		var existingValue = existingValues.GetNext()
		if counter < first || counter > last {
			arrayIndex++
			array.SetValue(arrayIndex, existingValue)
		} else {
			removedIndex++
			removed.SetValue(removedIndex, existingValue)
		}
	}

	// Update the internal array.
	v.array_ = array
	return removed
}

func (v *list_[V]) RemoveAll() {
	var arrayClass = ArrayClass[V]()
	v.array_ = arrayClass.Array(0)
}

// Searchable[V] Methods

func (v *list_[V]) ContainsValue(
	value V,
) bool {
	return v.GetIndex(value) > 0
}

func (v *list_[V]) ContainsAny(
	values Sequential[V],
) bool {
	var iterator = values.GetIterator()
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

func (v *list_[V]) ContainsAll(
	values Sequential[V],
) bool {
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var candidate = iterator.GetNext()
		if v.GetIndex(candidate) == 0 {
			// One of the values is missing.
			return false
		}
	}
	// Found all of the values.
	return true
}

func (v *list_[V]) GetIndex(
	value V,
) Index {
	var index Index
	var collatorClass = age.CollatorClass[V]()
	var compare = collatorClass.Collator().CompareValues
	var iterator = v.GetIterator()
	for iterator.HasNext() {
		index++
		var candidate = iterator.GetNext()
		if compare(candidate, value) {
			// Found the value.
			return index
		}
	}
	// The value was not found.
	return 0
}

// Sequential[V] Methods

func (v *list_[V]) IsEmpty() bool {
	return v.array_.IsEmpty()
}

func (v *list_[V]) GetSize() age.Size {
	var size = v.array_.GetSize()
	return size
}

func (v *list_[V]) AsArray() []V {
	var array = v.array_.AsArray()
	return array
}

func (v *list_[V]) GetIterator() age.IteratorLike[V] {
	var iterator = v.array_.GetIterator()
	return iterator
}

// Sortable[V] Methods

func (v *list_[V]) SortValues() {
	v.array_.SortValues()
}

func (v *list_[V]) SortValuesWithRanker(
	ranker age.RankingFunction[V],
) {
	v.array_.SortValuesWithRanker(ranker)
}

func (v *list_[V]) ReverseValues() {
	v.array_.ReverseValues()
}

func (v *list_[V]) ShuffleValues() {
	v.array_.ShuffleValues()
}

// Updatable[V] Methods

func (v *list_[V]) SetValue(
	index Index,
	value V,
) {
	v.array_.SetValue(index, value)
}

func (v *list_[V]) SetValues(
	index Index,
	values Sequential[V],
) {
	v.array_.SetValues(index, values)
}

// Stringer Methods

func (v *list_[V]) String() string {
	return uti.Format(v)
}

// PROTECTED INTERFACE

// Private Methods

// This private instance method normalizes the specified relative index.  The
// following transformation is performed:
//
//	[-size..-1] and [1..size] => [1..size]
//
// Notice that the specified index cannot be zero since zero is NOT an ordinal.
func (v *list_[V]) toNormalized(index Index) Index {
	var size = Index(v.GetSize())
	switch {
	case size == 0:
		// The list is empty.
		panic("Cannot index an empty List.")
	case index == 0:
		// Zero is not an ordinal.
		panic("Indices must be positive or negative ordinals, not zero.")
	case index < -size || index > size:
		// The index is outside the bounds of the specified index range.
		panic(fmt.Sprintf(
			"The specified index is outside the allowed ranges [-%v..-1] and [1..%v]: %v",
			size,
			size,
			index))
	case index < 0:
		// Convert a negative index.
		return index + size + 1
	case index > 0:
		// Leave it as it is.
		return index
	default:
		// This should never happen so time to...
		panic(fmt.Sprintf("Go compiler problem, unexpected index value: %v", index))
	}
}

// Instance Structure

type list_[V any] struct {
	// Declare the instance attributes.
	array_ ArrayLike[V]
}

// Class Structure

type listClass_[V any] struct {
	// Declare the class constants.
}

// Class Reference

var listMap_ = map[string]any{}
var listMutex_ syn.Mutex

func listClass[V any]() *listClass_[V] {
	// Generate the name of the bound class type.
	var class *listClass_[V]
	var name = fmt.Sprintf("%T", class)

	// Check for an existing bound class type.
	listMutex_.Lock()
	var value = listMap_[name]
	switch actual := value.(type) {
	case *listClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Add a new bound class type.
		class = &listClass_[V]{
			// Initialize the class constants.
		}
		listMap_[name] = class
	}
	listMutex_.Unlock()

	// Return a reference to the bound class type.
	return class
}
