/*******************************************************************************
 *   Copyright (c) 2009-2024 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package cdcn

import (
	col "github.com/craterdog/go-collection-framework/v3"
)

// CLASS NAMESPACE

// Private Class Namespace Type

type notationClass_ struct {
	defaultDepth int
}

// Private Class Namespace Reference

var notationClass = &notationClass_{
	defaultDepth: 8,
}

// Public Class Namespace Access

func NotationClass() col.NotationClassLike {
	return notationClass
}

// Public Class Constants

func (c *notationClass_) DefaultDepth() int {
	return c.defaultDepth
}

// Public Class Constructors

func (c *notationClass_) Make() col.NotationLike {
	var notation = c.MakeWithDepth(c.defaultDepth)
	return notation
}

func (c *notationClass_) MakeWithDepth(depth int) col.NotationLike {
	if depth < 0 || depth > c.defaultDepth {
		depth = c.defaultDepth
	}
	var notation = &notation_{
		formatter: FormatterClass().MakeWithDepth(depth),
		parser:    ParserClass().Make(),
	}
	return notation
}

// CLASS INSTANCES

// Private Class Type Definition

type notation_ struct {
	formatter col.FormatterLike
	parser    col.ParserLike
}

// Public Interface

func (v *notation_) FormatCollection(collection col.Collection) string {
	return v.formatter.FormatCollection(collection)
}

func (v *notation_) ParseSource(source string) col.Collection {
	return v.parser.ParseSource(source)
}
