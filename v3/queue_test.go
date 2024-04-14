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

package collection_test

import (
	col "github.com/craterdog/go-collection-framework/v3"
	not "github.com/craterdog/go-collection-framework/v3/cdcn"
	ass "github.com/stretchr/testify/assert"
	syn "sync"
	tes "testing"
)

func TestQueueConstructor(t *tes.T) {
	var notation = not.Notation().Make()
	var Queue = col.Queue[int64]()
	var _ = Queue.MakeFromSource("[ ](Queue)", notation)
	var _ = Queue.MakeFromSource("[1, 2, 3](Queue)", notation)
}

func TestQueueConstructors(t *tes.T) {
	var Queue = col.Queue[int]()
	var queue1 = Queue.MakeFromArray([]int{1, 2, 3})
	var queue2 = Queue.MakeFromSequence(queue1)
	ass.Equal(t, queue1.AsArray(), queue2.AsArray())
}

func TestQueueWithConcurrency(t *tes.T) {
	// Create a wait group for synchronization.
	var group = new(syn.WaitGroup)
	defer group.Wait()

	// Create a new queue with a specific capacity.
	var queue = col.Queue[int]().MakeWithCapacity(12)
	ass.Equal(t, 12, queue.GetCapacity())
	ass.True(t, queue.IsEmpty())
	ass.Equal(t, 0, queue.GetSize())

	// Add some values to the queue.
	for i := 1; i < 10; i++ {
		queue.AddValue(i)
	}
	ass.Equal(t, 9, queue.GetSize())

	// Remove values from the queue in the background.
	group.Add(1)
	go func() {
		defer group.Done()
		var value int
		var ok = true
		for i := 1; ok; i++ {
			value, ok = queue.RemoveHead()
			if ok {
				ass.Equal(t, i, value)
			}
		}
		queue.RemoveAll()
	}()

	// Add some more values to the queue.
	for i := 10; i < 101; i++ {
		queue.AddValue(i)
	}
	queue.CloseQueue()
}

func TestQueueWithFork(t *tes.T) {
	// Create a wait group for synchronization.
	var group = new(syn.WaitGroup)
	defer group.Wait()

	// Create a new queue with a fan out of two.
	var Queue = col.Queue[int]()
	var input = Queue.MakeWithCapacity(3)
	var outputs = Queue.Fork(group, input, 2)

	// Remove values from the output queues in the background.
	var readOutput = func(output col.QueueLike[int], name string) {
		defer group.Done()
		var value int
		var ok bool = true
		for i := 1; ok; i++ {
			value, ok = output.RemoveHead()
			if ok {
				ass.Equal(t, i, value)
			}
		}
	}
	group.Add(2)
	var iterator = outputs.GetIterator()
	for iterator.HasNext() {
		var output = iterator.GetNext()
		go readOutput(output, "output")
	}

	// Add values to the input queue.
	for i := 1; i < 11; i++ {
		input.AddValue(i)
	}
	input.CloseQueue()
}

func TestQueueWithInvalidFanOut(t *tes.T) {
	// Create a wait group for synchronization.
	var group = new(syn.WaitGroup)
	defer group.Wait()

	// Create a new queue with an invalid fan out.
	var Queue = col.Queue[int]()
	var input = Queue.MakeWithCapacity(3)
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The fan out size for a queue must be greater than one.", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	Queue.Fork(group, input, 1) // Should panic here.
}

func TestQueueWithSplitAndJoin(t *tes.T) {
	// Create a wait group for synchronization.
	var group = new(syn.WaitGroup)
	defer group.Wait()

	// Create a new queue with a split of five outputs and a join back to one.
	var Queue = col.Queue[int]()
	var input = Queue.MakeWithCapacity(3)
	var split = Queue.Split(group, input, 5)
	var output = Queue.Join(group, split)

	// Remove values from the output queue in the background.
	group.Add(1)
	go func() {
		defer group.Done()
		var value int
		var ok bool = true
		for i := 1; ok; i++ {
			value, ok = output.RemoveHead()
			if ok {
				ass.Equal(t, i, value)
			}
		}
	}()

	// Add values to the input queue.
	for i := 1; i < 21; i++ {
		input.AddValue(i)
	}
	input.CloseQueue()
}

func TestQueueWithInvalidSplit(t *tes.T) {
	// Create a wait group for synchronization.
	var group = new(syn.WaitGroup)
	defer group.Wait()

	// Create a new queue with an invalid fan out.
	var Queue = col.Queue[int]()
	var input = Queue.MakeWithCapacity(3)
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The size of the split must be greater than one.", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	Queue.Split(group, input, 1) // Should panic here.
}

func TestQueueWithInvalidJoin(t *tes.T) {
	// Create a wait group for synchronization.
	var group = new(syn.WaitGroup)
	defer group.Wait()

	// Create a new queue with an invalid fan out.
	var inputs = col.List[col.QueueLike[int]]().Make()
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The number of input queues for a join must be at least one.", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	var Queue = col.Queue[int]()
	Queue.Join(group, inputs) // Should panic here.
	defer group.Done()
}
