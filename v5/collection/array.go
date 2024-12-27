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

func ArrayClass[V any]() ArrayClassLike[V] {
	return arrayClass[V]()
}

// Constructor Methods

func (c *arrayClass_[V]) Array(
	values []V,
) ArrayLike[V] {
	var array = uti.CopyArray(values)
	return array_[V](array)
}

func (c *arrayClass_[V]) ArrayWithSize(
	size age.Size,
) ArrayLike[V] {
	var array = make([]V, int(size)) // All values initialized to zero.
	return array_[V](array)
}

func (c *arrayClass_[V]) ArrayFromSequence(
	values Sequential[V],
) ArrayLike[V] {
	var array = values.AsArray() // This returns a copy of the array.
	return array_[V](array)
}

// Constant Methods

// Function Methods

// INSTANCE INTERFACE

// Principal Methods

func (v array_[V]) GetClass() ArrayClassLike[V] {
	return arrayClass[V]()
}

// Attribute Methods

// Accessible[V] Methods

func (v array_[V]) GetValue(
	index Index,
) V {
	var goIndex = v.toZeroBased(index)
	return v[goIndex]
}

func (v array_[V]) GetValues(
	first Index,
	last Index,
) Sequential[V] {
	var goFirst = v.toZeroBased(first)
	var goLast = v.toZeroBased(last)
	var values = v[goFirst : goLast+1]
	var array = uti.CopyArray(values)
	return array_[V](array)
}

// Sequential[V] Methods

func (v array_[V]) IsEmpty() bool {
	return len(v) == 0
}

func (v array_[V]) GetSize() age.Size {
	var size = age.Size(len(v))
	return size
}

func (v array_[V]) AsArray() []V {
	var array = uti.CopyArray(v)
	return array
}

func (v array_[V]) GetIterator() age.IteratorLike[V] {
	var array = uti.CopyArray(v)
	var iteratorClass = age.IteratorClass[V]()
	var iterator = iteratorClass.Iterator(array)
	return iterator
}

// Sortable[V] Methods

func (v array_[V]) SortValues() {
	if v.GetSize() > 1 {
		var sorterClass = age.SorterClass[V]()
		var sorter = sorterClass.Sorter()
		sorter.SortValues(v)
	}
}

func (v array_[V]) SortValuesWithRanker(
	ranker age.RankingFunction[V],
) {
	if v.GetSize() > 1 {
		var sorterClass = age.SorterClass[V]()
		var sorter = sorterClass.SorterWithRanker(ranker)
		sorter.SortValues(v)
	}
}

func (v array_[V]) ReverseValues() {
	if v.GetSize() > 1 {
		var sorterClass = age.SorterClass[V]()
		var sorter = sorterClass.Sorter()
		sorter.ReverseValues(v)
	}
}

func (v array_[V]) ShuffleValues() {
	if v.GetSize() > 1 {
		var sorterClass = age.SorterClass[V]()
		var sorter = sorterClass.Sorter()
		sorter.ShuffleValues(v)
	}
}

// Updatable[V] Methods

func (v array_[V]) SetValue(
	index Index,
	value V,
) {
	var goIndex = v.toZeroBased(index)
	v[goIndex] = value
}

func (v array_[V]) SetValues(
	index Index,
	values Sequential[V],
) {
	var goIndex = v.toZeroBased(index)
	copy(v[goIndex:], values.AsArray())
}

// Stringer Methods

func (v array_[V]) String() string {
	return uti.Format(v)
}

// PROTECTED INTERFACE

// Private Methods

// This private instance method normalizes a relative ORDINAL-based index into
// this array to match the Go (ZERO-based) indexing.  The following
// transformation is performed:
//
//	[-size..-1] and [1..size] => [0..size)
//
// Notice that the specified index cannot be zero since zero is NOT an ordinal.
func (v array_[V]) toZeroBased(index Index) int {
	var size = Index(v.GetSize())
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

type array_[V any] []V

// Class Structure

type arrayClass_[V any] struct {
	// Declare the class constants.
}

// Class Reference

var arrayMap_ = map[string]any{}
var arrayMutex_ syn.Mutex

func arrayClass[V any]() *arrayClass_[V] {
	// Generate the name of the bound class type.
	var class *arrayClass_[V]
	var name = fmt.Sprintf("%T", class)

	// Check for an existing bound class type.
	arrayMutex_.Lock()
	var value = arrayMap_[name]
	switch actual := value.(type) {
	case *arrayClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Add a new bound class type.
		class = &arrayClass_[V]{
			// Initialize the class constants.
		}
		arrayMap_[name] = class
	}
	arrayMutex_.Unlock()

	// Return a reference to the bound class type.
	return class
}
