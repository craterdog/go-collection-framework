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
	syn "sync"
)

// CLASS ACCESS

// Reference

var iteratorClass = map[string]any{}
var iteratorMutex syn.Mutex

// Function

func Iterator[V Value]() IteratorClassLike[V] {
	// Generate the name of the bound class type.
	var class IteratorClassLike[V]
	var name = fmt.Sprintf("%T", class)

	// Check for existing bound class type.
	iteratorMutex.Lock()
	var value = iteratorClass[name]
	switch actual := value.(type) {
	case *iteratorClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Add a new bound class type.
		class = &iteratorClass_[V]{
			// This class does not define any constants.
		}
		iteratorClass[name] = class
	}
	iteratorMutex.Unlock()

	// Return a reference to the bound class type.
	return class
}

// CLASS METHODS

// Target

type iteratorClass_[V Value] struct {
	// This class does not define any constants.
}

// Constructors

func (c *iteratorClass_[V]) Make(sequence Sequential[V]) IteratorLike[V] {
	var values = sequence.AsArray() // The returned Go array is immutable.
	var size = len(values)
	var iterator = &iterator_[V]{
		size:   size,
		values: values,
	}
	return iterator
}

// INSTANCE METHODS

// Target

type iterator_[V Value] struct {
	size   int // So we can safely cache the size.
	slot   int // The initial slot is zero.
	values []V // The Go array of values is immutable.
}

// Public

func (v *iterator_[V]) GetNext() V {
	var result V
	if v.slot < v.size {
		v.slot = v.slot + 1
		result = v.values[v.slot-1] // convert to ZERO based indexing
	}
	return result
}

func (v *iterator_[V]) GetPrevious() V {
	var result V
	if v.slot > 0 {
		result = v.values[v.slot-1] // convert to ZERO based indexing
		v.slot = v.slot - 1
	}
	return result
}

func (v *iterator_[V]) GetSlot() int {
	return v.slot
}

func (v *iterator_[V]) HasNext() bool {
	return v.slot < v.size
}

func (v *iterator_[V]) HasPrevious() bool {
	return v.slot > 0
}

func (v *iterator_[V]) ToEnd() {
	v.slot = v.size
}

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

func (v *iterator_[V]) ToStart() {
	v.slot = 0
}
