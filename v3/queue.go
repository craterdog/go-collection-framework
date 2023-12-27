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

// Private Class Namespace Type

type queueClass_[V Value] struct {
	defaultCapacity int
}

// Private Class Namespace References

var queueClass = map[string]any{}

// Public Class Namespace Access

func QueueClass[V Value]() QueueClassLike[V] {
	var class *queueClass_[V]
	var key = fmt.Sprintf("%T", class) // The name of the bound class type.
	var value = queueClass[key]
	switch actual := value.(type) {
	case *queueClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Create a new bound class type.
		class = &queueClass_[V]{
			defaultCapacity: 16,
		}
		queueClass[key] = class
	}
	return class
}

// Public Class Constants

func (c *queueClass_[V]) GetDefaultCapacity() int {
	return c.defaultCapacity
}

// Public Class Constructors

func (c *queueClass_[V]) Empty() QueueLike[V] {
	var queue = c.WithCapacity(c.defaultCapacity)
	return queue
}

func (c *queueClass_[V]) FromArray(values []V) QueueLike[V] {
	var array = ArrayClass[V]().FromArray(values)
	var queue = c.FromSequence(array)
	return queue
}

func (c *queueClass_[V]) FromSequence(values Sequential[V]) QueueLike[V] {
	var queue = c.Empty()
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		queue.AddValue(value)
	}
	return queue
}

func (c *queueClass_[V]) FromString(values string) QueueLike[V] {
	// First we parse it as a collection of any type value.
	var cdcn = CDCNClass().Default()
	var collection = cdcn.ParseCollection(values).(Sequential[Value])

	// Then we convert it to a queue of type V.
	var queue = c.Empty()
	var iterator = collection.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext().(V)
		queue.AddValue(value)
	}
	return queue
}

func (c *queueClass_[V]) WithCapacity(capacity int) QueueLike[V] {
	if capacity < 1 {
		capacity = c.defaultCapacity
	}
	var available = make(chan bool, capacity)
	var values = ListClass[V]().Empty()
	var queue = &queue_[V]{
		available: available,
		capacity:  capacity,
		values:    values,
	}
	return queue
}

// Public Class Functions

