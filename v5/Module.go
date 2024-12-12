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

// DEFAULT CONSTRUCTORS

// Agent

func Collator[V any]() age.CollatorLike[V] {
	return age.CollatorClass[V]().Make()
}

func Iterator[V any](
	values []V,
) age.IteratorLike[V] {
	return age.IteratorClass[V]().Make(
		values,
	)
}

func Sorter[V any]() age.SorterLike[V] {
	return age.SorterClass[V]().Make()
}

// Collection

func Array[V any](
	size Size,
) col.ArrayLike[V] {
	return col.ArrayClass[V]().Make(
		size,
	)
}

func Association[K comparable, V any](
	key K,
	value V,
) col.AssociationLike[K, V] {
	return col.AssociationClass[K, V]().Make(
		key,
		value,
	)
}

func Catalog[K comparable, V any]() col.CatalogLike[K, V] {
	return col.CatalogClass[K, V]().Make()
}

func List[V any]() col.ListLike[V] {
	return col.ListClass[V]().Make()
}

func Map[K comparable, V any]() col.MapLike[K, V] {
	return col.MapClass[K, V]().Make()
}

func Queue[V any]() col.QueueLike[V] {
	return col.QueueClass[V]().Make()
}

func Set[V any]() col.SetLike[V] {
	return col.SetClass[V]().Make()
}

func Stack[V any]() col.StackLike[V] {
	return col.StackClass[V]().Make()
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
		collator = class.MakeWithMaximumDepth(maximumDepth)
	default:
		collator = class.Make()
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
		iterator = class.Make(values)
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
		sorter = class.MakeWithRanker(ranker)
	default:
		sorter = class.Make()
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
	var association = class.Make(key, value)
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
		array = class.Make(size)
	case len(values) > 0:
		array = class.MakeFromArray(values)
	case sequence != nil:
		array = class.MakeFromSequence(sequence)
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
		catalog = class.MakeFromArray(associations)
	case len(mappings) > 0:
		catalog = class.MakeFromMap(mappings)
	case sequence != nil:
		catalog = class.MakeFromSequence(sequence)
	default:
		catalog = class.Make()
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
		list = class.MakeFromArray(values)
	case sequence != nil:
		list = class.MakeFromSequence(sequence)
	default:
		list = class.Make()
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
		map_ = class.MakeFromArray(associations)
	case len(mappings) > 0:
		map_ = class.MakeFromMap(mappings)
	case sequence != nil:
		map_ = class.MakeFromSequence(sequence)
	default:
		map_ = class.Make()
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
		queue = class.MakeWithCapacity(capacity)
	case len(values) > 0:
		queue = class.MakeFromArray(values)
	case sequence != nil:
		queue = class.MakeFromSequence(sequence)
	default:
		queue = class.Make()
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
		set = class.MakeWithCollator(collator)
		switch {
		case len(values) > 0:
			for _, value := range values {
				set.AddValue(value)
			}
		case sequence != nil:
			set.AddValues(sequence)
		}
	case len(values) > 0:
		set = class.MakeFromArray(values)
	case sequence != nil:
		set = class.MakeFromSequence(sequence)
	default:
		set = class.Make()
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
		stack = class.MakeWithCapacity(capacity)
	case len(values) > 0:
		stack = class.MakeFromArray(values)
	case sequence != nil:
		stack = class.MakeFromSequence(sequence)
	default:
		stack = class.Make()
	}
	return stack
}
