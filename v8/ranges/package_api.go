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
Package "ranges" declares a set of classes that maintain a range of primitive
values of a generic type:
  - Interval (a finite discrete range)
  - Spectrum (an infinite discrete range)
  - Continuum (an infinite continuous range)

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
package ranges

import (
	col "github.com/craterdog/go-collection-framework/v8/collections"
)

// TYPE DECLARATIONS

/*
Bracket is a constrained type representing the inclusiveness of a bounded
collection.
*/
type Bracket uint8

const (
	Inclusive Bracket = iota
	Exclusive
)

// FUNCTIONAL DECLARATIONS

// CLASS DECLARATIONS

/*
ContinuumClassLike[V Continuous] is a class interface that declares the
complete set of class constructors, constants and functions that must be
supported by each concrete continuum-like class.

A continuum-like class defines two endpoints for an infinite continuous
sequence of elements.  The endpoints may be inclusive (denoted by a square
bracket) or exclusive (denoted by a round bracket).
*/
type ContinuumClassLike[V Continuous] interface {
	// Constructor Methods
	Continuum(
		left Bracket,
		minimum V,
		maximum V,
		right Bracket,
	) ContinuumLike[V]
}

/*
IntervalClassLike[V Discrete] is a class interface that declares the
complete set of class constructors, constants and functions that must be
supported by each concrete interval-like class.

An interval-like class defines two endpoints for a finite discrete sequence of
elements.  The endpoints may be inclusive (denoted by a square bracket) or
exclusive (denoted by a round bracket).
*/
type IntervalClassLike[V Discrete] interface {
	// Constructor Methods
	Interval(
		left Bracket,
		minimum V,
		maximum V,
		right Bracket,
	) IntervalLike[V]
}

/*
SpectrumClassLike[V Ordered[V]] is a class interface that
declares the complete set of class constructors, constants and functions that
must be supported by each concrete spectrum-like class.

A spectrum-like class defines two endpoints for an infinite discrete sequence
of elements.  The endpoints may be inclusive (denoted by a square bracket) or
exclusive (denoted by a round bracket).
*/
type SpectrumClassLike[V Ordered[V]] interface {
	// Constructor Methods
	Spectrum(
		left Bracket,
		minimum V,
		maximum V,
		right Bracket,
	) SpectrumLike[V]
}

// INSTANCE DECLARATIONS

/*
ContinuumLike[V Continuous] is an instance interface that declares the complete
set of principal, attribute and aspect methods that must be supported by each
instance of a concrete continuum-like class.
*/
type ContinuumLike[V Continuous] interface {
	// Principal Methods
	GetClass() ContinuumClassLike[V]

	// Aspect Interfaces
	Bounded[V]
	col.Searchable[V]
}

/*
IntervalLike[V Discrete] is an instance interface that declares the complete set
of principal, attribute and aspect methods that must be supported by each
instance of a concrete interval-like class.
*/
type IntervalLike[V Discrete] interface {
	// Principal Methods
	GetClass() IntervalClassLike[V]

	// Aspect Interfaces
	col.Accessible[V]
	Bounded[V]
	col.Searchable[V]
	col.Sequential[V]
}

/*
SpectrumLike[V Ordered[V]] is an instance interface that declares the complete
set of principal, attribute and aspect methods that must be supported by each
instance of a concrete spectrum-like class.
*/
type SpectrumLike[V Ordered[V]] interface {
	// Principal Methods
	GetClass() SpectrumClassLike[V]

	// Aspect Interfaces
	Bounded[V]
	col.Searchable[V]
}

// ASPECT DECLARATIONS

/*
Bounded[V any] is an aspect interface that declares a set of method signatures
that must be supported by each instance of an bounded concrete class.

A bounded class maintains the endpoints for a sequence of generic typed
primitive components.
*/
type Bounded[V any] interface {
	GetLeft() Bracket
	SetLeft(
		bracket Bracket,
	)
	GetMinimum() V
	SetMinimum(
		minimum V,
	)
	GetMaximum() V
	SetMaximum(
		minimum V,
	)
	GetRight() Bracket
	SetRight(
		bracket Bracket,
	)
}

/*
Continuous is an aspect interface that defines a set of method signatures that
must be supported by each instance of a continuous class.
*/
type Continuous interface {
	AsFloat() float64
	AsSource() string
	IsDefined() bool
}

/*
Discrete is an aspect interface that defines a set of method signatures that
must be supported by each instance of a discrete class.
*/
type Discrete interface {
	AsInteger() int
	AsSource() string
	IsDefined() bool
}

/*
Ordered[V any] is an aspect interface that declares a set of method signatures
that must be supported by each instance of an ordered class.
*/
type Ordered[V any] interface {
	AsSource() string
	IsBefore(
		value V,
	) bool
}
