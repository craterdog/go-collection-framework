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

package collection_test

import (
	age "github.com/craterdog/go-collection-framework/v7/agent"
	col "github.com/craterdog/go-collection-framework/v7/collection"
	ass "github.com/stretchr/testify/assert"
	syn "sync"
	tes "testing"
)

type Integer int

func TestCatalogConstructors(t *tes.T) {
	var class = col.CatalogClass[rune, int64]()
	class.Catalog()
	class.CatalogFromArray([]col.AssociationLike[rune, int64]{})
	class.CatalogFromMap(map[rune]int64{})
	var sequence = class.CatalogFromMap(map[rune]int64{
		'a': 1,
		'b': 2,
		'c': 3,
	})
	var catalog = class.CatalogFromSequence(sequence)
	ass.Equal(t, sequence.AsArray(), catalog.AsArray())
}

func TestCatalogsWithStringsAndIntegers(t *tes.T) {
	var catalogCollatorClass = age.CollatorClass[col.CatalogLike[string, int]]()
	var catalogCollator = catalogCollatorClass.Collator()
	var keys = col.ListClass[string]().ListFromArray([]string{"foo", "bar"})
	var associationClass = col.AssociationClass[string, int]()
	var association1 = associationClass.Association("foo", 1)
	var association2 = associationClass.Association("bar", 2)
	var association3 = associationClass.Association("baz", 3)
	var catalogClass = col.CatalogClass[string, int]()
	var catalog = catalogClass.Catalog()
	ass.True(t, catalog.IsEmpty())
	ass.True(t, catalog.GetSize() == 0)
	ass.Equal(t, []string{}, catalog.GetKeys().AsArray())
	ass.Equal(t, []col.AssociationLike[string, int]{}, catalog.AsArray())
	var iterator = catalog.GetIterator()
	ass.False(t, iterator.HasNext())
	ass.False(t, iterator.HasPrevious())
	iterator.ToStart()
	iterator.ToEnd()
	catalog.SortValues()
	catalog.ShuffleValues()
	catalog.RemoveAll()
	catalog.SetValue(association1.GetKey(), association1.GetValue())
	ass.False(t, catalog.IsEmpty())
	ass.True(t, catalog.GetSize() == 1)
	catalog.SetValue(association2.GetKey(), association2.GetValue())
	catalog.SetValue(association3.GetKey(), association3.GetValue())
	ass.True(t, catalog.GetSize() == 3)
	var catalog2 = catalogClass.CatalogFromSequence(catalog)
	ass.True(t, catalogCollator.CompareValues(catalog, catalog2))
	var m = col.CatalogClass[string, int]().CatalogFromMap(map[string]int{
		"foo": 1,
		"bar": 2,
		"baz": 3,
	})
	var associationCollatorClass = age.CollatorClass[col.AssociationLike[string, int]]()
	var associationCollator = associationCollatorClass.Collator()
	var catalog3 = catalogClass.CatalogFromSequence(m)
	catalog2.SortValues()
	catalog3.SortValuesWithRanker(associationCollator.RankValues)
	ass.True(t, catalogCollator.CompareValues(catalog2, catalog3))
	iterator = catalog.GetIterator()
	ass.True(t, iterator.HasNext())
	ass.False(t, iterator.HasPrevious())
	ass.Equal(t, association1, iterator.GetNext())
	ass.True(t, iterator.HasPrevious())
	iterator.ToEnd()
	ass.False(t, iterator.HasNext())
	ass.True(t, iterator.HasPrevious())
	ass.Equal(t, association3, iterator.GetPrevious())
	ass.True(t, iterator.HasNext())
	ass.Equal(t, []string{"foo", "bar", "baz"}, catalog.GetKeys().AsArray())
	ass.Equal(t, 3, int(catalog.GetValue("baz")))
	catalog.SetValue("bar", 5)
	ass.Equal(t, []int{1, 5}, catalog.GetValues(keys).AsArray())
	catalog.SortValues()
	ass.Equal(t, []string{"bar", "baz", "foo"}, catalog.GetKeys().AsArray())
	catalog.ReverseValues()
	ass.Equal(t, []string{"foo", "baz", "bar"}, catalog.GetKeys().AsArray())
	catalog.ReverseValues()
	ass.Equal(t, []int{1, 5}, catalog.RemoveValues(keys).AsArray())
	ass.True(t, catalog.GetSize() == 1)
	ass.Equal(t, 3, int(catalog.RemoveValue("baz")))
	ass.True(t, catalog.IsEmpty())
	ass.True(t, catalog.GetSize() == 0)
	catalog.RemoveAll()
	ass.True(t, catalog.IsEmpty())
	ass.True(t, catalog.GetSize() == 0)
}

