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
// constants, constructors and functions for the Catalog class namespace.
// NOTE: the Go language requires the key type here support the "comparable"
// interface so we must narrow it down from "any".
type catalogClass_[K comparable, V Value] struct {
	// This class defines no constants.
}

// This private constant defines a map to hold all the singleton references to
// the type specific Catalog class namespaces.
var catalogClassSingletons = map[string]any{}

// This public function returns the singleton reference to a type specific
// Catalog class namespace.  It also initializes any class constants as needed.
// NOTE: the Go language requires the key type here support the "comparable"
// interface so we must narrow it down from "any".
func Catalog[K comparable, V Value]() *catalogClass_[K, V] {
	var class *catalogClass_[K, V]
	var key = fmt.Sprintf("%T", class) // The name of the bound class type.
	var value = catalogClassSingletons[key]
	switch actual := value.(type) {
	case *catalogClass_[K, V]:
		// This bound class type already exists.
		class = actual
	default:
		// Create a new bound class type.
		class = &catalogClass_[K, V]{
			// This class defines no constants.
		}
		catalogClassSingletons[key] = class
	}
	return class
}

// CLASS CONSTRUCTORS

// This public class constructor creates a new empty Catalog.
func (c *catalogClass_[K, V]) Empty() CatalogLike[K, V] {
	var keys = map[Key]Binding[K, V]{}
	var associations = List[Binding[K, V]]().Empty()
	var catalog = &catalog_[K, V]{associations, keys}
	return catalog
}

// This public class constructor creates a new Catalog from the specified Go
// array of values.
func (c *catalogClass_[K, V]) FromArray(
	associations []Binding[K, V],
) CatalogLike[K, V] {
	var array = Array[Binding[K, V]]().FromArray(associations)
	var catalog = c.FromSequence(array)
	return catalog
}

// This public class constructor creates a new Catalog from the specified
// Go map of associations.
func (c *catalogClass_[K, V]) FromMap(
	associations map[K]V,
) CatalogLike[K, V] {
	var catalog = c.Empty()
	for key, value := range associations {
		catalog.SetValue(key, value)
	}
	return catalog
}

// This public class constructor creates a new Catalog from the specified
// sequence of associations.
func (c *catalogClass_[K, V]) FromSequence(
	associations Sequential[Binding[K, V]],
) CatalogLike[K, V] {
	var catalog = c.Empty()
	var iterator = associations.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		var key = association.GetKey()
		var value = association.GetValue()
		catalog.SetValue(key, value)
	}
	return catalog
}

// This public class constructor creates a new Catalog from the specified string
// containing the CDCN definition for the Catalog.
func (c *catalogClass_[K, V]) FromString(source string) CatalogLike[K, V] {
	// First we parse it as a collection of any type value.
	var collection = Parser().ParseCollection([]byte(source)).(Sequential[Binding[Key, Value]])

	// Then we convert it to a Catalog of type Binding[K, V].
	var catalog = c.Empty()
	var iterator = collection.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		var key = association.GetKey().(K)
		var value = association.GetValue().(V)
		catalog.SetValue(key, value)
	}
	return catalog
}

// CLASS FUNCTIONS

// This public class function returns a new Catalog containing only the
// associations that are in the specified Catalog that have the specified keys.
// The associations in the resulting Catalog will be in the same order as the
// specified keys.
func (c *catalogClass_[K, V]) Extract(
	catalog CatalogLike[K, V],
	keys Sequential[K],
) CatalogLike[K, V] {
	var result = c.Empty()
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		var key = iterator.GetNext()
		var value = catalog.GetValue(key)
		result.SetValue(key, value)
	}
	return result
}

