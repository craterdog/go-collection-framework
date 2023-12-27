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

// The Go language requires the key type here support the "comparable"
// interface so we must narrow it down from type Key (i.e. "any").
type mapClass_[K comparable, V Value] struct {
	// This class defines no constants.
}

// Private Namespace Reference(s)

var mapClass = map[string]any{}

// Public Namespace Access

func MapClass[K comparable, V Value]() MapClassLike[K, V] {
	var class *mapClass_[K, V]
	var key = fmt.Sprintf("%T", class) // The name of the bound class type.
	var value = mapClass[key]
	switch actual := value.(type) {
	case *mapClass_[K, V]:
		// This bound class type already exists.
		class = actual
	default:
		// Create a new bound class type.
		class = &mapClass_[K, V]{
			// This class defines no constants.
		}
		mapClass[key] = class
	}
	return class
}

// Public Class Constructors

func (c *mapClass_[K, V]) Empty() MapLike[K, V] {
	return map_[K, V](map[K]V{})
}

func (c *mapClass_[K, V]) FromArray(associations []Binding[K, V]) MapLike[K, V] {
	var size = len(associations)
	var duplicate = make(map[K]V, size)
	for _, association := range associations {
		var key = association.GetKey()
		var value = association.GetValue()
		duplicate[key] = value
	}
	return map_[K, V](duplicate)
}

func (c *mapClass_[K, V]) FromMap(associations map[K]V) MapLike[K, V] {
	var size = len(associations)
	var duplicate = make(map[K]V, size)
	for key, value := range associations {
		duplicate[key] = value
	}
	return map_[K, V](duplicate)
}

func (c *mapClass_[K, V]) FromSequence(
	associations Sequential[Binding[K, V]],
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

func (c *mapClass_[K, V]) FromString(associations string) MapLike[K, V] {
	// First we parse it as a collection of any type value.
	var cdcn = CDCNClass().Default()
	var collection = cdcn.ParseCollection(associations).(Sequential[Binding[Key, Value]])

	// Then we convert it to a Map of type Binding[K, V].
	var map_ = c.Empty()
	var iterator = collection.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		var key = association.GetKey().(K)
		var value = association.GetValue().(V)
		map_.SetValue(key, value)
	}
	return map_
}

// CLASS TYPE

// Private Class Type Definition

type map_[K comparable, V Value] map[K]V

// Associative Interface

func (v map_[K, V]) GetKeys() Sequential[K] {
	var keys = ListClass[K]().Empty()
	var iterator = v.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		keys.AppendValue(association.GetKey())
	}
	return keys
}

func (v map_[K, V]) GetValue(key K) V {
	var value = v[key]
	return value
}

func (v map_[K, V]) GetValues(keys Sequential[K]) Sequential[V] {
	var values = ListClass[V]().Empty()
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		var key = iterator.GetNext()
		values.AppendValue(v.GetValue(key))
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

func (v map_[K, V]) RemoveValue(key K) V {
	var value, exists = v[key]
	if exists {
		delete(v, key)
	}
	return value
}

func (v map_[K, V]) RemoveValues(keys Sequential[K]) Sequential[V] {
	var values = ListClass[V]().Empty()
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		var key = iterator.GetNext()
		values.AppendValue(v.RemoveValue(key))
	}
	return values
}

func (v map_[K, V]) SetValue(key K, value V) {
	v[key] = value
}

// Sequential Interface

func (v map_[K, V]) AsArray() []Binding[K, V] {
	var size = len(v)
	var result = make([]Binding[K, V], size)
	var index = 0
	for key, value := range v {
		var association = AssociationClass[K, V]().FromPair(key, value)
		result[index] = association
		index++
	}
	return result
}

func (v map_[K, V]) GetIterator() Ratcheted[Binding[K, V]] {
	var array = ArrayClass[Binding[K, V]]().FromArray(v.AsArray())
	var iterator = array.GetIterator()
	return iterator
}

func (v map_[K, V]) GetSize() int {
	return len(v)
}

func (v map_[K, V]) IsEmpty() bool {
	return len(v) == 0
}

// Private Interface

// This public class method is used by Go to generate a string from a Map.
func (v map_[K, V]) String() string {
	var cdcn = CDCNClass().Default()
	return cdcn.FormatCollection(v)
}
