/*******************************************************************************
 *   Copyright (c) 2009-2023 Crater Dog Technologies™.  All Rights Reserved.   *
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

func TestStackWithSmallCapacity(t *tes.T) {
	var stack = col.StackWithCapacity[int](1)
	stack.AddValue(1)
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "Attempted to add a value onto a stack that has reached its capacity: 1\nvalue: 2\nstack: [\n\t1\n](stack)\n", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	stack.AddValue(2) // This should panic.
}

func TestEmptyStackRetrieval(t *tes.T) {
	var stack = col.Stack[int]()
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "Attempted to retrieve the top of an empty stack!", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	stack.GetTop() // This should panic.
}

func TestEmptyStackRemoval(t *tes.T) {
	var stack = col.Stack[int]()
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "Attempted to remove the top of an empty stack!", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	stack.RemoveTop() // This should panic.
}

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
