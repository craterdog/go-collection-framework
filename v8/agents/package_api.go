/*
................................................................................
.    Copyright (c) 2009-2025 Crater Dog Technologies.  All Rights Reserved.    .
................................................................................
.  DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               .
.                                                                              .
.  This code is free software; you can redistribute it and/or modify it under  .
.  the terms of The MIT License (MIT), as published by the Open Source         .
.  Initiative. (See https://opensource.org/license/MIT)                        .
................................................................................
*/

/*
Package "agents" declares a set of agents that operate on values that have a
generic type.  They are used by the collection classes declared in this Go
module.

For detailed documentation on this package refer to the wiki:
  - https://github.com/craterdog/go-collection-framework/wiki

This package follows the Crater Dog Technologies™ Go Coding Conventions located
here:
  - https://github.com/craterdog/go-development-tools/wiki/Coding-Conventions

Additional concrete implementations of the classes declared by this package can
be developed and used seamlessly since the interface declarations only depend on
other interfaces and intrinsic types—and the class implementations only depend
on interfaces, not on each other.
*/
package agents

import ()

// TYPE DECLARATIONS

/*
Rank is a constrained type representing the possible rankings for two values.
*/
type Rank uint8

const (
	LesserRank Rank = iota
	EqualRank
	GreaterRank
)

// FUNCTIONAL DECLARATIONS

/*
RankingFunction[V any] is a functional type that declares the signature for any
function that can determine the relative ranking of two values.
*/
type RankingFunction[V any] func(
	first V,
	second V,
) Rank

// CLASS DECLARATIONS

/*
CollatorClassLike[V any] is a class interface that declares the complete set
of class constructors, constants and functions that must be supported by each
concrete collator-like class.

A collator-like class is capable of recursively comparing and ranking two values
of any type.  An optional maximum depth may be specified that limits the depth
of the structures being collated to avoid possible infinite recursion.

The default maximum depth is 16.
*/
type CollatorClassLike[V any] interface {
	// Constructor Methods
	Collator() CollatorLike[V]
	CollatorWithMaximumDepth(
		maximumDepth uint,
	) CollatorLike[V]
}

/*
SorterClassLike[V any] is a class interface that declares the complete set
of class constructors, constants and functions that must be supported by each
concrete sorter-like class.

A sorter-like class implements a specific sorting algorithm.  It uses a ranking
function to correlate the values.  If no ranking function is specified the
values are sorted into their "natural" ordering by type of value.
*/
type SorterClassLike[V any] interface {
	// Constructor Methods
	Sorter() SorterLike[V]
	SorterWithRanker(
		ranker RankingFunction[V],
	) SorterLike[V]
}

// INSTANCE DECLARATIONS

/*
CollatorLike[V any] is an instance interface that declares the complete set of
principal, attribute and aspect methods that must be supported by each
instance of a concrete collator-like class.
*/
type CollatorLike[V any] interface {
	// Principal Methods
	GetClass() CollatorClassLike[V]
	CompareValues(
		first V,
		second V,
	) bool
	RankValues(
		first V,
		second V,
	) Rank

	// Attribute Methods
	GetMaximumDepth() uint
}

/*
SorterLike[V any] is an instance interface that declares the complete set of
principal, attribute and aspect methods that must be supported by each
instance of a concrete sorter-like class.
*/
type SorterLike[V any] interface {
	// Principal Methods
	GetClass() SorterClassLike[V]
	SortValues(
		values []V,
	)
	ReverseValues(
		values []V,
	)
	ShuffleValues(
		values []V,
	)

	// Attribute Methods
	GetRanker() RankingFunction[V]
}

// ASPECT DECLARATIONS
