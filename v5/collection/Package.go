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
Package "collection" declares a set of collection classes that maintain values
of a generic type:
  - Array (an extended Go array)
  - Map (an extended Go map)
  - List (a sortable list)
  - Catalog (a sortable map of key-value associations)
  - Set (an ordered set)
  - Stack (a LIFO)
  - Queue (a blocking FIFO)

For detailed documentation on this package refer to the wiki:
  - https://github.com/craterdog/go-collection-framework/wiki

This package follows the Crater Dog Technologies™ Go Coding Conventions located
here:
  - https://github.com/craterdog/go-class-model/wiki

Additional concrete implementations of the classes declared by this package can
be developed and used seamlessly since the interface declarations only depend on
other interfaces and intrinsic types—and the class implementations only depend
on interfaces, not on each other.
*/
package collection

import (
	age "github.com/craterdog/go-collection-framework/v5/agent"
)

// TYPE DECLARATIONS

/*
Index is a constrained type representing the positive (or negative) ORDINAL
index of a value in a sequence.  The indices are ordinal rather than zero-based
which never really made sense except for pointer offsets.  What is the "zeroth
value" anyway?  It's the "first value", right?  So we start a fresh...

This approach allows for positive indices starting at the beginning of a
sequence—and negative indices starting at the end of the sequence, as follows:

	    1           2           3             N
	[value 1] . [value 2] . [value 3] ... [value N]
	   -N        -(N-1)      -(N-2)          -1

Notice that because the indices are ordinal based, the positive and negative
indices are symmetrical.  An index can NEVER be zero.
*/
type Index int

// FUNCTIONAL DECLARATIONS

// CLASS DECLARATIONS

/*
ArrayClassLike[V any] is a class interface that declares the complete set of
class constructors, constants and functions that must be supported by each
concrete array-like class.

An array-like class maintains a fixed length indexed sequence of values.  Each
value is associated with an implicit positive integer index.  An array-like
class uses ORDINAL based indexing rather than the more common—and
nonsensical—ZERO based indexing scheme (see the description of what this means
in the Accessible[V] aspect definition).

This type provides a higher level abstraction for the intrinsic Go array type.
*/
type ArrayClassLike[V any] interface {
	// Constructor Methods
	Make(
		size age.Size,
	) ArrayLike[V]
	MakeFromArray(
		values []V,
	) ArrayLike[V]
	MakeFromSequence(
		values Sequential[V],
	) ArrayLike[V]
}

/*
AssociationClassLike[K comparable, V any] is a class interface that declares
the complete set of class constructors, constants and functions that must be
supported by each concrete association-like class.

An association-like class captures the relationship between a generic typed
key-value pair.
*/
type AssociationClassLike[K comparable, V any] interface {
	// Constructor Methods
	Make(
		key K,
		value V,
	) AssociationLike[K, V]
}

