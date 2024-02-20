/*******************************************************************************
 *   Copyright (c) 2009-2024 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package cdcn

import (
	//fmt "fmt"
	col "github.com/craterdog/go-collection-framework/v3"
	reg "regexp"
	sts "strings"
)

// CLASS ACCESS

// Reference

var scannerClass = &scannerClass_{
	matchers: map[TokenType]*reg.Regexp{
		BooleanToken:     reg.MustCompile(`^(?:` + boolean_ + `)`),
		ComplexToken:     reg.MustCompile(`^(?:` + complex_ + `)`),
		ContextToken:     reg.MustCompile(`^(?:` + context_ + `)`),
		DelimiterToken:   reg.MustCompile(`^(?:` + delimiter_ + `)`),
		EOLToken:         reg.MustCompile(`^(?:` + eol_ + `)`),
		FloatToken:       reg.MustCompile(`^(?:` + float_ + `)`),
		HexadecimalToken: reg.MustCompile(`^(?:` + hexadecimal_ + `)`),
		IntegerToken:     reg.MustCompile(`^(?:` + integer_ + `)`),
		NilToken:         reg.MustCompile(`^(?:` + nil_ + `)`),
		RuneToken:        reg.MustCompile(`^(?:` + rune_ + `)`),
		SpaceToken:       reg.MustCompile(`^(?:` + space_ + `)`),
		StringToken:      reg.MustCompile(`^(?:` + string_ + `)`),
	},
}

// Function

func Scanner() ScannerClassLike {
	return scannerClass
}

// CLASS METHODS

// Target

type scannerClass_ struct {
	matchers map[TokenType]*reg.Regexp
}

// Constructors

func (c *scannerClass_) Make(
	source string,
	tokens col.QueueLike[TokenLike],
) ScannerLike {
	var scanner = &scanner_{
		line:     1,
		position: 1,
		runes:    []rune(source),
		tokens:   tokens,
	}
	go scanner.scanTokens() // Start scanning tokens in the background.
	return scanner
}

// Functions

func (c *scannerClass_) MatchToken(
	tokenType TokenType,
	text string,
) col.ListLike[string] {
	var matcher = c.matchers[tokenType]
	var matches = matcher.FindStringSubmatch(text)
	return col.List[string]().MakeFromArray(matches)
}

// INSTANCE METHODS

// Target

type scanner_ struct {
	first    int // A zero based index of the first possible rune in the next token.
	line     int // The line number in the source string of the next rune.
	next     int // A zero based index of the next possible rune in the next token.
	position int // The position in the current line of the next rune.
	runes    []rune
	tokens   col.QueueLike[TokenLike]
}

// Private

func (v *scanner_) emitToken(tokenType TokenType) {
	var tokenValue = string(v.runes[v.first:v.next])
	switch tokenValue {
	case "\x00":
		tokenValue = "<NULL>"
	case "\a":
		tokenValue = "<BELL>"
	case "\b":
		tokenValue = "<BKSP>"
	case "\t":
		tokenValue = "<HTAB>"
	case "\f":
		tokenValue = "<FMFD>"
	case "\n":
		tokenValue = "<EOLN>"
	case "\r":
		tokenValue = "<CRTN>"
	case "\v":
		tokenValue = "<VTAB>"
	}
	var token = Token().Make(v.line, v.position, tokenType, tokenValue)
	//fmt.Println(token) // Uncomment when debugging.
	v.tokens.AddValue(token) // Will block if queue is full.
}

func (v *scanner_) foundEOF() {
	v.emitToken(EOFToken)
}

func (v *scanner_) foundError() {
	v.next++
	v.emitToken(ErrorToken)
}

func (v *scanner_) foundToken(tokenType TokenType) bool {
	var text = string(v.runes[v.next:])
	var matches = Scanner().MatchToken(tokenType, text)
	if !matches.IsEmpty() {
		var match = matches.GetValue(1)
		var token = []rune(match)
		v.next += len(token)
		if tokenType != SpaceToken {
			v.emitToken(tokenType)
		}
		var count = sts.Count(match, "\n")
		if count > 0 {
			v.line += count
			v.position = v.indexOfLastEOL(token)
		} else {
			v.position += v.next - v.first
		}
		v.first = v.next
		return true
	}
	return false
}

func (v *scanner_) indexOfLastEOL(runes []rune) int {
	var length = len(runes)
	for index := length; index > 0; index-- {
		if runes[index-1] == '\n' {
			return length - index + 1
		}
	}
	return 0
}

func (v *scanner_) scanTokens() {
loop:
	for v.next < len(v.runes) {
		switch {
		case v.foundToken(BooleanToken):
		case v.foundToken(ComplexToken):
		case v.foundToken(ContextToken):
		case v.foundToken(DelimiterToken):
		case v.foundToken(EOLToken):
		case v.foundToken(FloatToken):
		case v.foundToken(HexadecimalToken):
		case v.foundToken(IntegerToken):
		case v.foundToken(NilToken):
		case v.foundToken(RuneToken):
		case v.foundToken(SpaceToken):
		case v.foundToken(StringToken):
		default:
			v.foundError()
			break loop
		}
	}
	v.foundEOF()
	v.tokens.CloseQueue()
}

/*
These private constants define the regular expression sub-patterns that make up
all token types.  Unfortunately there is no way to make them private to the
scanner class namespace since they must be TRUE Go constants to be initialized
in this way.  We append an underscore to each name to lessen the chance of a
name collision with other private Go class constants in this package.
*/
const (
	base10_      = `[0-9]`
	base16_      = `[0-9a-f]`
	boolean_     = `false|true`
	complex_     = `\((` + float_ + `)` + sign_ + `(` + float_ + `)i\)`
	context_     = `[Aa]rray|Catalog|List|[Mm]ap|Queue|Set|Stack`
	delimiter_   = `\[|\]|\(|\)|:|,`
	eol_         = `\n`
	escape_      = `\\(?:(?:` + unicode_ + `)|[abfnrtv'"\\])`
	exponent_    = `[eE]` + sign_ + ordinal_
	float_       = sign_ + `?(?:` + scalar_ + `)(?:` + exponent_ + `)?`
	fraction_    = `\.` + base10_ + `+`
	hexadecimal_ = `0x` + base16_ + `+`
	integer_     = zero_ + `|` + sign_ + `?` + ordinal_
	nil_         = `nil`
	ordinal_     = `[1-9][0-9]*`
	rune_        = `'(` + escape_ + `|[^'` + eol_ + `])'`
	scalar_      = `(?:` + zero_ + `|` + ordinal_ + `)` + fraction_
	sign_        = `[+-]`
	space_       = `[ ]+`
	string_      = `"(` + escape_ + `|[^"` + eol_ + `])*"`
	unicode_     = `x` + base16_ + `{2}|u` + base16_ + `{4}|U` + base16_ + `{8}`
	zero_        = `0`
)
