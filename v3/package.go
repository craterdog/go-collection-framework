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
Package collections defines a set of simple, pragmatic abstract types and
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
  - https://github.com/craterdog/go-class-framework/wiki

Additional implementations of the classes provided by this package can be
developed and used seamlessly since the interface definitions only depend on
other interfaces and primitive types; and the class implementations only depend
on interfaces, not on each other.
*/
package collections

// TYPES

// Specializations

/*
Collection is a generic type representing any type of collections of values.
*/
type Collection any

/*
Key is a generic type representing any type of associative key.
*/
type Key any

/*
Primitive is a generic type representing any type of Go primitive value.
*/
type Primitive any

/*
Value is a generic type representing any type of value.
*/
type Value any

// Functionals

/*
ComparingFunction defines the signature for any function that can determine
whether or not two specified values are equal to each other.  The meaning of
"equality" is determined by the specific function that implements this
signature.
*/
type ComparingFunction func(first Value, second Value) bool

/*
RankingFunction defines the signature for any function that can determine
the relative ordering of two values. The result must be one of the following:

	-1: The first value is less than the second value.
	 0: The first value is equal to the second value.
	 1: The first value is more than the second value.

The meaning of "less" and "more" is determined by the specific function that
implements this signature.
*/
type RankingFunction func(first Value, second Value) int

// INTERFACES

// Aspects

/*
Accessible[V Value] defines the set of method signatures that must be
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
type Accessible[V Value] interface {
	// Methods
	GetValue(index int) V
	GetValues(first int, last int) Sequential[V]
}

/*
Associative[K Key, V Value] defines the set of method signatures that must be
supported by all sequences of key-value associations.
*/
type Associative[K Key, V Value] interface {
	// Methods
	GetKeys() Sequential[K]
	GetValue(key K) V
	GetValues(keys Sequential[K]) Sequential[V]
	RemoveAll()
	RemoveValue(key K) V
	RemoveValues(keys Sequential[K]) Sequential[V]
	SetValue(key K, value V)
}

/*
Canonical defines the set of method signatures that must be supported by all
canonical notations.
*/
type Canonical interface {
	// Methods
	FormatCollection(collection Collection) string
	ParseSource(source string) Collection
}

/*
Expandable[V Value] defines the set of method signatures that must be supported
by all sequences that allow new values to be appended, inserted and removed.
*/
type Expandable[V Value] interface {
	// Methods
	AppendValue(value V)
	AppendValues(values Sequential[V])
	InsertValue(slot int, value V)
	InsertValues(slot int, values Sequential[V])
	RemoveAll()
	RemoveValue(index int) V
	RemoveValues(first int, last int) Sequential[V]
}

/*
Flexible[V Value] defines the set of method signatures that must be supported by
all sequences of values that allow new values to be added and existing values to
be removed.
*/
type Flexible[V Value] interface {
	// Methods
	AddValue(value V)
	AddValues(values Sequential[V])
	RemoveAll()
	RemoveValue(value V)
	RemoveValues(values Sequential[V])
}

/*
Limited[V Value] defines the set of method signatures that must be supported by
all sequences of values that allow new values to be added and limit the total
number of values in the sequence.
*/
type Limited[V Value] interface {
	// Methods
	AddValue(value V)
	RemoveAll()
}

/*
Searchable[V Value] defines the set of method signatures that must be supported
by all searchable sequences of values.
*/
type Searchable[V Value] interface {
	// Methods
	ContainsAll(values Sequential[V]) bool
	ContainsAny(values Sequential[V]) bool
	ContainsValue(value V) bool
	GetIndex(value V) int
}

/*
Sequential[V Value] defines the set of method signatures that must be supported
by all sequences of values.
*/
type Sequential[V Value] interface {
	// Methods
	AsArray() []V
	GetIterator() IteratorLike[V]
	GetSize() int
	IsEmpty() bool
}

/*
Sortable[V Value] defines the set of method signatures that must be supported by
all sequences whose values may be reordered using various sorting algorithms.
*/
type Sortable[V Value] interface {
	// Methods
	ReverseValues()
	ShuffleValues()
	SortValues()
	SortValuesWithRanker(ranker RankingFunction)
}

