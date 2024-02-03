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

// CLASS ACCESS

// Reference

var catalogClass = map[string]any{}

// Function

func Catalog[K comparable, V Value]() CatalogClassLike[K, V] {
	// Generate the name of the bound class type.
	var class CatalogClassLike[K, V]
	var name = fmt.Sprintf("%T", class)

	// Check for existing bound class type.
	var value = catalogClass[name]
	switch actual := value.(type) {
	case *catalogClass_[K, V]:
		// This bound class type already exists.
		class = actual
	default:
		// Add a new bound class type.
		class = &catalogClass_[K, V]{
			// This class defines no constants.
		}
		catalogClass[name] = class
	}

	// Return a reference to the bound class type.
	return class
}

// CLASS METHODS

// Target

/*
The Go language requires the key type here support the "comparable" interface so
we must narrow it down from type Key (i.e "any").
*/
type catalogClass_[K comparable, V Value] struct {
	// This class defines no constants.
}

// Constructors

func (c *catalogClass_[K, V]) Make() CatalogLike[K, V] {
	var keys = map[K]AssociationLike[K, V]{}
	var associations = List[AssociationLike[K, V]]().Make()
	var catalog = &catalog_[K, V]{associations, keys}
	return catalog
}

func (c *catalogClass_[K, V]) MakeFromArray(
	associations []AssociationLike[K, V],
) CatalogLike[K, V] {
	var array = Array[AssociationLike[K, V]]().MakeFromArray(associations)
	var catalog = c.MakeFromSequence(array)
	return catalog
}

func (c *catalogClass_[K, V]) MakeFromMap(
	associations map[K]V,
) CatalogLike[K, V] {
	var catalog = c.Make()
	for key, value := range associations {
		catalog.SetValue(key, value)
	}
	return catalog
}

func (c *catalogClass_[K, V]) MakeFromSequence(
	associations Sequential[AssociationLike[K, V]],
) CatalogLike[K, V] {
	var catalog = c.Make()
	var iterator = associations.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		var key = association.GetKey()
		var value = association.GetValue()
		catalog.SetValue(key, value)
	}
	return catalog
}

func (c *catalogClass_[K, V]) MakeFromSource(
	source string,
	notation NotationLike,
) CatalogLike[K, V] {
	// First we parse it as a collection of any type value.
	var collection = notation.ParseSource(source).(Sequential[AssociationLike[Key, Value]])

	// Then we convert it to a catalog of type AssociationLike[K, V].
	var catalog = c.Make()
	var iterator = collection.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		var key = association.GetKey().(K)
		var value = association.GetValue().(V)
		catalog.SetValue(key, value)
	}
	return catalog
}

// Functions

/*
This public class function returns a new catalog containing only the
associations that are in the specified catalog that have the specified keys.
The associations in the resulting catalog will be in the same order as the
specified keys.
*/
func (c *catalogClass_[K, V]) Extract(
	catalog CatalogLike[K, V],
	keys Sequential[K],
) CatalogLike[K, V] {
	var result = c.Make()
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		var key = iterator.GetNext()
		var value = catalog.GetValue(key)
		result.SetValue(key, value)
	}
	return result
}

/*
This public class function returns a new catalog containing all of the
associations that are in the specified Catalogs in the order that they appear in
each catalog.  If a key is present in both Catalogs, the value of the key from
the second catalog takes precedence.
*/
func (c *catalogClass_[K, V]) Merge(
	first CatalogLike[K, V],
	second CatalogLike[K, V],
) CatalogLike[K, V] {
	var catalog = c.MakeFromSequence(first)
	var iterator = second.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		var key = association.GetKey()
		var value = association.GetValue()
		catalog.SetValue(key, value)
	}
	return catalog
}

// INSTANCE METHODS

// Target

type catalog_[K comparable, V Value] struct {
	associations ListLike[AssociationLike[K, V]]
	keys         map[K]AssociationLike[K, V]
}

// Associative

func (v *catalog_[K, V]) GetKeys() Sequential[K] {
	var keys = List[K]().Make()
	var iterator = v.associations.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		keys.AppendValue(association.GetKey())
	}
	return keys
}

func (v *catalog_[K, V]) GetValue(key K) V {
	var value V // Set the return value to its zero value.
	var association, exists = v.keys[key]
	if exists {
		// Extract the value.
		value = association.GetValue()
	}
	return value
}

func (v *catalog_[K, V]) GetValues(keys Sequential[K]) Sequential[V] {
	var values = List[V]().Make()
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		var key = iterator.GetNext()
		values.AppendValue(v.GetValue(key))
	}
	return values
}

func (v *catalog_[K, V]) RemoveAll() {
	v.keys = map[K]AssociationLike[K, V]{}
	v.associations.RemoveAll()
}

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

func (v *catalog_[K, V]) RemoveValues(keys Sequential[K]) Sequential[V] {
	var values = List[V]().Make()
	var iterator = Iterator[K]().Make(keys)
	for iterator.HasNext() {
		var key = iterator.GetNext()
		values.AppendValue(v.RemoveValue(key))
	}
	return values
}

func (v *catalog_[K, V]) SetValue(key K, value V) {
	var association, exists = v.keys[key]
	if exists {
		// Set the value of an existing association.
		association.SetValue(value)
	} else {
		// Add a new association.
		association = Association[K, V]().Make(key, value)
		v.associations.AppendValue(association)
		v.keys[key] = association
	}
}

// Sequential

func (v *catalog_[K, V]) AsArray() []AssociationLike[K, V] {
	return v.associations.AsArray()
}

func (v *catalog_[K, V]) GetIterator() IteratorLike[AssociationLike[K, V]] {
	return v.associations.GetIterator()
}

func (v *catalog_[K, V]) GetSize() int {
	return v.associations.GetSize()
}

func (v *catalog_[K, V]) IsEmpty() bool {
	return v.associations.IsEmpty()
}

// Sortable

func (v *catalog_[K, V]) ReverseValues() {
	v.associations.ReverseValues()
}

func (v *catalog_[K, V]) ShuffleValues() {
	v.associations.ShuffleValues()
}

func (v *catalog_[K, V]) SortValues() {
	v.associations.SortValues()
}

func (v *catalog_[K, V]) SortValuesWithRanker(ranker RankingFunction) {
	v.associations.SortValuesWithRanker(ranker)
}

// Stringer

func (v *catalog_[K, V]) String() string {
	var formatter = Formatter().Make()
	return formatter.FormatCollection(v)
}
