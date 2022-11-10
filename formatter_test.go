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
	col "github.com/craterdog/go-collection-framework"
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

func TestFormatterWithIndentation(t *tes.T) {
	var formatter = col.Formatter(1)
	ass.Equal(t, 1, formatter.GetIndentation())
}

func TestFormatterWithPrimitives(t *tes.T) {
	s0 := col.FormatValue(nil)
	ass.Equal(t, "<nil>", s0)
	fmt.Println("\nNil: " + s0)

	s1 := col.FormatValue(v1)
	ass.Equal(t, "true", s1)
	fmt.Println("\nBoolean: " + s1)

	s2 := col.FormatValue(v2)
	ass.Equal(t, "0xa", s2)
	fmt.Println("\nUnsigned: " + s2)

	s3 := col.FormatValue(v3)
	ass.Equal(t, "42", s3)
	fmt.Println("\nInteger: " + s3)

	s4 := col.FormatValue(v4)
	ass.Equal(t, "0.125", s4)
	fmt.Println("\nFloat: " + s4)

	s5 := col.FormatValue(v5)
	ass.Equal(t, "(1+1i)", s5)
	fmt.Println("\nComplex: " + s5)

	s6 := col.FormatValue(v6)
	ass.Equal(t, "'☺'", s6)
	fmt.Println("\nRune: " + s6)

	s7 := col.FormatValue(v7)
	ass.Equal(t, "\"Hello World!\"", s7)
	fmt.Println("\nString: " + s7)
}

func TestFormatterWithEmptyArray(t *tes.T) {
	array := []any{}
	s := col.FormatValue(array)
	ass.Equal(t, "[ ](array)", s)
	fmt.Println("\nEmpty Array: " + s)
}

func TestFormatterWithArrayOfAny(t *tes.T) {
	array := []any{v1, v2, v3, v4, v5, v6, v7, v8}
	s := col.FormatValue(array)
	ass.Equal(t, "(array)", s[len(s)-7:])
	fmt.Println("\nArray of Any: " + s)
}

func TestFormatterWithArrayOfIntegers(t *tes.T) {
	array := []int{1, 2, 3, 4}
	s := col.FormatValue(array)
	ass.Equal(t, "(array)", s[len(s)-7:])
	fmt.Println("\nArray of Integers: " + s)
}

func TestFormatterWithEmptyMap(t *tes.T) {
	mapp := map[any]any{}
	s := col.FormatValue(mapp)
	ass.Equal(t, "[:](map)", s)
	fmt.Println("\nEmpty Map: " + s)
}

func TestFormatterWithMapOfAnyToAny(t *tes.T) {
	mapp := map[any]any{
		k1: v1,
		k2: v2,
		k3: v3,
		k4: v4,
		k5: v5,
		k6: v6,
		k7: v7,
		k8: v8}
	s := col.FormatValue(mapp)
	ass.Equal(t, "(map)", s[len(s)-5:])
	fmt.Println("\nMap of Any to Any: " + s)
}

func TestFormatterWithMapOfStringToInteger(t *tes.T) {
	mapp := map[string]int{
		"first":  1,
		"second": 2,
		"third":  3}
	s := col.FormatValue(mapp)
	ass.Equal(t, "(map)", s[len(s)-5:])
	fmt.Println("\nMap of String to Integer: " + s)
}

func TestFormatterWithEmptyList(t *tes.T) {
	list := col.List[any]()
	s := col.FormatValue(list)
	ass.Equal(t, "[ ](list)", s)
	fmt.Println("\nEmpty List: " + s)
}

func TestFormatterWithListOfAny(t *tes.T) {
	list := col.ListFromArray[any]([]any{v1, v2, v3, v4, v5, v6, v7, v8})
	s := col.FormatValue(list)
	ass.Equal(t, "(list)", s[len(s)-6:])
	fmt.Println("\nList of Any: " + s)
}

func TestFormatterWithListOfBoolean(t *tes.T) {
	list := col.ListFromArray[bool]([]bool{k1, v1})
	s := col.FormatValue(list)
	fmt.Println("\nList of Boolean: " + s)
}

