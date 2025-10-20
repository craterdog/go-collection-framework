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
	fra "github.com/craterdog/go-collection-framework/v8"
	ass "github.com/stretchr/testify/assert"
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
	fra.ListClass[string]().Concatenate(list, list)
	var association = fra.Association[string, int]("A", 1)
	var catalog = fra.Catalog[string, int]()
	fra.CatalogFromArray[string, int]([]fra.AssociationLike[string, int]{association})
	fra.CatalogFromMap[string, int](catalog.AsMap())
	fra.CatalogFromSequence[string, int](catalog)
	fra.CatalogClass[string, int]().Extract(catalog, list)
	fra.CatalogClass[string, int]().Merge(catalog, catalog)
	fra.Queue[string]()
	fra.QueueWithCapacity[string](8)
	var queue = fra.QueueFromArray[string](list.AsArray())
	fra.QueueFromSequence[string](queue)
	var group fra.Synchronized = new(syn.WaitGroup)
	defer group.Wait()
	var queues = fra.QueueClass[string]().Fork(group, queue, 2)
	fra.QueueClass[string]().Split(group, queue, 2)
	fra.QueueClass[string]().Join(group, queues)
	queue.CloseChannel()
	var set = fra.Set[string]()
	fra.SetWithCollator[string](set.GetCollator())
	fra.SetFromArray[string](set.AsArray())
	fra.SetFromSequence[string](set)
	fra.SetClass[string]().And(set, set)
	fra.SetClass[string]().Ior(set, set)
	fra.SetClass[string]().San(set, set)
	fra.SetClass[string]().Xor(set, set)
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
	fmt.Println("The union of the two sets is:", Set.Ior(set1, set2))
	fmt.Println()

	// Find the logical difference between the two sets.
	fmt.Println("The first set minus the second set is:", Set.San(set1, set2))
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
	var group fra.Synchronized = new(syn.WaitGroup)
	defer group.Wait()

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
	group.Go(func() {
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
	})

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

// AGENTS

// Tilde Types
type Boolean bool
type Byte byte
type Integer int
type String string

type Foolish interface {
	GetFoo() int
	GetBar() string
	GetNil() Foolish
}

func FooBar(foo int, bar string, baz Foolish) Foolish {
	return &foobar{foo, bar, baz}
}

// Encapsulated Type
type foobar struct {
	foo int
	bar string
	Baz Foolish
}

func (v *foobar) GetFoo() int     { return v.foo }
func (v foobar) GetFoo2() int     { return v.foo }
func (v *foobar) GetBar() string  { return v.bar }
func (v foobar) GetBar2() string  { return v.bar }
func (v *foobar) GetNil() Foolish { return nil }
func (v foobar) GetNil2() Foolish { return nil }

// Pure Structure
type Fuz struct {
	Bar string
}

func TestRank(t *tes.T) {
	ass.Equal(t, "LesserRank", fra.LesserRank.String())
	ass.Equal(t, "EqualRank", fra.EqualRank.String())
	ass.Equal(t, "GreaterRank", fra.GreaterRank.String())
}

func TestCompareMaximum(t *tes.T) {
	var collator = fra.CollatorClass[any]().CollatorWithMaximumDepth(1)
	var list = fra.ListClass[any]().ListFromArray([]any{"foo", []int{1, 2, 3}})
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The maximum traversal depth was exceeded: 1", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	collator.CompareValues(list, list)
}

func TestRankMaximum(t *tes.T) {
	var collator = fra.CollatorClass[any]().CollatorWithMaximumDepth(1)
	var list = fra.ListClass[any]().ListFromArray([]any{"foo", []int{1, 2, 3}})
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The maximum traversal depth was exceeded: 1", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	collator.RankValues(list, list)
}

func TestComparison(t *tes.T) {
	var collator = fra.CollatorClass[any]().Collator()

	// Nil
	var ShouldBeNil any

	ass.True(t, collator.CompareValues(nil, nil))
	ass.True(t, collator.CompareValues(nil, ShouldBeNil))
	ass.True(t, collator.CompareValues(ShouldBeNil, ShouldBeNil))
	ass.True(t, collator.CompareValues(ShouldBeNil, nil))

	// Boolean
	var False = false
	var True = true
	var ShouldBeFalse bool

	ass.True(t, collator.CompareValues(ShouldBeFalse, False))
	ass.False(t, collator.CompareValues(True, ShouldBeFalse))

	ass.False(t, collator.CompareValues(False, True))
	ass.True(t, collator.CompareValues(False, False))
	ass.False(t, collator.CompareValues(True, False))
	ass.True(t, collator.CompareValues(True, True))

	// Byte
	var Zero byte = 0x00
	var One byte = 0x01
	var ShouldBeZero byte

	ass.True(t, collator.CompareValues(ShouldBeZero, Zero))
	ass.False(t, collator.CompareValues(One, ShouldBeZero))

	ass.False(t, collator.CompareValues(Zero, One))
	ass.True(t, collator.CompareValues(Zero, Zero))
	ass.False(t, collator.CompareValues(One, Zero))
	ass.True(t, collator.CompareValues(One, One))

	// Integer
	var Zilch = 0
	var Two = 2
	var Three = 3
	var ShouldBeZilch int

	ass.True(t, collator.CompareValues(ShouldBeZilch, Zilch))
	ass.False(t, collator.CompareValues(Two, ShouldBeZilch))

	ass.False(t, collator.CompareValues(Two, Three))
	ass.True(t, collator.CompareValues(Two, Two))
	ass.False(t, collator.CompareValues(Three, Two))
	ass.True(t, collator.CompareValues(Three, Three))

	// Float
	var Negligible = 0.0
	var Fourth = 0.25
	var Half = 0.5
	var ShouldBeNegligible float64

	ass.True(t, collator.CompareValues(ShouldBeNegligible, Negligible))
	ass.False(t, collator.CompareValues(Half, ShouldBeNegligible))

	ass.False(t, collator.CompareValues(Fourth, Half))
	ass.True(t, collator.CompareValues(Fourth, Fourth))
	ass.False(t, collator.CompareValues(Half, Fourth))
	ass.True(t, collator.CompareValues(Half, Half))

	// Complex
	var Origin = 0 + 0i
	var PiOver4 = 1 + 1i
	var PiOver2 = 1 + 0i
	var ShouldBeOrigin complex128

	ass.True(t, collator.CompareValues(ShouldBeOrigin, Origin))
	ass.False(t, collator.CompareValues(PiOver4, ShouldBeOrigin))

	ass.False(t, collator.CompareValues(PiOver4, PiOver2))
	ass.True(t, collator.CompareValues(PiOver4, PiOver4))
	ass.False(t, collator.CompareValues(PiOver2, PiOver4))
	ass.True(t, collator.CompareValues(PiOver2, PiOver2))

	// Rune
	var Null = rune(0)
	var Sad = '☹'
	var Happy = '☺'
	var ShouldBeNull rune

	ass.True(t, collator.CompareValues(ShouldBeNull, Null))
	ass.False(t, collator.CompareValues(Sad, ShouldBeNull))

	ass.False(t, collator.CompareValues(Happy, Sad))
	ass.True(t, collator.CompareValues(Happy, Happy))
	ass.False(t, collator.CompareValues(Sad, Happy))
	ass.True(t, collator.CompareValues(Sad, Sad))

	// String
	var Empty = ""
	var Hello = "Hello"
	var World = "World"
	var ShouldBeEmpty string

	ass.True(t, collator.CompareValues(ShouldBeEmpty, Empty))
	ass.False(t, collator.CompareValues(Hello, ShouldBeEmpty))

	ass.False(t, collator.CompareValues(World, Hello))
	ass.True(t, collator.CompareValues(World, World))
	ass.False(t, collator.CompareValues(Hello, World))
	ass.True(t, collator.CompareValues(Hello, Hello))

	// Array
	var Universe = "Universe"
	var a0 = []any{}
	var a1 = []any{Hello, World}
	var a2 = []any{Hello, Universe}
	var aNil []any

	ass.True(t, collator.CompareValues(aNil, aNil))
	ass.False(t, collator.CompareValues(aNil, a0))
	ass.False(t, collator.CompareValues(a0, aNil))
	ass.True(t, collator.CompareValues(a0, a0))

	ass.False(t, collator.CompareValues(a1, a2))
	ass.True(t, collator.CompareValues(a1, a1))
	ass.False(t, collator.CompareValues(a2, a1))
	ass.True(t, collator.CompareValues(a2, a2))

	// Map
	var m0 = map[any]any{}
	var m1 = map[any]any{
		One: True,
		Two: World}
	var m2 = map[any]any{
		One: True,
		Two: Hello}
	var m3 = map[any]any{
		One: nil,
		Two: Hello}
	var mNil map[any]any

	ass.True(t, collator.CompareValues(mNil, mNil))
	ass.False(t, collator.CompareValues(mNil, m0))
	ass.False(t, collator.CompareValues(m0, mNil))
	ass.True(t, collator.CompareValues(m0, m0))

	ass.False(t, collator.CompareValues(m1, m2))
	ass.True(t, collator.CompareValues(m1, m1))
	ass.False(t, collator.CompareValues(m2, m1))
	ass.True(t, collator.CompareValues(m2, m2))
	ass.False(t, collator.CompareValues(m2, m3))

	// Struct
	var f0 Foolish
	var f1 = FooBar(1, "one", nil)
	var f2 = FooBar(1, "one", nil)
	var f3 = FooBar(2, "two", nil)
	var f4 = Fuz{"two"}
	var f5 = Fuz{"two"}
	var f6 = Fuz{"three"}
	ass.True(t, collator.CompareValues(f0, f0))
	ass.False(t, collator.CompareValues(f0, f1))
	ass.True(t, collator.CompareValues(f1, f1))
	ass.True(t, collator.CompareValues(f1, f2))
	ass.False(t, collator.CompareValues(f2, f3))
	ass.True(t, collator.CompareValues(f4, f4))
	ass.True(t, collator.CompareValues(f4, f5))
	ass.False(t, collator.CompareValues(f5, f6))
	ass.True(t, collator.CompareValues(&f4, &f4))
	ass.True(t, collator.CompareValues(&f4, &f5))
	ass.False(t, collator.CompareValues(&f5, &f6))
}

