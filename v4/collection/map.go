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
	sts "strings"
	syn "sync"
)

// CLASS ACCESS

// Reference

var mapClass = map[string]any{}
var mapMutex syn.Mutex

// Function

// NOTE:
// The Go language requires the key type here support the "comparable" interface
// so we must narrow it down from type Key (i.e. ."any").
func Map[K comparable, V Value](notation NotationLike) MapClassLike[K, V] {
	// Validate the notation argument.
	if notation == nil {
		panic("A notation must be specified when creating this class.")
	}

	// Generate the name of the bound class type.
	var class MapClassLike[K, V]
	var name = fmt.Sprintf("%T-%T", class, notation)

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

// NOTE:
// The Go language requires the key type here support the "comparable" interface
// so we must narrow it down from type Key (i.e. ."any").
type mapClass_[K comparable, V Value] struct {
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

func (c *mapClass_[K, V]) MakeFromSource(source string) MapLike[K, V] {
	// First we parse it as a collection of any type value.
	var collection = c.notation_.ParseSource(source).(Sequential[AssociationLike[Key, Value]])

	// Next we must convert each value explicitly to type AssociationLike[K, V].
	var anys = collection.AsArray()
	var array = make([]AssociationLike[K, V], len(anys))
	for index, association := range anys {
		var key = association.GetKey().(K)
		var value = association.GetValue().(V)
		array[index] = Association[K, V]().MakeWithAttributes(key, value)
	}

	// Then we can create the stack from the type AssociationLike[K, V] array.
	return c.MakeFromArray(array)
}

// INSTANCE METHODS

// Target

// NOTE:
// The Go language requires the key type here support the "comparable" interface
// so we must narrow it down from type Key (i.e. ."any").
type map_[K comparable, V Value] map[K]V

// Associative

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

func (v map_[K, V]) GetValue(key K) V {
	var value = v[key]
	return value
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

func (v map_[K, V]) SetValue(key K, value V) {
	v[key] = value
}

// Sequential

func (v map_[K, V]) AsArray() []AssociationLike[K, V] {
	var size = len(v)
	var array = make([]AssociationLike[K, V], size)
	var index = 0
	for key, value := range v {
		var association = Association[K, V]().MakeWithAttributes(key, value)
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

func (v map_[K, V]) GetSize() int {
	return len(v)
}

func (v map_[K, V]) IsEmpty() bool {
	return len(v) == 0
}

// Stringer

// NOTE:
// Since this class only extends the primitive map type it cannot have any
// attributes assigned to it.  This means that we have no way of accessing its
// notation.  So we cannot use the notation specific formatter to generate the
// string value for this map and must generate it manually.  This is only a
// problem when this method is called directly—as done by the fmt.Sprintf()
// method.  The formatters themselves can handle the formatting of maps just
// fine.
func (v map_[K, V]) String() string {
	var string_ = "["
	if v.IsEmpty() {
		string_ += (":")
	} else {
		var builder sts.Builder
		var iterator = v.GetIterator()
		for iterator.HasNext() {
			var association = iterator.GetNext()
			var key = association.GetKey()
			var value = association.GetValue()
			builder.WriteString(fmt.Sprintf("%#v: %#v, ", key, value))
		}
		var last = builder.Len() - 2
		string_ += builder.String()[:last]
	}
	string_ += "](Map)\n"
	return string_
}
