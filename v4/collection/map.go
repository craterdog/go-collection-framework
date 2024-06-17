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
	age "github.com/craterdog/go-collection-framework/v4/agent"
	syn "sync"
)

// CLASS ACCESS

// Reference

var mapClass = map[string]any{}
var mapMutex syn.Mutex

// Function

func Map[K comparable, V any](notation NotationLike) MapClassLike[K, V] {
	// Generate the name of the bound class type.
	var class *mapClass_[K, V]
	var name = fmt.Sprintf("%T", class)

	// Check for existing bound class type.
	mapMutex.Lock()
	var value = mapClass[name]
	switch actual := value.(type) {
	case *mapClass_[K, V]:
		// This bound class type already exists.
		class = actual
	default:
		// Add a new bound class type.
		class = &mapClass_[K, V]{
			// Initialize the class constants.
			notation_: notation,
		}
		mapClass[name] = class
	}
	mapMutex.Unlock()

	// Return a reference to the bound class type.
	return class
}

// CLASS METHODS

// Target

type mapClass_[K comparable, V any] struct {
	// Define the class constants.
	notation_ NotationLike
}

// Constants

func (c *mapClass_[K, V]) Notation() NotationLike {
	return c.notation_
}

// Constructors

func (c *mapClass_[K, V]) Make() MapLike[K, V] {
	return map_[K, V](map[K]V{})
}

func (c *mapClass_[K, V]) MakeFromArray(associations []AssociationLike[K, V]) MapLike[K, V] {
	var size = len(associations)
	var duplicate = make(map[K]V, size)
	for _, association := range associations {
		var key = association.GetKey()
		var value = association.GetValue()
		duplicate[key] = value
	}
	return map_[K, V](duplicate)
}

func (c *mapClass_[K, V]) MakeFromMap(associations map[K]V) MapLike[K, V] {
	var size = len(associations)
	var duplicate = make(map[K]V, size)
	for key, value := range associations {
		duplicate[key] = value
	}
	return map_[K, V](duplicate)
}

func (c *mapClass_[K, V]) MakeFromSequence(
	associations Sequential[AssociationLike[K, V]],
) MapLike[K, V] {
	var size = associations.GetSize()
	var iterator = associations.GetIterator()
	var duplicate = make(map[K]V, size)
	for index := 0; index < size; index++ {
		var association = iterator.GetNext()
		var key = association.GetKey()
		var value = association.GetValue()
		duplicate[key] = value
	}
	return map_[K, V](duplicate)
}

// INSTANCE METHODS

// Target

type map_[K comparable, V any] map[K]V

// Attributes

func (v map_[K, V]) GetClass() MapClassLike[K, V] {
	return Map[K, V](nil)
}

// Associative

func (v map_[K, V]) GetValue(key K) V {
	var value = v[key]
	return value
}

func (v map_[K, V]) SetValue(key K, value V) {
	v[key] = value
}

func (v map_[K, V]) GetKeys() Sequential[K] {
	var size = len(v)
	var array = make([]K, size)
	var keys = array_[K](array)
	var index = 1
	var iterator = v.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		var key = association.GetKey()
		keys.SetValue(index, key)
		index++
	}
	return keys
}

func (v map_[K, V]) GetValues(keys Sequential[K]) Sequential[V] {
	var size = keys.GetSize()
	var array = make([]V, size)
	var values = array_[V](array)
	var index = 1
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		var key = iterator.GetNext()
		var value = v.GetValue(key)
		values.SetValue(index, value)
		index++
	}
	return values
}

func (v map_[K, V]) RemoveValue(key K) V {
	var value, exists = v[key]
	if exists {
		delete(v, key)
	}
	return value
}

func (v map_[K, V]) RemoveValues(keys Sequential[K]) Sequential[V] {
	var size = keys.GetSize()
	var array = make([]V, size)
	var values = array_[V](array)
	var index = 1
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		var key = iterator.GetNext()
		var value = v.RemoveValue(key)
		values.SetValue(index, value)
		index++
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

// Sequential

func (v map_[K, V]) IsEmpty() bool {
	return len(v) == 0
}

func (v map_[K, V]) GetSize() int {
	return len(v)
}

func (v map_[K, V]) AsArray() []AssociationLike[K, V] {
	var size = len(v)
	var array = make([]AssociationLike[K, V], size)
	var index = 0
	for key, value := range v {
		var association = Association[K, V](v.GetClass().Notation()).MakeWithAttributes(key, value)
		array[index] = association
		index++
	}
	return array
}

func (v map_[K, V]) GetIterator() age.IteratorLike[AssociationLike[K, V]] {
	var array = v.AsArray() // This copies the internal array.
	var associations = array_[AssociationLike[K, V]](array)
	var iterator = associations.GetIterator()
	return iterator
}

// Stringer

func (v map_[K, V]) String() string {
	return v.GetClass().Notation().FormatValue(v)
}
