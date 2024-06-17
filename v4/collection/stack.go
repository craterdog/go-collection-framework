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
	age "github.com/craterdog/go-collection-framework/v4/agent"
	syn "sync"
)

// CLASS ACCESS

// Reference

var stackClass = map[string]any{}
var stackMutex syn.Mutex

// Function

func Stack[V any](notation NotationLike) StackClassLike[V] {
	// Generate the name of the bound class type.
	var class *stackClass_[V]
	var name = fmt.Sprintf("%T", class)

	// Check for existing bound class type.
	stackMutex.Lock()
	var value = stackClass[name]
	switch actual := value.(type) {
	case *stackClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Add a new bound class type.
		class = &stackClass_[V]{
			// Initialize the class constants.
			notation_:        notation,
			defaultCapacity_: 16,
		}
		stackClass[name] = class
	}
	stackMutex.Unlock()

	// Return a reference to the bound class type.
	return class
}

// CLASS METHODS

// Target

type stackClass_[V any] struct {
	// Define the class constants.
	notation_        NotationLike
	defaultCapacity_ uint
}

// Constants

func (c *stackClass_[V]) Notation() NotationLike {
	return c.notation_
}

func (c *stackClass_[V]) DefaultCapacity() uint {
	return c.defaultCapacity_
}

// Constructors

func (c *stackClass_[V]) Make() StackLike[V] {
	var list = List[V](c.notation_).Make()
	return &stack_[V]{
		class_:    c,
		capacity_: c.defaultCapacity_,
		values_:   list,
	}
}

func (c *stackClass_[V]) MakeWithCapacity(capacity uint) StackLike[V] {
	if capacity < 1 {
		panic("A stack must have a capacity greater than zero.")
	}
	var list = List[V](c.notation_).Make()
	return &stack_[V]{
		class_:    c,
		capacity_: capacity,
		values_:   list,
	}
}

func (c *stackClass_[V]) MakeFromArray(values []V) StackLike[V] {
	var list = List[V](c.notation_).MakeFromArray(values)
	return &stack_[V]{
		class_:    c,
		capacity_: c.defaultCapacity_,
		values_:   list,
	}
}

func (c *stackClass_[V]) MakeFromSequence(values Sequential[V]) StackLike[V] {
	var list = List[V](c.notation_).MakeFromSequence(values)
	return &stack_[V]{
		class_:    c,
		capacity_: c.defaultCapacity_,
		values_:   list,
	}
}

// INSTANCE METHODS

// Target

type stack_[V any] struct {
	class_    StackClassLike[V]
	capacity_ uint
	values_   ListLike[V]
}

// Attributes

func (v *stack_[V]) GetClass() StackClassLike[V] {
	return v.class_
}

func (v *stack_[V]) GetCapacity() uint {
	return v.capacity_
}

// Limited

func (v *stack_[V]) AddValue(value V) {
	if uint(v.values_.GetSize()) == v.capacity_ {
		panic("Attempted to add a value onto a stack that has reached its capacity.")
	}
	v.values_.InsertValue(0, value)
}

func (v *stack_[V]) RemoveAll() {
	v.values_.RemoveAll()
}

// Sequential

func (v *stack_[V]) IsEmpty() bool {
	return v.values_.IsEmpty()
}

func (v *stack_[V]) GetSize() int {
	return v.values_.GetSize()
}

func (v *stack_[V]) AsArray() []V {
	return v.values_.AsArray()
}

func (v *stack_[V]) GetIterator() age.IteratorLike[V] {
	return v.values_.GetIterator()
}

// Stringer

func (v *stack_[V]) String() string {
	return v.GetClass().Notation().FormatValue(v)
}

// Public

func (v *stack_[V]) RemoveTop() V {
	if v.values_.IsEmpty() {
		panic("Attempted to remove the top of an empty stack!")
	}
	return v.values_.RemoveValue(1)
}
