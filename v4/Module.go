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
Package "module" defines a universal constructor for each class that is exported
by this module.  Each constructor delegates the actual construction process to
one of the classes defined in a subpackage for this module.

For detailed documentation on this entire module refer to the wiki:
  - https://github.com/craterdog/go-collection-framework/wiki

This package follows the Crater Dog Technologies™ (craterdog) Go Coding
Conventions located here:
  - https://github.com/craterdog/go-model-framework/wiki

Most of the classes defined in this module utilize a notation class to handle
the parsing and formatting of instances of each class using a specific notation.
The default notation is Crater Dog Collection Notation™ (CDCN), but others like
JSON and XML could be supported as well.
*/
package module

import (
	fmt "fmt"
	age "github.com/craterdog/go-collection-framework/v4/agent"
	cdc "github.com/craterdog/go-collection-framework/v4/cdcn"
	col "github.com/craterdog/go-collection-framework/v4/collection"
	//jso "github.com/craterdog/go-collection-framework/v4/json"
	//xml "github.com/craterdog/go-collection-framework/v4/xml"
	ref "reflect"
)

// TYPE ALIASES

// Agents

type (
	NotationLike = col.NotationLike
)

/*
NOTE:
The Go language does not currently support aliases for generic types.  If it
ever does, the following lines should be uncommented:

// Collections

type (
	AssociationLike = col.AssociationLike
	ArrayLike       = col.ArrayLike
	MapLike         = col.MapLike
	ListLike        = col.ListLike
	SetLike         = col.SetLike
	CatalogLike     = col.CatalogLike
	QueueLike       = col.QueueLike
	StackLike       = col.StackLike
)
*/

// GLOBAL FUNCTIONS

// Notations

func ParseSource(source string) (value any) {
	var notation = cdc.Notation().Make()
	return notation.ParseSource(source)
}

func FormatValue(value any) (source string) {
	var notation = cdc.Notation().Make()
	return notation.FormatValue(value)
}

// Agents

func IsUndefined(value any) bool {
	var collator = age.Collator[any]().Make()
	return collator.IsUndefined(value)
}

// UNIVERSAL CONSTRUCTORS

// Notations

func CDCN(arguments ...any) NotationLike {
	if len(arguments) > 0 {
		panic("The CDCN constructor does not take any arguments.")
	}
	var notation = cdc.Notation().Make()
	return notation
}

func JSON(arguments ...any) NotationLike {
	panic("The JSON notation is not yet supported.")
	//var notation = jso.Notation().Make()
	//return notation
}

func XML(arguments ...any) NotationLike {
	panic("The XML notation is not yet supported.")
	//var notation = xml.Notation().Make()
	//return notation
}

// Collections

func Association[K comparable, V any](arguments ...any) col.AssociationLike[K, V] {
	// Initialize the possible arguments.
	var notation = CDCN()
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
			var notationType = ref.TypeOf((*col.NotationLike)(nil)).Elem()
			var reflectedType = ref.TypeOf(argument)
			switch {
			case reflectedType.Implements(notationType):
				notation = argument.(col.NotationLike)
			default:
				var message = fmt.Sprintf(
					"Unknown argument type passed into the association constructor: %T\n",
					actual,
				)
				panic(message)
			}
		}
	}

	// Call the right constructor.
	var class = col.Association[K, V](notation)
	if !ref.ValueOf(key).IsValid() || !ref.ValueOf(value).IsValid() {
		panic("The constructor for an association requires a key and value.")
	}
	var association = class.MakeWithAttributes(key, value)
	return association
}

