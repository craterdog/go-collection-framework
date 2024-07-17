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
Package "agent" defines a set of agents that operate on values of a generic
type.

This package follows the Crater Dog Technologies™ (craterdog) Go Coding
Conventions located here:
  - https://github.com/craterdog/go-model-framework/wiki

Additional implementations of the classes provided by this package can be
developed and used seamlessly since the interface definitions only depend on
other interfaces and intrinsic types; and the class implementations only depend
on interfaces, not on each other.
*/
package agent

// Types

/*
Rank is a constrained type representing the possible rankings for two values.
*/
type Rank uint8

const (
	LesserRank Rank = iota
	EqualRank
	GreaterRank
)

// Functionals

/*
RankingFunction[V any] defines the signature for any function that can determine
the relative ordering of two values. The result must be one of the following:

	LesserRank: The first value is less than the second value.
	EqualRank: The first value is equal to the second value.
	GreaterRank: The first value is more than the second value.

The meaning of "lesser" and "greater" is determined by the specific function
that implements this signature.
*/
type RankingFunction[V any] func(
	first V,
	second V,
) Rank

// Classes

/*
CollatorClassLike[V any] defines the set of class constructors, constants and
functions that must be supported by all collator-class-like classes.
*/
type CollatorClassLike[V any] interface {
	// Constructors
	Make() CollatorLike[V]
	MakeWithMaximum(maximum int) CollatorLike[V]

	// Constants
	DefaultMaximum() int
}

/*
InspectorClassLike defines the set of class constructors, constants and
functions that must be supported by all inspector-class-like classes.
*/
type InspectorClassLike interface {
	// Constructors
	Make() InspectorLike
}

/*
IteratorClassLike[V any] defines the set of class constructors, constants and
functions that must be supported by all iterator-class-like classes.
*/
type IteratorClassLike[V any] interface {
	// Constructors
	MakeFromArray(values []V) IteratorLike[V]
}

/*
SorterClassLike[V any] defines the set of class constructors, constants and
functions that must be supported by all sorter-class-like classes.
*/
type SorterClassLike[V any] interface {
	// Constructors
	Make() SorterLike[V]
	MakeWithRanker(ranker RankingFunction[V]) SorterLike[V]

	// Constants
	DefaultRanker() RankingFunction[V]
}

// Instances

/*
CollatorLike[V any] defines the set of abstractions and methods that must be
supported by all collator-like instances.  A collator-like class is capable of
comparing and ranking two values of any type.
*/
type CollatorLike[V any] interface {
	// Attributes
	GetClass() CollatorClassLike[V]
	GetDepth() int
	GetMaximum() int

	// Methods
	CompareValues(
		first V,
		second V,
	) bool
	RankValues(
		first V,
		second V,
	) Rank
}

/*
InspectorLike[V any] defines the set of abstractions and methods that must be
supported by all inspector-like instances.  An inspector-like class is capable
of determining whether or not values possess certain properties.
*/
type InspectorLike interface {
	// Attributes
	GetClass() InspectorClassLike

	// Methods
	ImplementsAspect(
		value any,
		aspect any,
	) bool
	IsDefined(value any) bool
}

/*
IteratorLike[V any] defines the set of abstractions and methods that must be
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
type IteratorLike[V any] interface {
	// Attributes
	GetClass() IteratorClassLike[V]
	GetSize() int
	GetSlot() int

	// Methods
	IsEmpty() bool
	ToStart()
	ToSlot(slot int)
	ToEnd()
	GetNext() V
	GetPrevious() V
	HasNext() bool
	HasPrevious() bool
}

/*
SorterLike[V any] defines the set of abstractions and methods that must be
supported by all sorter-like instances.  A sorter-like class implements a
specific sorting algorithm.

This type is parameterized as follows:
  - V is any type of value.

A sorter-like class uses a ranking function to order the values.  If no ranking
function is specified the values are sorted into their natural order.
*/
type SorterLike[V any] interface {
	// Attributes
	GetClass() SorterClassLike[V]
	GetRanker() RankingFunction[V]

	// Abstractions
	Systematic[V]
}

// Aspects

/*
Systematic[V any] defines the set of method signatures that must be supported
by all systematic sorting agents.
*/
type Systematic[V any] interface {
	SortValues(values []V)
	ReverseValues(values []V)
	ShuffleValues(values []V)
}
