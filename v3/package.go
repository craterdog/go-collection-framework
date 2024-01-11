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

The package defines a set of simple, pragmatic abstract classes and interfaces
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

// Functional Types

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
// supported by all canonical notations.
type Canonical interface {
	FormatCollection(collection Collection) string
	ParseCollection(collection string) Collection
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
// supported by all sequences of values that allow new values to be added and
// existing values to be removed.
type Flexible[V Value] interface {
	AddValue(value V)
	AddValues(values Sequential[V])
	RemoveAll()
	RemoveValue(value V)
	RemoveValues(values Sequential[V])
}

// This abstract interface defines the set of method signatures that must be
// supported by all sequences of values that allow new values to be added and
// limit the total number of values in the sequence.
type Limited[V Value] interface {
	AddValue(value V)
	GetCapacity() int
	RemoveAll()
}

// This abstract interface defines the set of method signatures that must be
// supported by all searchable sequences of values.
type Searchable[V Value] interface {
	ContainsAll(values Sequential[V]) bool
	ContainsAny(values Sequential[V]) bool
	ContainsValue(value V) bool
	GetIndex(value V) int
}

// This abstract interface defines the set of method signatures that must be
// supported by all sequences of values.
type Sequential[V Value] interface {
	AsArray() []V
	GetIterator() IteratorLike[V]
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

// Abstract Classes

// This abstract class defines the set of abstract interfaces that must be
// supported by all array-like classes.  An array-like type maintains a fixed
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
	// Supported Interfaces
	Accessible[V]
	Sequential[V]
	Sortable[V]
	Updatable[V]
}

// This abstract class defines the set of abstract interfaces that must be
// supported by all association-like classes.  An association-like type maintains
// information about a key-value association.
//
// This type is parameterized as follows:
//   - K is a primitive type of key.
//   - V is any type of value.
//
// This type is used by catalog-like classes to maintain their associations.
type AssociationLike[K Key, V Value] interface {
	// Public Interface
	GetKey() K
	GetValue() V
	SetValue(value V)
}

// This abstract class defines the set of abstract interfaces that must be
// supported by all catalog-like classes.  A catalog-like type maintains a
// sequence of key-value associations.
//
// This type is parameterized as follows:
//   - K is a primitive type of key.
//   - V is any type of entity.
//
// A catalog-like type can use any association-like type key-value association.
type CatalogLike[K Key, V Value] interface {
	// Supported Interfaces
	Associative[K, V]
	Sequential[AssociationLike[K, V]]
	Sortable[AssociationLike[K, V]]
}

// This abstract class defines the set of abstract interfaces that must be
// supported by all collator-like classes.  A collator-like type is capable of
// comparing and ranking two values the any type.
type CollatorLike interface {
	// Public Interface
	CompareValues(first Value, second Value) bool
	RankValues(first Value, second Value) int
}

// This abstract class defines the set of abstract interfaces that must be
// supported by all iterator-like classes.  An iterator-like type can be used to
// move forward and backward over the values in a sequence.  It implements the
// Gang of Four (GoF) Iterator Design Pattern:
//
//	https://en.wikipedia.org/wiki/Iterator_pattern
//
// A iterator agent locks into the slots that reside between each value in the
// sequence:
//
//	    [value 1] . [value 2] . [value 3] ... [value N]
//	  ^           ^           ^                         ^
//	slot 0      slot 1      slot 2                    slot N
//
// It moves from slot to slot and has access to the values (if they exist) on
// each side of the slot.  At each slot an iterator has access to the previous
// value and next value in the sequence (assuming they exist). The slot at the
// start of the sequence has no PREVIOUS value, and the slot at the end of the
// sequence has no NEXT value.
//
// This type is parameterized as follows:
//   - V is any type of value.
//
// An iterator-like type is supported by all collection types.
type IteratorLike[V Value] interface {
	// Public Interface
	GetNext() V
	GetPrevious() V
	GetSlot() int
	HasNext() bool
	HasPrevious() bool
	ToEnd()
	ToSlot(slot int)
	ToStart()
}

// This abstract class defines the set of abstract interfaces that must be
// supported by all list-like classes.  A list-like type maintains a dynamic
// sequence of values which can grow or shrink as needed.  Each value is
// associated with an implicit positive integer index. An array-like type
// uses ORDINAL based indexing rather than the more common—and
// nonsensical—ZERO based indexing scheme (see the description of what this
// means in the Accessible interface definition).
//
// This type is parameterized as follows:
//   - V is any type of value.
//
// All comparison and ranking of values in the sequence is done using the
// default collator.
type ListLike[V Value] interface {
	// Supported Interfaces
	Accessible[V]
	Expandable[V]
	Searchable[V]
	Sequential[V]
	Sortable[V]
	Updatable[V]
}

// This abstract class defines the set of abstract interfaces that must be
// supported by all map-like classes.  A map-like type extends the primitive Go
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
	// Supported Interfaces
	Associative[K, V]
	Sequential[AssociationLike[K, V]]
}

