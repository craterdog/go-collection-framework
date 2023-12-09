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
	col "github.com/craterdog/go-collection-framework/v3"
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

func TestListsWithStrings(t *tes.T) {
	var Array = col.Array[string]()
	var foo = Array.FromArray([]string{"foo"})
	var bar = Array.FromArray([]string{"bar"})
	var baz = Array.FromArray([]string{"baz"})
	var foz = Array.FromArray([]string{"foz"})
	var barbaz = Array.FromArray([]string{"bar", "baz"})
	var bazbaz = Array.FromArray([]string{"baz", "baz"})
	var foobar = Array.FromArray([]string{"foo", "bar"})
	var baxbaz = Array.FromArray([]string{"bax", "baz"})
	var baxbez = Array.FromArray([]string{"bax", "bez"})
	var barfoobax = Array.FromArray([]string{"bar", "foo", "bax"})
	var foobazbar = Array.FromArray([]string{"foo", "baz", "bar"})
	var foobarbaz = Array.FromArray([]string{"foo", "bar", "baz"})
	var barbazfoo = Array.FromArray([]string{"bar", "baz", "foo"})
	var List = col.List[string]()
	var list = List.FromNothing()
	ass.True(t, list.IsEmpty())
	ass.Equal(t, 0, list.GetSize())
	ass.False(t, list.ContainsValue("bax"))
	ass.Equal(t, []string{}, list.AsArray())
	var iterator = list.GetIterator()
	ass.False(t, iterator.HasNext())
	ass.False(t, iterator.HasPrevious())
	iterator.ToStart()
	iterator.ToEnd()
	list.ShuffleValues()
	list.SortValues()
	list.RemoveAll()                              //       [ ]
	list.AppendValue("foo")                       //       ["foo"]
	ass.False(t, list.IsEmpty())                  //       ["foo"]
	ass.Equal(t, 1, list.GetSize())               //       ["foo"]
	ass.Equal(t, "foo", string(list.GetValue(1))) //       ["foo"]
	list.AppendValues(barbaz)                     //       ["foo", "bar", "baz"]
	ass.Equal(t, 3, list.GetSize())               //       ["foo", "bar", "baz"]
	ass.Equal(t, "foo", string(list.GetValue(1))) //       ["foo", "bar", "baz"]
	ass.Equal(t, barbaz.AsArray(), list.GetValues(2, 3).AsArray())
	ass.Equal(t, foo.AsArray(), list.GetValues(1, 1).AsArray())
	var list2 = List.FromSequence(list)
	ass.True(t, col.CompareValues(list, list2))
	var array = Array.FromArray([]string{"foo", "bar", "baz"})
	var list3 = List.FromSequence(array)
	list2.SortValues()
	list3.SortValues()
	ass.True(t, col.CompareValues(list2, list3))
	iterator = list.GetIterator()                       // ["foo", "bar", "baz"]
	ass.True(t, iterator.HasNext())                     // ["foo", "bar", "baz"]
	ass.False(t, iterator.HasPrevious())                // ["foo", "bar", "baz"]
	ass.Equal(t, "foo", string(iterator.GetNext()))     // ["foo", "bar", "baz"]
	ass.True(t, iterator.HasPrevious())                 // ["foo", "bar", "baz"]
	iterator.ToEnd()                                    // ["foo", "bar", "baz"]
	ass.False(t, iterator.HasNext())                    // ["foo", "bar", "baz"]
	ass.True(t, iterator.HasPrevious())                 // ["foo", "bar", "baz"]
	ass.Equal(t, "baz", string(iterator.GetPrevious())) // ["foo", "bar", "baz"]
	ass.True(t, iterator.HasNext())                     // ["foo", "bar", "baz"]
	list.ShuffleValues()                                // [ ?, ?, ? ]
	list.RemoveAll()                                    // [ ]
	ass.True(t, list.IsEmpty())                         // [ ]
	ass.Equal(t, 0, list.GetSize())                     // [ ]
	list.InsertValue(0, "baz")                          // ["baz"]
	ass.Equal(t, 1, list.GetSize())                     // ["baz"]
	ass.Equal(t, "baz", string(list.GetValue(-1)))      // ["baz"]
	list.InsertValues(0, foobar)                        // ["foo", "bar", "baz"]
	ass.Equal(t, 3, list.GetSize())                     // ["foo", "bar", "baz"]
	ass.Equal(t, "foo", string(list.GetValue(-3)))      // ["foo", "bar", "baz"]
	ass.Equal(t, "bar", string(list.GetValue(-2)))      // ["foo", "bar", "baz"]
	ass.Equal(t, "baz", string(list.GetValue(-1)))      // ["foo", "bar", "baz"]
	list.ReverseValues()                                // ["baz", "bar", "foo"]
	ass.Equal(t, "foo", string(list.GetValue(-1)))      // ["baz", "bar", "foo"]
	ass.Equal(t, "bar", string(list.GetValue(-2)))      // ["baz", "bar", "foo"]
	ass.Equal(t, "baz", string(list.GetValue(-3)))      // ["baz", "bar", "foo"]
	list.ReverseValues()                                // ["foo", "bar", "baz"]
	ass.Equal(t, 0, list.GetIndex("foz"))               // ["foo", "bar", "baz"]
	ass.Equal(t, 3, list.GetIndex("baz"))               // ["foo", "bar", "baz"]
	ass.True(t, list.ContainsValue("baz"))              // ["foo", "bar", "baz"]
	ass.False(t, list.ContainsValue("bax"))             // ["foo", "bar", "baz"]
	ass.True(t, list.ContainsAny(baxbaz))               // ["foo", "bar", "baz"]
	ass.False(t, list.ContainsAny(baxbez))              // ["foo", "bar", "baz"]
	ass.True(t, list.ContainsAll(barbaz))               // ["foo", "bar", "baz"]
	ass.False(t, list.ContainsAll(baxbaz))              // ["foo", "bar", "baz"]
	list.SetValue(3, "bax")                             // ["foo", "bar", "bax"]
	list.InsertValues(3, baz)                           // ["foo", "bar", "bax", "baz"]
	ass.Equal(t, 4, list.GetSize())                     // ["foo", "bar", "bax", "baz"]
	ass.Equal(t, "baz", string(list.GetValue(4)))       // ["foo", "bar", "bax", "baz"]
	list.InsertValue(4, "bar")                          // ["foo", "bar", "bax", "baz", "bar"]
	ass.Equal(t, 5, list.GetSize())                     // ["foo", "bar", "bax", "baz", "bar"]
	ass.Equal(t, "bar", string(list.GetValue(5)))       // ["foo", "bar", "bax", "baz", "bar"]
	list.InsertValue(2, "foo")                          // ["foo", "bar", "foo", "bax", "baz", "bar"]
	ass.Equal(t, 6, list.GetSize())                     // ["foo", "bar", "foo", "bax", "baz", "bar"]
	ass.Equal(t, "bar", string(list.GetValue(2)))       // ["foo", "bar", "foo", "bax", "baz", "bar"]
	ass.Equal(t, "foo", string(list.GetValue(3)))       // ["foo", "bar", "foo", "bax", "baz", "bar"]
	ass.Equal(t, "bax", string(list.GetValue(4)))       // ["foo", "bar", "foo", "bax", "baz", "bar"]
	ass.Equal(t, bar.AsArray(), list.GetValues(6, 6).AsArray())
	list.InsertValues(5, baz)                     //       ["foo", "bar", "foo", "bax", "baz", "baz", "bar"]
	ass.Equal(t, 7, list.GetSize())               //       ["foo", "bar", "foo", "bax", "baz", "baz", "bar"]
	ass.Equal(t, "bax", string(list.GetValue(4))) //       ["foo", "bar", "foo", "bax", "baz", "baz", "bar"]
	ass.Equal(t, "baz", string(list.GetValue(5))) //       ["foo", "bar", "foo", "bax", "baz", "baz", "bar"]
	ass.Equal(t, "baz", string(list.GetValue(6))) //       ["foo", "bar", "foo", "bax", "baz", "baz", "bar"]
	ass.Equal(t, barfoobax.AsArray(), list.GetValues(2, -4).AsArray())
	list.SetValues(2, foobazbar) //                        ["foo", "foo", "baz", "bar", "baz", "baz", "bar"]
	ass.Equal(t, foobazbar.AsArray(), list.GetValues(2, -4).AsArray())
	list.SetValues(-1, foz)
	ass.Equal(t, "foz", string(list.GetValue(-1))) //      ["foo", "foo", "baz", "bar", "baz", "baz", "foz"]
	list.SortValues()                              //      ["bar", "baz", "baz", "baz", "foo", "foo", "foz"]

	ass.Equal(t, bazbaz.AsArray(), list.RemoveValues(2, -5).AsArray()) // ["bar", "baz", "foo", "foo", "foz"]
	ass.Equal(t, barbaz.AsArray(), list.RemoveValues(1, 2).AsArray())  // ["foo", "foo", "foz"]
	ass.Equal(t, "foz", string(list.RemoveValue(-1)))                  // ["foo", "foo"]
	ass.Equal(t, 2, list.GetSize())                                    // ["foo", "foo"]
	list.RemoveAll()                                                   // [ ]
	ass.Equal(t, 0, list.GetSize())                                    // [ ]
	list.SortValues()                                                  // [ ]
	list.AppendValues(foobarbaz)                                       // ["foo", "bar", "baz"]
	list.SortValues()                                                  // ["bar", "baz", "foo"]
	ass.Equal(t, barbazfoo.AsArray(), list.AsArray())                  // ["bar", "baz", "foo"]
	list.RemoveAll()                                                   // [ ]
	list.AppendValue("foo")                                            // ["foo"]
	list.SortValues()                                                  // ["foo"]
	ass.Equal(t, 1, list.GetSize())                                    // ["foo"]
	ass.Equal(t, "foo", string(list.GetValue(1)))                      // ["foo"]
	list.AppendValue("bar")                                            // ["foo", "bar"]
	list.SortValues()                                                  // ["bar", "foo"]
	ass.Equal(t, 2, list.GetSize())                                    // ["bar", "foo"]
	ass.Equal(t, "bar", string(list.GetValue(1)))                      // ["bar", "foo"]
}

