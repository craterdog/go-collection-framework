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

func TestCatalogsWithStringsAndIntegers(t *tes.T) {
	var keys = col.ListFromArray([]string{"foo", "bar"})
	var association1 = col.Association[string, int]("foo", 1)
	var association2 = col.Association[string, int]("bar", 2)
	var association3 = col.Association[string, int]("baz", 3)
	var associations = col.ListFromArray([]col.Binding[string, int]{association2, association3})
	var catalog = col.Catalog[string, int]()
	ass.True(t, catalog.IsEmpty())
	ass.Equal(t, 0, catalog.GetSize())
	ass.Equal(t, []string{}, catalog.GetKeys().AsArray())
	ass.Equal(t, []col.Binding[string, int]{}, catalog.AsArray())
	var iterator = col.Iterator[col.Binding[string, int]](catalog)
	ass.False(t, iterator.HasNext())
	ass.False(t, iterator.HasPrevious())
	iterator.ToStart()
	iterator.ToEnd()
	catalog.SortValues()
	catalog.ShuffleValues()
	catalog.RemoveAll()
	catalog.AddAssociation(association1)
	ass.False(t, catalog.IsEmpty())
	ass.Equal(t, 1, catalog.GetSize())
	catalog.AddAssociations(associations)
	ass.Equal(t, 3, catalog.GetSize())
	var catalog2 = col.CatalogFromSequence[string, int](catalog)
	ass.True(t, col.CompareValues(catalog, catalog2))
	var m = col.Map[string, int](map[string]int{
		"foo": 1,
		"bar": 2,
		"baz": 3,
	})
	var catalog3 = col.CatalogFromSequence[string, int](m)
	catalog2.SortValues()
	catalog3.SortValues()
	ass.True(t, col.CompareValues(catalog2, catalog3))
	iterator = col.Iterator[col.Binding[string, int]](catalog)
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
	var association1 = col.Association[string, int]("foo", 1)
	var association2 = col.Association[string, int]("bar", 2)
	var association3 = col.Association[string, int]("baz", 3)
	var catalog1 = col.Catalog[string, int]()
	catalog1.AddAssociation(association1)
	catalog1.AddAssociation(association2)
	var catalog2 = col.Catalog[string, int]()
	catalog2.AddAssociation(association2)
	catalog2.AddAssociation(association3)
	var catalog3 = col.Merge(catalog1, catalog2)
	var catalog4 = col.Catalog[string, int]()
	catalog4.AddAssociation(association1)
	catalog4.AddAssociation(association2)
	catalog4.AddAssociation(association3)
	ass.True(t, col.CompareValues(catalog3, catalog4))
}

func TestCatalogsWithExtract(t *tes.T) {
	var keys col.Sequential[string] = col.ListFromArray([]string{"foo", "baz"})
	var association1 = col.Association[string, int]("foo", 1)
	var association2 = col.Association[string, int]("bar", 2)
	var association3 = col.Association[string, int]("baz", 3)
	var catalog1 = col.Catalog[string, int]()
	catalog1.AddAssociation(association1)
	catalog1.AddAssociation(association2)
	catalog1.AddAssociation(association3)
	var catalog2 = col.Extract(catalog1, keys)
	var catalog3 = col.Catalog[string, int]()
	catalog3.AddAssociation(association1)
	catalog3.AddAssociation(association3)
	ass.True(t, col.CompareValues(catalog2, catalog3))
}

func TestCatalogsWithEmptyCatalogs(t *tes.T) {
	var keys col.Sequential[int] = col.List[int]()
	var catalog1 = col.Catalog[int, string]()
	var catalog2 = col.Catalog[int, string]()
	var catalog3 = col.Merge(catalog1, catalog2)
	var catalog4 = col.Extract(catalog1, keys)
	ass.True(t, col.CompareValues(catalog1, catalog2))
	ass.True(t, col.CompareValues(catalog2, catalog3))
	ass.True(t, col.CompareValues(catalog3, catalog4))
	ass.True(t, col.CompareValues(catalog4, catalog1))
}
