/*
................................................................................
.    Copyright (c) 2009-2024 Crater Dog Technologies.  All Rights Reserved.    .
................................................................................
.  DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               .
.                                                                              .
.  This code is free software; you can redistribute it and/or modify it under  .
.  the terms of The MIT License (MIT), as published by the Open Source         .
.  Initiative. (See http://opensource.org/licenses/MIT)                        .
................................................................................
*/

package collections

import (
	fmt "fmt"
)

// ARRAY IMPLEMENTATION

// This type defines the structure and methods associated with a native array of
// values. Each value is associated with an implicit positive integer index. The
// array uses ORDINAL based indexing rather than ZERO based indexing (see the
// description of what this means in the Sequential interface definition).
// This type is parameterized as follows:
//   - V is any type of value.
type Array[V Value] []V

// This constructor creates a new array from the specified sequence of values.
func ArrayFromSequence[V Value](sequence Sequential[V]) []V {
	var v = sequence.AsArray()
	return v
}

// SEQUENTIAL INTERFACE

// This method determines whether or not this array is empty.
func (v Array[V]) IsEmpty() bool {
	return len(v) == 0
}

// This method returns the number of values contained in this array.
func (v Array[V]) GetSize() int {
	return len(v)
}

// This method returns all the values in this array. The values retrieved are in
// the same order as they are in the array.
func (v Array[V]) AsArray() []V {
	var length = len(v)
	var result = make([]V, length)
	copy(result, v)
	return result
}

// ACCESSIBLE INTERFACE

// This method retrieves from this array the value that is associated with the
// specified index.
func (v Array[V]) GetValue(index int) V {
	index = v.GoIndex(index)
	return v[index]
}

// This method retrieves from this array all values from the first index through
// the last index (inclusive).
func (v Array[V]) GetValues(first int, last int) Sequential[V] {
	first = v.GoIndex(first)
	last = v.GoIndex(last)
	var result = Array[V](v[first : last+1])
	return result
}

// UPDATABLE INTERFACE

// This method sets the value in this array that is associated with the specified
// index to be the specified value.
func (v Array[V]) SetValue(index int, value V) {
	index = v.GoIndex(index)
	v[index] = value
}

// This method sets the values in this array starting with the specified index
// to the specified values.
func (v Array[V]) SetValues(index int, values Sequential[V]) {
	index = v.GoIndex(index)
	copy(v[index:], values.AsArray())
}

// GO INTERFACE

// This method is used by Go to generate a string from an array.
func (v Array[V]) String() string {
	return FormatCollection(v)
}

// This method normalizes an index to match the Go (zero based) indexing. The
// following transformation is performed:
//
//	[-length..-1] and [1..length] => [0..length)
//
// Notice that the specified index cannot be zero since zero is not an ORDINAL.
func (v Array[V]) GoIndex(index int) int {
	var length = len(v)
	switch {
	case length == 0:
		// The array is empty.
		panic("Cannot index an empty array.")
	case index == 0:
		// Zero is not an ordinal.
		panic("Indices must be positive or negative ordinals, not zero.")
	case index < -length || index > length:
		// The index is outside the bounds of the specified range.
		panic(fmt.Sprintf(
			"The specified index is outside the allowed ranges [-%v..-1] and [1..%v]: %v",
			length,
			length,
			index))
	case index < 0:
		// Convert a negative index.
		return index + length
	case index > 0:
		// Convert a positive index.
		return index - 1
	default:
		// This should never happen so time to...
		panic(fmt.Sprintf("Go compiler problem, unexpected index value: %v", index))
	}
}
