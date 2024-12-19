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
	fmt "fmt"
	age "github.com/craterdog/go-collection-framework/v5/agent"
	col "github.com/craterdog/go-collection-framework/v5/collection"
	uti "github.com/craterdog/go-missing-utilities/v2"
	ref "reflect"
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

// Agent

func AnyCollator[V any](
	arguments ...any,
) age.CollatorLike[V] {
	// Initialize the possible arguments.
	var maximumDepth Size

	// Process the actual arguments.
	for _, argument := range arguments {
		switch actual := argument.(type) {
		case Size:
			maximumDepth = actual
		default:
			var message = fmt.Sprintf(
				"An unknown argument type passed into the collator constructor: %T\n",
				actual,
			)
			panic(message)
		}
	}

	// Call the correct constructor.
	var class = age.CollatorClass[V]()
	var collator age.CollatorLike[V]
	switch {
	case maximumDepth > 0:
		collator = class.CollatorWithMaximumDepth(maximumDepth)
	default:
		collator = class.Collator()
	}
	return collator
}

func AnyIterator[V any](
	arguments ...any,
) age.IteratorLike[V] {
	// Initialize the possible arguments.
	var values []V

	// Process the actual arguments.
	for _, argument := range arguments {
		switch actual := argument.(type) {
		case []V:
			values = actual
		default:
			var message = fmt.Sprintf(
				"An unknown argument type passed into the iterator constructor: %T\n",
				actual,
			)
			panic(message)
		}
	}

	// Call the correct constructor.
	var class = age.IteratorClass[V]()
	var iterator age.IteratorLike[V]
	switch {
	case uti.IsDefined(values):
		iterator = class.Iterator(values)
	default:
		var message = "At least one argument is required by the iterator constructor."
		panic(message)
	}
	return iterator
}

func AnySorter[V any](
	arguments ...any,
) age.SorterLike[V] {
	// Initialize the possible arguments.
	var ranker age.RankingFunction[V]

	// Process the actual arguments.
	for _, argument := range arguments {
		switch actual := argument.(type) {
		case age.RankingFunction[V]:
			ranker = actual
		default:
			var message = fmt.Sprintf(
				"An unknown argument type passed into the sorter constructor: %T\n",
				actual,
			)
			panic(message)
		}
	}

	// Call the correct constructor.
	var class = age.SorterClass[V]()
	var sorter age.SorterLike[V]
	switch {
	case uti.IsDefined(ranker):
		sorter = class.SorterWithRanker(ranker)
	default:
		sorter = class.Sorter()
	}
	return sorter
}

// Collection

func AnyAssociation[K comparable, V any](arguments ...any) col.AssociationLike[K, V] {
	// Initialize the possible arguments.
	var key K
	var value V

	// Process the actual arguments.
	for _, argument := range arguments {
		switch actual := argument.(type) {
		case K:
			key = actual
		case V:
			value = actual
		default:
			var message = fmt.Sprintf(
				"Unknown argument type passed into the association constructor: %T\n",
				actual,
			)
			panic(message)
		}
	}

	// Call the correct constructor.
	if uti.IsUndefined(key) || uti.IsUndefined(value) {
		panic("The constructor for an association requires a key and value.")
	}
	var class = col.AssociationClass[K, V]()
	var association = class.Association(key, value)
	return association
}

func AnyArray[V any](arguments ...any) col.ArrayLike[V] {
	// Initialize the possible arguments.
	var size Size
	var values []V
	var sequence col.Sequential[V]

	// Process the actual arguments.
	for _, argument := range arguments {
		switch actual := argument.(type) {
		case int:
			size = Size(uint(actual))
		case uint:
			size = Size(actual)
		case Size:
			size = actual
		case []V:
			values = actual
		default:
			var sequenceType = ref.TypeOf((*col.Sequential[V])(nil)).Elem()
			var reflectedType = ref.TypeOf(argument)
			switch {
			case reflectedType.Implements(sequenceType):
				sequence = argument.(col.Sequential[V])
			default:
				var message = fmt.Sprintf(
					"Unknown argument type passed into the array constructor: %T\n",
					actual,
				)
				panic(message)
			}
		}
	}

	// Call the correct constructor.
	var class = col.ArrayClass[V]()
	var array col.ArrayLike[V]
	switch {
	case size > 0:
		array = class.Array(size)
	case len(values) > 0:
		array = class.ArrayFromArray(values)
	case sequence != nil:
		array = class.ArrayFromSequence(sequence)
	default:
		panic("The constructor for an array requires an argument.")
	}
	return array
}