/*
Synchronized defines the set of method signatures that must be supported by all
synchronized groups of threads.
*/
type Synchronized interface {
	// Methods
	Add(delta int)
	Done()
	Wait()
}

/*
Systematic[V Value] defines the set of method signatures that must be supported
by all systematic sorting agents.
*/
type Systematic[V Value] interface {
	// Methods
	ReverseValues(values []V)
	ShuffleValues(values []V)
	SortValues(values []V)
}

/*
Updatable[V Value] defines the set of method signatures that must be supported
by all updatable sequences of values.
*/
type Updatable[V Value] interface {
	// Methods
	SetValue(index int, value V)
	SetValues(index int, values Sequential[V])
}

// Classes

/*
ArrayClassLike[V Value] defines the set of class constants, constructors and
functions that must be supported by all array-class-like classes.
*/
type ArrayClassLike[V Value] interface {
	// Constructors
	MakeFromArray(values []V) ArrayLike[V]
	MakeFromSequence(values Sequential[V]) ArrayLike[V]
	MakeFromSize(size int) ArrayLike[V]
	MakeFromSource(source string, notation NotationLike) ArrayLike[V]
}

/*
AssociationClassLike[K Key, V Value] defines the set of class constants,
constructors and functions that must be supported by all
association-class-like classes.
*/
type AssociationClassLike[K Key, V Value] interface {
	// Constructors
	MakeWithAttributes(key K, value V) AssociationLike[K, V]
}

/*
CatalogClassLike[K comparable, V Value] defines the set of class constants,
constructors and functions that must be supported by all catalog-class-like
classes.  The following functions are supported:

Extract() returns a new catalog containing only the associations that are in
the specified catalog that have the specified keys.  The associations in the
resulting catalog will be in the same order as the specified keys.

Merge() returns a new catalog containing all of the associations that are in
the specified Catalogs in the order that they appear in each catalog.  If a
key is present in both Catalogs, the value of the key from the second
catalog takes precedence.
*/
type CatalogClassLike[K comparable, V Value] interface {
	// Constructors
	Make() CatalogLike[K, V]
	MakeFromArray(associations []AssociationLike[K, V]) CatalogLike[K, V]
	MakeFromMap(associations map[K]V) CatalogLike[K, V]
	MakeFromSequence(associations Sequential[AssociationLike[K, V]]) CatalogLike[K, V]
	MakeFromSource(source string, notation NotationLike) CatalogLike[K, V]

	// Functions
	Extract(catalog CatalogLike[K, V], keys Sequential[K]) CatalogLike[K, V]
	Merge(first CatalogLike[K, V], second CatalogLike[K, V]) CatalogLike[K, V]
}

/*
CollatorClassLike defines the set of class constants, constructors and functions
that must be supported by all collator-class-like classes.
*/
type CollatorClassLike interface {
	// Constants
	DefaultMaximum() int

	// Constructors
	Make() CollatorLike
	MakeWithMaximum(maximum int) CollatorLike
}

/*
FormatterClassLike defines the set of class constants, constructors and
functions that must be supported by all formatter-class-like classes.
*/
type FormatterClassLike interface {
	// Constants
	DefaultMaximum() int

	// Constructors
	Make() FormatterLike
	MakeWithMaximum(maximum int) FormatterLike
}

/*
IteratorClassLike[V Value] defines the set of class constants, constructors and
functions that must be supported by all iterator-class-like classes.
*/
type IteratorClassLike[V Value] interface {
	// Constructors
	MakeFromSequence(values Sequential[V]) IteratorLike[V]
}

/*
ListClassLike[V Value] defines the set of class constants, constructors and
functions that must be supported by all list-class-like classes.  The following
functions are supported:

Concatenate() combines two lists into a new list containing all values in both
lists.  The order of the values in each list is preserved in the new list.
*/
type ListClassLike[V Value] interface {
	// Constructors
	Make() ListLike[V]
	MakeFromArray(values []V) ListLike[V]
	MakeFromSequence(values Sequential[V]) ListLike[V]
	MakeFromSource(source string, notation NotationLike) ListLike[V]

	// Functions
	Concatenate(first ListLike[V], second ListLike[V]) ListLike[V]
}

