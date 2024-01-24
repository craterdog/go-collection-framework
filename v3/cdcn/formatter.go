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

// Public Class Namespace Access

/*
FormatterClass defines an implementation of a formatter-like class that uses
Crater Dog Collection Notation™ (CDCN) for formatting collections.  Since the
go-collection-framework uses the same formatting notation this class simply
delegates to the collection framework's formatter class.
*/
func FormatterClass() col.FormatterClassLike {
	return col.FormatterClass()
}
