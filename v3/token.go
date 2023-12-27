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
	fmt "fmt"
)

// CLASS NAMESPACE

// Private Class Namespace Type

type tokenClass_ struct {
	boolean_   string
	complex_   string
	delimiter_ string
	eof_       string
	eol_       string
	error_     string
	float_     string
	integer_   string
	nil_       string
	rune_      string
	string_    string
	type_      string
	unsigned_  string
}

// Private Class Namespace Reference

var tokenClass = &tokenClass_{
	boolean_:   "Boolean",
	complex_:   "Complex",
	delimiter_: "Delimiter",
	eof_:       "EOF",
	eol_:       "EOL",
	error_:     "Error",
	float_:     "Float",
	integer_:   "Integer",
	nil_:       "Nil",
	rune_:      "Rune",
	string_:    "String",
	type_:      "Type",
	unsigned_:  "Unsigned",
}

// Public Class Namespace Access

func TokenClass() *tokenClass_ {
	return tokenClass
}

// Public Class Constants

func (c *tokenClass_) GetBoolean() string {
	return c.boolean_
}

func (c *tokenClass_) GetComplex() string {
	return c.complex_
}

func (c *tokenClass_) GetDelimiter() string {
	return c.delimiter_
}

func (c *tokenClass_) GetEOF() string {
	return c.eof_
}

func (c *tokenClass_) GetEOL() string {
	return c.eol_
}

func (c *tokenClass_) GetError() string {
	return c.error_
}

func (c *tokenClass_) GetFloat() string {
	return c.float_
}

func (c *tokenClass_) GetInteger() string {
	return c.integer_
}

func (c *tokenClass_) GetNil() string {
	return c.nil_
}

func (c *tokenClass_) GetRune() string {
	return c.rune_
}

func (c *tokenClass_) GetString() string {
	return c.string_
}

func (c *tokenClass_) GetType() string {
	return c.type_
}

func (c *tokenClass_) GetUnsigned() string {
	return c.unsigned_
}

// Public Class Constructors

func (c *tokenClass_) FromContext(
	line int,
	position int,
	tokenType string,
	tokenValue string,
) *token_ {
	var token = &token_{
		line_:     line,
		position_: position,
		type_:     tokenType,
		value_:    tokenValue,
	}
	return token
}

// CLASS INSTANCES

// Private Class Type Definition

type token_ struct {
	line_     int // The line number of the token in the lexical context.
	position_ int // The position in the line of the first rune of the token.
	type_     string
	value_    string
}

// Private Interface

func (v *token_) GetLine() int {
	return v.line_
}

func (v *token_) GetPosition() int {
	return v.position_
}

func (v *token_) GetType() string {
	return v.type_
}

func (v *token_) GetValue() string {
	return v.value_
}

// This public class method returns the canonical string version of this token.
func (v *token_) String() string {
	var s string
	switch {
	case v.type_ == TokenClass().eof_:
		s = "<EOF>"
	case v.type_ == TokenClass().eol_:
		s = "<EOL>"
	case len(v.value_) > 60:
		s = fmt.Sprintf("%.60q...", v.value_)
	default:
		s = fmt.Sprintf("%q", v.value_)
	}
	return fmt.Sprintf(
		"Token [type: %s, line: %d, position: %d]: %s",
		v.type_, v.line_, v.position_, s)
}
