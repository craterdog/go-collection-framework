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

/*
Package "collection" defines a set of simple, pragmatic abstract types and
interfaces for Go based collections of values. It also provide an efficient and
compact implementation of the following collection classes based on these
abstractions:
  - Array (extended Go array)
  - Map (extended Go map)
  - List (a sortable list)
  - Catalog (a sortable map)
  - Set (an ordered set)
  - Stack (a LIFO)
  - Queue (a blocking FIFO)

For detailed documentation on this package refer to the wiki:
  - https://github.com/craterdog/go-collection-framework/wiki

This package follows the Crater Dog Technologies™ (craterdog) Go Coding
Conventions located here:
  - https://github.com/craterdog/go-model-framework/wiki

Additional implementations of the classes provided by this package can be
developed and used seamlessly since the interface definitions only depend on
other interfaces and primitive types; and the class implementations only depend
on interfaces, not on each other.
*/
package collection

import (
	age "github.com/craterdog/go-collection-framework/v4/agent"
)

// Classes

/*
ArrayClassLike[V any] is a class interface that defines the complete set of
class constructors, constants and functions that must be supported by each
concrete array-like class.
*/
type ArrayClassLike[V any] interface {
	// Constructors
	MakeWithSize(size uint) ArrayLike[V]
	MakeFromArray(values []V) ArrayLike[V]
	MakeFromSequence(values Sequential[V]) ArrayLike[V]

	// Constants
	Notation() NotationLike
}

/*
AssociationClassLike[K comparable, V any] is a class interface that defines
the complete set of class constructors, constants and functions that must be
supported by each concrete association-like class.
*/
type AssociationClassLike[
	K comparable,
	V any,
] interface {
	// Constructors
	MakeWithAttributes(
		key K,
		value V,
	) AssociationLike[K, V]

	// Constants
	Notation() NotationLike
}

/*
CatalogClassLike[K comparable, V any] is a class interface that defines the
complete set of class constructors, constants and functions that must be
supported by each concrete catalog-like class.

The following functions are supported:

Extract() returns a new catalog containing only the associations that are in
the specified catalog that have the specified keys.  The associations in the
resulting catalog will be in the same order as the specified keys.

Merge() returns a new catalog containing all of the associations that are in
the specified Catalogs in the order that they appear in each catalog.  If a
key is present in both Catalogs, the value of the key from the second
catalog takes precedence.
*/
type CatalogClassLike[
	K comparable,
	V any,
] interface {
	// Constructors
	Make() CatalogLike[K, V]
	MakeFromArray(associations []AssociationLike[K, V]) CatalogLike[K, V]
	MakeFromMap(associations map[K]V) CatalogLike[K, V]
	MakeFromSequence(associations Sequential[AssociationLike[K, V]]) CatalogLike[K, V]

	// Constants
	Notation() NotationLike

	// Functions
	Extract(
		catalog CatalogLike[K, V],
		keys Sequential[K],
	) CatalogLike[K, V]
	Merge(
		first CatalogLike[K, V],
		second CatalogLike[K, V],
	) CatalogLike[K, V]
}

/*
ListClassLike[V any] is a class interface that defines the complete set of
class constructors, constants and functions that must be supported by each
concrete list-like class.

The following functions are supported:

Concatenate() combines two lists into a new list containing all values in both
lists.  The order of the values in each list is preserved in the new list.
*/
type ListClassLike[V any] interface {
	// Constructors
	Make() ListLike[V]
	MakeFromArray(values []V) ListLike[V]
	MakeFromSequence(values Sequential[V]) ListLike[V]

	// Constants
	Notation() NotationLike

	// Functions
	Concatenate(
		first ListLike[V],
		second ListLike[V],
	) ListLike[V]
}

/*
MapClassLike[K comparable, V any] is a class interface that defines the
complete set of class constructors, constants and functions that must be
supported by each concrete map-like class.
*/
type MapClassLike[
	K comparable,
	V any,
] interface {
	// Constructors
	Make() MapLike[K, V]
	MakeFromArray(associations []AssociationLike[K, V]) MapLike[K, V]
	MakeFromMap(associations map[K]V) MapLike[K, V]
	MakeFromSequence(associations Sequential[AssociationLike[K, V]]) MapLike[K, V]

	// Constants
	Notation() NotationLike
}

/*
NotationClassLike is a class interface that defines the complete set of class
constructors, constants and functions that must be supported by each concrete
notation-like class.
*/
type NotationClassLike interface {
	// Constructors
	Make() NotationLike
}

