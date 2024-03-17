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
	col "github.com/craterdog/go-collection-framework/v2"
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

func TestSortingEmpty(t *tes.T) {
	var empty = []any{}
	col.SortArray(empty, col.RankValues)
}

func TestSortingIntegers(t *tes.T) {
	var unsorted = []int{4, 3, 1, 5, 2}
	var sorted = []int{1, 2, 3, 4, 5}
	col.SortArray(unsorted, col.RankValues)
	ass.Equal(t, sorted, unsorted)
}

func TestSortingStrings(t *tes.T) {
	var unsorted = []string{"alpha", "beta", "gamma", "delta"}
	var sorted = []string{"alpha", "beta", "delta", "gamma"}
	col.SortArray(unsorted, col.RankValues)
	ass.Equal(t, sorted, unsorted)
}
