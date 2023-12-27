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
	//fmt "fmt"
	reg "regexp"
)

// CLASS NAMESPACE

// Private Class Namespace Type

type scannerClass_ struct {
	booleanMatcher   *reg.Regexp
	complexMatcher   *reg.Regexp
	delimiterMatcher *reg.Regexp
	eolMatcher       *reg.Regexp
	floatMatcher     *reg.Regexp
	integerMatcher   *reg.Regexp
	nilMatcher       *reg.Regexp
	runeMatcher      *reg.Regexp
	spaceMatcher     *reg.Regexp
	stringMatcher    *reg.Regexp
	typeMatcher      *reg.Regexp
	unsignedMatcher  *reg.Regexp
}

// Private Class Namespace Reference

var scannerClass = &scannerClass_{
	booleanMatcher:   reg.MustCompile(`^(?:` + boolean_ + `)`),
	complexMatcher:   reg.MustCompile(`^(?:` + complex_ + `)`),
	delimiterMatcher: reg.MustCompile(`^(?:` + delimiter_ + `)`),
	eolMatcher:       reg.MustCompile(`^(?:` + eol_ + `)`),
	floatMatcher:     reg.MustCompile(`^(?:` + float_ + `)`),
	integerMatcher:   reg.MustCompile(`^(?:` + integer_ + `)`),
	nilMatcher:       reg.MustCompile(`^(?:` + nil_ + `)`),
	runeMatcher:      reg.MustCompile(`^(?:` + rune_ + `)`),
	spaceMatcher:     reg.MustCompile(`^(?:` + space_ + `)`),
	stringMatcher:    reg.MustCompile(`^(?:` + string_ + `)`),
	typeMatcher:      reg.MustCompile(`^(?:` + type_ + `)`),
	unsignedMatcher:  reg.MustCompile(`^(?:` + unsigned_ + `)`),
}

// Public Class Namespace Access

func ScannerClass() *scannerClass_ {
	return scannerClass
}

// Public Class Constants

func (c *scannerClass_) GetBooleanMatcher() *reg.Regexp {
	return c.booleanMatcher
}

func (c *scannerClass_) GetComplexMatcher() *reg.Regexp {
	return c.complexMatcher
}

func (c *scannerClass_) GetDelimiterMatcher() *reg.Regexp {
	return c.delimiterMatcher
}

func (c *scannerClass_) GetEOLMatcher() *reg.Regexp {
	return c.eolMatcher
}

func (c *scannerClass_) GetFloatMatcher() *reg.Regexp {
	return c.floatMatcher
}

func (c *scannerClass_) GetIntegerMatcher() *reg.Regexp {
	return c.integerMatcher
}

func (c *scannerClass_) GetNilMatcher() *reg.Regexp {
	return c.nilMatcher
}

func (c *scannerClass_) GetRuneMatcher() *reg.Regexp {
	return c.runeMatcher
}

func (c *scannerClass_) GetStringMatcher() *reg.Regexp {
	return c.stringMatcher
}

func (c *scannerClass_) GetTypeMatcher() *reg.Regexp {
	return c.typeMatcher
}

func (c *scannerClass_) GetUnsignedMatcher() *reg.Regexp {
	return c.unsignedMatcher
}

// Public Class Constructors

func (c *scannerClass_) FromSource(
	source string,
	tokens chan *token_,
) *scanner_ {
	var scanner = &scanner_{
		line:     1,
		position: 1,
		runes:    []rune(source),
		tokens:   tokens,
	}
	go scanner.scanTokens() // Start scanning tokens in the background.
	return scanner
}

// Public Class Functions

// This public class function returns, for the specified text, a Go array of
// the matching subgroups for a Go boolean primitive. The first string in the Go
// array is the entire matched string.
func (c *scannerClass_) MatchBoolean(text string) []string {
	return c.booleanMatcher.FindStringSubmatch(text)
}

