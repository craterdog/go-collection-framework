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
	fmt "fmt"
	not "github.com/craterdog/go-collection-framework/v4/cdcn"
	col "github.com/craterdog/go-collection-framework/v4/collection"
	ass "github.com/stretchr/testify/assert"
	osx "os"
	sts "strings"
	tes "testing"
)

const collectionTests = "../test/input/"

func TestFormatMaximum(t *tes.T) {
	var notation = not.Notation().Make()
	var Formatter = not.Formatter()
	var formatter = Formatter.MakeWithMaximum(0)
	var array = col.Array[any](notation).MakeFromArray([]any{1, []any{1, 2, []any{1, 2, 3}}})
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
	var formatter = not.Formatter().MakeWithMaximum(8)
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

func TestCollectionRoundtrips(t *tes.T) {
	var notation = not.Notation().Make()
	var files, err = osx.ReadDir(collectionTests)
	if err != nil {
		var message = fmt.Sprintf("Could not find the %s directory.", collectionTests)
		panic(message)
	}
	for _, file := range files {
		var filename = collectionTests + file.Name()
		if sts.HasSuffix(filename, ".cdcn") {
			fmt.Println(filename)
			var bytes, err = osx.ReadFile(filename)
			if err != nil {
				panic(err)
			}
			var expected = string(bytes)
			var collection = notation.ParseSource(expected)
			var actual = notation.FormatCollection(collection)
			if !sts.HasPrefix(file.Name(), "map") {
				// Skip maps since they are non-deterministic.
				ass.Equal(t, expected, actual)
				bytes = []byte(actual)
				err = osx.WriteFile(filename, bytes, 0644)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

func TestParseBadFirst(t *tes.T) {
	var parser = not.Parser().Make()
	var source = `bad[ ](array)
`
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(
				t,
				"An unexpected token was received by the parser: Token [type: error, line: 1, position: 1]: \"b\"\n\x1b[36m0001: bad[ ](array)\n \x1b[32m>>>──⌃\x1b[36m\n0002: \n\x1b[0m\n",
				e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	var _ = parser.ParseSource(source)
}

func TestParseBadMiddle(t *tes.T) {
	var parser = not.Parser().Make()
	var source = `[bad](array)
`
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(
				t,
				"An unexpected token was received by the parser: Token [type: error, line: 1, position: 2]: \"b\"\n\x1b[36m0001: [bad](array)\n \x1b[32m>>>───⌃\x1b[36m\n0002: \n\x1b[0m\n",
				e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	var _ = parser.ParseSource(source)
}

func TestParseBadEnd(t *tes.T) {
	var parser = not.Parser().Make()
	var source = `[ ](array)bad
`
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(
				t,
				"An unexpected token was received by the parser: Token [type: error, line: 1, position: 11]: \"b\"\n\x1b[36m0001: [ ](array)bad\n \x1b[32m>>>────────────⌃\x1b[36m\n0002: \n\x1b[0m\n",
				e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	var _ = parser.ParseSource(source)
}
