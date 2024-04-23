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
	fmt "fmt"
	not "github.com/craterdog/go-collection-framework/v3/cdcn"
	col "github.com/craterdog/go-collection-framework/v3/collection"
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

func TestStackConstructor(t *tes.T) {
	var notation = not.Notation().Make()
	var Stack = col.Stack[int64](notation)
	var empty = Stack.MakeFromSource("[ ](Stack)")
	ass.Equal(t, "[ ](Stack)\n", fmt.Sprintf("%v", empty))
	var full = Stack.MakeFromSource("[1, 2, 3](Stack)")
	ass.Equal(t, `[
    1
    2
    3
](Stack)
`, fmt.Sprintf("%v", full))
}

func TestStackConstructors(t *tes.T) {
	var notation = not.Notation().Make()
	var Stack = col.Stack[int64](notation)
	var stack1 = Stack.MakeFromArray([]int64{1, 2, 3})
	var stack2 = Stack.MakeFromSequence(stack1)
	ass.Equal(t, stack1.AsArray(), stack2.AsArray())
}

func TestStackWithSmallCapacity(t *tes.T) {
	var notation = not.Notation().Make()
	var stack = col.Stack[int](notation).MakeWithCapacity(1)
	stack.AddValue(1)
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "Attempted to add a value onto a stack that has reached its capacity: 1\nvalue: 2\nstack: [1](Stack)\n", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	stack.AddValue(2) // This should panic.
}

func TestEmptyStackRemoval(t *tes.T) {
	var notation = not.Notation().Make()
	var stack = col.Stack[int](notation).Make()
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
	var notation = not.Notation().Make()
	var stack = col.Stack[string](notation).Make()
	ass.True(t, stack.IsEmpty())
	ass.Equal(t, 0, stack.GetSize())
	stack.RemoveAll()
	stack.AddValue("foo")
	stack.AddValue("bar")
	stack.AddValue("baz")
	ass.Equal(t, 3, stack.GetSize())
	ass.Equal(t, "baz", string(stack.RemoveTop()))
	ass.Equal(t, 2, stack.GetSize())
	ass.Equal(t, "bar", string(stack.RemoveTop()))
	ass.Equal(t, 1, stack.GetSize())
	stack.RemoveAll()
}
