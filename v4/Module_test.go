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

package module_test

import (
	col "github.com/craterdog/go-collection-framework/v4"
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

func TestCDCNConstructor(t *tes.T) {
	var _ = col.CDCN()
}

func TestXMLConstructor(t *tes.T) {
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The XML notation is not yet supported.", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	var _ = col.XML() // This should panic.
}

func TestJSONConstructor(t *tes.T) {
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The JSON notation is not yet supported.", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	var _ = col.JSON() // This should panic.
}

func TestCollectionConstructors(t *tes.T) {
	var notation = col.CDCN()
	var array = col.Array[string](3, []string{"foo", "bar"})
	var map_ = col.Map[string, int64](`["alpha": 1, "beta": 2, "gamma": 3](Map)`)
	var list = col.List[int64](notation, "[1, 2, 3](List)")
	var set = col.Set[int64](array)
	var catalog = col.Catalog[string, int64](map_)
	var associations = col.List[col.AssociationLike[string, int64]](catalog)
	var association = col.Association[string, int64]("delta", 4)
	associations.AppendValue(association)
	var _ = col.Stack[int64](4, list)
	var _ = col.Queue[int64](notation, 5, set)
}
