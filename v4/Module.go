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
XML and JSON can be supported as well.
*/
package module

import (
	age "github.com/craterdog/go-collection-framework/v4/agent"
	not "github.com/craterdog/go-collection-framework/v4/cdcn"
	col "github.com/craterdog/go-collection-framework/v4/collection"
)

// UNIVERSAL CONSTRUCTORS

// Collections

func Array[V col.Value](arguments ...any) col.ArrayLike[V] {
	// Initialize the possible arguments.
	var notation col.NotationLike = not.Notation().Make()
	var values []V
	var sequence col.Sequential[V]
	var source string
	var size int

	// Process the actual arguments.
	for _, argument := range arguments {
		switch actual := argument.(type) {
		case col.NotationLike:
			notation = actual
		case []V:
			values = actual
		case col.Sequential[V]:
			sequence = actual
		case string:
			source = actual
		case int:
			size = actual
		}
	}

	// Call the right constructor.
	var class = col.Array[V](notation)
	var array col.ArrayLike[V]
	switch {
	case len(values) > 0:
		array = class.MakeFromArray(values)
	case sequence != nil:
		array = class.MakeFromSequence(sequence)
	case len(source) > 0:
		array = class.MakeFromSource(source)
	case size > 0:
		array = class.MakeFromSize(size)
	default:
		panic("The constructor for an array requires an argument.")
	}
	return array
}

func Association[K col.Key, V col.Value](key K, value V) col.AssociationLike[K, V] {
	var class = col.Association[K, V]()
	var association col.AssociationLike[K, V] = class.MakeWithAttributes(key, value)
	return association
}

func Catalog[K comparable, V col.Value](arguments ...any) col.CatalogLike[K, V] {
	// Initialize the possible arguments.
	var notation col.NotationLike = not.Notation().Make()
	var associations []col.AssociationLike[K, V]
	var mappings map[K]V
	var sequence col.Sequential[col.AssociationLike[K, V]]
	var source string

	// Process the actual arguments.
	for _, argument := range arguments {
		switch actual := argument.(type) {
		case col.NotationLike:
			notation = actual
		case []col.AssociationLike[K, V]:
			associations = actual
		case map[K]V:
			mappings = actual
		case col.Sequential[col.AssociationLike[K, V]]:
			sequence = actual
		case string:
			source = actual
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
		catalog = class.MakeFromSource(source)
	default:
		catalog = class.Make()
	}
	return catalog
}

func List[V col.Value](arguments ...any) col.ListLike[V] {
	// Initialize the possible arguments.
	var notation col.NotationLike = not.Notation().Make()
	var values []V
	var sequence col.Sequential[V]
	var source string

	// Process the actual arguments.
	for _, argument := range arguments {
		switch actual := argument.(type) {
		case col.NotationLike:
			notation = actual
		case []V:
			values = actual
		case col.Sequential[V]:
			sequence = actual
		case string:
			source = actual
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
		list = class.MakeFromSource(source)
	default:
		list = class.Make()
	}
	return list
}

func Map[K comparable, V col.Value](arguments ...any) col.MapLike[K, V] {
	// Initialize the possible arguments.
	var notation col.NotationLike = not.Notation().Make()
	var associations []col.AssociationLike[K, V]
	var mappings map[K]V
	var sequence col.Sequential[col.AssociationLike[K, V]]
	var source string

	// Process the actual arguments.
	for _, argument := range arguments {
		switch actual := argument.(type) {
		case col.NotationLike:
			notation = actual
		case []col.AssociationLike[K, V]:
			associations = actual
		case map[K]V:
			mappings = actual
		case col.Sequential[col.AssociationLike[K, V]]:
			sequence = actual
		case string:
			source = actual
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
		map_ = class.MakeFromSource(source)
	default:
		map_ = class.Make()
	}
	return map_
}

func Queue[V col.Value](arguments ...any) col.QueueLike[V] {
	// Initialize the possible arguments.
	var notation col.NotationLike = not.Notation().Make()
	var values []V
	var sequence col.Sequential[V]
	var source string
	var capacity int

	// Process the actual arguments.
	for _, argument := range arguments {
		switch actual := argument.(type) {
		case col.NotationLike:
			notation = actual
		case []V:
			values = actual
		case col.Sequential[V]:
			sequence = actual
		case string:
			source = actual
		case int:
			capacity = actual
		}
	}

	// Call the right constructor.
	var class = col.Queue[V](notation)
	var queue col.QueueLike[V]
	switch {
	case len(values) > 0:
		queue = class.MakeFromArray(values)
	case sequence != nil:
		queue = class.MakeFromSequence(sequence)
	case len(source) > 0:
		queue = class.MakeFromSource(source)
	case capacity > 0:
		queue = class.MakeWithCapacity(capacity)
	default:
		queue = class.Make()
	}
	return queue
}

func Set[V col.Value](arguments ...any) col.SetLike[V] {
	// Initialize the possible arguments.
	var notation col.NotationLike = not.Notation().Make()
	var values []V
	var sequence col.Sequential[V]
	var source string
	var collator age.CollatorLike[V]

	// Process the actual arguments.
	for _, argument := range arguments {
		switch actual := argument.(type) {
		case col.NotationLike:
			notation = actual
		case []V:
			values = actual
		case col.Sequential[V]:
			sequence = actual
		case string:
			source = actual
		case age.CollatorLike[V]:
			collator = actual
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
			set.AddValues(class.MakeFromSource(source))
		}
	case len(values) > 0:
		set = class.MakeFromArray(values)
	case sequence != nil:
		set = class.MakeFromSequence(sequence)
	case len(source) > 0:
		set = class.MakeFromSource(source)
	default:
		set = class.Make()
	}
	return set
}

func Stack[V col.Value](arguments ...any) col.StackLike[V] {
	// Initialize the possible arguments.
	var notation col.NotationLike = not.Notation().Make()
	var values []V
	var sequence col.Sequential[V]
	var source string
	var capacity int

	// Process the actual arguments.
	for _, argument := range arguments {
		switch actual := argument.(type) {
		case col.NotationLike:
			notation = actual
		case []V:
			values = actual
		case col.Sequential[V]:
			sequence = actual
		case string:
			source = actual
		case int:
			capacity = actual
		}
	}

	// Call the right constructor.
	var class = col.Stack[V](notation)
	var stack col.StackLike[V]
	switch {
	case len(values) > 0:
		stack = class.MakeFromArray(values)
	case sequence != nil:
		stack = class.MakeFromSequence(sequence)
	case len(source) > 0:
		stack = class.MakeFromSource(source)
	case capacity > 0:
		stack = class.MakeWithCapacity(capacity)
	default:
		stack = class.Make()
	}
	return stack
}
