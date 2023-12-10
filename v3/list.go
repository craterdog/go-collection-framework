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
// constants, constructors and functions for the list class namespace.
type listClass_[V Value] struct {
	// This class defines no constants.
}

// This private constant defines a map to hold all the singleton references to
// the type specific list namespaces.
var listClassSingletons = map[string]any{}

// This public function returns the singleton reference to a type specific
// list namespace.  It also initializes any class constants as needed.
func List[V Value]() *listClass_[V] {
	var class *listClass_[V]
	var key = fmt.Sprintf("%T", class)
	var value = listClassSingletons[key]
	switch actual := value.(type) {
	case *listClass_[V]:
		class = actual
	default:
		class = &listClass_[V]{
			// This class defines no constants.
		}
		listClassSingletons[key] = class
	}
	return class
}

// CLASS CONSTRUCTORS

// This public class constructor creates a new empty list.
// The list uses the natural comparing function.
func (c *listClass_[V]) FromNothing() ListLike[V] {
	var Array = Array[V]() // Retrieve the array namespace.
	var array = Array.FromSize(0)
	var compare = CompareValues
	var list = &list_[V]{array, compare}
	return list
}

// This public class constructor creates a new empty list that uses the
// specified comparing function.
func (c *listClass_[V]) FromComparer(compare ComparingFunction) ListLike[V] {
	var Array = Array[V]() // Retrieve the array namespace.
	var array = Array.FromSize(0)
	var list = &list_[V]{array, compare}
	return list
}

// This public class constructor creates a new list from the specified sequence.
// The list uses the natural compare function.
func (c *listClass_[V]) FromSequence(sequence Sequential[V]) ListLike[V] {
	var list = c.FromNothing()
	var iterator = sequence.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		list.AppendValue(value)
	}
	return list
}

// CLASS FUNCTIONS

// This public class function returns the concatenation of the two specified
// lists.
func (c *listClass_[V]) Concatenate(first, second ListLike[V]) ListLike[V] {
	var list = c.FromNothing()
	list.AppendValues(first)
	list.AppendValues(second)
	return list
}

// CLASS TYPE

// Encapsulated Type

// This private class type encapsulates a Go structure containing private
// attributes that can only be accessed and manipulated using methods that
// implement the list-like abstract type.  Each value in a list is associated
// with an implicit positive integer index. The list uses ORDINAL based indexing
// rather than ZERO based indexing (see the description of what this means in
// the sequential interface definition).  The comparison of values in the list
// use a configurable comparison function.
// This type is parameterized as follows:
//   - V is any type of value.
type list_[V Value] struct {
	ArrayLike[V]
	compare ComparingFunction
}

// Searchable Interface

// This public class method returns the comparing function for this list.
func (v *list_[V]) GetComparer() ComparingFunction {
	return v.compare
}

// This method returns the index of the FIRST occurrence of the specified value in
// this list, or zero if this list does not contain the value.
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

// This method determines whether or not this list contains the specified value.
func (v *list_[V]) ContainsValue(value V) bool {
	return v.GetIndex(value) > 0
}

// This method determines whether or not this list contains ANY of the specified
// values.
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

// This method determines whether or not this list contains ALL of the specified
// values.
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

// Expandable Interface

// This method appends the specified value to the end of this list.
func (v *list_[V]) AppendValue(value V) {

	// Create a new bigger array.
	var size = v.GetSize() + 1
	var Array = Array[V]()
	var array = Array.FromSize(size)

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
	v.ArrayLike = array
}

// This method appends the specified values to the end of this list.
func (v *list_[V]) AppendValues(values Sequential[V]) {

	// Create a new bigger array.
	var size = v.GetSize() + values.GetSize()
	var Array = Array[V]()
	var array = Array.FromSize(size)

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
	v.ArrayLike = array
}

// This method inserts the specified value into this list in the specified
// slot between the existing values.
func (v *list_[V]) InsertValue(slot int, value V) {

	// Create a new larger array.
	var size = v.GetSize() + 1
	var Array = Array[V]()
	var array = Array.FromSize(size)

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
	v.ArrayLike = array
}

// This method inserts the specified values into this list in the specified
// slot between existing values.
func (v *list_[V]) InsertValues(slot int, values Sequential[V]) {

	// Create a new bigger array.
	var size = v.GetSize() + values.GetSize()
	var Array = Array[V]()
	var array = Array.FromSize(size)

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
	v.ArrayLike = array
}

// This method removes the value at the specified index from this list. The
// removed value is returned.
func (v *list_[V]) RemoveValue(index int) V {

	// Create a new smaller array.
	var removed = v.GetValue(index)
	var size = v.GetSize() - 1
	var Array = Array[V]()
	var array = Array.FromSize(size)

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
	v.ArrayLike = array

	return removed
}

// This method removes the values in the specified index range from this list.
// The removed values are returned.
func (v *list_[V]) RemoveValues(first int, last int) Sequential[V] {

	// Create two smaller arrays.
	first = v.normalized(first)
	last = v.normalized(last)
	var delta = last - first + 1
	var size = v.GetSize() - delta
	var Array = Array[V]()
	var removed = Array.FromSize(delta)
	var array = Array.FromSize(size)

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
	v.ArrayLike = array

	return removed
}

// This method removes all values from this list.
func (v *list_[V]) RemoveAll() {
	var Array = Array[V]()
	v.ArrayLike = Array.FromSize(0)
}

// Private Interface

// This private class method is used by Go to generate a canonical string for
// the list.
func (v *list_[V]) String() string {
	return FormatCollection(v)
}

// This private class method normalizes the specified index.  The following
// transformation is performed:
// [-size..-1] and [1..size] => [1..size]
func (v *list_[V]) normalized(index int) int {
	var size = v.GetSize()
	switch {
	case size == 0:
		// The list is empty.
		panic("Cannot index an empty list.")
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
