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
// constants, constructors and functions for the List class namespace.
type listClass_[V Value] struct {
	// This class defines no constants.
}

// This private constant defines a map to hold all the singleton references to
// the type specific List class namespaces.
var listClassSingletons = map[string]any{}

// This public function returns the singleton reference to a type specific
// List class namespace.  It also initializes any class constants as needed.
func List[V Value]() *listClass_[V] {
	var class *listClass_[V]
	var key = fmt.Sprintf("%T", class) // The name of the bound class type.
	var value = listClassSingletons[key]
	switch actual := value.(type) {
	case *listClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Create a new bound class type.
		class = &listClass_[V]{
			// This class defines no constants.
		}
		listClassSingletons[key] = class
	}
	return class
}

// CLASS CONSTRUCTORS

// This public class constructor creates a new empty List.
// The List uses the natural comparing function.
func (c *listClass_[V]) Empty() ListLike[V] {
	var compare = Collator().CompareValues
	var list = c.WithComparer(compare)
	return list
}

// This public class constructor creates a new List from the specified sequence.
// The List uses the natural compare function.
func (c *listClass_[V]) FromSequence(sequence Sequential[V]) ListLike[V] {
	var list = c.Empty()
	var iterator = sequence.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		list.AppendValue(value)
	}
	return list
}

// This public class constructor creates a new empty List that uses the
// specified comparing function.
func (c *listClass_[V]) WithComparer(compare ComparingFunction) ListLike[V] {
	var Array = Array[V]() // Retrieve the array namespace.
	var values = Array.WithSize(0)
	var list = &list_[V]{
		compare: compare,
		values:  values,
	}
	return list
}

// CLASS FUNCTIONS

// This public class function returns the concatenation of the two specified
// lists.
func (c *listClass_[V]) Concatenate(first, second ListLike[V]) ListLike[V] {
	var list = c.Empty()
	list.AppendValues(first)
	list.AppendValues(second)
	return list
}

// CLASS TYPE

// Encapsulated Type

// This private class type encapsulates a Go structure containing private
// attributes that can only be accessed and manipulated using methods that
// implement the list-like abstract type.  Each value in a List is associated
// with an implicit positive integer index. The List uses ORDINAL based indexing
// rather than ZERO based indexing (see the description of what this means in
// the sequential interface definition).  The comparison of values in the List
// use a configurable comparison function.
// This type is parameterized as follows:
//   - V is any type of value.
type list_[V Value] struct {
	compare ComparingFunction
	values  ArrayLike[V]
}

// Accessible Interface

// This public class method retrieves from this array the value that is
// associated with the specified index.
func (v *list_[V]) GetValue(index int) V {
	return v.values.GetValue(index)
}

// This public class method retrieves from this array all values from the first
// index through the last index (inclusive).
func (v *list_[V]) GetValues(first int, last int) Sequential[V] {
	return v.values.GetValues(first, last)
}

// Expandable Interface

// This public class method appends the specified value to the end of this List.
func (v *list_[V]) AppendValue(value V) {

	// Create a new bigger array.
	var size = v.GetSize() + 1
	var Array = Array[V]()
	var array = Array.WithSize(size)

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
	v.values = array
}

// This public class method appends the specified values to the end of this
// List.
func (v *list_[V]) AppendValues(values Sequential[V]) {

	// Create a new bigger array.
	var size = v.GetSize() + values.GetSize()
	var Array = Array[V]()
	var array = Array.WithSize(size)

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
	v.values = array
}

// This public class method inserts the specified value into this List in the
// specified slot between the existing values.
func (v *list_[V]) InsertValue(slot int, value V) {

	// Create a new larger array.
	var size = v.GetSize() + 1
	var Array = Array[V]()
	var array = Array.WithSize(size)

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
	v.values = array
}

