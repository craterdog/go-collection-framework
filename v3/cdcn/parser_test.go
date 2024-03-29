/*
................................................................................
.    Copyright (c) 2009-2024 Crater Dog Technologies.  All Rights Reserved.    .
................................................................................
.  DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               .
.                                                                              .
.  This code is free software; you can redistribute it and/or modify it under  .
.  the terms of The MIT License (MIT), as published by the Open Source         .
.  Initiative. (See http://opensource.org/licenses/MIT)                        .
................................................................................
*/

package cdcn_test

import (
	not "github.com/craterdog/go-collection-framework/v3/cdcn"
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

func TestParseBadFirst(t *tes.T) {
	var parser = not.Parser().Make()
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
				"An unexpected token was received by the parser: Token [type: Error, line: 1, position: 2]: \"b\"\n\x1b[36m0001: [bad](array)\n \x1b[32m>>>───⌃\x1b[36m\n0002: \n\x1b[0m\n",
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
				"An unexpected token was received by the parser: Token [type: Error, line: 1, position: 11]: \"b\"\n\x1b[36m0001: [ ](array)bad\n \x1b[32m>>>────────────⌃\x1b[36m\n0002: \n\x1b[0m\n",
				e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	var _ = parser.ParseSource(source)
}

func TestParseExtraEOL(t *tes.T) {
	var parser = not.Parser().Make()
	var source = `[ ](array)

`
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(
				t,
				"An unexpected token was received by the parser: Token [type: EOL, line: 1, position: 11]: \"<EOLN>\"\n\x1b[36m0001: [ ](array)\n \x1b[32m>>>────────────⌃\x1b[36m\n0002: \n\x1b[0m\nWas expecting 'EOF' from:\n  \x1b[32msource: \x1b[33mcollection EOF  ! Terminated with an end-of-file marker.\x1b[0m\n\n  \x1b[32mcollection: \x1b[33m\"[\" (associations | values) \"]\" \"(\" Context \")\"\x1b[0m\n\n",
				e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	var _ = parser.ParseSource(source)
}
