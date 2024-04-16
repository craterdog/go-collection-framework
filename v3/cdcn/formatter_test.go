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

package cdcn_test

import (
	cdc "github.com/craterdog/go-collection-framework/v3/cdcn"
	col "github.com/craterdog/go-collection-framework/v3/collection"
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

func TestFormatMaximum(t *tes.T) {
	var Formatter = cdc.Formatter()
	var formatter = Formatter.MakeWithMaximum(0)
	var array = col.Array[any]().MakeFromArray([]any{1, []any{1, 2, []any{1, 2, 3}}})
	var s = formatter.FormatCollection(array)
	ass.Equal(t, "[...](Array)\n", s)
	formatter = Formatter.MakeWithMaximum(1)
	s = formatter.FormatCollection(array)
	ass.Equal(t, "[\n    1\n    [...](array)\n](Array)\n", s)
	formatter = Formatter.MakeWithMaximum(2)
	s = formatter.FormatCollection(array)
	ass.Equal(t, "[\n    1\n    [\n        1\n        2\n        [...](array)\n    ](array)\n](Array)\n", s)
}

func TestFormatInvalidType(t *tes.T) {
	var formatter = cdc.Formatter().MakeWithMaximum(8)
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
