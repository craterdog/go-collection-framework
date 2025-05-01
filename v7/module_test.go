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
	fra "github.com/craterdog/go-collection-framework/v7"
	col "github.com/craterdog/go-collection-framework/v7/collection"
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

	// Create a new list from an array.
	var list = fra.ListFromArray[string](
		[]string{"bar", "foo", "bax"},
	)
	fmt.Println("The initialized list:", list)
	fmt.Println()

	// Add some more values to the list.
	list.AppendValue("fuz")
	list.AppendValue("box")
	fmt.Println("The augmented list:", list)
	fmt.Println()

	// Change a value in the list.
	list.SetValue(2, "foz")
	fmt.Println("The updated list:", list)
	fmt.Println()

	// Insert a new value at the beginning of the list (slot 0).
	list.InsertValue(0, "bax")
	fmt.Println("The updated list:", list)
	fmt.Println()

	// Insert a new value at the end of the list (slot N).
	list.InsertValue(6, "bux")
	fmt.Println("The updated list:", list)
	fmt.Println()

	// Sort the values in the list.
	list.SortValues()
	fmt.Println("The sorted list:", list)
	fmt.Println()

	// Randomly shuffle the values in the list.
	list.ShuffleValues()
	fmt.Println("The shuffled list:", list)
	fmt.Println()

	// Remove a value from the list.
	list.RemoveValue(4)
	fmt.Println("The shortened list:", list)
	fmt.Println()

	// Remove all values from the list.
	list.RemoveAll()
	fmt.Println("The empty list:", list)
	fmt.Println()
}

func TestSetExampleCode(t *tes.T) {
	fmt.Println("SET EXAMPLE:")

	// Create two sets with overlapping values.
	var set1 = fra.SetFromArray[string](
		[]string{"alpha", "beta", "gamma"},
	)
	fmt.Println("The first set is:", set1)
	fmt.Println()
	var set2 = fra.SetFromArray[string](
		[]string{"beta", "gamma", "delta"},
	)
	fmt.Println("The second set is:", set2)
	fmt.Println()

	// Find the logical union of the two sets.
	var Set = set1.GetClass()
	fmt.Println("The union of the two sets is:", Set.Or(set1, set2))
	fmt.Println()

	// Find the logical difference between the two sets.
	fmt.Println("The first set minus the second set is:", Set.Sans(set1, set2))
	fmt.Println()

	// Find the logical intersection of the two sets.
	fmt.Println("The intersection of the two sets is:", Set.And(set1, set2))
	fmt.Println()

	// Find the logical exclusion of the two sets.
	fmt.Println("The exclusion of the two sets is:", Set.Xor(set1, set2))
	fmt.Println()

	// Add an existing value to the first set.
	set1.AddValue("beta")
	fmt.Println("Adding an existing value to a set does not change it:", set1)
	fmt.Println()
}

func TestStackExampleCode(t *tes.T) {
	fmt.Println("STACK EXAMPLE:")

	// Create a new empty stack.
	var stack = fra.Stack[string]()
	fmt.Println("The empty stack:", stack)
	fmt.Println()

	// Add some values to it.
	stack.AddValue("foo")
	fmt.Println("The stack with one value on it:", stack)
	fmt.Println()
	stack.AddValue("bar")
	fmt.Println("The stack with two values on it:", stack)
	fmt.Println()
	stack.AddValue("baz")
	fmt.Println("The stack with three values on it:", stack)
	fmt.Println()

	// Remove the top value from the stack.
	var top = stack.RemoveLast()
	fmt.Println("The top value was:", top)
	fmt.Println()
	fmt.Println("The stack with only two values on it:", stack)
	fmt.Println()

	// Remove all values from the stack.
	stack.RemoveAll()
	fmt.Println("The stack with no more values on it:", stack)
	fmt.Println()
	var isEmpty = stack.IsEmpty()
	fmt.Println("The stack is now empty?", isEmpty)
	fmt.Println()
}

