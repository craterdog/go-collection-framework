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
	fmt "fmt"
	col "github.com/craterdog/go-collection-framework/v3"
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

// Define the zero value keys.
var k1 bool
var k2 uint
var k3 int
var k4 float64
var k5 complex128
var k6 rune
var k7 string
var k8 string = "array"

// Define the values.
var v1 bool = true
var v2 uint = 0xa
var v3 int = 42
var v4 float64 = 0.125
var v5 complex128 = 1 + 1i
var v6 rune = '☺'
var v7 string = "Hello World!"
var v8 = []int{1, 2, 3}

func TestFormatEmptyArray(t *tes.T) {
	var array = []any{}
	var s = col.FormatCollection(array)
	ass.Equal(t, "[ ](array)", s)
	fmt.Println("\nEmpty Array: " + s)
}

func TestFormatArrayOfAny(t *tes.T) {
	var array = []any{v1, v2, v3, v4, v5, v6, v7, v8}
	var s = col.FormatCollection(array)
	ass.Equal(t, "(array)", s[len(s)-7:])
	fmt.Println("\nArray of Any: " + s)
}

func TestFormatArrayOfIntegers(t *tes.T) {
	var array = []int{1, 2, 3, 4}
	var s = col.FormatCollection(array)
	ass.Equal(t, "(array)", s[len(s)-7:])
	fmt.Println("\nArray of Integers: " + s)
}

func TestFormatEmptyMap(t *tes.T) {
	var mapp = map[any]any{}
	var s = col.FormatCollection(mapp)
	ass.Equal(t, "[:](map)", s)
	fmt.Println("\nEmpty Map: " + s)
}

func TestFormatMapOfAnyToAny(t *tes.T) {
	var mapp = map[any]any{
		k1: v1,
		k2: v2,
		k3: v3,
		k4: v4,
		k5: v5,
		k6: v6,
		k7: v7,
		k8: v8}
	var s = col.FormatCollection(mapp)
	ass.Equal(t, "(map)", s[len(s)-5:])
	fmt.Println("\nMap of Any to Any: " + s)
}

func TestFormatMapOfStringToInteger(t *tes.T) {
	var mapp = map[string]int{
		"first":  1,
		"second": 2,
		"third":  3}
	var s = col.FormatCollection(mapp)
	ass.Equal(t, "(map)", s[len(s)-5:])
	fmt.Println("\nMap of String to Integer: " + s)
}

func TestFormatEmptyList(t *tes.T) {
	var List = col.List[any]()
	var list = List.FromNothing()
	var s = col.FormatCollection(list)
	ass.Equal(t, "[ ](list)", s)
	fmt.Println("\nEmpty List: " + s)
}

func TestFormatListOfAny(t *tes.T) {
	var Array = col.Array[any]()
	var array = Array.FromArray([]any{v1, v2, v3, v4, v5, v6, v7, v8})
	var List = col.List[any]()
	var list = List.FromSequence(array)
	var s = col.FormatCollection(list)
	ass.Equal(t, s, fmt.Sprintf("%s", list))
	ass.Equal(t, "(list)", s[len(s)-6:])
	fmt.Println("\nList of Any: " + s)
}

func TestFormatListOfBoolean(t *tes.T) {
	var Array = col.Array[bool]()
	var array = Array.FromArray([]bool{k1, v1})
	var List = col.List[bool]()
	var list = List.FromSequence(array)
	var s = col.FormatCollection(list)
	fmt.Println("\nList of Boolean: " + s)
}

func TestFormatSetOfAny(t *tes.T) {
	var Array = col.Array[any]()
	var array = Array.FromArray([]any{v1, v2, v3, v4, v5, v6, v7, v8})
	var Set = col.Set[any]()
	var set = Set.FromSequence(array)
	var s = col.FormatCollection(set)
	ass.Equal(t, s, fmt.Sprintf("%s", set))
	ass.Equal(t, "(set)", s[len(s)-5:])
	fmt.Println("\nSet of Any: " + s)
}

