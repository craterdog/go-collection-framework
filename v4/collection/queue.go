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

var queueClass = map[string]any{}
var queueMutex syn.Mutex

// Function

func Queue[V any](notation NotationLike) QueueClassLike[V] {
	// Validate the notation argument.
	if notation == nil {
		panic("A notation must be specified when creating this class.")
	}

	// Generate the name of the bound class type.
	var class QueueClassLike[V]
	var name = fmt.Sprintf("%T-%T", class, notation)

	// Check for existing bound class type.
	queueMutex.Lock()
	var value = queueClass[name]
	switch actual := value.(type) {
	case *queueClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Add a new bound class type.
		class = &queueClass_[V]{
			notation_:        notation,
			defaultCapacity_: 16,
		}
		queueClass[name] = class
	}
	queueMutex.Unlock()

	// Return a reference to the bound class type.
	return class
}

// CLASS METHODS

// Target

type queueClass_[V any] struct {
	notation_        NotationLike
	defaultCapacity_ uint
}

// Constants

func (c *queueClass_[V]) Notation() NotationLike {
	return c.notation_
}

func (c *queueClass_[V]) DefaultCapacity() uint {
	return c.defaultCapacity_
}

// Constructors

func (c *queueClass_[V]) Make() QueueLike[V] {
	return c.MakeWithCapacity(c.defaultCapacity_)
}

func (c *queueClass_[V]) MakeWithCapacity(capacity uint) QueueLike[V] {
	if capacity < 1 {
		capacity = c.defaultCapacity_
	}
	var available = make(chan bool, capacity)
	var values = List[V](c.notation_).Make()
	return &queue_[V]{
		class_:     c,
		available_: available,
		capacity_:  capacity,
		values_:    values,
	}
}

func (c *queueClass_[V]) MakeFromArray(values []V) QueueLike[V] {
	var array = Array[V](c.notation_).MakeFromArray(values)
	return c.MakeFromSequence(array)
}

func (c *queueClass_[V]) MakeFromSequence(values Sequential[V]) QueueLike[V] {
	var queue = c.Make()
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		queue.AddValue(value) // This call handles the synchronization.
	}
	return queue
}

func (c *queueClass_[V]) MakeFromSource(source string) QueueLike[V] {
	// First we parse it as a collection of any type value.
	var collection = c.notation_.ParseSource(source).(Sequential[any])

	// Next we must convert each value explicitly to type V.
	var anys = collection.AsArray()
	var array = make([]V, len(anys))
	for index, value := range anys {
		array[index] = value.(V)
	}

	// Then we can create the stack from the type V array.
	return c.MakeFromArray(array)
}

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
	var outputs = List[QueueLike[V]](c.notation_).Make()
	var i uint
	for ; i < size; i++ {
		outputs.AppendValue(c.MakeWithCapacity(capacity))
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
	var outputs = List[QueueLike[V]](c.notation_).Make()
	var i uint
	for ; i < size; i++ {
		outputs.AppendValue(c.MakeWithCapacity(capacity))
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
	var output = c.MakeWithCapacity(capacity)

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

// INSTANCE METHODS

// Target

// NOTE:
// If the Go "chan" type ever supports snapshots of its state, the underlying
// list can be removed and the channel modified to pass the values instead of
// the availability. Currently, the underlying list is only required by the
// "AsArray()" instance method.
type queue_[V any] struct {
	class_     QueueClassLike[V]
	available_ chan bool
	capacity_  uint
	mutex_     syn.Mutex
	values_    ListLike[V]
}

// Attributes

func (v *queue_[V]) GetClass() QueueClassLike[V] {
	return v.class_
}

func (v *queue_[V]) GetCapacity() uint {
	return v.capacity_
}

// Limited

func (v *queue_[V]) AddValue(value V) {
	v.mutex_.Lock()
	v.values_.AppendValue(value)
	v.mutex_.Unlock()
	v.available_ <- true // The queue will block if at capacity.
}

func (v *queue_[V]) RemoveAll() {
	v.mutex_.Lock()
	v.available_ = make(chan bool, v.capacity_)
	v.values_ = List[V](v.class_.Notation()).Make()
	v.mutex_.Unlock()
}

// Sequential

func (v *queue_[V]) IsEmpty() bool {
	v.mutex_.Lock()
	var result = len(v.available_) == 0
	v.mutex_.Unlock()
	return result
}

func (v *queue_[V]) GetSize() int {
	v.mutex_.Lock()
	var size = len(v.available_)
	v.mutex_.Unlock()
	return size
}

func (v *queue_[V]) AsArray() []V {
	v.mutex_.Lock()
	var array = v.values_.AsArray()
	v.mutex_.Unlock()
	return array
}

func (v *queue_[V]) GetIterator() age.IteratorLike[V] {
	v.mutex_.Lock()
	var iterator = v.values_.GetIterator()
	v.mutex_.Unlock()
	return iterator
}

// Stringer

func (v *queue_[V]) String() string {
	var notation = v.class_.Notation()
	return notation.FormatCollection(v)
}

// Public

func (v *queue_[V]) RemoveHead() (V, bool) {
	// Default the return value to the zero value for type V.
	var head V
	var ok bool

	// Remove the head value from the queue if one exists.
	_, ok = <-v.available_ // Will block until a value is available.
	if ok {
		v.mutex_.Lock()
		head = v.values_.RemoveValue(1)
		v.mutex_.Unlock()
	}

	// Return the results
	return head, ok
}

func (v *queue_[V]) CloseQueue() {
	v.mutex_.Lock()
	close(v.available_)
	// No more values can be placed on the queue.
	v.mutex_.Unlock()
}
