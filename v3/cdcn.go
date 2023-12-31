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

// Private Class Namespace Reference

var cdcnClass = &cdcnClass_{
	defaultDepth: 8,
}

// Public Class Namespace Access

func CDCNClass() NotationClassLike {
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
		formatter: FormatterClass().WithDepth(depth),
		parser:    ParserClass().CDCN(),
	}
	return cdcn
}

// CLASS INSTANCES

// Private Class Type Definition

type cdcn_ struct {
	formatter *formatter_
	parser    *parser_
}

// Public Interface

func (v *cdcn_) FormatCollection(collection Collection) string {
	return v.formatter.FormatCollection(collection)
}

func (v *cdcn_) ParseCollection(collection string) Collection {
	return v.parser.ParseCollection(collection)
}