func TestFormatSetOfSet(t *tes.T) {
	var Array = col.Array[any]()
	var array1 = Array.FromArray([]any{v1, v2, v3, v4, v5, v6, v7, v8})
	var array2 = Array.FromArray([]any{k1, k2, k3, k4, k5, k6, k7, k8})
	var Set = col.Set[any]()
	var set1 = Set.FromSequence(array1)
	var set2 = Set.FromSequence(array2)
	var set = Set.FromNothing()
	set.AddValue(set1)
	set.AddValue(set2)
	var s = col.FormatCollection(set)
	fmt.Println("\nSet of Set: " + s)
}

func TestFormatStackOfAny(t *tes.T) {
	var stack = col.Stack[any]()
	stack.AddValue(v1)
	stack.AddValue(v2)
	stack.AddValue(v3)
	stack.AddValue(v4)
	stack.AddValue(v5)
	stack.AddValue(v6)
	stack.AddValue(v7)
	var s = col.FormatCollection(stack)
	ass.Equal(t, s, fmt.Sprintf("%s", stack))
	ass.Equal(t, "(stack)", s[len(s)-7:])
	fmt.Println("\nStack: " + s)
}

func TestFormatQueueOfAny(t *tes.T) {
	var queue col.FIFO[any] = col.Queue[any]()
	queue.AddValue(v1)
	queue.AddValue(v2)
	queue.AddValue(v3)
	queue.AddValue(v4)
	queue.AddValue(v5)
	queue.AddValue(v6)
	queue.AddValue(v7)
	var s = col.FormatCollection(queue)
	ass.Equal(t, s, fmt.Sprintf("%s", queue))
	ass.Equal(t, "(queue)", s[len(s)-7:])
	fmt.Println("\nQueue: " + s)
}

func TestFormatEmptyCatalog(t *tes.T) {
	var Catalog = col.Catalog[any, any]()
	var catalog = Catalog.FromNothing()
	var s = col.FormatCollection(catalog)
	ass.Equal(t, "[:](catalog)", s)
	fmt.Println("\nEmpty Catalog: " + s)
}

func TestFormatCatalogOfAnyToAny(t *tes.T) {
	var Catalog = col.Catalog[any, any]()
	var catalog = Catalog.FromNothing()
	catalog.SetValue(k1, v1)
	catalog.SetValue(k2, v2)
	catalog.SetValue(k3, v3)
	catalog.SetValue(k4, v4)
	catalog.SetValue(k5, v5)
	catalog.SetValue(k6, v6)
	catalog.SetValue(k7, v7)
	var s = col.FormatCollection(catalog)
	ass.Equal(t, s, fmt.Sprintf("%s", catalog))
	ass.Equal(t, "(catalog)", s[len(s)-9:])
	fmt.Println("\nCatalog: " + s)
}

func TestFormatCatalogOfStringToAny(t *tes.T) {
	var Catalog = col.Catalog[string, any]()
	var catalog = Catalog.FromNothing()
	catalog.SetValue("key1", v1)
	catalog.SetValue("key2", v2)
	catalog.SetValue("key3", v3)
	catalog.SetValue("key4", v4)
	catalog.SetValue("key5", v5)
	catalog.SetValue("key6", v6)
	catalog.SetValue("key7", v7)
	var s = col.FormatCollection(catalog)
	fmt.Println("\nCatalog of String to Any: " + s)
}

func TestFormatCatalogOfStringToInteger(t *tes.T) {
	var Catalog = col.Catalog[string, int]()
	var catalog = Catalog.FromNothing()
	catalog.SetValue("key1", 1)
	catalog.SetValue("key2", 2)
	catalog.SetValue("key3", 3)
	catalog.SetValue("key4", 4)
	catalog.SetValue("key5", 5)
	catalog.SetValue("key6", 6)
	catalog.SetValue("key7", 7)
	var s = col.FormatCollection(catalog)
	fmt.Println("\nCatalog of String to Integer: " + s)
}

func TestFormatInvalidType(t *tes.T) {
	var s struct{}
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "Attempted to format:\n    value: {}\n    type: struct {}\n    kind: struct\n", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	col.FormatCollection(s) // This should panic.
}
