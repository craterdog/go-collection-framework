/*******************************************************************************
 *   Copyright (c) 2009-2024 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package collections

import (
	reg "regexp"
	sts "strings"
	utf "unicode/utf8"
)

// CLASS NAMESPACE

// This private type defines the namespace structure associated with the
// constants, constructors and functions for the Scanner class namespace.
type scannerClass_ struct {
	booleanMatcher   *reg.Regexp
	complexMatcher   *reg.Regexp
	contextMatcher   *reg.Regexp
	delimiterMatcher *reg.Regexp
	floatMatcher     *reg.Regexp
	integerMatcher   *reg.Regexp
	nilMatcher       *reg.Regexp
	runeMatcher      *reg.Regexp
	stringMatcher    *reg.Regexp
	unsignedMatcher  *reg.Regexp
}

// These private constants define regular expression sub-patterns.
const (
	base10    = `[0-9]`
	base16    = `[0-9a-f]`
	boolean   = `false|true`
	complex_  = `\((` + float + `)` + sign + `(` + float + `)i\)`
	context   = `[Aa]rray|Catalog|List|[Mm]ap|Queue|Set|Stack`
	delimiter = `\[|\]|\(|\)|:|,`
	eol       = `\n`
	escape    = `\\(?:(?:` + unicode + `)|[abfnrtv'"\\])`
	exponent  = `[eE]` + sign + ordinal
	float     = sign + `?(?:` + scalar + `)(?:` + exponent + `)?`
	fraction  = `\.` + base10 + `+`
	integer   = zero + `|` + sign + `?` + ordinal
	nil_      = `nil`
	ordinal   = `[1-9][0-9]*`
	rune_     = `'(` + escape + `|[^'` + eol + `])'`
	scalar    = `(?:` + zero + `|` + ordinal + `)` + fraction
	sign      = `[+-]`
	string_   = `"(` + escape + `|[^"` + eol + `])*"`
	unicode   = `u` + base16 + `{4}|U` + base16 + `{8}`
	unsigned  = `0x` + base16 + `+`
	zero      = `0`
)

// This private constant defines the singleton reference to the Scanner
// class namespace.  It also initializes any class constants as needed.
var scannerClassSingleton = &scannerClass_{
	booleanMatcher:   reg.MustCompile(`^(?:` + boolean + `)`),
	complexMatcher:   reg.MustCompile(`^(?:` + complex_ + `)`),
	contextMatcher:   reg.MustCompile(`^(?:` + context + `)`),
	delimiterMatcher: reg.MustCompile(`^(?:` + delimiter + `)`),
	floatMatcher:     reg.MustCompile(`^(?:` + float + `)`),
	integerMatcher:   reg.MustCompile(`^(?:` + integer + `)`),
	nilMatcher:       reg.MustCompile(`^(?:` + nil_ + `)`),
	runeMatcher:      reg.MustCompile(`^(?:` + rune_ + `)`),
	stringMatcher:    reg.MustCompile(`^(?:` + string_ + `)`),
	unsignedMatcher:  reg.MustCompile(`^(?:` + unsigned + `)`),
}

// This public function returns the singleton reference to the Scanner
// class namespace.
func Scanner() *scannerClass_ {
	return scannerClassSingleton
}

// CLASS CONSTANTS

// This public class constant represents a regular expression that can be used
// to match strings containing boolean values.
func (c *scannerClass_) BooleanMatcher() *reg.Regexp {
	return c.booleanMatcher
}

// This public class constant represents a regular expression that can be used
// to match strings containing complex values.
func (c *scannerClass_) ComplexMatcher() *reg.Regexp {
	return c.complexMatcher
}

// This public class constant represents a regular expression that can be used
// to match strings containing context values.
func (c *scannerClass_) ContextMatcher() *reg.Regexp {
	return c.contextMatcher
}

// This public class constant represents a regular expression that can be used
// to match strings containing delimiter values.
func (c *scannerClass_) DelimiterMatcher() *reg.Regexp {
	return c.delimiterMatcher
}

// This public class constant represents a regular expression that can be used
// to match strings containing float values.
func (c *scannerClass_) FloatMatcher() *reg.Regexp {
	return c.floatMatcher
}

// This public class constant represents a regular expression that can be used
// to match strings containing integer values.
func (c *scannerClass_) IntegerMatcher() *reg.Regexp {
	return c.integerMatcher
}

// This public class constant represents a regular expression that can be used
// to match strings containing nil values.
func (c *scannerClass_) NilMatcher() *reg.Regexp {
	return c.nilMatcher
}

// This public class constant represents a regular expression that can be used
// to match strings containing rune values.
func (c *scannerClass_) RuneMatcher() *reg.Regexp {
	return c.runeMatcher
}

// This public class constant represents a regular expression that can be used
// to match strings containing string values.
func (c *scannerClass_) StringMatcher() *reg.Regexp {
	return c.stringMatcher
}

// This public class constant represents a regular expression that can be used
// to match strings containing unsigned values.
func (c *scannerClass_) UnsignedMatcher() *reg.Regexp {
	return c.unsignedMatcher
}

// CLASS CONSTRUCTORS

// This public class constructor creates a new Scanner from the specified
// source bytes.
func (c *scannerClass_) FromSource(
	source []byte,
	tokens chan TokenLike,
) ScannerLike {
	var scanner = &scanner_{
		line:     1,
		position: 1,
		source:   source,
		tokens:   tokens,
	}
	go scanner.scanTokens() // Start scanning tokens in the background.
	return scanner
}

// CLASS FUNCTIONS

// This public class function returns, for the specified string, an array of the
// matching subgroups for a Go boolean primitive. The first string in the array
// is the entire matched string.
func (c *scannerClass_) MatchBoolean(string_ string) []string {
	return c.booleanMatcher.FindStringSubmatch(string_)
}

// This public class function returns, for the specified string, an array of the
// matching subgroups for a Go complex primitive. The first string in the array
// is the entire matched string.
func (c *scannerClass_) MatchComplex(string_ string) []string {
	return c.complexMatcher.FindStringSubmatch(string_)
}

// This public class function returns, for the specified string, an array of the
// matching subgroups for a collection context. The first string in the array
// is the entire matched string.
func (c *scannerClass_) MatchContext(string_ string) []string {
	return c.contextMatcher.FindStringSubmatch(string_)
}

// This public class function returns, for the specified string, an array of the
// matching subgroups for a delimiter.
func (c *scannerClass_) MatchDelimiter(string_ string) []string {
	return c.delimiterMatcher.FindStringSubmatch(string_)
}

// This public class function returns, for the specified string, an array of the
// matching subgroups for a Go float primitive. The first string in the array
// is the entire matched string.
func (c *scannerClass_) MatchFloat(string_ string) []string {
	return c.floatMatcher.FindStringSubmatch(string_)
}

// This public class function returns, for the specified string, an array of the
// matching subgroups for a Go integer primitive. The first string in the array
// is the entire matched string.
func (c *scannerClass_) MatchInteger(string_ string) []string {
	return c.integerMatcher.FindStringSubmatch(string_)
}

// This public class function returns, for the specified string, an array of the
// matching subgroups for a Go nil primitive. The first string in the array
// is the entire matched string.
func (c *scannerClass_) MatchNil(string_ string) []string {
	return c.nilMatcher.FindStringSubmatch(string_)
}

// This public class function returns, for the specified string, an array of the
// matching subgroups for a Go rune primitive. The first string in the array
// is the entire matched string.
func (c *scannerClass_) MatchRune(string_ string) []string {
	return c.runeMatcher.FindStringSubmatch(string_)
}

// This public class function returns, for the specified string, an array of the
// matching subgroups for a Go string primitive. The first string in the array
// is the entire matched string.
func (c *scannerClass_) MatchString(string_ string) []string {
	return c.stringMatcher.FindStringSubmatch(string_)
}

// This public class function returns, for the specified string, an array of the
// matching subgroups for a Go unsigned primitive. The first string in the array
// is the entire matched string.
func (c *scannerClass_) MatchUnsigned(string_ string) []string {
	return c.unsignedMatcher.FindStringSubmatch(string_)
}

// CLASS TYPE

// Encapsulated Type

// This private class type encapsulates a Go structure containing private
// attributes that can only be accessed and manipulated using methods that
// implement the scanner-like abstract type.
//
// The source bytes can be viewed like this:
//
//	| byte 0 | byte 1 | byte 2 | byte 3 | byte 4 | byte 5 | ... | byte N-1 |
//	| rune 0 |      rune 1     |      rune 2     | rune 3 | ... | rune R-1 |
//
// Runes can be one to eight bytes long.
type scanner_ struct {
	firstByte int // A zero based index of the first possible byte in the next token.
	line      int // The line number in the source string of the next rune.
	nextByte  int // A zero based index of the next possible byte in the next token.
	position  int // The position in the current line of the next token.
	source    []byte
	tokens    chan TokenLike
}

// Private Interface

// This private class method adds a token of the specified type with the current
// Scanner information to the token channel. It then resets the first byte index
// to the next byte index position. It returns the token type of the type added
// to the channel.
func (v *scanner_) emitToken(tokenType string) string {
	var tokenValue = string(v.source[v.firstByte:v.nextByte])
	if tokenType == Token().TypeEOF() {
		tokenValue = "<EOF>"
	}
	if tokenType == Token().TypeError() {
		switch tokenValue {
		case "\a":
			tokenValue = "<BELL>"
		case "\b":
			tokenValue = "<BKSP>"
		case "\t":
			tokenValue = "<TAB>"
		case "\f":
			tokenValue = "<FF>"
		case "\r":
			tokenValue = "<CR>"
		case "\v":
			tokenValue = "<VTAB>"
		}
	}
	var token = Token().FromContext(v.line, v.position, tokenType, tokenValue)
	//fmt.Println(token)
	v.tokens <- token
	v.firstByte = v.nextByte
	v.position += sts.Count(tokenValue, "") - 1 // Add the number of runes in the token.
	return tokenType
}

// This private class method adds a boolean token with the current Scanner
// information to the token channel. It returns true if a boolean token was found.
func (v *scanner_) foundBoolean() bool {
	var string_ = string(v.source[v.nextByte:])
	var matches = Scanner().MatchBoolean(string_)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(Token().TypeBoolean())
		return true
	}
	return false
}

// This private class method adds a complex number token with the current
// Scanner information to the token channel. It returns true if a complex number
// token was found.
func (v *scanner_) foundComplex() bool {
	var string_ = string(v.source[v.nextByte:])
	var matches = Scanner().MatchComplex(string_)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(Token().TypeComplex())
		return true
	}
	return false
}

// This private class method adds a context token with the current Scanner
// information to the token channel. It returns true if a type token was found.
func (v *scanner_) foundContext() bool {
	var string_ = string(v.source[v.nextByte:])
	var matches = Scanner().MatchContext(string_)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(Token().TypeContext())
		return true
	}
	return false
}

// This private class method adds a delimiter token with the current Scanner
// information to the token channel. It returns true if a delimiter token was
// found.
func (v *scanner_) foundDelimiter() bool {
	var string_ = string(v.source[v.nextByte:])
	var matches = Scanner().MatchDelimiter(string_)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(Token().TypeDelimiter())
		return true
	}
	return false
}

// This private class method adds an EOF token with the current Scanner
// information to the token channel. It returns true if an EOF marker was found.
func (v *scanner_) foundEOF() bool {
	// The last byte in a POSIX standard file must be an EOL character.
	var string_ = string(v.source[v.nextByte:])
	if !sts.HasPrefix(string_, EOL) {
		return false
	}
	v.nextByte++
	v.line++
	// Now make sure there are no more bytes.
	if v.nextByte != len(v.source) {
		v.nextByte--
		v.line--
		return false
	}
	v.emitToken(Token().TypeEOF())
	return true
}

// This private class method adds an EOL token with the current Scanner
// information to the token channel. It returns true if an EOL token was found.
func (v *scanner_) foundEOL() bool {
	var string_ = string(v.source[v.nextByte:])
	if !sts.HasPrefix(string_, EOL) {
		return false
	}
	v.nextByte++
	v.emitToken(Token().TypeEOL())
	v.line++
	v.position = 1
	return true
}

// This private class method adds an error token with the current Scanner
// information to the token channel. It always returns true.
func (v *scanner_) foundError() bool {
	var bytes = v.source[v.nextByte:]
	var _, width = utf.DecodeRune(bytes)
	v.nextByte += width
	v.emitToken(Token().TypeError())
	return true
}

// This private class method adds a floating point token with the current
// Scanner information to the token channel. It returns true if a floating point
// token was found.
func (v *scanner_) foundFloat() bool {
	var string_ = string(v.source[v.nextByte:])
	var matches = Scanner().MatchFloat(string_)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(Token().TypeFloat())
		return true
	}
	return false
}

// This private class method adds a integer token with the current Scanner
// information to the token channel. It returns true if a integer token was
// found.
func (v *scanner_) foundInteger() bool {
	var string_ = string(v.source[v.nextByte:])
	var matches = Scanner().MatchInteger(string_)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(Token().TypeInteger())
		return true
	}
	return false
}

// This private class method adds a nil token with the current Scanner
// information to the token channel. It returns true if a nil token was found.
func (v *scanner_) foundNil() bool {
	var string_ = string(v.source[v.nextByte:])
	var matches = Scanner().MatchNil(string_)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(Token().TypeNil())
		return true
	}
	return false
}

// This private class method adds a rune token with the current Scanner
// information to the token channel. It returns true if a rune token was found.
func (v *scanner_) foundRune() bool {
	var string_ = string(v.source[v.nextByte:])
	var matches = Scanner().MatchRune(string_)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(Token().TypeRune())
		return true
	}
	return false
}

// This private class method adds a string token with the current Scanner
// information to the token channel. It returns true if a string token was
// found.
func (v *scanner_) foundString() bool {
	var string_ = string(v.source[v.nextByte:])
	var matches = Scanner().MatchString(string_)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(Token().TypeString())
		return true
	}
	return false
}

// This private class method adds an unsigned integer token with the current
// Scanner information to the token channel. It returns true if an unsigned
// integer token was found.
func (v *scanner_) foundUnsigned() bool {
	var string_ = string(v.source[v.nextByte:])
	var matches = Scanner().MatchUnsigned(string_)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(Token().TypeUnsigned())
		return true
	}
	return false
}

// This private class method attempts to scan any token starting with the next
// rune in the source array. It checks for each type of token as the cases for
// the switch statement. If that token type is found, this method returns true
// and skips the rest of the cases.  If no valid token is found, or a TypeEOF is
// found this method returns false.
func (v *scanner_) scanToken() bool {
	v.skipSpaces()
	switch {
	case v.foundBoolean():
	case v.foundComplex():
	case v.foundContext():
	case v.foundDelimiter():
	case v.foundEOF():
		// We are at the end of the source array.
		return false
	case v.foundEOL():
	case v.foundFloat():
	case v.foundNil():
	case v.foundRune():
	case v.foundString():
	case v.foundUnsigned():
	case v.foundInteger(): // Must be after all other numeric types.
	case v.foundError(): // Must be last.
		// No valid token was found.
		return false
	}
	return true
}

// This private class method continues scanning tokens from the source array
// until an error occurs or the end of file is reached. It then closes the token
// channel.
func (v *scanner_) scanTokens() {
	for v.scanToken() {
	}
	close(v.tokens)
}

// This private class method scans through any spaces in the source array and
// sets the next byte index to the next non-space rune.
func (v *scanner_) skipSpaces() {
	if v.nextByte < len(v.source) {
		for {
			if v.source[v.nextByte] != ' ' {
				break
			}
			v.nextByte++
			v.position++
		}
		v.firstByte = v.nextByte
	}
}