func TestTildeTypes(t *tes.T) {
	var collator = fra.CollatorClass[any]().Collator()

	// Boolean
	var False = Boolean(false)
	var True = Boolean(true)
	var ShouldBeFalse Boolean

	ass.Equal(t, fra.EqualRank, collator.RankValues(ShouldBeFalse, ShouldBeFalse))
	ass.Equal(t, fra.LesserRank, collator.RankValues(ShouldBeFalse, True))
	ass.Equal(t, fra.EqualRank, collator.RankValues(False, ShouldBeFalse))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(True, ShouldBeFalse))
	ass.Equal(t, fra.EqualRank, collator.RankValues(ShouldBeFalse, False))
	ass.Equal(t, fra.LesserRank, collator.RankValues(False, True))
	ass.Equal(t, fra.EqualRank, collator.RankValues(False, False))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(True, False))
	ass.Equal(t, fra.EqualRank, collator.RankValues(True, True))

	// Byte
	var Zero = Byte(0)
	var One = Byte(1)
	var TFF = Byte(255)
	var ShouldBeZero Byte

	ass.True(t, collator.CompareValues(ShouldBeZero, Zero))
	ass.False(t, collator.CompareValues(One, ShouldBeZero))

	ass.False(t, collator.CompareValues(One, TFF))
	ass.True(t, collator.CompareValues(One, One))
	ass.False(t, collator.CompareValues(TFF, One))
	ass.True(t, collator.CompareValues(TFF, TFF))

	// Integer
	var Zilch = Integer(0)
	var Two = Integer(2)
	var Three = Integer(3)
	var ShouldBeZilch Integer

	ass.True(t, collator.CompareValues(ShouldBeZilch, Zilch))
	ass.False(t, collator.CompareValues(Two, ShouldBeZilch))

	ass.False(t, collator.CompareValues(Two, Three))
	ass.True(t, collator.CompareValues(Two, Two))
	ass.False(t, collator.CompareValues(Three, Two))
	ass.True(t, collator.CompareValues(Three, Three))

	// String
	var Empty = String("")
	var Hello = String("Hello")
	var World = String("World")
	var ShouldBeEmpty String

	ass.True(t, collator.CompareValues(ShouldBeEmpty, Empty))
	ass.False(t, collator.CompareValues(Hello, ShouldBeEmpty))

	ass.False(t, collator.CompareValues(World, Hello))
	ass.True(t, collator.CompareValues(World, World))
	ass.False(t, collator.CompareValues(Hello, World))
	ass.True(t, collator.CompareValues(Hello, Hello))
}

func TestCompareRecursiveArrays(t *tes.T) {
	var collator = fra.CollatorClass[any]().Collator()
	var list = fra.ListClass[any]().ListFromArray(
		[]any{0},
	)
	list.SetValue(1, list) // Now it is recursive.
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The maximum traversal depth was exceeded: 16", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	collator.CompareValues(list, list) // This should panic.
}

func TestCompareRecursiveMaps(t *tes.T) {
	var collator = fra.CollatorClass[any]().Collator()
	var catalog = fra.CatalogClass[string, any]().CatalogFromMap(
		map[string]any{
			"first": 1,
		},
	)
	catalog.SetValue("first", catalog) // Now it is recursive.
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The maximum traversal depth was exceeded: 16", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	collator.CompareValues(catalog, catalog) // This should panic.
}

func TestNilRanking(t *tes.T) {
	var collator = fra.CollatorClass[any]().Collator()
	var ShouldBeNil any
	ass.Equal(t, fra.EqualRank, collator.RankValues(nil, nil))
	ass.Equal(t, fra.EqualRank, collator.RankValues(nil, ShouldBeNil))
	ass.Equal(t, fra.EqualRank, collator.RankValues(ShouldBeNil, ShouldBeNil))
	ass.Equal(t, fra.EqualRank, collator.RankValues(ShouldBeNil, nil))
}

