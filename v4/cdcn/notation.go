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

package cdcn

import (
	col "github.com/craterdog/go-collection-framework/v4/collection"
)

// CLASS ACCESS

// Reference

var notationClass = &notationClass_{
	// Initialize the class constants.
}

// Function

func Notation() col.NotationClassLike {
	return notationClass
}

// CLASS METHODS

// Target

type notationClass_ struct {
	// Define the class constants.
}

// Constructors

func (c *notationClass_) Make() col.NotationLike {
	return &notation_{
		class_:     c,
		formatter_: Formatter().Make(),
		parser_:    Parser().Make(),
	}
}

// INSTANCE METHODS

// Target

type notation_ struct {
	class_     col.NotationClassLike
	formatter_ FormatterLike
	parser_    ParserLike
}

// Attributes

func (v *notation_) GetClass() col.NotationClassLike {
	return v.class_
}

// Canonical

func (v *notation_) FormatValue(value any) (source string) {
	source = v.formatter_.FormatValue(value)
	return source
}

func (v *notation_) ParseSource(source string) (value any) {
	value = v.parser_.ParseSource(source)
	return value
}
