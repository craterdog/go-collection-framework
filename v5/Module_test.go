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

package module_test

import (
	fmt "fmt"
	fra "github.com/craterdog/go-collection-framework/v5"
	syn "sync"
	tes "testing"
)

func TestModuleExampleCode(t *tes.T) {
	fmt.Println("MODULE EXAMPLE:")

	// Create a new association.
	var association = fra.Association[string, int]("A", 1)
	fmt.Println(association)

	// Create a new array collection of size 3.
	var array = fra.ArrayWithSize[int64](3)
	fmt.Println(array)

	// Create a new array collection from an intrinsic Go array of values.
	array = fra.Array[int64]([]int64{1, 2, 3})
	fmt.Println(array)

	// Create a new array collection from a sequence of values.
	array = fra.ArrayFromSequence[int64](array)
	fmt.Println(array)

	// Create a new empty map collection.
	var map_ = fra.Map[string, int64](nil)
	fmt.Println(map_)

	// Create a new map collection from an intrinsic Go map of associations.
	map_ = fra.Map[string, int64](map[string]int64{
		"one": 1,
		"two": 2,
	})
	fmt.Println(map_)

	// Create a new map collection from an intrinsic Go array of associations.
	map_ = fra.MapFromArray[string, int64](map_.AsArray())
	fmt.Println(map_)

	// Create a new map collection from a sequence of associations.
	map_ = fra.MapFromSequence[string, int64](map_)
	fmt.Println(map_)

	// Create a new empty list collection.
	var list = fra.List[string]()
	fmt.Println(list)

	// Create a new list collection using an intrinsic Go array of values.
	list = fra.ListFromArray[string]([]string{
		"Hello",
		"World",
	})
	fmt.Println(list)

	// Create a new list collection using a sequence of values.
	list = fra.ListFromSequence[string](map_.GetKeys())
	fmt.Println(list)

	// Create a new empty catalog collection.
	var catalog = fra.Catalog[string, int64]()
	fmt.Println(catalog)

	// Create a new catalog collection from a map collection.
	catalog = fra.CatalogFromSequence[string, int64](map_)
	catalog.SetValue("delta", 4)
	catalog.SortValues()
	fmt.Println(catalog)

	// Iterate through the catalog associations.
	var iterator = catalog.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		fmt.Println(association)
	}
	fmt.Println()

	// Create a new empty set collection.
	var set = fra.Set[string]()
	fmt.Println(set)

	// Create a new set collection from an intrinsic Go array of values.
	set = fra.SetFromArray[string]([]string{"c", "a", "b"})
	fmt.Println(set)

	// Create a new set collection from the keys in a catalog collection.
	set = fra.SetFromSequence[string](catalog.GetKeys())
	fmt.Println(set)

	// Create a new empty stack collection with a capacity of 4.
	var stack = fra.StackWithCapacity[string](4)
	fmt.Println(stack)

	// Create a new stack collection from a list collection.
	stack = fra.StackFromSequence[string](list)
	fmt.Println(stack)

	// Create a new empty queue collection with a capacity of 5.
	var queue = fra.QueueWithCapacity[string](5)
	fmt.Println(queue)

	// Create a new queue collection from a set collection.
	queue = fra.QueueFromSequence[string](set)
	fmt.Println(queue)
	fmt.Println()
}

func TestArrayExampleCode(t *tes.T) {
	fmt.Println("ARRAY EXAMPLE:")

	// Create a new wrapped array using the universal constructor.
	var array = fra.Array[string]([]string{"foo", "bar", "baz"})
	fmt.Println("The array is:", array)

	// Retrieve the first value in the array.
	var first = array.GetValue(1)
	fmt.Println("The first value is:", first)

	// Retrieve the last value in the array using negative indexing.
	var last = array.GetValue(-1)
	fmt.Println("The last value is:", last)

	// Set the second value in the array to a new value.
	array.SetValue(2, "fuz")
	var second = array.GetValue(2)
	fmt.Println("The second value is now:", second)
	fmt.Println("The updated array is:", array)

	// Sort the values in the array.
	array.SortValues()
	fmt.Println("The sorted array is:", array)
	fmt.Println()
}

func TestMapExampleCode(t *tes.T) {
	fmt.Println("MAP EXAMPLE:")

	// Create a new wrapped map using the universal constructor.
	var map_ = fra.Map[string, int](map[string]int{
		"foo": 1,
		"bar": 2,
	})
	fmt.Println("The initial map is:", map_)

	// Add a new association.
	map_.SetValue("baz", 3)
	fmt.Println("The augmented map is now:", map_)

	// Retrieve a couple values.
	var one = map_.GetValue("foo")
	fmt.Println("The \"foo\" value is:", one)
	var three = map_.GetValue("baz")
	fmt.Println("The \"baz\" value is:", three)

	// Change the value of an association.
	map_.SetValue("bar", 42)
	var answer = map_.GetValue("bar")
	fmt.Println("The \"bar\" value is now:", answer)
	fmt.Println("The updated map is:", map_)

	// Remove an association.
	map_.RemoveValue("baz")
	fmt.Println("The smaller map is:", map_)

	// Empty the map.
	map_.RemoveAll()
	fmt.Println("The empty map is:", map_)
	fmt.Println()
}

func TestListExampleCode(t *tes.T) {
	fmt.Println("LIST EXAMPLE:")

	// Create a new list using the universal constructor.
	var list = fra.ListFromArray[string]([]string{"bar", "foo", "bax"})
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
	var set1 = fra.SetFromArray[string]([]string{"alpha", "beta", "gamma"})
	fmt.Println("The first set is:", set1)
	var set2 = fra.SetFromArray[string]([]string{"beta", "gamma", "delta"})
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
	var catalog = fra.CatalogFromMap[string, int64](map[string]int64{
		"foo": 1,
		"bar": 2,
		"baz": 3,
	})
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
	var list = fra.ListFromArray[string]([]string{"foo", "bar", "baz"})
	var iterator = list.GetIterator()

	// Iterate over the values in order.
	fmt.Println("The list values in order:")
	for iterator.HasNext() {
		var value = iterator.GetNext()
		fmt.Println("    value:", value)
	}

	// Go to a specific value in the list.
	iterator.ToSlot(2)
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
