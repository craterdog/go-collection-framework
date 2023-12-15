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
// constants, constructors and functions for the Array class namespace.
type arrayClass_[V Value] struct {
	// This class defines no constants.
}

// This private constant defines a map to hold all the singleton references to
// the type specific Array namespaces.
var arrayClassSingletons = map[string]any{}

// This public function returns the singleton reference to a type specific
// Array class namespace.  It also initializes any class constants as needed.
func Array[V Value]() *arrayClass_[V] {
	var class *arrayClass_[V]
	var key = fmt.Sprintf("%T", class) // The name of the bound class type.
	var value = arrayClassSingletons[key]
	switch actual := value.(type) {
	case *arrayClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Create a new bound class type.
		class = &arrayClass_[V]{
			// This class defines no constants.
		}
		arrayClassSingletons[key] = class
	}
	return class
}

// CLASS CONSTRUCTORS

// This public class constructor creates a new Array from the specified
// Go array of values.
func (c *arrayClass_[V]) FromArray(array []V) ArrayLike[V] {
	var length = len(array)
	var duplicate = make([]V, length)
	copy(duplicate, array)
	return array_[V](duplicate)
}

// This public class constructor creates a new Array of the specified size.
func (c *arrayClass_[V]) WithSize(size int) ArrayLike[V] {
	var array = make([]V, size) // All values initialized to zero.
	return array_[V](array)
}

// CLASS TYPE

// Extended Type

// This private class type extends the primitive Go array data type and defines
// the methods that implement the array-like abstract type.  Each value is
// associated with an implicit positive integer index. The Array uses ORDINAL
// based indexing rather than ZERO based indexing (see the description of what
// this means in the Sequential interface definition).
// This type is parameterized as follows:
//   - V is any type of value.
type array_[V Value] []V

// Accessible Interface

// This public class method retrieves from this Array the value that is
// associated with the specified index.
func (v array_[V]) GetValue(index int) V {
	index = v.goIndex(index)
	return v[index]
}

// This public class method retrieves from this Array all values from the first
// index through the last index (inclusive).
func (v array_[V]) GetValues(first int, last int) Sequential[V] {
	first = v.goIndex(first)
	last = v.goIndex(last)
	var sequence = v[first : last+1]
	var array = Array[V]().FromArray(sequence) // This copies the underlying array.
	return array
}

// Sequential Interface

// This public class method returns all the values in this Array. The values
// retrieved are in the same order as they are in the Array.
func (v array_[V]) AsArray() []V {
	var length = len(v)
	var array = make([]V, length)
	copy(array, v)
	return array
}

// This public class method generates for this Array an iterator that can be
// used to traverse its values.
func (v array_[V]) GetIterator() Ratcheted[V] {
	var Iterator = Iterator[V]()
	var iterator = Iterator.FromSequence(v)
	return iterator
}

// This public class method returns the number of values contained in this
// Array.
func (v array_[V]) GetSize() int {
	return len(v)
}

// This public class method determines whether or not this Array is empty.
func (v array_[V]) IsEmpty() bool {
	return len(v) == 0
}

// Sortable Interface

// This public class method reverses the order of all values in this list.
func (v array_[V]) ReverseValues() {
	var Sorter = Sorter[V]()
	Sorter.ReverseValues(v)
}

// This public class method pseudo-randomly shuffles the values in this list.
func (v array_[V]) ShuffleValues() {
	var Sorter = Sorter[V]()
	Sorter.ShuffleValues(v)
}

// This public class method sorts the values in this list using the default
// ranking function.
func (v array_[V]) SortValues() {
	v.SortValuesWithRanker(Collator().RankValues)
}

// This public class method sorts the values in this list using the specified
// ranking function.
func (v array_[V]) SortValuesWithRanker(ranker RankingFunction) {
	if v.GetSize() > 1 {
		var Sorter = Sorter[V]()
		Sorter.SortValues(v, ranker)
	}
}

// Updatable Interface

// This public class method sets the value in this Array that is associated
// with the specified index to be the specified value.
func (v array_[V]) SetValue(index int, value V) {
	index = v.goIndex(index)
	v[index] = value
}

// This public class method sets the values in this Array starting with the
// specified index to the specified values.
func (v array_[V]) SetValues(index int, values Sequential[V]) {
	// The full index range must be in bounds.
	var size = values.GetSize()
	var first = v.goIndex(index)
	var last = v.goIndex(index+size-1) + 1
	copy(v[first:last], values.AsArray())
}

// Private Interface

// This public class method is used by Go to generate a string from an Array.
func (v array_[V]) String() string {
	return Formatter().FormatCollection(v)
}

// This private class method normalizes an index into this Array to match the
// Go (zero based) indexing. The following transformation is performed:
//
//	[-size..-1] and [1..size] => [0..size)
//
// Notice that the specified index cannot be zero since zero is not an ORDINAL.
func (v array_[V]) goIndex(index int) int {
	var size = v.GetSize()
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
		return index + size
	case index > 0:
		// Convert a positive index.
		return index - 1
	default:
		// This should never happen so time to...
		panic(fmt.Sprintf("Go compiler problem, unexpected index value: %v", index))
	}
}
