/*******************************************************************************
 *   Copyright (c) 2009-2023 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package collections

import (
	byt "bytes"
	fmt "fmt"
	reg "regexp"
	sts "strings"
	utf "unicode/utf8"
)

// TOKENS

// This integer type is used as a type identifier for each token.
type TokenType int

// This enumeration defines all possible token types including the error token.
const (
	TokenError TokenType = iota
	TokenBoolean
	TokenComplex
	TokenContext
	TokenDelimiter
	TokenEOF
	TokenEOL
	TokenFloat
	TokenInteger
	TokenNil
	TokenRune
	TokenString
	TokenUnsigned
)

// This method returns the string representation for each token type.
func (v TokenType) String() string {
	return [...]string{
		"Error",
		"Boolean",
		"Complex",
		"Context",
		"Delimiter",
		"EOF",
		"EOL",
		"Float",
		"Integer",
		"Nil",
		"Rune",
		"String",
		"Unsigned",
	}[v]
}

// This type defines the structure and methods for each token returned by the
// scanner.
type Token struct {
	Type     TokenType
	Value    string
	Line     int // The line number of the token in the input string.
	Position int // The position in the line of the first rune of the token.
}

// This method returns the canonical string version of this token.
func (v Token) String() string {
	var s string
	switch {
	case v.Type == TokenEOF:
		s = "<EOF>"
	case v.Type == TokenEOL:
		s = "<EOL>"
	case len(v.Value) > 60:
		s = fmt.Sprintf("%.60q...", v.Value)
	default:
		s = fmt.Sprintf("%q", v.Value)
	}
	return fmt.Sprintf("Token [type: %s, line: %d, position: %d]: %s", v.Type, v.Line, v.Position, s)
}

// SCANNER

// This constructor creates a new scanner initialized with the specified array
// of bytes. The scanner will scan in tokens matching Go primitive types.
func Scanner(source []byte, tokens chan Token) *scanner {
	var v = &scanner{source: source, line: 1, position: 1, tokens: tokens}
	go v.scanTokens() // Start scanning in the background.
	return v
}

// This type defines the structure and methods for the scanner agent. The source
// bytes can be viewed like this:
//
//   | byte 0 | byte 1 | byte 2 | byte 3 | byte 4 | byte 5 | ... | byte N-1 |
//   | rune 0 |      rune 1     |      rune 2     | rune 3 | ... | rune R-1 |
//
// Runes can be one to eight bytes long.

type scanner struct {
	source    []byte
	firstByte int // The zero based index of the first possible byte in the next token.
	nextByte  int // The zero based index of the next possible byte in the next token.
	line      int // The line number in the source string of the next rune.
	position  int // The position in the current line of the first rune in the next token.
	tokens    chan Token
}

// This method continues scanning tokens from the source array until an error
// occurs or the end of file is reached. It then closes the token channel.
func (v *scanner) scanTokens() {
	for v.scanToken() {
	}
	close(v.tokens)
}

// This method attempts to scan any token starting with the next rune in the
// source array. It checks for each type of token as the cases for the switch
// statement. If that token type is found, this method returns true and skips
// the rest of the cases.  If no valid token is found, or a TokenEOF is found
// this method returns false.
func (v *scanner) scanToken() bool {
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
	case v.foundInteger():  // Must be after all other numeric types.
	default:
		// No valid token was found.
		v.foundError()
		return false
	}
	return true
}

// This method scans through any spaces in the source array and sets the next
// byte index to the next non-white-space rune.
func (v *scanner) skipSpaces() {
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

// This method adds a token of the specified type with the current scanner
// information to the token channel. It then resets the first byte index to the
// next byte index position. It returns the token type of the type added to the
// channel.
func (v *scanner) emitToken(tType TokenType) TokenType {
	var tValue = string(v.source[v.firstByte:v.nextByte])
	if tType == TokenEOF {
		tValue = "<EOF>"
	}
	if tType == TokenError {
		switch tValue {
		case "\a":
			tValue = "<BELL>"
		case "\b":
			tValue = "<BKSP>"
		case "\t":
			tValue = "<TAB>"
		case "\f":
			tValue = "<FF>"
		case "\r":
			tValue = "<CR>"
		case "\v":
			tValue = "<VTAB>"
		}
	}
	var token = Token{tType, tValue, v.line, v.position}
	//fmt.Println(token)
	v.tokens <- token
	v.firstByte = v.nextByte
	v.position += sts.Count(tValue, "") - 1 // Add the number of runes in the token.
	return tType
}

// This method adds a boolean token with the current scanner information to the
// token channel. It returns true if a boolean token was found.
func (v *scanner) foundBoolean() bool {
	var s = v.source[v.nextByte:]
	var matches = scanBoolean(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenBoolean)
		return true
	}
	return false
}

// This method adds a context token with the current scanner information
// to the token channel. It returns true if a type token was found.
func (v *scanner) foundContext() bool {
	var s = v.source[v.nextByte:]
	var matches = scanContext(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenContext)
		return true
	}
	return false
}

// This method adds a complex number token with the current scanner information
// to the token channel. It returns true if a complex number token was found.
func (v *scanner) foundComplex() bool {
	var s = v.source[v.nextByte:]
	var matches = scanComplex(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenComplex)
		return true
	}
	return false
}

// This method adds a delimiter token with the current scanner information to the
// token channel. It returns true if a delimiter token was found.
func (v *scanner) foundDelimiter() bool {
	var s = v.source[v.nextByte:]
	var matches = scanDelimiter(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenDelimiter)
		return true
	}
	return false
}

// This method adds an error token with the current scanner information to the token
// channel.
func (v *scanner) foundError() {
	var bytes = v.source[v.nextByte:]
	var _, width = utf.DecodeRune(bytes)
	v.nextByte += width
	v.emitToken(TokenError)
}

// This method adds an EOF token with the current scanner information to the token
// channel. It returns true if an EOF marker was found.
func (v *scanner) foundEOF() bool {
	// The last byte in a POSIX standard file must be an EOL character.
	var s = v.source[v.nextByte:]
	if !byt.HasPrefix(s, []byte(EOL)) {
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
	v.emitToken(TokenEOF)
	return true
}

// This method adds an EOL token with the current scanner information to the token
// channel. It returns true if an EOL token was found.
func (v *scanner) foundEOL() bool {
	var s = v.source[v.nextByte:]
	if !byt.HasPrefix(s, []byte(EOL)) {
		return false
	}
	v.nextByte++
	v.emitToken(TokenEOL)
	v.line++
	v.position = 1
	return true
}

// This method adds a floating point token with the current scanner information
// to the token channel. It returns true if a floating point token was found.
func (v *scanner) foundFloat() bool {
	var s = v.source[v.nextByte:]
	var matches = scanFloat(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenFloat)
		return true
	}
	return false
}

// This method adds a integer token with the current scanner information to the
// token channel. It returns true if a integer token was found.
func (v *scanner) foundInteger() bool {
	var s = v.source[v.nextByte:]
	var matches = scanInteger(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenInteger)
		return true
	}
	return false
}

// This method adds a nil token with the current scanner information to the
// token channel. It returns true if a nil token was found.
func (v *scanner) foundNil() bool {
	var s = v.source[v.nextByte:]
	var matches = scanNil(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenNil)
		return true
	}
	return false
}

// This method adds a rune token with the current scanner information to the
// token channel. It returns true if a rune token was found.
func (v *scanner) foundRune() bool {
	var s = v.source[v.nextByte:]
	var matches = scanRune(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenRune)
		return true
	}
	return false
}

// This method adds a string token with the current scanner information to the
// token channel. It returns true if a string token was found.
func (v *scanner) foundString() bool {
	var s = v.source[v.nextByte:]
	var matches = scanString(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenString)
		return true
	}
	return false
}

// This method adds an unsigned integer token with the current scanner
// information to the token channel. It returns true if an unsigned integer
// token was found.
func (v *scanner) foundUnsigned() bool {
	var s = v.source[v.nextByte:]
	var matches = scanUnsigned(s)
	if len(matches) > 0 {
		v.nextByte += len(matches[0])
		v.emitToken(TokenUnsigned)
		return true
	}
	return false
}

// This scanner is used for matching boolean primitives.
var booleanScanner = reg.MustCompile(`^(?:` + boolean + `)`)

// This function returns for the specified string an array of the matching
// subgroups for a boolean primitive. The first string in the array is the
// entire matched string.
func scanBoolean(v []byte) []string {
	return bytesToStrings(booleanScanner.FindSubmatch(v))
}

// This scanner is used for matching contexts.
var contextScanner = reg.MustCompile(`^(?:` + context + `)`)

// This function returns for the specified string an array of the matching
// subgroups for a context. The first string in the array is the
// entire matched string.
func scanContext(v []byte) []string {
	return bytesToStrings(contextScanner.FindSubmatch(v))
}

// This scanner is used for matching complex primitives.
var complexScanner = reg.MustCompile(`^(?:` + complex_ + `)`)

// This function returns for the specified string an array of the matching
// subgroups for a complex primitive. The first string in the array is the
// entire matched string.
func scanComplex(v []byte) []string {
	return bytesToStrings(complexScanner.FindSubmatch(v))
}

// This function returns for the specified string an array of the matching
// subgroups for a delimiter. The first string in the array is the entire
// matched string.
func scanDelimiter(v []byte) []string {
	var result []string
	for _, delimiter := range delimiters {
		if byt.HasPrefix(v, delimiter) {
			result = append(result, string(delimiter))
		}
	}
	return result
}

// This scanner is used for matching float primitives.
var floatScanner = reg.MustCompile(`^(?:` + float + `)`)

// This function returns for the specified string an array of the matching
// subgroups for a float primitive. The first string in the array is the
// entire matched string.
func scanFloat(v []byte) []string {
	return bytesToStrings(floatScanner.FindSubmatch(v))
}

// This scanner is used for matching integer primitives.
var integerScanner = reg.MustCompile(`^(?:` + integer + `)`)

// This function returns for the specified string an array of the matching
// subgroups for an integer primitive. The first string in the array is the
// entire matched string.
func scanInteger(v []byte) []string {
	return bytesToStrings(integerScanner.FindSubmatch(v))
}

// This scanner is used for matching nil primitives.
var nilScanner = reg.MustCompile(`^(?:` + nil_ + `)`)

// This function returns for the specified string an array of the matching
// subgroups for a nil primitive. The first string in the array is the
// entire matched string.
func scanNil(v []byte) []string {
	return bytesToStrings(nilScanner.FindSubmatch(v))
}

// This scanner is used for matching rune primitives.
var runeScanner = reg.MustCompile(`^(?:` + rune_ + `)`)

// This function returns for the specified string an array of the matching
// subgroups for a rune primitive. The first string in the array is the
// entire matched string.
func scanRune(v []byte) []string {
	return bytesToStrings(runeScanner.FindSubmatch(v))
}

// This scanner is used for matching string primitives.
var stringScanner = reg.MustCompile(`^(?:` + string_ + `)`)

// This function returns for the specified string an array of the matching
// subgroups for a string primitive. The first string in the array is the
// entire matched string.
func scanString(v []byte) []string {
	return bytesToStrings(stringScanner.FindSubmatch(v))
}

// This scanner is used for matching unsigned integer primitives.
var unsignedScanner = reg.MustCompile(`^(?:` + unsigned + `)`)

// This function returns for the specified string an array of the matching
// subgroups for an unsigned integer primitive. The first string in the array is
// the entire matched string.
func scanUnsigned(v []byte) []string {
	return bytesToStrings(unsignedScanner.FindSubmatch(v))
}

// CONSTANT DEFINITIONS

// These constants define the POSIX standard representations.
const (
	EOF = "\n"  // Must be last byte in a file.
	EOL = "\n"
)

// These constant definitions capture regular expression subpatterns.
const (
	base10    = `[0-9]`
	base16    = `[0-9a-f]`
	boolean   = `false|true`
	complex_  = `\((` + float + `)` + sign + `(` + float + `)i\)`
	context   = `array|catalog|list|map|queue|set|stack`
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

// This array contains the set of delimiters that may be used to separate the
// Go primitive types.
var delimiters = [][]byte{
	[]byte("]"),
	[]byte("["),
	[]byte(")"),
	[]byte("("),
	[]byte(":"),
	[]byte(","),
}

// PRIVATE FUNCTIONS

func bytesToStrings(bytes [][]byte) []string {
	var strings = make([]string, len(bytes))
	for index, array := range bytes {
		strings[index] = string(array)
	}
	return strings
}
