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

// This private type defines the namespace structure associated with the
// constants, constructors and functions for the Association class namespace.
type associationClass_[K Key, V Value] struct {
	// This class has no class constants.
}

// This private constant defines a map to hold all the singleton references to
// the type specific Association class namespaces.
var associationClassSingletons = map[string]any{}

// This public function returns the singleton reference to a type specific
// Association class namespace.  It also initializes any class constants as
// needed.
func Association[K Key, V Value]() *associationClass_[K, V] {
	var class *associationClass_[K, V]
	var key = fmt.Sprintf("%T", class) // The name of the bound class type.
	var value = associationClassSingletons[key]
	switch actual := value.(type) {
	case *associationClass_[K, V]:
		// This bound class type already exists.
		class = actual
	default:
		// Create a new bound class type.
		class = &associationClass_[K, V]{
			// This class has no class constants.
		}
		associationClassSingletons[key] = class
	}
	return class
}

// CLASS CONSTRUCTORS

// This public class constructor creates a new Association from the specified
// key and value.
func (c *associationClass_[K, V]) FromPair(key K, value V) AssociationLike[K, V] {
	var association = &association_[K, V]{key, value}
	return association
}

// CLASS TYPE

// This private class structure encapsulates a Go structure containing private
// attributes that can only be accessed and manipulated using methods that
// implement the association-like abstract type.  The attributes maintain the
// information about a key-value pair. This type is parameterized as follows:
//   - K is a primitive type of key.
//   - V is any type of value.
//
// This structure is used by the catalog class to maintain its associations.
type association_[K Key, V Value] struct {
	key   K
	value V
}

// Binding Interface

// This public class method returns the key for this Association.
func (v *association_[K, V]) GetKey() K {
	return v.key
}

// This public class method returns the value for this Association.
func (v *association_[K, V]) GetValue() V {
	return v.value
}

// This public class method sets the value of this Association to a new value.
func (v *association_[K, V]) SetValue(value V) {
	v.value = value
}

// Private Interface

// This public class method returns a canonical string for this Association.
func (v *association_[K, V]) String() string {
	return Formatter().FormatValue(v)
}
