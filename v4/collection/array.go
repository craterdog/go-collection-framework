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
	sts "strings"
	syn "sync"
)

// CLASS ACCESS

// Reference

var arrayClass = map[string]any{}
var arrayMutex syn.Mutex

// Function

func Array[V any](notation NotationLike) ArrayClassLike[V] {
	// Validate the notation argument.
	if notation == nil {
		panic("A notation must be specified when creating this class.")
	}

	// Generate the name of the bound class type.
	var class ArrayClassLike[V]
	var name = fmt.Sprintf("%T-%T", class, notation)

	// Check for existing bound class type.
	arrayMutex.Lock()
	var value = arrayClass[name]
	switch actual := value.(type) {
	case *arrayClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Add a new bound class type.
		class = &arrayClass_[V]{
			notation_: notation,
		}
		arrayClass[name] = class
	}
	arrayMutex.Unlock()

	// Return a reference to the bound class type.
	return class
}

// CLASS METHODS

// Target

type arrayClass_[V any] struct {
	notation_ NotationLike
}

// Constants

func (c *arrayClass_[V]) Notation() NotationLike {
	return c.notation_
}

// Constructors

func (c *arrayClass_[V]) MakeWithSize(size uint) ArrayLike[V] {
	var array = make([]V, size) // All values initialized to zero.
	return array_[V](array)
}

func (c *arrayClass_[V]) MakeFromArray(values []V) ArrayLike[V] {
	var size = len(values)
	var array = make([]V, size)
	copy(array, values)
	return array_[V](array)
}

func (c *arrayClass_[V]) MakeFromSequence(values Sequential[V]) ArrayLike[V] {
	var size = values.GetSize()
	var iterator = values.GetIterator()
	var array = make([]V, size)
	for index := 0; index < size; index++ {
		var value = iterator.GetNext()
		array[index] = value
	}
	return array_[V](array)
}

func (c *arrayClass_[V]) MakeFromSource(source string) ArrayLike[V] {
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

// INSTANCE METHODS

// Target

type array_[V any] []V

// Accessible

func (v array_[V]) GetValue(index int) V {
	index = v.toZeroBased(index)
	return v[index]
}

func (v array_[V]) GetValues(first int, last int) Sequential[V] {
	first = v.toZeroBased(first)
	last = v.toZeroBased(last)
	var values = v[first : last+1]
	var size = last - first + 1
	var array = make([]V, size)
	copy(array, values)
	return array_[V](array)
}

// Sequential

func (v array_[V]) IsEmpty() bool {
	return len(v) == 0
}

func (v array_[V]) GetSize() int {
	return len(v)
}

func (v array_[V]) AsArray() []V {
	var size = len(v)
	var array = make([]V, size)
	copy(array, v)
	return array
}

func (v array_[V]) GetIterator() age.IteratorLike[V] {
	var iterator = age.Iterator[V]().MakeFromArray(v.AsArray())
	return iterator
}

// Sortable

func (v array_[V]) SortValues() {
	var collator = age.Collator[V]().Make()
	var ranker = collator.RankValues
	v.SortValuesWithRanker(ranker)
}

func (v array_[V]) SortValuesWithRanker(ranker age.RankingFunction[V]) {
	if v.GetSize() > 1 {
		var sorter = age.Sorter[V]().MakeWithRanker(ranker)
		sorter.SortValues(v)
	}
}

func (v array_[V]) ReverseValues() {
	var sorter = age.Sorter[V]().Make()
	sorter.ReverseValues(v)
}

func (v array_[V]) ShuffleValues() {
	var sorter = age.Sorter[V]().Make()
	sorter.ShuffleValues(v)
}

// Updatable

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

// Stringer

// NOTE:
// Since this class only extends the primitive array type it cannot have any
// attributes assigned to it.  This means that we have no way of accessing its
// notation.  So we cannot use the notation specific formatter to generate the
// string value for this array and must generate it manually.  This is only a
// problem when this method is called directly—as done by the fmt.Sprintf()
// method.  The formatters themselves can handle the formatting of arrays just
// fine.
func (v array_[V]) String() string {
	var string_ = "["
	if v.IsEmpty() {
		string_ += (" ")
	} else {
		var builder sts.Builder
		var iterator = v.GetIterator()
		for iterator.HasNext() {
			var value = iterator.GetNext()
			builder.WriteString(fmt.Sprintf("%#v, ", value))
		}
		var last = builder.Len() - 2
		string_ += builder.String()[:last]
	}
	string_ += "](Array)\n"
	return string_
}

// Private

// This private instance method normalizes a relative ORDINAL-based index into
// this array to match the Go (ZERO-based) indexing.  The following
// transformation is performed:
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