/*
QueueClassLike[V any] is a class interface that defines the complete set of
class constructors, constants and functions that must be supported by each
concrete queue-like class.

The following functions are supported:

Fork() connects the output of the specified input Queue with a number of new
output queues specified by the size parameter and returns a sequence of the new
output queues. Each value added to the input queue will be added automatically
to ALL of the output queues. This pattern is useful when a set of DIFFERENT
operations needs to occur for every value and each operation can be done in
parallel.

Split() connects the output of the specified input Queue with the number of
output queues specified by the size parameter and returns a sequence of the new
output queues. Each value added to the input queue will be added automatically
to ONE of the output queues. This pattern is useful when a SINGLE operation
needs to occur for each value and the operation can be done on the values in
parallel.  The results can then be consolidated later on using the Join()
function.

Join() connects the outputs of the specified sequence of input queues with a new
output queue returns the new output queue. Each value removed from each input
queue will automatically be added to the output queue.  This pattern is useful
when the results of the processing with a Split() function need to be
consolidated into a single queue.
*/
type QueueClassLike[V any] interface {
	// Constructors
	Make() QueueLike[V]
	MakeWithCapacity(capacity uint) QueueLike[V]
	MakeFromArray(values []V) QueueLike[V]
	MakeFromSequence(values Sequential[V]) QueueLike[V]

	// Constants
	Notation() NotationLike
	DefaultCapacity() uint

	// Functions
	Fork(
		group Synchronized,
		input QueueLike[V],
		size uint,
	) Sequential[QueueLike[V]]
	Split(
		group Synchronized,
		input QueueLike[V],
		size uint,
	) Sequential[QueueLike[V]]
	Join(
		group Synchronized,
		inputs Sequential[QueueLike[V]],
	) QueueLike[V]
}

/*
SetClassLike[V any] is a class interface that defines the complete set of
class constructors, constants and functions that must be supported by each
concrete set-like class.

The following functions are supported:

And() returns a new set containing the values that are both of the specified
sets.

Or() returns a new set containing the values that are in either of the specified
sets.

Sans() returns a new set containing the values that are in the first specified
set but not in the second specified set.

Xor() returns a new set containing the values that are in the first specified
set or the second specified set but not both.
*/
type SetClassLike[V any] interface {
	// Constructors
	Make() SetLike[V]
	MakeWithCollator(collator age.CollatorLike[V]) SetLike[V]
	MakeFromArray(values []V) SetLike[V]
	MakeFromSequence(values Sequential[V]) SetLike[V]

	// Constants
	Notation() NotationLike

	// Functions
	And(
		first SetLike[V],
		second SetLike[V],
	) SetLike[V]
	Or(
		first SetLike[V],
		second SetLike[V],
	) SetLike[V]
	Sans(
		first SetLike[V],
		second SetLike[V],
	) SetLike[V]
	Xor(
		first SetLike[V],
		second SetLike[V],
	) SetLike[V]
}

/*
StackClassLike[V any] is a class interface that defines the complete set of
class constructors, constants and functions that must be supported by each
concrete stack-like class.
*/
type StackClassLike[V any] interface {
	// Constructors
	Make() StackLike[V]
	MakeWithCapacity(capacity uint) StackLike[V]
	MakeFromArray(values []V) StackLike[V]
	MakeFromSequence(values Sequential[V]) StackLike[V]

	// Constants
	Notation() NotationLike
	DefaultCapacity() uint
}

// Instances

/*
ArrayLike[V any] is an instance interface that defines the complete set of
instance attributes, abstractions and methods that must be supported by each
instance of a concrete array-like class.  An array-like class maintains a fixed
length indexed sequence of values.  Each value is associated with an implicit
positive integer index. An array-like class uses ORDINAL based indexing rather
than the more common—and nonsensical—ZERO based indexing scheme (see the
description of what this means in the Accessible interface definition).

This type is parameterized as follows:
  - V is any type of value.

This type essentially provides a higher level abstraction for the primitive Go
array type.
*/
type ArrayLike[V any] interface {
	// Attributes
	GetClass() ArrayClassLike[V]

	// Abstractions
	Accessible[V]
	Sequential[V]
	Sortable[V]
	Updatable[V]
}