func TestBooleanRanking(t *tes.T) {
	var collator = fra.CollatorClass[any]().Collator()
	var False = false
	var True = true
	var ShouldBeFalse bool
	ass.Equal(t, fra.EqualRank, collator.RankValues(ShouldBeFalse, ShouldBeFalse))
	ass.Equal(t, fra.LesserRank, collator.RankValues(ShouldBeFalse, True))
	ass.Equal(t, fra.EqualRank, collator.RankValues(False, ShouldBeFalse))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(True, ShouldBeFalse))
	ass.Equal(t, fra.EqualRank, collator.RankValues(ShouldBeFalse, False))
	ass.Equal(t, fra.LesserRank, collator.RankValues(False, True))
	ass.Equal(t, fra.EqualRank, collator.RankValues(False, False))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(True, False))
	ass.Equal(t, fra.EqualRank, collator.RankValues(True, True))
}

func TestByteRanking(t *tes.T) {
	var collator = fra.CollatorClass[any]().Collator()
	var Zero byte = 0x00
	var One byte = 0x01
	var ShouldBeZero byte
	ass.Equal(t, fra.EqualRank, collator.RankValues(ShouldBeZero, ShouldBeZero))
	ass.Equal(t, fra.LesserRank, collator.RankValues(ShouldBeZero, One))
	ass.Equal(t, fra.EqualRank, collator.RankValues(Zero, ShouldBeZero))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(One, ShouldBeZero))
	ass.Equal(t, fra.EqualRank, collator.RankValues(ShouldBeZero, Zero))
	ass.Equal(t, fra.LesserRank, collator.RankValues(Zero, One))
	ass.Equal(t, fra.EqualRank, collator.RankValues(Zero, Zero))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(One, Zero))
	ass.Equal(t, fra.EqualRank, collator.RankValues(One, One))
}

func TestIntegerRanking(t *tes.T) {
	var collator = fra.CollatorClass[any]().Collator()
	var Zilch = 0
	var Two = 2
	var Three = 3
	var ShouldBeZilch int
	ass.Equal(t, fra.EqualRank, collator.RankValues(ShouldBeZilch, ShouldBeZilch))
	ass.Equal(t, fra.LesserRank, collator.RankValues(ShouldBeZilch, Two))
	ass.Equal(t, fra.EqualRank, collator.RankValues(Zilch, ShouldBeZilch))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(Two, ShouldBeZilch))
	ass.Equal(t, fra.EqualRank, collator.RankValues(ShouldBeZilch, Zilch))
	ass.Equal(t, fra.LesserRank, collator.RankValues(Two, Three))
	ass.Equal(t, fra.EqualRank, collator.RankValues(Two, Two))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(Three, Two))
	ass.Equal(t, fra.EqualRank, collator.RankValues(Three, Three))
}

func TestFloatRanking(t *tes.T) {
	var collator = fra.CollatorClass[any]().Collator()
	var Negligible = 0.0
	var Fourth = 0.25
	var Half = 0.5
	var ShouldBeNegligible float64
	ass.Equal(t, fra.EqualRank, collator.RankValues(ShouldBeNegligible, ShouldBeNegligible))
	ass.Equal(t, fra.LesserRank, collator.RankValues(ShouldBeNegligible, Half))
	ass.Equal(t, fra.EqualRank, collator.RankValues(Negligible, ShouldBeNegligible))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(Half, ShouldBeNegligible))
	ass.Equal(t, fra.EqualRank, collator.RankValues(ShouldBeNegligible, Negligible))
	ass.Equal(t, fra.LesserRank, collator.RankValues(Fourth, Half))
	ass.Equal(t, fra.EqualRank, collator.RankValues(Fourth, Fourth))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(Half, Fourth))
	ass.Equal(t, fra.EqualRank, collator.RankValues(Half, Half))
}

func TestComplexRanking(t *tes.T) {
	var collator = fra.CollatorClass[any]().Collator()
	var Origin complex128
	var One = 1 + 0i
	var Pi = -1 + 0i
	var PiOver2 = 0 + 1i
	var PiOver4 = 1 + 1i

	ass.Equal(t, fra.EqualRank, collator.RankValues(Origin, Origin))
	ass.Equal(t, fra.LesserRank, collator.RankValues(Origin, One))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(Origin, Pi))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(Origin, PiOver2))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(Origin, PiOver4))

	ass.Equal(t, fra.GreaterRank, collator.RankValues(One, Origin))
	ass.Equal(t, fra.EqualRank, collator.RankValues(One, One))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(One, Pi))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(One, PiOver2))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(One, PiOver4))

	ass.Equal(t, fra.GreaterRank, collator.RankValues(Pi, Origin))
	ass.Equal(t, fra.LesserRank, collator.RankValues(Pi, One))
	ass.Equal(t, fra.EqualRank, collator.RankValues(Pi, Pi))
	ass.Equal(t, fra.LesserRank, collator.RankValues(Pi, PiOver2))
	ass.Equal(t, fra.LesserRank, collator.RankValues(Pi, PiOver4))

	ass.Equal(t, fra.GreaterRank, collator.RankValues(PiOver2, Origin))
	ass.Equal(t, fra.LesserRank, collator.RankValues(PiOver2, One))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(PiOver2, Pi))
	ass.Equal(t, fra.EqualRank, collator.RankValues(PiOver2, PiOver2))
	ass.Equal(t, fra.LesserRank, collator.RankValues(PiOver2, PiOver4))

	ass.Equal(t, fra.GreaterRank, collator.RankValues(PiOver4, Origin))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(PiOver4, One))
	ass.Equal(t, fra.LesserRank, collator.RankValues(PiOver4, Pi))
	ass.Equal(t, fra.LesserRank, collator.RankValues(PiOver4, PiOver2))
	ass.Equal(t, fra.EqualRank, collator.RankValues(PiOver4, PiOver4))
}

func TestRuneRanking(t *tes.T) {
	var collator = fra.CollatorClass[any]().Collator()
	var Null = rune(0)
	var Sad = '☹'
	var Happy = '☺'
	var ShouldBeNull rune
	ass.Equal(t, fra.EqualRank, collator.RankValues(ShouldBeNull, ShouldBeNull))
	ass.Equal(t, fra.LesserRank, collator.RankValues(ShouldBeNull, Sad))
	ass.Equal(t, fra.EqualRank, collator.RankValues(Null, ShouldBeNull))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(Sad, ShouldBeNull))
	ass.Equal(t, fra.EqualRank, collator.RankValues(ShouldBeNull, Null))
	ass.Equal(t, fra.LesserRank, collator.RankValues(Sad, Happy))
	ass.Equal(t, fra.EqualRank, collator.RankValues(Sad, Sad))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(Happy, Sad))
	ass.Equal(t, fra.EqualRank, collator.RankValues(Happy, Happy))
}

