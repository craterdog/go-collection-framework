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

package module_test

import (
	fmt "fmt"
	fra "github.com/craterdog/go-collection-framework/v5"
	col "github.com/craterdog/go-collection-framework/v5/collection"
	syn "sync"
	tes "testing"
)

func TestModuleFunctions(t *tes.T) {
	fra.Collator[any]()
	fra.CollatorWithMaximumDepth[any](8)
	fra.Iterator[any]([]any{"foo", 5})
	var sorter = fra.Sorter[any]()
	fra.SorterWithRanker[any](sorter.GetRanker())
	fra.List[string]()
	var list = fra.ListFromArray[string]([]string{"A"})
	fra.ListFromSequence[string](list)
	fra.ListConcatenate[string](list, list)
	var association = fra.Association[string, int]("A", 1)
	var catalog = fra.Catalog[string, int]()
	fra.CatalogFromArray[string, int]([]col.AssociationLike[string, int]{association})
	fra.CatalogFromMap[string, int](catalog.AsMap())
	fra.CatalogFromSequence[string, int](catalog)
	fra.CatalogExtract[string, int](catalog, list)
	fra.CatalogMerge[string, int](catalog, catalog)
	fra.Queue[string]()
	fra.QueueWithCapacity[string](8)
	var queue = fra.QueueFromArray[string](list.AsArray())
	fra.QueueFromSequence[string](queue)
	var group = new(syn.WaitGroup)
	defer group.Wait()
	var queues = fra.QueueFork[string](group, queue, 2)
	fra.QueueSplit[string](group, queue, 2)
	fra.QueueJoin[string](group, queues)
	queue.CloseChannel()
	var set = fra.Set[string]()
	fra.SetWithCollator[string](set.GetCollator())
	fra.SetFromArray[string](set.AsArray())
	fra.SetFromSequence[string](set)
	fra.SetAnd[string](set, set)
	fra.SetOr[string](set, set)
	fra.SetSans[string](set, set)
	fra.SetXor[string](set, set)
	fra.Stack[string]()
	fra.StackWithCapacity[string](8)
	fra.StackFromArray[string](list.AsArray())
	fra.StackFromSequence[string](list)
}

func TestModuleExampleCode(t *tes.T) {
	fmt.Println("MODULE EXAMPLE:")

	// Create an empty list.
	var list = fra.List[string]()
	fmt.Printf("An empty list: %v\n", list)
	fmt.Println()

	// Create a list using an intrinsic Go array of values.
	list = fra.ListFromArray[string](
		[]string{"Hello", "World"},
	)
	fmt.Printf("A list: %v\n", list)
	fmt.Println()

	// Create an empty catalog.
	var catalog = fra.Catalog[string, int64]()
	fmt.Printf("An empty catalog: %v\n", catalog)
	fmt.Println()

	// Create a catalog from an intrinsic Go map.
	catalog = fra.CatalogFromMap[string, int64](
		map[string]int64{
			"alpha": 1,
			"beta":  2,
			"gamma": 3,
		},
	)
	fmt.Printf("A catalog: %v\n", catalog)
	fmt.Println()

	// Create a list of the catalog keys.
	list = fra.ListFromSequence[string](catalog.GetKeys())
	fmt.Printf("A list of keys: %v\n", list)
	fmt.Println()

	// Create a set from an intrinsic Go array of values.
	var set = fra.SetFromArray[string](
		[]string{"a", "b", "r", "a", "c", "a", "d", "a", "b", "r", "a"},
	)
	fmt.Printf("A set: %v\n", set)
	fmt.Println()

	// Create an empty stack with a capacity of 4.
	var stack = fra.StackWithCapacity[string](4)
	fmt.Printf("An empty stack: %v\n", stack)
	fmt.Println()

	// Create a stack containing the values from a list.
	stack = fra.StackFromSequence[string](list)
	fmt.Printf("A stack: %v\n", stack)
	fmt.Println()

	// Create an empty queue with a capacity of 5.
	var queue = fra.QueueWithCapacity[string](5)
	fmt.Printf("An empty queue: %v\n", queue)
	fmt.Println()

	// Create a queue containing the values from a set.
	queue = fra.QueueFromSequence[string](set)
	fmt.Printf("A queue: %v\n", queue)
	fmt.Println()
}

func TestListExampleCode(t *tes.T) {
	fmt.Println("LIST EXAMPLE:")

	// Create a new list using the universal constructor.
	var list = fra.ListFromArray[string](
		[]string{"bar", "foo", "bax"},
	)
	fmt.Println("The initialized list:", list)

	// Add some more values to the list.
	list.AppendValue("fuz")
	list.AppendValue("box")
	fmt.Println("The augmented list:", list)

	// Change a value in the list.
	list.SetValue(2, "foz")
	fmt.Println("The updated list:", list)

	// Insert a new value at the beginning of the list (slot 0).
	list.InsertValue(0, "bax")
	fmt.Println("The updated list:", list)

	// Insert a new value at the end of the list (slot N).
	list.InsertValue(6, "bux")
	fmt.Println("The updated list:", list)

	// Sort the values in the list.
	list.SortValues()
	fmt.Println("The sorted list:", list)

	// Randomly shuffle the values in the list.
	list.ShuffleValues()
	fmt.Println("The shuffled list:", list)

	// Remove a value from the list.
	list.RemoveValue(4)
	fmt.Println("The shortened list:", list)

	// Remove all values from the list.
	list.RemoveAll()
	fmt.Println("The empty list:", list)
	fmt.Println()
}

