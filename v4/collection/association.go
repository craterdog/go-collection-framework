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

package collection

import (
	fmt "fmt"
	syn "sync"
)

// CLASS ACCESS

// Reference

var associationClass = map[string]any{}
var associationMutex syn.Mutex

// Function

func Association[K comparable, V any](notation NotationLike) AssociationClassLike[K, V] {
	// Generate the name of the bound class type.
	var class *associationClass_[K, V]
	var name = fmt.Sprintf("%T", class)

	// Check for existing bound class type.
	associationMutex.Lock()
	var value = associationClass[name]
	switch actual := value.(type) {
	case *associationClass_[K, V]:
		// This bound class type already exists.
		class = actual
	default:
		// Add a new bound class type.
		class = &associationClass_[K, V]{
			// Initialize the class constants.
			notation_: notation,
		}
		associationClass[name] = class
	}
	associationMutex.Unlock()

	// Return a reference to the bound class type.
	return class
}

// CLASS METHODS

// Target

type associationClass_[K comparable, V any] struct {
	// Define the class constants.
	notation_ NotationLike
}

// Constants

func (c *associationClass_[K, V]) Notation() NotationLike {
	return c.notation_
}

// Constructors

func (c *associationClass_[K, V]) Make(
	key K,
	value V,
) AssociationLike[K, V] {
	return &association_[K, V]{
		class_: c,
		key_:   key,
		value_: value,
	}
}

// INSTANCE METHODS

// Target

type association_[K comparable, V any] struct {
	class_ AssociationClassLike[K, V]
	key_   K
	value_ V
}

// Attributes

func (v *association_[K, V]) GetClass() AssociationClassLike[K, V] {
	return v.class_
}

func (v *association_[K, V]) GetKey() K {
	return v.key_
}

func (v *association_[K, V]) GetValue() V {
	return v.value_
}

func (v *association_[K, V]) SetValue(value V) {
	v.value_ = value
}

// Stringer

func (v *association_[K, V]) String() string {
	return v.GetClass().Notation().FormatValue(v)
}
