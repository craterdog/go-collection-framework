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
	syn "sync"
)

// CLASS ACCESS

// Reference

var mapClass = map[string]any{}
var mapMutex syn.Mutex

// Function

func Map[K comparable, V Value]() MapClassLike[K, V] {
	// Generate the name of the bound class type.
	var class MapClassLike[K, V]
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
			// This class defines no constants.
		}
		mapClass[name] = class
	}
	mapMutex.Unlock()

	// Return a reference to the bound class type.
	return class
}

// CLASS METHODS

// Target

/*
NOTE:
The Go language requires the key type here support the "comparable" interface so
we must narrow it down from type Key (i.e. "any").
*/
type mapClass_[K comparable, V Value] struct {
	// This class defines no constants.
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

func (c *mapClass_[K, V]) MakeFromSource(
	source string,
	notation NotationLike,
) MapLike[K, V] {
	// First we parse it as a collection of any type value.
	var collection = notation.ParseSource(source).(Sequential[AssociationLike[Key, Value]])

	// Then we convert it to a Map of type AssociationLike[K, V].
	var map_ = c.Make()
	var iterator = collection.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		var key = association.GetKey().(K)
		var value = association.GetValue().(V)
		map_.SetValue(key, value)
	}
	return map_
}

// INSTANCE METHODS

// Target

type map_[K comparable, V Value] map[K]V

// Associative

func (v map_[K, V]) GetKeys() Sequential[K] {
	var keys = List[K]().Make()
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
	var values = List[V]().Make()
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
	var values = List[V]().Make()
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

// Sequential

func (v map_[K, V]) AsArray() []AssociationLike[K, V] {
	var size = len(v)
	var result = make([]AssociationLike[K, V], size)
	var index = 0
	for key, value := range v {
		var association = Association[K, V]().MakeWithAttributes(key, value)
		result[index] = association
		index++
	}
	return result
}

func (v map_[K, V]) GetIterator() IteratorLike[AssociationLike[K, V]] {
	var array = Array[AssociationLike[K, V]]().MakeFromArray(v.AsArray())
	var iterator = array.GetIterator()
	return iterator
}

func (v map_[K, V]) GetSize() int {
	return len(v)
}

func (v map_[K, V]) IsEmpty() bool {
	return len(v) == 0
}

// Stringer

func (v map_[K, V]) String() string {
	var formatter = Formatter().Make()
	return formatter.FormatCollection(v)
}
