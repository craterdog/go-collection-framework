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
var k9 string = "map"

// Define the values.
var v1 bool = true
var v2 uint = 0xa
var v3 int = 42
var v4 float64 = 0.125
var v5 complex128 = 1 + 1i
var v6 rune = '☺'
var v7 string = "Hello World!"
var v8 = []int{1, 2, 3}
var v9 = map[string]int{"one": 1, "two": 2, "three": 3}

func TestFormatPrimitives(t *tes.T) {
	var Formatter = col.Formatter()
	ass.Equal(t, "true", Formatter.FormatValue(v1))
	ass.Equal(t, "0xa", Formatter.FormatValue(v2))
	ass.Equal(t, "42", Formatter.FormatValue(v3))
	ass.Equal(t, "0.125", Formatter.FormatValue(v4))
	ass.Equal(t, "(1.0+1.0i)", Formatter.FormatValue(v5))
	ass.Equal(t, "'☺'", Formatter.FormatValue(v6))
	ass.Equal(t, "\"Hello World!\"", Formatter.FormatValue(v7))
}

func TestFormatterConstants(t *tes.T) {
	ass.Equal(t, 8, col.Formatter().DefaultDepth())
}

func TestFormatSpecificDepths(t *tes.T) {
	var array = col.Array[any]().FromArray([]any{1, []any{1, 2, []any{1, 2, 3}}})
	var collator = col.Formatter().WithSpecificDepth(0)
	var s = collator.FormatCollection(array)
	ass.Equal(t, "[...](Array)\n", s)
	collator = col.Formatter().WithSpecificDepth(1)
	s = collator.FormatCollection(array)
	ass.Equal(t, "[\n    1\n    [...](array)\n](Array)\n", s)
	collator = col.Formatter().WithSpecificDepth(2)
	s = collator.FormatCollection(array)
	ass.Equal(t, "[\n    1\n    [\n        1\n        2\n        [...](array)\n    ](array)\n](Array)\n", s)
}

func TestFormatEmptyArray(t *tes.T) {
	var array = []any{}
	var s = col.Formatter().FormatCollection(array)
	ass.Equal(t, "[ ](array)\n", s)
	fmt.Println("Empty Array: " + s)
}

func TestFormatArrayOfAny(t *tes.T) {
	var array = []any{v1, v2, v3, v4, v5, v6, v7, v8, v9}
	var Array = col.Array[any]().FromArray(array)
	var s = col.Formatter().FormatCollection(Array)
	ass.Equal(t, "(Array)\n", s[len(s)-8:])
	fmt.Println("Array of Any: " + s)
}

func TestFormatArrayOfIntegers(t *tes.T) {
	var array = []int{1, 2, 3, 4}
	var s = col.Formatter().FormatCollection(array)
	ass.Equal(t, "(array)\n", s[len(s)-8:])
	fmt.Println("Array of Integers: " + s)
}

func TestFormatEmptyMap(t *tes.T) {
	var mapp = map[any]any{}
	var s = col.Formatter().FormatCollection(mapp)
	ass.Equal(t, "[:](map)\n", s)
	fmt.Println("Empty Map: " + s)
}

func TestFormatMapOfAnyToAny(t *tes.T) {
	var map_ = map[any]any{
		k1: v1,
		k2: v2,
		k3: v3,
		k4: v4,
		k5: v5,
		k6: v6,
		k7: v7,
		k8: v8,
		k9: v9,
	}
	var Map = col.Map[any, any]().FromMap(map_)
	var s = col.Formatter().FormatCollection(Map)
	ass.Equal(t, "(Map)\n", s[len(s)-6:])
	fmt.Println("Map of Any to Any: " + s)
}

func TestFormatMapOfStringToInteger(t *tes.T) {
	var mapp = map[string]int{
		"first":  1,
		"second": 2,
		"third":  3}
	var s = col.Formatter().FormatCollection(mapp)
	ass.Equal(t, "(map)\n", s[len(s)-6:])
	fmt.Println("Map of String to Integer: " + s)
}

func TestFormatEmptyList(t *tes.T) {
	var List = col.List[any]()
	var list = List.Empty()
	var s = col.Formatter().FormatCollection(list)
	ass.Equal(t, "[ ](List)\n", s)
	fmt.Println("Empty List: " + s)
}

func TestFormatListOfAny(t *tes.T) {
	var Array = col.Array[any]()
	var array = Array.FromArray([]any{v1, v2, v3, v4, v5, v6, v7, v8})
	var List = col.List[any]()
	var list = List.FromSequence(array)
	var s = col.Formatter().FormatCollection(list)
	ass.Equal(t, s, fmt.Sprintf("%s", list))
	ass.Equal(t, "(List)\n", s[len(s)-7:])
	fmt.Println("List of Any: " + s)
}

