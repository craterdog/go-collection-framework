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

// CLASS NAMESPACE

// This private type defines the namespace structure associated with the
// constants, constructors and functions for the Queue class namespace.
type queueClass_[V Value] struct {
	defaultCapacity int
}

// This private constant defines a map to hold all the singleton references to
// the type specific Queue class namespaces.
var queueClassSingletons = map[string]any{}

// This public function returns the singleton reference to a type specific
// Queue class namespace.  It also initializes any class constants as needed.
func Queue[V Value]() *queueClass_[V] {
	var class *queueClass_[V]
	var key = fmt.Sprintf("%T", class) // The name of the bound class type.
	var value = queueClassSingletons[key]
	switch actual := value.(type) {
	case *queueClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Create a new bound class type.
		class = &queueClass_[V]{
			defaultCapacity: 16,
		}
		queueClassSingletons[key] = class
	}
	return class
}

// CLASS CONSTANTS

// This public class constant represents the default maximum capacity for a
// Queue.
func (c *queueClass_[V]) DefaultCapacity() int {
	return c.defaultCapacity
}

// CLASS CONSTRUCTORS

// This public class constructor creates a new empty Queue with the default
// capacity which is 16.
func (c *queueClass_[V]) Empty() QueueLike[V] {
	var queue = c.WithCapacity(c.defaultCapacity)
	return queue
}

// This public class constructor creates a new Queue from the specified Go array
// of values.
func (c *queueClass_[V]) FromArray(values []V) QueueLike[V] {
	var array = Array[V]().FromArray(values)
	var queue = c.FromSequence(array)
	return queue
}

// This public class constructor creates a new Queue from the specified
// sequence of values. The Queue uses the default capacity which is 16.
func (c *queueClass_[V]) FromSequence(values Sequential[V]) QueueLike[V] {
	var queue = c.Empty()
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		queue.AddValue(value)
	}
	return queue
}

// This public class constructor creates a new empty Queue with the specified
// capacity.
func (c *queueClass_[V]) WithCapacity(capacity int) QueueLike[V] {
	if capacity < 1 {
		capacity = c.defaultCapacity
	}
	var available = make(chan bool, capacity)
	var values = List[V]().Empty()
	var queue = &queue_[V]{
		available: available,
		values:    values,
	}
	return queue
}

// CLASS FUNCTIONS

// This public class function connects the output of the specified input Queue
// with a number of new output queues specified by the size parameter and
// returns a sequence of the new output queues. Each value added to the input
// Queue will be added automatically to ALL of the output queues. This pattern
// is useful when a set of DIFFERENT operations needs to occur for every value
// and each operation can be done in parallel.
func (c *queueClass_[V]) Fork(wg *syn.WaitGroup, input FIFO[V], size int) Sequential[FIFO[V]] {
	// Validate the arguments.
	if size < 2 {
		panic("The fan out size for a Queue must be greater than one.")
	}

	// Create the new output queues.
	var capacity = input.GetCapacity()
	var outputs = List[FIFO[V]]().Empty()
	for i := 0; i < size; i++ {
		outputs.AppendValue(c.WithCapacity(capacity))
	}

	// Connect up the input Queue to the output queues in a separate go-routine.
	wg.Add(1)
	go func() {
		// Make sure the wait group is decremented on termination.
		defer wg.Done()

		// Write each value read from the input Queue to each output Queue.
		var iterator = outputs.GetIterator()
		for {
			// Read from the input Queue.
			var value, ok = input.RemoveHead() // Will block when empty.
			if !ok {
				break // The input Queue has been closed.
			}

			// Write to all output queues.
			iterator.ToStart()
			for iterator.HasNext() {
				var output = iterator.GetNext()
				output.AddValue(value) // Will block when full.
			}
		}

		// Close all output queues.
		iterator.ToStart()
		for iterator.HasNext() {
			var output = iterator.GetNext()
			output.CloseQueue()
		}
	}()

	return outputs
}

// This public class function connects the outputs of the specified sequence of
// input queues with a new output Queue returns the new output Queue. Each value
// removed from each input Queue will automatically be added to the output
// Queue. This pattern is useful when the results of the processing with a
// Split() function need to be consolidated into a single Queue.
func (c *queueClass_[V]) Join(wg *syn.WaitGroup, inputs Sequential[FIFO[V]]) FIFO[V] {
	// Validate the arguments.
	if inputs == nil || inputs.IsEmpty() {
		panic("The number of input queues for a join must be at least one.")
	}

	// Create the new output Queue.
	var iterator = inputs.GetIterator()
	var capacity = iterator.GetNext().GetCapacity()
	var output = c.WithCapacity(capacity)

	// Connect up the input queues to the output Queue.
	wg.Add(1)
	go func() {
		// Make sure the wait group is decremented on termination.
		defer wg.Done()

		// Take turns reading from each input Queue and writing to the output Queue.
		iterator.ToStart()
		for {
			var input = iterator.GetNext()
			var value, ok = input.RemoveHead() // Will block when empty.
			if !ok {
				break // The input Queue has been closed.
			}
			output.AddValue(value) // Will block when full.
			if !iterator.HasNext() {
				iterator.ToStart()
			}
		}

		// Close the output Queue.
		output.CloseQueue()
	}()

	return output
}