func TestFormatterWithSetOfAny(t *tes.T) {
	set := col.SetFromArray[any]([]any{v1, v2, v3, v4, v5, v6, v7, v8})
	s := col.FormatValue(set)
	ass.Equal(t, "(set)", s[len(s)-5:])
	fmt.Println("\nSet of Any: " + s)
}

func TestFormatterWithSetOfSet(t *tes.T) {
	set1 := col.SetFromArray[any]([]any{v1, v2, v3, v4, v5, v6, v7, v8})
	set2 := col.SetFromArray[any]([]any{k1, k2, k3, k4, k5, k6, k7, k8})
	set := col.SetFromArray[col.SetLike[any]]([]col.SetLike[any]{set1, set2})
	s := col.FormatValue(set)
	fmt.Println("\nSet of Set: " + s)
}

func TestFormatterWithStackOfAny(t *tes.T) {
	stack := col.Stack[any]()
	stack.AddValue(v1)
	stack.AddValue(v2)
	stack.AddValue(v3)
	stack.AddValue(v4)
	stack.AddValue(v5)
	stack.AddValue(v6)
	stack.AddValue(v7)
	s := col.FormatValue(stack)
	ass.Equal(t, "(stack)", s[len(s)-7:])
	fmt.Println("\nStack: " + s)
}

func TestFormatterWithQueueOfAny(t *tes.T) {
	var queue col.FIFO[any] = col.Queue[any]()
	queue.AddValue(v1)
	queue.AddValue(v2)
	queue.AddValue(v3)
	queue.AddValue(v4)
	queue.AddValue(v5)
	queue.AddValue(v6)
	queue.AddValue(v7)
	s := col.FormatValue(queue)
	ass.Equal(t, "(queue)", s[len(s)-7:])
	fmt.Println("\nQueue: " + s)
}

func TestFormatterWithAssociationOfAnyToAny(t *tes.T) {
	association := col.Association[any, any]("foo", 5)
	s := col.FormatValue(association)
	fmt.Println("\nAssociation of Any to Any: " + s)
}

func TestFormatterWithAssociationOfStringToInteger(t *tes.T) {
	association := col.Association[string, int]("bar", 42)
	s := col.FormatValue(association)
	fmt.Println("\nAssociation of String to Integer: " + s)
}

func TestFormatterWithEmptyCatalog(t *tes.T) {
	catalog := col.Catalog[any, any]()
	s := col.FormatValue(catalog)
	ass.Equal(t, "[:](catalog)", s)
	fmt.Println("\nEmpty Catalog: " + s)
}

func TestFormatterWithCatalogOfAnyToAny(t *tes.T) {
	catalog := col.Catalog[any, any]()
	catalog.SetValue(k1, v1)
	catalog.SetValue(k2, v2)
	catalog.SetValue(k3, v3)
	catalog.SetValue(k4, v4)
	catalog.SetValue(k5, v5)
	catalog.SetValue(k6, v6)
	catalog.SetValue(k7, v7)
	s := col.FormatValue(catalog)
	ass.Equal(t, "(catalog)", s[len(s)-9:])
	fmt.Println("\nCatalog: " + s)
}

func TestFormatterWithCatalogOfStringToAny(t *tes.T) {
	catalog := col.Catalog[string, any]()
	catalog.SetValue("key1", v1)
	catalog.SetValue("key2", v2)
	catalog.SetValue("key3", v3)
	catalog.SetValue("key4", v4)
	catalog.SetValue("key5", v5)
	catalog.SetValue("key6", v6)
	catalog.SetValue("key7", v7)
	s := col.FormatValue(catalog)
	fmt.Println("\nCatalog of String to Any: " + s)
}

func TestFormatterWithCatalogOfStringToInteger(t *tes.T) {
	catalog := col.Catalog[string, int]()
	catalog.SetValue("key1", 1)
	catalog.SetValue("key2", 2)
	catalog.SetValue("key3", 3)
	catalog.SetValue("key4", 4)
	catalog.SetValue("key5", 5)
	catalog.SetValue("key6", 6)
	catalog.SetValue("key7", 7)
	s := col.FormatValue(catalog)
	fmt.Println("\nCatalog of String to Integer: " + s)
}