func AnyCatalog[K comparable, V any](arguments ...any) col.CatalogLike[K, V] {
	// Initialize the possible arguments.
	var associations []col.AssociationLike[K, V]
	var mappings map[K]V
	var sequence col.Sequential[col.AssociationLike[K, V]]

	// Process the actual arguments.
	for _, argument := range arguments {
		switch actual := argument.(type) {
		case []col.AssociationLike[K, V]:
			associations = actual
		case map[K]V:
			mappings = actual
		default:
			var sequenceType = ref.TypeOf((*col.Sequential[col.AssociationLike[K, V]])(nil)).Elem()
			var reflectedType = ref.TypeOf(argument)
			switch {
			case reflectedType.Implements(sequenceType):
				sequence = argument.(col.Sequential[col.AssociationLike[K, V]])
			default:
				var message = fmt.Sprintf(
					"Unknown argument type passed into the catalog constructor: %T\n",
					actual,
				)
				panic(message)
			}
		}
	}

	// Call the correct constructor.
	var class = col.CatalogClass[K, V]()
	var catalog col.CatalogLike[K, V]
	switch {
	case len(associations) > 0:
		catalog = class.CatalogFromArray(associations)
	case len(mappings) > 0:
		catalog = class.CatalogFromMap(mappings)
	case sequence != nil:
		catalog = class.CatalogFromSequence(sequence)
	default:
		catalog = class.Catalog()
	}
	return catalog
}

func AnyList[V any](arguments ...any) col.ListLike[V] {
	// Initialize the possible arguments.
	var values []V
	var sequence col.Sequential[V]

	// Process the actual arguments.
	for _, argument := range arguments {
		switch actual := argument.(type) {
		case []V:
			values = actual
		default:
			var sequenceType = ref.TypeOf((*col.Sequential[V])(nil)).Elem()
			var reflectedType = ref.TypeOf(argument)
			switch {
			case reflectedType.Implements(sequenceType):
				sequence = argument.(col.Sequential[V])
			default:
				var message = fmt.Sprintf(
					"Unknown argument type passed into the list constructor: %T\n",
					actual,
				)
				panic(message)
			}
		}
	}

	// Call the correct constructor.
	var class = col.ListClass[V]()
	var list col.ListLike[V]
	switch {
	case len(values) > 0:
		list = class.ListFromArray(values)
	case sequence != nil:
		list = class.ListFromSequence(sequence)
	default:
		list = class.List()
	}
	return list
}

func AnyMap[K comparable, V any](arguments ...any) col.MapLike[K, V] {
	// Initialize the possible arguments.
	var associations []col.AssociationLike[K, V]
	var mappings map[K]V
	var sequence col.Sequential[col.AssociationLike[K, V]]

	// Process the actual arguments.
	for _, argument := range arguments {
		switch actual := argument.(type) {
		case []col.AssociationLike[K, V]:
			associations = actual
		case map[K]V:
			mappings = actual
		default:
			var sequenceType = ref.TypeOf((*col.Sequential[col.AssociationLike[K, V]])(nil)).Elem()
			var reflectedType = ref.TypeOf(argument)
			switch {
			case reflectedType.Implements(sequenceType):
				sequence = argument.(col.Sequential[col.AssociationLike[K, V]])
			default:
				var message = fmt.Sprintf(
					"Unknown argument type passed into the map constructor: %T\n",
					actual,
				)
				panic(message)
			}
		}
	}

	// Call the correct constructor.
	var class = col.MapClass[K, V]()
	var map_ col.MapLike[K, V]
	switch {
	case len(associations) > 0:
		map_ = class.MapFromArray(associations)
	case len(mappings) > 0:
		map_ = class.MapFromMap(mappings)
	case sequence != nil:
		map_ = class.MapFromSequence(sequence)
	default:
		map_ = class.Map()
	}
	return map_
}

