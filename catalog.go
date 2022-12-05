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

// ASSOCIATION IMPLEMENTATION

// This constructor creates a new association with the specified key and value.
func Association[K Key, V Value](key K, value V) AssociationLike[K, V] {
	return &association[K, V]{key, value}
}

// This type defines the structure and methods associated with a key-value
// pair. This type is parameterized as follows:
//   - K is a primitive type of key.
//   - V is any type of entity.
//
// This structure is used by the catalog type to maintain its associations.
type association[K Key, V Value] struct {
	key   K
	value V
}

// STRINGER INTERFACE

// This method returns a canonical string for this association.
func (v *association[K, V]) String() string {
	return FormatValue(v)
}

// BINDING INTERFACE

// This method returns the key for this association.
func (v *association[K, V]) GetKey() K {
	return v.key
}

// This method returns the value for this association.
func (v *association[K, V]) GetValue() V {
	return v.value
}

// This method sets the value of this association to a new value.
func (v *association[K, V]) SetValue(value V) {
	v.value = value
}

// CATALOG IMPLEMENTATION

// This constructor creates a new empty catalog.
func Catalog[K Key, V Value]() CatalogLike[K, V] {
	var keys = map[Key]Binding[K, V]{}
	var associations = List[Binding[K, V]]()
	return &catalog[K, V]{associations, associations, keys}
}

// This constructor creates a new catalog from the specified array of
// associations.
func CatalogFromArray[K Key, V Value](array []Binding[K, V]) CatalogLike[K, V] {
	var v = Catalog[K, V]()
	for _, association := range array {
		v.AddAssociation(association)
	}
	return v
}

// This constructor creates a new catalog from the specified sequence of
// associations.
func CatalogFromSequence[K Key, V Value](sequence Sequential[Binding[K, V]]) CatalogLike[K, V] {
	var v = Catalog[K, V]()
	var iterator = Iterator(sequence)
	for iterator.HasNext() {
		var association = iterator.GetNext()
		v.AddAssociation(association)
	}
	return v
}

// This function returns a new catalog containing all of the associations
// that are in the specified catalogs in the order that they appear in each
// catalog.
func Merge[K Key, V Value](first, second CatalogLike[K, V]) CatalogLike[K, V] {
	var result = Catalog[K, V]()
	result.AddAssociations(first)
	result.AddAssociations(second)
	return result
}

// This function returns a new catalog containing only the associations
// that are in the specified catalog that have the specified keys. The
// associations in the resulting catalog will be in the same order as the
// specified keys.
func Extract[K Key, V Value](catalog CatalogLike[K, V], keys Sequential[K]) CatalogLike[K, V] {
	var result = Catalog[K, V]()
	var iterator = Iterator(keys)
	for iterator.HasNext() {
		var key = iterator.GetNext()
		var value = catalog.GetValue(key)
		result.SetValue(key, value)
	}
	return result
}

// This type defines the structure and methods associated with a catalog of
// key-value pair associations. This type is parameterized as follows:
//   - K is a primitive type of key.
//   - V is any type of entity.
type catalog[K Key, V Value] struct {
	// Note: The delegated methods don't see the real collection type.
	Sequential[Binding[K, V]]
	associations ListLike[Binding[K, V]]
	keys         map[Key]Binding[K, V]
}

// STRINGER INTERFACE

func (v *catalog[K, V]) String() string {
	return FormatValue(v)
}

// ASSOCIATIVE INTERFACE

// This method appends the specified association to the end of this catalog.
func (v *catalog[K, V]) AddAssociation(association Binding[K, V]) {
	var key = association.GetKey()
	var value = association.GetValue()
	v.SetValue(key, value) // This copies the association.
}

// This method appends the specified associations to the end of this catalog.
func (v *catalog[K, V]) AddAssociations(associations Sequential[Binding[K, V]]) {
	var iterator = Iterator(associations)
	for iterator.HasNext() {
		var association = iterator.GetNext()
		v.AddAssociation(association)
	}
}

// This method returns the keys for this catalog.
func (v *catalog[K, V]) GetKeys() Sequential[K] {
	var keys = List[K]()
	var iterator = Iterator[Binding[K, V]](v.associations)
	for iterator.HasNext() {
		var association = iterator.GetNext()
		keys.AddValue(association.GetKey())
	}
	return keys
}

// This method returns the values associated with the specified keys for this
// catalog. The values are returned in the same order as the keys in the
// catalog.
func (v *catalog[K, V]) GetValues(keys Sequential[K]) Sequential[V] {
	var values = List[V]()
	var iterator = Iterator(keys)
	for iterator.HasNext() {
		var key = iterator.GetNext()
		values.AddValue(v.GetValue(key))
	}
	return values
}

// This method returns the value that is associated with the specified key in
// this catalog.
func (v *catalog[K, V]) GetValue(key K) V {
	var value V // Set the return value to its zero value.
	var association, exists = v.keys[key]
	if exists {
		// Extract the value.
		value = association.GetValue()
	}
	return value
}

// This method sets the value associated with the specified key to the
// specified value.
func (v *catalog[K, V]) SetValue(key K, value V) {
	var association, exists = v.keys[key]
	if exists {
		// Set the value of an existing association.
		association.SetValue(value)
	} else {
		// Add a new association.
		association = Association[K, V](key, value)
		v.associations.AddValue(association)
		v.keys[key] = association
	}
}

// This method removes the association associated with the specified key from the
// catalog and returns it.
func (v *catalog[K, V]) RemoveValue(key K) V {
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

// This method removes the associations associated with the specified keys from
// the catalog and returns the removed values.
func (v *catalog[K, V]) RemoveValues(keys Sequential[K]) Sequential[V] {
	var values = List[V]()
	var iterator = Iterator(keys)
	for iterator.HasNext() {
		var key = iterator.GetNext()
		values.AddValue(v.RemoveValue(key))
	}
	return values
}

// This method removes all associations from this catalog.
func (v *catalog[K, V]) RemoveAll() {
	v.keys = map[Key]Binding[K, V]{}
	v.associations.RemoveAll()
}

// SORTABLE INTERFACE

// This method sorts this catalog using the canonical rank function to compare
// the keys.
func (v *catalog[K, V]) SortValues() {
	v.associations.SortValues()
}

// This method sorts this catalog using the specified rank function to compare
// the keys.
func (v *catalog[K, V]) SortValuesWithRanker(rank RankingFunction) {
	v.associations.SortValuesWithRanker(rank)
}

// This method reverses the order of all associations in this catalog.
func (v *catalog[K, V]) ReverseValues() {
	v.associations.ReverseValues()
}

// This method pseudo-randomly shuffles the order of all associations in this
// catalog.
func (v *catalog[K, V]) ShuffleValues() {
	v.associations.ShuffleValues()
}
