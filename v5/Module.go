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
┌────────────────────────────────── WARNING ───────────────────────────────────┐
│              This "Module.go" file was automatically generated.              │
│      Updates to any part of this file—other than the Module Description      │
│             and the Global Functions sections may be overwritten.            │
└──────────────────────────────────────────────────────────────────────────────┘

Package "module" declares type aliases for the commonly used types declared in
the packages contained in this module.  It also provides a default constructor
for each commonly used class that is exported by the module.  Each constructor
delegates the actual construction process to its corresponding concrete class
declared in the corresponding package contained within this module.

For detailed documentation on this entire module refer to the wiki:
  - https://github.com/craterdog/go-collection-framework//wiki
*/
package module

import (
	age "github.com/craterdog/go-collection-framework/v5/agent"
	col "github.com/craterdog/go-collection-framework/v5/collection"
)

// TYPE ALIASES

// Agent

type (
	Rank = age.Rank
	Size = age.Size
	Slot = age.Slot
)

const (
	LesserRank  = age.LesserRank
	EqualRank   = age.EqualRank
	GreaterRank = age.GreaterRank
)

// Collection

type (
	Index = col.Index
)

type (
	Synchronized = col.Synchronized
)

// CLASS CONSTRUCTORS

// Agent/Collator

func Collator[V any]() age.CollatorLike[V] {
	return age.CollatorClass[V]().Collator()
}

func CollatorWithMaximumDepth[V any](
	maximumDepth age.Size,
) age.CollatorLike[V] {
	return age.CollatorClass[V]().CollatorWithMaximumDepth(
		maximumDepth,
	)
}

// Agent/Iterator

func Iterator[V any](
	array []V,
) age.IteratorLike[V] {
	return age.IteratorClass[V]().Iterator(
		array,
	)
}

// Agent/Sorter

func Sorter[V any]() age.SorterLike[V] {
	return age.SorterClass[V]().Sorter()
}

func SorterWithRanker[V any](
	ranker age.RankingFunction[V],
) age.SorterLike[V] {
	return age.SorterClass[V]().SorterWithRanker(
		ranker,
	)
}

// Collection/Array

func Array[V any](
	size age.Size,
) col.ArrayLike[V] {
	return col.ArrayClass[V]().Array(
		size,
	)
}

func ArrayFromArray[V any](
	values []V,
) col.ArrayLike[V] {
	return col.ArrayClass[V]().ArrayFromArray(
		values,
	)
}

func ArrayFromSequence[V any](
	values col.Sequential[V],
) col.ArrayLike[V] {
	return col.ArrayClass[V]().ArrayFromSequence(
		values,
	)
}

// Collection/Association

func Association[K comparable, V any](
	key K,
	value V,
) col.AssociationLike[K, V] {
	return col.AssociationClass[K, V]().Association(
		key,
		value,
	)
}

// Collection/Catalog

func Catalog[K comparable, V any]() col.CatalogLike[K, V] {
	return col.CatalogClass[K, V]().Catalog()
}

func CatalogFromArray[K comparable, V any](
	associations []col.AssociationLike[K, V],
) col.CatalogLike[K, V] {
	return col.CatalogClass[K, V]().CatalogFromArray(
		associations,
	)
}

func CatalogFromMap[K comparable, V any](
	associations map[K]V,
) col.CatalogLike[K, V] {
	return col.CatalogClass[K, V]().CatalogFromMap(
		associations,
	)
}

func CatalogFromSequence[K comparable, V any](
	associations col.Sequential[col.AssociationLike[K, V]],
) col.CatalogLike[K, V] {
	return col.CatalogClass[K, V]().CatalogFromSequence(
		associations,
	)
}

// Collection/List

func List[V any]() col.ListLike[V] {
	return col.ListClass[V]().List()
}

func ListFromArray[V any](
	values []V,
) col.ListLike[V] {
	return col.ListClass[V]().ListFromArray(
		values,
	)
}

func ListFromSequence[V any](
	values col.Sequential[V],
) col.ListLike[V] {
	return col.ListClass[V]().ListFromSequence(
		values,
	)
}

// Collection/Map

func Map[K comparable, V any]() col.MapLike[K, V] {
	return col.MapClass[K, V]().Map()
}

func MapFromArray[K comparable, V any](
	associations []col.AssociationLike[K, V],
) col.MapLike[K, V] {
	return col.MapClass[K, V]().MapFromArray(
		associations,
	)
}

func MapFromMap[K comparable, V any](
	associations map[K]V,
) col.MapLike[K, V] {
	return col.MapClass[K, V]().MapFromMap(
		associations,
	)
}

func MapFromSequence[K comparable, V any](
	associations col.Sequential[col.AssociationLike[K, V]],
) col.MapLike[K, V] {
	return col.MapClass[K, V]().MapFromSequence(
		associations,
	)
}

// Collection/Queue

func Queue[V any]() col.QueueLike[V] {
	return col.QueueClass[V]().Queue()
}

func QueueWithCapacity[V any](
	capacity age.Size,
) col.QueueLike[V] {
	return col.QueueClass[V]().QueueWithCapacity(
		capacity,
	)
}

func QueueFromArray[V any](
	values []V,
) col.QueueLike[V] {
	return col.QueueClass[V]().QueueFromArray(
		values,
	)
}

func QueueFromSequence[V any](
	values col.Sequential[V],
) col.QueueLike[V] {
	return col.QueueClass[V]().QueueFromSequence(
		values,
	)
}

// Collection/Set

func Set[V any]() col.SetLike[V] {
	return col.SetClass[V]().Set()
}

func SetWithCollator[V any](
	collator age.CollatorLike[V],
) col.SetLike[V] {
	return col.SetClass[V]().SetWithCollator(
		collator,
	)
}

func SetFromArray[V any](
	values []V,
) col.SetLike[V] {
	return col.SetClass[V]().SetFromArray(
		values,
	)
}

func SetFromSequence[V any](
	values col.Sequential[V],
) col.SetLike[V] {
	return col.SetClass[V]().SetFromSequence(
		values,
	)
}

// Collection/Stack

func Stack[V any]() col.StackLike[V] {
	return col.StackClass[V]().Stack()
}

func StackWithCapacity[V any](
	capacity age.Size,
) col.StackLike[V] {
	return col.StackClass[V]().StackWithCapacity(
		capacity,
	)
}

func StackFromArray[V any](
	values []V,
) col.StackLike[V] {
	return col.StackClass[V]().StackFromArray(
		values,
	)
}

func StackFromSequence[V any](
	values col.Sequential[V],
) col.StackLike[V] {
	return col.StackClass[V]().StackFromSequence(
		values,
	)
}

// GLOBAL FUNCTIONS
