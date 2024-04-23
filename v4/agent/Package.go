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
other interfaces and primitive types; and the class implementations only depend
on interfaces, not on each other.
*/
package agent

// Types

/*
Value is a generic type representing any type of value.
*/
type Value any

// Functionals

/*
RankingFunction[V Value] defines the signature for any function that can
determine the relative ordering of two values. The result must be one of the
following:

	-1: The first value is less than the second value.
	 0: The first value is equal to the second value.
	 1: The first value is more than the second value.

The meaning of "less" and "more" is determined by the specific function that
implements this signature.
*/
type RankingFunction[V Value] func(
	first V,
	second V,
) int

// Aspects

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

// Classes

/*
CollatorClassLike[V Value] defines the set of class constants, constructors and
functions that must be supported by all collator-class-like classes.
*/
type CollatorClassLike[V Value] interface {
	// Constants
	DefaultMaximum() int

	// Constructors
	Make() CollatorLike[V]
	MakeWithMaximum(maximum int) CollatorLike[V]
}

/*
IteratorClassLike[V Value] defines the set of class constants, constructors and
functions that must be supported by all iterator-class-like classes.
*/
type IteratorClassLike[V Value] interface {
	// Constructors
	MakeFromArray(values []V) IteratorLike[V]
}

/*
SorterClassLike[V Value] defines the set of class constants, constructors and
functions that must be supported by all sorter-class-like classes.
*/
type SorterClassLike[V Value] interface {
	// Constants
	DefaultRanker() RankingFunction[V]

	// Constructors
	Make() SorterLike[V]
	MakeWithRanker(ranker RankingFunction[V]) SorterLike[V]
}

// Instances

/*
CollatorLike[V Value] defines the set of abstractions and methods that must be
supported by all collator-like instances.  A collator-like class is capable of
comparing and ranking two values of any type.
*/
type CollatorLike[V Value] interface {
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
	) int
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
	// Attributes
	GetClass() IteratorClassLike[V]

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
	GetClass() SorterClassLike[V]
	GetRanker() RankingFunction[V]

	// Abstractions
	Systematic[V]
}
