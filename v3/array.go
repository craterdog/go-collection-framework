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

type arrayClass_[V Value] struct {
	// This class defines no constants.
}

// Private Namespace Reference(s)

var arrayClass = map[string]any{}

// Public Namespace Access

func Array[V Value]() ArrayClassLike[V] {
	var class *arrayClass_[V]
	var key = fmt.Sprintf("%T", class) // The name of the bound class type.
	var value = arrayClass[key]
	switch actual := value.(type) {
	case *arrayClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Create a new bound class type.
		class = &arrayClass_[V]{
			// This class defines no constants.
		}
		arrayClass[key] = class
	}
	return class
}

// Public Class Constructors

func (c *arrayClass_[V]) FromArray(values []V) ArrayLike[V] {
	var size = len(values)
	var array = make([]V, size)
	copy(array, values)
	return array_[V](array)
}

func (c *arrayClass_[V]) FromSequence(values Sequential[V]) ArrayLike[V] {
	var size = values.GetSize()
	var iterator = values.GetIterator()
	var array = make([]V, size)
	for index := 0; index < size; index++ {
		var value = iterator.GetNext()
		array[index] = value
	}
	return array_[V](array)
}

func (c *arrayClass_[V]) FromString(values string) ArrayLike[V] {
	// First we parse it as a collection of any type value.
	var collection = CDCN().Default().ParseCollection(values).(Sequential[Value])

	// Then we convert it to an Array of type V.
	var array = c.WithSize(collection.GetSize())
	var index int
	var iterator = collection.GetIterator()
	for iterator.HasNext() {
		index++
		var value = iterator.GetNext().(V)
		array.SetValue(index, value)
	}
	return array
}

func (c *arrayClass_[V]) WithSize(size int) ArrayLike[V] {
	var array = make([]V, size) // All values initialized to zero.
	return array_[V](array)
}

// CLASS TYPE

// Private Class Type Definition

type array_[V Value] []V

// Accessible Interface

func (v array_[V]) GetValue(index int) V {
	index = v.toZeroBased(index)
	return v[index]
}

func (v array_[V]) GetValues(first int, last int) Sequential[V] {
	first = v.toZeroBased(first)
	last = v.toZeroBased(last)
	var sequence = v[first : last+1]
	var array = Array[V]().FromArray(sequence) // This copies the underlying Go array.
	return array
}

// Sequential Interface

func (v array_[V]) AsArray() []V {
	var length = len(v)
	var array = make([]V, length)
	copy(array, v)
	return array
}

func (v array_[V]) GetIterator() Ratcheted[V] {
	var Iterator = Iterator[V]()
	var iterator = Iterator.FromSequence(v)
	return iterator
}

func (v array_[V]) GetSize() int {
	return len(v)
}

func (v array_[V]) IsEmpty() bool {
	return len(v) == 0
}

// Sortable Interface

func (v array_[V]) ReverseValues() {
	var Sorter = Sorter[V]().Default()
	Sorter.ReverseValues(v)
}

func (v array_[V]) ShuffleValues() {
	var Sorter = Sorter[V]().Default()
	Sorter.ShuffleValues(v)
}

func (v array_[V]) SortValues() {
	var ranker = Collator().Default().RankValues
	v.SortValuesWithRanker(ranker)
}

func (v array_[V]) SortValuesWithRanker(ranker RankingFunction) {
	if v.GetSize() > 1 {
		var Sorter = Sorter[V]().WithRanker(ranker)
		Sorter.SortValues(v)
	}
}

// Updatable Interface

func (v array_[V]) SetValue(index int, value V) {
	index = v.toZeroBased(index)
	v[index] = value
}

func (v array_[V]) SetValues(index int, values Sequential[V]) {
	// The full index range must be in bounds.
	var size = values.GetSize()
	var first = v.toZeroBased(index)
	var last = v.toZeroBased(index+size-1) + 1
	copy(v[first:last], values.AsArray())
}

// Private Interface

// This public class method is used by Go to generate a string from an Array.
func (v array_[V]) String() string {
	return CDCN().Default().FormatCollection(v)
}

// This private class method normalizes a relative ORDINAL-based index into this
// Array to match the Go (ZERO-based) indexing. The following transformation is
// performed:
//
//	[-size..-1] and [1..size] => [0..size)
//
// Notice that the specified index cannot be zero since zero is NOT an ordinal.
func (v array_[V]) toZeroBased(index int) int {
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