func TestFormatListOfBoolean(t *tes.T) {
	var Array = col.Array[bool]()
	var array = Array.FromArray([]bool{k1, v1})
	var List = col.List[bool]()
	var list = List.FromSequence(array)
	var s = col.Formatter().FormatCollection(list)
	fmt.Println("List of Boolean: " + s)
}

func TestFormatSetOfAny(t *tes.T) {
	var Array = col.Array[any]()
	var array = Array.FromArray([]any{v1, v2, v3, v4, v5, v6, v7, v8})
	var Set = col.Set[any]()
	var set = Set.FromSequence(array)
	var s = col.Formatter().FormatCollection(set)
	ass.Equal(t, s, fmt.Sprintf("%s", set))
	ass.Equal(t, "(Set)\n", s[len(s)-6:])
	fmt.Println("Set of Any: " + s)
}

func TestFormatSetOfSet(t *tes.T) {
	var Array = col.Array[any]()
	var array1 = Array.FromArray([]any{v1, v2, v3, v4, v5, v6, v7, v8})
	var array2 = Array.FromArray([]any{k1, k2, k3, k4, k5, k6, k7, k8})
	var Set = col.Set[any]()
	var set1 = Set.FromSequence(array1)
	var set2 = Set.FromSequence(array2)
	var set = Set.Empty()
	set.AddValue(set1)
	set.AddValue(set2)
	var s = col.Formatter().FormatCollection(set)
	fmt.Println("Set of Set: " + s)
}

func TestFormatStackOfAny(t *tes.T) {
	var Stack = col.Stack[any]()
	var stack = Stack.Empty()
	stack.AddValue(v1)
	stack.AddValue(v2)
	stack.AddValue(v3)
	stack.AddValue(v4)
	stack.AddValue(v5)
	stack.AddValue(v6)
	stack.AddValue(v7)
	var s = col.Formatter().FormatCollection(stack)
	ass.Equal(t, s, fmt.Sprintf("%s", stack))
	ass.Equal(t, "(Stack)\n", s[len(s)-8:])
	fmt.Println("Stack: " + s)
}

func TestFormatQueueOfAny(t *tes.T) {
	var Queue = col.Queue[any]()
	var queue = Queue.Empty()
	queue.AddValue(v1)
	queue.AddValue(v2)
	queue.AddValue(v3)
	queue.AddValue(v4)
	queue.AddValue(v5)
	queue.AddValue(v6)
	queue.AddValue(v7)
	var s = col.Formatter().FormatCollection(queue)
	ass.Equal(t, s, fmt.Sprintf("%s", queue))
	ass.Equal(t, "(Queue)\n", s[len(s)-8:])
	fmt.Println("Queue: " + s)
}

func TestFormatEmptyCatalog(t *tes.T) {
	var Catalog = col.Catalog[any, any]()
	var catalog = Catalog.Empty()
	var s = col.Formatter().FormatCollection(catalog)
	ass.Equal(t, "[:](Catalog)\n", s)
	fmt.Println("Empty Catalog: " + s)
}

func TestFormatCatalogOfAnyToAny(t *tes.T) {
	var Catalog = col.Catalog[any, any]()
	var catalog = Catalog.Empty()
	catalog.SetValue(k1, v1)
	catalog.SetValue(k2, v2)
	catalog.SetValue(k3, v3)
	catalog.SetValue(k4, v4)
	catalog.SetValue(k5, v5)
	catalog.SetValue(k6, v6)
	catalog.SetValue(k7, v7)
	var s = col.Formatter().FormatCollection(catalog)
	ass.Equal(t, s, fmt.Sprintf("%s", catalog))
	ass.Equal(t, "(Catalog)\n", s[len(s)-10:])
	fmt.Println("Catalog: " + s)
}

func TestFormatCatalogOfStringToAny(t *tes.T) {
	var Catalog = col.Catalog[string, any]()
	var catalog = Catalog.Empty()
	catalog.SetValue("key1", v1)
	catalog.SetValue("key2", v2)
	catalog.SetValue("key3", v3)
	catalog.SetValue("key4", v4)
	catalog.SetValue("key5", v5)
	catalog.SetValue("key6", v6)
	catalog.SetValue("key7", v7)
	var s = col.Formatter().FormatCollection(catalog)
	fmt.Println("\nCatalog of String to Any: " + s)
}

func TestFormatCatalogOfStringToInteger(t *tes.T) {
	var Catalog = col.Catalog[string, int]()
	var catalog = Catalog.Empty()
	catalog.SetValue("key1", 1)
	catalog.SetValue("key2", 2)
	catalog.SetValue("key3", 3)
	catalog.SetValue("key4", 4)
	catalog.SetValue("key5", 5)
	catalog.SetValue("key6", 6)
	catalog.SetValue("key7", 7)
	var s = col.Formatter().FormatCollection(catalog)
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
	col.Formatter().FormatCollection(s) // This should panic.
}