/*
MapClassLike[K comparable, V Value] defines the set of class constants,
constructors and functions that must be supported by all map-class-like
classes.
*/
type MapClassLike[K comparable, V Value] interface {
	// Constructors
	Make() MapLike[K, V]
	MakeFromArray(associations []AssociationLike[K, V]) MapLike[K, V]
	MakeFromMap(associations map[K]V) MapLike[K, V]
	MakeFromSequence(associations Sequential[AssociationLike[K, V]]) MapLike[K, V]
	MakeFromSource(source string, notation NotationLike) MapLike[K, V]
}

/*
NotationClassLike defines the set of class constants, constructors and
functions that must be supported by all notation-class-like classes.
*/
type NotationClassLike interface {
	// Constructors
	Make() NotationLike
}

/*
ParserClassLike defines the set of class constants, constructors and functions
that must be supported by all parser-class-like classes.
*/
type ParserClassLike interface {
	// Constructors
	Make() ParserLike
}

/*
QueueClassLike[V Value] defines the set of class constants, constructors and
functions that must be supported by all queue-class-like classes.  The following
functions are supported:

Fork() connects the output of the specified input Queue with a number of new
output queues specified by the size parameter and returns a sequence of the new
output queues. Each value added to the input queue will be added automatically
to ALL of the output queues. This pattern is useful when a set of DIFFERENT
operations needs to occur for every value and each operation can be done in
parallel.

Join() connects the outputs of the specified sequence of input queues with a new
output queue returns the new output queue. Each value removed from each input
queue will automatically be added to the output queue.  This pattern is useful
when the results of the processing with a Split() function need to be
consolidated into a single queue.

Split() connects the output of the specified input Queue with the number of
output queues specified by the size parameter and returns a sequence of the new
output queues. Each value added to the input queue will be added automatically
to ONE of the output queues. This pattern is useful when a SINGLE operation
needs to occur for each value and the operation can be done on the values in
parallel.  The results can then be consolidated later on using the Join()
function.
*/
type QueueClassLike[V Value] interface {
	// Constants
	DefaultCapacity() int

	// Constructors
	Make() QueueLike[V]
	MakeFromArray(values []V) QueueLike[V]
	MakeFromSequence(values Sequential[V]) QueueLike[V]
	MakeFromSource(source string, notation NotationLike) QueueLike[V]
	MakeWithCapacity(capacity int) QueueLike[V]

	// Functions
	Fork(
		group Synchronized,
		input QueueLike[V],
		size int,
	) Sequential[QueueLike[V]]
	Join(group Synchronized, inputs Sequential[QueueLike[V]]) QueueLike[V]
	Split(
		group Synchronized,
		input QueueLike[V],
		size int,
	) Sequential[QueueLike[V]]
}

/*
SetClassLike[V Value] defines the set of class constants, constructors and
functions that must be supported by all set-class-like classes.  The following
functions are supported:

And() returns a new set containing the values that are both of the specified
sets.

Or() returns a new set containing the values that are in either of the specified
sets.

Sans() returns a new set containing the values that are in the first specified
set but not in the second specified set.

Xor() returns a new set containing the values that are in the first specified
set or the second specified set but not both.
*/
type SetClassLike[V Value] interface {
	// Constructors
	Make() SetLike[V]
	MakeFromArray(values []V) SetLike[V]
	MakeFromSequence(values Sequential[V]) SetLike[V]
	MakeFromSource(source string, notation NotationLike) SetLike[V]
	MakeWithCollator(collator CollatorLike) SetLike[V]

	// Functions
	And(first SetLike[V], second SetLike[V]) SetLike[V]
	Or(first SetLike[V], second SetLike[V]) SetLike[V]
	Sans(first SetLike[V], second SetLike[V]) SetLike[V]
	Xor(first SetLike[V], second SetLike[V]) SetLike[V]
}