func TestQueueExampleCode(t *tes.T) {
	fmt.Println("QUEUE EXAMPLE:")

	// Create a wait group for synchronization.
	var wg = new(syn.WaitGroup)
	defer wg.Wait()

	// Create a new queue with a specific capacity.
	var queue = fra.QueueWithCapacity[int](12)
	fmt.Println("The empty queue:", queue)
	fmt.Println()

	// Add some values to the queue.
	for i := 1; i < 10; i++ {
		queue.AddValue(i)
		fmt.Println("Added value:", i)
	}
	fmt.Println()
	fmt.Println("The partially filled queue:", queue)
	fmt.Println()

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

	// Create a new catalog from a map.
	var catalog = fra.CatalogFromMap[string, int64](
		map[string]int64{
			"foo": 1,
			"bar": 2,
			"baz": 3,
		},
	)
	fmt.Println("The initialized catalog:", catalog)
	fmt.Println()

	// Add a new association to the catalog.
	catalog.SetValue("fuz", 4)
	fmt.Println("The updated catalog:", catalog)
	fmt.Println()

	// Sort the associations in the catalog by key.
	catalog.SortValues()
	fmt.Println("The sorted catalog:", catalog)
	fmt.Println()

	// List the keys for the catalog.
	var keys = catalog.GetKeys()
	fmt.Println("The keys for the catalog:", keys)
	fmt.Println()

	catalog.ReverseValues()
	fmt.Println("The reversed catalog:", catalog)
	fmt.Println()

	// Retrieve a value from the catalog.
	var value = catalog.GetValue("bar")
	fmt.Println("The value for the \"bar\" key is", value)
	fmt.Println()

	// Remove a value from the catalog.
	catalog.RemoveValue("foo")
	fmt.Println("The smaller catalog:", catalog)
	fmt.Println()

	// Change an existing value in the catalog.
	catalog.SetValue("baz", 5)
	fmt.Println("The updated catalog:", catalog)
	fmt.Println()
}

func TestCollatorExampleCode(t *tes.T) {
	fmt.Println("COLLATOR EXAMPLE:")

	// Create a collator with the default maximum depth.
	var collator = fra.Collator[any]()

	// Collate two strings.
	var s1 = "first"
	var s2 = "second"
	fmt.Println(
		"The first and second strings are equal:",
		collator.CompareValues(s1, s2),
	)
	fmt.Println(
		"The first string is ranked before the second strings:",
		collator.RankValues(s1, s2) == fra.LesserRank,
	)
	fmt.Println()

	// Collate two arrays.
	var a1 = []int{1, 2, 3}
	var a2 = []int{1, 3, 2}
	fmt.Println(
		"The first and second arrays are equal:",
		collator.CompareValues(a1, a2),
	)
	fmt.Println(
		"The first array is ranked before the second array:",
		collator.RankValues(a1, a2) == fra.LesserRank,
	)
	fmt.Println()
}

func TestSorterExampleCode(t *tes.T) {
	fmt.Println("SORTER EXAMPLE:")

	// Create a sorter with the default (natural) ranker.
	var sorter = fra.Sorter[string]()

	// Create an array.
	var array = []string{
		"alpha",
		"beta",
		"gamma",
		"delta",
	}
	fmt.Println("The initial ordering of the values:", array)
	fmt.Println()

	// Sort the values in alphabetical order.
	sorter.SortValues(array)
	fmt.Println("The values in alphabetical order:", array)
	fmt.Println()

	// Sort the values in reverse order.
	sorter.ReverseValues(array)
	fmt.Println("The values in reverse order:", array)
	fmt.Println()

	// Shuffle the order of the values.
	sorter.ShuffleValues(array)
	fmt.Println("The values in random order:", array)
	fmt.Println()

	// Sort the values with a custom ranking function.
	sorter = fra.SorterWithRanker[string](
		func(first, second string) fra.Rank {
			switch {
			case first < second:
				return fra.GreaterRank
			case first > second:
				return fra.LesserRank
			default:
				return fra.EqualRank
			}
		},
	)
	sorter.SortValues(array)
	fmt.Println("The values in custom order:", array)
	fmt.Println()
}

func TestIteratorExampleCode(t *tes.T) {
	fmt.Println("ITERATOR EXAMPLE:")

	// Create a list from an array.
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
	fmt.Println()

	// Go to a specific value in the list.
	iterator.SetSlot(2)
	var value = iterator.GetPrevious()
	fmt.Println("The second value in the list is:", value)
	fmt.Println()

	// Iterate over the values in reverse order.
	fmt.Println("The list values in reverse order:")
	iterator.ToEnd()
	for iterator.HasPrevious() {
		var value = iterator.GetPrevious()
		fmt.Println("    value:", value)
	}
	fmt.Println()
}