// This public class function connects the output of the specified input Queue
// with a number of new output queues specified by the size parameter and
// returns a sequence of the new output queues. Each value added to the input
// queue will be added automatically to ALL of the output queues. This pattern
// is useful when a set of DIFFERENT operations needs to occur for every value
// and each operation can be done in parallel.
func (c *queueClass_[V]) Fork(
	group Synchronized,
	input QueueLike[V],
	size int,
) Sequential[QueueLike[V]] {
	// Validate the arguments.
	if size < 2 {
		panic("The fan out size for a queue must be greater than one.")
	}

	// Create the new output queues.
	var capacity = input.GetCapacity()
	var outputs = ListClass[QueueLike[V]]().Empty()
	for i := 0; i < size; i++ {
		outputs.AppendValue(c.WithCapacity(capacity))
	}

	// Connect up the input queue to the output queues in a separate go-routine.
	group.Add(1)
	go func() {
		// Make sure the wait group is decremented on termination.
		defer group.Done()

		// Write each value read from the input queue to each output queue.
		var iterator = outputs.GetIterator()
		for {
			// Read from the input queue.
			var value, ok = input.RemoveHead() // Will block when empty.
			if !ok {
				break // The input queue has been closed.
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
// input queues with a new output queue returns the new output queue. Each value
// removed from each input queue will automatically be added to the output
// queue. This pattern is useful when the results of the processing with a
// Split() function need to be consolidated into a single queue.
func (c *queueClass_[V]) Join(
	group Synchronized,
	inputs Sequential[QueueLike[V]],
) QueueLike[V] {
	// Validate the arguments.
	if inputs == nil || inputs.IsEmpty() {
		panic("The number of input queues for a join must be at least one.")
	}

	// Create the new output queue.
	var iterator = inputs.GetIterator()
	var capacity = iterator.GetNext().GetCapacity()
	var output = c.WithCapacity(capacity)

	// Connect up the input queues to the output queue.
	group.Add(1)
	go func() {
		// Make sure the wait group is decremented on termination.
		defer group.Done()

		// Take turns reading from each input queue and writing to the output queue.
		iterator.ToStart()
		for {
			var input = iterator.GetNext()
			var value, ok = input.RemoveHead() // Will block when empty.
			if !ok {
				break // The input queue has been closed.
			}
			output.AddValue(value) // Will block when full.
			if !iterator.HasNext() {
				iterator.ToStart()
			}
		}

		// Close the output queue.
		output.CloseQueue()
	}()

	return output
}

// This public class function connects the output of the specified input Queue
// with the number of output queues specified by the size parameter and returns
// a sequence of the new output queues. Each value added to the input queue will
// be added automatically to ONE of the output queues. This pattern is useful
// when a SINGLE operation needs to occur for each value and the operation can
// be done on the values in parallel. The results can then be consolidated later
// on using the Join() function.
func (c *queueClass_[V]) Split(
	group Synchronized,
	input QueueLike[V],
	size int,
) Sequential[QueueLike[V]] {
	// Validate the arguments.
	if size < 2 {
		panic("The size of the split must be greater than one.")
	}

	// Create the new output queues.
	var capacity = input.GetCapacity()
	var outputs = ListClass[QueueLike[V]]().Empty()
	for i := 0; i < size; i++ {
		outputs.AppendValue(c.WithCapacity(capacity))
	}

	// Connect up the input queue to the output queues.
	group.Add(1)
	go func() {
		// Make sure the wait group is decremented on termination.
		defer group.Done()

		// Take turns reading from the input queue and writing to each output queue.
		var iterator = outputs.GetIterator()
		for {
			// Read from the input queue.
			var value, ok = input.RemoveHead() // Will block when empty.
			if !ok {
				break // The input queue has been closed.
			}

			// Write to the next output queue.
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

// CLASS INSTANCES

// Private Class Type Definition

type queue_[V Value] struct {
	available chan bool
	capacity  int
	mutex     syn.Mutex
	values    ListLike[V]
}

// NOTE: If the Go "chan" type ever supports snapshots of its state, the
// underlying list can be removed and the channel modified to pass the values
// instead of the availability. Currently, the underlying list is only required
// by the "AsArray()" class method.

// FIFO Interface

func (v *queue_[V]) AddValue(value V) {
	v.mutex.Lock()
	v.values.AppendValue(value)
	v.mutex.Unlock()
	v.available <- true // The queue will block if at capacity.
}

func (v *queue_[V]) CloseQueue() {
	v.mutex.Lock()
	close(v.available) // No more values can be placed on the queue.
	v.mutex.Unlock()
}

func (v *queue_[V]) GetCapacity() int {
	return v.capacity
}

func (v *queue_[V]) RemoveHead() (V, bool) {
	// Default the return value to the zero value for type V.
	var head V
	var ok bool

	// Remove the head value from the queue if one exists.
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

func (v *queue_[V]) AsArray() []V {
	v.mutex.Lock()
	var array = v.values.AsArray()
	v.mutex.Unlock()
	return array
}

func (v *queue_[V]) GetIterator() Ratcheted[V] {
	v.mutex.Lock()
	var iterator = v.values.GetIterator()
	v.mutex.Unlock()
	return iterator
}

func (v *queue_[V]) GetSize() int {
	v.mutex.Lock()
	var size = len(v.available)
	v.mutex.Unlock()
	return size
}

func (v *queue_[V]) IsEmpty() bool {
	v.mutex.Lock()
	var result = len(v.available) == 0
	v.mutex.Unlock()
	return result
}

// Private Interface

// This public class method is used by Go to generate a canonical string for
// the queue.
func (v *queue_[V]) String() string {
	var cdcn = CDCNClass().Default()
	return cdcn.FormatCollection(v)
}
