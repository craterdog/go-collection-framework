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
	col "github.com/craterdog/go-collection-framework/collections"
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

func TestStacksWithStrings(t *tes.T) {
	var stack = col.Stack[string]()
	ass.True(t, stack.IsEmpty())
	ass.Equal(t, 0, stack.GetSize())
	stack.RemoveAll()
	stack.AddValue("foo")
	stack.AddValue("bar")
	stack.AddValue("baz")
	ass.Equal(t, 3, stack.GetSize())
	ass.Equal(t, "baz", string(stack.GetTop()))
	ass.Equal(t, "baz", string(stack.RemoveTop()))
	ass.Equal(t, 2, stack.GetSize())
	ass.Equal(t, "bar", string(stack.GetTop()))
	ass.Equal(t, "bar", string(stack.RemoveTop()))
	ass.Equal(t, 1, stack.GetSize())
	ass.Equal(t, "foo", string(stack.GetTop()))
	stack.RemoveAll()
}
