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

package agent_test

import (
	age "github.com/craterdog/go-collection-framework/v4/agent"
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

func TestSortingEmpty(t *tes.T) {
	var collator = age.Collator[any]().Make()
	var ranker = collator.RankValues
	var sorter = age.Sorter[any]().MakeWithRanker(ranker)
	var empty = []any{}
	sorter.SortValues(empty)
}

func TestSortingIntegers(t *tes.T) {
	var collator = age.Collator[int]().Make()
	var ranker = collator.RankValues
	var sorter = age.Sorter[int]().MakeWithRanker(ranker)
	var unsorted = []int{4, 3, 1, 5, 2}
	var sorted = []int{1, 2, 3, 4, 5}
	sorter.SortValues(unsorted)
	ass.Equal(t, sorted, unsorted)
}

func TestSortingStrings(t *tes.T) {
	var collator = age.Collator[string]().Make()
	var ranker = collator.RankValues
	var sorter = age.Sorter[string]().MakeWithRanker(ranker)
	var unsorted = []string{"alpha", "beta", "gamma", "delta"}
	var sorted = []string{"alpha", "beta", "delta", "gamma"}
	sorter.SortValues(unsorted)
	ass.Equal(t, sorted, unsorted)
}