// This abstract class defines the set of abstract interfaces that must be
// supported by all notation-like classes.  A notation-like type can be used to
// parse and format collections using a canonical notation like XML, JSON and
// CDCN (Crater Dog Collection Notation™).
type NotationLike interface {
	// Supported Interfaces
	Canonical
}

// This abstract class defines the set of abstract interfaces that must be
// supported by all queue-like classes.  A queue-like type implements FIFO
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
	// Supported Interfaces
	Sequential[V]
	Limited[V]

	// Public Interface
	CloseQueue()
	RemoveHead() (head V, ok bool)
}

// This abstract class defines the set of abstract interfaces that must be
// supported by all set-like classes.  A set-like type maintains an ordered
// sequence of values which can grow or shrink as needed.
//
// This type is parameterized as follows:
//   - V is any type of value.
//
// The order of the values is determined by a configurable RankingFunction.
type SetLike[V Value] interface {
	// Supported Interfaces
	Accessible[V]
	Flexible[V]
	Searchable[V]
	Sequential[V]

	// Public Interface
	GetCollator() CollatorLike
}

// This abstract class defines the set of abstract interfaces that must be
// supported by all stack-like classes.  A stack-like type implements LIFO
// (i.e. last-in-first-out) semantics.
//
// This type is parameterized as follows:
//   - V is any type of value.
//
// A stack-like type enforces a maximum depth and will panic if that depth is
// exceeded.  It will also panic on attempts to remove a value when it is empty.
type StackLike[V Value] interface {
	// Supported Interfaces
	Sequential[V]
	Limited[V]

	// Public Interface
	RemoveTop() V
}

// Abstract Namespaces

// This abstract namespace defines the set of class constants, constructors and
// functions that must be supported by all array-class-like namespaces.
type ArrayClassLike[V Value] interface {
	// Constructors
	FromArray(values []V) ArrayLike[V]
	FromSequence(values Sequential[V]) ArrayLike[V]
	FromString(values string) ArrayLike[V]
	WithSize(size int) ArrayLike[V]
}

// This abstract namespace defines the set of class constants, constructors and
// functions that must be supported by all association-class-like namespaces.
type AssociationClassLike[K Key, V Value] interface {
	// Constructors
	FromPair(key K, value V) AssociationLike[K, V]
}

// This abstract namespace defines the set of class constants, constructors and
// functions that must be supported by all catalog-class-like namespaces.
type CatalogClassLike[K comparable, V Value] interface {
	// Constructors
	Empty() CatalogLike[K, V]
	FromArray(associations []AssociationLike[K, V]) CatalogLike[K, V]
	FromMap(associations map[K]V) CatalogLike[K, V]
	FromSequence(associations Sequential[AssociationLike[K, V]]) CatalogLike[K, V]
	FromString(associations string) CatalogLike[K, V]

	// Functions
	Extract(catalog CatalogLike[K, V], keys Sequential[K]) CatalogLike[K, V]
	Merge(first, second CatalogLike[K, V]) CatalogLike[K, V]
}