func TestStringRanking(t *tes.T) {
	var collator = fra.CollatorClass[any]().Collator()
	var Empty = ""
	var Hello = "Hello"
	var World = "World"
	var ShouldBeEmpty string
	ass.Equal(t, fra.EqualRank, collator.RankValues(ShouldBeEmpty, ShouldBeEmpty))
	ass.Equal(t, fra.LesserRank, collator.RankValues(ShouldBeEmpty, Hello))
	ass.Equal(t, fra.EqualRank, collator.RankValues(Empty, ShouldBeEmpty))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(Hello, ShouldBeEmpty))
	ass.Equal(t, fra.EqualRank, collator.RankValues(ShouldBeEmpty, Empty))
	ass.Equal(t, fra.LesserRank, collator.RankValues(Hello, World))
	ass.Equal(t, fra.EqualRank, collator.RankValues(Hello, Hello))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(World, Hello))
	ass.Equal(t, fra.EqualRank, collator.RankValues(World, World))
}

func TestArrayRanking(t *tes.T) {
	var collator = fra.CollatorClass[any]().Collator()
	var Hello = "Hello"
	var World = "World"
	var Universe = "Universe"
	var a0 = []any{}
	var a1 = []any{Hello, World}
	var a2 = []any{Hello, Universe}
	var a3 = []any{Hello, World, Universe}
	var a4 = []any{Hello, Universe, World}
	var aNil []any
	ass.Equal(t, fra.EqualRank, collator.RankValues(aNil, aNil))
	ass.Equal(t, fra.LesserRank, collator.RankValues(aNil, a0))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(a0, aNil))
	ass.Equal(t, fra.EqualRank, collator.RankValues(a0, a0))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(a1, aNil))
	ass.Equal(t, fra.LesserRank, collator.RankValues(a2, a1))
	ass.Equal(t, fra.EqualRank, collator.RankValues(a2, a2))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(a1, a2))
	ass.Equal(t, fra.EqualRank, collator.RankValues(a1, a1))
	ass.Equal(t, fra.LesserRank, collator.RankValues(a2, a3))
	ass.Equal(t, fra.EqualRank, collator.RankValues(a2, a2))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(a3, a2))
	ass.Equal(t, fra.EqualRank, collator.RankValues(a3, a3))
	ass.Equal(t, fra.LesserRank, collator.RankValues(a4, a1))
	ass.Equal(t, fra.EqualRank, collator.RankValues(a4, a4))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(a1, a4))
	ass.Equal(t, fra.EqualRank, collator.RankValues(a1, a1))
}

func TestMapRanking(t *tes.T) {
	var collator = fra.CollatorClass[any]().Collator()
	var True = true
	var One byte = 0x01
	var Two = 2
	var Three = "three"
	var Hello = "Hello"
	var World = "World"
	var Universe = "Universe"
	var m0 = map[any]any{}
	var m1 = map[any]any{
		One: True,
		Two: World}
	var m2 = map[any]any{
		One: True,
		Two: Hello}
	var m3 = map[any]any{
		One:   True,
		Two:   World,
		Three: Universe}
	var m4 = map[any]any{
		One:   True,
		Two:   Universe,
		Three: World}
	var mNil map[any]any
	ass.Equal(t, fra.EqualRank, collator.RankValues(mNil, mNil))
	ass.Equal(t, fra.LesserRank, collator.RankValues(mNil, m0))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(m0, mNil))
	ass.Equal(t, fra.EqualRank, collator.RankValues(m0, m0))
	ass.Equal(t, fra.LesserRank, collator.RankValues(m2, m1))
	ass.Equal(t, fra.EqualRank, collator.RankValues(m2, m2))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(m1, m2))
	ass.Equal(t, fra.EqualRank, collator.RankValues(m1, m1))
	ass.Equal(t, fra.LesserRank, collator.RankValues(m2, m3))
	ass.Equal(t, fra.EqualRank, collator.RankValues(m2, m2))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(m3, m2))
	ass.Equal(t, fra.EqualRank, collator.RankValues(m3, m3))
	ass.Equal(t, fra.LesserRank, collator.RankValues(m4, m1))
	ass.Equal(t, fra.EqualRank, collator.RankValues(m4, m4))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(m1, m4))
	ass.Equal(t, fra.EqualRank, collator.RankValues(m1, m1))
}

func TestStructRanking(t *tes.T) {
	var collator = fra.CollatorClass[any]().Collator()
	var f1 = FooBar(1, "one", nil)
	var f2 = FooBar(1, "two", nil)
	var f3 = FooBar(2, "two", nil)
	var f4 = Fuz{"two"}
	var f5 = Fuz{"two"}
	var f6 = Fuz{"three"}
	ass.Equal(t, fra.EqualRank, collator.RankValues(f1, f1))
	ass.Equal(t, fra.LesserRank, collator.RankValues(f1, f2))
	ass.Equal(t, fra.LesserRank, collator.RankValues(f2, f3))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(f3, f1))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(f3, f2))
	ass.Equal(t, fra.EqualRank, collator.RankValues(f4, f4))
	ass.Equal(t, fra.EqualRank, collator.RankValues(f4, f5))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(f5, f6))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(f3, &f4))
	ass.Equal(t, fra.EqualRank, collator.RankValues(&f4, &f4))
	ass.Equal(t, fra.EqualRank, collator.RankValues(&f4, &f5))
	ass.Equal(t, fra.GreaterRank, collator.RankValues(&f5, &f6))
}

func TestTildeArrays(t *tes.T) {
	var collator = fra.CollatorClass[String]().Collator()
	var ranker = collator.RankValues
	var sorter = fra.SorterClass[String]().SorterWithRanker(ranker)
	var alpha = String("alpha")
	var beta = String("beta")
	var gamma = String("gamma")
	var delta = String("delta")
	var array = []String{alpha, beta, gamma, delta}
	sorter.SortValues(array)
	ass.Equal(t, alpha, array[0])
	ass.Equal(t, beta, array[1])
	ass.Equal(t, delta, array[2])
	ass.Equal(t, gamma, array[3])
}

func TestRankRecursiveArrays(t *tes.T) {
	var collator = fra.CollatorClass[any]().Collator()
	var list = fra.ListClass[any]().ListFromArray(
		[]any{0},
	)
	list.SetValue(1, list) // Now it is recursive.
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The maximum traversal depth was exceeded: 16", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	collator.RankValues(list, list) // This should panic.
}

func TestRankRecursiveMaps(t *tes.T) {
	var collator = fra.CollatorClass[any]().Collator()
	var catalog = fra.CatalogClass[string, any]().CatalogFromMap(
		map[string]any{
			"first": 1,
		},
	)
	catalog.SetValue("first", catalog) // Now it is recursive.
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The maximum traversal depth was exceeded: 16", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	collator.RankValues(catalog, catalog) // This should panic.
}

func TestIteratorsWithLists(t *tes.T) {
	var list = fra.ListClass[int]().ListFromArray([]int{1, 2, 3, 4, 5})
	list = fra.ListClass[int]().ListFromSequence(list)
	var iterator = list.GetIterator()
	ass.False(t, iterator.IsEmpty())
	ass.True(t, iterator.GetSize() == 5)
	ass.True(t, iterator.GetSlot() == 0)
	ass.False(t, iterator.HasPrevious())
	ass.True(t, iterator.HasNext())
	ass.Equal(t, 1, iterator.GetNext())
	ass.True(t, iterator.HasPrevious())
	ass.True(t, iterator.HasNext())
	ass.Equal(t, 1, iterator.GetPrevious())
	iterator.SetSlot(2)
	ass.True(t, iterator.HasPrevious())
	ass.True(t, iterator.HasNext())
	ass.Equal(t, 3, iterator.GetNext())
	ass.False(t, iterator.IsEmpty())
	ass.True(t, iterator.GetSize() == 5)
	ass.True(t, iterator.GetSlot() == 3)
	iterator.ToEnd()
	ass.True(t, iterator.HasPrevious())
	ass.False(t, iterator.HasNext())
	ass.Equal(t, 5, iterator.GetPrevious())
	iterator.ToStart()
	ass.False(t, iterator.HasPrevious())
	ass.True(t, iterator.HasNext())
	ass.Equal(t, 1, iterator.GetNext())
}

