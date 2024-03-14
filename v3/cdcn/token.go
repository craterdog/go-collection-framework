/*
................................................................................
.    Copyright (c) 2009-2024 Crater Dog Technologies™.  All Rights Reserved.   .
................................................................................
.  DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               .
.                                                                              .
.  This code is free software; you can redistribute it and/or modify it under  .
.  the terms of The MIT License (MIT), as published by the Open Source         .
.  Initiative. (See http://opensource.org/licenses/MIT)                        .
................................................................................
*/

package cdcn

import (
	fmt "fmt"
)

// CLASS ACCESS

// Reference

var tokenClass = &tokenClass_{
	strings_: map[TokenType]string{
		BooleanToken:     "Boolean",
		ComplexToken:     "Complex",
		ContextToken:     "Context",
		DelimiterToken:   "Delimiter",
		EOFToken:         "EOF",
		EOLToken:         "EOL",
		ErrorToken:       "Error",
		FloatToken:       "Float",
		HexadecimalToken: "Hexadecimal",
		IntegerToken:     "Integer",
		NilToken:         "Nil",
		RuneToken:        "Rune",
		SpaceToken:       "Space",
		StringToken:      "String",
	},
}

// Function

func Token() TokenClassLike {
	return tokenClass
}

// CLASS METHODS

// Target

type tokenClass_ struct {
	strings_ map[TokenType]string
}

// Constructors

func (c *tokenClass_) MakeWithAttributes(
	line int,
	position int,
	type_ TokenType,
	value string,
) TokenLike {
	return &token_{
		line_:     line,
		position_: position,
		type_:     type_,
		value_:    value,
	}
}

// Functions

func (c *tokenClass_) AsString(type_ TokenType) string {
	return c.strings_[type_]
}

// INSTANCE METHODS

// Target

type token_ struct {
	line_     int // The line number of the token in the source string.
	position_ int // The position in the line of the first rune of the token.
	type_     TokenType
	value_    string
}

// Stringer

func (v *token_) String() string {
	var s = fmt.Sprintf("%q", v.value_)
	if len(s) > 40 {
		s = fmt.Sprintf("%.40q...", v.value_)
	}
	return fmt.Sprintf("Token [type: %s, line: %d, position: %d]: %s",
		Token().AsString(v.type_),
		v.line_,
		v.position_,
		s,
	)
}

// Public

func (v *token_) GetLine() int {
	return v.line_
}

func (v *token_) GetPosition() int {
	return v.position_
}

func (v *token_) GetType() TokenType {
	return v.type_
}

func (v *token_) GetValue() string {
	return v.value_
}