// This public class function returns, for the specified text, a Go array of
// the matching subgroups for a Go complex primitive. The first string in the Go
// array is the entire matched string.
func (c *scannerClass_) MatchComplex(text string) []string {
	return c.complexMatcher.FindStringSubmatch(text)
}

// This public class function returns, for the specified text, a Go array of
// the matching subgroups for a delimiter.
func (c *scannerClass_) MatchDelimiter(text string) []string {
	return c.delimiterMatcher.FindStringSubmatch(text)
}

// This public class function returns, for the specified text, a Go array of
// the matching subgroups for an end-of-line (EOL).
func (c *scannerClass_) MatchEOL(text string) []string {
	return c.eolMatcher.FindStringSubmatch(text)
}

// This public class function returns, for the specified text, a Go array of
// the matching subgroups for a Go float primitive. The first string in the Go
// array is the entire matched string.
func (c *scannerClass_) MatchFloat(text string) []string {
	return c.floatMatcher.FindStringSubmatch(text)
}

// This public class function returns, for the specified text, a Go array of
// the matching subgroups for a Go integer primitive. The first string in the Go
// array is the entire matched string.
func (c *scannerClass_) MatchInteger(text string) []string {
	return c.integerMatcher.FindStringSubmatch(text)
}

// This public class function returns, for the specified text, a Go array of
// the matching subgroups for a Go nil primitive. The first string in the Go
// array is the entire matched string.
func (c *scannerClass_) MatchNil(text string) []string {
	return c.nilMatcher.FindStringSubmatch(text)
}

// This public class function returns, for the specified text, a Go array of
// the matching subgroups for a Go rune primitive. The first string in the Go
// array is the entire matched string.
func (c *scannerClass_) MatchRune(text string) []string {
	return c.runeMatcher.FindStringSubmatch(text)
}

// This public class function returns, for the specified text, a Go array of
// the matching subgroups for a Go string primitive. The first string in the Go
// array is the entire matched string.
func (c *scannerClass_) MatchString(text string) []string {
	return c.stringMatcher.FindStringSubmatch(text)
}

// This public class function returns, for the specified text, a Go array of
// the matching subgroups for a collection type. The first string in the Go
// array is the entire matched string.
func (c *scannerClass_) MatchType(text string) []string {
	return c.typeMatcher.FindStringSubmatch(text)
}

// This public class function returns, for the specified text, a Go array of
// the matching subgroups for a Go unsigned primitive. The first string in the
// Go array is the entire matched string.
func (c *scannerClass_) MatchUnsigned(text string) []string {
	return c.unsignedMatcher.FindStringSubmatch(text)
}

// CLASS INSTANCES

// Private Class Type Definition

type scanner_ struct {
	first    int // A zero based index of the first possible rune in the next token.
	line     int // The line number in the source text of the next rune.
	next     int // A zero based index of the next possible rune in the next token.
	position int // The position in the current line of the next rune.
	runes    []rune
	tokens   chan *token_
}

// Private Interface

// This private class method adds a token of the specified type with the current
// scanner information to the token channel. It then resets the first rune index
// to the next rune index position. It returns the token type of the type added
// to the channel.
func (v *scanner_) emitToken(tokenType string) string {
	var tokenValue = string(v.runes[v.first:v.next])
	switch tokenValue {
	case "\a":
		tokenValue = "<BELL>"
	case "\b":
		tokenValue = "<BKSP>"
	case "\t":
		tokenValue = "<TAB>"
	case "\f":
		tokenValue = "<FF>"
	case "\n":
		tokenValue = "<EOL>"
	case "\r":
		tokenValue = "<CR>"
	case "\v":
		tokenValue = "<VTAB>"
	}
	var token = TokenClass().FromContext(v.line, v.position, tokenType, tokenValue)
	//fmt.Println(token) // Uncomment when debugging.
	v.tokens <- token
	v.position += v.next - v.first
	v.first = v.next
	return tokenType
}