func TestSetExampleCode(t *tes.T) {
	fmt.Println("SET EXAMPLE:")

	// Create two sets with overlapping values using the universal constructor.
	var set1 = fra.SetFromArray[string](
		[]string{"alpha", "beta", "gamma"},
	)
	fmt.Println("The first set is:", set1)
	var set2 = fra.SetFromArray[string](
		[]string{"beta", "gamma", "delta"},
	)
	fmt.Println("The second set is:", set2)

	// Find the logical union of the two sets.
	var Set = set1.GetClass()
	fmt.Println("The union of the two sets is:", Set.Or(set1, set2))

	// Find the logical difference between the two sets.
	fmt.Println("The first set minus the second set is:", Set.Sans(set1, set2))

	// Find the logical intersection of the two sets.
	fmt.Println("The intersection of the two sets is:", Set.And(set1, set2))

	// Find the logical exclusion of the two sets.
	fmt.Println("The exclusion of the two sets is:", Set.Xor(set1, set2))

	// Add an existing value to the first set.
	set1.AddValue("beta")
	fmt.Println("Adding an existing value to a set does not change it:", set1)
	fmt.Println()
}

func TestStackExampleCode(t *tes.T) {
	fmt.Println("STACK EXAMPLE:")

	// Create a new empty stack using the universal constructor.
	var stack = fra.Stack[string]()
	fmt.Println("The empty stack:", stack)

	// Add some values to it.
	stack.AddValue("foo")
	fmt.Println("The stack with one value on it:", stack)
	stack.AddValue("bar")
	fmt.Println("The stack with two values on it:", stack)
	stack.AddValue("baz")
	fmt.Println("The stack with three values on it:", stack)

	// Remove the top value from the stack.
	var top = stack.RemoveLast()
	fmt.Println("The top value was:", top)
	fmt.Println("The stack with only two values on it:", stack)

	// Remove all values from the stack.
	stack.RemoveAll()
	fmt.Println("The stack with no more values on it:", stack)
	var isEmpty = stack.IsEmpty()
	fmt.Println("The stack is now empty?", isEmpty)
	fmt.Println()
}

func TestQueueExampleCode(t *tes.T) {
	fmt.Println("QUEUE EXAMPLE:")

	// Create a wait group for synchronization.
	var wg = new(syn.WaitGroup)
	defer wg.Wait()

	// Create a new queue with a specific capacity using the universal constructor.
	var queue = fra.QueueWithCapacity[int](12)
	fmt.Println("The empty queue:", queue)

	// Add some values to the queue.
	for i := 1; i < 10; i++ {
		queue.AddValue(i)
		fmt.Println("Added value:", i)
	}
	fmt.Println("The partially filled queue:", queue)

	// Remove values from the queue in the background.
	wg.Add(1)
	go func() {
		defer wg.Done()
		var value int
		var ok = true
		for i := 1; ok; i++ {
			value, ok = queue.RemoveFirst()
			if ok {
				fmt.Println("Removed value:", value)
			}
		}
		fmt.Println("The closed queue:", queue)
		fmt.Println()
	}()

	// Add some more values to the queue.
	for i := 10; i < 31; i++ {
		queue.AddValue(i)
		fmt.Println("Added value:", i)
	}
	queue.CloseChannel()
}

func TestCatalogExampleCode(t *tes.T) {
	fmt.Println("CATALOG EXAMPLE:")

	// Create a new catalog using the universal constructor.
	var catalog = fra.CatalogFromMap[string, int64](
		map[string]int64{
			"foo": 1,
			"bar": 2,
			"baz": 3,
		},
	)
	fmt.Println("The initialized catalog:", catalog)

	// Add a new association to the catalog.
	catalog.SetValue("fuz", 4)
	fmt.Println("The updated catalog:", catalog)

	// Sort the associations in the catalog by key.
	catalog.SortValues()
	fmt.Println("The sorted catalog:", catalog)

	// List the keys for the catalog.
	var keys = catalog.GetKeys()
	fmt.Println("The keys for the catalog:", keys)

	catalog.ReverseValues()
	fmt.Println("The reversed catalog:", catalog)

	// Retrieve a value from the catalog.
	var value = catalog.GetValue("bar")
	fmt.Println("The value for the \"bar\" key is", value)

	// Remove a value from the catalog.
	catalog.RemoveValue("foo")
	fmt.Println("The smaller catalog:", catalog)

	// Change an existing value in the catalog.
	catalog.SetValue("baz", 5)
	fmt.Println("The updated catalog:", catalog)
	fmt.Println()
}

func TestIteratorExampleCode(t *tes.T) {
	fmt.Println("ITERATOR EXAMPLE:")

	// Create a list using the universal constructor and an iterator for it.
	var list = fra.ListFromArray[string](
		[]string{"foo", "bar", "baz"},
	)
	var iterator = list.GetIterator()

	// Iterate over the values in order.
	fmt.Println("The list values in order:")
	for iterator.HasNext() {
		var value = iterator.GetNext()
		fmt.Println("    value:", value)
	}

	// Go to a specific value in the list.
	iterator.SetSlot(2)
	var value = iterator.GetPrevious()
	fmt.Println("The second value in the list is:", value)

	// Iterate over the values in reverse order.
	fmt.Println("The list values in reverse order:")
	iterator.ToEnd()
	for iterator.HasPrevious() {
		var value = iterator.GetPrevious()
		fmt.Println("    value:", value)
	}
	fmt.Println()
}
