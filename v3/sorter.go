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
	ran "crypto/rand"
	fmt "fmt"
	big "math/big"
)

// CLASS NAMESPACE

// This private type defines the namespace structure associated with the
// constants, constructors and functions for the Sorter class namespace.
type sorterClass_[V Value] struct {
	defaultRanker RankingFunction
}

// This private constant defines a map to hold all the singleton references to
// the type specific Sorter class namespaces.
var sorterClassSingletons = map[string]any{}

// This public function returns the singleton reference to a type specific
// Sorter class namespace.  It also initializes any class constants as needed.
func Sorter[V Value]() *sorterClass_[V] {
	var class *sorterClass_[V]
	var key = fmt.Sprintf("%T", class) // The name of the bound class type.
	var value = sorterClassSingletons[key]
	switch actual := value.(type) {
	case *sorterClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Create a new bound class type.
		class = &sorterClass_[V]{
			defaultRanker: Collator().RankValues,
		}
		sorterClassSingletons[key] = class
	}
	return class
}

// CLASS CONSTANTS

// This public class constant represents the default ranking function.
func (c *sorterClass_[V]) DefaultRanker() RankingFunction {
	return c.defaultRanker
}

// CLASS CONSTRUCTORS

// This public class constructor creates a new Sorter that can be used to sort
// an array using the specified ranking function.
func (c *sorterClass_[V]) WithDefaultRanker() SorterLike[V] {
	var sorter = &sorter_[V]{
		rank: c.defaultRanker,
	}
	return sorter
}

// This public class constructor creates a new Sorter that can be used to sort
// an array using the specified ranking function.
func (c *sorterClass_[V]) WithSpecifiedRanker(ranker RankingFunction) SorterLike[V] {
	var sorter = &sorter_[V]{
		rank: ranker,
	}
	return sorter
}

// CLASS FUNCTIONS

// This public class function reverses the order of the values in the specified
// array.
func (c *sorterClass_[V]) ReverseValues(array []V) {
	var v = c.WithDefaultRanker()
	v.ReverseValues(array)
}

// This public class function randomly shuffles the values in the specified
// array.
func (c *sorterClass_[V]) ShuffleValues(array []V) {
	var v = c.WithDefaultRanker()
	v.ShuffleValues(array)
}

// This public class function sorts the values in the specified array using the
// specified ranking function.
func (c *sorterClass_[V]) SortValues(array []V, ranker RankingFunction) {
	var v = c.WithSpecifiedRanker(ranker)
	v.SortValues(array)
}

// CLASS TYPE

// Encapsulated Type

// This private class type encapsulates a Go structure containing private
// attributes that can only be accessed and manipulated using methods that
// implement the sorter-like abstract type.
type sorter_[V Value] struct {
	rank RankingFunction
}

// Systematic Interface

// This public class method reverses the order of the values in the specified
// array in place.
func (v *sorter_[V]) ReverseValues(array []V) {
	var length = len(array)
	var half = length / 2 // Rounds down to the nearest integer.
	for index := 0; index < half; index++ {
		array[index], array[length-index-1] = array[length-index-1], array[index]
	}
}

// This public class method randomly shuffles the values in the specified array
// in place.
func (v *sorter_[V]) ShuffleValues(array []V) {
	var size = len(array)
	for i := 0; i < size; i++ {
		var r = v.randomIndex(size)
		array[i], array[r] = array[r], array[i]
	}
}

// This public class method sorts the values in the specified array in place
// using an iterative merge sort along with the ranking function associated with
// this Sorter.  The algorithm is documented here:
//   - https://en.wikipedia.org/wiki/Merge_sort#Bottom-up_implementation
//
// This iterative approach saves on memory allocation by swapping between two
// arrays of the same size rather than allocating new arrays for each sub-array.
// This results in stable O[nlog(n)] time and O[n] space performance.
func (v *sorter_[V]) SortValues(array []V) {
	// Create a buffer array.
	var length = len(array)
	var buffer = make([]V, length)
	copy(buffer, array) // Make a copy of the original unsorted array.

	// Iterate through sub-array widths of 2, 4, 8, ... length.
	for width := 1; width < length; width *= 2 {

		// Split the buffer array into two arrays.
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
				array[left:right])
		}

		// Swap the two arrays.
		buffer, array = array, buffer
	}

	// Synchronize the two arrays.
	copy(array, buffer) // Both arrays are now sorted.
}

// Private Interface

// This private class method is used for the merging part of the merge sort
// algorithm.
func (v *sorter_[V]) mergeArrays(left []V, right []V, merged []V) {
	var leftIndex = 0
	var leftLength = len(left)
	var rightIndex = 0
	var rightLength = len(right)
	var mergedIndex = 0
	var mergedLength = len(merged)

	// Work our way through filling the entire merged array.
	for mergedIndex < mergedLength {

		// Check to see if both left and right arrays still have values.
		if leftIndex < leftLength && rightIndex < rightLength {

			// Copy the next smallest value to the merged array.
			if v.rank(left[leftIndex], right[rightIndex]) < 0 {
				merged[mergedIndex] = left[leftIndex]
				leftIndex++
			} else {
				merged[mergedIndex] = right[rightIndex]
				rightIndex++
			}

		} else if leftIndex < leftLength {
			// Copy the rest of the left array to the merged array.
			copy(merged[mergedIndex:], left[leftIndex:])
			leftIndex++

		} else {
			// Copy the rest of the right array to the merged array.
			copy(merged[mergedIndex:], right[rightIndex:])
			rightIndex++
		}
		mergedIndex++
	}
}

// This private class method generates a cryptographically secure random index
// in the range [0..size).
func (v *sorter_[V]) randomIndex(size int) int {
	var random, err = ran.Int(ran.Reader, big.NewInt(int64(size)))
	if err != nil {
		// There was an issue with the underlying OS so time to...
		panic("Unable to generate a random index:\n" + err.Error())
	}
	return int(random.Int64())
}