// This public class method inserts the specified values into this List in the
// specified slot between existing values.
func (v *list_[V]) InsertValues(slot int, values Sequential[V]) {

	// Create a new bigger array.
	var size = v.GetSize() + values.GetSize()
	var Array = Array[V]()
	var array = Array.WithSize(size)

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
	v.values = array
}

// This public class method removes all values from this List.
func (v *list_[V]) RemoveAll() {
	var Array = Array[V]()
	v.values = Array.WithSize(0)
}

// This public class method removes the value at the specified index from this
// List. The removed value is returned.
func (v *list_[V]) RemoveValue(index int) V {

	// Create a new smaller array.
	var removed = v.GetValue(index)
	var size = v.GetSize() - 1
	var Array = Array[V]()
	var array = Array.WithSize(size)

	// Copy the remaining values into the new array.
	var counter = v.normalized(index)
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

	// Update the internal array.
	v.values = array

	return removed
}

// This public class method removes the values in the specified index range from
// this List.  The removed values are returned.
func (v *list_[V]) RemoveValues(first int, last int) Sequential[V] {

	// Create two smaller arrays.
	first = v.normalized(first)
	last = v.normalized(last)
	var delta = last - first + 1
	var size = v.GetSize() - delta
	var Array = Array[V]()
	var removed = Array.WithSize(delta)
	var array = Array.WithSize(size)

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
	v.values = array

	return removed
}

// Searchable Interface

// This public class method determines whether or not this List contains ALL of
// the specified values.
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

// This public class method determines whether or not this List contains ANY of
// the specified values.
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

// This public class method determines whether or not this List contains the
// specified value.
func (v *list_[V]) ContainsValue(value V) bool {
	return v.GetIndex(value) > 0
}

// This public class method returns the comparing function for this List.
func (v *list_[V]) GetComparer() ComparingFunction {
	return v.compare
}

// This public class method returns the index of the FIRST occurrence of the
// specified value in this list, or zero if this List does not contain the
// value.
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

// This public class method returns all the values in this array. The values
// retrieved are in the same order as they are in the array.
func (v *list_[V]) AsArray() []V {
	return v.values.AsArray()
}

// This public class method generates for this array an iterator that can be
// used to traverse its values.
func (v *list_[V]) GetIterator() Ratcheted[V] {
	return v.values.GetIterator()
}

// This public class method returns the number of values contained in this
// array.
func (v *list_[V]) GetSize() int {
	return v.values.GetSize()
}

// This public class method determines whether or not this array is empty.
func (v *list_[V]) IsEmpty() bool {
	return v.values.IsEmpty()
}

// Sortable Interface

// This public class method reverses the order of all values in this List.
func (v *list_[V]) ReverseValues() {
	v.values.ReverseValues()
}

// This public class method pseudo-randomly shuffles the values in this List.
func (v *list_[V]) ShuffleValues() {
	v.values.ShuffleValues()
}

// This public class method sorts the values in this List using the natural
// ranking function.
func (v *list_[V]) SortValues() {
	v.values.SortValues()
}

// This public class method sorts the values in this List using the specified
// ranking function.
func (v *list_[V]) SortValuesWithRanker(ranker RankingFunction) {
	v.values.SortValuesWithRanker(ranker)
}

// Updatable Interface

// This public class method sets the value in this array that is associated
// with the specified index to be the specified value.
func (v *list_[V]) SetValue(index int, value V) {
	v.values.SetValue(index, value)
}

// This public class method sets the values in this array starting with the
// specified index to the specified values.
func (v *list_[V]) SetValues(index int, values Sequential[V]) {
	v.values.SetValues(index, values)
}

// Private Interface

// This public class method is used by Go to generate a canonical string for
// the List.
func (v *list_[V]) String() string {
	return Formatter().FormatCollection(v)
}

// This private class method normalizes the specified index.  The following
// transformation is performed:
// [-size..-1] and [1..size] => [1..size]
func (v *list_[V]) normalized(index int) int {
	var size = v.GetSize()
	switch {
	case size == 0:
		// The List is empty.
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
