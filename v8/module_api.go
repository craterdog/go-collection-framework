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
┌────────────────────────────────── WARNING ───────────────────────────────────┐
│         This "module_api.go" file was automatically generated using:         │
│            https://github.com/craterdog/go-development-tools/wiki            │
│                                                                              │
│      Updates to any part of this file—other than the Module Description      │
│             and the Global Functions sections may be overwritten.            │
└──────────────────────────────────────────────────────────────────────────────┘

Package "module" declares type aliases for the commonly used types declared in
the packages contained in this module.  It also provides constructors for each
commonly used class that is exported by the module.  Each constructor delegates
the actual construction process to its corresponding concrete class declared in
the corresponding package contained within this module.

For detailed documentation on this entire module refer to the wiki:
  - https://github.com/craterdog/go-collection-framework/wiki
*/
package module

import (
	age "github.com/craterdog/go-collection-framework/v8/agents"
	col "github.com/craterdog/go-collection-framework/v8/collections"
)

// TYPE ALIASES

// Agents

type (
	Event = age.Event
	Rank  = age.Rank
)

const (
	LesserRank  = age.LesserRank
	EqualRank   = age.EqualRank
	GreaterRank = age.GreaterRank
)

type (
	RankingFunction[V any] = age.RankingFunction[V]
)

type (
	CollatorClassLike[V any] = age.CollatorClassLike[V]
	IteratorClassLike[V any] = age.IteratorClassLike[V]
	SorterClassLike[V any]   = age.SorterClassLike[V]
)

type (
	CollatorLike[V any] = age.CollatorLike[V]
	IteratorLike[V any] = age.IteratorLike[V]
	SorterLike[V any]   = age.SorterLike[V]
)

// Collections

type (
	AssociationClassLike[K comparable, V any] = col.AssociationClassLike[K, V]
	CatalogClassLike[K comparable, V any]     = col.CatalogClassLike[K, V]
	ListClassLike[V any]                      = col.ListClassLike[V]
	QueueClassLike[V any]                     = col.QueueClassLike[V]
	SetClassLike[V any]                       = col.SetClassLike[V]
	StackClassLike[V any]                     = col.StackClassLike[V]
)

type (
	AssociationLike[K comparable, V any] = col.AssociationLike[K, V]
	CatalogLike[K comparable, V any]     = col.CatalogLike[K, V]
	ListLike[V any]                      = col.ListLike[V]
	QueueLike[V any]                     = col.QueueLike[V]
	SetLike[V any]                       = col.SetLike[V]
	StackLike[V any]                     = col.StackLike[V]
)

type (
	Accessible[V any]                = col.Accessible[V]
	Associative[K comparable, V any] = col.Associative[K, V]
	Elastic[V any]                   = col.Elastic[V]
	Fifo[V any]                      = col.Fifo[V]
	Lifo[V any]                      = col.Lifo[V]
	Malleable[V any]                 = col.Malleable[V]
	Searchable[V any]                = col.Searchable[V]
	Sequential[V any]                = col.Sequential[V]
	Sortable[V any]                  = col.Sortable[V]
	Synchronized                     = col.Synchronized
	Updatable[V any]                 = col.Updatable[V]
)

// CLASS ACCESSORS

// Agents

func CollatorClass[V any]() CollatorClassLike[V] {
	return age.CollatorClass[V]()
}

func Collator[V any]() CollatorLike[V] {
	return CollatorClass[V]().Collator()
}

func CollatorWithMaximumDepth[V any](
	maximumDepth uint,
) CollatorLike[V] {
	return CollatorClass[V]().CollatorWithMaximumDepth(
		maximumDepth,
	)
}

func IteratorClass[V any]() IteratorClassLike[V] {
	return age.IteratorClass[V]()
}

func Iterator[V any](
	array []V,
) IteratorLike[V] {
	return IteratorClass[V]().Iterator(
		array,
	)
}

func SorterClass[V any]() SorterClassLike[V] {
	return age.SorterClass[V]()
}

