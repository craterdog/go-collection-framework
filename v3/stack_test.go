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

func TestStackConstructors(t *tes.T) {
	var Stack = col.Stack[int]()
	var stack1 = Stack.FromArray([]int{1, 2, 3})
	var stack2 = Stack.FromSequence(stack1)
	ass.Equal(t, stack1.AsArray(), stack2.AsArray())
}

func TestStackWithSmallCapacity(t *tes.T) {
	var Stack = col.Stack[int]()
	var stack = Stack.WithCapacity(1)
	stack.AddValue(1)
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "Attempted to add a value onto a Stack that has reached its capacity: 1\nvalue: 2\nstack: [1](Stack)\n\n", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	stack.AddValue(2) // This should panic.
}

func TestEmptyStackRetrieval(t *tes.T) {
	var Stack = col.Stack[int]()
	var stack = Stack.Empty()
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "Attempted to retrieve the top of an empty Stack!", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	stack.GetTop() // This should panic.
}

func TestEmptyStackRemoval(t *tes.T) {
	var Stack = col.Stack[int]()
	var stack = Stack.Empty()
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "Attempted to remove the top of an empty Stack!", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	stack.RemoveTop() // This should panic.
}

func TestStacksWithStrings(t *tes.T) {
	var Stack = col.Stack[string]()
	var stack = Stack.Empty()
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
