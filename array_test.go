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
	col "github.com/craterdog/go-collection-framework"
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

func TestEmptyArrays(t *tes.T) {
	var array = col.Array[string]([]string{})
	ass.True(t, array.IsEmpty())
	ass.Equal(t, 0, array.GetSize())
	ass.False(t, array.ContainsValue("bax"))
	ass.Equal(t, 0, array.GetIndex("baz"))
	ass.Equal(t, []string{}, array.AsArray())
	var iterator = col.Iterator[string](array)
	ass.False(t, iterator.HasNext())
	ass.False(t, iterator.HasPrevious())
	iterator.ToStart()
	iterator.ToEnd()
}

func TestArraysWithStrings(t *tes.T) {
	var array = col.Array[string]([]string{"foo", "bar", "baz"})
	var foobar = col.Array[string]([]string{"foo", "bar"})
	ass.False(t, array.IsEmpty())
	ass.Equal(t, 3, array.GetSize())
	ass.Equal(t, "foo", array.GetValue(1))
	ass.Equal(t, foobar, array.GetValues(1, 2))
	ass.Equal(t, 3, array.GetIndex("baz"))
	ass.False(t, array.ContainsValue("bax"))
	ass.True(t, array.ContainsValue("bar"))
	ass.True(t, array.ContainsAny(foobar))
	ass.True(t, array.ContainsAll(foobar))
	array.SetValue(2, "bax")
	ass.True(t, array.ContainsValue("bax"))
	ass.True(t, array.ContainsAny(foobar))
	ass.False(t, array.ContainsAll(foobar))
	ass.Equal(t, []string{"foo", "bax", "baz"}, array.AsArray())
	array.SetValues(2, foobar)
	ass.Equal(t, []string{"foo", "foo", "bar"}, array.AsArray())
}