func TestCatalogsWithMerge(t *tes.T) {
	var collatorClass = age.CollatorClass[col.CatalogLike[string, int]]()
	var collator = collatorClass.Collator()
	var associationClass = col.AssociationClass[string, int]()
	var association1 = associationClass.Association("foo", 1)
	var association2 = associationClass.Association("bar", 2)
	var association3 = associationClass.Association("baz", 3)
	var catalogClass = col.CatalogClass[string, int]()
	var catalog1 = catalogClass.Catalog()
	catalog1.SetValue(association1.GetKey(), association1.GetValue())
	catalog1.SetValue(association2.GetKey(), association2.GetValue())
	var catalog2 = catalogClass.Catalog()
	catalog2.SetValue(association2.GetKey(), association2.GetValue())
	catalog2.SetValue(association3.GetKey(), association3.GetValue())
	var catalog3 = catalogClass.Merge(catalog1, catalog2)
	var catalog4 = catalogClass.Catalog()
	catalog4.SetValue(association1.GetKey(), association1.GetValue())
	catalog4.SetValue(association2.GetKey(), association2.GetValue())
	catalog4.SetValue(association3.GetKey(), association3.GetValue())
	ass.True(t, collator.CompareValues(catalog3, catalog4))
}

func TestCatalogsWithExtract(t *tes.T) {
	var keys = col.ListClass[string]().ListFromArray([]string{"foo", "baz"})
	var associationClass = col.AssociationClass[string, int]()
	var association1 = associationClass.Association("foo", 1)
	var association2 = associationClass.Association("bar", 2)
	var association3 = associationClass.Association("baz", 3)
	var catalogClass = col.CatalogClass[string, int]()
	var catalog1 = catalogClass.Catalog()
	catalog1.SetValue(association1.GetKey(), association1.GetValue())
	catalog1.SetValue(association2.GetKey(), association2.GetValue())
	catalog1.SetValue(association3.GetKey(), association3.GetValue())
	var catalog2 = catalogClass.Extract(catalog1, keys)
	var catalog3 = catalogClass.Catalog()
	catalog3.SetValue(association1.GetKey(), association1.GetValue())
	catalog3.SetValue(association3.GetKey(), association3.GetValue())
	var collatorClass = age.CollatorClass[col.CatalogLike[string, int]]()
	var collator = collatorClass.Collator()
	ass.True(t, collator.CompareValues(catalog2, catalog3))
	var catalog4 = catalogClass.CatalogFromArray([]col.AssociationLike[string, int]{
		association1,
		association2,
		association3,
	})
	ass.True(t, collator.CompareValues(catalog1, catalog4))
}

func TestCatalogsWithEmptyCatalogs(t *tes.T) {
	var keys = col.ListClass[int]().List()
	var catalogClass = col.CatalogClass[int, string]()
	var catalog1 = catalogClass.Catalog()
	var catalog2 = catalogClass.Catalog()
	var catalog3 = catalogClass.Merge(catalog1, catalog2)
	var catalog4 = catalogClass.Extract(catalog1, keys)
	var collator = age.CollatorClass[col.CatalogLike[int, string]]().Collator()
	ass.True(t, collator.CompareValues(catalog1, catalog2))
	ass.True(t, collator.CompareValues(catalog2, catalog3))
	ass.True(t, collator.CompareValues(catalog3, catalog4))
	ass.True(t, collator.CompareValues(catalog4, catalog1))
}

func TestListConstructors(t *tes.T) {
	var listClass = col.ListClass[int64]()
	listClass.List()
	var sequence = listClass.ListFromArray([]int64{1, 2, 3})
	var list = listClass.ListFromSequence(sequence)
	ass.Equal(t, sequence.AsArray(), list.AsArray())
}

