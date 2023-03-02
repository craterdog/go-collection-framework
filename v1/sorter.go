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
	big "math/big"
)

// SORTER INTERFACE

// This function reverses the order of the values in the specified array.
func ReverseArray[V Value](array []V) {
	var v = Sorter[V](nil)
	v.ReverseArray(array)
}

// This function randomly shuffles the values in the specified array.
func ShuffleArray[V Value](array []V) {
	var v = Sorter[V](nil)
	v.ShuffleArray(array)
}

// This function sorts the values in the specified array using the specified
// ranking function.
func SortArray[V Value](array []V, rank RankingFunction) {
	var v = Sorter[V](rank)
	v.SortArray(array)
}

// This constructor creates a new instance of a sorter that can be used to
// sort an array using a specific ranking function.
func Sorter[V Value](rank RankingFunction) SorterLike[V] {
	return &sorter[V]{rank: rank}
}

// SORTER IMPLEMENTATION

// This type defines the structure and methods for a merge sorter agent.
type sorter[V Value] struct {
	rank RankingFunction
}

// This method reverses the order of the values in the specified array in place.
func (v *sorter[V]) ReverseArray(array []V) {
	v.reverseArray(array)
}

// This method randomly shuffles the values in the specified array in place.
func (v *sorter[V]) ShuffleArray(array []V) {
	v.shuffleArray(array)
}

// This method sorts the values in the specified array in place.
func (v *sorter[V]) SortArray(array []V) {
	v.sortArray(array)
}

// This method reverses the the order of the values in the specified array.
func (v *sorter[V]) reverseArray(array []V) {
	var length = len(array)
	var half = length / 2 // Rounds down to the nearest integer.
	for index := 0; index < half; index++ {
		array[index], array[length-index-1] = array[length-index-1], array[index]
		//var first = array[index]
		//var last = array[length - index - 1]
		//array[index] = last
		//array[length - index - 1] = first
	}
}

// This method pseudo-randomly shuffles the values in the specified array.
func (v *sorter[V]) shuffleArray(array []V) {
	var size = len(array)
	for i := 0; i < size; i++ {
		var r = randomIndex(size)
		array[i], array[r] = array[r], array[i]
	}
}

// This function generates a cryptographically secure random index in the
// range [0..size).
func randomIndex(size int) int {
	var random, err = ran.Int(ran.Reader, big.NewInt(int64(size)))
	if err != nil {
		// There was an issue with the underlying OS so time to...
		panic("Unable to generate a random index:\n" + err.Error())
	}
	return int(random.Int64())
}

// The following methods are used to sort an array using the ranking function
// that is associated with the sorter. These methods implement a merge sort
// that has been optimized to be iterative rather than recursive. They also
// save on memory allocation by swapping between two arrays of the same size
// rather than allocating new arrays for each subarray.  This results in stable
// O[nlog(n)] time and O[n] space performance. The algorithm is documented here:
//   - https://en.wikipedia.org/wiki/Merge_sort#Bottom-up_implementation

func (v *sorter[V]) sortArray(array []V) {
	// Create a buffer array.
	var length = len(array)
	var buffer = make([]V, length)
	copy(buffer, array) // Make a copy of the original unsorted array.

	// Iterate through subarray widths of 2, 4, 8, ... length.
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

			// Sort and merge the subarrays.
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

func (v *sorter[V]) mergeArrays(left []V, right []V, merged []V) {
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
