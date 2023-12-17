/*******************************************************************************
 *   Copyright (c) 2009-2024 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package collections

import (
	fmt "fmt"
)

// CLASS NAMESPACE

// This private type defines the namespace structure associated with the
// constants, constructors and functions for the Iterator class namespace.
type iteratorClass_[V Value] struct {
	// This class does not define any constants.
}

// This private constant defines a map to hold all the singleton references to
// the type specific Iterator class namespaces.
var iteratorClassSingletons = map[string]any{}

// This public function returns the singleton reference to a type specific
// Iterator class namespace.  It also initializes any class constants as needed.
func Iterator[V Value]() *iteratorClass_[V] {
	var class *iteratorClass_[V]
	var key = fmt.Sprintf("%T", class) // The name of the bound class type.
	var value = iteratorClassSingletons[key]
	switch actual := value.(type) {
	case *iteratorClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Create a new bound class type.
		class = &iteratorClass_[V]{
			// This class does not define any constants.
		}
		iteratorClassSingletons[key] = class
	}
	return class
}

// CLASS CONSTRUCTORS

// This public class constructor creates a new Iterator from the specified
// sequence value.  The Iterator that can be used to traverse the values in the
// specified sequence.
func (c *iteratorClass_[V]) FromSequence(sequence Sequential[V]) IteratorLike[V] {
	var values = sequence.AsArray() // The returned array is immutable.
	var size = len(values)
	var iterator = &iterator_[V]{
		size:   size,
		values: values,
	}
	return iterator
}

// CLASS TYPE

// Encapsulated Type

// This private class type encapsulates a Go structure containing private
// attributes that can only be accessed and manipulated using methods that
// implement the iterator-like abstract type.
type iterator_[V Value] struct {
	size   int // So we can safely cache the size.
	slot   int // The initial slot is zero.
	values []V // The array of values is immutable.
}

// Ratcheted Interface

// This public class method retrieves the value after the current slot.
func (v *iterator_[V]) GetNext() V {
	var result V
	if v.slot < v.size {
		v.slot = v.slot + 1
		result = v.values[v.slot-1] // convert to ZERO based indexing
	}
	return result
}

// This public class method retrieves the value before the current slot.
func (v *iterator_[V]) GetPrevious() V {
	var result V
	if v.slot > 0 {
		result = v.values[v.slot-1] // convert to ZERO based indexing
		v.slot = v.slot - 1
	}
	return result
}

// This public class method returns the current slot between values that this
// Iterator is currently locked into.
func (v *iterator_[V]) GetSlot() int {
	return v.slot
}

// This public class method determines whether or not there is a value after the
// current slot.
func (v *iterator_[V]) HasNext() bool {
	return v.slot < v.size
}

// This public class method determines whether or not there is a value before
// the current slot.
func (v *iterator_[V]) HasPrevious() bool {
	return v.slot > 0
}

// This public class method moves this Iterator to the slot after the last
// value.
func (v *iterator_[V]) ToEnd() {
	v.slot = v.size
}

// This public class method moves this Iterator to the specified slot between
// values.
func (v *iterator_[V]) ToSlot(slot int) {
	if slot > v.size {
		slot = v.size
	}
	if slot < -v.size {
		slot = -v.size
	}
	if slot < 0 {
		slot = slot + v.size + 1
	}
	v.slot = slot
}

// This public class method moves this Iterator to the slot before the first
// value.
func (v *iterator_[V]) ToStart() {
	v.slot = 0
}
