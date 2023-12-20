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
// constants, constructors and functions for the Map class namespace.
// NOTE: the Go language requires the key type here support the "comparable"
// interface so we must narrow it down from "any".
type mapClass_[K comparable, V Value] struct {
	// This class defines no constants.
}

// This private constant defines a Map to hold all the singleton references to
// the type specific Map class namespaces.
var mapClassSingletons = map[string]any{}

// This public function returns the singleton reference to a type specific
// Map class namespace.  It also initializes any class constants as needed.
// NOTE: the Go language requires the key type here support the "comparable"
// interface so we must narrow it down from "any".
func Map[K comparable, V Value]() *mapClass_[K, V] {
	var class *mapClass_[K, V]
	var key = fmt.Sprintf("%T", class) // The name of the bound class type.
	var value = mapClassSingletons[key]
	switch actual := value.(type) {
	case *mapClass_[K, V]:
		// This bound class type already exists.
		class = actual
	default:
		// Create a new bound class type.
		class = &mapClass_[K, V]{
			// This class defines no constants.
		}
		mapClassSingletons[key] = class
	}
	return class
}

// CLASS CONSTRUCTORS

// This public class constructor creates a new empty Map.
func (c *mapClass_[K, V]) Empty() MapLike[K, V] {
	return map_[K, V](map[K]V{})
}

// This public class constructor creates a new Map from the specified
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

// This public class constructor creates a new Map from the specified
// Go map of associations.
func (c *mapClass_[K, V]) FromMap(associations map[K]V) MapLike[K, V] {
	var length = len(associations)
	var duplicate = make(map[K]V, length)
	for key, value := range associations {
		duplicate[key] = value
	}
	return map_[K, V](duplicate)
}

// This public class constructor creates a new Map from the specified string
// containing the CDCN definition for the Map.
func (c *mapClass_[K, V]) FromString(source string) MapLike[K, V] {
	// First we parse it as a collection of any type value.
	var collection = Parser().ParseCollection([]byte(source)).(Sequential[Binding[Key, Value]])

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

// Extended Type

// This private class type extends the primitive Go "map[K]V" data type and
// defines the methods that implement the map-like abstract type.  The Map
// extended type manages key-value pair associations.
// This type is parameterized as follows:
//   - K is a primitive type of key.
//   - V is any type of entity.
type map_[K comparable, V Value] map[K]V

// Associative Interface

// This public class method returns the keys for this Map.
func (v map_[K, V]) GetKeys() Sequential[K] {
	var keys = List[K]().Empty()
	var iterator = v.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		keys.AppendValue(association.GetKey())
	}
	return keys
}

// This public class method returns the value that is associated with the
// specified key in this Map.
func (v map_[K, V]) GetValue(key K) V {
	var value = v[key]
	return value
}

// This public class method returns the values associated with the specified
// keys for this Map. The values are returned in the same order as the keys in
// the Map.
func (v map_[K, V]) GetValues(keys Sequential[K]) Sequential[V] {
	var values = List[V]().Empty()
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		var key = iterator.GetNext()
		values.AppendValue(v.GetValue(key))
	}
	return values
}

// This public class method removes all associations from this Map.
func (v map_[K, V]) RemoveAll() {
	var keys = v.GetKeys()
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		var key = iterator.GetNext()
		delete(v, key)
	}
}

// This public class method removes the association associated with the
// specified key from the Map and returns it.
func (v map_[K, V]) RemoveValue(key K) V {
	var value, exists = v[key]
	if exists {
		delete(v, key)
	}
	return value
}

// This public class method removes the associations associated with the
// specified keys from the Map and returns the removed values.
func (v map_[K, V]) RemoveValues(keys Sequential[K]) Sequential[V] {
	var values = List[V]().Empty()
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		var key = iterator.GetNext()
		values.AppendValue(v.RemoveValue(key))
	}
	return values
}

// This public class method sets the value associated with the specified key to
// the specified value.
func (v map_[K, V]) SetValue(key K, value V) {
	v[key] = value
}

// Sequential Interface

// This public class method returns all the associations that are in this Map.
// The associations retrieved are in the same order as they are in the Map.
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

// This public class method generates for this Map an iterator that can be used
// to traverse its associations.
func (v map_[K, V]) GetIterator() Ratcheted[Binding[K, V]] {
	var Array = Array[Binding[K, V]]()
	var array = Array.FromArray(v.AsArray())
	var iterator = array.GetIterator()
	return iterator
}

// This public class method returns the number of values contained in this Map.
func (v map_[K, V]) GetSize() int {
	return len(v)
}

// This public class method determines whether or not this Map is empty.
func (v map_[K, V]) IsEmpty() bool {
	return len(v) == 0
}

// Private Interface

// This public class method is used by Go to generate a string from a Map.
func (v map_[K, V]) String() string {
	return Formatter().FormatCollection(v)
}
