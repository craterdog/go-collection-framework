/*
................................................................................
.    Copyright (c) 2009-2024 Crater Dog Technologies.  All Rights Reserved.    .
................................................................................
.  DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               .
.                                                                              .
.  This code is free software; you can redistribute it and/or modify it under  .
.  the terms of The MIT License (MIT), as published by the Open Source         .
.  Initiative. (See https://opensource.org/license/MIT)                        .
................................................................................
*/

package collection_test

import (
	fmt "fmt"
	not "github.com/craterdog/go-collection-framework/v4/cdcn"
	col "github.com/craterdog/go-collection-framework/v4/collection"
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

func TestMapConstructors(t *tes.T) {
	var notation = not.Notation().Make()
	var Map = col.Map[rune, int64](notation)
	var empty = Map.MakeFromArray([]col.AssociationLike[rune, int64]{})
	ass.Equal(t, "[:](Map)\n", fmt.Sprintf("%v", empty))
	var _ = Map.MakeFromMap(map[rune]int64{})
	var sequence = Map.MakeFromMap(map[rune]int64{'a': 1, 'b': 2, 'c': 3})
	var _ = Map.MakeFromSequence(sequence)
	var _ = Map.MakeFromSource("['a': 1, 'b': 2, 'c': 3](Map)")
}

func TestEmptyMaps(t *tes.T) {
	var notation = not.Notation().Make()
	var m = col.Map[string, int](notation).MakeFromMap(map[string]int{})
	ass.True(t, m.IsEmpty())
	ass.Equal(t, 0, m.GetSize())
	ass.Equal(t, []string{}, m.GetKeys().AsArray())
	ass.Equal(t, []col.AssociationLike[string, int]{}, m.AsArray())
	var iterator = m.GetIterator()
	ass.False(t, iterator.HasNext())
	ass.False(t, iterator.HasPrevious())
	iterator.ToStart()
	iterator.ToEnd()
	m.RemoveAll()
}

func TestMapsWithStringsAndIntegers(t *tes.T) {
	var notation = not.Notation().Make()
	var Association = col.Association[string, int]()
	var association1 = Association.MakeWithAttributes("foo", 1)
	var association2 = Association.MakeWithAttributes("bar", 2)
	var association3 = Association.MakeWithAttributes("baz", 3)
	var Map = col.Map[string, int](notation)
	var m = Map.MakeFromArray([]col.AssociationLike[string, int]{
		association1,
		association2,
		association3,
	})
	ass.Equal(t, 1, int(m.GetValue("foo")))
	ass.Equal(t, 2, int(m.GetValue("bar")))
	ass.Equal(t, 3, int(m.GetValue("baz")))
	m = Map.MakeFromMap(map[string]int{})
	m.SetValue(association1.GetKey(), association1.GetValue())
	ass.False(t, m.IsEmpty())
	ass.Equal(t, 1, m.GetSize())
	m.SetValue(association2.GetKey(), association2.GetValue())
	m.SetValue(association3.GetKey(), association3.GetValue())
	ass.Equal(t, 3, m.GetSize())
	ass.Equal(t, 3, int(m.GetValue("baz")))
	m.SetValue("bar", 5)
	var keys = col.Array[string](notation).MakeFromArray([]string{"foo", "bar"})
	ass.Equal(t, []int{1, 5}, m.GetValues(keys).AsArray())
	ass.Equal(t, []int{1, 5}, m.RemoveValues(keys).AsArray())
	ass.Equal(t, 1, m.GetSize())
	ass.Equal(t, 3, int(m.RemoveValue("baz")))
	ass.True(t, m.IsEmpty())
	ass.Equal(t, 0, m.GetSize())
}
