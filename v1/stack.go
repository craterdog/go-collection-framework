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

// STACK IMPLEMENTATION

// This constructor creates a new empty stack with the default capacity.
// The default capacity is 16 values.
func Stack[V Value]() StackLike[V] {
	return StackWithCapacity[V](0)
}

// This constructor creates a new empty stack with the specified capacity.
func StackWithCapacity[V Value](capacity int) StackLike[V] {
	// Groom the arguments.
	if capacity < 1 {
		capacity = 16 // The default value.
	}

	// Return an empty stack.
	var values = List[V]()
	return &stack[V]{values, values, capacity}
}

// This type defines the structure and methods associated with a stack of
// values. A stack implements last-in-first-out semantics.
// This type is parameterized as follows:
//   - V is any type of value.
type stack[V Value] struct {
	// Note: The delegated methods don't see the real collection type.
	Sequential[V]
	values   ListLike[V]
	capacity int
}

// STRINGER INTERFACE

func (v *stack[V]) String() string {
	return FormatValue(v)
}

// LIFO INTERFACE

// This method retrieves the capacity of this stack.
func (v *stack[V]) GetCapacity() int {
	return v.capacity
}

// This method adds the specified value to the top of this stack.
func (v *stack[V]) AddValue(value V) {
	if v.values.GetSize() == v.capacity {
		panic(fmt.Sprintf(
			"Attempted to add a value onto a stack that has reached its capacity: %v\nvalue: %v\nstack: %v\n",
			v.capacity,
			value,
			v))
	}
	v.values.AddValue(value)
}

// This method adds the specified values to the top of this stack.
func (v *stack[V]) AddValues(values Sequential[V]) {
	var iterator = Iterator(values)
	for iterator.HasNext() {
		var value = iterator.GetNext()
		v.AddValue(value) // We must call this explicitly to get the capacity check.
	}
}

// This method retrieves from this stack the value that is on top of it.
func (v *stack[V]) GetTop() V {
	if v.values.IsEmpty() {
		panic("Attempted to retrieve the top of an empty stack!")
	}
	return v.values.GetValue(-1)
}

// This method removes from this stack the value that is on top of it.
func (v *stack[V]) RemoveTop() V {
	if v.values.IsEmpty() {
		panic("Attempted to remove the top of an empty stack!")
	}
	return v.values.RemoveValue(-1)
}

// This method removes all values from this stack.
func (v *stack[V]) RemoveAll() {
	v.values.RemoveAll()
}