func Array[V any](arguments ...any) col.ArrayLike[V] {
	// Initialize the possible arguments.
	var notation = CDCN()
	var size uint
	var values []V
	var sequence col.Sequential[V]
	var source string

	// Process the actual arguments.
	for _, argument := range arguments {
		switch actual := argument.(type) {
		case int:
			size = uint(actual)
		case uint:
			size = actual
		case []V:
			values = actual
		case string:
			source = actual
		default:
			var notationType = ref.TypeOf((*col.NotationLike)(nil)).Elem()
			var sequenceType = ref.TypeOf((*col.Sequential[V])(nil)).Elem()
			var reflectedType = ref.TypeOf(argument)
			switch {
			case reflectedType.Implements(notationType):
				notation = argument.(col.NotationLike)
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

	// Call the right constructor.
	var class = col.Array[V](notation)
	var array col.ArrayLike[V]
	switch {
	case size > 0:
		array = class.MakeWithSize(size)
	case len(values) > 0:
		array = class.MakeFromArray(values)
	case sequence != nil:
		array = class.MakeFromSequence(sequence)
	case len(source) > 0:
		var collection = notation.ParseSource(source).(col.Sequential[any])
		// Convert the values to their real type.
		size = uint(collection.GetSize())
		array = class.MakeWithSize(size)
		var index int = 0
		var iterator = collection.GetIterator()
		for iterator.HasNext() {
			var value = iterator.GetNext().(V)
			array.SetValue(index, value)
			index++
		}
	default:
		panic("The constructor for an array requires an argument.")
	}
	return array
}

func Catalog[K comparable, V any](arguments ...any) col.CatalogLike[K, V] {
	// Initialize the possible arguments.
	var notation = CDCN()
	var associations []col.AssociationLike[K, V]
	var mappings map[K]V
	var sequence col.Sequential[col.AssociationLike[K, V]]
	var source string

	// Process the actual arguments.
	for _, argument := range arguments {
		switch actual := argument.(type) {
		case []col.AssociationLike[K, V]:
			associations = actual
		case map[K]V:
			mappings = actual
		case string:
			source = actual
		default:
			var notationType = ref.TypeOf((*col.NotationLike)(nil)).Elem()
			var sequenceType = ref.TypeOf((*col.Sequential[col.AssociationLike[K, V]])(nil)).Elem()
			var reflectedType = ref.TypeOf(argument)
			switch {
			case reflectedType.Implements(notationType):
				notation = argument.(col.NotationLike)
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

	// Call the right constructor.
	var class = col.Catalog[K, V](notation)
	var catalog col.CatalogLike[K, V]
	switch {
	case len(associations) > 0:
		catalog = class.MakeFromArray(associations)
	case len(mappings) > 0:
		catalog = class.MakeFromMap(mappings)
	case sequence != nil:
		catalog = class.MakeFromSequence(sequence)
	case len(source) > 0:
		catalog = class.Make()
		var collection = notation.ParseSource(source).(col.Sequential[col.AssociationLike[any, any]])
		// Convert the values to their real type.
		var iterator = collection.GetIterator()
		for iterator.HasNext() {
			var association = iterator.GetNext()
			var key = association.GetKey().(K)
			var value = association.GetValue().(V)
			catalog.SetValue(key, value)
		}
	default:
		catalog = class.Make()
	}
	return catalog
}

func List[V any](arguments ...any) col.ListLike[V] {
	// Initialize the possible arguments.
	var notation = CDCN()
	var values []V
	var sequence col.Sequential[V]
	var source string

	// Process the actual arguments.
	for _, argument := range arguments {
		switch actual := argument.(type) {
		case []V:
			values = actual
		case string:
			source = actual
		default:
			var notationType = ref.TypeOf((*col.NotationLike)(nil)).Elem()
			var sequenceType = ref.TypeOf((*col.Sequential[V])(nil)).Elem()
			var reflectedType = ref.TypeOf(argument)
			switch {
			case reflectedType.Implements(notationType):
				notation = argument.(col.NotationLike)
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

	// Call the right constructor.
	var class = col.List[V](notation)
	var list col.ListLike[V]
	switch {
	case len(values) > 0:
		list = class.MakeFromArray(values)
	case sequence != nil:
		list = class.MakeFromSequence(sequence)
	case len(source) > 0:
		list = class.Make()
		var collection = notation.ParseSource(source).(col.Sequential[any])
		// Convert the values to their real type.
		var iterator = collection.GetIterator()
		for iterator.HasNext() {
			var value = iterator.GetNext().(V)
			list.AppendValue(value)
		}
	default:
		list = class.Make()
	}
	return list
}

func Map[K comparable, V any](arguments ...any) col.MapLike[K, V] {
	// Initialize the possible arguments.
	var notation = CDCN()
	var associations []col.AssociationLike[K, V]
	var mappings map[K]V
	var sequence col.Sequential[col.AssociationLike[K, V]]
	var source string

	// Process the actual arguments.
	for _, argument := range arguments {
		switch actual := argument.(type) {
		case []col.AssociationLike[K, V]:
			associations = actual
		case map[K]V:
			mappings = actual
		case string:
			source = actual
		default:
			var notationType = ref.TypeOf((*col.NotationLike)(nil)).Elem()
			var sequenceType = ref.TypeOf((*col.Sequential[col.AssociationLike[K, V]])(nil)).Elem()
			var reflectedType = ref.TypeOf(argument)
			switch {
			case reflectedType.Implements(notationType):
				notation = argument.(col.NotationLike)
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

	// Call the right constructor.
	var class = col.Map[K, V](notation)
	var map_ col.MapLike[K, V]
	switch {
	case len(associations) > 0:
		map_ = class.MakeFromArray(associations)
	case len(mappings) > 0:
		map_ = class.MakeFromMap(mappings)
	case sequence != nil:
		map_ = class.MakeFromSequence(sequence)
	case len(source) > 0:
		map_ = class.Make()
		var collection = notation.ParseSource(source).(col.Sequential[col.AssociationLike[any, any]])
		// Convert the values to their real type.
		var iterator = collection.GetIterator()
		for iterator.HasNext() {
			var association = iterator.GetNext()
			var key = association.GetKey().(K)
			var value = association.GetValue().(V)
			map_.SetValue(key, value)
		}
	default:
		map_ = class.Make()
	}
	return map_
}

func Queue[V any](arguments ...any) col.QueueLike[V] {
	// Initialize the possible arguments.
	var notation = CDCN()
	var capacity uint
	var values []V
	var sequence col.Sequential[V]
	var source string

	// Process the actual arguments.
	for _, argument := range arguments {
		switch actual := argument.(type) {
		case int:
			capacity = uint(actual)
		case uint:
			capacity = actual
		case []V:
			values = actual
		case string:
			source = actual
		default:
			var notationType = ref.TypeOf((*col.NotationLike)(nil)).Elem()
			var sequenceType = ref.TypeOf((*col.Sequential[V])(nil)).Elem()
			var reflectedType = ref.TypeOf(argument)
			switch {
			case reflectedType.Implements(notationType):
				notation = argument.(col.NotationLike)
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

	// Call the right constructor.
	var class = col.Queue[V](notation)
	var queue col.QueueLike[V]
	switch {
	case capacity > 0:
		queue = class.MakeWithCapacity(capacity)
	case len(values) > 0:
		queue = class.MakeFromArray(values)
	case sequence != nil:
		queue = class.MakeFromSequence(sequence)
	case len(source) > 0:
		queue = class.Make()
		var collection = notation.ParseSource(source).(col.Sequential[any])
		// Convert the values to their real type.
		var iterator = collection.GetIterator()
		for iterator.HasNext() {
			var value = iterator.GetNext().(V)
			queue.AddValue(value)
		}
	default:
		queue = class.Make()
	}
	return queue
}

func Set[V any](arguments ...any) col.SetLike[V] {
	// Initialize the possible arguments.
	var notation = CDCN()
	var values []V
	var sequence col.Sequential[V]
	var source string
	var collator age.CollatorLike[V]

	// Process the actual arguments.
	for _, argument := range arguments {
		switch actual := argument.(type) {
		case []V:
			values = actual
		case string:
			source = actual
		case age.CollatorLike[V]:
			collator = actual
		default:
			var notationType = ref.TypeOf((*col.NotationLike)(nil)).Elem()
			var sequenceType = ref.TypeOf((*col.Sequential[V])(nil)).Elem()
			var reflectedType = ref.TypeOf(argument)
			switch {
			case reflectedType.Implements(notationType):
				notation = argument.(col.NotationLike)
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

	// Call the right constructor.
	var class = col.Set[V](notation)
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
		case len(source) > 0:
			var collection = notation.ParseSource(source).(col.Sequential[any])
			// Convert the values to their real type.
			var iterator = collection.GetIterator()
			for iterator.HasNext() {
				var value = iterator.GetNext().(V)
				set.AddValue(value)
			}
		}
	case len(values) > 0:
		set = class.MakeFromArray(values)
	case sequence != nil:
		set = class.MakeFromSequence(sequence)
	case len(source) > 0:
		set = class.Make()
		var collection = notation.ParseSource(source).(col.Sequential[any])
		// Convert the values to their real type.
		var iterator = collection.GetIterator()
		for iterator.HasNext() {
			var value = iterator.GetNext().(V)
			set.AddValue(value)
		}
	default:
		set = class.Make()
	}
	return set
}

func Stack[V any](arguments ...any) col.StackLike[V] {
	// Initialize the possible arguments.
	var notation = CDCN()
	var capacity uint
	var values []V
	var sequence col.Sequential[V]
	var source string

	// Process the actual arguments.
	for _, argument := range arguments {
		switch actual := argument.(type) {
		case int:
			capacity = uint(actual)
		case uint:
			capacity = actual
		case []V:
			values = actual
		case string:
			source = actual
		default:
			var notationType = ref.TypeOf((*col.NotationLike)(nil)).Elem()
			var sequenceType = ref.TypeOf((*col.Sequential[V])(nil)).Elem()
			var reflectedType = ref.TypeOf(argument)
			switch {
			case reflectedType.Implements(notationType):
				notation = argument.(col.NotationLike)
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

	// Call the right constructor.
	var class = col.Stack[V](notation)
	var stack col.StackLike[V]
	switch {
	case capacity > 0:
		stack = class.MakeWithCapacity(capacity)
	case len(values) > 0:
		stack = class.MakeFromArray(values)
	case sequence != nil:
		stack = class.MakeFromSequence(sequence)
	case len(source) > 0:
		stack = class.Make()
		var collection = notation.ParseSource(source).(col.Sequential[any])
		// Convert the values to their real type.
		var iterator = collection.GetIterator()
		for iterator.HasNext() {
			var value = iterator.GetNext().(V)
			stack.AddValue(value)
		}
	default:
		stack = class.Make()
	}
	return stack
}
