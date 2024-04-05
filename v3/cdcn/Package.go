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

/*
Package "cdcn" defines a set of classes that provide an implementation of the
notation-like abstract class for parsing and formatting source files containing
Crater Dog Collection Notation™ (CDCN).  The complete language grammar for CDCN
is located here:
  - https://github.com/craterdog/go-collection-framework/blob/main/v3/cdcn/Grammar.cdsn

This package follows the Crater Dog Technologies™ (craterdog) Go Coding
Conventions located here:
  - https://github.com/craterdog/go-model-framework/wiki

Additional implementations of the classes provided by this package can be
developed and used seamlessly since the interface definitions only depend on
other interfaces and primitive types; and the class implementations only depend
on interfaces, not on each other.
*/
package cdcn

import (
	col "github.com/craterdog/go-collection-framework/v3"
)

// Types

/*
TokenType is a constrained type representing any token type recognized by a
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
ScannerClassLike is a class interface that defines the set of class
constants, constructors and functions that must be supported by each
scanner-class-like concrete class.  The following functions are supported:

FormatToken() returns a formatted string containing the attributes of the token.

MatchToken() a list of strings representing any matches found in the specified
text of the specified token type using the regular expression defined for that
token type.  If the regular expression contains submatch patterns the matching
substrings are returned as additional values in the list.
*/
type ScannerClassLike interface {
	// Constructors
	Make(
		source string,
		tokens col.QueueLike[TokenLike],
	) ScannerLike

	// Functions
	FormatToken(token TokenLike) string
	MatchToken(
		type_ TokenType,
		text string,
	) col.ListLike[string]
}

/*
TokenClassLike is a class interface that defines the set of class
constants, constructors and functions that must be supported by each
token-class-like concrete class.
*/
type TokenClassLike interface {
	// Constructors
	MakeWithAttributes(
		line int,
		position int,
		type_ TokenType,
		value string,
	) TokenLike
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
TokenLike is an instance interface that defines the complete set of
abstractions and methods that must be supported by each instance of a
token-like concrete class.
*/
type TokenLike interface {
	// Attributes
	GetLine() int
	GetPosition() int
	GetType() TokenType
	GetValue() string
}