/*
AssociationLike[K comparable, V any] is an instance interface that defines
the complete set of instance attributes, abstractions and methods that must be
supported by each instance of a concrete association-like class.

This type is parameterized as follows:
  - K is a primitive type of key.
  - V is any type of value.

This type is used by catalog-like instances to maintain their associations.
*/
type AssociationLike[
	K comparable,
	V any,
] interface {
	// Attributes
	GetClass() AssociationClassLike[K, V]
	GetKey() K
	GetValue() V
	SetValue(value V)
}

/*
CatalogLike[K comparable, V any] is an instance interface that defines the
complete set of instance attributes, abstractions and methods that must be
supported by each instance of a concrete catalog-like class.

This type is parameterized as follows:
  - K is a primitive type of key.
  - V is any type of entity.

A catalog-like class can use any association-like class key-value association.
*/
type CatalogLike[
	K comparable,
	V any,
] interface {
	// Attributes
	GetClass() CatalogClassLike[K, V]

	// Abstractions
	Associative[K, V]
	Sequential[AssociationLike[K, V]]
	Sortable[AssociationLike[K, V]]
}

/*
ListLike[V any] is an instance interface that defines the complete set of
instance attributes, abstractions and methods that must be supported by each
instance of a concrete list-like class.  A list-like class maintains a dynamic
sequence of values which can grow or shrink as needed.  Each value is associated
with an implicit positive integer index. An array-like class uses ORDINAL based
indexing rather than the more common—and nonsensical—ZERO based indexing scheme
(see the description of what this means in the Accessible interface definition).

This type is parameterized as follows:
  - V is any type of value.

All comparison and ranking of values in the sequence is done using the default
collator.
*/
type ListLike[V any] interface {
	// Attributes
	GetClass() ListClassLike[V]

	// Abstractions
	Accessible[V]
	Expandable[V]
	Searchable[V]
	Sequential[V]
	Sortable[V]
	Updatable[V]
}

/*
MapLike[K comparable, V any] is an instance interface that defines the
complete set of instance attributes, abstractions and methods that must be
supported by each instance of a concrete map-like class.  A map-like class
extends the primitive Go map type and maintains a sequence of key-value
associations.  The order of the key-value associations in a primitive Go map is
random, even for two Go maps containing the same key-value associations.

This type is parameterized as follows:
  - K is a primitive type of key.
  - V is any type of entity.

A map-like class can use any association-like class key-value association.
*/
type MapLike[
	K comparable,
	V any,
] interface {
	// Attributes
	GetClass() MapClassLike[K, V]

	// Abstractions
	Associative[K, V]
	Sequential[AssociationLike[K, V]]
}

/*
NotationLike is an instance interface that defines the complete set of instance
attributes, abstractions and methods that must be supported by each instance of
a concrete notation-like class.  A notation-like class can be used to parse and
format collections using a canonical notation like XML, JSON and CDCN (Crater
Dog Collection Notation™).
*/
type NotationLike interface {
	// Attributes
	GetClass() NotationClassLike

	// Abstractions
	Canonical
}

/*
QueueLike[V any] is an instance interface that defines the complete set of
instance attributes, abstractions and methods that must be supported by each
instance of a concrete queue-like class.  A queue-like class implements FIFO
(i.e. first-in-first-out) semantics.

This type is parameterized as follows:
  - V is any type of value.

A queue-like class is generally used by multiple go-routines at the same time
and therefore enforces synchronized access.  A queue-like class enforces a
maximum length and will block on attempts to add a value it is full.  It will
also block on attempts to remove a value when it is empty.
*/
type QueueLike[V any] interface {
	// Attributes
	GetClass() QueueClassLike[V]
	GetCapacity() uint

	// Abstractions
	Limited[V]
	Sequential[V]

	// Methods
	RemoveHead() (
		head V,
		ok bool,
	)
	CloseQueue()
}

/*
SetLike[V any] is an instance interface that defines the complete set of
instance attributes, abstractions and methods that must be supported by each
instance of a concrete set-like class.  A set-like class maintains an ordered
sequence of values which can grow or shrink as needed.

This type is parameterized as follows:
  - V is any type of value.

The order of the values is determined by a configurable CollatorLike[V] agent.
*/
type SetLike[V any] interface {
	// Attributes
	GetClass() SetClassLike[V]
	GetCollator() age.CollatorLike[V]

	// Abstractions
	Accessible[V]
	Flexible[V]
	Searchable[V]
	Sequential[V]
}

