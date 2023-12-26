/*******************************************************************************
 *   Copyright (c) 2009-2024 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package collections

// CLASS NAMESPACE

// Private Class Namespace Type

type cdcnClass_ struct {
	defaultDepth int
}

// Private Namespace Reference(s)

var cdcnClass = &cdcnClass_{
	defaultDepth: 8,
}

// Public Namespace Access

func CDCN() NotationClassLike {
	return cdcnClass
}

// Public Class Constants

func (c *cdcnClass_) GetDefaultDepth() int {
	return c.defaultDepth
}

// Public Class Constructors

func (c *cdcnClass_) Default() NotationLike {
	var cdcn = c.WithDepth(c.defaultDepth)
	return cdcn
}

func (c *cdcnClass_) WithDepth(depth int) NotationLike {
	if depth < 0 || depth > c.defaultDepth {
		depth = c.defaultDepth
	}
	var cdcn = &cdcn_{
		formatter: Formatter().WithDepth(depth),
		parser:    Parser().CDCN(),
	}
	return cdcn
}

// CLASS TYPE

// Private Class Type Definition

type cdcn_ struct {
	formatter *formatter_
	parser    *parser_
}

// Standardized Interface

func (v *cdcn_) FormatCollection(collection Collection) string {
	return v.formatter.FormatCollection(collection)
}

// Stringent Interface

func (v *cdcn_) ParseCollection(source string) Collection {
	return v.parser.ParseCollection(source)
}