// This public class function connects the output of the specified input Queue
// with the number of output queues specified by the size parameter and returns
// a sequence of the new output queues. Each value added to the input Queue will
// be added automatically to ONE of the output queues. This pattern is useful
// when a SINGLE operation needs to occur for each value and the operation can
// be done on the values in parallel. The results can then be consolidated later
// on using the Join() function.
func (c *queueClass_[V]) Split(wg *syn.WaitGroup, input FIFO[V], size int) Sequential[FIFO[V]] {
	// Validate the arguments.
	if size < 2 {
		panic("The size of the split must be greater than one.")
	}

	// Create the new output queues.
	var capacity = input.GetCapacity()
	var outputs = List[FIFO[V]]().Empty()
	for i := 0; i < size; i++ {
		outputs.AppendValue(c.WithCapacity(capacity))
	}

	// Connect up the input Queue to the output queues.
	wg.Add(1)
	go func() {
		// Make sure the wait group is decremented on termination.
		defer wg.Done()

		// Take turns reading from the input Queue and writing to each output Queue.
		var iterator = outputs.GetIterator()
		for {
			// Read from the input Queue.
			var value, ok = input.RemoveHead() // Will block when empty.
			if !ok {
				break // The input Queue has been closed.
			}

			// Write to the next output Queue.
			var output = iterator.GetNext()
			output.AddValue(value) // Will block when full.
			if !iterator.HasNext() {
				iterator.ToStart()
			}
		}

		// Close all output queues.
		iterator.ToStart()
		for iterator.HasNext() {
			var output = iterator.GetNext()
			output.CloseQueue()
		}
	}()

	return outputs
}

// CLASS TYPE

// Encapsulated Type

// This private class type encapsulates a Go structure containing private
// attributes that can only be accessed and manipulated using methods that
// implement the queue-like abstract type.  A Queue implements
// first-in-first-out semantics. It is generally used by multiple go-routines at
// the same time and therefore enforces synchronized access.
// This type is parameterized as follows:
//   - V is any type of value.
//
// NOTE: If the Go "chan" type ever supports snapshots of its state, the
// underlying list can be removed and the channel modified to pass the values
// instead of the availability. Currently, the underlying list is only required
// by the "AsArray()" method.
type queue_[V Value] struct {
	available chan bool
	mutex     syn.Mutex
	values    ListLike[V]
}

// FIFO Interface

// This public class method adds the specified value to the end of this Queue.
func (v *queue_[V]) AddValue(value V) {
	v.mutex.Lock()
	v.values.AppendValue(value)
	v.mutex.Unlock()
	v.available <- true // The Queue will block if at capacity.
}

// This public class method closes the Queue so no more values can be placed on
// it.
func (v *queue_[V]) CloseQueue() {
	v.mutex.Lock()
	close(v.available) // No more values can be placed on the Queue.
	v.mutex.Unlock()
}

// This public class method retrieves the capacity of this Queue.
func (v *queue_[V]) GetCapacity() int {
	return cap(v.available) // The channel capacity is static.
}

// This public class method removes from this Queue the value that is at the
// head of it. It returns the removed value and a ", ok" value as the result.
func (v *queue_[V]) RemoveHead() (V, bool) {
	// Default the return value to the zero value for type V.
	var head V
	var ok bool

	// Remove the head value from the Queue if one exists.
	_, ok = <-v.available // Will block until a value is available.
	if ok {
		v.mutex.Lock()
		head = v.values.RemoveValue(1)
		v.mutex.Unlock()
	}

	// Return the results
	return head, ok
}

// Sequential Interface

// This public class method returns all the values in this Queue. The values
// retrieved are in the same order as they are in the Queue.
func (v *queue_[V]) AsArray() []V {
	v.mutex.Lock()
	var array = v.values.AsArray()
	v.mutex.Unlock()
	return array
}

// This public class method generates for this Queue an iterator that can be
// used to traverse its values.
func (v *queue_[V]) GetIterator() Ratcheted[V] {
	v.mutex.Lock()
	var iterator = v.values.GetIterator()
	v.mutex.Unlock()
	return iterator
}

// This public class method returns the number of values contained in this
// Queue.
func (v *queue_[V]) GetSize() int {
	v.mutex.Lock()
	var size = len(v.available)
	v.mutex.Unlock()
	return size
}

// This public class method determines whether or not this Queue is empty.
func (v *queue_[V]) IsEmpty() bool {
	v.mutex.Lock()
	var result = len(v.available) == 0
	v.mutex.Unlock()
	return result
}

// Private Interface

// This public class method is used by Go to generate a canonical string for
// the Queue.
func (v *queue_[V]) String() string {
	return Formatter().FormatCollection(v)
}
