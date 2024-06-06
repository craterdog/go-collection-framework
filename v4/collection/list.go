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

var listClass = map[string]any{}
var listMutex syn.Mutex

// Function

func List[V any](notation NotationLike) ListClassLike[V] {
	// Validate the notation argument.
	if notation == nil {
		panic("A notation must be specified when creating this class.")
	}

	// Generate the name of the bound class type.
	var class ListClassLike[V]
	var name = fmt.Sprintf("%T-%T", class, notation)

	// Check for existing bound class type.
	listMutex.Lock()
	var value = listClass[name]
	switch actual := value.(type) {
	case *listClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Add a new bound class type.
		class = &listClass_[V]{
			notation_: notation,
		}
		listClass[name] = class
	}
	listMutex.Unlock()

	// Return a reference to the bound class type.
	return class
}

// CLASS METHODS

// Target

type listClass_[V any] struct {
	notation_ NotationLike
}

// Constants

func (c *listClass_[V]) Notation() NotationLike {
	return c.notation_
}

// Constructors

func (c *listClass_[V]) Make() ListLike[V] {
	var values = Array[V](c.notation_).MakeFromSize(0)
	return &list_[V]{
		class_:  c,
		values_: values,
	}
}

func (c *listClass_[V]) MakeFromArray(values []V) ListLike[V] {
	var array = Array[V](c.notation_).MakeFromArray(values)
	return c.MakeFromSequence(array)
}

func (c *listClass_[V]) MakeFromSequence(values Sequential[V]) ListLike[V] {
	var list = c.Make()
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		list.AppendValue(value)
	}
	return list
}

func (c *listClass_[V]) MakeFromSource(source string) ListLike[V] {
	// First we parse it as a collection of any type value.
	var collection = c.notation_.ParseSource(source).(Sequential[any])

	// Next we must convert each value explicitly to type V.
	var anys = collection.AsArray()
	var array = make([]V, len(anys))
	for index, value := range anys {
		array[index] = value.(V)
	}

	// Then we can create the stack from the type V array.
	return c.MakeFromArray(array)
}

// Functions

func (c *listClass_[V]) Concatenate(first, second ListLike[V]) ListLike[V] {
	var list = c.Make()
	list.AppendValues(first)
	list.AppendValues(second)
	return list
}

// INSTANCE METHODS

// Target

type list_[V any] struct {
	class_  ListClassLike[V]
	values_ ArrayLike[V]
}

// Attributes

func (v *list_[V]) GetClass() ListClassLike[V] {
	return v.class_
}

// Accessible

func (v *list_[V]) GetValue(index int) V {
	return v.values_.GetValue(index)
}

func (v *list_[V]) GetValues(first int, last int) Sequential[V] {
	return v.values_.GetValues(first, last)
}

// Expandable

func (v *list_[V]) InsertValue(slot int, value V) {

	// Create a new larger array.
	var size = v.GetSize() + 1
	var array = Array[V](v.class_.Notation()).MakeFromSize(size)

	// Copy the values into the new array.
	var iterator = v.GetIterator()
	var index int
	for index < size {
		if index == slot {
			index++
			array.SetValue(index, value)
		} else {
			var existing = iterator.GetNext()
			index++
			array.SetValue(index, existing)
		}
	}

	// Update the internal array.
	v.values_ = array
}

func (v *list_[V]) InsertValues(slot int, values Sequential[V]) {

	// Create a new larger array.
	var size = v.GetSize() + values.GetSize()
	var array = Array[V](v.class_.Notation()).MakeFromSize(size)

	// Copy the values into the new array.
	var iterator = v.GetIterator()
	var index int
	for index < size {
		if index == slot {
			var iterator2 = values.GetIterator()
			for iterator2.HasNext() {
				index++
				var value = iterator2.GetNext()
				array.SetValue(index, value)
			}
		} else {
			var existing = iterator.GetNext()
			index++
			array.SetValue(index, existing)
		}
	}

	// Update the internal array.
	v.values_ = array
}

func (v *list_[V]) AppendValue(value V) {

	// Create a new larger array.
	var size = v.GetSize() + 1
	var array = Array[V](v.class_.Notation()).MakeFromSize(size)

	// Copy the existing values into the new array.
	var index int
	var iterator = v.GetIterator()
	for iterator.HasNext() {
		index++
		var existing = iterator.GetNext()
		array.SetValue(index, existing)
	}

	// Append the new value to the new array.
	index++
	array.SetValue(index, value)

	// Update the internal array.
	v.values_ = array
}

