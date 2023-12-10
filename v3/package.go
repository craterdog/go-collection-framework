/*******************************************************************************
 *   Copyright (c) 2009-2023 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

/*
This package defines a set of simple, pragmatic interfaces for collections of
sequential values. It also provides efficient and compact implementations of the
following collection classes based on those interfaces:
  - Array (native Go array)
  - Catalog (a sortable map)
  - List (a sortable list)
  - Map (native Go map)
  - Queue (a blocking FIFO)
  - Set (an ordered set)
  - Stack (LIFO)

Additional implementations of these collection classes can be defined and used
seamlessly since the interface definitions only depend on other interfaces and
native types; and the class implementations only depend on interfaces, not on
each other.

For detailed documentation on this package refer to the wiki:

	https://github.com/craterdog/go-collection-framework/wiki
*/
package collections

// PACKAGE TYPES

// Specialized Types

// This specialized type definition represents a specialization of the primitive
// Go any data type.  This type is used generically to represent collections of
// any type.
type Collection any

// This specialized type definition represents a specialization of the primitive
// Go any data type.  This type is used generically to represent elements used
// as keys.
type Key any

// This specialized type definition represents a specialization of the primitive
// Go any data type.  This type is used generically to represent primitive
// components.
type Primitive any

// This specialized type definition represents a specialization of the primitive
// Go any data type.  This type is used generically to represent components used
// as values.
type Value any

// Function Types

// This function type defines the signature for any function that can determine
// whether or not two specified values are equal to each other.
type ComparingFunction func(first Value, second Value) bool

// This function type defines the signature for any function that can determine
// the relative ordering of two specified values. The result must be one of
// the following:
//   - -1: The first value is less than the second value.
//   - 0: The first value is equal to the second value.
//   - 1: The first value is more than the second value.
type RankingFunction func(first Value, second Value) int

// This function type defines the signature for any function that can sort an
// array of values using a ranking function.
type SortingFunction[V Value] func(array []V, ranker RankingFunction)

// PACKAGE ABSTRACTIONS

// Abstract Interfaces

// This abstract interface defines the set of method signatures that must be
// supported by all sequences whose values can be accessed using indices. The
// indices of an accessible sequence are ORDINAL rather than ZERO based (which
// is "SO last century"). This allows for positive indices starting at the
// beginning of the sequence, and negative indices starting at the end of the
// sequence as follows:
//
//	    1           2           3             N
//	[value 1] . [value 2] . [value 3] ... [value N]
//	   -N        -(N-1)      -(N-2)          -1
//
// Notice that because the indices are ordinal based, the positive and negative
// indices are symmetrical.
type Accessible[V Value] interface {
	GetValue(index int) V
	GetValues(first int, last int) Sequential[V]
}

// This abstract interface defines the set of method signatures that must be
// supported by all sequences of key-value pair associations.
type Associative[K Key, V Value] interface {
	GetKeys() Sequential[K]
	GetValues(keys Sequential[K]) Sequential[V]
	GetValue(key K) V
	SetValue(key K, value V)
	RemoveValue(key K) V
	RemoveValues(keys Sequential[K]) Sequential[V]
	RemoveAll()
}

// This abstract interface defines the set of method signatures that must be
// supported by all binding associations.  It binds a read-only key with a
// setable value.
type Binding[K Key, V Value] interface {
	GetKey() K
	GetValue() V
	SetValue(value V)
}

// This abstract interface defines the set of method signatures that must be
// supported by all discerning agents that can compare and rank two values.
type Discerning interface {
	CompareValues(first Value, second Value) bool
	RankValues(first Value, second Value) int
}

// This abstract interface defines the set of method signatures that must be
// supported by all sequences that allow new values to be appended, inserted
// and removed.
type Expandable[V Value] interface {
	AppendValue(value V)
	AppendValues(values Sequential[V])
	InsertValue(slot int, value V)
	InsertValues(slot int, values Sequential[V])
	RemoveValue(index int) V
	RemoveValues(first int, last int) Sequential[V]
	RemoveAll()
}

// This abstract interface defines the set of method signatures that must be
// supported by all sequences whose values are accessed using first-in-first-out
// (FIFO) semantics.
type FIFO[V Value] interface {
	GetCapacity() int
	AddValue(value V)
	RemoveHead() (head V, ok bool)
	CloseQueue()
}

