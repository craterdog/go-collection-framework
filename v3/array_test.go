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

func TestEmptyArray(t *tes.T) {
	var Array = col.Array[string]()
	var array = Array.WithSize(0)
	ass.True(t, array.IsEmpty())
	ass.Equal(t, 0, array.GetSize())
	ass.Equal(t, []string{}, array.AsArray())
	var iterator = array.GetIterator()
	ass.False(t, iterator.HasNext())
	ass.False(t, iterator.HasPrevious())
	iterator.ToStart()
	iterator.ToEnd()
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "Cannot index an empty array.", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	Array.GoIndex(array, 0) // This should panic.
}

func TestArrayIndexOfZero(t *tes.T) {
	var Array = col.Array[int]()
	var array = Array.FromArray([]int{1, 2, 3})
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "Indices must be positive or negative ordinals, not zero.", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	Array.GoIndex(array, 0) // This should panic.
}

func TestArrayWithStrings(t *tes.T) {
	var Array = col.Array[string]()
	var array = Array.FromArray([]string{"foo", "bar", "baz"})
	var foobar = Array.FromArray([]string{"foo", "bar"})
	ass.False(t, array.IsEmpty())
	ass.Equal(t, 3, array.GetSize())
	ass.Equal(t, "foo", array.GetValue(1))
	ass.Equal(t, foobar, array.GetValues(1, 2))
	array.SetValue(2, "bax")
	ass.Equal(t, []string{"foo", "bax", "baz"}, array.AsArray())
	array.SetValues(2, foobar)
	ass.Equal(t, []string{"foo", "foo", "bar"}, array.AsArray())
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The specified index is outside the allowed ranges [-3..-1] and [1..3]: 4", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	Array.GoIndex(array, 4) // This should panic.
}

func TestArrayWithIntegers(t *tes.T) {
	var Array = col.Array[int]()
	var array = Array.FromArray([]int{1, 2, 3})
	for index, value := range array.AsArray() {
		ass.Equal(t, index, Array.GoIndex(array, value))
		ass.Equal(t, index, Array.GoIndex(array, value-4))
	}
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The specified index is outside the allowed ranges [-3..-1] and [1..3]: -4", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	Array.GoIndex(array, -4) // This should panic.
}