func Sorter[V any]() SorterLike[V] {
	return SorterClass[V]().Sorter()
}

func SorterWithRanker[V any](
	ranker age.RankingFunction[V],
) SorterLike[V] {
	return SorterClass[V]().SorterWithRanker(
		ranker,
	)
}

// Collections

func AssociationClass[K comparable, V any]() AssociationClassLike[K, V] {
	return col.AssociationClass[K, V]()
}

func Association[K comparable, V any](
	key K,
	value V,
) AssociationLike[K, V] {
	return AssociationClass[K, V]().Association(
		key,
		value,
	)
}

func CatalogClass[K comparable, V any]() CatalogClassLike[K, V] {
	return col.CatalogClass[K, V]()
}

func Catalog[K comparable, V any]() CatalogLike[K, V] {
	return CatalogClass[K, V]().Catalog()
}

func CatalogFromArray[K comparable, V any](
	associations []col.AssociationLike[K, V],
) CatalogLike[K, V] {
	return CatalogClass[K, V]().CatalogFromArray(
		associations,
	)
}

func CatalogFromMap[K comparable, V any](
	associations map[K]V,
) CatalogLike[K, V] {
	return CatalogClass[K, V]().CatalogFromMap(
		associations,
	)
}

func CatalogFromSequence[K comparable, V any](
	associations col.Sequential[col.AssociationLike[K, V]],
) CatalogLike[K, V] {
	return CatalogClass[K, V]().CatalogFromSequence(
		associations,
	)
}

func ListClass[V any]() ListClassLike[V] {
	return col.ListClass[V]()
}

func List[V any]() ListLike[V] {
	return ListClass[V]().List()
}

func ListFromArray[V any](
	values []V,
) ListLike[V] {
	return ListClass[V]().ListFromArray(
		values,
	)
}

func ListFromSequence[V any](
	values col.Sequential[V],
) ListLike[V] {
	return ListClass[V]().ListFromSequence(
		values,
	)
}

func QueueClass[V any]() QueueClassLike[V] {
	return col.QueueClass[V]()
}

func Queue[V any]() QueueLike[V] {
	return QueueClass[V]().Queue()
}

func QueueWithCapacity[V any](
	capacity uint,
) QueueLike[V] {
	return QueueClass[V]().QueueWithCapacity(
		capacity,
	)
}

func QueueFromArray[V any](
	values []V,
) QueueLike[V] {
	return QueueClass[V]().QueueFromArray(
		values,
	)
}

func QueueFromSequence[V any](
	values col.Sequential[V],
) QueueLike[V] {
	return QueueClass[V]().QueueFromSequence(
		values,
	)
}

func SetClass[V any]() SetClassLike[V] {
	return col.SetClass[V]()
}

func Set[V any]() SetLike[V] {
	return SetClass[V]().Set()
}

func SetWithCollator[V any](
	collator age.CollatorLike[V],
) SetLike[V] {
	return SetClass[V]().SetWithCollator(
		collator,
	)
}

func SetFromArray[V any](
	values []V,
) SetLike[V] {
	return SetClass[V]().SetFromArray(
		values,
	)
}

func SetFromSequence[V any](
	values col.Sequential[V],
) SetLike[V] {
	return SetClass[V]().SetFromSequence(
		values,
	)
}

func StackClass[V any]() StackClassLike[V] {
	return col.StackClass[V]()
}

func Stack[V any]() StackLike[V] {
	return StackClass[V]().Stack()
}

func StackWithCapacity[V any](
	capacity uint,
) StackLike[V] {
	return StackClass[V]().StackWithCapacity(
		capacity,
	)
}

func StackFromArray[V any](
	values []V,
) StackLike[V] {
	return StackClass[V]().StackFromArray(
		values,
	)
}

func StackFromSequence[V any](
	values col.Sequential[V],
) StackLike[V] {
	return StackClass[V]().StackFromSequence(
		values,
	)
}

// GLOBAL FUNCTIONS