func (v *list_[V]) AppendValues(values Sequential[V]) {

	// Create a new larger array.
	var size = v.GetSize() + values.GetSize()
	var array = Array[V](v.class_.Notation()).MakeFromSize(size)

	// Copy the existing values into the new array.
	var index int
	var iterator = v.GetIterator()
	for iterator.HasNext() {
		index++
		var existing = iterator.GetNext()
		array.SetValue(index, existing)
	}

	// Append the new values to the new array.
	iterator = values.GetIterator()
	for iterator.HasNext() {
		index++
		var value = iterator.GetNext()
		array.SetValue(index, value)
	}

	// Update the internal array.
	v.values_ = array
}

func (v *list_[V]) RemoveValue(index int) V {

	// Create a new smaller array.
	var removed = v.GetValue(index)
	var size = v.GetSize() - 1
	var array = Array[V](v.class_.Notation()).MakeFromSize(size)

	// Copy the remaining values into the new array.
	var counter = v.toNormalized(index)
	index = 1
	var iterator = v.GetIterator()
	for iterator.HasNext() {
		counter--
		var value = iterator.GetNext()
		if counter == 0 {
			continue // Skip this value.
		}
		array.SetValue(index, value)
		index++
	}

	// Update the internal array.
	v.values_ = array

	return removed
}

func (v *list_[V]) RemoveValues(first int, last int) Sequential[V] {

	// Create two smaller arrays.
	first = v.toNormalized(first)
	last = v.toNormalized(last)
	var delta = last - first + 1
	var size = v.GetSize() - delta
	var Array = Array[V](v.class_.Notation())
	var removed = Array.MakeFromSize(delta)
	var array = Array.MakeFromSize(size)

	// Split the existing values into the two new arrays.
	var counter int
	var arrayIndex int
	var removedIndex int
	var iterator = v.GetIterator()
	for iterator.HasNext() {
		counter++
		var existing = iterator.GetNext()
		if counter < first || counter > last {
			arrayIndex++
			array.SetValue(arrayIndex, existing)
		} else {
			removedIndex++
			removed.SetValue(removedIndex, existing)
		}
	}

	// Update the internal array.
	v.values_ = array

	return removed
}

func (v *list_[V]) RemoveAll() {
	v.values_ = Array[V](v.class_.Notation()).MakeFromSize(0)
}

// Searchable

func (v *list_[V]) ContainsValue(value V) bool {
	return v.GetIndex(value) > 0
}

func (v *list_[V]) ContainsAny(values Sequential[V]) bool {
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

func (v *list_[V]) ContainsAll(values Sequential[V]) bool {
	var iterator = values.GetIterator()
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

func (v *list_[V]) GetIndex(value V) int {
	var compare = age.Collator[V]().Make().CompareValues
	for index, candidate := range v.AsArray() {
		if compare(candidate, value) {
			// Found the value.
			return index + 1 // Convert to an ORDINAL based index.
		}
	}
	// The value was not found.
	return 0
}

// Sequential

func (v *list_[V]) IsEmpty() bool {
	return v.values_.IsEmpty()
}

func (v *list_[V]) GetSize() int {
	return v.values_.GetSize()
}

func (v *list_[V]) AsArray() []V {
	return v.values_.AsArray()
}

func (v *list_[V]) GetIterator() age.IteratorLike[V] {
	return v.values_.GetIterator()
}

// Sortable

func (v *list_[V]) SortValues() {
	v.values_.SortValues()
}

func (v *list_[V]) SortValuesWithRanker(ranker age.RankingFunction[V]) {
	v.values_.SortValuesWithRanker(ranker)
}

func (v *list_[V]) ReverseValues() {
	v.values_.ReverseValues()
}

func (v *list_[V]) ShuffleValues() {
	v.values_.ShuffleValues()
}

// Updatable

func (v *list_[V]) SetValue(index int, value V) {
	v.values_.SetValue(index, value)
}

func (v *list_[V]) SetValues(index int, values Sequential[V]) {
	v.values_.SetValues(index, values)
}

// Stringer

func (v *list_[V]) String() string {
	var notation = v.class_.Notation()
	return notation.FormatCollection(v)
}

// Private

// This private instance method normalizes the specified relative index.  The
// following transformation is performed:
//
//	[-size..-1] and [1..size] => [1..size]
//
// Notice that the specified index cannot be zero since zero is NOT an ordinal.
func (v *list_[V]) toNormalized(index int) int {
	var size = v.GetSize()
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