/*
SorterClassLike[V Value] defines the set of class constants, constructors and
functions that must be supported by all sorter-class-like classes.
*/
type SorterClassLike[V Value] interface {
	// Constants
	DefaultRanker() RankingFunction

	// Constructors
	Make() SorterLike[V]
	MakeWithRanker(ranker RankingFunction) SorterLike[V]
}

/*
StackClassLike[V Value] defines the set of class constants, constructors and
functions that must be supported by all stack-class-like classes.
*/
type StackClassLike[V Value] interface {
	// Constants
	DefaultCapacity() int

	// Constructors
	Make() StackLike[V]
	MakeFromArray(values []V) StackLike[V]
	MakeFromSequence(values Sequential[V]) StackLike[V]
	MakeFromSource(source string, notation NotationLike) StackLike[V]
	MakeWithCapacity(capacity int) StackLike[V]
}

// Instances

/*
ArrayLike[V Value] defines the set of abstractions and methods that must be
supported by all array-like instances.  An array-like class maintains a fixed
length indexed sequence of values.  Each value is associated with an implicit
positive integer index. An array-like class uses ORDINAL based indexing rather
than the more common—and nonsensical—ZERO based indexing scheme (see the
description of what this means in the Accessible interface definition).

This type is parameterized as follows:
  - V is any type of value.

This type essentially provides a higher level abstraction for the primitive Go
array type.
*/
type ArrayLike[V Value] interface {
	// Abstractions
	Accessible[V]
	Sequential[V]
	Sortable[V]
	Updatable[V]
}

/*
AssociationLike[K Key, V Value] defines the set of abstractions and methods that
must be supported by all association-like instances.  An association-like class
maintains information about a key-value association.

This type is parameterized as follows:
  - K is a primitive type of key.
  - V is any type of value.

This type is used by catalog-like instances to maintain their associations.
*/
type AssociationLike[K Key, V Value] interface {
	// Attributes
	GetKey() K
	GetValue() V
	SetValue(value V)
}

/*
CatalogLike[K Key, V Value] defines the set of abstractions and methods that
must be supported by all catalog-like instances.  A catalog-like class maintains
a sequence of key-value associations.

This type is parameterized as follows:
  - K is a primitive type of key.
  - V is any type of entity.

A catalog-like class can use any association-like class key-value association.
*/
type CatalogLike[K Key, V Value] interface {
	// Abstractions
	Associative[K, V]
	Sequential[AssociationLike[K, V]]
	Sortable[AssociationLike[K, V]]
}

/*
CollatorLike defines the set of abstractions and methods that must be supported
by all collator-like instances.  A collator-like class is capable of comparing
and ranking two values of any type.
*/
type CollatorLike interface {
	// Attributes
	GetDepth() int
	GetMaximum() int

	// Methods
	CompareValues(first Value, second Value) bool
	RankValues(first Value, second Value) int
}

/*
FormatterLike defines the set of abstractions and methods that must be supported
by all formatter-like instances.
*/
type FormatterLike interface {
	// Attributes
	GetDepth() int
	GetMaximum() int

	// Methods
	FormatCollection(collection Collection) string
}

/*
IteratorLike[V Value] defines the set of abstractions and methods that must be
supported by all iterator-like instances.  An iterator-like class can be used to
move forward and backward over the values in a sequence.  It implements the Gang
of Four (GoF) Iterator Design Pattern:
  - https://en.wikipedia.org/wiki/Iterator_pattern

A iterator agent locks into the slots that reside between each value in the
sequence:

	    [value 1] . [value 2] . [value 3] ... [value N]
	  ^           ^           ^                         ^
	slot 0      slot 1      slot 2                    slot N

It moves from slot to slot and has access to the values (if they exist) on each
side of the slot.  At each slot an iterator has access to the previous value
and next value in the sequence (assuming they exist). The slot at the start of
the sequence has no PREVIOUS value, and the slot at the end of the sequence has
no NEXT value.

This type is parameterized as follows:
  - V is any type of value.

An iterator-like class is supported by all collection types.
*/
type IteratorLike[V Value] interface {
	// Methods
	GetNext() V
	GetPrevious() V
	GetSlot() int
	HasNext() bool
	HasPrevious() bool
	ToEnd()
	ToSlot(slot int)
	ToStart()
}

