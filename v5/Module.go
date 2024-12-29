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

// Collection/Catalog

func CatalogExtract[K comparable, V any](
	catalog col.CatalogLike[K, V],
	keys col.Sequential[K],
) col.CatalogLike[K, V] {
	var catalogClass = col.CatalogClass[K, V]()
	return catalogClass.Extract(catalog, keys)
}

func CatalogMerge[K comparable, V any](
	first col.CatalogLike[K, V],
	second col.CatalogLike[K, V],
) col.CatalogLike[K, V] {
	var catalogClass = col.CatalogClass[K, V]()
	return catalogClass.Merge(first, second)
}

// Collection/List

func ListConcatenate[V any](
	first col.ListLike[V],
	second col.ListLike[V],
) col.ListLike[V] {
	var listClass = col.ListClass[V]()
	return listClass.Concatenate(first, second)
}

// Collection/Queue

func QueueFork[V any](
	group col.Synchronized,
	input col.QueueLike[V],
	size age.Size,
) col.Sequential[col.QueueLike[V]] {
	var queueClass = col.QueueClass[V]()
	return queueClass.Fork(group, input, size)
}

func QueueSplit[V any](
	group col.Synchronized,
	input col.QueueLike[V],
	size age.Size,
) col.Sequential[col.QueueLike[V]] {
	var queueClass = col.QueueClass[V]()
	return queueClass.Split(group, input, size)
}

func QueueJoin[V any](
	group col.Synchronized,
	inputs col.Sequential[col.QueueLike[V]],
) col.QueueLike[V] {
	var queueClass = col.QueueClass[V]()
	return queueClass.Join(group, inputs)
}

// Collection/Set

func SetAnd[V any](
	first col.SetLike[V],
	second col.SetLike[V],
) col.SetLike[V] {
	var setClass = col.SetClass[V]()
	return setClass.And(first, second)
}

func SetOr[V any](
	first col.SetLike[V],
	second col.SetLike[V],
) col.SetLike[V] {
	var setClass = col.SetClass[V]()
	return setClass.Or(first, second)
}

func SetSans[V any](
	first col.SetLike[V],
	second col.SetLike[V],
) col.SetLike[V] {
	var setClass = col.SetClass[V]()
	return setClass.Sans(first, second)
}

func SetXor[V any](
	first col.SetLike[V],
	second col.SetLike[V],
) col.SetLike[V] {
	var setClass = col.SetClass[V]()
	return setClass.Xor(first, second)
}
