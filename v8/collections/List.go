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

func ListClass[V any]() ListClassLike[V] {
	return listClass[V]()
}

// Constructor Methods

func (c *listClass_[V]) List() ListLike[V] {
	var instance = &list_[V]{
		array_: []V{},
	}
	return instance
}

func (c *listClass_[V]) ListFromArray(
	values []V,
) ListLike[V] {
	var instance = &list_[V]{
		array_: uti.CopyArray(values),
	}
	return instance
}

func (c *listClass_[V]) ListFromSequence(
	values Sequential[V],
) ListLike[V] {
	var instance = &list_[V]{
		array_: values.AsArray(),
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
	index int,
) V {
	var size = v.GetSize()
	var slot = uti.RelativeToCardinal(index, size)
	var value = v.array_[slot]
	return value
}

func (v *list_[V]) GetValues(
	first int,
	last int,
) Sequential[V] {
	var size = v.GetSize()
	var goFirst = uti.RelativeToCardinal(first, size)
	var goLast = uti.RelativeToCardinal(last, size) + 1
	var values = listClass[V]().ListFromArray(v.array_[goFirst:goLast])
	return values
}

func (v *list_[V]) GetIndex(
	value V,
) int {
	var index int
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

// Malleable[V] Methods

func (v *list_[V]) InsertValue(
	slot uint,
	value V,
) {
	// Create a new larger array.
	var size = v.GetSize() + 1
	var array = make([]V, size)

	// Copy the values into the new array.
	copy(array, v.array_[:slot])
	array[slot] = value
	copy(array[slot+1:], v.array_[slot:])

	// Update the internal array.
	v.array_ = array
}

func (v *list_[V]) InsertValues(
	slot uint,
	values Sequential[V],
) {
	// Create a new larger array.
	var newValues = values.AsArray()
	var delta = uti.ArraySize(newValues)
	var size = uti.ArraySize(v.array_) + delta
	var array = make([]V, size)

	// Copy the values into the new array.
	copy(array, v.array_[:slot])
	copy(array[slot:], newValues)
	copy(array[slot+delta:], v.array_[slot:])

	// Update the internal array.
	v.array_ = array
}

func (v *list_[V]) AppendValue(
	value V,
) {
	// Create a new larger array.
	var size = len(v.array_) + 1
	var array = make([]V, size)

	// Copy the values into the new array.
	copy(array, v.array_)
	array[size-1] = value

	// Update the internal array.
	v.array_ = array
}

func (v *list_[V]) AppendValues(
	values Sequential[V],
) {
	// Create a new larger array.
	var newValues = values.AsArray()
	var size = len(v.array_) + len(newValues)
	var array = make([]V, size)

	// Copy the values into the new array.
	copy(array, v.array_)
	copy(array[len(v.array_):], newValues)

	// Update the internal array.
	v.array_ = array
}

func (v *list_[V]) RemoveValue(
	index int,
) V {
	// Convert to zero-based index.
	var size = v.GetSize()
	var slot = uti.RelativeToCardinal(index, size)

	// Create a new smaller array.
	size--
	var array = make([]V, size)

	// Copy the values into the new array.
	copy(array, v.array_[:slot])
	var removed = v.array_[slot]
	copy(array[slot:], v.array_[slot+1:])

	// Update the internal array.
	v.array_ = array
	return removed
}

func (v *list_[V]) RemoveValues(
	first int,
	last int,
) Sequential[V] {
	// Create two smaller arrays.
	var size = v.GetSize()
	var goFirst = uti.RelativeToCardinal(first, size)
	var goLast = uti.RelativeToCardinal(last, size) + 1
	var delta = uint(goLast - goFirst)
	size -= delta
	var array = make([]V, size)
	var removed = make([]V, delta)

	// Copy the existing values into the two new arrays.
	copy(array, v.array_[:goFirst])
	copy(removed, v.array_[goFirst:goLast])
	copy(array[goFirst:], v.array_[goLast:])

	// Update the internal array.
	v.array_ = array

	// Return a list of the removed values.
	var values = listClass[V]().ListFromArray(removed)
	return values
}

func (v *list_[V]) RemoveAll() {
	v.array_ = []V{}
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

// Sequential[V] Methods

func (v *list_[V]) IsEmpty() bool {
	return len(v.array_) == 0
}

func (v *list_[V]) GetSize() uint {
	return uti.ArraySize(v.array_)
}

func (v *list_[V]) AsArray() []V {
	var array = uti.CopyArray(v.array_)
	return array
}

func (v *list_[V]) GetIterator() age.IteratorLike[V] {
	var array = uti.CopyArray(v.array_)
	var iteratorClass = age.IteratorClass[V]()
	var iterator = iteratorClass.Iterator(array)
	return iterator
}

// Sortable[V] Methods

func (v *list_[V]) SortValues() {
	var sorterClass = age.SorterClass[V]()
	var sorter = sorterClass.Sorter()
	sorter.SortValues(v.array_)
}

func (v *list_[V]) SortValuesWithRanker(
	ranker age.RankingFunction[V],
) {
	var sorterClass = age.SorterClass[V]()
	var sorter = sorterClass.SorterWithRanker(ranker)
	sorter.SortValues(v.array_)
}

func (v *list_[V]) ReverseValues() {
	var sorterClass = age.SorterClass[V]()
	var sorter = sorterClass.Sorter()
	sorter.ReverseValues(v.array_)
}

func (v *list_[V]) ShuffleValues() {
	var sorterClass = age.SorterClass[V]()
	var sorter = sorterClass.Sorter()
	sorter.ShuffleValues(v.array_)
}

// Updatable[V] Methods

func (v *list_[V]) SetValue(
	index int,
	value V,
) {
	var size = v.GetSize()
	var slot = uti.RelativeToCardinal(index, size)
	v.array_[slot] = value
}

func (v *list_[V]) SetValues(
	index int,
	values Sequential[V],
) {
	var size = v.GetSize()
	var slot = uti.RelativeToCardinal(index, size)
	var newValues = values.AsArray()
	copy(v.array_[slot:], newValues)
}

// PROTECTED INTERFACE

func (v *list_[V]) String() string {
	return uti.Format(v)
}

// Private Methods

// Instance Structure

type list_[V any] struct {
	// Declare the instance attributes.
	array_ []V
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