// This abstract namespace defines the set of class constants, constructors and
// functions that must be supported by all collator-class-like namespaces.
type CollatorClassLike interface {
	// Constants
	GetDefaultDepth() int

	// Constructors
	Default() CollatorLike
	WithDepth(depth int) CollatorLike
}

// This abstract namespace defines the set of class constants, constructors and
// functions that must be supported by all iterator-class-like namespaces.
type IteratorClassLike[V Value] interface {
	// Constructors
	FromSequence(sequence Sequential[V]) IteratorLike[V]
}

// This abstract namespace defines the set of class constants, constructors and
// functions that must be supported by all list-class-like namespaces.
type ListClassLike[V Value] interface {
	// Constructors
	Empty() ListLike[V]
	FromArray(values []V) ListLike[V]
	FromSequence(values Sequential[V]) ListLike[V]
	FromString(values string) ListLike[V]

	// Functions
	Concatenate(first, second ListLike[V]) ListLike[V]
}

// This abstract namespace defines the set of class constants, constructors and
// functions that must be supported by all map-class-like namespaces.
type MapClassLike[K comparable, V Value] interface {
	// Constructors
	Empty() MapLike[K, V]
	FromArray(associations []AssociationLike[K, V]) MapLike[K, V]
	FromMap(associations map[K]V) MapLike[K, V]
	FromSequence(associations Sequential[AssociationLike[K, V]]) MapLike[K, V]
	FromString(associations string) MapLike[K, V]
}

// This abstract namespace defines the set of class constants, constructors and
// functions that must be supported by all notation-class-like namespaces.
type NotationClassLike interface {
	// Constants
	GetDefaultDepth() int

	// Constructors
	Default() NotationLike
	WithDepth(depth int) NotationLike
}

// This abstract namespace defines the set of class constants, constructors and
// functions that must be supported by all queue-class-like namespaces.
type QueueClassLike[V Value] interface {
	// Constants
	GetDefaultCapacity() int

	// Constructors
	Empty() QueueLike[V]
	FromArray(values []V) QueueLike[V]
	FromSequence(values Sequential[V]) QueueLike[V]
	FromString(values string) QueueLike[V]
	WithCapacity(capacity int) QueueLike[V]

	// Functions
	Fork(group Synchronized, input QueueLike[V], size int) Sequential[QueueLike[V]]
	Join(group Synchronized, inputs Sequential[QueueLike[V]]) QueueLike[V]
	Split(group Synchronized, input QueueLike[V], size int) Sequential[QueueLike[V]]
}

// This abstract namespace defines the set of class constants, constructors and
// functions that must be supported by all set-class-like namespaces.
type SetClassLike[V Value] interface {
	// Constructors
	Empty() SetLike[V]
	FromArray(values []V) SetLike[V]
	FromSequence(values Sequential[V]) SetLike[V]
	FromSequenceWithCollator(values Sequential[V], collator CollatorLike) SetLike[V]
	FromString(values string) SetLike[V]
	WithCollator(collator CollatorLike) SetLike[V]

	// Functions
	And(first, second SetLike[V]) SetLike[V]
	Or(first, second SetLike[V]) SetLike[V]
	Sans(first, second SetLike[V]) SetLike[V]
	Xor(first, second SetLike[V]) SetLike[V]
}

// This abstract namespace defines the set of class constants, constructors and
// functions that must be supported by all stack-class-like namespaces.
type StackClassLike[V Value] interface {
	// Constants
	GetDefaultCapacity() int

	// Constructors
	Empty() StackLike[V]
	FromArray(values []V) StackLike[V]
	FromSequence(values Sequential[V]) StackLike[V]
	FromString(values string) StackLike[V]
	WithCapacity(capacity int) StackLike[V]
}