func TestListsWithStrings(t *tes.T) {
	var listClass = col.ListClass[string]()
	var collator = age.CollatorClass[col.ListLike[string]]().Collator()
	var foo = listClass.ListFromArray([]string{"foo"})
	var bar = listClass.ListFromArray([]string{"bar"})
	var baz = listClass.ListFromArray([]string{"baz"})
	var foz = listClass.ListFromArray([]string{"foz"})
	var barbaz = listClass.ListFromArray([]string{"bar", "baz"})
	var bazbaz = listClass.ListFromArray([]string{"baz", "baz"})
	var foobar = listClass.ListFromArray([]string{"foo", "bar"})
	var baxbaz = listClass.ListFromArray([]string{"bax", "baz"})
	var baxbez = listClass.ListFromArray([]string{"bax", "bez"})
	var barfoobax = listClass.ListFromArray([]string{"bar", "foo", "bax"})
	var foobazbar = listClass.ListFromArray([]string{"foo", "baz", "bar"})
	var foobarbaz = listClass.ListFromArray([]string{"foo", "bar", "baz"})
	var barbazfoo = listClass.ListFromArray([]string{"bar", "baz", "foo"})
	var list = listClass.List()
	ass.True(t, list.IsEmpty())
	ass.True(t, list.GetSize() == 0)
	ass.False(t, list.ContainsValue("bax"))
	ass.Equal(t, []string{}, list.AsArray())
	var iterator = list.GetIterator()
	ass.False(t, iterator.HasNext())
	ass.False(t, iterator.HasPrevious())
	iterator.ToStart()
	iterator.ToEnd()
	list.ShuffleValues()
	list.SortValues()
	list.RemoveAll()                      //       [ ]
	list.AppendValue("foo")               //       ["foo"]
	ass.False(t, list.IsEmpty())          //       ["foo"]
	ass.True(t, list.GetSize() == 1)      //       ["foo"]
	ass.Equal(t, "foo", list.GetValue(1)) //       ["foo"]
	list.AppendValues(barbaz)             //       ["foo", "bar", "baz"]
	ass.True(t, list.GetSize() == 3)      //       ["foo", "bar", "baz"]
	ass.Equal(t, "foo", list.GetValue(1)) //       ["foo", "bar", "baz"]
	ass.True(t, collator.CompareValues(listClass.ListFromArray(list.AsArray()), list))
	ass.Equal(t, barbaz.AsArray(), list.GetValues(2, 3).AsArray())
	ass.Equal(t, foo.AsArray(), list.GetValues(1, 1).AsArray())
	var list2 = listClass.ListFromSequence(list)
	ass.True(t, collator.CompareValues(list, list2))
	var array = listClass.ListFromArray([]string{"foo", "bar", "baz"})
	var list3 = listClass.ListFromSequence(array)
	list2.SortValues()
	list3.SortValues()
	ass.True(t, collator.CompareValues(list2, list3))
	iterator = list.GetIterator()               // ["foo", "bar", "baz"]
	ass.True(t, iterator.HasNext())             // ["foo", "bar", "baz"]
	ass.False(t, iterator.HasPrevious())        // ["foo", "bar", "baz"]
	ass.Equal(t, "foo", iterator.GetNext())     // ["foo", "bar", "baz"]
	ass.True(t, iterator.HasPrevious())         // ["foo", "bar", "baz"]
	iterator.ToEnd()                            // ["foo", "bar", "baz"]
	ass.False(t, iterator.HasNext())            // ["foo", "bar", "baz"]
	ass.True(t, iterator.HasPrevious())         // ["foo", "bar", "baz"]
	ass.Equal(t, "baz", iterator.GetPrevious()) // ["foo", "bar", "baz"]
	ass.True(t, iterator.HasNext())             // ["foo", "bar", "baz"]
	list.ShuffleValues()                        // [ ?, ?, ? ]
	list.RemoveAll()                            // [ ]
	ass.True(t, list.IsEmpty())                 // [ ]
	ass.True(t, list.GetSize() == 0)            // [ ]
	list.InsertValue(0, "baz")                  // ["baz"]
	ass.True(t, list.GetSize() == 1)            // ["baz"]
	ass.Equal(t, "baz", list.GetValue(-1))      // ["baz"]
	list.InsertValues(0, foobar)                // ["foo", "bar", "baz"]
	ass.True(t, list.GetSize() == 3)            // ["foo", "bar", "baz"]
	ass.Equal(t, "foo", list.GetValue(-3))      // ["foo", "bar", "baz"]
	ass.Equal(t, "bar", list.GetValue(-2))      // ["foo", "bar", "baz"]
	ass.Equal(t, "baz", list.GetValue(-1))      // ["foo", "bar", "baz"]
	list.ReverseValues()                        // ["baz", "bar", "foo"]
	ass.Equal(t, "foo", list.GetValue(-1))      // ["baz", "bar", "foo"]
	ass.Equal(t, "bar", list.GetValue(-2))      // ["baz", "bar", "foo"]
	ass.Equal(t, "baz", list.GetValue(-3))      // ["baz", "bar", "foo"]
	list.ReverseValues()                        // ["foo", "bar", "baz"]
	ass.True(t, list.GetIndex("foz") == 0)      // ["foo", "bar", "baz"]
	ass.True(t, list.GetIndex("baz") == 3)      // ["foo", "bar", "baz"]
	ass.True(t, list.ContainsValue("baz"))      // ["foo", "bar", "baz"]
	ass.False(t, list.ContainsValue("bax"))     // ["foo", "bar", "baz"]
	ass.True(t, list.ContainsAny(baxbaz))       // ["foo", "bar", "baz"]
	ass.False(t, list.ContainsAny(baxbez))      // ["foo", "bar", "baz"]
	ass.True(t, list.ContainsAll(barbaz))       // ["foo", "bar", "baz"]
	ass.False(t, list.ContainsAll(baxbaz))      // ["foo", "bar", "baz"]
	list.SetValue(3, "bax")                     // ["foo", "bar", "bax"]
	list.InsertValues(3, baz)                   // ["foo", "bar", "bax", "baz"]
	ass.True(t, list.GetSize() == 4)            // ["foo", "bar", "bax", "baz"]
	ass.Equal(t, "baz", list.GetValue(4))       // ["foo", "bar", "bax", "baz"]
	list.InsertValue(4, "bar")                  // ["foo", "bar", "bax", "baz", "bar"]
	ass.True(t, list.GetSize() == 5)            // ["foo", "bar", "bax", "baz", "bar"]
	ass.Equal(t, "bar", list.GetValue(5))       // ["foo", "bar", "bax", "baz", "bar"]
	list.InsertValue(2, "foo")                  // ["foo", "bar", "foo", "bax", "baz", "bar"]
	ass.True(t, list.GetSize() == 6)            // ["foo", "bar", "foo", "bax", "baz", "bar"]
	ass.Equal(t, "bar", list.GetValue(2))       // ["foo", "bar", "foo", "bax", "baz", "bar"]
	ass.Equal(t, "foo", list.GetValue(3))       // ["foo", "bar", "foo", "bax", "baz", "bar"]
	ass.Equal(t, "bax", list.GetValue(4))       // ["foo", "bar", "foo", "bax", "baz", "bar"]
	ass.Equal(t, bar.AsArray(), list.GetValues(6, 6).AsArray())
	list.InsertValues(5, baz)             //       ["foo", "bar", "foo", "bax", "baz", "baz", "bar"]
	ass.True(t, list.GetSize() == 7)      //       ["foo", "bar", "foo", "bax", "baz", "baz", "bar"]
	ass.Equal(t, "bax", list.GetValue(4)) //       ["foo", "bar", "foo", "bax", "baz", "baz", "bar"]
	ass.Equal(t, "baz", list.GetValue(5)) //       ["foo", "bar", "foo", "bax", "baz", "baz", "bar"]
	ass.Equal(t, "baz", list.GetValue(6)) //       ["foo", "bar", "foo", "bax", "baz", "baz", "bar"]
	ass.Equal(t, barfoobax.AsArray(), list.GetValues(2, -4).AsArray())
	list.SetValues(2, foobazbar) //                        ["foo", "foo", "baz", "bar", "baz", "baz", "bar"]
	ass.Equal(t, foobazbar.AsArray(), list.GetValues(2, -4).AsArray())
	list.SetValues(-1, foz)
	ass.Equal(t, "foz", list.GetValue(-1)) //      ["foo", "foo", "baz", "bar", "baz", "baz", "foz"]
	list.SortValues()                      //      ["bar", "baz", "baz", "baz", "foo", "foo", "foz"]

	ass.Equal(t, bazbaz.AsArray(), list.RemoveValues(2, -5).AsArray()) // ["bar", "baz", "foo", "foo", "foz"]
	ass.Equal(t, barbaz.AsArray(), list.RemoveValues(1, 2).AsArray())  // ["foo", "foo", "foz"]
	ass.Equal(t, "foz", list.RemoveValue(-1))                          // ["foo", "foo"]
	ass.True(t, list.GetSize() == 2)                                   // ["foo", "foo"]
	list.RemoveAll()                                                   // [ ]
	ass.True(t, list.GetSize() == 0)                                   // [ ]
	list.SortValues()                                                  // [ ]
	list.AppendValues(foobarbaz)                                       // ["foo", "bar", "baz"]
	list.SortValues()                                                  // ["bar", "baz", "foo"]
	ass.Equal(t, barbazfoo.AsArray(), list.AsArray())                  // ["bar", "baz", "foo"]
	list.RemoveAll()                                                   // [ ]
	list.AppendValue("foo")                                            // ["foo"]
	list.SortValues()                                                  // ["foo"]
	ass.True(t, list.GetSize() == 1)                                   // ["foo"]
	ass.Equal(t, "foo", list.GetValue(1))                              // ["foo"]
	list.AppendValue("bar")                                            // ["foo", "bar"]
	list.SortValues()                                                  // ["bar", "foo"]
	ass.True(t, list.GetSize() == 2)                                   // ["bar", "foo"]
	ass.Equal(t, "bar", list.GetValue(1))                              // ["bar", "foo"]
}

