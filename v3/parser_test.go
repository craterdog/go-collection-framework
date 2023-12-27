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

func TestParseBadFirst(t *tes.T) {
	var parser = col.ParserClass().CDCN()
	var source = `bad[ ](array)
`
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(
				t,
				"An unexpected token was received by the parser: Token [type: Error, line: 1, position: 1]: \"b\"\n\x1b[36m0001: bad[ ](array)\n \x1b[32m>>>──⌃\x1b[36m\n0002: \n\x1b[0m\n",
				e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	var _ = parser.ParseCollection(source)
}

func TestParseBadMiddle(t *tes.T) {
	var parser = col.ParserClass().CDCN()
	var source = `[bad](array)
`
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(
				t,
				"An unexpected token was received by the parser: Token [type: Error, line: 1, position: 2]: \"b\"\n\x1b[36m0001: [bad](array)\n \x1b[32m>>>───⌃\x1b[36m\n0002: \n\x1b[0m\n",
				e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	var _ = parser.ParseCollection(source)
}

func TestParseBadEnd(t *tes.T) {
	var parser = col.ParserClass().CDCN()
	var source = `[ ](array)bad
`
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(
				t,
				"An unexpected token was received by the parser: Token [type: Error, line: 1, position: 11]: \"b\"\n\x1b[36m0001: [ ](array)bad\n \x1b[32m>>>────────────⌃\x1b[36m\n0002: \n\x1b[0m\n",
				e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	var _ = parser.ParseCollection(source)
}

func TestParseExtraEOL(t *tes.T) {
	var parser = col.ParserClass().CDCN()
	var source = `[ ](array)

`
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(
				t,
				"An unexpected token was received by the parser: Token [type: EOL, line: 1, position: 11]: <EOL>\n\x1b[36m0001: [ ](array)\n \x1b[32m>>>────────────⌃\x1b[36m\n0002: \n\x1b[0m\nWas expecting 'EOF' from:\n  \x1b[32m$source: \x1b[33mcollection EOF  ! EOF is the end-of-file marker.\x1b[0m\n\n  \x1b[32m$collection: \x1b[33m\"[\" (associations | values) \"]\" context\x1b[0m\n\n",
				e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	var _ = parser.ParseCollection(source)
}
