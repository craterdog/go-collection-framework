/*******************************************************************************
 *   Copyright (c) 2009-2024 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package collections_test

import (
	col "github.com/craterdog/go-collection-framework/v3"
	ass "github.com/stretchr/testify/assert"
	syn "sync"
	tes "testing"
)

func TestQueueConstructors(t *tes.T) {
	var Queue = col.Queue[int]()
	var queue1 = Queue.FromArray([]int{1, 2, 3})
	var queue2 = Queue.FromSequence(queue1)
	ass.Equal(t, queue1.AsArray(), queue2.AsArray())
}

func TestQueueWithConcurrency(t *tes.T) {
	// Create a wait group for synchronization.
	var wg = new(syn.WaitGroup)
	defer wg.Wait()

	// Create a new Queue with a specific capacity.
	var Queue = col.Queue[int]()
	var queue = Queue.WithCapacity(12)
	ass.Equal(t, 12, queue.GetCapacity())
	ass.True(t, queue.IsEmpty())
	ass.Equal(t, 0, queue.GetSize())

	// Add some values to the Queue.
	for i := 1; i < 10; i++ {
		queue.AddValue(i)
	}
	ass.Equal(t, 9, queue.GetSize())

	// Remove values from the Queue in the background.
	wg.Add(1)
	go func() {
		defer wg.Done()
		var value int
		var ok = true
		for i := 1; ok; i++ {
			value, ok = queue.RemoveHead()
			if ok {
				ass.Equal(t, i, value)
			}
		}
	}()

	// Add some more values to the Queue.
	for i := 10; i < 101; i++ {
		queue.AddValue(i)
	}
	queue.CloseQueue()
}

func TestQueueWithFork(t *tes.T) {
	// Create a wait group for synchronization.
	var wg = new(syn.WaitGroup)
	defer wg.Wait()

	// Create a new Queue with a fan out of two.
	var Queue = col.Queue[int]()
	var input = Queue.WithCapacity(3)
	var outputs = Queue.Fork(wg, input, 2)

	// Remove values from the output queues in the background.
	var readOutput = func(output col.FIFO[int], name string) {
		defer wg.Done()
		var value int
		var ok bool = true
		for i := 1; ok; i++ {
			value, ok = output.RemoveHead()
			if ok {
				ass.Equal(t, i, value)
			}
		}
	}
	wg.Add(2)
	var iterator = outputs.GetIterator()
	for iterator.HasNext() {
		var output = iterator.GetNext()
		go readOutput(output, "output")
	}

	// Add values to the input Queue.
	for i := 1; i < 11; i++ {
		input.AddValue(i)
	}
	input.CloseQueue()
}

func TestQueueWithInvalidFanOut(t *tes.T) {
	// Create a wait group for synchronization.
	var wg = new(syn.WaitGroup)
	defer wg.Wait()

	// Create a new Queue with an invalid fan out.
	var Queue = col.Queue[int]()
	var input = Queue.WithCapacity(3)
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The fan out size for a Queue must be greater than one.", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	Queue.Fork(wg, input, 1) // Should panic here.
}

func TestQueueWithSplitAndJoin(t *tes.T) {
	// Create a wait group for synchronization.
	var wg = new(syn.WaitGroup)
	defer wg.Wait()

	// Create a new Queue with a split of five outputs and a join back to one.
	var Queue = col.Queue[int]()
	var input = Queue.WithCapacity(3)
	var split = Queue.Split(wg, input, 5)
	var output = Queue.Join(wg, split)

	// Remove values from the output Queue in the background.
	wg.Add(1)
	go func() {
		defer wg.Done()
		var value int
		var ok bool = true
		for i := 1; ok; i++ {
			value, ok = output.RemoveHead()
			if ok {
				ass.Equal(t, i, value)
			}
		}
	}()

	// Add values to the input Queue.
	for i := 1; i < 21; i++ {
		input.AddValue(i)
	}
	input.CloseQueue()
}

func TestQueueWithInvalidSplit(t *tes.T) {
	// Create a wait group for synchronization.
	var wg = new(syn.WaitGroup)
	defer wg.Wait()

	// Create a new Queue with an invalid fan out.
	var Queue = col.Queue[int]()
	var input = Queue.WithCapacity(3)
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The size of the split must be greater than one.", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	Queue.Split(wg, input, 1) // Should panic here.
}

func TestQueueWithInvalidJoin(t *tes.T) {
	// Create a wait group for synchronization.
	var wg = new(syn.WaitGroup)
	defer wg.Wait()

	// Create a new Queue with an invalid fan out.
	var List = col.List[col.FIFO[int]]()
	var inputs = List.Empty()
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The number of input queues for a join must be at least one.", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	var Queue = col.Queue[int]()
	Queue.Join(wg, inputs) // Should panic here.
	defer wg.Done()
}
