/*******************************************************************************
 *   Copyright (c) 2009-2024 Crater Dog Technologies™.  All Rights Reserved.   *
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
	not "github.com/craterdog/go-collection-framework/v3/cdcn"
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

func TestArrayConstructors(t *tes.T) {
	var notation = not.Notation().Make()
	var Array = col.Array[int64]()
	var sequence = Array.MakeFromArray([]int64{1, 2, 3})
	var _ = Array.MakeFromSequence(sequence)
	var _ = Array.MakeFromSource("[ ](Array)", notation)
	var _ = Array.MakeFromSource("[1, 2, 3](Array)", notation)
}

func TestEmptyArray(t *tes.T) {
	var array = col.Array[string]().MakeFromSize(0)
	ass.True(t, array.IsEmpty())
	ass.Equal(t, 0, array.GetSize())
	ass.Equal(t, []string{}, array.AsArray())
	array.SortValues()
	var iterator = array.GetIterator()
	ass.False(t, iterator.HasNext())
	ass.False(t, iterator.HasPrevious())
	iterator.ToStart()
	iterator.ToEnd()
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "Cannot index an empty Array.", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	array.GetValue(1) // This should panic.
}

func TestArrayWithSize(t *tes.T) {
	var array = col.Array[string]().MakeFromSize(3)
	ass.False(t, array.IsEmpty())
	ass.Equal(t, 3, array.GetSize())
	ass.Equal(t, []string{"", "", ""}, array.AsArray())
	var iterator = array.GetIterator()
	ass.True(t, iterator.HasNext())
	ass.False(t, iterator.HasPrevious())
	iterator.ToStart()
	iterator.ToEnd()
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The specified index is outside the allowed ranges [-3..-1] and [1..3]: 4", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	array.SetValues(2, array) // This should panic.
}

func TestArrayIndexOfZero(t *tes.T) {
	var array = col.Array[int]().MakeFromArray([]int{1, 2, 3})
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "Indices must be positive or negative ordinals, not zero.", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	array.GetValue(0) // This should panic.
}

func TestArrayWithStrings(t *tes.T) {
	var collator = col.Collator().Make()
	var Array = col.Array[string]()
	var array = Array.MakeFromArray([]string{"foo", "bar", "baz"})
	var foobar = Array.MakeFromArray([]string{"foo", "bar"})
	ass.False(t, array.IsEmpty())
	ass.Equal(t, 3, array.GetSize())
	ass.Equal(t, "foo", array.GetValue(1))
	ass.Equal(t, foobar, array.GetValues(1, 2))
	array.SetValue(2, "bax")
	array.ShuffleValues()
	array.SortValuesWithRanker(collator.RankValues)
	ass.Equal(t, []string{"bax", "baz", "foo"}, array.AsArray())
	array.SetValues(2, foobar)
	ass.Equal(t, []string{"bax", "foo", "bar"}, array.AsArray())
	array.ReverseValues()
	ass.Equal(t, []string{"bar", "foo", "bax"}, array.AsArray())
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The specified index is outside the allowed ranges [-3..-1] and [1..3]: 4", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	array.SetValue(4, "bax") // This should panic.
}

func TestArrayWithIntegers(t *tes.T) {
	var array = col.Array[int]().MakeFromArray([]int{1, 2, 3})
	for index, value := range array.AsArray() {
		ass.Equal(t, index, array.GetValue(value)-1)
		ass.Equal(t, index, array.GetValue(value-4)-1)
	}
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The specified index is outside the allowed ranges [-3..-1] and [1..3]: -4", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	array.GetValue(-4) // This should panic.
}
