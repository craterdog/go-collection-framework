/*
................................................................................
.    Copyright (c) 2009-2024 Crater Dog Technologies.  All Rights Reserved.    .
................................................................................
.  DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               .
.                                                                              .
.  This code is free software; you can redistribute it and/or modify it under  .
.  the terms of The MIT License (MIT), as published by the Open Source         .
.  Initiative. (See https://opensource.org/license/MIT)                        .
................................................................................
*/

/*
Package "cdcn" provides a set of classes that provide an implementation of the
notation-like abstract class for parsing and formatting source files containing
Crater Dog Collection Notation™ (CDCN).  The complete language syntax for CDCN
is located here:
  - https://github.com/craterdog/go-collection-framework/blob/main/v4/cdcn/Syntax.cdsn

This package follows the Crater Dog Technologies™ Go Coding Conventions located
here:
  - https://github.com/craterdog/go-model-framework/wiki

Additional concrete implementations of the classes defined by this package can
be developed and used seamlessly since the interface definitions only depend on
other interfaces and intrinsic types—and the class implementations only depend
on interfaces, not on each other.
*/
package cdcn

import (
	col "github.com/craterdog/go-collection-framework/v4/collection"
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
	TypeToken
)

// Classes

/*
FormatterClassLike is a class interface that defines the complete set of
class constructors, constants and functions that must be supported by each
concrete formatter-like class.
*/
type FormatterClassLike interface {
	// Constructors
	Make() FormatterLike
	MakeWithMaximum(maximum int) FormatterLike

	// Constants
	DefaultMaximum() int
}

/*
ParserClassLike is a class interface that defines the complete set of
class constructors, constants and functions that must be supported by each
concrete parser-like class.
*/
type ParserClassLike interface {
	// Constructors
	Make() ParserLike
}

/*
ScannerClassLike is a class interface that defines the complete set of
class constructors, constants and functions that must be supported by each
concrete scanner-like class.  The following functions are supported:

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
TokenClassLike is a class interface that defines the complete set of
class constructors, constants and functions that must be supported by each
concrete token-like class.
*/
type TokenClassLike interface {
	// Constructors
	Make(
		line int,
		position int,
		type_ TokenType,
		value string,
	) TokenLike
}

// Instances

/*
FormatterLike is an instance interface that defines the complete set of
instance attributes, abstractions and methods that must be supported by each
instance of a concrete formatter-like class.
*/
type FormatterLike interface {
	// Attributes
	GetClass() FormatterClassLike
	GetDepth() int
	GetMaximum() int

	// Methods
	FormatValue(value any) (source string)
}

/*
ParserLike is an instance interface that defines the complete set of
instance attributes, abstractions and methods that must be supported by each
instance of a concrete parser-like class.
*/
type ParserLike interface {
	// Attributes
	GetClass() ParserClassLike

	// Methods
	ParseSource(source string) (value any)
}

/*
ScannerLike is an instance interface that defines the complete set of
instance attributes, abstractions and methods that must be supported by each
instance of a concrete scanner-like class.
*/
type ScannerLike interface {
	// Attributes
	GetClass() ScannerClassLike
}

/*
TokenLike is an instance interface that defines the complete set of
instance attributes, abstractions and methods that must be supported by each
instance of a concrete token-like class.
*/
type TokenLike interface {
	// Attributes
	GetClass() TokenClassLike
	GetLine() int
	GetPosition() int
	GetType() TokenType
	GetValue() string
}
