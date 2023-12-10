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

// This private type defines the namespace structure associated with the constants,
// constructors and functions for the catalog class namespace.
type catalogClass_[K Key, V Value] struct {
	// This class defines no constants.
}

// This private constant defines a map to hold all the singleton references to
// the type specific catalog namespaces.
var catalogClassSingletons = map[string]any{}

// This public function returns the singleton reference to a type specific
// catalog namespace.  It also initializes any class constants as needed.
func Catalog[K Key, V Value]() *catalogClass_[K, V] {
	var class *catalogClass_[K, V]
	var key = fmt.Sprintf("%T", class)
	var value = catalogClassSingletons[key]
	switch actual := value.(type) {
	case *catalogClass_[K, V]:
		class = actual
	default:
		class = &catalogClass_[K, V]{
			// This class defines no constants.
		}
		catalogClassSingletons[key] = class
	}
	return class
}

// CLASS CONSTRUCTORS

// This public class constructor creates a new empty catalog.
func (c *catalogClass_[K, V]) FromNothing() CatalogLike[K, V] {
	var keys = map[Key]Binding[K, V]{}
	var List = List[Binding[K, V]]()
	var associations = List.FromNothing()
	var catalog = &catalog_[K, V]{associations, associations, keys}
	return catalog
}

// This public class constructor creates a new catalog from the specified
// array of associations.
func (c *catalogClass_[K, V]) FromArray(array []Binding[K, V]) CatalogLike[K, V] {
	var catalog = c.FromNothing()
	for _, association := range array {
		var key = association.GetKey()
		var value = association.GetValue()
		catalog.SetValue(key, value)
	}
	return catalog
}

// This public class constructor creates a new catalog from the specified
// sequence of associations.
func (c *catalogClass_[K, V]) FromSequence(sequence Sequential[Binding[K, V]]) CatalogLike[K, V] {
	var catalog = c.FromNothing()
	var iterator = sequence.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		var key = association.GetKey()
		var value = association.GetValue()
		catalog.SetValue(key, value)
	}
	return catalog
}

// CLASS FUNCTIONS

// This public class function returns a new catalog containing all of the
// associations that are in the specified catalogs in the order that they appear
// in each catalog.  If a key is present in both catalogs, the value of the key
// from the second catalog takes precedence.
func (c *catalogClass_[K, V]) Merge(first, second CatalogLike[K, V]) CatalogLike[K, V] {
	var catalog = c.FromSequence(first)
	var iterator = second.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		var key = association.GetKey()
		var value = association.GetValue()
		catalog.SetValue(key, value)
	}
	return catalog
}

// This public class function returns a new catalog containing only the
// associations that are in the specified catalog that have the specified keys.
// The associations in the resulting catalog will be in the same order as the
// specified keys.
func (c *catalogClass_[K, V]) Extract(catalog CatalogLike[K, V], keys Sequential[K]) CatalogLike[K, V] {
	var result = c.FromNothing()
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		var key = iterator.GetNext()
		var value = catalog.GetValue(key)
		result.SetValue(key, value)
	}
	return result
}

// CLASS TYPE

// Encapsulated Type

// This private class type encapsulates a Go structure containing private
// attributes that can only be accessed and manipulated using methods that
// implement the catalog-like abstract type.  A catalog-like type maintains
// key-value pair associations. This type is parameterized as follows:
//   - K is a primitive type of key.
//   - V is any type of entity.
type catalog_[K Key, V Value] struct {
	// Note: The delegated methods don't see the real collection type.
	Sequential[Binding[K, V]]
	associations ListLike[Binding[K, V]]
	keys         map[Key]Binding[K, V]
}

// Associative Interface

// This public class method returns the keys for this catalog.
func (v *catalog_[K, V]) GetKeys() Sequential[K] {
	var List = List[K]()
	var keys = List.FromNothing()
	var iterator = v.associations.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		keys.AppendValue(association.GetKey())
	}
	return keys
}

// This public class method returns the values associated with the specified
// keys for this catalog. The values are returned in the same order as the keys
// in the catalog.
func (v *catalog_[K, V]) GetValues(keys Sequential[K]) Sequential[V] {
	var List = List[V]()
	var values = List.FromNothing()
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		var key = iterator.GetNext()
		values.AppendValue(v.GetValue(key))
	}
	return values
}

// This public class method returns the value that is associated with the
// specified key in this catalog.
func (v *catalog_[K, V]) GetValue(key K) V {
	var value V // Set the return value to its zero value.
	var association, exists = v.keys[key]
	if exists {
		// Extract the value.
		value = association.GetValue()
	}
	return value
}

// This public class method sets the value associated with the specified key to
// the specified value.
func (v *catalog_[K, V]) SetValue(key K, value V) {
	var association, exists = v.keys[key]
	if exists {
		// Set the value of an existing association.
		association.SetValue(value)
	} else {
		// Add a new association.
		var Association = Association[K, V]()
		association = Association.FromPair(key, value)
		v.associations.AppendValue(association)
		v.keys[key] = association
	}
}

// This public class method removes the association associated with the
// specified key from the catalog and returns it.
func (v *catalog_[K, V]) RemoveValue(key K) V {
	var old V // Set the return value to its zero value.
	var association, exists = v.keys[key]
	if exists {
		var index = v.associations.GetIndex(association)
		v.associations.RemoveValue(index)
		old = association.GetValue()
		delete(v.keys, key)
	}
	return old
}

// This public class method removes the associations associated with the
// specified keys from the catalog and returns the removed values.
func (v *catalog_[K, V]) RemoveValues(keys Sequential[K]) Sequential[V] {
	var List = List[V]()
	var values = List.FromNothing()
	var Iterator = Iterator[K]()
	var iterator = Iterator.FromSequence(keys)
	for iterator.HasNext() {
		var key = iterator.GetNext()
		values.AppendValue(v.RemoveValue(key))
	}
	return values
}

// This public class method removes all associations from this catalog.
func (v *catalog_[K, V]) RemoveAll() {
	v.keys = map[Key]Binding[K, V]{}
	v.associations.RemoveAll()
}

// Sortable Interface

// This public class method sorts this catalog using the canonical rank function
// to compare the keys.
func (v *catalog_[K, V]) SortValues() {
	v.associations.SortValues()
}

// This public class method sorts this catalog using the specified ranking
// function to compare the keys.
func (v *catalog_[K, V]) SortValuesWithRanker(ranking RankingFunction) {
	v.associations.SortValuesWithRanker(ranking)
}

// This public class method reverses the order of all associations in this
// catalog.
func (v *catalog_[K, V]) ReverseValues() {
	v.associations.ReverseValues()
}

// This public class method pseudo-randomly shuffles the order of all
// associations in this catalog.
func (v *catalog_[K, V]) ShuffleValues() {
	v.associations.ShuffleValues()
}

// Go Stringer Interface

// This public class method returns the canonical string for this catalog.
func (v *catalog_[K, V]) String() string {
	return FormatCollection(v)
}
