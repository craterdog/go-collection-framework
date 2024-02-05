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

	https://github.com/craterdog/go-collection-framework/blob/main/v3/cdcn/grammar.cdcn

This package follows the Crater Dog Technologies™ (craterdog) Go Coding
Conventions located here:

	https://github.com/craterdog/go-coding-conventions/wiki

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

// Enumerated

/*
This enumerated type represents all possible values for a CDCN token type.
*/
const (
	TypeError col.TokenType = iota
	TypeBoolean
	TypeComplex
	TypeContext
	TypeDelimiter
	TypeEOF
	TypeEOL
	TypeFloat
	TypeHexadecimal
	TypeInteger
	TypeNil
	TypeRune
	TypeSpace
	TypeString
)