func TestSortingEmpty(t *tes.T) {
	var collator = fra.CollatorClass[any]().Collator()
	var ranker = collator.RankValues
	var sorter = fra.SorterClass[any]().SorterWithRanker(ranker)
	var empty = []any{}
	sorter.SortValues(empty)
}

func TestSortingIntegers(t *tes.T) {
	var collator = fra.CollatorClass[int]().Collator()
	var ranker = collator.RankValues
	var sorter = fra.SorterClass[int]().SorterWithRanker(ranker)
	var unsorted = []int{4, 3, 1, 5, 2}
	var sorted = []int{1, 2, 3, 4, 5}
	sorter.SortValues(unsorted)
	ass.Equal(t, sorted, unsorted)
}

func TestSortingStrings(t *tes.T) {
	var collator = fra.CollatorClass[string]().Collator()
	var ranker = collator.RankValues
	var sorter = fra.SorterClass[string]().SorterWithRanker(ranker)
	var unsorted = []string{"alpha", "beta", "gamma", "delta"}
	var sorted = []string{"alpha", "beta", "delta", "gamma"}
	sorter.SortValues(unsorted)
	ass.Equal(t, sorted, unsorted)
}

var (
	invalid fra.State = fra.ControllerClass().Invalid()
	state1  fra.State = "$State1"
	state2  fra.State = "$State2"
	state3  fra.State = "$State3"
)

var (
	initialized fra.Event = "$Initialized"
	processed   fra.Event = "$Processed"
	finalized   fra.Event = "$Finalized"
)

func TestController(t *tes.T) {
	var events = []fra.Event{initialized, processed, finalized}
	var transitions = map[fra.State]fra.Transitions{
		state1: fra.Transitions{state2, invalid, invalid},
		state2: fra.Transitions{invalid, state2, state3},
		state3: fra.Transitions{invalid, invalid, invalid},
	}

	var controller = fra.Controller(events, transitions, state1)
	ass.Equal(t, state1, controller.GetState())
	ass.Equal(t, state2, controller.ProcessEvent(initialized))
	ass.Equal(t, state2, controller.ProcessEvent(processed))
	ass.Equal(t, state3, controller.ProcessEvent(finalized))
	controller.SetState(state1)
	ass.Equal(t, state1, controller.GetState())
}

// COLLECTIONS

func TestCatalogConstructors(t *tes.T) {
	var class = fra.CatalogClass[rune, int64]()
	class.Catalog()
	class.CatalogFromArray([]fra.AssociationLike[rune, int64]{})
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
	var catalogCollator = fra.Collator[fra.CatalogLike[string, int]]()
	var keys = fra.ListClass[string]().ListFromArray([]string{"foo", "bar"})
	var association1 = fra.Association("foo", 1)
	var association2 = fra.Association("bar", 2)
	var association3 = fra.Association("baz", 3)
	var catalog = fra.Catalog[string, int]()
	ass.True(t, catalog.IsEmpty())
	ass.True(t, catalog.GetSize() == 0)
	ass.Equal(t, []string{}, catalog.GetKeys().AsArray())
	ass.Equal(t, []fra.AssociationLike[string, int]{}, catalog.AsArray())
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
	var catalog2 = fra.CatalogFromSequence(catalog)
	ass.True(t, catalogCollator.CompareValues(catalog, catalog2))
	var m = fra.CatalogClass[string, int]().CatalogFromMap(map[string]int{
		"foo": 1,
		"bar": 2,
		"baz": 3,
	})
	var associationCollator = fra.Collator[fra.AssociationLike[string, int]]()
	var catalog3 = fra.CatalogFromSequence(m)
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
	ass.Equal(t, 3, catalog.GetValue("baz"))
	catalog.SetValue("bar", 5)
	ass.Equal(t, []int{1, 5}, catalog.GetValues(keys).AsArray())
	catalog.SortValues()
	ass.Equal(t, []string{"bar", "baz", "foo"}, catalog.GetKeys().AsArray())
	catalog.ReverseValues()
	ass.Equal(t, []string{"foo", "baz", "bar"}, catalog.GetKeys().AsArray())
	catalog.ReverseValues()
	ass.Equal(t, []int{1, 5}, catalog.RemoveValues(keys).AsArray())
	ass.True(t, catalog.GetSize() == 1)
	ass.Equal(t, 3, catalog.RemoveValue("baz"))
	ass.True(t, catalog.IsEmpty())
	ass.True(t, catalog.GetSize() == 0)
	catalog.RemoveAll()
	ass.True(t, catalog.IsEmpty())
	ass.True(t, catalog.GetSize() == 0)
}

func TestCatalogsWithMerge(t *tes.T) {
	var collator = fra.Collator[fra.CatalogLike[string, int]]()
	var association1 = fra.Association("foo", 1)
	var association2 = fra.Association("bar", 2)
	var association3 = fra.Association("baz", 3)
	var catalog1 = fra.Catalog[string, int]()
	catalog1.SetValue(association1.GetKey(), association1.GetValue())
	catalog1.SetValue(association2.GetKey(), association2.GetValue())
	var catalog2 = fra.Catalog[string, int]()
	catalog2.SetValue(association2.GetKey(), association2.GetValue())
	catalog2.SetValue(association3.GetKey(), association3.GetValue())
	var catalog3 = fra.CatalogClass[string, int]().Merge(catalog1, catalog2)
	var catalog4 = fra.Catalog[string, int]()
	catalog4.SetValue(association1.GetKey(), association1.GetValue())
	catalog4.SetValue(association2.GetKey(), association2.GetValue())
	catalog4.SetValue(association3.GetKey(), association3.GetValue())
	ass.True(t, collator.CompareValues(catalog3, catalog4))
}

func TestCatalogsWithExtract(t *tes.T) {
	var keys = fra.ListClass[string]().ListFromArray([]string{"foo", "baz"})
	var association1 = fra.Association("foo", 1)
	var association2 = fra.Association("bar", 2)
	var association3 = fra.Association("baz", 3)
	var catalog1 = fra.Catalog[string, int]()
	catalog1.SetValue(association1.GetKey(), association1.GetValue())
	catalog1.SetValue(association2.GetKey(), association2.GetValue())
	catalog1.SetValue(association3.GetKey(), association3.GetValue())
	var catalog2 = fra.CatalogClass[string, int]().Extract(catalog1, keys)
	var catalog3 = fra.Catalog[string, int]()
	catalog3.SetValue(association1.GetKey(), association1.GetValue())
	catalog3.SetValue(association3.GetKey(), association3.GetValue())
	var collator = fra.Collator[fra.CatalogLike[string, int]]()
	ass.True(t, collator.CompareValues(catalog2, catalog3))
	var catalog4 = fra.CatalogFromArray([]fra.AssociationLike[string, int]{
		association1,
		association2,
		association3,
	})
	ass.True(t, collator.CompareValues(catalog1, catalog4))
}