// This private class method adds a boolean token with the current Scanner
// information to the token channel. It returns true if a boolean token was found.
func (v *scanner_) foundBoolean() bool {
	var text = string(v.runes[v.next:])
	var matches = scannerClass.booleanMatcher.FindStringSubmatch(text)
	if len(matches) > 0 {
		v.next += len([]rune(matches[0]))
		v.emitToken(TokenClass().GetBoolean())
		return true
	}
	return false
}

// This private class method adds a complex number token with the current
// scanner information to the token channel. It returns true if a complex number
// token was found.
func (v *scanner_) foundComplex() bool {
	var text = string(v.runes[v.next:])
	var matches = scannerClass.complexMatcher.FindStringSubmatch(text)
	if len(matches) > 0 {
		v.next += len([]rune(matches[0]))
		v.emitToken(TokenClass().GetComplex())
		return true
	}
	return false
}

// This private class method adds a delimiter token with the current Scanner
// information to the token channel. It returns true if a delimiter token was
// found.
func (v *scanner_) foundDelimiter() bool {
	var text = string(v.runes[v.next:])
	var matches = scannerClass.delimiterMatcher.FindStringSubmatch(text)
	if len(matches) > 0 {
		v.next += len([]rune(matches[0]))
		v.emitToken(TokenClass().GetDelimiter())
		return true
	}
	return false
}

// This private class method adds an EOL token with the current Scanner
// information to the token channel. It returns true if an EOL token was found.
func (v *scanner_) foundEOL() bool {
	var text = string(v.runes[v.next:])
	var matches = scannerClass.eolMatcher.FindStringSubmatch(text)
	if len(matches) > 0 {
		v.next += len([]rune(matches[0]))
		v.emitToken(TokenClass().GetEOL())
		v.line++
		v.position = 1
		return true
	}
	return false
}

// This private class method adds an error token with the current Scanner
// information to the token channel. It always returns true.
func (v *scanner_) foundError() {
	v.next++
	v.emitToken(TokenClass().GetError())
}

// This private class method adds a floating point token with the current
// scanner information to the token channel. It returns true if a floating point
// token was found.
func (v *scanner_) foundFloat() bool {
	var text = string(v.runes[v.next:])
	var matches = scannerClass.floatMatcher.FindStringSubmatch(text)
	if len(matches) > 0 {
		v.next += len([]rune(matches[0]))
		v.emitToken(TokenClass().GetFloat())
		return true
	}
	return false
}

// This private class method adds a integer token with the current Scanner
// information to the token channel. It returns true if a integer token was
// found.
func (v *scanner_) foundInteger() bool {
	var text = string(v.runes[v.next:])
	var matches = scannerClass.integerMatcher.FindStringSubmatch(text)
	if len(matches) > 0 {
		v.next += len([]rune(matches[0]))
		v.emitToken(TokenClass().GetInteger())
		return true
	}
	return false
}

// This private class method adds a nil token with the current Scanner
// information to the token channel. It returns true if a nil token was found.
func (v *scanner_) foundNil() bool {
	var text = string(v.runes[v.next:])
	var matches = scannerClass.nilMatcher.FindStringSubmatch(text)
	if len(matches) > 0 {
		v.next += len([]rune(matches[0]))
		v.emitToken(TokenClass().GetNil())
		return true
	}
	return false
}

// This private class method adds a rune token with the current Scanner
// information to the token channel. It returns true if a rune token was found.
func (v *scanner_) foundRune() bool {
	var text = string(v.runes[v.next:])
	var matches = scannerClass.runeMatcher.FindStringSubmatch(text)
	if len(matches) > 0 {
		v.next += len([]rune(matches[0]))
		v.emitToken(TokenClass().GetRune())
		return true
	}
	return false
}

// This private class method moves the scanner past any spaces in the source
// runes.  It does not add any tokens to the token channel.  It returns true if
// a string token was found.
func (v *scanner_) foundSpace() bool {
	var text = string(v.runes[v.next:])
	var matches = scannerClass.spaceMatcher.FindStringSubmatch(text)
	if len(matches) > 0 {
		v.next += len([]rune(matches[0]))
		v.position += v.next - v.first
		v.first = v.next
		return true
	}
	return false
}

