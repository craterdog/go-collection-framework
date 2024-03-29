/*
................................................................................
.    Copyright (c) 2009-2024 Crater Dog Technologies.  All Rights Reserved.    .
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
	col "github.com/craterdog/go-collection-framework/v3"
	not "github.com/craterdog/go-collection-framework/v3/cdcn"
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

func TestCatalogConstructors(t *tes.T) {
	var notation = not.Notation().Make()
	var Catalog = col.Catalog[rune, int64]()
	var _ = Catalog.MakeFromArray([]col.AssociationLike[rune, int64]{})
	var _ = Catalog.MakeFromSource("[:](Catalog)", notation)
	var _ = Catalog.MakeFromSource("['a': 1, 'b': 2, 'c': 3](Catalog)", notation)
	var _ = Catalog.MakeFromMap(map[rune]int64{})
	var _ = Catalog.MakeFromMap(map[rune]int64{
		'a': 1,
		'b': 2,
		'c': 3,
	})
}

func TestCatalogsWithStringsAndIntegers(t *tes.T) {
	var collator = col.Collator().Make()
	var keys = col.Array[string]().MakeFromArray([]string{"foo", "bar"})
	var Association = col.Association[string, int]()
	var association1 = Association.MakeWithAttributes("foo", 1)
	var association2 = Association.MakeWithAttributes("bar", 2)
	var association3 = Association.MakeWithAttributes("baz", 3)
	var Catalog = col.Catalog[string, int]()
	var catalog = Catalog.Make()
	ass.True(t, catalog.IsEmpty())
	ass.Equal(t, 0, catalog.GetSize())
	ass.Equal(t, []string{}, catalog.GetKeys().AsArray())
	ass.Equal(t, []col.AssociationLike[string, int]{}, catalog.AsArray())
	var iterator = catalog.GetIterator()
	ass.False(t, iterator.HasNext())
	ass.False(t, iterator.HasPrevious())
	iterator.ToStart()
	iterator.ToEnd()
	catalog.SortValues()
	catalog.ShuffleValues()
	catalog.RemoveAll()
	catalog.SetValue(association1.GetKey(), association1.GetValue())
	ass.False(t, catalog.IsEmpty())
	ass.Equal(t, 1, catalog.GetSize())
	catalog.SetValue(association2.GetKey(), association2.GetValue())
	catalog.SetValue(association3.GetKey(), association3.GetValue())
	ass.Equal(t, 3, catalog.GetSize())
	var catalog2 = Catalog.MakeFromSequence(catalog)
	ass.True(t, collator.CompareValues(catalog, catalog2))
	var m = col.Map[string, int]().MakeFromMap(map[string]int{
		"foo": 1,
		"bar": 2,
		"baz": 3,
	})
	var catalog3 = Catalog.MakeFromSequence(m)
	catalog2.SortValues()
	catalog3.SortValuesWithRanker(collator.RankValues)
	ass.True(t, collator.CompareValues(catalog2, catalog3))
	iterator = catalog.GetIterator()
	ass.True(t, iterator.HasNext())
	ass.False(t, iterator.HasPrevious())
	ass.Equal(t, association1, iterator.GetNext())
	ass.True(t, iterator.HasPrevious())
	iterator.ToEnd()
	ass.False(t, iterator.HasNext())
	ass.True(t, iterator.HasPrevious())
	ass.Equal(t, association3, iterator.GetPrevious())
	ass.True(t, iterator.HasNext())
	ass.Equal(t, []string{"foo", "bar", "baz"}, catalog.GetKeys().AsArray())
	ass.Equal(t, 3, int(catalog.GetValue("baz")))
	catalog.SetValue("bar", 5)
	ass.Equal(t, []int{1, 5}, catalog.GetValues(keys).AsArray())
	catalog.SortValues()
	ass.Equal(t, []string{"bar", "baz", "foo"}, catalog.GetKeys().AsArray())
	catalog.ReverseValues()
	ass.Equal(t, []string{"foo", "baz", "bar"}, catalog.GetKeys().AsArray())
	catalog.ReverseValues()
	ass.Equal(t, []int{1, 5}, catalog.RemoveValues(keys).AsArray())
	ass.Equal(t, 1, catalog.GetSize())
	ass.Equal(t, 3, int(catalog.RemoveValue("baz")))
	ass.True(t, catalog.IsEmpty())
	ass.Equal(t, 0, catalog.GetSize())
	catalog.RemoveAll()
	ass.True(t, catalog.IsEmpty())
	ass.Equal(t, 0, catalog.GetSize())
}

func TestCatalogsWithMerge(t *tes.T) {
	var collator = col.Collator().Make()
	var Association = col.Association[string, int]()
	var association1 = Association.MakeWithAttributes("foo", 1)
	var association2 = Association.MakeWithAttributes("bar", 2)
	var association3 = Association.MakeWithAttributes("baz", 3)
	var Catalog = col.Catalog[string, int]()
	var catalog1 = Catalog.Make()
	catalog1.SetValue(association1.GetKey(), association1.GetValue())
	catalog1.SetValue(association2.GetKey(), association2.GetValue())
	var catalog2 = Catalog.Make()
	catalog2.SetValue(association2.GetKey(), association2.GetValue())
	catalog2.SetValue(association3.GetKey(), association3.GetValue())
	var catalog3 = Catalog.Merge(catalog1, catalog2)
	var catalog4 = Catalog.Make()
	catalog4.SetValue(association1.GetKey(), association1.GetValue())
	catalog4.SetValue(association2.GetKey(), association2.GetValue())
	catalog4.SetValue(association3.GetKey(), association3.GetValue())
	ass.True(t, collator.CompareValues(catalog3, catalog4))
}

func TestCatalogsWithExtract(t *tes.T) {
	var collator = col.Collator().Make()
	var keys = col.Array[string]().MakeFromArray([]string{"foo", "baz"})
	var Association = col.Association[string, int]()
	var association1 = Association.MakeWithAttributes("foo", 1)
	var association2 = Association.MakeWithAttributes("bar", 2)
	var association3 = Association.MakeWithAttributes("baz", 3)
	var Catalog = col.Catalog[string, int]()
	var catalog1 = Catalog.Make()
	catalog1.SetValue(association1.GetKey(), association1.GetValue())
	catalog1.SetValue(association2.GetKey(), association2.GetValue())
	catalog1.SetValue(association3.GetKey(), association3.GetValue())
	var catalog2 = Catalog.Extract(catalog1, keys)
	var catalog3 = Catalog.Make()
	catalog3.SetValue(association1.GetKey(), association1.GetValue())
	catalog3.SetValue(association3.GetKey(), association3.GetValue())
	ass.True(t, collator.CompareValues(catalog2, catalog3))
	var catalog4 = Catalog.MakeFromArray([]col.AssociationLike[string, int]{
		association1,
		association2,
		association3,
	})
	ass.True(t, collator.CompareValues(catalog1, catalog4))
}

func TestCatalogsWithEmptyCatalogs(t *tes.T) {
	var collator = col.Collator().Make()
	var keys = col.Array[int]().MakeFromSize(0)
	var Catalog = col.Catalog[int, string]()
	var catalog1 = Catalog.Make()
	var catalog2 = Catalog.Make()
	var catalog3 = Catalog.Merge(catalog1, catalog2)
	var catalog4 = Catalog.Extract(catalog1, keys)
	ass.True(t, collator.CompareValues(catalog1, catalog2))
	ass.True(t, collator.CompareValues(catalog2, catalog3))
	ass.True(t, collator.CompareValues(catalog3, catalog4))
	ass.True(t, collator.CompareValues(catalog4, catalog1))
}
