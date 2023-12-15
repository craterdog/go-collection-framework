/*******************************************************************************
 *   Copyright (c) 2009-2023 Crater Dog Technologies™.  All Rights Reserved.   *
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
// constants, constructors and functions for the Stack class namespace.
type stackClass_[V Value] struct {
	defaultCapacity int
}

// This private constant defines a map to hold all the singleton references to
// the type specific Stack class namespaces.
var stackClassSingletons = map[string]any{}

// This public function returns the singleton reference to a type specific
// Stack class namespace.  It also initializes any class constants as needed.
func Stack[V Value]() *stackClass_[V] {
	var class *stackClass_[V]
	var key = fmt.Sprintf("%T", class) // The name of the bound class type.
	var value = stackClassSingletons[key]
	switch actual := value.(type) {
	case *stackClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Create a new bound class type.
		class = &stackClass_[V]{
			defaultCapacity: 16,
		}
		stackClassSingletons[key] = class
	}
	return class
}

// CLASS CONSTANTS

// This public class constant represents the default maximum capacity for a
// Stack.
func (c *stackClass_[V]) DefaultCapacity() int {
	return c.defaultCapacity
}

// CLASS CONSTRUCTORS

// This public class constructor creates a new empty Stack with the default
// capacity.
func (c *stackClass_[V]) Empty() StackLike[V] {
	var stack = c.WithCapacity(c.defaultCapacity)
	return stack
}

// This public class constructor creates a new empty Stack with the specified
// capacity.
func (c *stackClass_[V]) WithCapacity(capacity int) StackLike[V] {
	var List = List[V]()
	var values = List.Empty()
	if capacity < 1 {
		capacity = c.defaultCapacity
	}
	var stack = &stack_[V]{values, capacity}
	return stack
}

// This public class constructor creates a new Stack from the specified
// sequence. The Stack uses the default capacity.
func (c *stackClass_[V]) FromSequence(sequence Sequential[V]) StackLike[V] {
	var stack = c.Empty()
	var iterator = sequence.GetIterator()
	iterator.ToEnd()
	for iterator.HasPrevious() {
		var value = iterator.GetPrevious()
		stack.AddValue(value)
	}
	return stack
}

// CLASS TYPE

// Encapsulated Type

// This private class type encapsulates a Go structure containing private
// attributes that can only be accessed and manipulated using methods that
// implement the stack-like abstract type.  A Stack implements last-in-first-out
// semantics.
// This type is parameterized as follows:
//   - V is any type of value.
type stack_[V Value] struct {
	values   ListLike[V]
	capacity int
}

// LIFO Interface

// This public class method retrieves the capacity of this Stack.
func (v *stack_[V]) GetCapacity() int {
	return v.capacity
}

// This public class method adds the specified value to the top of this Stack.
func (v *stack_[V]) AddValue(value V) {
	if v.values.GetSize() == v.capacity {
		panic(fmt.Sprintf(
			"Attempted to add a value onto a Stack that has reached its capacity: %v\nvalue: %v\nstack: %v\n",
			v.capacity,
			value,
			v))
	}
	v.values.InsertValue(0, value)
}

// This public class method retrieves from this Stack the value that is on top of it.
func (v *stack_[V]) GetTop() V {
	if v.values.IsEmpty() {
		panic("Attempted to retrieve the top of an empty Stack!")
	}
	return v.values.GetValue(1)
}

// This public class method removes from this Stack the value that is on top of it.
func (v *stack_[V]) RemoveTop() V {
	if v.values.IsEmpty() {
		panic("Attempted to remove the top of an empty Stack!")
	}
	return v.values.RemoveValue(1)
}

// This public class method removes all values from this Stack.
func (v *stack_[V]) RemoveAll() {
	v.values.RemoveAll()
}

// Sequential Interface

// This public class method determines whether or not this Stack is empty.
func (v *stack_[V]) IsEmpty() bool {
	return v.values.IsEmpty()
}

// This public class method returns the number of values contained in this
// Stack.
func (v *stack_[V]) GetSize() int {
	return v.values.GetSize()
}

// This public class method returns all the values in this Stack. The values
// retrieved are in the same order as they are in the Stack.
func (v *stack_[V]) AsArray() []V {
	return v.values.AsArray()
}

// This public class method generates for this Stack an iterator that can be
// used to traverse its values.
func (v *stack_[V]) GetIterator() Ratcheted[V] {
	return v.values.GetIterator()
}

// Private Interface

// This public class method is used by Go to generate a canonical string for
// the Stack.
func (v *stack_[V]) String() string {
	return Formatter().FormatCollection(v)
}
