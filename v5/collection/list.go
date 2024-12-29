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
	index Index,
) V {
	var goIndex = v.toZeroBased(index)
	var value = v.array_[goIndex]
	return value
}

func (v *list_[V]) GetValues(
	first Index,
	last Index,
) Sequential[V] {
	var goFirst = v.toZeroBased(first)
	var goLast = v.toZeroBased(last) + 1 // In Go, the last index is exclusive.
	var values = listClass[V]().ListFromArray(v.array_[goFirst:goLast])
	return values
}

// Malleable[V] Methods

func (v *list_[V]) InsertValue(
	slot age.Slot,
	value V,
) {
	// Create a new larger array.
	var size = len(v.array_) + 1
	var array = make([]V, size)

	// Copy the values into the new array.
	var goIndex = int(slot) // The Go index is after the matching slot.
	copy(array, v.array_[:goIndex])
	array[goIndex] = value
	copy(array[goIndex+1:], v.array_[goIndex:])

	// Update the internal array.
	v.array_ = array
}

func (v *list_[V]) InsertValues(
	slot age.Slot,
	values Sequential[V],
) {
	// Create a new larger array.
	var newValues = values.AsArray()
	var size = len(v.array_) + len(newValues)
	var array = make([]V, size)

	// Copy the values into the new array.
	var goIndex = int(slot) // The Go index is after the matching slot.
	copy(array, v.array_[:goIndex])
	copy(array[goIndex:], newValues)
	copy(array[goIndex+len(newValues):], v.array_[goIndex:])

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
	index Index,
) V {
	// Create a new smaller array.
	var size = len(v.array_) - 1
	var array = make([]V, size)

	// Copy the values into the new array.
	var goIndex = v.toZeroBased(index)
	copy(array, v.array_[:goIndex])
	var removed = v.array_[goIndex]
	copy(array[goIndex:], v.array_[goIndex+1:])

	// Update the internal array.
	v.array_ = array
	return removed
}

func (v *list_[V]) RemoveValues(
	first Index,
	last Index,
) Sequential[V] {
	// Create two smaller arrays.
	var goFirst = v.toZeroBased(first)
	var goLast = v.toZeroBased(last) + 1 // In Go, the last index is exclusive.
	var delta = goLast - goFirst
	var size = len(v.array_) - delta
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
	return len(v.array_) == 0
}

func (v *list_[V]) GetSize() age.Size {
	return age.Size(len(v.array_))
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
	index Index,
	value V,
) {
	var goIndex = v.toZeroBased(index)
	v.array_[goIndex] = value
}

func (v *list_[V]) SetValues(
	index Index,
	values Sequential[V],
) {
	var goIndex = v.toZeroBased(index)
	var newValues = values.AsArray()
	copy(v.array_[goIndex:], newValues)
}

// Stringer Methods

func (v *list_[V]) String() string {
	return uti.Format(v)
}

// PROTECTED INTERFACE

// Private Methods

// This private instance method transforms a relative (ORDINAL-based) index into
// the corresponding Go (ZERO-based) index.  The following transformation is
// performed:
//
//	[-size..-1] and [1..size] => [0..size)
//
// Notice that the specified relative index cannot be zero since zero is NOT an
// ordinal.
func (v *list_[V]) toZeroBased(index Index) int {
	var size = Index(len(v.array_))
	switch {
	case size == 0:
		// The Array is empty.
		panic("Cannot index an empty Array.")
	case index == 0:
		// Zero is not an ordinal.
		panic("Indices must be positive or negative ordinals, not zero.")
	case index < -size || index > size:
		// The index is outside the bounds of the specified range.
		panic(fmt.Sprintf(
			"The specified index is outside the allowed ranges [-%v..-1] and [1..%v]: %v",
			size,
			size,
			index))
	case index < 0:
		// Convert a negative index.
		return int(index + size)
	case index > 0:
		// Convert a positive index.
		return int(index - 1)
	default:
		// This should never happen so time to...
		panic(fmt.Sprintf("Go compiler problem, unexpected index value: %v", index))
	}
}

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
