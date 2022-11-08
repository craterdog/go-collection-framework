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
	col "github.com/craterdog/go-collection-framework/collections"
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

	// Add values to the queue in bulk.
	for i := 1; i < 10; i++ {
		queue.AddValue(i)
	}
	ass.Equal(t, 9, queue.GetSize())

	// Remove values from the queue in the background.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i < 101; i++ {
			var value, _ = queue.RemoveHead()
			ass.Equal(t, i, value)
		}
	}()

	// Add more values to the queue.
	for i := 10; i < 101; i++ {
		queue.AddValue(i)
	}
}

func TestQueueWithFork(t *tes.T) {
	var queues = col.Queues[int]()

	// Create a wait group for synchronization.
	var wg = new(syn.WaitGroup)
	defer wg.Wait()

	// Create a new queue with a fan out of two.
	var input = col.QueueWithCapacity[int](3)
	var outputs = queues.Fork(wg, input, 2)

	// Remove values from the output queues in the background.
	var readOutput = func(output col.FIFO[int], name string) {
		defer wg.Done()
		var value, ok = output.RemoveHead()
		for i := 1; ok; i++ {
			ass.Equal(t, i, value)
			value, ok = output.RemoveHead()
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
	var queues = col.Queues[int]()

	// Create a wait group for synchronization.
	var wg = new(syn.WaitGroup)
	defer wg.Wait()

	// Create a new queue with a split of five outputs and a join back to one.
	var input = col.QueueWithCapacity[int](3)
	var split = queues.Split(wg, input, 5)
	var output = queues.Join(wg, split)

	// Remove values from the output queue in the background.
	wg.Add(1)
	go func() {
		defer wg.Done()
		var value, ok = output.RemoveHead()
		for i := 1; ok; i++ {
			ass.Equal(t, i, value)
			value, ok = output.RemoveHead()
		}
	}()

	// Add values to the input queue.
	for i := 1; i < 21; i++ {
		input.AddValue(i)
	}
	input.CloseQueue()
}
