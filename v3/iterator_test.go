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
	col "github.com/craterdog/go-collection-framework/v3"
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

func TestIteratorsWithLists(t *tes.T) {
	var array = col.Array[int]().MakeFromArray([]int{1, 2, 3, 4, 5})
	var list = col.List[int]().MakeFromSequence(array)
	var iterator = list.GetIterator()
	ass.False(t, iterator.HasPrevious())
	ass.True(t, iterator.HasNext())
	ass.Equal(t, 1, iterator.GetNext())
	ass.True(t, iterator.HasPrevious())
	ass.True(t, iterator.HasNext())
	ass.Equal(t, 1, iterator.GetPrevious())
	iterator.ToSlot(2)
	ass.True(t, iterator.HasPrevious())
	ass.True(t, iterator.HasNext())
	ass.Equal(t, 3, iterator.GetNext())
	iterator.ToEnd()
	ass.True(t, iterator.HasPrevious())
	ass.False(t, iterator.HasNext())
	ass.Equal(t, 5, iterator.GetPrevious())
	iterator.ToStart()
	ass.False(t, iterator.HasPrevious())
	ass.True(t, iterator.HasNext())
	ass.Equal(t, 1, iterator.GetNext())
}