func TestListsWithTildes(t *tes.T) {
	var array = col.ListClass[Integer]().ListFromArray([]Integer{3, 1, 4, 5, 9, 2})
	var list = col.ListClass[Integer]().ListFromSequence(array)
	ass.False(t, list.IsEmpty())        // [3,1,4,5,9,2]
	ass.True(t, list.GetSize() == 6)    // [3,1,4,5,9,2]
	ass.True(t, list.GetValue(1) == 3)  // [3,1,4,5,9,2]
	ass.True(t, list.GetValue(-1) == 2) // [3,1,4,5,9,2]
	list.SortValues()                   // [1,2,3,4,5,9]
	ass.True(t, list.GetSize() == 6)    // [1,2,3,4,5,9]
	ass.True(t, list.GetValue(3) == 3)  // [1,2,3,4,5,9]
}

func TestListsWithConcatenate(t *tes.T) {
	var listClass = col.ListClass[int]()
	var collator = age.CollatorClass[col.ListLike[int]]().Collator()
	var onetwothree = listClass.ListFromArray([]int{1, 2, 3})
	var fourfivesix = listClass.ListFromArray([]int{4, 5, 6})
	var onethrusix = listClass.ListFromArray([]int{1, 2, 3, 4, 5, 6})
	var list1 = listClass.List()
	list1.AppendValues(onetwothree)
	var list2 = listClass.List()
	list2.AppendValues(fourfivesix)
	var list3 = listClass.Concatenate(list1, list2)
	var list4 = listClass.List()
	list4.AppendValues(onethrusix)
	ass.True(t, collator.CompareValues(list3, list4))
}

