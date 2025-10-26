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

package collections

import (
	fmt "fmt"
	uti "github.com/craterdog/go-missing-utilities/v8"
	syn "sync"
)

// CLASS INTERFACE

// Access Function

func QueueClass[V any]() QueueClassLike[V] {
	return queueClass[V]()
}

// Constructor Methods

func (c *queueClass_[V]) Queue() QueueLike[V] {
	var instance = c.QueueWithCapacity(0) // Request the default capacity.
	return instance
}

func (c *queueClass_[V]) QueueWithCapacity(
	capacity uint,
) QueueLike[V] {
	if capacity < 1 {
		capacity = 16 // This is the default capacity.
	}
	var available = make(chan bool, capacity)
	var listClass = ListClass[V]()
	var values = listClass.List()
	var instance = &queue_[V]{
		// Initialize the instance attributes.
		available_: available,
		capacity_:  capacity,
		values_:    values,
	}
	return instance
}

func (c *queueClass_[V]) QueueFromArray(
	values []V,
) QueueLike[V] {
	var queue = c.Queue()
	for _, value := range values {
		queue.AddValue(value)
	}
	return queue
}

func (c *queueClass_[V]) QueueFromSequence(
	values Sequential[V],
) QueueLike[V] {
	var queue = c.Queue()
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		queue.AddValue(value)
	}
	return queue
}

// Constant Methods

// Function Methods

// Functions

func (c *queueClass_[V]) Fork(
	group Synchronized,
	input QueueLike[V],
	size uint,
) Sequential[QueueLike[V]] {
	// Validate the arguments.
	if size < 2 {
		panic("The fan out size for a queue must be greater than one.")
	}

	// Create the new output queues.
	var capacity = input.GetCapacity()
	var listClass = ListClass[QueueLike[V]]()
	var outputs = listClass.List()
	var counter uint
	for ; counter < size; counter++ {
		outputs.AppendValue(c.QueueWithCapacity(capacity))
	}

	// Connect up the input queue to the output queues in a separate go-routine.
	group.Go(func() {
		// Write each value read from the input queue to each output queue.
		var iterator = outputs.GetIterator()
		for {
			// Read from the input queue.
			var value, ok = input.RemoveFirst() // Will block when empty.
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
			output.CloseChannel()
		}
	})

	return outputs
}

func (c *queueClass_[V]) Split(
	group Synchronized,
	input QueueLike[V],
	size uint,
) Sequential[QueueLike[V]] {
	// Validate the arguments.
	if size < 2 {
		panic("The size of the split must be greater than one.")
	}

	// Create the new output queues.
	var capacity = input.GetCapacity()
	var listClass = ListClass[QueueLike[V]]()
	var outputs = listClass.List()
	var counter uint
	for ; counter < size; counter++ {
		outputs.AppendValue(c.QueueWithCapacity(capacity))
	}

	// Connect up the input queue to the output queues.
	group.Go(func() {
		// Take turns reading from the input queue and writing to each output queue.
		var iterator = outputs.GetIterator()
		for {
			// Read from the input queue.
			var value, ok = input.RemoveFirst() // Will block when empty.
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
			output.CloseChannel()
		}
	})

	return outputs
}

func (c *queueClass_[V]) Join(
	group Synchronized,
	inputs Sequential[QueueLike[V]],
) QueueLike[V] {
	// Validate the arguments.
	if !uti.IsDefined(inputs) || inputs.IsEmpty() {
		panic("The number of input queues for a join must be at least one.")
	}

	// Create the new output queue.
	var iterator = inputs.GetIterator()
	var capacity = iterator.GetNext().GetCapacity()
	var output = c.QueueWithCapacity(capacity)

	// Connect up the input queues to the output queue.
	group.Go(func() {
		// Take turns reading from each input queue and writing to the output queue.
		iterator.ToStart()
		for {
			var input = iterator.GetNext()
			var value, ok = input.RemoveFirst() // Will block when empty.
			if !ok {
				break // The input queue has been closed.
			}
			output.AddValue(value) // Will block when full.
			if !iterator.HasNext() {
				iterator.ToStart()
			}
		}

		// Close the output queue.
		output.CloseChannel()
	})

	return output
}

// INSTANCE INTERFACE

// Principal Methods

func (v *queue_[V]) GetClass() QueueClassLike[V] {
	return queueClass[V]()
}

// Attribute Methods

func (v *queue_[V]) GetCapacity() uint {
	return v.capacity_
}

// Fifo[V] Methods

func (v *queue_[V]) AddValue(
	value V,
) {
	v.mutex_.Lock()
	v.values_.AppendValue(value)
	v.mutex_.Unlock()
	v.available_ <- true // The queue will block if at capacity.
}

func (v *queue_[V]) RemoveFirst() (
	first V,
	ok bool,
) {
	// Remove the first value from the queue if one exists.
	_, ok = <-v.available_ // Will block until a value is available.
	if ok {
		v.mutex_.Lock()
		first = v.values_.RemoveValue(1)
		v.mutex_.Unlock()
	}
	return
}

func (v *queue_[V]) RemoveAll() {
	v.mutex_.Lock()
	v.available_ = make(chan bool, v.capacity_)
	var listClass = ListClass[V]()
	v.values_ = listClass.List()
	v.mutex_.Unlock()
}

func (v *queue_[V]) CloseChannel() {
	v.mutex_.Lock()
	close(v.available_)
	// No more values can be placed on the queue.
	v.mutex_.Unlock()
}

// Sequential[V] Methods

func (v *queue_[V]) IsEmpty() bool {
	v.mutex_.Lock()
	var result = len(v.available_) == 0
	v.mutex_.Unlock()
	return result
}

func (v *queue_[V]) GetSize() uint {
	v.mutex_.Lock()
	var size = uint(len(v.available_))
	v.mutex_.Unlock()
	return size
}

func (v *queue_[V]) AsArray() []V {
	v.mutex_.Lock()
	var array = v.values_.AsArray()
	v.mutex_.Unlock()
	return array
}

func (v *queue_[V]) GetIterator() uti.IteratorLike[V] {
	v.mutex_.Lock()
	var iterator = v.values_.GetIterator()
	v.mutex_.Unlock()
	return iterator
}

// PROTECTED INTERFACE

func (v *queue_[V]) String() string {
	return uti.Format(v)
}

// Private Methods

// Instance Structure

// NOTE:
// If the Go "chan" type ever supports snapshots of its state, the underlying
// list and mutex can be removed and the channel modified to pass the values
// instead of the availability. Currently, the underlying list is only required
// by the "AsArray()" instance method.
type queue_[V any] struct {
	// Declare the instance attributes.
	available_ chan bool
	capacity_  uint
	mutex_     syn.Mutex
	values_    ListLike[V]
}

// Class Structure

type queueClass_[V any] struct {
	// Declare the class constants.
}

// Class Reference

var queueMap_ = map[string]any{}
var queueMutex_ syn.Mutex

func queueClass[V any]() *queueClass_[V] {
	// Generate the name of the bound class type.
	var class *queueClass_[V]
	var name = fmt.Sprintf("%T", class)

	// Check for an existing bound class type.
	queueMutex_.Lock()
	var value = queueMap_[name]
	switch actual := value.(type) {
	case *queueClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Add a new bound class type.
		class = &queueClass_[V]{
			// Initialize the class constants.
		}
		queueMap_[name] = class
	}
	queueMutex_.Unlock()

	// Return a reference to the bound class type.
	return class
}
