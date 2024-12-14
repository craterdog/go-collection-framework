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
	uti "github.com/craterdog/go-missing-utilities/v2"
	syn "sync"
)

// CLASS INTERFACE

// Access Function

func AssociationClass[K comparable, V any]() AssociationClassLike[K, V] {
	return associationClassReference[K, V]()
}

// Constructor Methods

func (c *associationClass_[K, V]) Association(
	key K,
	value V,
) AssociationLike[K, V] {
	if uti.IsUndefined(key) {
		panic("The \"key\" attribute is required by this class.")
	}
	if uti.IsUndefined(value) {
		panic("The \"value\" attribute is required by this class.")
	}
	var instance = &association_[K, V]{
		// Initialize the instance attributes.
		key_:   key,
		value_: value,
	}
	return instance
}

// Constant Methods

// Function Methods

// INSTANCE INTERFACE

// Principal Methods

func (v *association_[K, V]) GetClass() AssociationClassLike[K, V] {
	return associationClassReference[K, V]()
}

// Attribute Methods

func (v *association_[K, V]) GetKey() K {
	return v.key_
}

func (v *association_[K, V]) GetValue() V {
	return v.value_
}

func (v *association_[K, V]) SetValue(
	value V,
) {
	if uti.IsUndefined(value) {
		panic("The \"value\" attribute is required by this class.")
	}
	v.value_ = value
}

// Stringer Methods

func (v *association_[K, V]) String() string {
	var result = uti.Format(v.GetKey())
	result += ": "
	result += uti.Format(v.GetValue())
	return result
}

// PROTECTED INTERFACE

// Private Methods

// Instance Structure

type association_[K comparable, V any] struct {
	// Declare the instance attributes.
	key_   K
	value_ V
}

// Class Structure

type associationClass_[K comparable, V any] struct {
	// Declare the class constants.
}

// Class Reference

var associationMap_ = map[string]any{}
var associationMutex_ syn.Mutex

func associationClassReference[K comparable, V any]() *associationClass_[K, V] {
	// Generate the name of the bound class type.
	var class *associationClass_[K, V]
	var name = fmt.Sprintf("%T", class)

	// Check for an existing bound class type.
	associationMutex_.Lock()
	var value = associationMap_[name]
	switch actual := value.(type) {
	case *associationClass_[K, V]:
		// This bound class type already exists.
		class = actual
	default:
		// Add a new bound class type.
		class = &associationClass_[K, V]{
			// Initialize the class constants.
		}
		associationMap_[name] = class
	}
	associationMutex_.Unlock()

	// Return a reference to the bound class type.
	return class
}
