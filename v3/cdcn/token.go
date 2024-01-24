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
	fmt "fmt"
	col "github.com/craterdog/go-collection-framework/v3"
)

// CLASS NAMESPACE

// Private Class Namespace Type

type tokenClass_ struct {
	strings map[col.TokenType]string
}

// Private Class Namespace Reference

var tokenClass = &tokenClass_{
	strings: map[col.TokenType]string{
		TypeBoolean:     "Boolean",
		TypeComplex:     "Complex",
		TypeContext:     "Context",
		TypeDelimiter:   "Delimiter",
		TypeEOF:         "EOF",
		TypeEOL:         "EOL",
		TypeError:       "Error",
		TypeFloat:       "Float",
		TypeHexadecimal: "Hexadecimal",
		TypeInteger:     "Integer",
		TypeNil:         "Nil",
		TypeRune:        "Rune",
		TypeSpace:       "Space",
		TypeString:      "String",
	},
}

// Public Class Namespace Access

func TokenClass() col.TokenClassLike {
	return tokenClass
}

// Public Class Constructors

func (c *tokenClass_) Make(
	line int,
	position int,
	tokenType col.TokenType,
	tokenValue string,
) col.TokenLike {
	var token = &token_{
		line:       line,
		position:   position,
		tokenType:  tokenType,
		tokenValue: tokenValue,
	}
	return token
}

// Public Class Functions

func (c *tokenClass_) AsString(tokenType col.TokenType) string {
	return c.strings[tokenType]
}

// CLASS INSTANCES

// Private Class Type Definition

type token_ struct {
	line       int // The line number of the token in the source string.
	position   int // The position in the line of the first rune of the token.
	tokenType  col.TokenType
	tokenValue string
}

// Stringer Interface

func (v *token_) String() string {
	var s = fmt.Sprintf("%q", v.tokenValue)
	if len(s) > 40 {
		s = fmt.Sprintf("%.40q...", v.tokenValue)
	}
	return fmt.Sprintf("Token [type: %s, line: %d, position: %d]: %s",
		TokenClass().AsString(v.tokenType),
		v.line,
		v.position,
		s,
	)
}

// Public Interface

func (v *token_) GetLine() int {
	return v.line
}

func (v *token_) GetPosition() int {
	return v.position
}

func (v *token_) GetType() col.TokenType {
	return v.tokenType
}

func (v *token_) GetValue() string {
	return v.tokenValue
}
