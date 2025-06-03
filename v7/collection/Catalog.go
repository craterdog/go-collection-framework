/*
................................................................................
.    Copyright (c) 2009-2025 Crater Dog Technologies.  All Rights Reserved.    .
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
	age "github.com/craterdog/go-collection-framework/v7/agent"
	uti "github.com/craterdog/go-missing-utilities/v7"
	syn "sync"
)

// CLASS INTERFACE

// Access Function

func CatalogClass[K comparable, V any]() CatalogClassLike[K, V] {
	return catalogClass[K, V]()
}

// Constructor Methods

func (c *catalogClass_[K, V]) Catalog() CatalogLike[K, V] {
	var listClass = ListClass[AssociationLike[K, V]]()
	var keys = map[K]AssociationLike[K, V]{}
	var associations = listClass.List()
	var instance = &catalog_[K, V]{
		// Initialize the instance attributes.
		keys_:         keys,
		associations_: associations,
	}
	return instance
}

func (c *catalogClass_[K, V]) CatalogFromArray(
	associations []AssociationLike[K, V],
) CatalogLike[K, V] {
	var catalog = c.Catalog()
	for _, association := range associations {
		var key = association.GetKey()
		var value = association.GetValue()
		catalog.SetValue(key, value)
	}
	return catalog
}

func (c *catalogClass_[K, V]) CatalogFromMap(
	associations map[K]V,
) CatalogLike[K, V] {
	// NOTE:
	// The ordering of the key-value associations in the specified intrinsic Go
	// map data type is non-deterministic, even using the same associations
	// across multiple runs.  To make this constructor deterministic we sort the
	// specified map associations using their "natural" ordering.
	var catalog = c.Catalog()
	for key, value := range associations {
		catalog.SetValue(key, value)
	}
	catalog.SortValues()
	return catalog
}

func (c *catalogClass_[K, V]) CatalogFromSequence(
	associations Sequential[AssociationLike[K, V]],
) CatalogLike[K, V] {
	var catalog = c.Catalog()
	var iterator = associations.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		var key = association.GetKey()
		var value = association.GetValue()
		catalog.SetValue(key, value)
	}
	return catalog
}

// Constant Methods

// Function Methods

func (c *catalogClass_[K, V]) Extract(
	catalog CatalogLike[K, V],
	keys Sequential[K],
) CatalogLike[K, V] {
	var result = c.Catalog()
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		var key = iterator.GetNext()
		var value = catalog.GetValue(key)
		result.SetValue(key, value)
	}
	return result
}

func (c *catalogClass_[K, V]) Merge(
	first CatalogLike[K, V],
	second CatalogLike[K, V],
) CatalogLike[K, V] {
	var catalog = c.CatalogFromSequence(first)
	var iterator = second.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		var key = association.GetKey()
		var value = association.GetValue()
		catalog.SetValue(key, value)
	}
	return catalog
}

// INSTANCE INTERFACE

// Principal Methods

func (v *catalog_[K, V]) GetClass() CatalogClassLike[K, V] {
	return catalogClass[K, V]()
}

// Attribute Methods

// Associative[K, V] Methods

func (v *catalog_[K, V]) AsMap() map[K]V {
	var map_ = map[K]V{}
	var iterator = v.associations_.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		var key = association.GetKey()
		var value = association.GetValue()
		map_[key] = value
	}
	return map_
}

func (v *catalog_[K, V]) GetValue(
	key K,
) V {
	var value V // Set the return value to its zero value.
	var association, exists = v.keys_[key]
	if exists {
		// Extract the value.
		value = association.GetValue()
	}
	return value
}

func (v *catalog_[K, V]) SetValue(
	key K,
	value V,
) {
	var association, exists = v.keys_[key]
	if exists {
		// Set the value of an existing association.
		association.SetValue(value)
	} else {
		// Add a new association.
		var associationClass = AssociationClass[K, V]()
		association = associationClass.Association(key, value)
		v.associations_.AppendValue(association)
		v.keys_[key] = association
	}
}

func (v *catalog_[K, V]) GetKeys() Sequential[K] {
	var listClass = ListClass[K]()
	var keys = listClass.List()
	var iterator = v.associations_.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		keys.AppendValue(association.GetKey())
	}
	return keys
}

func (v *catalog_[K, V]) GetValues(
	keys Sequential[K],
) Sequential[V] {
	var listClass = ListClass[V]()
	var values = listClass.List()
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		var key = iterator.GetNext()
		values.AppendValue(v.GetValue(key))
	}
	return values
}

func (v *catalog_[K, V]) RemoveValue(
	key K,
) V {
	var old V // Set the return value to its zero value.
	var association, exists = v.keys_[key]
	if exists {
		var index = v.associations_.GetIndex(association)
		v.associations_.RemoveValue(index)
		old = association.GetValue()
		delete(v.keys_, key)
	}
	return old
}

func (v *catalog_[K, V]) RemoveValues(
	keys Sequential[K],
) Sequential[V] {
	var listClass = ListClass[V]()
	var values = listClass.List()
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		var key = iterator.GetNext()
		values.AppendValue(v.RemoveValue(key))
	}
	return values
}

func (v *catalog_[K, V]) RemoveAll() {
	v.keys_ = map[K]AssociationLike[K, V]{}
	v.associations_.RemoveAll()
}

// Sequential[AssociationLike[K, V]] Methods

func (v *catalog_[K, V]) IsEmpty() bool {
	return v.associations_.IsEmpty()
}

func (v *catalog_[K, V]) GetSize() uti.Cardinal {
	var size = v.associations_.GetSize()
	return size
}

func (v *catalog_[K, V]) AsArray() []AssociationLike[K, V] {
	var array = v.associations_.AsArray()
	return array
}

func (v *catalog_[K, V]) GetIterator() age.IteratorLike[AssociationLike[K, V]] {
	var iterator = v.associations_.GetIterator()
	return iterator
}

// Sortable[AssociationLike[K, V]] Methods

func (v *catalog_[K, V]) SortValues() {
	v.associations_.SortValues()
}

func (v *catalog_[K, V]) SortValuesWithRanker(
	ranker age.RankingFunction[AssociationLike[K, V]],
) {
	v.associations_.SortValuesWithRanker(ranker)
}

func (v *catalog_[K, V]) ReverseValues() {
	v.associations_.ReverseValues()
}

func (v *catalog_[K, V]) ShuffleValues() {
	v.associations_.ShuffleValues()
}

// Stringer Methods

func (v *catalog_[K, V]) String() string {
	return uti.Format(v)
}

// PROTECTED INTERFACE

// Private Methods

// Instance Structure

type catalog_[K comparable, V any] struct {
	// Declare the instance attributes.
	associations_ ListLike[AssociationLike[K, V]]
	keys_         map[K]AssociationLike[K, V]
}

// Class Structure

type catalogClass_[K comparable, V any] struct {
	// Declare the class constants.
}

// Class Reference

var catalogMap_ = map[string]any{}
var catalogMutex_ syn.Mutex

func catalogClass[K comparable, V any]() *catalogClass_[K, V] {
	// Generate the name of the bound class type.
	var class *catalogClass_[K, V]
	var name = fmt.Sprintf("%T", class)

	// Check for an existing bound class type.
	catalogMutex_.Lock()
	var value = catalogMap_[name]
	switch actual := value.(type) {
	case *catalogClass_[K, V]:
		// This bound class type already exists.
		class = actual
	default:
		// Add a new bound class type.
		class = &catalogClass_[K, V]{
			// Initialize the class constants.
		}
		catalogMap_[name] = class
	}
	catalogMutex_.Unlock()

	// Return a reference to the bound class type.
	return class
}
