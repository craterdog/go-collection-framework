/*******************************************************************************
 *   Copyright (c) 2009-2024 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package cdcn_test

import (
	col "github.com/craterdog/go-collection-framework/v3"
	not "github.com/craterdog/go-collection-framework/v3/cdcn"
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

func TestFormatDepths(t *tes.T) {
	var Formatter = not.FormatterClass()
	var formatter = Formatter.MakeWithDepth(0)
	var array = col.ArrayClass[any]().MakeFromArray([]any{1, []any{1, 2, []any{1, 2, 3}}})
	var s = formatter.FormatCollection(array)
	ass.Equal(t, "[...](Array)", s)
	formatter = Formatter.MakeWithDepth(1)
	s = formatter.FormatCollection(array)
	ass.Equal(t, "[\n    1\n    [...](array)\n](Array)", s)
	formatter = Formatter.MakeWithDepth(2)
	s = formatter.FormatCollection(array)
	ass.Equal(t, "[\n    1\n    [\n        1\n        2\n        [...](array)\n    ](array)\n](Array)", s)
}

func TestFormatInvalidType(t *tes.T) {
	var formatter = not.FormatterClass().MakeWithDepth(8)
	var s struct{}
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "Attempted to format:\n    value: {}\n    type: struct {}\n    kind: struct\n", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	formatter.FormatCollection(s) // This should panic.
}
