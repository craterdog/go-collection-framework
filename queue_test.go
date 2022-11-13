/*******************************************************************************
 *   Copyright (c) 2009-2022 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package collections_test

import (
	col "github.com/craterdog/go-collection-framework"
	ass "github.com/stretchr/testify/assert"
	syn "sync"
	tes "testing"
)

func TestQueueWithConcurrency(t *tes.T) {
	// Create a wait group for synchronization.
	var wg = new(syn.WaitGroup)
	defer wg.Wait()

	// Create a new queue with a specific capacity.
	var queue = col.QueueWithCapacity[int](12)
	ass.Equal(t, 12, queue.GetCapacity())
	ass.True(t, queue.IsEmpty())
	ass.Equal(t, 0, queue.GetSize())

	// Add some values to the queue.
	for i := 1; i < 10; i++ {
		queue.AddValue(i)
	}
	ass.Equal(t, 9, queue.GetSize())

	// Remove values from the queue in the background.
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

	// Add some more values to the queue.
	for i := 10; i < 101; i++ {
		queue.AddValue(i)
	}
	queue.CloseQueue()
}

func TestQueueWithFork(t *tes.T) {
	// Create a wait group for synchronization.
	var wg = new(syn.WaitGroup)
	defer wg.Wait()

	// Create a new queue with a fan out of two.
	var input col.FIFO[int] = col.QueueWithCapacity[int](3)
	var outputs = col.Fork(wg, input, 2)

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
	var iterator = col.Iterator(outputs)
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

func TestQueueWithSplitAndJoin(t *tes.T) {
	// Create a wait group for synchronization.
	var wg = new(syn.WaitGroup)
	defer wg.Wait()

	// Create a new queue with a split of five outputs and a join back to one.
	var input col.FIFO[int] = col.QueueWithCapacity[int](3)
	var split = col.Split(wg, input, 5)
	var output = col.Join(wg, split)

	// Remove values from the output queue in the background.
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

	// Add values to the input queue.
	for i := 1; i < 21; i++ {
		input.AddValue(i)
	}
	input.CloseQueue()
}