func AnyQueue[V any](arguments ...any) col.QueueLike[V] {
	// Initialize the possible arguments.
	var capacity Size
	var values []V
	var sequence col.Sequential[V]

	// Process the actual arguments.
	for _, argument := range arguments {
		switch actual := argument.(type) {
		case int:
			capacity = Size(uint(actual))
		case uint:
			capacity = Size(actual)
		case Size:
			capacity = actual
		case []V:
			values = actual
		default:
			var sequenceType = ref.TypeOf((*col.Sequential[V])(nil)).Elem()
			var reflectedType = ref.TypeOf(argument)
			switch {
			case reflectedType.Implements(sequenceType):
				sequence = argument.(col.Sequential[V])
			default:
				var message = fmt.Sprintf(
					"Unknown argument type passed into the queue constructor: %T\n",
					actual,
				)
				panic(message)
			}
		}
	}

	// Call the correct constructor.
	var class = col.QueueClass[V]()
	var queue col.QueueLike[V]
	switch {
	case capacity > 0:
		queue = class.QueueWithCapacity(capacity)
	case len(values) > 0:
		queue = class.QueueFromArray(values)
	case sequence != nil:
		queue = class.QueueFromSequence(sequence)
	default:
		queue = class.Queue()
	}
	return queue
}

func AnySet[V any](arguments ...any) col.SetLike[V] {
	// Initialize the possible arguments.
	var values []V
	var sequence col.Sequential[V]
	var collator age.CollatorLike[V]

	// Process the actual arguments.
	for _, argument := range arguments {
		switch actual := argument.(type) {
		case []V:
			values = actual
		case age.CollatorLike[V]:
			collator = actual
		default:
			var sequenceType = ref.TypeOf((*col.Sequential[V])(nil)).Elem()
			var reflectedType = ref.TypeOf(argument)
			switch {
			case reflectedType.Implements(sequenceType):
				sequence = argument.(col.Sequential[V])
			default:
				var message = fmt.Sprintf(
					"Unknown argument type passed into the set constructor: %T\n",
					actual,
				)
				panic(message)
			}
		}
	}

	// Call the correct constructor.
	var class = col.SetClass[V]()
	var set col.SetLike[V]
	switch {
	case collator != nil:
		set = class.SetWithCollator(collator)
		switch {
		case len(values) > 0:
			for _, value := range values {
				set.AddValue(value)
			}
		case sequence != nil:
			set.AddValues(sequence)
		}
	case len(values) > 0:
		set = class.SetFromArray(values)
	case sequence != nil:
		set = class.SetFromSequence(sequence)
	default:
		set = class.Set()
	}
	return set
}

func AnyStack[V any](arguments ...any) col.StackLike[V] {
	// Initialize the possible arguments.
	var capacity Size
	var values []V
	var sequence col.Sequential[V]

	// Process the actual arguments.
	for _, argument := range arguments {
		switch actual := argument.(type) {
		case int:
			capacity = Size(uint(actual))
		case uint:
			capacity = Size(actual)
		case Size:
			capacity = actual
		case []V:
			values = actual
		default:
			var sequenceType = ref.TypeOf((*col.Sequential[V])(nil)).Elem()
			var reflectedType = ref.TypeOf(argument)
			switch {
			case reflectedType.Implements(sequenceType):
				sequence = argument.(col.Sequential[V])
			default:
				var message = fmt.Sprintf(
					"Unknown argument type passed into the stack constructor: %T\n",
					actual,
				)
				panic(message)
			}
		}
	}

	// Call the correct constructor.
	var class = col.StackClass[V]()
	var stack col.StackLike[V]
	switch {
	case capacity > 0:
		stack = class.StackWithCapacity(capacity)
	case len(values) > 0:
		stack = class.StackFromArray(values)
	case sequence != nil:
		stack = class.StackFromSequence(sequence)
	default:
		stack = class.Stack()
	}
	return stack
}