/*
ListLike[V Value] defines the set of abstractions and methods that must be
supported by all list-like instances.  A list-like class maintains a dynamic
sequence of values which can grow or shrink as needed.  Each value is associated
with an implicit positive integer index. An array-like class uses ORDINAL based
indexing rather than the more common—and nonsensical—ZERO based indexing scheme
(see the description of what this means in the Accessible interface definition).

This type is parameterized as follows:
  - V is any type of value.

All comparison and ranking of values in the sequence is done using the default
collator.
*/
type ListLike[V Value] interface {
	// Abstractions
	Accessible[V]
	Expandable[V]
	Searchable[V]
	Sequential[V]
	Sortable[V]
	Updatable[V]
}

/*
MapLike[K Key, V Value] defines the set of abstractions and methods that must be
supported by all map-like instances.  A map-like class extends the primitive Go
map type and maintains a sequence of key-value associations.  The order of the
key-value associations in a primitive Go map is random, even for two Go maps
containing the same key-value associations.

This type is parameterized as follows:
  - K is a primitive type of key.
  - V is any type of entity.

A map-like class can use any association-like class key-value association.
*/
type MapLike[K Key, V Value] interface {
	// Abstractions
	Associative[K, V]
	Sequential[AssociationLike[K, V]]
}

/*
NotationLike defines the set of abstractions and methods that must be supported
by all notation-like instances.  A notation-like class can be used to parse and
format collections using a canonical notation like XML, JSON and  CDCN
(Crater Dog Collection Notation™).
*/
type NotationLike interface {
	// Abstractions
	Canonical
}

/*
ParserLike defines the set of abstractions and methods that must be supported by
all parser-like instances.
*/
type ParserLike interface {
	// Methods
	ParseSource(source string) Collection
}

/*
QueueLike[V Value] defines the set of abstractions and methods that must be
supported by all queue-like instances.  A queue-like class implements FIFO
(i.e.  first-in-first-out) semantics.

This type is parameterized as follows:
  - V is any type of value.

A queue-like class is generally used by multiple go-routines at the same time
and therefore enforces synchronized access.  A queue-like class enforces a
maximum length and will block on attempts to add a value it is full.  It will
also block on attempts to remove a value when it is empty.
*/
type QueueLike[V Value] interface {
	// Attributes
	GetCapacity() int

	// Abstractions
	Limited[V]
	Sequential[V]

	// Methods
	CloseQueue()
	RemoveHead() (head V, ok bool)
}

/*
ScannerLike defines the set of abstractions and methods that must be supported
by all scanner-like instances.
*/
type ScannerLike interface {
}

/*
SetLike[V Value] defines the set of abstractions and methods that must be
supported by all set-like instances.  A set-like class maintains an ordered
sequence of values which can grow or shrink as needed.

This type is parameterized as follows:
  - V is any type of value.

The order of the values is determined by a configurable RankingFunction.
*/
type SetLike[V Value] interface {
	// Attributes
	GetCollator() CollatorLike

	// Abstractions
	Accessible[V]
	Flexible[V]
	Searchable[V]
	Sequential[V]
}

/*
SorterLike[V Value] defines the set of abstractions and methods that must be
supported by all sorter-like instances.  A sorter-like class implements a
specific sorting algorithm.

This type is parameterized as follows:
  - V is any type of value.

A sorter-like class uses a ranking function to order the values.  If no ranking
function is specified the values are sorted into their natural order.
*/
type SorterLike[V Value] interface {
	// Attributes
	GetRanker() RankingFunction

	// Abstractions
	Systematic[V]
}

/*
StackLike[V Value] defines the set of abstractions and methods that must be
supported by all stack-like instances.  A stack-like class implements LIFO
(i.e.  last-in-first-out) semantics.

This type is parameterized as follows:
  - V is any type of value.

A stack-like class enforces a maximum depth and will panic if that depth is
exceeded.  It will also panic on attempts to remove a value when it is empty.
*/
type StackLike[V Value] interface {
	// Attributes
	GetCapacity() int

	// Abstractions
	Limited[V]
	Sequential[V]

	// Methods
	RemoveTop() V
}
