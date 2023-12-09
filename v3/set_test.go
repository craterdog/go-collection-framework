/*******************************************************************************
 *   Copyright (c) 2009-2022 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package collections_test

import (
	col "github.com/craterdog/go-collection-framework/v3"
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

func TestSetsWithStrings(t *tes.T) {
	var Array = col.Array[string]()
	var empty = []string{}
	var bazbar = Array.FromArray([]string{"baz", "bar"})
	var bazfoo = Array.FromArray([]string{"baz", "foo"})
	var baxbaz = Array.FromArray([]string{"bax", "baz"})
	var baxbez = Array.FromArray([]string{"bax", "bez"})
	var barbaz = Array.FromArray([]string{"bar", "baz"})
	var bar = Array.FromArray([]string{"bar"})
	var set = col.Set[string]()                                   // [ ]
	ass.True(t, set.IsEmpty())                                    // [ ]
	ass.Equal(t, 0, set.GetSize())                                // [ ]
	ass.False(t, set.ContainsValue("bax"))                        // [ ]
	ass.Equal(t, empty, set.AsArray())                            // [ ]
	var iterator = set.GetIterator()                              // [ ]
	ass.False(t, iterator.HasNext())                              // [ ]
	ass.False(t, iterator.HasPrevious())                          // [ ]
	iterator.ToStart()                                            // [ ]
	iterator.ToEnd()                                              // [ ]
	set.RemoveAll()                                               // [ ]
	set.RemoveValue("foo")                                        // [ ]
	set.AddValue("foo")                                           // ["foo"]
	ass.False(t, set.IsEmpty())                                   // ["foo"]
	ass.Equal(t, 1, set.GetSize())                                // ["foo"]
	ass.Equal(t, "foo", string(set.GetValue(1)))                  // ["foo"]
	ass.Equal(t, 0, set.GetIndex("baz"))                          // ["foo"]
	ass.True(t, set.ContainsValue("foo"))                         // ["foo"]
	ass.False(t, set.ContainsValue("bax"))                        // ["foo"]
	set.AddValues(bazbar)                                         // ["bar", "baz", "foo"]
	ass.Equal(t, 3, set.GetSize())                                // ["bar", "baz", "foo"]
	ass.Equal(t, 2, set.GetIndex("baz"))                          // ["bar", "baz", "foo"]
	ass.Equal(t, "bar", string(set.GetValue(1)))                  // ["bar", "baz", "foo"]
	ass.Equal(t, bazfoo.AsArray(), set.GetValues(2, 3).AsArray()) // ["bar", "baz", "foo"]
	ass.Equal(t, bar.AsArray(), set.GetValues(1, 1).AsArray())    // ["bar", "baz", "foo"]
	var set2 = col.SetFromSequence[string](set)                   // ["bar", "baz", "foo"]
	ass.True(t, col.CompareValues(set, set2))                     // ["bar", "baz", "foo"]
	var array = Array.FromArray([]string{"foo", "bar", "baz"})    // ["bar", "baz", "foo"]
	var set3 = col.SetFromSequence[string](array)                 // ["bar", "baz", "foo"]
	ass.True(t, col.CompareValues(set2, set3))                    // ["bar", "baz", "foo"]
	iterator = set.GetIterator()                                  // ["bar", "baz", "foo"]
	ass.True(t, iterator.HasNext())                               // ["bar", "baz", "foo"]
	ass.False(t, iterator.HasPrevious())                          // ["bar", "baz", "foo"]
	ass.Equal(t, "bar", string(iterator.GetNext()))               // ["bar", "baz", "foo"]
	ass.True(t, iterator.HasPrevious())                           // ["bar", "baz", "foo"]
	iterator.ToEnd()                                              // ["bar", "baz", "foo"]
	ass.False(t, iterator.HasNext())                              // ["bar", "baz", "foo"]
	ass.True(t, iterator.HasPrevious())                           // ["bar", "baz", "foo"]
	ass.Equal(t, "foo", string(iterator.GetPrevious()))           // ["bar", "baz", "foo"]
	ass.True(t, iterator.HasNext())                               // ["bar", "baz", "foo"]
	ass.True(t, set.ContainsValue("baz"))                         // ["bar", "baz", "foo"]
	ass.False(t, set.ContainsValue("bax"))                        // ["bar", "baz", "foo"]
	ass.True(t, set.ContainsAny(baxbaz))                          // ["bar", "baz", "foo"]
	ass.False(t, set.ContainsAny(baxbez))                         // ["bar", "baz", "foo"]
	ass.True(t, set.ContainsAll(barbaz))                          // ["bar", "baz", "foo"]
	ass.False(t, set.ContainsAll(baxbaz))                         // ["bar", "baz", "foo"]
	set.RemoveAll()                                               // [ ]
	ass.True(t, set.IsEmpty())                                    // [ ]
	ass.Equal(t, 0, set.GetSize())                                // [ ]
}

func TestSetsWithIntegers(t *tes.T) {
	var Array = col.Array[int]()
	var array = Array.FromArray([]int{3, 1, 4, 5, 9, 2})
	var set = col.Set[int]()               // [ ]
	set.AddValues(array)                   // [1,2,3,4,5,9]
	ass.False(t, set.IsEmpty())            // [1,2,3,4,5,9]
	ass.Equal(t, 6, set.GetSize())         // [1,2,3,4,5,9]
	ass.Equal(t, 1, int(set.GetValue(1)))  // [1,2,3,4,5,9]
	ass.Equal(t, 9, int(set.GetValue(-1))) // [1,2,3,4,5,9]
	set.RemoveValue(6)                     // [1,2,3,4,5,9]
	ass.Equal(t, 6, set.GetSize())         // [1,2,3,4,5,9]
	set.RemoveValue(3)                     // [1,2,4,5,9]
	ass.Equal(t, 5, set.GetSize())         // [1,2,4,5,9]
	ass.Equal(t, 4, int(set.GetValue(3)))  // [1,2,4,5,9]
}

func TestSetsWithTildas(t *tes.T) {
	var Array = col.Array[Integer]()
	var array = Array.FromArray([]Integer{3, 1, 4, 5, 9, 2})
	var set = col.Set[Integer]()           // [ ]
	set.AddValues(array)                   // [1,2,3,4,5,9]
	ass.False(t, set.IsEmpty())            // [1,2,3,4,5,9]
	ass.Equal(t, 6, set.GetSize())         // [1,2,3,4,5,9]
	ass.Equal(t, 1, int(set.GetValue(1)))  // [1,2,3,4,5,9]
	ass.Equal(t, 9, int(set.GetValue(-1))) // [1,2,3,4,5,9]
	set.RemoveValue(6)                     // [1,2,3,4,5,9]
	ass.Equal(t, 6, set.GetSize())         // [1,2,3,4,5,9]
	set.RemoveValue(3)                     // [1,2,4,5,9]
	ass.Equal(t, 5, set.GetSize())         // [1,2,4,5,9]
	ass.Equal(t, 4, int(set.GetValue(3)))  // [1,2,4,5,9]
}

func TestSetsWithSets(t *tes.T) {
	var Array = col.Array[int]()
	var array1 = Array.FromArray([]int{3, 1, 4, 5, 9, 2})
	var array2 = Array.FromArray([]int{7, 1, 4, 5, 9, 2})
	var set1 = col.Set[int]()
	set1.AddValues(array1)
	var set2 = col.Set[int]()
	set2.AddValues(array2)
	var set = col.Set[col.SetLike[int]]()
	var Array2 = col.Array[col.SetLike[int]]()
	var array3 = Array2.FromArray([]col.SetLike[int]{set1, set2})
	set.AddValues(array3)
	ass.False(t, set.IsEmpty())
	ass.Equal(t, 2, set.GetSize())
	ass.Equal(t, set1, set.GetValue(1))
	ass.Equal(t, set2, set.GetValue(-1))
	set.RemoveValue(set1)
	ass.Equal(t, 1, set.GetSize())
	set.RemoveAll()
	ass.Equal(t, 0, set.GetSize())
}

func TestSetWithNot(t *tes.T) {
	var set = col.Set[int]()
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "Not(set) is meaningless, use Sans(fullSet, set) instead.", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	col.Not(set) // This should panic.
}

func TestSetsWithAnd(t *tes.T) {
	var Array = col.Array[int]()
	var array1 = Array.FromArray([]int{3, 1, 2})
	var array2 = Array.FromArray([]int{3, 2, 4})
	var array3 = Array.FromArray([]int{3, 2})
	var set1 = col.Set[int]()
	set1.AddValues(array1)
	var set2 = col.Set[int]()
	set2.AddValues(array2)
	var set3 = col.And(set1, set2)
	var set4 = col.Set[int]()
	set4.AddValues(array3)
	ass.True(t, col.CompareValues(set3, set4))
}

func TestSetsWithSans(t *tes.T) {
	var Array = col.Array[int]()
	var array1 = Array.FromArray([]int{3, 1, 2})
	var array2 = Array.FromArray([]int{3, 2, 4})
	var array3 = Array.FromArray([]int{1})
	var set1 = col.Set[int]()
	set1.AddValues(array1)
	var set2 = col.Set[int]()
	set2.AddValues(array2)
	var set3 = col.Sans(set1, set2)
	var set4 = col.Set[int]()
	set4.AddValues(array3)
	ass.True(t, col.CompareValues(set3, set4))
}

func TestSetsWithOr(t *tes.T) {
	var Array = col.Array[int]()
	var array1 = Array.FromArray([]int{3, 1, 5})
	var array2 = Array.FromArray([]int{6, 2, 4})
	var array3 = Array.FromArray([]int{1, 3, 5, 6, 2, 4})
	var set1 = col.Set[int]()
	set1.AddValues(array1)
	var set2 = col.Set[int]()
	set2.AddValues(array2)
	var set3 = col.Or(set1, set2)
	ass.True(t, set3.ContainsAll(set1))
	ass.True(t, set3.ContainsAll(set2))
	var set4 = col.Set[int]()
	set4.AddValues(array3)
	ass.True(t, col.CompareValues(set3, set4))
}

func TestSetsWithXor(t *tes.T) {
	var Array = col.Array[int]()
	var array1 = Array.FromArray([]int{2, 3, 1, 5})
	var array2 = Array.FromArray([]int{6, 2, 5, 4})
	var array3 = Array.FromArray([]int{1, 3, 4, 6})
	var set1 = col.Set[int]()
	set1.AddValues(array1)
	var set2 = col.Set[int]()
	set2.AddValues(array2)
	var set3 = col.Xor(set1, set2)
	var set4 = col.Set[int]()
	set4.AddValues(array3)
	ass.True(t, col.CompareValues(set3, set4))
}

func TestSetsWithEmptySets(t *tes.T) {
	var set1 = col.Set[int]()
	var set2 = col.Set[int]()
	var set3 = col.And(set1, set2)
	var set4 = col.Sans(set1, set2)
	var set5 = col.Or(set1, set2)
	var set6 = col.Xor(set1, set2)
	ass.True(t, col.CompareValues(set3, set4))
	ass.True(t, col.CompareValues(set4, set5))
	ass.True(t, col.CompareValues(set5, set6))
	ass.True(t, col.CompareValues(set6, set1))
}