// This private class method adds a string token with the current Scanner
// information to the token channel. It returns true if a string token was
// found.
func (v *scanner_) foundString() bool {
	var text = string(v.runes[v.next:])
	var matches = scannerClass.stringMatcher.FindStringSubmatch(text)
	if len(matches) > 0 {
		v.next += len([]rune(matches[0]))
		v.emitToken(TokenClass().GetString())
		return true
	}
	return false
}

// This private class method adds a collection type token with the current
// scanner information to the token channel. It returns true if a collection
// type token was found.
func (v *scanner_) foundType() bool {
	var text = string(v.runes[v.next:])
	var matches = scannerClass.typeMatcher.FindStringSubmatch(text)
	if len(matches) > 0 {
		v.next += len([]rune(matches[0]))
		v.emitToken(TokenClass().GetType())
		return true
	}
	return false
}

// This private class method adds an unsigned integer token with the current
// scanner information to the token channel. It returns true if an unsigned
// integer token was found.
func (v *scanner_) foundUnsigned() bool {
	var text = string(v.runes[v.next:])
	var matches = scannerClass.unsignedMatcher.FindStringSubmatch(text)
	if len(matches) > 0 {
		v.next += len([]rune(matches[0]))
		v.emitToken(TokenClass().GetUnsigned())
		return true
	}
	return false
}

// This private class method continues scanning tokens from the source text
// until the of the source text is reached.  It attempts to scan each token
// type until a matching token is found.  That token is then added to the token
// channel to be read by the parser.  If no valid token type is found an error
// token is added to the channel.  Once an error is encountered, or the end of
// the source text is reached, an end-of-file (EOF) token is added to the
// token channel and the channel is closed.
func (v *scanner_) scanTokens() {
loop:
	for v.next < len(v.runes) {
		switch {
		case v.foundSpace(): // Remove any whitespace.
		case v.foundBoolean():
		case v.foundComplex():
		case v.foundDelimiter():
		case v.foundEOL():
		case v.foundFloat():
		case v.foundNil():
		case v.foundRune():
		case v.foundString():
		case v.foundType():
		case v.foundUnsigned():
		case v.foundInteger(): // Must be after the other numeric types.
		default:
			v.foundError()
			break loop
		}
	}
	v.emitToken(TokenClass().GetEOF())
	close(v.tokens)
}

// These private constants define the regular expression sub-patterns that make
// up all token types.  Unfortunately there is no way to make them private to
// the scanner class namespace since they must be TRUE Go constants to be
// initialized in this way.  We add an underscore to lessen the chance of a name
// collision with other private Go class constants.
const (
	base10_    = `[0-9]`
	base16_    = `[0-9a-f]`
	boolean_   = `false|true`
	complex_   = `\((` + float_ + `)` + sign_ + `(` + float_ + `)i\)`
	delimiter_ = `\[|\]|\(|\)|:|,`
	eol_       = `\n`
	escape_    = `\\(?:(?:` + unicode_ + `)|[abfnrtv'"\\])`
	exponent_  = `[eE]` + sign_ + ordinal_
	float_     = sign_ + `?(?:` + scalar_ + `)(?:` + exponent_ + `)?`
	fraction_  = `\.` + base10_ + `+`
	integer_   = zero_ + `|` + sign_ + `?` + ordinal_
	nil_       = `nil`
	ordinal_   = `[1-9][0-9]*`
	rune_      = `'(` + escape_ + `|[^'` + eol_ + `])'`
	scalar_    = `(?:` + zero_ + `|` + ordinal_ + `)` + fraction_
	sign_      = `[+-]`
	space_     = `[ ]+`
	string_    = `"(` + escape_ + `|[^"` + eol_ + `])*"`
	type_      = `[Aa]rray|Catalog|List|[Mm]ap|Queue|Set|Stack`
	unicode_   = `x` + base16_ + `{2}|u` + base16_ + `{4}|U` + base16_ + `{8}`
	unsigned_  = `0x` + base16_ + `+`
	zero_      = `0`
)