func TestCatalogsWithEmptyCatalogs(t *tes.T) {
	var keys = fra.ListClass[int]().List()
	var catalog1 = fra.Catalog[int, string]()
	var catalog2 = fra.Catalog[int, string]()
	var catalog3 = fra.CatalogClass[int, string]().Merge(catalog1, catalog2)
	var catalog4 = fra.CatalogClass[int, string]().Extract(catalog1, keys)
	var collator = fra.Collator[fra.CatalogLike[int, string]]()
	ass.True(t, collator.CompareValues(catalog1, catalog2))
	ass.True(t, collator.CompareValues(catalog2, catalog3))
	ass.True(t, collator.CompareValues(catalog3, catalog4))
	ass.True(t, collator.CompareValues(catalog4, catalog1))
}

func TestListConstructors(t *tes.T) {
	fra.List[int64]()
	var sequence = fra.ListFromArray([]int64{1, 2, 3})
	var list = fra.ListFromSequence(sequence)
	ass.Equal(t, sequence.AsArray(), list.AsArray())
}

func TestListsWithStrings(t *tes.T) {
	var collator = fra.CollatorClass[fra.ListLike[string]]().Collator()
	var foo = fra.ListFromArray([]string{"foo"})
	var bar = fra.ListFromArray([]string{"bar"})
	var baz = fra.ListFromArray([]string{"baz"})
	var foz = fra.ListFromArray([]string{"foz"})
	var barbaz = fra.ListFromArray([]string{"bar", "baz"})
	var bazbaz = fra.ListFromArray([]string{"baz", "baz"})
	var foobar = fra.ListFromArray([]string{"foo", "bar"})
	var baxbaz = fra.ListFromArray([]string{"bax", "baz"})
	var baxbez = fra.ListFromArray([]string{"bax", "bez"})
	var barfoobax = fra.ListFromArray([]string{"bar", "foo", "bax"})
	var foobazbar = fra.ListFromArray([]string{"foo", "baz", "bar"})
	var foobarbaz = fra.ListFromArray([]string{"foo", "bar", "baz"})
	var barbazfoo = fra.ListFromArray([]string{"bar", "baz", "foo"})
	var list = fra.List[string]()
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
	ass.True(t, collator.CompareValues(fra.ListFromArray(list.AsArray()), list))
	ass.Equal(t, barbaz.AsArray(), list.GetValues(2, 3).AsArray())
	ass.Equal(t, foo.AsArray(), list.GetValues(1, 1).AsArray())
	var list2 = fra.ListFromSequence(list)
	ass.True(t, collator.CompareValues(list, list2))
	var array = fra.ListFromArray([]string{"foo", "bar", "baz"})
	var list3 = fra.ListFromSequence(array)
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
	var array = fra.ListFromArray([]Integer{3, 1, 4, 5, 9, 2})
	var list = fra.ListFromSequence(array)
	ass.False(t, list.IsEmpty())        // [3,1,4,5,9,2]
	ass.True(t, list.GetSize() == 6)    // [3,1,4,5,9,2]
	ass.True(t, list.GetValue(1) == 3)  // [3,1,4,5,9,2]
	ass.True(t, list.GetValue(-1) == 2) // [3,1,4,5,9,2]
	list.SortValues()                   // [1,2,3,4,5,9]
	ass.True(t, list.GetSize() == 6)    // [1,2,3,4,5,9]
	ass.True(t, list.GetValue(3) == 3)  // [1,2,3,4,5,9]
}

func TestListsWithConcatenate(t *tes.T) {
	var collator = fra.CollatorClass[fra.ListLike[int]]().Collator()
	var onetwothree = fra.ListFromArray([]int{1, 2, 3})
	var fourfivesix = fra.ListFromArray([]int{4, 5, 6})
	var onethrusix = fra.ListFromArray([]int{1, 2, 3, 4, 5, 6})
	var list1 = fra.List[int]()
	list1.AppendValues(onetwothree)
	var list2 = fra.List[int]()
	list2.AppendValues(fourfivesix)
	var list3 = fra.ListClass[int]().Concatenate(list1, list2)
	var list4 = fra.List[int]()
	list4.AppendValues(onethrusix)
	ass.True(t, collator.CompareValues(list3, list4))
}

func TestListsWithEmptyLists(t *tes.T) {
	var collator = fra.Collator[fra.ListLike[int]]()
	var empty = fra.List[int]()
	var list = fra.ListClass[int]().Concatenate(empty, empty)
	ass.True(t, collator.CompareValues(empty, empty))
	ass.True(t, collator.CompareValues(list, empty))
	ass.True(t, collator.CompareValues(empty, list))
	ass.True(t, collator.CompareValues(list, list))
}

func TestQueueConstructors(t *tes.T) {
	fra.Queue[int64]()
	fra.QueueWithCapacity[int64](5)
	var sequence = fra.QueueFromArray([]int64{1, 2, 3})
	var queue = fra.QueueFromSequence(sequence)
	ass.Equal(t, sequence.AsArray(), queue.AsArray())
}

func TestQueueWithConcurrency(t *tes.T) {
	// Create a wait group for synchronization.
	var group fra.Synchronized = new(syn.WaitGroup)
	defer group.Wait()

	// Create a new queue with a specific capacity.
	var queue = fra.QueueWithCapacity[int](12)
	ass.True(t, queue.GetCapacity() == 12)
	ass.True(t, queue.IsEmpty())
	ass.True(t, queue.GetSize() == 0)

	// Add some values to the queue.
	for i := 1; i < 10; i++ {
		queue.AddValue(i)
	}
	ass.True(t, queue.GetSize() == 9)

	// Remove values from the queue in the background.
	group.Go(func() {
		var value int
		var ok = true
		for i := 1; ok; i++ {
			value, ok = queue.RemoveFirst()
			if ok {
				ass.Equal(t, i, value)
			}
		}
		queue.RemoveAll()
	})

	// Add some more values to the queue.
	for i := 10; i < 101; i++ {
		queue.AddValue(i)
	}
	queue.CloseChannel()
}

func TestQueueWithFork(t *tes.T) {
	// Create a wait group for synchronization.
	var group fra.Synchronized = new(syn.WaitGroup)
	defer group.Wait()

	// Create a new queue with a fan out of two.
	var input = fra.QueueWithCapacity[int](3)
	var outputs = fra.QueueClass[int]().Fork(group, input, 2)

	// Remove values from the output queues in the background.
	var iterator = outputs.GetIterator()
	for iterator.HasNext() {
		var output = iterator.GetNext()
		group.Go(func() {
			var value int
			var ok = true
			for i := 1; ok; i++ {
				value, ok = output.RemoveFirst()
				if ok {
					ass.Equal(t, i, value)
				}
			}
		})
	}

	// Add values to the input queue.
	for i := 1; i < 11; i++ {
		input.AddValue(i)
	}
	input.CloseChannel()
}

