/*******************************************************************************
 *   Copyright (c) 2009-2022 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package collections

// MAP IMPLEMENTATION

// This type defines the structure and methods associated with a native map of
// key-value pair associations. This type is parameterized as follows:
//   - K is a primitive type of key.
//   - V is any type of entity.
type Map[K comparable, V Value] map[K]V

// This constructor creates a new catalog from the specified sequence of
// associations.
func MapFromSequence[K comparable, V Value](sequence Sequential[Binding[K, V]]) map[K]V {
	var v = make(map[K]V)
	var iterator = sequence.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		var key = association.GetKey()
		var value = association.GetValue()
		v[key] = value
	}
	return v
}

// SEQUENTIAL INTERFACE

// This method determines whether or not this map is empty.
func (v Map[K, V]) IsEmpty() bool {
	return len(v) == 0
}

// This method returns the number of values contained in this map.
func (v Map[K, V]) GetSize() int {
	return len(v)
}

// This method returns all the associations that are in this map. The
// associations retrieved are in the same order as they are in the map.
func (v Map[K, V]) AsArray() []Binding[K, V] {
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

// This public class method generates for this map an iterator that can be
// used to traverse its associations.
func (v Map[K, V]) GetIterator() Ratcheted[Binding[K, V]] {
	var Array = Array[Binding[K, V]]()
	var array = Array.FromArray(v.AsArray())
	var iterator = array.GetIterator()
	return iterator
}

// ASSOCIATIVE INTERFACE

// This method returns the keys for this map.
func (v Map[K, V]) GetKeys() Sequential[K] {
	var List = List[K]()
	var keys = List.FromNothing()
	var iterator = v.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		keys.AppendValue(association.GetKey())
	}
	return keys
}

// This method returns the values associated with the specified keys for this
// map. The values are returned in the same order as the keys in the map.
func (v Map[K, V]) GetValues(keys Sequential[K]) Sequential[V] {
	var List = List[V]()
	var values = List.FromNothing()
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		var key = iterator.GetNext()
		values.AppendValue(v.GetValue(key))
	}
	return values
}

// This method returns the value that is associated with the specified key in
// this map.
func (v Map[K, V]) GetValue(key K) V {
	var value = v[key]
	return value
}

// This method sets the value associated with the specified key to the
// specified value.
func (v Map[K, V]) SetValue(key K, value V) {
	v[key] = value
}

// This method removes the association associated with the specified key from the
// map and returns it.
func (v Map[K, V]) RemoveValue(key K) V {
	var value, exists = v[key]
	if exists {
		delete(v, key)
	}
	return value
}

// This method removes the associations associated with the specified keys from
// the map and returns the removed values.
func (v Map[K, V]) RemoveValues(keys Sequential[K]) Sequential[V] {
	var List = List[V]()
	var values = List.FromNothing()
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		var key = iterator.GetNext()
		values.AppendValue(v.RemoveValue(key))
	}
	return values
}

// This method removes all associations from this map.
func (v Map[K, V]) RemoveAll() {
	var keys = v.GetKeys()
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		var key = iterator.GetNext()
		delete(v, key)
	}
}

// GO INTERFACE

func (v Map[K, V]) String() string {
	return FormatCollection(v)
}
