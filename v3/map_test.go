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
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

func TestMapConstructors(t *tes.T) {
	var Map = col.MapClass[rune, int64]()
	var _ = Map.FromArray([]col.Binding[rune, int64]{})
	var _ = Map.FromMap(map[rune]int64{})
	var sequence = Map.FromMap(map[rune]int64{'a': 1, 'b': 2, 'c': 3})
	var _ = Map.FromSequence(sequence)
	var _ = Map.FromString("[:](Map)")
	var _ = Map.FromString("['a': 1, 'b': 2, 'c': 3](Map)")
}

func TestEmptyMaps(t *tes.T) {
	var m = col.MapClass[string, int]().FromMap(map[string]int{})
	ass.True(t, m.IsEmpty())
	ass.Equal(t, 0, m.GetSize())
	ass.Equal(t, []string{}, m.GetKeys().AsArray())
	ass.Equal(t, []col.Binding[string, int]{}, m.AsArray())
	var iterator = m.GetIterator()
	ass.False(t, iterator.HasNext())
	ass.False(t, iterator.HasPrevious())
	iterator.ToStart()
	iterator.ToEnd()
	m.RemoveAll()
}

func TestMapsWithStringsAndIntegers(t *tes.T) {
	var Association = col.AssociationClass[string, int]()
	var association1 = Association.FromPair("foo", 1)
	var association2 = Association.FromPair("bar", 2)
	var association3 = Association.FromPair("baz", 3)
	var Map = col.MapClass[string, int]()
	var m = Map.FromArray([]col.Binding[string, int]{
		association1,
		association2,
		association3,
	})
	ass.Equal(t, 1, int(m.GetValue("foo")))
	ass.Equal(t, 2, int(m.GetValue("bar")))
	ass.Equal(t, 3, int(m.GetValue("baz")))
	m = Map.FromMap(map[string]int{})
	m.SetValue(association1.GetKey(), association1.GetValue())
	ass.False(t, m.IsEmpty())
	ass.Equal(t, 1, m.GetSize())
	m.SetValue(association2.GetKey(), association2.GetValue())
	m.SetValue(association3.GetKey(), association3.GetValue())
	ass.Equal(t, 3, m.GetSize())
	ass.Equal(t, 3, int(m.GetValue("baz")))
	m.SetValue("bar", 5)
	var keys = col.ArrayClass[string]().FromArray([]string{"foo", "bar"})
	ass.Equal(t, []int{1, 5}, m.GetValues(keys).AsArray())
	ass.Equal(t, []int{1, 5}, m.RemoveValues(keys).AsArray())
	ass.Equal(t, 1, m.GetSize())
	ass.Equal(t, 3, int(m.RemoveValue("baz")))
	ass.True(t, m.IsEmpty())
	ass.Equal(t, 0, m.GetSize())
}
