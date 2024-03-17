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

import (
	col "github.com/craterdog/go-collection-framework/v3"
)

// CLASS ACCESS

// Reference

var notationClass = &notationClass_{
	// This class defines no constants.
}

// Function

func Notation() col.NotationClassLike {
	return notationClass
}

// CLASS METHODS

// Target

type notationClass_ struct {
	// This class defines no constants.
}

// Constructors

func (c *notationClass_) Make() col.NotationLike {
	return &notation_{
		formatter_: Formatter().Make(),
		parser_:    Parser().Make(),
	}
}

// INSTANCE METHODS

// Target

type notation_ struct {
	formatter_ col.FormatterLike
	parser_    ParserLike
}

// Canonical

func (v *notation_) FormatCollection(collection col.Collection) string {
	return v.formatter_.FormatCollection(collection)
}

func (v *notation_) ParseSource(source string) col.Collection {
	return v.parser_.ParseSource(source)
}
