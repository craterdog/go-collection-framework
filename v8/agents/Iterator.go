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

package agents

import (
	fmt "fmt"
	uti "github.com/craterdog/go-missing-utilities/v8"
	syn "sync"
)

// CLASS INTERFACE

// Access Function

func IteratorClass[V any]() IteratorClassLike[V] {
	return iteratorClass[V]()
}

// Constructor Methods

func (c *iteratorClass_[V]) Iterator(
	array []V,
) IteratorLike[V] {
	if uti.IsUndefined(array) {
		panic("The \"array\" attribute is required by this class.")
	}
	var instance = &iterator_[V]{
		// Initialize the instance attributes.
		size_:   uti.ArraySize(array),
		values_: array,
	}
	return instance
}

// Constant Methods

// Function Methods

// INSTANCE INTERFACE

// Principal Methods

func (v *iterator_[V]) GetClass() IteratorClassLike[V] {
	return iteratorClass[V]()
}

func (v *iterator_[V]) IsEmpty() bool {
	return v.size_ == 0
}

func (v *iterator_[V]) ToStart() {
	v.slot_ = 0
}

func (v *iterator_[V]) ToEnd() {
	v.slot_ = v.size_
}

func (v *iterator_[V]) HasPrevious() bool {
	return v.slot_ > 0
}

func (v *iterator_[V]) GetPrevious() V {
	var result_ V
	if v.slot_ > 0 {
		result_ = v.values_[v.slot_-1] // convert to ZERO based indexing
		v.slot_ = v.slot_ - 1
	}
	return result_
}

func (v *iterator_[V]) HasNext() bool {
	return v.slot_ < v.size_
}

func (v *iterator_[V]) GetNext() V {
	var result_ V
	if v.slot_ < v.size_ {
		v.slot_ = v.slot_ + 1
		result_ = v.values_[v.slot_-1] // convert to ZERO based indexing
	}
	return result_
}

// Attribute Methods

func (v *iterator_[V]) GetSize() uint {
	return v.size_
}

func (v *iterator_[V]) GetSlot() uint {
	return v.slot_
}

func (v *iterator_[V]) SetSlot(
	slot uint,
) {
	if slot > v.size_ {
		slot = v.size_
	}
	v.slot_ = slot
}

// PROTECTED INTERFACE

// Private Methods

// Instance Structure

type iterator_[V any] struct {
	// Declare the instance attributes.
	slot_   uint
	size_   uint
	values_ []V
}

// Class Structure

type iteratorClass_[V any] struct {
	// Declare the class constants.
}

// Class Reference

var iteratorMap_ = map[string]any{}
var iteratorMutex_ syn.Mutex

func iteratorClass[V any]() *iteratorClass_[V] {
	// Generate the name of the bound class type.
	var class *iteratorClass_[V]
	var name = fmt.Sprintf("%T", class)

	// Check for an existing bound class type.
	iteratorMutex_.Lock()
	var value = iteratorMap_[name]
	switch actual := value.(type) {
	case *iteratorClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Add a new bound class type.
		class = &iteratorClass_[V]{
			// Initialize the class constants.
		}
		iteratorMap_[name] = class
	}
	iteratorMutex_.Unlock()

	// Return a reference to the bound class type.
	return class
}