func TestListsWithEmptyLists(t *tes.T) {
	var collator = age.CollatorClass[col.ListLike[int]]().Collator()
	var listClass = col.ListClass[int]()
	var empty = listClass.List()
	var list = listClass.Concatenate(empty, empty)
	ass.True(t, collator.CompareValues(empty, empty))
	ass.True(t, collator.CompareValues(list, empty))
	ass.True(t, collator.CompareValues(empty, list))
	ass.True(t, collator.CompareValues(list, list))
}

func TestQueueConstructors(t *tes.T) {
	var queueClass = col.QueueClass[int64]()
	queueClass.Queue()
	queueClass.QueueWithCapacity(5)
	var sequence = queueClass.QueueFromArray([]int64{1, 2, 3})
	var queue = queueClass.QueueFromSequence(sequence)
	ass.Equal(t, sequence.AsArray(), queue.AsArray())
}

func TestQueueWithConcurrency(t *tes.T) {
	// Create a wait group for synchronization.
	var group = new(syn.WaitGroup)
	defer group.Wait()

	// Create a new queue with a specific capacity.
	var queue = col.QueueClass[int]().QueueWithCapacity(12)
	ass.True(t, queue.GetCapacity() == 12)
	ass.True(t, queue.IsEmpty())
	ass.True(t, queue.GetSize() == 0)

	// Add some values to the queue.
	for i := 1; i < 10; i++ {
		queue.AddValue(i)
	}
	ass.True(t, queue.GetSize() == 9)

	// Remove values from the queue in the background.
	group.Add(1)
	go func() {
		defer group.Done()
		var value int
		var ok = true
		for i := 1; ok; i++ {
			value, ok = queue.RemoveFirst()
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
	queue.CloseChannel()
}

func TestQueueWithFork(t *tes.T) {
	// Create a wait group for synchronization.
	var group = new(syn.WaitGroup)
	defer group.Wait()

	// Create a new queue with a fan out of two.
	var queueClass = col.QueueClass[int]()
	var input = queueClass.QueueWithCapacity(3)
	var outputs = queueClass.Fork(group, input, 2)

	// Remove values from the output queues in the background.
	var readOutput = func(output col.QueueLike[int], name string) {
		defer group.Done()
		var value int
		var ok bool = true
		for i := 1; ok; i++ {
			value, ok = output.RemoveFirst()
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
	input.CloseChannel()
}

func TestQueueWithInvalidFanOut(t *tes.T) {
	// Create a wait group for synchronization.
	var group = new(syn.WaitGroup)
	defer group.Wait()

	// Create a new queue with an invalid fan out.
	var queueClass = col.QueueClass[int]()
	var input = queueClass.QueueWithCapacity(3)
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The fan out size for a queue must be greater than one.", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	queueClass.Fork(group, input, 1) // Should panic here.
}

func TestQueueWithSplitAndJoin(t *tes.T) {
	// Create a wait group for synchronization.
	var group = new(syn.WaitGroup)
	defer group.Wait()

	// Create a new queue with a split of five outputs and a join back to one.
	var queueClass = col.QueueClass[int]()
	var input = queueClass.QueueWithCapacity(3)
	var split = queueClass.Split(group, input, 5)
	var output = queueClass.Join(group, split)

	// Remove values from the output queue in the background.
	group.Add(1)
	go func() {
		defer group.Done()
		var value int
		var ok bool = true
		for i := 1; ok; i++ {
			value, ok = output.RemoveFirst()
			if ok {
				ass.Equal(t, i, value)
			}
		}
	}()

	// Add values to the input queue.
	for i := 1; i < 21; i++ {
		input.AddValue(i)
	}
	input.CloseChannel()
}

func TestQueueWithInvalidSplit(t *tes.T) {
	// Create a wait group for synchronization.
	var group = new(syn.WaitGroup)
	defer group.Wait()

	// Create a new queue with an invalid fan out.
	var queueClass = col.QueueClass[int]()
	var input = queueClass.QueueWithCapacity(3)
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The size of the split must be greater than one.", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	queueClass.Split(group, input, 1) // Should panic here.
}

func TestQueueWithInvalidJoin(t *tes.T) {
	// Create a wait group for synchronization.
	var group = new(syn.WaitGroup)
	defer group.Wait()

	// Create a new queue with an invalid fan out.
	var inputs = col.ListClass[col.QueueLike[int]]().List()
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The number of input queues for a join must be at least one.", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	var queueClass = col.QueueClass[int]()
	queueClass.Join(group, inputs) // Should panic here.
	defer group.Done()
}

func TestSetConstructors(t *tes.T) {
	var collatorClass = age.CollatorClass[int64]()
	var setClass = col.SetClass[int64]()
	setClass.Set()
	setClass.SetWithCollator(collatorClass.Collator())
	var sequence = setClass.SetFromArray([]int64{1, 2, 3})
	var set = setClass.SetFromSequence(sequence)
	ass.Equal(t, sequence.AsArray(), set.AsArray())
}

func TestSetsWithStrings(t *tes.T) {
	var collator = age.CollatorClass[col.SetLike[string]]().Collator()
	var listClass = col.ListClass[string]()
	var empty = []string{}
	var bazbar = listClass.ListFromArray([]string{"baz", "bar"})
	var bazfoo = listClass.ListFromArray([]string{"baz", "foo"})
	var baxbaz = listClass.ListFromArray([]string{"bax", "baz"})
	var baxbez = listClass.ListFromArray([]string{"bax", "bez"})
	var barbaz = listClass.ListFromArray([]string{"bar", "baz"})
	var bar = listClass.ListFromArray([]string{"bar"})
	var setClass = col.SetClass[string]()
	var set = setClass.Set()                                           // [ ]
	ass.True(t, set.IsEmpty())                                         // [ ]
	ass.True(t, set.GetSize() == 0)                                    // [ ]
	ass.False(t, set.ContainsValue("bax"))                             // [ ]
	ass.Equal(t, empty, set.AsArray())                                 // [ ]
	var iterator = set.GetIterator()                                   // [ ]
	ass.False(t, iterator.HasNext())                                   // [ ]
	ass.False(t, iterator.HasPrevious())                               // [ ]
	iterator.ToStart()                                                 // [ ]
	iterator.ToEnd()                                                   // [ ]
	set.RemoveAll()                                                    // [ ]
	set.RemoveValue("foo")                                             // [ ]
	set.AddValue("foo")                                                // ["foo"]
	ass.False(t, set.IsEmpty())                                        // ["foo"]
	ass.True(t, set.GetSize() == 1)                                    // ["foo"]
	ass.Equal(t, "foo", set.GetValue(1))                               // ["foo"]
	ass.True(t, set.GetIndex("baz") == 0)                              // ["foo"]
	ass.True(t, set.ContainsValue("foo"))                              // ["foo"]
	ass.False(t, set.ContainsValue("bax"))                             // ["foo"]
	set.AddValues(bazbar)                                              // ["bar", "baz", "foo"]
	ass.True(t, set.GetSize() == 3)                                    // ["bar", "baz", "foo"]
	ass.True(t, set.GetIndex("baz") == 2)                              // ["bar", "baz", "foo"]
	ass.Equal(t, "bar", set.GetValue(1))                               // ["bar", "baz", "foo"]
	ass.Equal(t, bazfoo.AsArray(), set.GetValues(2, 3).AsArray())      // ["bar", "baz", "foo"]
	ass.Equal(t, bar.AsArray(), set.GetValues(1, 1).AsArray())         // ["bar", "baz", "foo"]
	var set2 = setClass.SetFromSequence(set)                           // ["bar", "baz", "foo"]
	ass.True(t, collator.CompareValues(set, set2))                     // ["bar", "baz", "foo"]
	var array = listClass.ListFromArray([]string{"foo", "bar", "baz"}) // ["bar", "baz", "foo"]
	var set3 = setClass.SetFromSequence(array)                         // ["bar", "baz", "foo"]
	ass.True(t, collator.CompareValues(set2, set3))                    // ["bar", "baz", "foo"]
	iterator = set.GetIterator()                                       // ["bar", "baz", "foo"]
	ass.True(t, iterator.HasNext())                                    // ["bar", "baz", "foo"]
	ass.False(t, iterator.HasPrevious())                               // ["bar", "baz", "foo"]
	ass.Equal(t, "bar", string(iterator.GetNext()))                    // ["bar", "baz", "foo"]
	ass.True(t, iterator.HasPrevious())                                // ["bar", "baz", "foo"]
	iterator.ToEnd()                                                   // ["bar", "baz", "foo"]
	ass.False(t, iterator.HasNext())                                   // ["bar", "baz", "foo"]
	ass.True(t, iterator.HasPrevious())                                // ["bar", "baz", "foo"]
	ass.Equal(t, "foo", string(iterator.GetPrevious()))                // ["bar", "baz", "foo"]
	ass.True(t, iterator.HasNext())                                    // ["bar", "baz", "foo"]
	ass.True(t, set.ContainsValue("baz"))                              // ["bar", "baz", "foo"]
	ass.False(t, set.ContainsValue("bax"))                             // ["bar", "baz", "foo"]
	ass.True(t, set.ContainsAny(baxbaz))                               // ["bar", "baz", "foo"]
	ass.False(t, set.ContainsAny(baxbez))                              // ["bar", "baz", "foo"]
	ass.True(t, set.ContainsAll(barbaz))                               // ["bar", "baz", "foo"]
	ass.False(t, set.ContainsAll(baxbaz))                              // ["bar", "baz", "foo"]
	set.RemoveAll()                                                    // [ ]
	ass.True(t, set.IsEmpty())                                         // [ ]
	ass.True(t, set.GetSize() == 0)                                    // [ ]
}

func TestSetsWithIntegers(t *tes.T) {
	var listClass = col.ListClass[int]()
	var array = listClass.ListFromArray([]int{3, 1, 4, 5, 9, 2})
	var set = col.SetClass[int]().Set() // [ ]
	set.AddValues(array)                // [1,2,3,4,5,9]
	ass.False(t, set.IsEmpty())         // [1,2,3,4,5,9]
	ass.True(t, set.GetSize() == 6)     // [1,2,3,4,5,9]
	ass.True(t, set.GetValue(1) == 1)   // [1,2,3,4,5,9]
	ass.True(t, set.GetValue(-1) == 9)  // [1,2,3,4,5,9]
	set.RemoveValue(6)                  // [1,2,3,4,5,9]
	ass.True(t, set.GetSize() == 6)     // [1,2,3,4,5,9]
	set.RemoveValue(3)                  // [1,2,4,5,9]
	ass.True(t, set.GetSize() == 5)     // [1,2,4,5,9]
	ass.True(t, set.GetValue(3) == 4)   // [1,2,4,5,9]
}

func TestSetsWithTildes(t *tes.T) {
	var listClass = col.ListClass[Integer]()
	var array = listClass.ListFromArray([]Integer{3, 1, 4, 5, 9, 2})
	var set = col.SetClass[Integer]().Set() // [ ]
	set.AddValues(array)                    // [1,2,3,4,5,9]
	ass.False(t, set.IsEmpty())             // [1,2,3,4,5,9]
	ass.True(t, set.GetSize() == 6)         // [1,2,3,4,5,9]
	ass.True(t, set.GetValue(1) == 1)       // [1,2,3,4,5,9]
	ass.True(t, set.GetValue(-1) == 9)      // [1,2,3,4,5,9]
	set.RemoveValue(6)                      // [1,2,3,4,5,9]
	ass.True(t, set.GetSize() == 6)         // [1,2,3,4,5,9]
	set.RemoveValue(3)                      // [1,2,4,5,9]
	ass.True(t, set.GetSize() == 5)         // [1,2,4,5,9]
	ass.True(t, set.GetValue(3) == 4)       // [1,2,4,5,9]
}

func TestSetsWithSets(t *tes.T) {
	var listClass = col.ListClass[int]()
	var array1 = listClass.ListFromArray([]int{3, 1, 4, 5, 9, 2})
	var array2 = listClass.ListFromArray([]int{7, 1, 4, 5, 9, 2})
	var setClass = col.SetClass[int]()
	var set1 = setClass.Set()
	set1.AddValues(array1)
	var set2 = setClass.Set()
	set2.AddValues(array2)
	var set = col.SetClass[col.SetLike[int]]().Set()
	set.AddValue(set1)
	set.AddValue(set2)
	ass.False(t, set.IsEmpty())
	ass.True(t, set.GetSize() == 2)
	ass.Equal(t, set1, set.GetValue(1))
	ass.Equal(t, set2, set.GetValue(-1))
	set.RemoveValue(set1)
	ass.True(t, set.GetSize() == 1)
	set.RemoveAll()
	ass.True(t, set.GetSize() == 0)
}

func TestSetsWithAnd(t *tes.T) {
	var collator = age.CollatorClass[col.SetLike[int]]().Collator()
	var listClass = col.ListClass[int]()
	var array1 = listClass.ListFromArray([]int{3, 1, 2})
	var array2 = listClass.ListFromArray([]int{3, 2, 4})
	var array3 = listClass.ListFromArray([]int{3, 2})
	var setClass = col.SetClass[int]()
	var set1 = setClass.Set()
	set1.AddValues(array1)
	var set2 = setClass.Set()
	set2.AddValues(array2)
	var set3 = setClass.And(set1, set2)
	var set4 = setClass.Set()
	set4.AddValues(array3)
	ass.True(t, collator.CompareValues(set3, set4))
}

func TestSetsWithSans(t *tes.T) {
	var collator = age.CollatorClass[col.SetLike[int]]().Collator()
	var listClass = col.ListClass[int]()
	var array1 = listClass.ListFromArray([]int{3, 1, 2})
	var array2 = listClass.ListFromArray([]int{3, 2, 4})
	var array3 = listClass.ListFromArray([]int{1})
	var setClass = col.SetClass[int]()
	var set1 = setClass.Set()
	set1.AddValues(array1)
	var set2 = setClass.Set()
	set2.AddValues(array2)
	var set3 = setClass.Sans(set1, set2)
	var set4 = setClass.Set()
	set4.AddValues(array3)
	ass.True(t, collator.CompareValues(set3, set4))
}

func TestSetsWithOr(t *tes.T) {
	var collator = age.CollatorClass[col.SetLike[int]]().Collator()
	var listClass = col.ListClass[int]()
	var array1 = listClass.ListFromArray([]int{3, 1, 5})
	var array2 = listClass.ListFromArray([]int{6, 2, 4})
	var array3 = listClass.ListFromArray([]int{1, 3, 5, 6, 2, 4})
	var setClass = col.SetClass[int]()
	var set1 = setClass.Set()
	set1.AddValues(array1)
	var set2 = setClass.Set()
	set2.AddValues(array2)
	var set3 = setClass.Or(set1, set2)
	ass.True(t, set3.ContainsAll(set1))
	ass.True(t, set3.ContainsAll(set2))
	var set4 = setClass.Set()
	set4.AddValues(array3)
	ass.True(t, collator.CompareValues(set3, set4))
}

func TestSetsWithXor(t *tes.T) {
	var collator = age.CollatorClass[col.SetLike[int]]().Collator()
	var listClass = col.ListClass[int]()
	var array1 = listClass.ListFromArray([]int{2, 3, 1, 5})
	var array2 = listClass.ListFromArray([]int{6, 2, 5, 4})
	var array3 = listClass.ListFromArray([]int{1, 3, 4, 6})
	var setClass = col.SetClass[int]()
	var set1 = setClass.Set()
	set1.AddValues(array1)
	var set2 = setClass.Set()
	set2.AddValues(array2)
	var set3 = setClass.Xor(set1, set2)
	var set4 = setClass.Set()
	set4.AddValues(array3)
	ass.True(t, collator.CompareValues(set3, set4))
}

func TestSetsWithEmptySets(t *tes.T) {
	var collator = age.CollatorClass[col.SetLike[int]]().Collator()
	var setClass = col.SetClass[int]()
	var set1 = setClass.Set()
	var set2 = setClass.Set()
	var set3 = setClass.And(set1, set2)
	var set4 = setClass.Sans(set1, set2)
	var set5 = setClass.Or(set1, set2)
	var set6 = setClass.Xor(set1, set2)
	ass.True(t, collator.CompareValues(set3, set4))
	ass.True(t, collator.CompareValues(set4, set5))
	ass.True(t, collator.CompareValues(set5, set6))
	ass.True(t, collator.CompareValues(set6, set1))
}

func TestStackConstructors(t *tes.T) {
	var stackClass = col.StackClass[int64]()
	stackClass.Stack()
	stackClass.StackWithCapacity(5)
	var sequence = stackClass.StackFromArray([]int64{1, 2, 3})
	var stack = stackClass.StackFromSequence(sequence)
	ass.Equal(t, sequence.AsArray(), stack.AsArray())
}

func TestStackWithSmallCapacity(t *tes.T) {
	var stack = col.StackClass[int]().StackWithCapacity(1)
	stack.AddValue(1)
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "Attempted to add a value onto a stack that has reached its capacity.", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	stack.AddValue(2) // This should panic.
}

func TestEmptyStackRemoval(t *tes.T) {
	var stack = col.StackClass[int]().Stack()
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "Attempted to remove a value from an empty stack!", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	stack.RemoveLast() // This should panic.
}

func TestStacksWithStrings(t *tes.T) {
	var stack = col.StackClass[string]().Stack()
	ass.True(t, stack.IsEmpty())
	ass.True(t, stack.GetSize() == 0)
	stack.RemoveAll()
	stack.AddValue("foo")
	stack.AddValue("bar")
	stack.AddValue("baz")
	ass.True(t, stack.GetSize() == 3)
	ass.Equal(t, "baz", string(stack.RemoveLast()))
	ass.True(t, stack.GetSize() == 2)
	ass.Equal(t, "bar", string(stack.RemoveLast()))
	ass.True(t, stack.GetSize() == 1)
	stack.RemoveAll()
}
