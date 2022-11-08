/*******************************************************************************
 *   Copyright (c) 2009-2022 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package collections

// SEQUENTIAL INTERFACES

type Key any

type Value any

// This interface defines the methods supported by all sequences of values.
type Sequential[V Value] interface {
	IsEmpty() bool
	GetSize() int
	AsArray() []V
}

// This interface defines the methods supported by all sequences whose values can
// be indexed. The indices of an indexed sequence are ORDINAL rather than ZERO
// based (which is "SO last century"). This allows for positive indices starting
// at the beginning of the sequence, and negative indices starting at the end of
// the sequence as follows:
//
//	    1           2           3             N
//	[value 1] . [value 2] . [value 3] ... [value N]
//	   -N        -(N-1)      -(N-2)          -1
//
// Notice that because the indices are ordinal based, the positive and negative
// indices are symmetrical.
type Indexed[V Value] interface {
	SetComparer(compare ComparisonFunction)
	GetValue(index int) V
	GetValues(first int, last int) Sequential[V]
	GetIndex(value V) int
}

// This interface defines the methods supported by all searchable sequences of
// values.
type Searchable[V Value] interface {
	ContainsValue(value V) bool
	ContainsAny(values Sequential[V]) bool
	ContainsAll(values Sequential[V]) bool
}

// This interface defines the methods supported by all updatable sequences of
// values.
type Updatable[V Value] interface {
	SetValue(index int, value V)
	SetValues(index int, values Sequential[V])
}

// This interface defines the methods supported by all sequences of values that
// allow values to be added and removed.
type Flexible[V Value] interface {
	SetRanker(rank RankingFunction)
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
	ShuffleValues()
	SortValues()
	SortValuesWithRanker(ranker RankingFunction)
	ReverseValues()
}

// This interface defines the methods supported by all associative sequences
// whose values consist of key-value pair associations.
type Associative[K Key, V Value] interface {
	AddAssociation(association AssociationLike[K, V])
	AddAssociations(associations Sequential[AssociationLike[K, V]])
	GetKeys() Sequential[K]
	GetValue(key K) V
	GetValues(keys Sequential[K]) Sequential[V]
	SetValue(key K, value V)
	RemoveValue(key K) V
	RemoveValues(keys Sequential[K]) Sequential[V]
	RemoveAll()
	SortAssociations()
	SortAssociationsWithRanker(ranker RankingFunction)
	ReverseAssociations()
}

// This interface defines the methods supported by all sequences whose values
// are accessed using first-in-first-out (FIFO) semantics.
type FIFO[V Value] interface {
	GetCapacity() int
	AddValue(value V)
	AddValues(values Sequential[V])
	RemoveHead() (head V, ok bool)
	CloseQueue()
}

// This interface defines the methods supported by all sequences whose values
// are accessed using last-in-first-out (LIFO) semantics.
type LIFO[V Value] interface {
	GetCapacity() int
	AddValue(value V)
	AddValues(values Sequential[V])
	GetTop() V
	RemoveTop() V
	RemoveAll()
}

// COLLECTION INTERFACES

// This interface consolidates all the interfaces supported by array-like
// sequences.
type ArrayLike[V Value] interface {
	Sequential[V]
	Indexed[V]
	Searchable[V]
	Updatable[V]
}

// This interface defines the methods supported by all association-like types.
// An association associates a key with an setable value.
type AssociationLike[K Key, V Value] interface {
	GetKey() K
	GetValue() V
	SetValue(value V)
}

// This interface consolidates all the interfaces supported by catalog-like
// sequences.
type CatalogLike[K Key, V Value] interface {
	Sequential[AssociationLike[K, V]]
	Associative[K, V]
}

// This interface consolidates all the interfaces supported by list-like
// sequences.
type ListLike[V Value] interface {
	Sequential[V]
	Indexed[V]
	Searchable[V]
	Updatable[V]
	Malleable[V]
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
	Indexed[V]
	Searchable[V]
	Flexible[V]
}

// This interface consolidates all the interfaces supported by stack-like
// sequences.
type StackLike[V Value] interface {
	Sequential[V]
	LIFO[V]
}

// AGENT INTERFACES

// This interface defines the methods supported by all iterator-like agents that
// are capable of moving forward and backward over the values in a sequence. It
// is used to implement the GoF Iterator Pattern:
//   - https://en.wikipedia.org/wiki/Iterator_pattern
//
// An iterator locks into the slots that reside between each value in the
// sequence:
//
//	    [value 1] . [value 2] . [value 3] ... [value N]
//	  ^           ^           ^                         ^
//	slot 0      slot 1      slot 2                    slot N
//
// It moves from slot to slot and has access to the values (if they exist) on
// each side of the slot.
type IteratorLike[V Value] interface {
	GetSlot() int
	ToSlot(slot int)
	ToStart()
	ToEnd()
	HasPrevious() bool
	GetPrevious() V
	HasNext() bool
	GetNext() V
}

// This interface defines the methods supported by all collator-like agent
// types that can compare and rank two values.
type CollatorLike interface {
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

// This interface defines the methods supported by all sorter-like agents that
// can sort an array of values using a ranking function.
type SorterLike[V Value] interface {
	SortArray(array []V)
}

// This type defines the function signature for any function that can sort an
// array of values using a ranking function.
type Sort[V Value] func(array []V, rank RankingFunction)
