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

var stackClass = map[string]any{}
var stackMutex syn.Mutex

// Function

func Stack[V Value]() StackClassLike[V] {
	// Generate the name of the bound class type.
	var class StackClassLike[V]
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
			defaultCapacity: 16,
		}
		stackClass[name] = class
	}
	stackMutex.Unlock()

	// Return a reference to the bound class type.
	return class
}

// CLASS METHODS

// Target

type stackClass_[V Value] struct {
	defaultCapacity int
}

// Constants

func (c *stackClass_[V]) DefaultCapacity() int {
	return c.defaultCapacity
}

// Constructors

func (c *stackClass_[V]) Make() StackLike[V] {
	var list = List[V]().Make()
	var stack = &stack_[V]{
		capacity: c.defaultCapacity,
		values:   list,
	}
	return stack
}

func (c *stackClass_[V]) MakeFromArray(values []V) StackLike[V] {
	var list = List[V]().MakeFromArray(values)
	var stack = &stack_[V]{
		capacity: c.defaultCapacity,
		values:   list,
	}
	return stack
}

func (c *stackClass_[V]) MakeFromSequence(values Sequential[V]) StackLike[V] {
	var list = List[V]().MakeFromSequence(values)
	var stack = &stack_[V]{
		capacity: c.defaultCapacity,
		values:   list,
	}
	return stack
}

func (c *stackClass_[V]) MakeFromSource(
	source string,
	notation NotationLike,
) StackLike[V] {
	// First we parse it as a collection of any type value.
	var collection = notation.ParseSource(source).(Sequential[Value])

	// Next we must convert each value explicitly to type V.
	var anys = collection.AsArray()
	var array = make([]V, c.defaultCapacity)
	for index, value := range anys {
		array[index] = value.(V)
	}

	// Then we can create the stack from the type V array.
	var stack = c.MakeFromArray(array)
	return stack
}

func (c *stackClass_[V]) MakeWithCapacity(capacity int) StackLike[V] {
	if capacity < 1 {
		panic("A stack must have a capacity greater than zero.")
	}
	var list = List[V]().Make()
	var stack = &stack_[V]{
		capacity: capacity,
		values:   list,
	}
	return stack
}

// INSTANCE METHODS

// Target

type stack_[V Value] struct {
	capacity int
	values   ListLike[V]
}

// Limited

func (v *stack_[V]) AddValue(value V) {
	if v.values.GetSize() == v.capacity {
		panic(fmt.Sprintf(
			"Attempted to add a value onto a stack that has reached its capacity: %v\nvalue: %v\nstack: %v\n",
			v.capacity,
			value,
			v))
	}
	v.values.InsertValue(0, value)
}

func (v *stack_[V]) GetCapacity() int {
	return v.capacity
}

func (v *stack_[V]) RemoveAll() {
	v.values.RemoveAll()
}

// Sequential

func (v *stack_[V]) AsArray() []V {
	return v.values.AsArray()
}

func (v *stack_[V]) GetIterator() IteratorLike[V] {
	return v.values.GetIterator()
}

func (v *stack_[V]) GetSize() int {
	return v.values.GetSize()
}

func (v *stack_[V]) IsEmpty() bool {
	return v.values.IsEmpty()
}

// Stringer

func (v *stack_[V]) String() string {
	var formatter = Formatter().Make()
	return formatter.FormatCollection(v)
}

// Public

func (v *stack_[V]) RemoveTop() V {
	if v.values.IsEmpty() {
		panic("Attempted to remove the top of an empty stack!")
	}
	return v.values.RemoveValue(1)
}
