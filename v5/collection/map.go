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
	age "github.com/craterdog/go-collection-framework/v5/agent"
	uti "github.com/craterdog/go-missing-utilities/v2"
	syn "sync"
)

// CLASS INTERFACE

// Access Function

func MapClass[K comparable, V any]() MapClassLike[K, V] {
	return mapClass[K, V]()
}

// Constructor Methods

func (c *mapClass_[K, V]) Map(
	associations map[K]V,
) MapLike[K, V] {
	var instance = map_[K, V](map[K]V{})
	for key, value := range associations {
		instance[key] = value
	}
	return instance
}

func (c *mapClass_[K, V]) MapFromArray(
	associations []AssociationLike[K, V],
) MapLike[K, V] {
	var instance = map_[K, V](map[K]V{})
	for _, association := range associations {
		var key = association.GetKey()
		var value = association.GetValue()
		instance[key] = value
	}
	return instance
}

func (c *mapClass_[K, V]) MapFromSequence(
	associations Sequential[AssociationLike[K, V]],
) MapLike[K, V] {
	var instance = map_[K, V](map[K]V{})
	var iterator = associations.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		var key = association.GetKey()
		var value = association.GetValue()
		instance[key] = value
	}
	return instance
}

// Constant Methods

// Function Methods

// INSTANCE INTERFACE

// Principal Methods

func (v map_[K, V]) GetClass() MapClassLike[K, V] {
	return mapClass[K, V]()
}

// Attribute Methods

// Associative[K, V] Methods

func (v map_[K, V]) AsMap() map[K]V {
	var map_ = uti.CopyMap(v)
	return map_
}

func (v map_[K, V]) GetValue(
	key K,
) V {
	var value = v[key]
	return value
}

func (v map_[K, V]) SetValue(
	key K,
	value V,
) {
	v[key] = value
}

func (v map_[K, V]) GetKeys() Sequential[K] {
	var size = v.GetSize()
	var arrayClass = ArrayClass[K]()
	var keys = arrayClass.ArrayWithSize(size)
	var index Index
	for key := range v {
		index++
		keys.SetValue(index, key)
	}
	return keys
}

func (v map_[K, V]) GetValues(
	keys Sequential[K],
) Sequential[V] {
	var size = keys.GetSize()
	var arrayClass = ArrayClass[V]()
	var values = arrayClass.ArrayWithSize(size)
	var index Index
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		index++
		var key = iterator.GetNext()
		var value = v.GetValue(key)
		values.SetValue(index, value)
	}
	return values
}

func (v map_[K, V]) RemoveValue(
	key K,
) V {
	var value, exists = v[key]
	if exists {
		delete(v, key)
	}
	return value
}

func (v map_[K, V]) RemoveValues(
	keys Sequential[K],
) Sequential[V] {
	var size = keys.GetSize()
	var arrayClass = ArrayClass[V]()
	var values = arrayClass.ArrayWithSize(size)
	var index Index
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		index++
		var key = iterator.GetNext()
		var value = v.RemoveValue(key)
		values.SetValue(index, value)
	}
	return values
}

func (v map_[K, V]) RemoveAll() {
	var keys = v.GetKeys()
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		var key = iterator.GetNext()
		delete(v, key)
	}
}

// Sequential[AssociationLike[K, V]] Methods

func (v map_[K, V]) IsEmpty() bool {
	return len(v) == 0
}

func (v map_[K, V]) GetSize() age.Size {
	var size = age.Size(len(v))
	return size
}

func (v map_[K, V]) AsArray() []AssociationLike[K, V] {
	var size = len(v)
	var array = make([]AssociationLike[K, V], size)
	var index = 0
	var associationClass = AssociationClass[K, V]()
	for key, value := range v {
		var association = associationClass.Association(key, value)
		array[index] = association
		index++
	}
	return array
}

func (v map_[K, V]) GetIterator() age.IteratorLike[AssociationLike[K, V]] {
	var array = v.AsArray()
	var iteratorClass = age.IteratorClass[AssociationLike[K, V]]()
	var iterator = iteratorClass.Iterator(array)
	return iterator
}

// Stringer Methods

func (v map_[K, V]) String() string {
	return uti.Format(v)
}

// PROTECTED INTERFACE

// Private Methods

// Instance Structure

type map_[K comparable, V any] map[K]V

// Class Structure

type mapClass_[K comparable, V any] struct {
	// Declare the class constants.
}

// Class Reference

var mapMap_ = map[string]any{}
var mapMutex_ syn.Mutex

func mapClass[K comparable, V any]() *mapClass_[K, V] {
	// Generate the name of the bound class type.
	var class *mapClass_[K, V]
	var name = fmt.Sprintf("%T", class)

	// Check for an existing bound class type.
	mapMutex_.Lock()
	var value = mapMap_[name]
	switch actual := value.(type) {
	case *mapClass_[K, V]:
		// This bound class type already exists.
		class = actual
	default:
		// Add a new bound class type.
		class = &mapClass_[K, V]{
			// Initialize the class constants.
		}
		mapMap_[name] = class
	}
	mapMutex_.Unlock()

	// Return a reference to the bound class type.
	return class
}
