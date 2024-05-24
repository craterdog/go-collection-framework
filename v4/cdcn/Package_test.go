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
	not "github.com/craterdog/go-collection-framework/v4/cdcn"
	col "github.com/craterdog/go-collection-framework/v4/collection"
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

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
	parser.ParseSource(source)
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
	parser.ParseSource(source)
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
	parser.ParseSource(source)
}

func TestCollectionRoundtrip(t *tes.T) {
	var notation = not.Notation().Make()
	var collection = notation.ParseSource(data)
	var actual = notation.FormatCollection(collection)
	ass.Equal(t, data, actual)
}

const data = `[
    [
        true
        0xa
        42
        0.125
        (3.0-4.0i)
        '☺'
        "Hello World!"
        [
            1
            2
            3
        ](array)
        [
            [ ](array)
            [:](map)
        ](array)
        ["alpha": 1](map)
    ](Array)
    [ ](array)
    [
        1
        2
        3
        4
    ](array)
    [
        false: true
        0x0: 0xa
        0: 42
        0.0: 0.125
        (0.0+0.0i): (1.0+1.0i)
        '\x00': '☺'
        "": "Hello World!"
    ](Catalog)
    [:](Catalog)
    [
        "boolean": true
        "unsigned": 0xa
        "integer": 42
        "float": 0.125
        "complex": (3.0-4.0i)
        "rune": '☺'
        "string": "Hello World!"
        "empty": [
            [ ](array)
            [:](map)
        ](array)
        "array": [
            1
            2
            3
        ](array)
        "map": ["alpha": 1](map)
    ](Catalog)
    [
        "key1": 1
        "key2": 2
        "key3": 3
        "key4": 4
        "key5": 5
        "key6": 6
        "key7": 7
    ](Catalog)
    [
        true
        0xa
        42
        0.125
        (1.0+1.0i)
        '☺'
        "Hello World!"
        [
            1
            2
            3
        ](array)
    ](List)
    [
        false
        true
    ](List)
    [ ](List)
    [
        true
        0xa
        42
        0.125
        1.1E-100
        2.2E+200
        (3.0-4.0i)
        '☺'
        "Hello World!"
        [
            1
            2
            3
        ](array)
        [
            [ ](array)
            [:](map)
        ](array)
        ["alpha": 1](map)
    ](List)
    [
        0x0: 0xa
        '\x00': '☺'
        "array": [
            1
            2
            3
        ](array)
        0.0: 0.125
        false: true
        "": "Hello World!"
        0: 42
        (0.0+0.0i): (1.0+1.0i)
    ](Catalog)
    [:](map)
    [
        "second": 2
        "third": 3
        "first": 1
    ](Catalog)
    [
        true
        0xa
        42
        0.125
        (1.0+1.0i)
        '☺'
        "Hello World!"
    ](Queue)
    [
        1
        2
        3
        4
    ](Queue)
    [
        [
            1
            2
            3
        ](array)
        true
        (1.0+1.0i)
        0.125
        42
        '☺'
        "Hello World!"
        0xa
    ](Set)
    [
        [
            [
                1
                2
                3
            ](array)
            true
            (1.0+1.0i)
            0.125
            42
            '☺'
            "Hello World!"
            0xa
        ](Set)
        [
            false
            (0.0+0.0i)
            0.0
            0
            '\x00'
            ""
            "array"
            0x0
        ](Set)
    ](Set)
    [
        "blue"
        "green"
        "indigo"
        "orange"
        "red"
        "violet"
        "yellow"
    ](Set)
    [
        "Hello World!"
        '☺'
        (1.0+1.0i)
        0.125
        42
        0xa
        true
    ](Stack)
    [
        4
        3
        2
        1
    ](Stack)
](List)
`
