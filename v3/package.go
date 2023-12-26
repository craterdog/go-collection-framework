/*******************************************************************************
 *   Copyright (c) 2009-2024 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

/*
This package file defines the INTERFACE to this package.  Any additions to the
types defined in this file require a MINOR version change.  Any deletions from,
or changes to, the types defined in this file require a MAJOR version change.

The package defines a set of simple, pragmatic abstract types and interfaces
for Go based collections of values. It also provide an efficient and compact
implementation of the following collection classes based on these abstractions:
  - Array (extended Go array)
  - Map (extended Go map)
  - List (a sortable list)
  - Catalog (a sortable map)
  - Set (an ordered set)
  - Stack (a LIFO)
  - Queue (a blocking FIFO)

For detailed documentation on this package refer to the wiki:

	https://github.com/craterdog/go-collection-framework/wiki

This package follows the Crater Dog Technologies™ Go Coding Conventions
posted here:

	https://github.com/craterdog/go-coding-conventions/wiki

Additional implementations of these collection classes can be defined and used
seamlessly since the interface definitions only depend on other interfaces and
primitive types; and the class implementations only depend on interfaces, not on
each other.
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
// Go any data type.  This type is used generically to represent components used
// as values.
type Value any

// Function Types

// This function type defines the signature for any function that can determine
// whether or not two specified values are equal to each other.  The meaning of
// "equality" is determined by the specific function that implements this
// signature.
type ComparingFunction func(first Value, second Value) bool

// This function type defines the signature for any function that can determine
// the relative ordering of two values. The result must be one of the following:
//
//	-1: The first value is less than the second value.
//	 0: The first value is equal to the second value.
//	 1: The first value is more than the second value.
//
// The meaning of "less" and "more" is determined by the specific function that
// implements this signature.
type RankingFunction func(first Value, second Value) int

// PACKAGE ABSTRACTIONS

// Abstract Interfaces

// This abstract interface defines the set of method signatures that must be
// supported by all sequences whose values can be accessed using indices. The
// indices of an accessible sequence are ORDINAL rather than ZERO based—which
// never really made sense except for pointer offsets. What is the "zeroth
// value" anyway? It's the "first value", right?  So we start fresh...
//
// This approach allows for positive indices starting at the beginning of the
// sequence, and negative indices starting at the end of the sequence as
// follows:
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
// supported by all sequences of key-value associations.
type Associative[K Key, V Value] interface {
	GetKeys() Sequential[K]
	GetValue(key K) V
	GetValues(keys Sequential[K]) Sequential[V]
	RemoveAll()
	RemoveValue(key K) V
	RemoveValues(keys Sequential[K]) Sequential[V]
	SetValue(key K, value V)
}

// This abstract interface defines the set of method signatures that must be
// supported by all key-value associations.  An association binds a read-only
// key with a setable value.
type Binding[K Key, V Value] interface {
	GetKey() K
	GetValue() V
	SetValue(value V)
}

// This abstract interface defines the set of method signatures that must be
// supported by all sequences that allow new values to be appended, inserted
// and removed.
type Expandable[V Value] interface {
	AppendValue(value V)
	AppendValues(values Sequential[V])
	InsertValue(slot int, value V)
	InsertValues(slot int, values Sequential[V])
	RemoveAll()
	RemoveValue(index int) V
	RemoveValues(first int, last int) Sequential[V]
}

// This abstract interface defines the set of method signatures that must be
// supported by all sequences whose values are accessed using first-in-first-out
// (FIFO) semantics.
type FIFO[V Value] interface {
	AddValue(value V)
	CloseQueue()
	GetCapacity() int
	RemoveHead() (head V, ok bool)
}

// This abstract interface defines the set of method signatures that must be
// supported by all sequences of values that allow new values to be added and
// existing values to be removed.
type Flexible[V Value] interface {
	AddValue(value V)
	AddValues(values Sequential[V])
	GetRanker() RankingFunction
	RemoveAll()
	RemoveValue(value V)
	RemoveValues(values Sequential[V])
}

// This abstract interface defines the set of method signatures that must be
// supported by all sequences whose values are accessed using last-in-first-out
// (LIFO) semantics.
type LIFO[V Value] interface {
	AddValue(value V)
	GetCapacity() int
	GetTop() V
	RemoveAll()
	RemoveTop() V
}

// This abstract interface defines the set of method signatures that must be
// supported by all ratcheted agents that are capable of moving forward and
// backward over the values in a sequence.  It is used to implement the Gang
// of Four (GoF) Iterator Pattern:
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
	GetNext() V
	GetPrevious() V
	GetSlot() int
	HasNext() bool
	HasPrevious() bool
	ToEnd()
	ToSlot(slot int)
	ToStart()
}

// This abstract interface defines the set of method signatures that must be
// supported by all searchable sequences of values.
type Searchable[V Value] interface {
	ContainsAll(values Sequential[V]) bool
	ContainsAny(values Sequential[V]) bool
	ContainsValue(value V) bool
	GetComparer() ComparingFunction
	GetIndex(value V) int
}

// This abstract interface defines the set of method signatures that must be
// supported by all sequences of values.
type Sequential[V Value] interface {
	AsArray() []V
	GetIterator() Ratcheted[V]
	GetSize() int
	IsEmpty() bool
}

// This abstract interface defines the set of method signatures that must be
// supported by all sequences whose values may be reordered using various
// sorting algorithms.
type Sortable[V Value] interface {
	ReverseValues()
	ShuffleValues()
	SortValues()
	SortValuesWithRanker(ranker RankingFunction)
}

// This abstract interface defines the set of method signatures that must be
// supported by all standardized agents that can format collections using a
// standard string representation.
type Standardized interface {
	FormatCollection(collection Collection) string
}

// This abstract interface defines the set of method signatures that must be
// supported by all stringent agents that can parse source text containing
// a standard representation of collections of values.
type Stringent interface {
	ParseCollection(source string) Collection
}

// This abstract interface defines the set of method signatures that must be
// supported by all synchronized groups of threads.
type Synchronized interface {
	Add(delta int)
	Done()
	Wait()
}

// This abstract interface defines the set of method signatures that must be
// supported by all updatable sequences of values.
type Updatable[V Value] interface {
	SetValue(index int, value V)
	SetValues(index int, values Sequential[V])
}

// Abstract Types

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all array-class-like types.
type ArrayClassLike[V Value] interface {
	FromArray(values []V) ArrayLike[V]
	FromSequence(values Sequential[V]) ArrayLike[V]
	FromString(values string) ArrayLike[V]
	WithSize(size int) ArrayLike[V]
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all array-like types.  An array-like type maintains a fixed
// length indexed sequence of values.  Each value is associated with an implicit
// positive integer index. An array-like type uses ORDINAL based indexing rather
// than the more common—and nonsensical—ZERO based indexing scheme (see the
// description of what this means in the Accessible interface definition).
//
// This type is parameterized as follows:
//   - V is any type of value.
//
// This type essentially provides a higher level abstraction for the primitive
// Go array type.
type ArrayLike[V Value] interface {
	Accessible[V]
	Sequential[V]
	Sortable[V]
	Updatable[V]
}

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all association-class-like types.
type AssociationClassLike[K Key, V Value] interface {
	FromPair(key K, value V) AssociationLike[K, V]
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all association-like types.  An association-like type maintains
// information about a key-value association.
//
// This type is parameterized as follows:
//   - K is a primitive type of key.
//   - V is any type of value.
//
// This type is used by catalog-like types to maintain their associations.
type AssociationLike[K Key, V Value] interface {
	Binding[K, V]
}

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all catalog-class-like types.
type CatalogClassLike[K comparable, V Value] interface {
	Empty() CatalogLike[K, V]
	FromArray(associations []Binding[K, V]) CatalogLike[K, V]
	FromMap(associations map[K]V) CatalogLike[K, V]
	FromSequence(associations Sequential[Binding[K, V]]) CatalogLike[K, V]
	FromString(associations string) CatalogLike[K, V]
	Extract(catalog CatalogLike[K, V], keys Sequential[K]) CatalogLike[K, V]
	Merge(first, second CatalogLike[K, V]) CatalogLike[K, V]
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all catalog-like types.  A catalog-like type maintains a
// sequence of key-value associations.
//
// This type is parameterized as follows:
//   - K is a primitive type of key.
//   - V is any type of entity.
//
// A catalog-like type can use any association-like type key-value association.
type CatalogLike[K Key, V Value] interface {
	Associative[K, V]
	Sequential[Binding[K, V]]
	Sortable[Binding[K, V]]
}

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all iterator-class-like types.
type IteratorClassLike[V Value] interface {
	FromSequence(sequence Sequential[V]) IteratorLike[V]
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all iterator-like types.  An iterator-like type can be used to
// traverse a sequence of values in either direction (forward or backward). It
// uses a ratchet based system that locks into the SLOTS between the values in
// the sequence. At each slot an iterator has access to the previous value and
// next value in the sequence (assuming they exist). The slot at the start of
// the sequence has no PREVIOUS value, and the slot at the end of the sequence
// has no NEXT value.
//
// This type is parameterized as follows:
//   - V is any type of value.
//
// An iterator-like type is supported by all collection types.
type IteratorLike[V Value] interface {
	Ratcheted[V]
}

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all list-class-like types.
type ListClassLike[V Value] interface {
	Empty() ListLike[V]
	FromArray(values []V) ListLike[V]
	FromSequence(values Sequential[V]) ListLike[V]
	FromString(values string) ListLike[V]
	WithComparer(comparer ComparingFunction) ListLike[V]
	Concatenate(first, second ListLike[V]) ListLike[V]
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all list-like types.  A list-like type maintains a dynamic
// sequence of values which can grow or shrink as needed.  Each value is
// associated with an implicit positive integer index. An array-like type
// uses ORDINAL based indexing rather than the more common—and
// nonsensical—ZERO based indexing scheme (see the description of what this
// means in the Accessible interface definition).
//
// This type is parameterized as follows:
//   - V is any type of value.
//
// A comparison of the values in the sequence is done using a configurable
// ComparerFunction.
type ListLike[V Value] interface {
	Accessible[V]
	Expandable[V]
	Searchable[V]
	Sequential[V]
	Sortable[V]
	Updatable[V]
}

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all map-class-like types.
type MapClassLike[K comparable, V Value] interface {
	Empty() MapLike[K, V]
	FromArray(associations []Binding[K, V]) MapLike[K, V]
	FromMap(associations map[K]V) MapLike[K, V]
	FromSequence(associations Sequential[Binding[K, V]]) MapLike[K, V]
	FromString(associations string) MapLike[K, V]
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all map-like types.  A map-like type extends the primitive Go
// map type and maintains a sequence of key-value associations.  The order of
// the key-value associations in a primitive Go map is random, even for two Go
// maps containing the same key-value associations.
//
// This type is parameterized as follows:
//   - K is a primitive type of key.
//   - V is any type of entity.
//
// A map-like type can use any association-like type key-value association.
type MapLike[K Key, V Value] interface {
	Associative[K, V]
	Sequential[Binding[K, V]]
}

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all notation-class-like types.
type NotationClassLike interface {
	Default() NotationLike
	GetDefaultDepth() int
	WithDepth(depth int) NotationLike
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all notation-like types.  A notation-like type can be used to
// parse and format collections using a canonical notation like XML, JSON and
// CDCN (Crater Dog Collection Notation™).
type NotationLike interface {
	Standardized
	Stringent
}

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all queue-class-like types.
type QueueClassLike[V Value] interface {
	Empty() QueueLike[V]
	FromArray(values []V) QueueLike[V]
	FromSequence(values Sequential[V]) QueueLike[V]
	FromString(values string) QueueLike[V]
	GetDefaultCapacity() int
	WithCapacity(capacity int) QueueLike[V]
	Fork(group Synchronized, input QueueLike[V], size int) Sequential[QueueLike[V]]
	Join(group Synchronized, inputs Sequential[QueueLike[V]]) QueueLike[V]
	Split(group Synchronized, input QueueLike[V], size int) Sequential[QueueLike[V]]
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all queue-like types.  A queue-like type implements FIFO
// (i.e. first-in-first-out) semantics.
//
// This type is parameterized as follows:
//   - V is any type of value.
//
// A queue-like type is generally used by multiple go-routines at the same
// time and therefore enforces synchronized access.  A queue-like type enforces
// a maximum length and will block on attempts to add a value it is full.  It
// will also block on attempts to remove a value when it is empty.
type QueueLike[V Value] interface {
	FIFO[V]
	Sequential[V]
}

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all set-class-like types.
type SetClassLike[V Value] interface {
	Empty() SetLike[V]
	FromArray(values []V) SetLike[V]
	FromSequence(values Sequential[V]) SetLike[V]
	FromSequenceWithRanker(values Sequential[V], ranker RankingFunction) SetLike[V]
	FromString(values string) SetLike[V]
	WithRanker(ranker RankingFunction) SetLike[V]
	And(first, second SetLike[V]) SetLike[V]
	Or(first, second SetLike[V]) SetLike[V]
	Sans(first, second SetLike[V]) SetLike[V]
	Xor(first, second SetLike[V]) SetLike[V]
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all set-like types.  A set-like type maintains an ordered
// sequence of values which can grow or shrink as needed.
//
// This type is parameterized as follows:
//   - V is any type of value.
//
// The order of the values is determined by a configurable RankingFunction.
type SetLike[V Value] interface {
	Accessible[V]
	Flexible[V]
	Searchable[V]
	Sequential[V]
}

// This abstract type defines the set of class constants, constructors and
// functions that must be supported by all stack-class-like types.
type StackClassLike[V Value] interface {
	Empty() StackLike[V]
	FromArray(values []V) StackLike[V]
	FromSequence(values Sequential[V]) StackLike[V]
	FromString(values string) StackLike[V]
	GetDefaultCapacity() int
	WithCapacity(capacity int) StackLike[V]
}

// This abstract type defines the set of abstract interfaces that must be
// supported by all stack-like types.  A stack-like type implements LIFO
// (i.e. last-in-first-out) semantics.
//
// This type is parameterized as follows:
//   - V is any type of value.
//
// A stack-like type enforces a maximum depth and will panic if that depth is
// exceeded.  It will also panic on attempts to remove a value when it is empty.
type StackLike[V Value] interface {
	LIFO[V]
	Sequential[V]
}
