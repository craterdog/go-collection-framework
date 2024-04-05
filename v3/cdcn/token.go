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

package cdcn

import ()

// CLASS ACCESS

// Reference

var tokenClass = &tokenClass_{
	// This class has no private constants to initialize.
}

// Function

func Token() TokenClassLike {
	return tokenClass
}

// CLASS METHODS

// Target

type tokenClass_ struct {
	// This class has no private constants.
}

// Constants

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

// INSTANCE METHODS

// Target

type token_ struct {
	line_     int
	position_ int
	type_     TokenType
	value_    string
}

// Attributes

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

// Public

// Private
