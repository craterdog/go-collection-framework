/*******************************************************************************
 *   Copyright (c) 2009-2022 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

/*
This collections package defines a set of simple, pragmatic interfaces for
collections of sequential values. It also provides efficient and compact
implementations of the following collection types based on those interfaces:
  - array (native Go array)
  - catalog (a sortable map)
  - list (a sortable list)
  - map (native Go map)
  - queue (a blocking FIFO)
  - set (an ordered set)
  - stack (LIFO)

Additional implementations of these collection types can be defined and used
seemlessly since the interface definitions only depend on other interfaces and
native types; and the type implementations only depend on interfaces, not on
each other. For a full description of this package see the wiki documentation
at the github repository maintaining this package:

	https://github.com/craterdog/go-collection-framework/wiki
*/
package collections

// TYPE DEFINITIONS

type (
	Key        any
	Value      any
	Primitive  any
	Collection any
)

// INDIVIDUAL INTERFACES

// This interface defines the methods supported by all sequences of values.
type Sequential[V Value] interface {
	IsEmpty() bool
	GetSize() int
	AsArray() []V
}

// This interface defines the methods supported by all sequences whose values can
// be accessed using indices. The indices of an accessible sequence are ORDINAL
// rather than ZERO based (which is "SO last century"). This allows for positive
// indices starting at the beginning of the sequence, and negative indices
// starting at the end of the sequence as follows:
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

// This interface defines the methods supported by all updatable sequences of
// values.
type Updatable[V Value] interface {
	SetValue(index int, value V)
	SetValues(index int, values Sequential[V])
}

// This interface defines the methods supported by all searchable sequences of
// values.
type Searchable[V Value] interface {
	GetIndex(value V) int
	ContainsValue(value V) bool
	ContainsAny(values Sequential[V]) bool
	ContainsAll(values Sequential[V]) bool
}

// This interface defines the methods supported by all sequences of values that
// allow values to be added and removed.
type Flexible[V Value] interface {
	AddValue(value V)
	AddValues(values Sequential[V])
	RemoveValue(value V)
	RemoveValues(values Sequential[V])
	RemoveAll()
}

// This interface defines the methods supported by all sequences whose values may
// be modified, inserted, removed, or reordered.
type Malleable[V Value] interface {
	AddValue(value V)
	AddValues(values Sequential[V])
	InsertValue(slot int, value V)
	InsertValues(slot int, values Sequential[V])
	RemoveValue(index int) V
	RemoveValues(first int, last int) Sequential[V]
	RemoveAll()
}

// This interface defines the methods supported by all sequences whose values may
// be sorted using various sorting algorithms.
type Sortable[V Value] interface {
	SortValues()
	SortValuesWithRanker(rank RankingFunction)
	ReverseValues()
	ShuffleValues()
}

// This interface defines the methods supported by all binding associations.
// It binds a readonly key with a setable value.
type Binding[K Key, V Value] interface {
	GetKey() K
	GetValue() V
	SetValue(value V)
}

// This interface defines the methods supported by all associative sequences
// whose values consist of key-value pair associations.
type Associative[K Key, V Value] interface {
	GetKeys() Sequential[K]
	GetValues(keys Sequential[K]) Sequential[V]
	GetValue(key K) V
	SetValue(key K, value V)
	RemoveValue(key K) V
	RemoveValues(keys Sequential[K]) Sequential[V]
	RemoveAll()
}

// This interface defines the methods supported by all sequences whose values
// are accessed using first-in-first-out (FIFO) semantics.
type FIFO[V Value] interface {
	GetCapacity() int
	AddValue(value V)
	RemoveHead() (head V, ok bool)
	CloseQueue()
}

// This interface defines the methods supported by all sequences whose values
// are accessed using last-in-first-out (LIFO) semantics.
type LIFO[V Value] interface {
	GetCapacity() int
	AddValue(value V)
	GetTop() V
	RemoveTop() V
	RemoveAll()
}

// This interface defines the methods supported by all ratcheted agents that
// are capable of moving forward and backward over the values in a sequence. It
// is used to implement the GoF Iterator Pattern:
//   - https://en.wikipedia.org/wiki/Iterator_pattern
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

// This interface defines the methods supported by all canonical agents that can
// format any value in a standard way.
type Canonical interface {
	GetIndentation() int
	FormatValue(value Value)
	AppendString(s string)
	AppendNewline()
	GetResult() string
}

// This interface defines the methods supported by all discerning agent types
// that can compare and rank two values.
type Discerning interface {
	CompareValues(first Value, second Value) bool
	RankValues(first Value, second Value) int
}

// This type defines the function signature for any function that can determine
// whether or not two specified values are equal to each other.
type ComparisonFunction func(first Value, second Value) bool

// This type defines the function signature for any function that can determine
// the relative ordering of two specified values. The result must be one of
// the following:
//   - -1: The first value is less than the second value.
//   - 0: The first value is equal to the second value.
//   - 1: The first value is more than the second value.
type RankingFunction func(first Value, second Value) int

// This interface defines the methods supported by all systematic agents that
// can shuffle or sort an array of values using a ranking function.
type Systematic[V Value] interface {
	ReverseArray(array []V)
	ShuffleArray(array []V)
	SortArray(array []V)
}

// This type defines the function signature for any function that can sort an
// array of values using a ranking function.
type Sort[V Value] func(array []V, rank RankingFunction)

// CONSOLIDATED INTERFACES

// This interface consolidates all the interfaces supported by native array-like
// sequences.
type ArrayLike[V Value] interface {
	Sequential[V]
	Accessible[V]
	Updatable[V]
}

// This interface consolidates all the interfaces supported by native map-like
// sequences. Note, that the order of the key-value pairs on a native map is
// random, even for two maps containing the same keys.
type MapLike[K Key, V Value] interface {
	Sequential[Binding[K, V]]
	Associative[K, V]
}

// This interface defines the methods supported by all association-like types.
// An association binds a key with a value.
type AssociationLike[K Key, V Value] interface {
	Binding[K, V]
}

// This interface consolidates all the interfaces supported by catalog-like
// sequences.
type CatalogLike[K Key, V Value] interface {
	Sequential[Binding[K, V]]
	Associative[K, V]
	Sortable[Binding[K, V]]
}

// This interface consolidates all the interfaces supported by list-like
// sequences.
type ListLike[V Value] interface {
	Sequential[V]
	Accessible[V]
	Updatable[V]
	Searchable[V]
	Malleable[V]
	Sortable[V]
}

// This interface consolidates all of the interfaces supported by queue-like
// sequences.
type QueueLike[V Value] interface {
	Sequential[V]
	FIFO[V]
}

// This interface consolidates all the interfaces supported by set-like
// sequences.
type SetLike[V Value] interface {
	Sequential[V]
	Accessible[V]
	Searchable[V]
	Flexible[V]
}

// This interface consolidates all the interfaces supported by stack-like
// sequences.
type StackLike[V Value] interface {
	Sequential[V]
	LIFO[V]
}

// This interface defines the methods supported by all iterator-like types.
type IteratorLike[V Value] interface {
	Ratcheted[V]
}

// This interface defines the methods supported by all sorter-like types.
type FormatterLike interface {
	Canonical
}

// This interface defines the methods supported by all collator-like types.
type CollatorLike interface {
	Discerning
}

// This interface defines the methods supported by all sorter-like types.
type SorterLike[V Value] interface {
	Systematic[V]
}
