/*******************************************************************************
 *   Copyright (c) 2009-2024 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package collections_test

import (
	col "github.com/craterdog/go-collection-framework/v3"
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

func TestAssociations(t *tes.T) {
	var association = col.AssociationClass[string, int]().FromPair("foo", 1)
	ass.Equal(t, "foo", association.GetKey())
	ass.Equal(t, 1, association.GetValue())
	association.SetValue(2)
	ass.Equal(t, 2, association.GetValue())
}
