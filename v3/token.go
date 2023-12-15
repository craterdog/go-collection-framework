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
// constants, constructors and functions for the Token class namespace.
type tokenClass_ struct {
	typeBoolean   string
	typeComplex   string
	typeContext   string
	typeDelimiter string
	typeEOF       string
	typeEOL       string
	typeError     string
	typeFloat     string
	typeInteger   string
	typeNil       string
	typeRune      string
	typeString    string
	typeUnsigned  string
}

// This private constant defines the singleton reference to the Token
// class namespace.  It also initializes any class constants as needed.
var tokenClassSingleton = &tokenClass_{
	typeBoolean:   "Boolean",
	typeComplex:   "Complex",
	typeContext:   "Context",
	typeDelimiter: "Delimiter",
	typeEOF:       "EOF",
	typeEOL:       "EOL",
	typeError:     "Error",
	typeFloat:     "Float",
	typeInteger:   "Integer",
	typeNil:       "Nil",
	typeRune:      "Rune",
	typeString:    "String",
	typeUnsigned:  "Unsigned",
}

// This public function returns the singleton reference to the Token
// class namespace.
func Token() *tokenClass_ {
	return tokenClassSingleton
}

// CLASS CONSTANTS

// This public class constant represents a boolean Token type.
func (c *tokenClass_) TypeBoolean() string {
	return c.typeBoolean
}

// This public class constant represents a complex Token type.
func (c *tokenClass_) TypeComplex() string {
	return c.typeComplex
}

// This public class constant represents a context Token type.
func (c *tokenClass_) TypeContext() string {
	return c.typeContext
}

// This public class constant represents a delimiter Token type.
func (c *tokenClass_) TypeDelimiter() string {
	return c.typeDelimiter
}

// This public class constant represents an end-of-file Token type.
func (c *tokenClass_) TypeEOF() string {
	return c.typeEOF
}

// This public class constant represents an end-of-line Token type.
func (c *tokenClass_) TypeEOL() string {
	return c.typeEOL
}

// This public class constant represents an error Token type.
func (c *tokenClass_) TypeError() string {
	return c.typeError
}

// This public class constant represents a float Token type.
func (c *tokenClass_) TypeFloat() string {
	return c.typeFloat
}

// This public class constant represents an integer Token type.
func (c *tokenClass_) TypeInteger() string {
	return c.typeInteger
}

// This public class constant represents a nil Token type.
func (c *tokenClass_) TypeNil() string {
	return c.typeNil
}

// This public class constant represents a rune Token type.
func (c *tokenClass_) TypeRune() string {
	return c.typeRune
}

// This public class constant represents a string Token type.
func (c *tokenClass_) TypeString() string {
	return c.typeString
}

// This public class constant represents an unsigned Token type.
func (c *tokenClass_) TypeUnsigned() string {
	return c.typeUnsigned
}

// CLASS CONSTRUCTORS

// This public class constructor creates a new Token from the specified
// lexical context.
func (c *tokenClass_) FromContext(
	line int,
	position int,
	tokenType string,
	tokenValue string,
) TokenLike {
	var token = &token_{
		line:       line,
		position:   position,
		tokenType:  tokenType,
		tokenValue: tokenValue,
	}
	return token
}

// CLASS TYPE

// Encapsulated Type

// This private class type encapsulates a Go structure containing private
// attributes that can only be accessed and manipulated using methods that
// implement the token-like abstract type.  A Token captures lexical context.
type token_ struct {
	line       int // The line number of the Token in the lexical context.
	position   int // The position in the line of the first rune of the Token.
	tokenType  string
	tokenValue string
}

// Contextual Interface

// This public class method returns the line number of this Token.
func (v *token_) GetLine() int {
	return v.line
}

// This public class method returns the position of this Token in the line.
func (v *token_) GetPosition() int {
	return v.position
}

// This public class method returns the type of this Token.
func (v *token_) GetType() string {
	return v.tokenType
}

// This public class method returns the value of this Token.
func (v *token_) GetValue() string {
	return v.tokenValue
}

// Private Interface

// This public class method returns the canonical string version of this Token.
func (v *token_) String() string {
	var s string
	switch {
	case v.tokenType == Token().typeEOF:
		s = "<EOF>"
	case v.tokenType == Token().typeEOL:
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
