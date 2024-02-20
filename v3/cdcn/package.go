/*******************************************************************************
 *   Copyright (c) 2009-2024 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

/*
Package cdcn defines a set of classes that provide an implementation of the
notation-like abstract class for parsing and formatting source files containing
Crater Dog Collection Notation™ (CDCN).  The complete language grammar for CDCN
is located here:

  - https://github.com/craterdog/go-collection-framework/blob/main/v3/cdcn/grammar.cdcn

This package follows the Crater Dog Technologies™ (craterdog) Go Coding
Conventions located here:

  - https://github.com/craterdog/go-coding-conventions/wiki

Additional implementations of the classes provided by this package can be
developed and used seamlessly since the interface definitions only depend on
other interfaces and primitive types; and the class implementations only depend
on interfaces, not on each other.
*/
package cdcn

import (
	col "github.com/craterdog/go-collection-framework/v3"
)

// TYPES

// Specializations

/*
TokenType is a specialized type representing any token type recognized by a
scanner.
*/
type TokenType uint8

const (
	ErrorToken TokenType = iota
	BooleanToken
	ComplexToken
	ContextToken
	DelimiterToken
	EOFToken
	EOLToken
	FloatToken
	HexadecimalToken
	IntegerToken
	NilToken
	RuneToken
	SpaceToken
	StringToken
)

// INTERFACES

// Classes

/*
ParserClassLike defines the set of class constants, constructors and functions
that must be supported by all parser-class-like classes.
*/
type ParserClassLike interface {
	// Constructors
	Make() ParserLike
}

/*
ScannerClassLike defines the set of class constants, constructors and functions
that must be supported by all scanner-class-like classes.  The following
functions are supported:

MatchToken() a list of strings representing any matches found in the specified
text of the specified token type using the regular expression defined for that
token type.  If the regular expression contains submatch patterns the matching
substrings are returned as additional values in the list.
*/
type ScannerClassLike interface {
	// Constructors
	Make(document string, tokens col.QueueLike[TokenLike]) ScannerLike

	// Functions
	MatchToken(tokenType TokenType, text string) col.ListLike[string]
}

/*
TokenClassLike defines the set of class constants, constructors and functions
that must be supported by all token-class-like classes.  The following functions
are supported:

AsString() returns a string representing the specified token type.
*/
type TokenClassLike interface {
	// Constructors
	Make(
		line, position int,
		tokenType TokenType,
		tokenValue string,
	) TokenLike

	// Functions
	AsString(tokenType TokenType) string
}

// Instances

/*
ParserLike defines the set of abstractions and methods that must be supported by
all parser-like instances.
*/
type ParserLike interface {
	// Methods
	ParseSource(source string) col.Collection
}

/*
ScannerLike defines the set of abstractions and methods that must be supported
by all scanner-like instances.
*/
type ScannerLike interface {
}

/*
TokenLike defines the set of abstractions and methods that must be supported by
all token-like instances.
*/
type TokenLike interface {
	// Methods
	GetLine() int
	GetPosition() int
	GetType() TokenType
	GetValue() string
}