// This public class function returns a new Catalog containing all of the
// associations that are in the specified catalogs in the order that they appear
// in each Catalog.  If a key is present in both catalogs, the value of the key
// from the second Catalog takes precedence.
func (c *catalogClass_[K, V]) Merge(
	first CatalogLike[K, V],
	second CatalogLike[K, V],
) CatalogLike[K, V] {
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

// CLASS TYPE

// Encapsulated Type

// This private class type encapsulates a Go structure containing private
// attributes that can only be accessed and manipulated using methods that
// implement the catalog-like abstract type.  A catalog-like type maintains
// key-value pair associations. This type is parameterized as follows:
//   - K is a primitive type of key.
//   - V is any type of entity.
type catalog_[K Key, V Value] struct {
	associations ListLike[Binding[K, V]]
	keys         map[Key]Binding[K, V]
}

// Associative Interface

// This public class method returns the keys for this Catalog.
func (v *catalog_[K, V]) GetKeys() Sequential[K] {
	var keys = List[K]().Empty()
	var iterator = v.associations.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		keys.AppendValue(association.GetKey())
	}
	return keys
}

// This public class method returns the value that is associated with the
// specified key in this Catalog.
func (v *catalog_[K, V]) GetValue(key K) V {
	var value V // Set the return value to its zero value.
	var association, exists = v.keys[key]
	if exists {
		// Extract the value.
		value = association.GetValue()
	}
	return value
}

// This public class method returns the values associated with the specified
// keys for this Catalog. The values are returned in the same order as the keys
// in the Catalog.
func (v *catalog_[K, V]) GetValues(keys Sequential[K]) Sequential[V] {
	var values = List[V]().Empty()
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		var key = iterator.GetNext()
		values.AppendValue(v.GetValue(key))
	}
	return values
}

// This public class method removes all associations from this Catalog.
func (v *catalog_[K, V]) RemoveAll() {
	v.keys = map[Key]Binding[K, V]{}
	v.associations.RemoveAll()
}

// This public class method removes the association associated with the
// specified key from the Catalog and returns it.
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
// specified keys from the Catalog and returns the removed values.
func (v *catalog_[K, V]) RemoveValues(keys Sequential[K]) Sequential[V] {
	var values = List[V]().Empty()
	var Iterator = Iterator[K]()
	var iterator = Iterator.FromSequence(keys)
	for iterator.HasNext() {
		var key = iterator.GetNext()
		values.AppendValue(v.RemoveValue(key))
	}
	return values
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

// Sequential Interface

// This public class method returns all the associations in this Catalog. The
// associations
// retrieved are in the same order as they are in the Catalog.
func (v *catalog_[K, V]) AsArray() []Binding[K, V] {
	return v.associations.AsArray()
}

// This public class method generates for this Catalog an iterator that can be
// used to traverse its associations.
func (v *catalog_[K, V]) GetIterator() Ratcheted[Binding[K, V]] {
	return v.associations.GetIterator()
}

// This public class method returns the number of associations contained in this
// Catalog.
func (v *catalog_[K, V]) GetSize() int {
	return v.associations.GetSize()
}

// This public class method determines whether or not this Catalog is empty.
func (v *catalog_[K, V]) IsEmpty() bool {
	return v.associations.IsEmpty()
}

// Sortable Interface

// This public class method reverses the order of all associations in this
// Catalog.
func (v *catalog_[K, V]) ReverseValues() {
	v.associations.ReverseValues()
}

// This public class method pseudo-randomly shuffles the order of all
// associations in this Catalog.
func (v *catalog_[K, V]) ShuffleValues() {
	v.associations.ShuffleValues()
}

// This public class method sorts this Catalog using the canonical rank function
// to compare the keys.
func (v *catalog_[K, V]) SortValues() {
	v.associations.SortValues()
}

// This public class method sorts this Catalog using the specified ranking
// function to compare the keys.
func (v *catalog_[K, V]) SortValuesWithRanker(ranking RankingFunction) {
	v.associations.SortValuesWithRanker(ranking)
}

// Private Interface

// This public class method returns the canonical string for this Catalog.
func (v *catalog_[K, V]) String() string {
	return Formatter().FormatCollection(v)
}
