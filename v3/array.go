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

// This private type defines the namespace structure associated with the constants,
// constructors and functions for the array class namespace.
type arrayClass_[V Value] struct {
	// This class defines no constants.
}

// This private constant defines a map to hold all the singleton references to
// the type specific array namespaces.
var arrayClassSingletons = map[string]any{}

// This public function returns the singleton reference to a type specific
// array namespace.  It also initializes any class constants as needed.
func Array[V Value]() *arrayClass_[V] {
	var class *arrayClass_[V]
	var key = fmt.Sprintf("%T", class)
	var value = arrayClassSingletons[key]
	switch actual := value.(type) {
	case *arrayClass_[V]:
		class = actual
	default:
		class = &arrayClass_[V]{
			// This class defines no constants.
		}
		arrayClassSingletons[key] = class
	}
	return class
}

// CLASS CONSTRUCTORS

// This public class constructor creates a new array of the specified size.
func (c *arrayClass_[V]) WithSize(size int) ArrayLike[V] {
	var array = make([]V, size) // All values initialized to zero.
	return array_[V](array)
}

// This public class constructor creates a new array from the specified
// Go array of values.
func (c *arrayClass_[V]) FromArray(array []V) ArrayLike[V] {
	var length = len(array)
	var duplicate = make([]V, length)
	copy(duplicate, array)
	return array_[V](duplicate)
}

// CLASS FUNCTIONS

// This public class function normalizes an index into a sequence to match the
// Go (zero based) indexing. The following transformation is performed:
//
//	[-size..-1] and [1..size] => [0..size)
//
// Notice that the specified index cannot be zero since zero is not an ORDINAL.
func (c *arrayClass_[V]) GoIndex(sequence Sequential[V], index int) int {
	var size = sequence.GetSize()
	switch {
	case size == 0:
		// The array is empty.
		panic("Cannot index an empty array.")
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
		return index + size
	case index > 0:
		// Convert a positive index.
		return index - 1
	default:
		// This should never happen so time to...
		panic(fmt.Sprintf("Go compiler problem, unexpected index value: %v", index))
	}
}

// CLASS TYPE

// Extended Type

// This private class type extends the primitive Go array data type and defines
// the methods that implement the array-like abstract type.  Each value is
// associated with an implicit positive integer index. The array uses ORDINAL
// based indexing rather than ZERO based indexing (see the description of what
// this means in the Sequential interface definition).
// This type is parameterized as follows:
//   - V is any type of value.
type array_[V Value] []V

// Sequential Interface

// This public class method determines whether or not this array is empty.
func (v array_[V]) IsEmpty() bool {
	return len(v) == 0
}

// This public class method returns the number of values contained in this
// array.
func (v array_[V]) GetSize() int {
	return len(v)
}

// This public class method returns all the values in this array. The values
// retrieved are in the same order as they are in the array.
func (v array_[V]) AsArray() []V {
	var length = len(v)
	var array = make([]V, length)
	copy(array, v)
	return array
}

// This public class method generates for this array an iterator that can be
// used to traverse its values.
func (v array_[V]) GetIterator() Ratcheted[V] {
	var Iterator = Iterator[V]()
	var iterator = Iterator.FromSequence(v)
	return iterator
}

// Accessible Interface

// This public class method retrieves from this array the value that is
// associated with the specified index.
func (v array_[V]) GetValue(index int) V {
	var Array = Array[V]()
	index = Array.GoIndex(v, index)
	return v[index]
}

// This public class method retrieves from this array all values from the first
// index through the last index (inclusive).
func (v array_[V]) GetValues(first int, last int) Sequential[V] {
	var Array = Array[V]()
	first = Array.GoIndex(v, first)
	last = Array.GoIndex(v, last)
	var sequence = v[first : last+1]
	var array = Array.FromArray(sequence) // This copies the underlying array.
	return array
}

// Updatable Interface

// This public class method sets the value in this array that is associated
// with the specified index to be the specified value.
func (v array_[V]) SetValue(index int, value V) {
	var Array = Array[V]()
	index = Array.GoIndex(v, index)
	v[index] = value
}

// This public class method sets the values in this array starting with the
// specified index to the specified values.
func (v array_[V]) SetValues(index int, values Sequential[V]) {
	var Array = Array[V]()
	index = Array.GoIndex(v, index)
	copy(v[index:], values.AsArray())
}

// Sortable Interface

// This public class method sorts the values in this list using the natural
// ranking function.
func (v array_[V]) SortValues() {
	v.SortValuesWithRanker(RankValues)
}

// This public class method sorts the values in this list using the specified
// ranking function.
func (v array_[V]) SortValuesWithRanker(ranker RankingFunction) {
	if v.GetSize() > 1 {
		SortValues(v, ranker)
	}
}

// This public class method reverses the order of all values in this list.
func (v array_[V]) ReverseValues() {
	ReverseValues(v)
}

// This public class method pseudo-randomly shuffles the values in this list.
func (v array_[V]) ShuffleValues() {
	ShuffleValues(v)
}

// Go Interface

// This public class method is used by Go to generate a string from an array.
func (v array_[V]) String() string {
	return FormatCollection(v)
}
