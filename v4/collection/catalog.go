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
	syn "sync"
)

// CLASS ACCESS

// Reference

var catalogClass = map[string]any{}
var catalogMutex syn.Mutex

// Function

func Catalog[K comparable, V any](notation NotationLike) CatalogClassLike[K, V] {
	// Generate the name of the bound class type.
	var class *catalogClass_[K, V]
	var name = fmt.Sprintf("%T", class)

	// Check for existing bound class type.
	catalogMutex.Lock()
	var value = catalogClass[name]
	switch actual := value.(type) {
	case *catalogClass_[K, V]:
		// This bound class type already exists.
		class = actual
	default:
		// Add a new bound class type.
		class = &catalogClass_[K, V]{
			// Initialize the class constants.
			notation_: notation,
		}
		catalogClass[name] = class
	}
	catalogMutex.Unlock()

	// Return a reference to the bound class type.
	return class
}

// CLASS METHODS

// Target

type catalogClass_[K comparable, V any] struct {
	// Define the class constants.
	notation_ NotationLike
}

// Constants

func (c *catalogClass_[K, V]) Notation() NotationLike {
	return c.notation_
}

// Constructors

func (c *catalogClass_[K, V]) Make() CatalogLike[K, V] {
	var keys = map[K]AssociationLike[K, V]{}
	var associations = List[AssociationLike[K, V]](c.notation_).Make()
	return &catalog_[K, V]{
		c,
		associations,
		keys,
	}
}

func (c *catalogClass_[K, V]) MakeFromArray(
	associations []AssociationLike[K, V],
) CatalogLike[K, V] {
	var array = Array[AssociationLike[K, V]](c.notation_).MakeFromArray(associations)
	return c.MakeFromSequence(array)
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

// Functions

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

type catalog_[K comparable, V any] struct {
	class_        CatalogClassLike[K, V]
	associations_ ListLike[AssociationLike[K, V]]
	keys_         map[K]AssociationLike[K, V]
}

// Attributes

func (v *catalog_[K, V]) GetClass() CatalogClassLike[K, V] {
	return v.class_
}

// Associative

func (v *catalog_[K, V]) GetValue(key K) V {
	var value V // Set the return value to its zero value.
	var association, exists = v.keys_[key]
	if exists {
		// Extract the value.
		value = association.GetValue()
	}
	return value
}

func (v *catalog_[K, V]) SetValue(key K, value V) {
	var association, exists = v.keys_[key]
	if exists {
		// Set the value of an existing association.
		association.SetValue(value)
	} else {
		// Add a new association.
		association = Association[K, V](v.GetClass().Notation()).Make(key, value)
		v.associations_.AppendValue(association)
		v.keys_[key] = association
	}
}

func (v *catalog_[K, V]) GetKeys() Sequential[K] {
	var keys = List[K](v.GetClass().Notation()).Make()
	var iterator = v.associations_.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		keys.AppendValue(association.GetKey())
	}
	return keys
}

func (v *catalog_[K, V]) GetValues(keys Sequential[K]) Sequential[V] {
	var values = List[V](v.GetClass().Notation()).Make()
	var iterator = keys.GetIterator()
	for iterator.HasNext() {
		var key = iterator.GetNext()
		values.AppendValue(v.GetValue(key))
	}
	return values
}

func (v *catalog_[K, V]) RemoveValue(key K) V {
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

func (v *catalog_[K, V]) RemoveValues(keys Sequential[K]) Sequential[V] {
	var values = List[V](v.GetClass().Notation()).Make()
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

// Sequential

func (v *catalog_[K, V]) IsEmpty() bool {
	return v.associations_.IsEmpty()
}

func (v *catalog_[K, V]) GetSize() int {
	return v.associations_.GetSize()
}

func (v *catalog_[K, V]) AsArray() []AssociationLike[K, V] {
	return v.associations_.AsArray()
}

func (v *catalog_[K, V]) GetIterator() age.IteratorLike[AssociationLike[K, V]] {
	return v.associations_.GetIterator()
}

// Sortable

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

// Stringer

func (v *catalog_[K, V]) String() string {
	return v.GetClass().Notation().FormatValue(v)
}