func TestQueueWithInvalidFanOut(t *tes.T) {
	// Create a wait group for synchronization.
	var group fra.Synchronized = new(syn.WaitGroup)
	defer group.Wait()

	// Create a new queue with an invalid fan out.
	var input = fra.QueueWithCapacity[int](3)
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The fan out size for a queue must be greater than one.", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	fra.QueueClass[int]().Fork(group, input, 1) // Should panic here.
}

func TestQueueWithSplitAndJoin(t *tes.T) {
	// Create a wait group for synchronization.
	var group fra.Synchronized = new(syn.WaitGroup)
	defer group.Wait()

	// Create a new queue with a split of five outputs and a join back to one.
	var input = fra.QueueWithCapacity[int](3)
	var split = fra.QueueClass[int]().Split(group, input, 5)
	var output = fra.QueueClass[int]().Join(group, split)

	// Remove values from the output queue in the background.
	group.Go(func() {
		var value int
		var ok = true
		for i := 1; ok; i++ {
			value, ok = output.RemoveFirst()
			if ok {
				ass.Equal(t, i, value)
			}
		}
	})

	// Add values to the input queue.
	for i := 1; i < 21; i++ {
		input.AddValue(i)
	}
	input.CloseChannel()
}

func TestQueueWithInvalidSplit(t *tes.T) {
	// Create a wait group for synchronization.
	var group fra.Synchronized = new(syn.WaitGroup)
	defer group.Wait()

	// Create a new queue with an invalid fan out.
	var input = fra.QueueWithCapacity[int](3)
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The size of the split must be greater than one.", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	fra.QueueClass[int]().Split(group, input, 1) // Should panic here.
}

func TestQueueWithInvalidJoin(t *tes.T) {
	// Create a wait group for synchronization.
	var group fra.Synchronized = new(syn.WaitGroup)
	defer group.Wait()

	// Create a new queue with an invalid fan out.
	var inputs = fra.List[fra.QueueLike[int]]()
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The number of input queues for a join must be at least one.", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	fra.QueueClass[int]().Join(group, inputs) // Should panic here.
}

func TestSetConstructors(t *tes.T) {
	var collator = fra.Collator[int64]()
	fra.Set[int64]()
	fra.SetWithCollator(collator)
	var sequence = fra.SetFromArray([]int64{1, 2, 3})
	var set = fra.SetFromSequence(sequence)
	ass.Equal(t, sequence.AsArray(), set.AsArray())
}

func TestSetsWithStrings(t *tes.T) {
	var collator = fra.Collator[fra.SetLike[string]]()
	fra.List[string]()
	var empty = []string{}
	var bazbar = fra.ListFromArray([]string{"baz", "bar"})
	var bazfoo = fra.ListFromArray([]string{"baz", "foo"})
	var baxbaz = fra.ListFromArray([]string{"bax", "baz"})
	var baxbez = fra.ListFromArray([]string{"bax", "bez"})
	var barbaz = fra.ListFromArray([]string{"bar", "baz"})
	var bar = fra.ListFromArray([]string{"bar"})
	var set = fra.Set[string]()                                   // [ ]
	ass.True(t, set.IsEmpty())                                    // [ ]
	ass.True(t, set.GetSize() == 0)                               // [ ]
	ass.False(t, set.ContainsValue("bax"))                        // [ ]
	ass.Equal(t, empty, set.AsArray())                            // [ ]
	var iterator = set.GetIterator()                              // [ ]
	ass.False(t, iterator.HasNext())                              // [ ]
	ass.False(t, iterator.HasPrevious())                          // [ ]
	iterator.ToStart()                                            // [ ]
	iterator.ToEnd()                                              // [ ]
	set.RemoveAll()                                               // [ ]
	set.RemoveValue("foo")                                        // [ ]
	set.AddValue("foo")                                           // ["foo"]
	ass.False(t, set.IsEmpty())                                   // ["foo"]
	ass.True(t, set.GetSize() == 1)                               // ["foo"]
	ass.Equal(t, "foo", set.GetValue(1))                          // ["foo"]
	ass.True(t, set.GetIndex("baz") == 0)                         // ["foo"]
	ass.True(t, set.ContainsValue("foo"))                         // ["foo"]
	ass.False(t, set.ContainsValue("bax"))                        // ["foo"]
	set.AddValues(bazbar)                                         // ["bar", "baz", "foo"]
	ass.True(t, set.GetSize() == 3)                               // ["bar", "baz", "foo"]
	ass.True(t, set.GetIndex("baz") == 2)                         // ["bar", "baz", "foo"]
	ass.Equal(t, "bar", set.GetValue(1))                          // ["bar", "baz", "foo"]
	ass.Equal(t, bazfoo.AsArray(), set.GetValues(2, 3).AsArray()) // ["bar", "baz", "foo"]
	ass.Equal(t, bar.AsArray(), set.GetValues(1, 1).AsArray())    // ["bar", "baz", "foo"]
	var set2 = fra.SetFromSequence(set)                           // ["bar", "baz", "foo"]
	ass.True(t, collator.CompareValues(set, set2))                // ["bar", "baz", "foo"]
	var array = fra.ListFromArray([]string{"foo", "bar", "baz"})  // ["bar", "baz", "foo"]
	var set3 = fra.SetFromSequence(array)                         // ["bar", "baz", "foo"]
	ass.True(t, collator.CompareValues(set2, set3))               // ["bar", "baz", "foo"]
	iterator = set.GetIterator()                                  // ["bar", "baz", "foo"]
	ass.True(t, iterator.HasNext())                               // ["bar", "baz", "foo"]
	ass.False(t, iterator.HasPrevious())                          // ["bar", "baz", "foo"]
	ass.Equal(t, "bar", string(iterator.GetNext()))               // ["bar", "baz", "foo"]
	ass.True(t, iterator.HasPrevious())                           // ["bar", "baz", "foo"]
	iterator.ToEnd()                                              // ["bar", "baz", "foo"]
	ass.False(t, iterator.HasNext())                              // ["bar", "baz", "foo"]
	ass.True(t, iterator.HasPrevious())                           // ["bar", "baz", "foo"]
	ass.Equal(t, "foo", string(iterator.GetPrevious()))           // ["bar", "baz", "foo"]
	ass.True(t, iterator.HasNext())                               // ["bar", "baz", "foo"]
	ass.True(t, set.ContainsValue("baz"))                         // ["bar", "baz", "foo"]
	ass.False(t, set.ContainsValue("bax"))                        // ["bar", "baz", "foo"]
	ass.True(t, set.ContainsAny(baxbaz))                          // ["bar", "baz", "foo"]
	ass.False(t, set.ContainsAny(baxbez))                         // ["bar", "baz", "foo"]
	ass.True(t, set.ContainsAll(barbaz))                          // ["bar", "baz", "foo"]
	ass.False(t, set.ContainsAll(baxbaz))                         // ["bar", "baz", "foo"]
	set.RemoveAll()                                               // [ ]
	ass.True(t, set.IsEmpty())                                    // [ ]
	ass.True(t, set.GetSize() == 0)                               // [ ]
}

