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

type listClass_[V Value] struct {
	// This class defines no constants.
}

// Private Namespace Reference(s)

var listClass = map[string]any{}

// Public Namespace Access

func List[V Value]() ListClassLike[V] {
	var class *listClass_[V]
	var key = fmt.Sprintf("%T", class) // The name of the bound class type.
	var value = listClass[key]
	switch actual := value.(type) {
	case *listClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Create a new bound class type.
		class = &listClass_[V]{
			// This class defines no constants.
		}
		listClass[key] = class
	}
	return class
}

// Public Class Constructors

func (c *listClass_[V]) Empty() ListLike[V] {
	var comparer = Collator().Default().CompareValues
	var list = c.WithComparer(comparer)
	return list
}

func (c *listClass_[V]) FromArray(values []V) ListLike[V] {
	var array = Array[V]().FromArray(values)
	var list = c.FromSequence(array)
	return list
}

func (c *listClass_[V]) FromSequence(values Sequential[V]) ListLike[V] {
	var list = c.Empty()
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		list.AppendValue(value)
	}
	return list
}

func (c *listClass_[V]) FromString(values string) ListLike[V] {
	// First we parse it as a collection of any type value.
	var collection = CDCN().Default().ParseCollection(values).(Sequential[Value])

	// Then we convert it to a list of type V.
	var list = c.Empty()
	var iterator = collection.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext().(V)
		list.AppendValue(value)
	}
	return list
}

func (c *listClass_[V]) WithComparer(comparer ComparingFunction) ListLike[V] {
	var Array = Array[V]() // Retrieve the list namespace.
	var values = Array.WithSize(0)
	var list = &list_[V]{
		compare: comparer,
		values:  values,
	}
	return list
}

// Public Class Functions

// This public class function returns the concatenation of the two specified
// lists.
func (c *listClass_[V]) Concatenate(first, second ListLike[V]) ListLike[V] {
	var list = c.Empty()
	list.AppendValues(first)
	list.AppendValues(second)
	return list
}

// CLASS TYPE

// Private Class Type Definition

type list_[V Value] struct {
	compare ComparingFunction
	values  ArrayLike[V]
}

// Accessible Interface

func (v *list_[V]) GetValue(index int) V {
	return v.values.GetValue(index)
}

func (v *list_[V]) GetValues(first int, last int) Sequential[V] {
	return v.values.GetValues(first, last)
}

// Expandable Interface

func (v *list_[V]) AppendValue(value V) {

	// Create a new larger Array.
	var size = v.GetSize() + 1
	var Array = Array[V]()
	var array = Array.WithSize(size)

	// Copy the existing values into the new Array.
	var index int
	var iterator = v.GetIterator()
	for iterator.HasNext() {
		index++
		var existing = iterator.GetNext()
		array.SetValue(index, existing)
	}

	// Append the new value to the new Array.
	index++
	array.SetValue(index, value)

	// Update the internal Array.
	v.values = array
}

func (v *list_[V]) AppendValues(values Sequential[V]) {

	// Create a new larger Array.
	var size = v.GetSize() + values.GetSize()
	var Array = Array[V]()
	var array = Array.WithSize(size)

	// Copy the existing values into the new Array.
	var index int
	var iterator = v.GetIterator()
	for iterator.HasNext() {
		index++
		var existing = iterator.GetNext()
		array.SetValue(index, existing)
	}

	// Append the new values to the new Array.
	iterator = values.GetIterator()
	for iterator.HasNext() {
		index++
		var value = iterator.GetNext()
		array.SetValue(index, value)
	}

	// Update the internal Array.
	v.values = array
}

func (v *list_[V]) InsertValue(slot int, value V) {

	// Create a new larger Array.
	var size = v.GetSize() + 1
	var Array = Array[V]()
	var array = Array.WithSize(size)

	// Copy the values into the new Array.
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

	// Update the internal Array.
	v.values = array
}

func (v *list_[V]) InsertValues(slot int, values Sequential[V]) {

	// Create a new larger Array.
	var size = v.GetSize() + values.GetSize()
	var Array = Array[V]()
	var array = Array.WithSize(size)

	// Copy the values into the new Array.
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

	// Update the internal Array.
	v.values = array
}

func (v *list_[V]) RemoveAll() {
	var Array = Array[V]()
	v.values = Array.WithSize(0)
}

func (v *list_[V]) RemoveValue(index int) V {

	// Create a new smaller Array.
	var removed = v.GetValue(index)
	var size = v.GetSize() - 1
	var Array = Array[V]()
	var array = Array.WithSize(size)

	// Copy the remaining values into the new Array.
	var counter = v.toNormalized(index)
	index = 0
	var iterator = v.GetIterator()
	for iterator.HasNext() {
		counter--
		var value = iterator.GetNext()
		if counter == 0 {
			continue // Skip this value.
		}
		index++
		array.SetValue(index, value)
	}

	// Update the internal Array.
	v.values = array

	return removed
}

func (v *list_[V]) RemoveValues(first int, last int) Sequential[V] {

	// Create two smaller Arrays.
	first = v.toNormalized(first)
	last = v.toNormalized(last)
	var delta = last - first + 1
	var size = v.GetSize() - delta
	var Array = Array[V]()
	var removed = Array.WithSize(delta)
	var array = Array.WithSize(size)

	// Split the existing values into the two new Arrays.
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

	// Update the internal Array.
	v.values = array

	return removed
}

// Searchable Interface

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

func (v *list_[V]) ContainsValue(value V) bool {
	return v.GetIndex(value) > 0
}

func (v *list_[V]) GetComparer() ComparingFunction {
	return v.compare
}

func (v *list_[V]) GetIndex(value V) int {
	for index, candidate := range v.AsArray() {
		if v.compare(candidate, value) {
			// Found the value.
			return index + 1 // Convert to an ORDINAL based index.
		}
	}
	// The value was not found.
	return 0
}

// Sequential Interface

func (v *list_[V]) AsArray() []V {
	return v.values.AsArray()
}

func (v *list_[V]) GetIterator() Ratcheted[V] {
	return v.values.GetIterator()
}

func (v *list_[V]) GetSize() int {
	return v.values.GetSize()
}

func (v *list_[V]) IsEmpty() bool {
	return v.values.IsEmpty()
}

// Sortable Interface

func (v *list_[V]) ReverseValues() {
	v.values.ReverseValues()
}

func (v *list_[V]) ShuffleValues() {
	v.values.ShuffleValues()
}

func (v *list_[V]) SortValues() {
	v.values.SortValues()
}

func (v *list_[V]) SortValuesWithRanker(ranker RankingFunction) {
	v.values.SortValuesWithRanker(ranker)
}

// Updatable Interface

func (v *list_[V]) SetValue(index int, value V) {
	v.values.SetValue(index, value)
}

func (v *list_[V]) SetValues(index int, values Sequential[V]) {
	v.values.SetValues(index, values)
}

// Private Interface

// This public class method is used by Go to generate a canonical string for
// the list.
func (v *list_[V]) String() string {
	return CDCN().Default().FormatCollection(v)
}

// This private class method normalizes the specified index.  The following
// transformation is performed:
// [-size..-1] and [1..size] => [1..size]
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