/*
CatalogClassLike[K comparable, V any] is a class interface that declares the
complete set of class constructors, constants and functions that must be
supported by each concrete catalog-like class.

A catalog-like class maintains an ordered set of generic typed key-value
associations.

The following class functions are supported:

Extract() returns a new catalog containing only the associations that are in
the specified catalog that have the specified keys.  The associations in the
resulting catalog will be in the same order as the specified keys.

Merge() returns a new catalog containing all of the associations that are in
the specified Catalogs in the order that they appear in each catalog.  If a
key is present in both Catalogs, the value of the key from the second
catalog takes precedence.
*/
type CatalogClassLike[K comparable, V any] interface {
	// Constructor Methods
	Make() CatalogLike[K, V]
	MakeFromArray(
		associations []AssociationLike[K, V],
	) CatalogLike[K, V]
	MakeFromMap(
		associations map[K]V,
	) CatalogLike[K, V]
	MakeFromSequence(
		associations Sequential[AssociationLike[K, V]],
	) CatalogLike[K, V]

	// Function Methods
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
ListClassLike[V any] is a class interface that declares the complete set of
class constructors, constants and functions that must be supported by each
concrete list-like class.

A list-like class maintains a dynamic sequence of values which can grow or
shrink as needed.  Each value is associated with an implicit positive integer
index.  It uses ORDINAL based indexing rather than the more common—and
nonsensical—ZERO based indexing scheme (see the description of what
this means in the Accessible[V] interface definition).

The following class functions are supported:

Concatenate() combines two lists into a new list containing all values in both
lists.  The order of the values in each list is preserved in the new list.
*/
type ListClassLike[V any] interface {
	// Constructor Methods
	Make() ListLike[V]
	MakeFromArray(
		values []V,
	) ListLike[V]
	MakeFromSequence(
		values Sequential[V],
	) ListLike[V]

	// Function Methods
	Concatenate(
		first ListLike[V],
		second ListLike[V],
	) ListLike[V]
}

/*
MapClassLike[K comparable, V any] is a class interface that declares the
complete set of class constructors, constants and functions that must be
supported by each concrete map-like class.

A map-like class extends the intrinsic Go map data type and maintains a
sequence of generic typed key-value associations.  The ordering of the
key-value associations in an intrinsic Go map is random, even for two Go maps
containing the same key-value associations.
*/
type MapClassLike[K comparable, V any] interface {
	// Constructor Methods
	Make() MapLike[K, V]
	MakeFromArray(
		associations []AssociationLike[K, V],
	) MapLike[K, V]
	MakeFromMap(
		associations map[K]V,
	) MapLike[K, V]
	MakeFromSequence(
		associations Sequential[AssociationLike[K, V]],
	) MapLike[K, V]
}

/*
QueueClassLike[V any] is a class interface that declares the complete set of
class constructors, constants and functions that must be supported by each
concrete queue-like class.

A queue-like class is generally used by multiple go-routines at the same time
and therefore enforces synchronized, first-in-first-out (FIFO) symantics for
generic typed values.  An optional queue capacity may be specified.  A request
to add a value to a queue will block when the queue has reached its maximum
capacity.  It will also block on attempts to remove a value when it is empty.
The default capacity for a queue-like class is 16 values.

The following class functions are supported:

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
	// Constructor Methods
	Make() QueueLike[V]
	MakeWithCapacity(
		capacity age.Size,
	) QueueLike[V]
	MakeFromArray(
		values []V,
	) QueueLike[V]
	MakeFromSequence(
		values Sequential[V],
	) QueueLike[V]

	// Function Methods
	Fork(
		group Synchronized,
		input QueueLike[V],
		size age.Size,
	) Sequential[QueueLike[V]]
	Split(
		group Synchronized,
		input QueueLike[V],
		size age.Size,
	) Sequential[QueueLike[V]]
	Join(
		group Synchronized,
		inputs Sequential[QueueLike[V]],
	) QueueLike[V]
}

/*
SetClassLike[V any] is a class interface that declares the complete set of
class constructors, constants and functions that must be supported by each
concrete set-like class.

A set-like class maintains an ordered sequence of unique generic typed
values—which can grow or shrink as needed.  The order of the values is
determined by a configurable collator agent.

The following class functions are supported:

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
	// Constructor Methods
	Make() SetLike[V]
	MakeWithCollator(
		collator age.CollatorLike[V],
	) SetLike[V]
	MakeFromArray(
		values []V,
	) SetLike[V]
	MakeFromSequence(
		values Sequential[V],
	) SetLike[V]

	// Function Methods
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
StackClassLike[V any] is a class interface that declares the complete set of
class constructors, constants and functions that must be supported by each
concrete stack-like class.

A stack-like class supports last-in-first-out (LIFO) symantics for generic
typed values.  An optional stack capacity may be specified.  The default
capacity is 16 values.  Adding a value to a full stack will cause an error.
*/
type StackClassLike[V any] interface {
	// Constructor Methods
	Make() StackLike[V]
	MakeWithCapacity(
		capacity age.Size,
	) StackLike[V]
	MakeFromArray(
		values []V,
	) StackLike[V]
	MakeFromSequence(
		values Sequential[V],
	) StackLike[V]
}

// INSTANCE DECLARATIONS

/*
ArrayLike[V any] is an instance interface that declares the complete set of
principal, attribute and aspect methods that must be supported by each
instance of a concrete array-like class.
*/
type ArrayLike[V any] interface {
	// Principal Methods
	GetClass() ArrayClassLike[V]

	// Aspect Interfaces
	Accessible[V]
	Sequential[V]
	Sortable[V]
	Updatable[V]
}

/*
AssociationLike[K comparable, V any] is an instance interface that declares
the complete set of principal, attribute and aspect methods that must be
supported by each instance of a concrete association-like class.
*/
type AssociationLike[K comparable, V any] interface {
	// Principal Methods
	GetClass() AssociationClassLike[K, V]

	// Attribute Methods
	GetKey() K
	GetValue() V
	SetValue(
		value V,
	)
}

/*
CatalogLike[K comparable, V any] is an instance interface that declares
the complete set of principal, attribute and aspect methods that must be
supported by each instance of a concrete catalog-like class.
*/
type CatalogLike[K comparable, V any] interface {
	// Principal Methods
	GetClass() CatalogClassLike[K, V]

	// Aspect Interfaces
	Associative[K, V]
	Sequential[AssociationLike[K, V]]
	Sortable[AssociationLike[K, V]]
}

/*
ListLike[V any] is an instance interface that declares the complete set of
principal, attribute and aspect methods that must be supported by each
instance of a concrete list-like class.
*/
type ListLike[V any] interface {
	// Principal Methods
	GetClass() ListClassLike[V]

	// Aspect Interfaces
	Accessible[V]
	Malleable[V]
	Searchable[V]
	Sequential[V]
	Sortable[V]
	Updatable[V]
}

/*
MapLike[K comparable, V any] is an instance interface that declares
the complete set of principal, attribute and aspect methods that must be
supported by each instance of a concrete map-like class.
*/
type MapLike[K comparable, V any] interface {
	// Principal Methods
	GetClass() MapClassLike[K, V]

	// Aspect Interfaces
	Associative[K, V]
	Sequential[AssociationLike[K, V]]
}

/*
QueueLike[V any] is an instance interface that declares the complete set of
principal, attribute and aspect methods that must be supported by each
instance of a concrete queue-like class.
*/
type QueueLike[V any] interface {
	// Principal Methods
	GetClass() QueueClassLike[V]

	// Attribute Methods
	GetCapacity() age.Size

	// Aspect Interfaces
	Fifo[V]
	Sequential[V]
}

/*
SetLike[V any] is an instance interface that declares the complete set of
principal, attribute and aspect methods that must be supported by each
instance of a concrete set-like class.
*/
type SetLike[V any] interface {
	// Principal Methods
	GetClass() SetClassLike[V]

	// Attribute Methods
	GetCollator() age.CollatorLike[V]

	// Aspect Interfaces
	Accessible[V]
	Elastic[V]
	Searchable[V]
	Sequential[V]
}

/*
StackLike[V any] is an instance interface that declares the complete set of
principal, attribute and aspect methods that must be supported by each
instance of a concrete stack-like class.
*/
type StackLike[V any] interface {
	// Principal Methods
	GetClass() StackClassLike[V]

	// Attribute Methods
	GetCapacity() age.Size

	// Aspect Interfaces
	Lifo[V]
	Sequential[V]
}

// ASPECT DECLARATIONS

/*
Accessible[V any] is an aspect interface that declares a set of method
signatures that must be supported by each instance of an accessible concrete
class.

An accessible class maintains values that can be accessed using indices. The
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
	GetValue(
		index Index,
	) V
	GetValues(
		first Index,
		last Index,
	) Sequential[V]
}

/*
Associative[K comparable, V any] is an aspect interface that declares a set of
method signatures that must be supported by each instance of an associative
concrete class.

An associative class maintains a sequence of generic typed key-value
associations.
*/
type Associative[K comparable, V any] interface {
	AsMap() map[K]V
	GetValue(
		key K,
	) V
	SetValue(
		key K,
		value V,
	)
	GetKeys() Sequential[K]
	GetValues(
		keys Sequential[K],
	) Sequential[V]
	RemoveValue(
		key K,
	) V
	RemoveValues(
		keys Sequential[K],
	) Sequential[V]
	RemoveAll()
}

/*
Elastic[V any] is an aspect interface that declares a set of method signatures
that must be supported by each instance of an elastic concrete class.
*/
type Elastic[V any] interface {
	AddValue(
		value V,
	)
	AddValues(
		values Sequential[V],
	)
	RemoveValue(
		value V,
	)
	RemoveValues(
		values Sequential[V],
	)
	RemoveAll()
}

/*
Fifo[V any] is an aspect interface that declares a set of method signatures
that must be supported by each instance of a synchronized first-in-first-out
channel concrete class.
*/
type Fifo[V any] interface {
	AddValue(
		value V,
	)
	RemoveFirst() (
		first V,
		ok bool,
	)
	RemoveAll()
	CloseChannel()
}

/*
Lifo[V any] is an aspect interface that declares a set of method signatures
that must be supported by each instance of a last-in-first-out class.
*/
type Lifo[V any] interface {
	AddValue(
		value V,
	)
	RemoveLast() V
	RemoveAll()
}

/*
Malleable[V any] is an aspect interface that declares a set of method
signatures that must be supported by each instance of a malleable concrete
class.
*/
type Malleable[V any] interface {
	InsertValue(
		slot age.Slot,
		value V,
	)
	InsertValues(
		slot age.Slot,
		values Sequential[V],
	)
	AppendValue(
		value V,
	)
	AppendValues(
		values Sequential[V],
	)
	RemoveValue(
		index Index,
	) V
	RemoveValues(
		first Index,
		last Index,
	) Sequential[V]
	RemoveAll()
}

/*
Searchable[V any] is an aspect interface that declares a set of method
signatures that must be supported by each instance of a searchable concrete
class.
*/
type Searchable[V any] interface {
	ContainsValue(
		value V,
	) bool
	ContainsAny(
		values Sequential[V],
	) bool
	ContainsAll(
		values Sequential[V],
	) bool
	GetIndex(
		value V,
	) Index
}

/*
Sequential[V any] is an aspect interface that declares a set of method
signatures that must be supported by each instance of a sequential concrete
class.
*/
type Sequential[V any] interface {
	IsEmpty() bool
	GetSize() age.Size
	AsArray() []V
	GetIterator() age.IteratorLike[V]
}

/*
Sortable[V any] is an aspect interface that declares a set of method
signatures that must be supported by each instance of a sortable concrete
class.

A sortable class allows its sequence of values to be sorted using a specific
sorting algorithm.
*/
type Sortable[V any] interface {
	SortValues()
	SortValuesWithRanker(
		ranker age.RankingFunction[V],
	)
	ReverseValues()
	ShuffleValues()
}

/*
Synchronized is an aspect interface that declares a set of method signatures
that must be supported by each instance of a synchronized concrete class.
*/
type Synchronized interface {
	Add(
		delta int,
	)
	Wait()
	Done()
}

/*
Updatable[V any] is an aspect interface that declares a set of method
signatures that must be supported by each instance of an updatable concrete
class.
*/
type Updatable[V any] interface {
	SetValue(
		index Index,
		value V,
	)
	SetValues(
		index Index,
		values Sequential[V],
	)
}