func TestListsWithTildas(t *tes.T) {
	var Array = col.Array[Integer]()
	var array = Array.FromArray([]Integer{3, 1, 4, 5, 9, 2})
	var List = col.List[Integer]()
	var list = List.FromSequence(array)
	ass.False(t, list.IsEmpty())            // [3,1,4,5,9,2]
	ass.Equal(t, 6, list.GetSize())         // [3,1,4,5,9,2]
	ass.Equal(t, 3, int(list.GetValue(1)))  // [3,1,4,5,9,2]
	ass.Equal(t, 2, int(list.GetValue(-1))) // [3,1,4,5,9,2]
	list.SortValues()                       // [1,2,3,4,5,9]
	ass.Equal(t, 6, list.GetSize())         // [1,2,3,4,5,9]
	ass.Equal(t, 3, int(list.GetValue(3)))  // [1,2,3,4,5,9]
}

func BadCompare(first col.Value, second col.Value) bool {
	panic("KaPow!")
}

func TestListsWithComparer(t *tes.T) {
	var List = col.List[int]()
	var list = List.FromComparer(BadCompare)
	list.AppendValue(1)
	list.AppendValue(2)
	list.AppendValue(3)
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "KaPow!", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	list.GetIndex(2)
}

func TestListsWithConcatenate(t *tes.T) {
	var Array = col.Array[int]()
	var onetwothree = Array.FromArray([]int{1, 2, 3})
	var fourfivesix = Array.FromArray([]int{4, 5, 6})
	var onethrusix = Array.FromArray([]int{1, 2, 3, 4, 5, 6})
	var List = col.List[int]()
	var list1 = List.FromNothing()
	list1.AppendValues(onetwothree)
	var list2 = List.FromNothing()
	list2.AppendValues(fourfivesix)
	var list3 = List.Concatenate(list1, list2)
	var list4 = List.FromNothing()
	list4.AppendValues(onethrusix)
	ass.True(t, col.CompareValues(list3, list4))
}

func TestListsWithEmptyLists(t *tes.T) {
	var List = col.List[int]()
	var empty = List.FromNothing()
	var list = List.Concatenate(empty, empty)
	ass.True(t, col.CompareValues(empty, empty))
	ass.True(t, col.CompareValues(list, empty))
	ass.True(t, col.CompareValues(empty, list))
	ass.True(t, col.CompareValues(list, list))
}