/*
StackLike[V any] is an instance interface that defines the complete set of
instance attributes, abstractions and methods that must be supported by each
instance of a concrete stack-like class.  A stack-like class implements LIFO
(i.e. last-in-first-out) semantics.

This type is parameterized as follows:
  - V is any type of value.

A stack-like class enforces a maximum depth and will panic if that depth is
exceeded.  It will also panic on attempts to remove a value when it is empty.
*/
type StackLike[V any] interface {
	// Attributes
	GetClass() StackClassLike[V]
	GetCapacity() uint

	// Abstractions
	Limited[V]
	Sequential[V]

	// Methods
	RemoveTop() V
}

// Aspects

/*
Accessible[V any] defines the set of method signatures that must be
supported by all sequences whose values can be accessed using indices. The
indices of an accessible sequence are ORDINAL rather than ZERO based—which
never really made sense except for pointer offsets. What is the "zeroth
value" anyway? It's the "first value", right?  So we start fresh...

This approach allows for positive indices starting at the beginning of the
sequence, and negative indices starting at the end of the sequence as follows:

	    1           2           3             N
	[value 1] . [value 2] . [value 3] ... [value N]
	   -N        -(N-1)      -(N-2)          -1

Notice that because the indices are ordinal based, the positive and negative
indices are symmetrical.
*/
type Accessible[V any] interface {
	GetValue(index int) V
	GetValues(
		first int,
		last int,
	) Sequential[V]
}

/*
Associative[K comparable, V any] defines the set of method signatures that
must be supported by all sequences of key-value associations.
*/
type Associative[
	K comparable,
	V any,
] interface {
	GetValue(key K) V
	SetValue(
		key K,
		value V,
	)
	GetKeys() Sequential[K]
	GetValues(keys Sequential[K]) Sequential[V]
	RemoveValue(key K) V
	RemoveValues(keys Sequential[K]) Sequential[V]
	RemoveAll()
}

/*
Canonical defines the set of method signatures that must be supported by all
canonical notations.
*/
type Canonical interface {
	ParseSource(source string) (value any)
	FormatValue(value any) (source string)
}

/*
Expandable[V any] defines the set of method signatures that must be supported
by all sequences that allow new values to be appended, inserted and removed.
*/
type Expandable[V any] interface {
	InsertValue(
		slot uint,
		value V,
	)
	InsertValues(
		slot uint,
		values Sequential[V],
	)
	AppendValue(value V)
	AppendValues(values Sequential[V])
	RemoveValue(index int) V
	RemoveValues(
		first int,
		last int,
	) Sequential[V]
	RemoveAll()
}

/*
Flexible[V any] defines the set of method signatures that must be supported by
all sequences of values that allow new values to be added and existing values to
be removed.
*/
type Flexible[V any] interface {
	AddValue(value V)
	AddValues(values Sequential[V])
	RemoveValue(value V)
	RemoveValues(values Sequential[V])
	RemoveAll()
}

/*
Limited[V any] defines the set of method signatures that must be supported by
all sequences of values that allow new values to be added and limit the total
number of values in the sequence.
*/
type Limited[V any] interface {
	AddValue(value V)
	RemoveAll()
}

/*
Searchable[V any] defines the set of method signatures that must be supported
by all searchable sequences of values.
*/
type Searchable[V any] interface {
	ContainsValue(value V) bool
	ContainsAny(values Sequential[V]) bool
	ContainsAll(values Sequential[V]) bool
	GetIndex(value V) int
}

/*
Sequential[V any] defines the set of method signatures that must be supported
by all sequences of values.  Note that sizes should be of type "uint" but the Go
language does not allow arithmetic and comparison operations between "int" and
"uint" so we us "int" for the return type to make it easier to use.
*/
type Sequential[V any] interface {
	IsEmpty() bool
	GetSize() int
	AsArray() []V
	GetIterator() age.IteratorLike[V]
}

/*
Sortable[V any] defines the set of method signatures that must be supported by
all sequences whose values may be reordered using various sorting algorithms.
*/
type Sortable[V any] interface {
	SortValues()
	SortValuesWithRanker(ranker age.RankingFunction[V])
	ReverseValues()
	ShuffleValues()
}

/*
Synchronized defines the set of method signatures that must be supported by all
synchronized groups of threads.
*/
type Synchronized interface {
	Add(delta int)
	Wait()
	Done()
}

/*
Updatable[V any] defines the set of method signatures that must be supported
by all updatable sequences of values.
*/
type Updatable[V any] interface {
	SetValue(
		index int,
		value V,
	)
	SetValues(
		index int,
		values Sequential[V],
	)
}
