/*******************************************************************************
 *   Copyright (c) 2009-2023 Crater Dog Technologies™.  All Rights Reserved.   *
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
// constants, constructors and functions for the map class namespace.
// NOTE: the Go language requires the key type here support the "comparable"
// interface so we must narrow it down from "any".
type mapClass_[K comparable, V Value] struct {
	// This class defines no constants.
}

// This private constant defines a map to hold all the singleton references to
// the type specific map namespaces.
var mapClassSingletons = map[string]any{}

// This public function returns the singleton reference to a type specific
// map namespace.  It also initializes any class constants as needed.
// NOTE: the Go language requires the key type here support the "comparable"
// interface so we must narrow it down from "any".
func Map[K comparable, V Value]() *mapClass_[K, V] {
	var class *mapClass_[K, V]
	var key = fmt.Sprintf("%T", class)
	var value = mapClassSingletons[key]
	switch actual := value.(type) {
	case *mapClass_[K, V]:
		class = actual
	default:
		class = &mapClass_[K, V]{
			// This class defines no constants.
		}
		mapClassSingletons[key] = class
	}
	return class
}

// CLASS CONSTRUCTORS

// This public class constructor creates a new map from the specified
// Go array of associations.
func (c *mapClass_[K, V]) FromArray(associations []Binding[K, V]) MapLike[K, V] {
	var length = len(associations)
	var duplicate = make(map[K]V, length)
	for _, association := range associations {
		var key = association.GetKey()
		var value = association.GetValue()
		duplicate[key] = value
	}
	return map_[K, V](duplicate)
}

// This public class constructor creates a new map from the specified
// Go map of associations.
func (c *mapClass_[K, V]) FromMap(associations map[K]V) MapLike[K, V] {
	var length = len(associations)
	var duplicate = make(map[K]V, length)
	for key, value := range associations {
		duplicate[key] = value
	}
	return map_[K, V](duplicate)
}

// CLASS TYPE

// Extended Type

// This private class type extends the primitive Go "map[K]V" data type and
// defines the methods that implement the map-like abstract type.  The map
// extended type manages key-value pair associations.
// This type is parameterized as follows:
//   - K is a primitive type of key.
//   - V is any type of entity.
type map_[K comparable, V Value] map[K]V

// Sequential Interface

// This public class method determines whether or not this map is empty.
func (v map_[K, V]) IsEmpty() bool {
	return len(v) == 0
}

// This public class method returns the number of values contained in this map.
func (v map_[K, V]) GetSize() int {
	return len(v)
}

// This public class method returns all the associations that are in this map.
// The associations retrieved are in the same order as they are in the map.
func (v map_[K, V]) AsArray() []Binding[K, V] {
	var length = len(v)
	var result = make([]Binding[K, V], length)
	var index = 0
	for key, value := range v {
		var Association = Association[K, V]()
		var association = Association.FromPair(key, value)
		result[index] = association
		index++
	}
	return result
}

// This public class method generates for this map an iterator that can be used
// to traverse its associations.
func (v map_[K, V]) GetIterator() Ratcheted[Binding[K, V]] {
	var Array = Array[Binding[K, V]]()
	var array = Array.FromArray(v.AsArray())
	var iterator = array.GetIterator()
	return iterator
}

// Associative Interface

// This public class method returns the keys for this map.
func (v map_[K, V]) GetKeys() Sequential[K] {
	var List = List[K]()
	var keys = List.Empty()
	var iterator = v.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		keys.AppendValue(association.GetKey())
	}
	return keys
}

// This public class method returns the values associated with the specified
// keys for this map. The values are returned in the same order as the keys in
// the map.
func (v map_[K, V]) GetValues(keys Sequential[K]) Sequential[V] {
	var List = List[V]()
	var values = List.Empty()
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		var key = iterator.GetNext()
		values.AppendValue(v.GetValue(key))
	}
	return values
}

// This public class method returns the value that is associated with the
// specified key in this map.
func (v map_[K, V]) GetValue(key K) V {
	var value = v[key]
	return value
}

// This public class method sets the value associated with the specified key to
// the specified value.
func (v map_[K, V]) SetValue(key K, value V) {
	v[key] = value
}

// This public class method removes the association associated with the
// specified key from the map and returns it.
func (v map_[K, V]) RemoveValue(key K) V {
	var value, exists = v[key]
	if exists {
		delete(v, key)
	}
	return value
}

// This public class method removes the associations associated with the
// specified keys from the map and returns the removed values.
func (v map_[K, V]) RemoveValues(keys Sequential[K]) Sequential[V] {
	var List = List[V]()
	var values = List.Empty()
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		var key = iterator.GetNext()
		values.AppendValue(v.RemoveValue(key))
	}
	return values
}

// This public class method removes all associations from this map.
func (v map_[K, V]) RemoveAll() {
	var keys = v.GetKeys()
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		var key = iterator.GetNext()
		delete(v, key)
	}
}

// Private Interface

// This public class method is used by Go to generate a string from a map.
func (v map_[K, V]) String() string {
	return Formatter().FormatCollection(v)
}
