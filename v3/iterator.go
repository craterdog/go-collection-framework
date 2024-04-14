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

package collection

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

func (c *iteratorClass_[V]) MakeFromSequence(values Sequential[V]) IteratorLike[V] {
	var array = values.AsArray() // The returned Go array is immutable.
	var size = len(array)
	return &iterator_[V]{
		size_:  size,
		array_: array,
	}
}

// INSTANCE METHODS

// Target

type iterator_[V Value] struct {
	size_  int // So we can safely cache the size.
	slot_  int // The initial slot is zero.
	array_ []V // The Go array of values is immutable.
}

// Public

func (v *iterator_[V]) GetNext() V {
	var result V
	if v.slot_ < v.size_ {
		v.slot_ = v.slot_ + 1
		result = v.array_[v.slot_-1] // convert to ZERO based indexing
	}
	return result
}

func (v *iterator_[V]) GetPrevious() V {
	var result V
	if v.slot_ > 0 {
		result = v.array_[v.slot_-1] // convert to ZERO based indexing
		v.slot_ = v.slot_ - 1
	}
	return result
}

func (v *iterator_[V]) GetSlot() int {
	return v.slot_
}

func (v *iterator_[V]) HasNext() bool {
	return v.slot_ < v.size_
}

func (v *iterator_[V]) HasPrevious() bool {
	return v.slot_ > 0
}

func (v *iterator_[V]) ToEnd() {
	v.slot_ = v.size_
}

func (v *iterator_[V]) ToSlot(slot int) {
	if slot > v.size_ {
		slot = v.size_
	}
	if slot < -v.size_ {
		slot = -v.size_
	}
	if slot < 0 {
		slot = slot + v.size_ + 1
	}
	v.slot_ = slot
}

func (v *iterator_[V]) ToStart() {
	v.slot_ = 0
}