func TestSetsWithIntegers(t *tes.T) {
	var array = fra.ListFromArray([]int{3, 1, 4, 5, 9, 2})
	var set = fra.Set[int]()           // [ ]
	set.AddValues(array)               // [1,2,3,4,5,9]
	ass.False(t, set.IsEmpty())        // [1,2,3,4,5,9]
	ass.True(t, set.GetSize() == 6)    // [1,2,3,4,5,9]
	ass.True(t, set.GetValue(1) == 1)  // [1,2,3,4,5,9]
	ass.True(t, set.GetValue(-1) == 9) // [1,2,3,4,5,9]
	set.RemoveValue(6)                 // [1,2,3,4,5,9]
	ass.True(t, set.GetSize() == 6)    // [1,2,3,4,5,9]
	set.RemoveValue(3)                 // [1,2,4,5,9]
	ass.True(t, set.GetSize() == 5)    // [1,2,4,5,9]
	ass.True(t, set.GetValue(3) == 4)  // [1,2,4,5,9]
}

func TestSetsWithTildes(t *tes.T) {
	var array = fra.ListFromArray([]Integer{3, 1, 4, 5, 9, 2})
	var set = fra.Set[Integer]()       // [ ]
	set.AddValues(array)               // [1,2,3,4,5,9]
	ass.False(t, set.IsEmpty())        // [1,2,3,4,5,9]
	ass.True(t, set.GetSize() == 6)    // [1,2,3,4,5,9]
	ass.True(t, set.GetValue(1) == 1)  // [1,2,3,4,5,9]
	ass.True(t, set.GetValue(-1) == 9) // [1,2,3,4,5,9]
	set.RemoveValue(6)                 // [1,2,3,4,5,9]
	ass.True(t, set.GetSize() == 6)    // [1,2,3,4,5,9]
	set.RemoveValue(3)                 // [1,2,4,5,9]
	ass.True(t, set.GetSize() == 5)    // [1,2,4,5,9]
	ass.True(t, set.GetValue(3) == 4)  // [1,2,4,5,9]
}

func TestSetsWithSets(t *tes.T) {
	var array1 = fra.ListFromArray([]int{3, 1, 4, 5, 9, 2})
	var array2 = fra.ListFromArray([]int{7, 1, 4, 5, 9, 2})
	var set1 = fra.Set[int]()
	set1.AddValues(array1)
	var set2 = fra.Set[int]()
	set2.AddValues(array2)
	var set = fra.Set[fra.SetLike[int]]()
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
	var collator = fra.Collator[fra.SetLike[int]]()
	var array1 = fra.ListFromArray([]int{3, 1, 2})
	var array2 = fra.ListFromArray([]int{3, 2, 4})
	var array3 = fra.ListFromArray([]int{3, 2})
	var set1 = fra.Set[int]()
	set1.AddValues(array1)
	var set2 = fra.Set[int]()
	set2.AddValues(array2)
	var set3 = fra.SetClass[int]().And(set1, set2)
	var set4 = fra.Set[int]()
	set4.AddValues(array3)
	ass.True(t, collator.CompareValues(set3, set4))
}

func TestSetsWithSan(t *tes.T) {
	var collator = fra.Collator[fra.SetLike[int]]()
	var array1 = fra.ListFromArray([]int{3, 1, 2})
	var array2 = fra.ListFromArray([]int{3, 2, 4})
	var array3 = fra.ListFromArray([]int{1})
	var set1 = fra.Set[int]()
	set1.AddValues(array1)
	var set2 = fra.Set[int]()
	set2.AddValues(array2)
	var set3 = fra.SetClass[int]().San(set1, set2)
	var set4 = fra.Set[int]()
	set4.AddValues(array3)
	ass.True(t, collator.CompareValues(set3, set4))
}

func TestSetsWithIor(t *tes.T) {
	var collator = fra.Collator[fra.SetLike[int]]()
	var array1 = fra.ListFromArray([]int{3, 1, 5})
	var array2 = fra.ListFromArray([]int{6, 2, 4})
	var array3 = fra.ListFromArray([]int{1, 3, 5, 6, 2, 4})
	var set1 = fra.Set[int]()
	set1.AddValues(array1)
	var set2 = fra.Set[int]()
	set2.AddValues(array2)
	var set3 = fra.SetClass[int]().Ior(set1, set2)
	ass.True(t, set3.ContainsAll(set1))
	ass.True(t, set3.ContainsAll(set2))
	var set4 = fra.Set[int]()
	set4.AddValues(array3)
	ass.True(t, collator.CompareValues(set3, set4))
}

func TestSetsWithXor(t *tes.T) {
	var collator = fra.Collator[fra.SetLike[int]]()
	var array1 = fra.ListFromArray([]int{2, 3, 1, 5})
	var array2 = fra.ListFromArray([]int{6, 2, 5, 4})
	var array3 = fra.ListFromArray([]int{1, 3, 4, 6})
	var set1 = fra.Set[int]()
	set1.AddValues(array1)
	var set2 = fra.Set[int]()
	set2.AddValues(array2)
	var set3 = fra.SetClass[int]().Xor(set1, set2)
	var set4 = fra.Set[int]()
	set4.AddValues(array3)
	ass.True(t, collator.CompareValues(set3, set4))
}

func TestSetsWithEmptySets(t *tes.T) {
	var collator = fra.Collator[fra.SetLike[int]]()
	var set1 = fra.Set[int]()
	var set2 = fra.Set[int]()
	var set3 = fra.SetClass[int]().And(set1, set2)
	var set4 = fra.SetClass[int]().San(set1, set2)
	var set5 = fra.SetClass[int]().Ior(set1, set2)
	var set6 = fra.SetClass[int]().Xor(set1, set2)
	ass.True(t, collator.CompareValues(set3, set4))
	ass.True(t, collator.CompareValues(set4, set5))
	ass.True(t, collator.CompareValues(set5, set6))
	ass.True(t, collator.CompareValues(set6, set1))
}

func TestStackConstructors(t *tes.T) {
	fra.Stack[int64]()
	fra.StackWithCapacity[int64](5)
	var sequence = fra.StackFromArray([]int64{1, 2, 3})
	var stack = fra.StackFromSequence(sequence)
	ass.Equal(t, sequence.AsArray(), stack.AsArray())
}

func TestStackWithSmallCapacity(t *tes.T) {
	var stack = fra.StackWithCapacity[int](1)
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
	var stack = fra.Stack[int]()
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
	var stack = fra.Stack[string]()
	ass.True(t, stack.IsEmpty())
	ass.True(t, stack.GetSize() == 0)
	stack.RemoveAll()
	stack.AddValue("foo")
	stack.AddValue("bar")
	stack.AddValue("baz")
	ass.True(t, stack.GetSize() == 3)
	var last = stack.GetLast()
	ass.Equal(t, last, stack.RemoveLast())
	ass.True(t, stack.GetSize() == 2)
	ass.Equal(t, "bar", stack.RemoveLast())
	ass.True(t, stack.GetSize() == 1)
	stack.RemoveAll()
}
