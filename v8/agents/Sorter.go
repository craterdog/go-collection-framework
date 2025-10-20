/*
................................................................................
.    Copyright (c) 2009-2025 Crater Dog Technologies.  All Rights Reserved.    .
................................................................................
.  DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               .
.                                                                              .
.  This code is free software; you can redistribute it and/or modify it under  .
.  the terms of The MIT License (MIT), as published by the Open Source         .
.  Initiative. (See https://opensource.org/license/MIT)                        .
................................................................................
*/

package agents

import (
	ran "crypto/rand"
	fmt "fmt"
	uti "github.com/craterdog/go-missing-utilities/v8"
	big "math/big"
	syn "sync"
)

// CLASS INTERFACE

// Access Function

func SorterClass[V any]() SorterClassLike[V] {
	return sorterClass[V]()
}

// Constructor Methods

func (c *sorterClass_[V]) Sorter() SorterLike[V] {
	var instance = &sorter_[V]{
		// Initialize the instance attributes.
		ranker_: CollatorClass[V]().Collator().RankValues,
	}
	return instance
}

func (c *sorterClass_[V]) SorterWithRanker(
	ranker RankingFunction[V],
) SorterLike[V] {
	if uti.IsUndefined(ranker) {
		panic("The \"ranker\" attribute is required by this class.")
	}
	var instance = &sorter_[V]{
		// Initialize the instance attributes.
		ranker_: ranker,
	}
	return instance
}

// Constant Methods

// Function Methods

// INSTANCE INTERFACE

// Principal Methods

func (v *sorter_[V]) GetClass() SorterClassLike[V] {
	return sorterClass[V]()
}

func (v *sorter_[V]) SortValues(
	values []V,
) {
	// Sort the values in place using a merge sort.
	v.sortValues(values)
}

func (v *sorter_[V]) ReverseValues(
	values []V,
) {
	// Reverse the values in place using Go's multi-assignment capability.
	var length = len(values)
	var half = length / 2 // Rounds down to the nearest integer.
	for index := 0; index < half; index++ {
		values[index], values[length-index-1] = values[length-index-1], values[index]
	}
}

func (v *sorter_[V]) ShuffleValues(
	values []V,
) {
	// Shuffle the values in place using random index exchanges.
	var size = len(values)
	for i := 0; i < size; i++ {
		var r = v.randomizeIndex(size)
		values[i], values[r] = values[r], values[i]
	}
}

// Attribute Methods

func (v *sorter_[V]) GetRanker() RankingFunction[V] {
	return v.ranker_
}

// PROTECTED INTERFACE

// Private Methods

func (v *sorter_[V]) randomizeIndex(
	size int,
) int {
	// Generate a cryptographically secure random index in the range [0..size).
	var random, err = ran.Int(ran.Reader, big.NewInt(int64(size)))
	if err != nil {
		// There was an issue with the underlying OS so time to...
		panic("Unable to generate a random index:\n" + err.Error())
	}
	return int(random.Int64())
}

// NOTE:
// This method, and the mergeArrays method, together sort the values in the
// specified Go array in place using an iterative merge sort along with the
// ranking function associated with this sorter.  The algorithm is documented
// here:
//   - https://en.wikipedia.org/wiki/Merge_sort#Bottom-up_implementation
//
// This iterative approach saves on memory allocation by swapping between
// two Go arrays of the same size rather than allocating new Go arrays for
// each sub-array.  This results in stable O[nlog(n)] time and O[n] space
// performance.
func (v *sorter_[V]) sortValues(
	values []V,
) {
	// Create a buffer Go array.
	var length = len(values)
	var buffer = make([]V, length)
	copy(buffer, values) // Make a copy of the original unsorted Go array.

	// Iterate through sub-array widths of 2, 4, 8, ... length.
	for width := 1; width < length; width *= 2 {

		// Split the buffer Go array into two Go arrays.
		for left := 0; left < length; left += width * 2 {

			// Find the middle (it must be less than length).
			var middle = left + width
			if middle > length {
				middle = length
			}

			// Find the right side (it must be less than length).
			var right = middle + width
			if right > length {
				right = length
			}

			// Sort and merge the sub-arrays.
			v.mergeArrays(
				buffer[left:middle],
				buffer[middle:right],
				values[left:right])
		}

		// Swap the two Go arrays.
		buffer, values = values, buffer
	}

	// Synchronize the two Go arrays.
	copy(values, buffer) // Both Go arrays are now sorted.
}

func (v *sorter_[V]) mergeArrays(
	left []V,
	right []V,
	merged []V,
) {
	var leftIndex = 0
	var leftLength = len(left)
	var rightIndex = 0
	var rightLength = len(right)
	var mergedIndex = 0
	var mergedLength = len(merged)

	// Work our way through filling the entire merged Go array.
	for mergedIndex < mergedLength {

		// Check to see if both left and right Go arrays still have values.
		if leftIndex < leftLength && rightIndex < rightLength {

			// Copy the next smallest value to the merged Go array.
			if v.ranker_(left[leftIndex], right[rightIndex]) == LesserRank {
				merged[mergedIndex] = left[leftIndex]
				leftIndex++
			} else {
				merged[mergedIndex] = right[rightIndex]
				rightIndex++
			}

		} else if leftIndex < leftLength {
			// Copy the rest of the left Go array to the merged Go array.
			copy(merged[mergedIndex:], left[leftIndex:])
			leftIndex++

		} else {
			// Copy the rest of the right Go array to the merged Go array.
			copy(merged[mergedIndex:], right[rightIndex:])
			rightIndex++
		}
		mergedIndex++
	}
}

// Instance Structure

type sorter_[V any] struct {
	// Declare the instance attributes.
	ranker_ RankingFunction[V]
}

// Class Structure

type sorterClass_[V any] struct {
	// Declare the class constants.
}

// Class Reference

var sorterMap_ = map[string]any{}
var sorterMutex_ syn.Mutex

func sorterClass[V any]() *sorterClass_[V] {
	// Generate the name of the bound class type.
	var class *sorterClass_[V]
	var name = fmt.Sprintf("%T", class)

	// Check for an existing bound class type.
	sorterMutex_.Lock()
	var value = sorterMap_[name]
	switch actual := value.(type) {
	case *sorterClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Add a new bound class type.
		class = &sorterClass_[V]{
			// Initialize the class constants.
		}
		sorterMap_[name] = class
	}
	sorterMutex_.Unlock()

	// Return a reference to the bound class type.
	return class
}