// This abstract interface defines the set of method signatures that must be
// supported by all sequences of values that allow new values to be added and
// existing values to be removed.
type Flexible[V Value] interface {
	GetRanker() RankingFunction
	AddValue(value V)
	AddValues(values Sequential[V])
	RemoveValue(value V)
	RemoveValues(values Sequential[V])
	RemoveAll()
}

// This abstract interface defines the set of method signatures that must be
// supported by all sequences whose values are accessed using last-in-first-out
// (LIFO) semantics.
type LIFO[V Value] interface {
	GetCapacity() int
	AddValue(value V)
	GetTop() V
	RemoveTop() V
	RemoveAll()
}

// This abstract interface defines the set of method signatures that must be
// supported by all ratcheted agents that are capable of moving forward and
// backward over the values in a sequence.  It is used to implement the Gang
// of Four Iterator Pattern:
//
//	https://en.wikipedia.org/wiki/Iterator_pattern
//
// A ratcheted agent locks into the slots that reside between each value in the
// sequence:
//
//	    [value 1] . [value 2] . [value 3] ... [value N]
//	  ^           ^           ^                         ^
//	slot 0      slot 1      slot 2                    slot N
//
// It moves from slot to slot and has access to the values (if they exist) on
// each side of the slot.
type Ratcheted[V Value] interface {
	GetSlot() int
	ToSlot(slot int)
	ToStart()
	ToEnd()
	HasPrevious() bool
	GetPrevious() V
	HasNext() bool
	GetNext() V
}

// This abstract interface defines the set of method signatures that must be
// supported by all searchable sequences of values.
type Searchable[V Value] interface {
	GetComparer() ComparingFunction
	GetIndex(value V) int
	ContainsValue(value V) bool
	ContainsAny(values Sequential[V]) bool
	ContainsAll(values Sequential[V]) bool
}

// This abstract interface defines the set of method signatures that must be
// supported by all sequences of values.
type Sequential[V Value] interface {
	IsEmpty() bool
	GetSize() int
	AsArray() []V
	GetIterator() Ratcheted[V]
}

// This abstract interface defines the set of method signatures that must be
// supported by all sequences whose values may be reordered using various
// sorting algorithms.
type Sortable[V Value] interface {
	SortValues()
	SortValuesWithRanker(ranker RankingFunction)
	ReverseValues()
	ShuffleValues()
}

// This abstract interface defines the set of method signatures that must be
// supported by all systematic agents that can shuffle or sort an array of
// values.
type Systematic[V Value] interface {
	ReverseValues(array []V)
	ShuffleValues(array []V)
	SortValues(array []V)
}

// This abstract interface defines the set of method signatures that must be
// supported by all updatable sequences of values.
type Updatable[V Value] interface {
	SetValue(index int, value V)
	SetValues(index int, values Sequential[V])
}

// Abstract Types

// This abstract type defines the set of abstract interfaces that must be
// supported by all array-like types.
type ArrayLike[V Value] interface {
	Sequential[V]
	Accessible[V]
	Updatable[V]
	Sortable[V]
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all association-like types.
type AssociationLike[K Key, V Value] interface {
	Binding[K, V]
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all catalog-like types.
type CatalogLike[K Key, V Value] interface {
	Sequential[Binding[K, V]]
	Associative[K, V]
	Sortable[Binding[K, V]]
}

// This interface defines the methods supported by all collator-like types.
type CollatorLike interface {
	Discerning
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all iterator-like types.
type IteratorLike[V Value] interface {
	Ratcheted[V]
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all list-like types.
type ListLike[V Value] interface {
	Sequential[V]
	Accessible[V]
	Updatable[V]
	Searchable[V]
	Expandable[V]
	Sortable[V]
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all map-like types.  NOTE: The order of the key-value pairs in a
// native map is random, even for two maps containing the same keys.
type MapLike[K Key, V Value] interface {
	Sequential[Binding[K, V]]
	Associative[K, V]
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all queue-like types.
type QueueLike[V Value] interface {
	Sequential[V]
	FIFO[V]
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all set-like types.
type SetLike[V Value] interface {
	Sequential[V]
	Accessible[V]
	Searchable[V]
	Flexible[V]
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all sorter-like types.
type SorterLike[V Value] interface {
	Systematic[V]
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all stack-like types.
type StackLike[V Value] interface {
	Sequential[V]
	LIFO[V]
}
