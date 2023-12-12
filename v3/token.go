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
	fmt "fmt"
)

// CLASS NAMESPACE

// This private type defines the namespace structure associated with the
// constants, constructors and functions for the token class namespace.
type tokenClass_ struct {
	typeError     string
	typeBoolean   string
	typeComplex   string
	typeContext   string
	typeDelimiter string
	typeEOF       string
	typeEOL       string
	typeFloat     string
	typeInteger   string
	typeNil       string
	typeRune      string
	typeString    string
	typeUnsigned  string
}

// This private constant defines the singleton reference to the token
// class namespace.  It also initializes any class constants as needed.
var tokenClassSingleton = &tokenClass_{
	typeError:     "Error",
	typeBoolean:   "Boolean",
	typeComplex:   "Complex",
	typeContext:   "Context",
	typeDelimiter: "Delimiter",
	typeEOF:       "EOF",
	typeEOL:       "EOL",
	typeFloat:     "Float",
	typeInteger:   "Integer",
	typeNil:       "Nil",
	typeRune:      "Rune",
	typeString:    "String",
	typeUnsigned:  "Unsigned",
}

// This public function returns the singleton reference to the token
// class namespace.
func Token() *tokenClass_ {
	return tokenClassSingleton
}

// CLASS CONSTANTS

// This public class constant represents an error token type.
func (c *tokenClass_) TypeError() string {
	return c.typeError
}

// This public class constant represents a boolean token type.
func (c *tokenClass_) TypeBoolean() string {
	return c.typeBoolean
}

// This public class constant represents a complex token type.
func (c *tokenClass_) TypeComplex() string {
	return c.typeComplex
}

// This public class constant represents a context token type.
func (c *tokenClass_) TypeContext() string {
	return c.typeContext
}

// This public class constant represents a delimiter token type.
func (c *tokenClass_) TypeDelimiter() string {
	return c.typeDelimiter
}

// This public class constant represents an end-of-file token type.
func (c *tokenClass_) TypeEOF() string {
	return c.typeEOF
}

// This public class constant represents an end-of-line token type.
func (c *tokenClass_) TypeEOL() string {
	return c.typeEOL
}

// This public class constant represents a float token type.
func (c *tokenClass_) TypeFloat() string {
	return c.typeFloat
}

// This public class constant represents an integer token type.
func (c *tokenClass_) TypeInteger() string {
	return c.typeInteger
}

// This public class constant represents a nil token type.
func (c *tokenClass_) TypeNil() string {
	return c.typeNil
}

// This public class constant represents a rune token type.
func (c *tokenClass_) TypeRune() string {
	return c.typeRune
}

// This public class constant represents a string token type.
func (c *tokenClass_) TypeString() string {
	return c.typeString
}

// This public class constant represents an unsigned token type.
func (c *tokenClass_) TypeUnsigned() string {
	return c.typeUnsigned
}

// CLASS CONSTRUCTORS

// This public class constructor creates a new token from the specified
// lexical context.
func (c *tokenClass_) FromContext(
	tokenType string,
	tokenValue string,
	line int,
	position int,
) TokenLike {
	var token = &token_{
		tokenType:  tokenType,
		tokenValue: tokenValue,
		line:       line,
		position:   position,
	}
	return token
}

// CLASS TYPE

// Encapsulated Type

// This private class type encapsulates a Go structure containing private
// attributes that can only be accessed and manipulated using methods that
// implement the token-like abstract type.  A token captures lexical context.
type token_ struct {
	tokenType  string
	tokenValue string
	line       int // The line number of the token in the lexical context.
	position   int // The position in the line of the first rune of the token.
}

// Contextual Interface

// This public class method returns the type of this token.
func (v *token_) GetType() string {
	return v.tokenType
}

// This public class method returns the value of this token.
func (v *token_) GetValue() string {
	return v.tokenValue
}

// This public class method returns the line number of this token.
func (v *token_) GetLine() int {
	return v.line
}

// This public class method returns the position of this token in the line.
func (v *token_) GetPosition() int {
	return v.position
}

// Private Interface

// This public class method returns the canonical string version of this token.
func (v *token_) String() string {
	var Token = Token()
	var s string
	switch {
	case v.tokenType == Token.typeEOF:
		s = "<EOF>"
	case v.tokenType == Token.typeEOL:
		s = "<EOL>"
	case len(v.tokenValue) > 60:
		s = fmt.Sprintf("%.60q...", v.tokenValue)
	default:
		s = fmt.Sprintf("%q", v.tokenValue)
	}
	return fmt.Sprintf(
		"Token [type: %s, line: %d, position: %d]: %s",
		v.tokenType, v.line, v.position, s)
}
