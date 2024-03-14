/*
................................................................................
.    Copyright (c) 2009-2024 Crater Dog Technologies™.  All Rights Reserved.   .
................................................................................
.  DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               .
.                                                                              .
.  This code is free software; you can redistribute it and/or modify it under  .
.  the terms of The MIT License (MIT), as published by the Open Source         .
.  Initiative. (See http://opensource.org/licenses/MIT)                        .
................................................................................
*/

package collections_test

import (
	col "github.com/craterdog/go-collection-framework/v2"
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

func TestEmptyMaps(t *tes.T) {
	var m = col.Map[string, int](map[string]int{})
	ass.True(t, m.IsEmpty())
	ass.Equal(t, 0, m.GetSize())
	ass.Equal(t, []string{}, m.GetKeys().AsArray())
	ass.Equal(t, []col.Binding[string, int]{}, m.AsArray())
	var iterator = col.Iterator[col.Binding[string, int]](m)
	ass.False(t, iterator.HasNext())
	ass.False(t, iterator.HasPrevious())
	iterator.ToStart()
	iterator.ToEnd()
	m.RemoveAll()
}

func TestMapsWithStringsAndIntegers(t *tes.T) {
	var keys = col.ListFromArray([]string{"foo", "bar"})
	var association1 = col.Association[string, int]("foo", 1)
	var association2 = col.Association[string, int]("bar", 2)
	var association3 = col.Association[string, int]("baz", 3)
	var m = col.Map[string, int](map[string]int{})
	m.SetValue(association1.GetKey(), association1.GetValue())
	ass.False(t, m.IsEmpty())
	ass.Equal(t, 1, m.GetSize())
	m.SetValue(association2.GetKey(), association2.GetValue())
	m.SetValue(association3.GetKey(), association3.GetValue())
	ass.Equal(t, 3, m.GetSize())
	ass.Equal(t, 3, int(m.GetValue("baz")))
	m.SetValue("bar", 5)
	ass.Equal(t, []int{1, 5}, m.GetValues(keys).AsArray())
	ass.Equal(t, []int{1, 5}, m.RemoveValues(keys).AsArray())
	ass.Equal(t, 1, m.GetSize())
	ass.Equal(t, 3, int(m.RemoveValue("baz")))
	ass.True(t, m.IsEmpty())
	ass.Equal(t, 0, m.GetSize())
}
