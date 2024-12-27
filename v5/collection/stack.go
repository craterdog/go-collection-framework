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
	age "github.com/craterdog/go-collection-framework/v5/agent"
	uti "github.com/craterdog/go-missing-utilities/v2"
	syn "sync"
)

// CLASS INTERFACE

// Access Function

func StackClass[V any]() StackClassLike[V] {
	return stackClass[V]()
}

// Constructor Methods

func (c *stackClass_[V]) Stack() StackLike[V] {
	var listClass = ListClass[V]()
	var values = listClass.List()
	var instance = &stack_[V]{
		// Initialize the instance attributes.
		capacity_: c.defaultCapacity_,
		values_:   values,
	}
	return instance
}

func (c *stackClass_[V]) StackWithCapacity(
	capacity age.Size,
) StackLike[V] {
	if capacity < 1 {
		capacity = c.defaultCapacity_
	}
	var listClass = ListClass[V]()
	var values = listClass.List()
	var instance = &stack_[V]{
		// Initialize the instance attributes.
		capacity_: capacity,
		values_:   values,
	}
	return instance
}

func (c *stackClass_[V]) StackFromArray(
	array []V,
) StackLike[V] {
	var listClass = ListClass[V]()
	var values = listClass.ListFromArray(array)
	var instance = &stack_[V]{
		// Initialize the instance attributes.
		capacity_: c.defaultCapacity_,
		values_:   values,
	}
	return instance
}

func (c *stackClass_[V]) StackFromSequence(
	sequence Sequential[V],
) StackLike[V] {
	var listClass = ListClass[V]()
	var values = listClass.ListFromSequence(sequence)
	var instance = &stack_[V]{
		// Initialize the instance attributes.
		capacity_: c.defaultCapacity_,
		values_:   values,
	}
	return instance
}

// Constant Methods

// Function Methods

// INSTANCE INTERFACE

// Principal Methods

func (v *stack_[V]) GetClass() StackClassLike[V] {
	return stackClass[V]()
}

// Attribute Methods

func (v *stack_[V]) GetCapacity() age.Size {
	return v.capacity_
}

// Lifo[V] Methods

func (v *stack_[V]) AddValue(
	value V,
) {
	if v.values_.GetSize() == v.capacity_ {
		panic("Attempted to add a value onto a stack that has reached its capacity.")
	}
	v.values_.InsertValue(0, value)
}

func (v *stack_[V]) RemoveLast() V {
	if v.values_.IsEmpty() {
		panic("Attempted to remove a value from an empty stack!")
	}
	var last = v.values_.RemoveValue(1)
	return last
}

func (v *stack_[V]) RemoveAll() {
	v.values_.RemoveAll()
}

// Sequential[V] Methods

func (v *stack_[V]) IsEmpty() bool {
	return v.values_.IsEmpty()
}

func (v *stack_[V]) GetSize() age.Size {
	var size = v.values_.GetSize()
	return size
}

func (v *stack_[V]) AsArray() []V {
	var array = v.values_.AsArray()
	return array
}

func (v *stack_[V]) GetIterator() age.IteratorLike[V] {
	var iterator = v.values_.GetIterator()
	return iterator
}

// Stringer Methods

func (v *stack_[V]) String() string {
	return uti.Format(v)
}

// PROTECTED INTERFACE

// Private Methods

// Instance Structure

type stack_[V any] struct {
	// Declare the instance attributes.
	capacity_ age.Size
	values_   ListLike[V]
}

// Class Structure

type stackClass_[V any] struct {
	// Declare the class constants.
	defaultCapacity_ age.Size
}

// Class Reference

var stackMap_ = map[string]any{}
var stackMutex_ syn.Mutex

func stackClass[V any]() *stackClass_[V] {
	// Generate the name of the bound class type.
	var class *stackClass_[V]
	var name = fmt.Sprintf("%T", class)

	// Check for an existing bound class type.
	stackMutex_.Lock()
	var value = stackMap_[name]
	switch actual := value.(type) {
	case *stackClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Add a new bound class type.
		class = &stackClass_[V]{
			// Initialize the class constants.
			defaultCapacity_: 16,
		}
		stackMap_[name] = class
	}
	stackMutex_.Unlock()

	// Return a reference to the bound class type.
	return class
}
