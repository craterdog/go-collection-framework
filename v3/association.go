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

import (
	fmt "fmt"
)

// CLASS NAMESPACE

// Private Class Namespace Type

type associationClass_[K Key, V Value] struct {
	// This class has no class constants.
}

// Private Namespace Reference(s)

var associationClass = map[string]any{}

// Public Namespace Access

func Association[K Key, V Value]() AssociationClassLike[K, V] {
	var class *associationClass_[K, V]
	var key = fmt.Sprintf("%T", class) // The name of the bound class type.
	var value = associationClass[key]
	switch actual := value.(type) {
	case *associationClass_[K, V]:
		// This bound class type already exists.
		class = actual
	default:
		// Create a new bound class type.
		class = &associationClass_[K, V]{
			// This class has no class constants.
		}
		associationClass[key] = class
	}
	return class
}

// Public Class Constructors

func (c *associationClass_[K, V]) FromPair(key K, value V) AssociationLike[K, V] {
	var association = &association_[K, V]{key, value}
	return association
}

// CLASS TYPE

// Private Class Type Definition

type association_[K Key, V Value] struct {
	key   K
	value V
}

// Binding Interface

func (v *association_[K, V]) GetKey() K {
	return v.key
}

func (v *association_[K, V]) GetValue() V {
	return v.value
}

func (v *association_[K, V]) SetValue(value V) {
	v.value = value
}
