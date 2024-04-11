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
Package "cdcn" provides a set of classes that provide an implementation of the
notation-like abstract class for parsing and formatting source files containing
Crater Dog Collection Notation™ (CDCN).  The complete language grammar for CDCN
is located here:
  - https://github.com/craterdog/go-collection-framework/blob/main/v3/cdcn/Grammar.cdsn

This package follows the Crater Dog Technologies™ Go Coding Conventions located
here:
  - https://github.com/craterdog/go-model-framework/wiki

Additional concrete implementations of the classes defined by this package can
be developed and used seamlessly since the interface definitions only depend on
other interfaces and primitive types—and the class implementations only depend
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
AssociationClassLike is a class interface that defines the complete set of
class constants, constructors and functions that must be supported by each
concrete association-like class.
*/
type AssociationClassLike interface {
	// Constructors
	MakeWithAttributes(
		key KeyLike,
		value ValueLike,
	) AssociationLike
}

/*
AssociationsClassLike is a class interface that defines the complete set of
class constants, constructors and functions that must be supported by each
concrete associations-like class.
*/
type AssociationsClassLike interface {
	// Constructors
	MakeWithAssociations(associations col.ListLike[AssociationLike]) AssociationsLike
}

/*
CollectionClassLike is a class interface that defines the complete set of
class constants, constructors and functions that must be supported by each
concrete collection-like class.
*/
type CollectionClassLike interface {
	// Constructors
	MakeWithAttributes(
		associations AssociationsLike,
		values ValuesLike,
		context string,
	) CollectionLike
}

/*
FormatterClassLike is a class interface that defines the complete set of
class constants, constructors and functions that must be supported by each
concrete formatter-like class.
*/
type FormatterClassLike interface {
	// Constructors
	Make() FormatterLike
}

/*
KeyClassLike is a class interface that defines the complete set of
class constants, constructors and functions that must be supported by each
concrete key-like class.
*/
type KeyClassLike interface {
	// Constructors
	MakeWithPrimitive(primitive PrimitiveLike) KeyLike
}

/*
ParserClassLike is a class interface that defines the complete set of
class constants, constructors and functions that must be supported by each
concrete parser-like class.
*/
type ParserClassLike interface {
	// Constructors
	Make() ParserLike
}

/*
PrimitiveClassLike is a class interface that defines the complete set of
class constants, constructors and functions that must be supported by each
concrete primitive-like class.
*/
type PrimitiveClassLike interface {
	// Constructors
	MakeWithBoolean(boolean string) PrimitiveLike
	MakeWithComplex(complex_ string) PrimitiveLike
	MakeWithFloat(float string) PrimitiveLike
	MakeWithHexadecimal(hexadecimal string) PrimitiveLike
	MakeWithInteger(integer string) PrimitiveLike
	MakeWithNil(nil_ string) PrimitiveLike
	MakeWithRune(rune_ string) PrimitiveLike
	MakeWithString(string_ string) PrimitiveLike
}

/*
ScannerClassLike is a class interface that defines the complete set of
class constants, constructors and functions that must be supported by each
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
class constants, constructors and functions that must be supported by each
concrete token-like class.
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

/*
ValidatorClassLike is a class interface that defines the complete set of
class constants, constructors and functions that must be supported by each
concrete validator-like class.
*/
type ValidatorClassLike interface {
	// Constructors
	Make() ValidatorLike
}

/*
ValueClassLike is a class interface that defines the complete set of
class constants, constructors and functions that must be supported by each
concrete value-like class.
*/
type ValueClassLike interface {
	// Constructors
	MakeWithCollection(collection CollectionLike) ValueLike
	MakeWithPrimitive(primitive PrimitiveLike) ValueLike
}

/*
ValuesClassLike is a class interface that defines the complete set of
class constants, constructors and functions that must be supported by each
concrete values-like class.
*/
type ValuesClassLike interface {
	// Constructors
	MakeWithValues(values col.ListLike[ValueLike]) ValuesLike
}

// Instances

/*
AssociationLike is an instance interface that defines the complete set of
instance attributes, abstractions and methods that must be supported by each
instance of a concrete association-like class.
*/
type AssociationLike interface {
	// Attributes
	GetKey() KeyLike
	GetValue() ValueLike
}

/*
AssociationsLike is an instance interface that defines the complete set of
instance attributes, abstractions and methods that must be supported by each
instance of a concrete associations-like class.
*/
type AssociationsLike interface {
	// Attributes
	GetAssociations() col.ListLike[AssociationLike]
}

/*
CollectionLike is an instance interface that defines the complete set of
instance attributes, abstractions and methods that must be supported by each
instance of a concrete collection-like class.
*/
type CollectionLike interface {
	// Attributes
	GetAssociations() AssociationsLike
	GetValues() ValuesLike
	GetContext() string
}

/*
FormatterLike is an instance interface that defines the complete set of
instance attributes, abstractions and methods that must be supported by each
instance of a concrete formatter-like class.
*/
type FormatterLike interface {
	// Methods
	FormatCollection(collection CollectionLike) string
}

/*
KeyLike is an instance interface that defines the complete set of
instance attributes, abstractions and methods that must be supported by each
instance of a concrete key-like class.
*/
type KeyLike interface {
	// Attributes
	GetPrimitive() PrimitiveLike
}

/*
ParserLike is an instance interface that defines the complete set of
instance attributes, abstractions and methods that must be supported by each
instance of a concrete parser-like class.
*/
type ParserLike interface {
	// Methods
	ParseSource(source string) col.Collection
}

/*
PrimitiveLike is an instance interface that defines the complete set of
instance attributes, abstractions and methods that must be supported by each
instance of a concrete primitive-like class.
*/
type PrimitiveLike interface {
	// Attributes
	GetBoolean() string
	GetComplex() string
	GetFloat() string
	GetHexadecimal() string
	GetInteger() string
	GetNil() string
	GetRune() string
	GetString() string
}

/*
ScannerLike is an instance interface that defines the complete set of
instance attributes, abstractions and methods that must be supported by each
instance of a concrete scanner-like class.
*/
type ScannerLike interface {
}

/*
TokenLike is an instance interface that defines the complete set of
instance attributes, abstractions and methods that must be supported by each
instance of a concrete token-like class.
*/
type TokenLike interface {
	// Attributes
	GetLine() int
	GetPosition() int
	GetType() TokenType
	GetValue() string
}

/*
ValidatorLike is an instance interface that defines the complete set of
instance attributes, abstractions and methods that must be supported by each
instance of a concrete validator-like class.
*/
type ValidatorLike interface {
	// Methods
	ValidateCollection(collection CollectionLike)
}

/*
ValueLike is an instance interface that defines the complete set of
instance attributes, abstractions and methods that must be supported by each
instance of a concrete value-like class.
*/
type ValueLike interface {
	// Attributes
	GetCollection() CollectionLike
	GetPrimitive() PrimitiveLike
}

/*
ValuesLike is an instance interface that defines the complete set of
instance attributes, abstractions and methods that must be supported by each
instance of a concrete values-like class.
*/
type ValuesLike interface {
	// Attributes
	GetValues() col.ListLike[ValueLike]
}
