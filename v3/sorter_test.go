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

func TestSortingEmpty(t *tes.T) {
	var Sorter = col.Sorter[any]()
	var Collator = col.Collator()
	var empty = []any{}
	Sorter.SortValues(empty, Collator.RankValues)
}

func TestSortingIntegers(t *tes.T) {
	var Sorter = col.Sorter[int]()
	var Collator = col.Collator()
	var unsorted = []int{4, 3, 1, 5, 2}
	var sorted = []int{1, 2, 3, 4, 5}
	Sorter.SortValues(unsorted, Collator.RankValues)
	ass.Equal(t, sorted, unsorted)
}

func TestSortingStrings(t *tes.T) {
	var Sorter = col.Sorter[string]()
	var Collator = col.Collator()
	var unsorted = []string{"alpha", "beta", "gamma", "delta"}
	var sorted = []string{"alpha", "beta", "delta", "gamma"}
	Sorter.SortValues(unsorted, Collator.RankValues)
	ass.Equal(t, sorted, unsorted)
}
