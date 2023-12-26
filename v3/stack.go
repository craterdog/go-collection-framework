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

// Private Class Namespace Type

type stackClass_[V Value] struct {
	defaultCapacity int
}

// Private Namespace Reference(s)

var stackClass = map[string]any{}

// Public Namespace Access

func Stack[V Value]() StackClassLike[V] {
	var class *stackClass_[V]
	var key = fmt.Sprintf("%T", class) // The name of the bound class type.
	var value = stackClass[key]
	switch actual := value.(type) {
	case *stackClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Create a new bound class type.
		class = &stackClass_[V]{
			defaultCapacity: 16,
		}
		stackClass[key] = class
	}
	return class
}

// Public Class Constants

func (c *stackClass_[V]) GetDefaultCapacity() int {
	return c.defaultCapacity
}

// Public Class Constructors

func (c *stackClass_[V]) Empty() StackLike[V] {
	var stack = c.WithCapacity(c.defaultCapacity)
	return stack
}

func (c *stackClass_[V]) FromArray(values []V) StackLike[V] {
	var array = Array[V]().FromArray(values)
	var stack = c.FromSequence(array)
	return stack
}

func (c *stackClass_[V]) FromSequence(values Sequential[V]) StackLike[V] {
	var stack = c.Empty()
	var iterator = values.GetIterator()
	iterator.ToEnd() // Add the values in reverse order since it is a LIFO.
	for iterator.HasPrevious() {
		var value = iterator.GetPrevious()
		stack.AddValue(value)
	}
	return stack
}

func (c *stackClass_[V]) FromString(values string) StackLike[V] {
	// First we parse it as a collection of any type value.
	var collection = CDCN().Default().ParseCollection(values).(Sequential[Value])

	// Then we convert it to a stack of type V.
	var stack = c.Empty()
	var iterator = collection.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext().(V)
		stack.AddValue(value)
	}
	return stack
}

func (c *stackClass_[V]) WithCapacity(capacity int) StackLike[V] {
	var values = List[V]().Empty()
	if capacity < 1 {
		capacity = c.defaultCapacity
	}
	var stack = &stack_[V]{
		capacity,
		values,
	}
	return stack
}

// CLASS TYPE

// Private Class Type Definition

type stack_[V Value] struct {
	capacity int
	values   ListLike[V]
}

// LIFO Interface

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

func (v *stack_[V]) GetTop() V {
	if v.values.IsEmpty() {
		panic("Attempted to retrieve the top of an empty stack!")
	}
	return v.values.GetValue(1)
}

func (v *stack_[V]) RemoveAll() {
	v.values.RemoveAll()
}

func (v *stack_[V]) RemoveTop() V {
	if v.values.IsEmpty() {
		panic("Attempted to remove the top of an empty stack!")
	}
	return v.values.RemoveValue(1)
}

// Sequential Interface

func (v *stack_[V]) AsArray() []V {
	return v.values.AsArray()
}

func (v *stack_[V]) GetIterator() Ratcheted[V] {
	return v.values.GetIterator()
}

func (v *stack_[V]) GetSize() int {
	return v.values.GetSize()
}

func (v *stack_[V]) IsEmpty() bool {
	return v.values.IsEmpty()
}

// Private Interface

// This public class method is used by Go to generate a canonical string for
// the stack.
func (v *stack_[V]) String() string {
	return CDCN().Default().FormatCollection(v)
}
