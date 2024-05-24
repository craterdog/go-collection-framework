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
	col "github.com/craterdog/go-collection-framework/v4"
	ass "github.com/stretchr/testify/assert"
	syn "sync"
	tes "testing"
)

func TestCDCNConstructor(t *tes.T) {
	col.CDCN()
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The CDCN constructor does not take any arguments.", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	col.CDCN("dummy")
}

func TestXMLConstructor(t *tes.T) {
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The XML notation is not yet supported.", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	col.XML() // This should panic.
}

func TestJSONConstructor(t *tes.T) {
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The JSON notation is not yet supported.", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	col.JSON() // This should panic.
}

func TestModuleExampleCode(t *tes.T) {
	fmt.Println("MODULE EXAMPLE:")

	// Create a new array collection of size 3 from a primitive array.
	var array = col.Array[int64](3, []int64{1, 2, 3})
	fmt.Println(array)

	// Create a new map collection from a CDCN source string.
	var notation = col.CDCN()
	var map_ = col.Map[string, int64](notation, `["alpha": 1, "beta": 2, "gamma": 3](Map)`)
	fmt.Println(map_)

	// Create a new list collection using an array collection.
	var list = col.List[int64](array)
	fmt.Println(list)

	// Create a new catalog collection from a map collection.
	var catalog = col.Catalog[string, int64](map_)
	catalog.SetValue("delta", 4)
	catalog.SortValues()
	fmt.Println(catalog)

	// Iterate through the catalog associations.
	var iterator = catalog.GetIterator()
	for iterator.HasNext() {
		var association = iterator.GetNext()
		fmt.Println(association.GetKey(), association.GetValue())
	}
	fmt.Println()

	// Create a new set collection from the keys in a catalog collection.
	var set = col.Set[string](catalog.GetKeys())
	fmt.Println(set)

	// Create a new stack collection with a capacity of 4 from a list collection.
	var stack = col.Stack[int64](4, list)
	fmt.Println(stack)

	// Create a new queue collection with a capacity of 5 from a set collection.
	var queue = col.Queue[string](notation, 5, set)
	fmt.Println(queue)
}

func TestArrayExampleCode(t *tes.T) {
	fmt.Println("ARRAY EXAMPLE:")

	// Create a new wrapped array using the universal constructor.
	var array = col.Array[string]([]string{"foo", "bar", "baz"})
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
}

func TestMapExampleCode(t *tes.T) {
	fmt.Println("MAP EXAMPLE:")

	// Create a new wrapped map using the universal constructor.
	var map_ = col.Map[string, int](map[string]int{
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
}

func TestListExampleCode(t *tes.T) {
	fmt.Println("LIST EXAMPLE:")

	// Create a new list using the universal constructor.
	var list = col.List[string]([]string{"bar", "foo", "bax"})
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
}

func TestSetExampleCode(t *tes.T) {
	fmt.Println("SET EXAMPLE:")

	// Create two sets with overlapping values using the universal constructor.
	var set1 = col.Set[string]([]string{"alpha", "beta", "gamma"})
	fmt.Println("The first set is:", set1)
	var set2 = col.Set[string](`[
    "beta"
    "gamma"
    "delta"
](Set)`)
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
}

func TestStackExampleCode(t *tes.T) {
	fmt.Println("STACK EXAMPLE:")

	// Create a new empty stack using the universal constructor.
	var stack = col.Stack[string]()
	fmt.Println("The empty stack:", stack)

	// Add some values to it.
	stack.AddValue("foo")
	fmt.Println("The stack with one value on it:", stack)
	stack.AddValue("bar")
	fmt.Println("The stack with two values on it:", stack)
	stack.AddValue("baz")
	fmt.Println("The stack with three values on it:", stack)

	// Remove the top value from the stack.
	var top = stack.RemoveTop()
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
	var queue = col.Queue[int](12)
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
			value, ok = queue.RemoveHead()
			if ok {
				fmt.Println("Removed value:", value)
			}
		}
		fmt.Println("The closed queue:", queue)
	}()

	// Add some more values to the queue.
	for i := 10; i < 31; i++ {
		queue.AddValue(i)
		fmt.Println("Added value:", i)
	}
	queue.CloseQueue()
	fmt.Println()
}

func TestCatalogExampleCode(t *tes.T) {
	fmt.Println("CATALOG EXAMPLE:")

	// Create a new catalog using the universal constructor.
	var catalog = col.Catalog[string, int64](`[
    "foo": 1
    "bar": 2
    "baz": 3
](Catalog)`)
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
}

func TestIteratorExampleCode(t *tes.T) {
	fmt.Println("ITERATOR EXAMPLE:")

	// Create a list using the universal constructor and an iterator for it.
	var list = col.List[string]([]string{"foo", "bar", "baz"})
	var iterator = list.GetIterator()

	// Iterate over the values in order.
	fmt.Println("The list values in order:")
	for iterator.HasNext() {
		var value = iterator.GetNext()
		fmt.Println("\tvalue:", value)
	}
	fmt.Println()

	// Go to a specific value in the list.
	iterator.ToSlot(2)
	var value = iterator.GetPrevious()
	fmt.Println("The second value in the list is:", value)

	// Iterate over the values in reverse order.
	fmt.Println("The list values in reverse order:")
	iterator.ToEnd()
	for iterator.HasPrevious() {
		var value = iterator.GetPrevious()
		fmt.Println("\tvalue:", value)
	}
	fmt.Println()
}

func TestCDCNImportExampleCode(t *tes.T) {
	fmt.Println("CDCN IMPORT EXAMPLE:")

	// Define the source string for a catalog collection containing all the Go primitive types.
	var source = `[
    "boolean": true
    "unsigned": 0xa
    "integer": 42
    "float": 0.125
    "complex": (3.0-4.0i)
    "rune": '☺'
    "string": "Hello World!"
    "array": [
        1
        2
        3
    ](array)
    "map": [
        "alpha": 1
        "beta": 2
        "gamma": 3
    ](map)
](Catalog)`

	// Parse the source string containing the CDCN for the catalog.
	var cdcn = col.CDCN()
	var collection = cdcn.ParseSource(source)

	// Print back out the formatted CDCN string for the catalog and see if they match.
	fmt.Println(collection)
}

func TestCDCNExportExampleCode(t *tes.T) {
	fmt.Println("CDCN EXPORT EXAMPLE:")

	// Define "zero" values for each primitive Go type.
	var v1 bool
	var v2 uint
	var v3 int
	var v4 float64
	var v5 complex128
	var v6 rune
	var v7 string
	var v8 []string
	var v9 map[string]int

	// Create a list of the "zero" values using the universal constructor.
	var collection = col.List[any]([]any{v1, v2, v3, v4, v5, v6, v7, v8, v9})

	// Print the list using the canonical CDCN format.
	var cdcn = col.CDCN()
	var source = cdcn.FormatCollection(collection)
	fmt.Println(source)

	// Attempt to print the list using the canonical JSON format.
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("The JSON notation is not yet supported.")
			fmt.Println()
		}
	}()
	var json = col.JSON()
	source = json.FormatCollection(collection)
	fmt.Println(source)
}
