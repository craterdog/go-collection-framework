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
	col "github.com/craterdog/go-collection-framework"
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

func TestListsWithStrings(t *tes.T) {
	var foo = col.ListFromArray([]string{"foo"})
	var bar = col.ListFromArray([]string{"bar"})
	var baz = col.ListFromArray([]string{"baz"})
	var foz = col.ListFromArray([]string{"foz"})
	var barbaz = col.ListFromArray([]string{"bar", "baz"})
	var bazbaz = col.ListFromArray([]string{"baz", "baz"})
	var foobar = col.ListFromArray([]string{"foo", "bar"})
	var baxbaz = col.ListFromArray([]string{"bax", "baz"})
	var baxbez = col.ListFromArray([]string{"bax", "bez"})
	var barfoobax = col.ListFromArray([]string{"bar", "foo", "bax"})
	var foobazbar = col.ListFromArray([]string{"foo", "baz", "bar"})
	var foobarbaz = col.ListFromArray([]string{"foo", "bar", "baz"})
	var barbazfoo = col.ListFromArray([]string{"bar", "baz", "foo"})
	var list = col.List[string]()
	ass.True(t, list.IsEmpty())
	ass.Equal(t, 0, list.GetSize())
	ass.False(t, list.ContainsValue("bax"))
	ass.Equal(t, []string{}, list.AsArray())
	var iterator = col.Iterator[string](list)
	ass.False(t, iterator.HasNext())
	ass.False(t, iterator.HasPrevious())
	iterator.ToStart()
	iterator.ToEnd()
	list.ShuffleValues()
	list.SortValues()
	list.RemoveAll()                              //       [ ]
	list.AddValue("foo")                          //       ["foo"]
	ass.False(t, list.IsEmpty())                  //       ["foo"]
	ass.Equal(t, 1, list.GetSize())               //       ["foo"]
	ass.Equal(t, "foo", string(list.GetValue(1))) //       ["foo"]
	list.AddValues(barbaz)                        //       ["foo", "bar", "baz"]
	ass.Equal(t, 3, list.GetSize())               //       ["foo", "bar", "baz"]
	ass.Equal(t, "foo", string(list.GetValue(1))) //       ["foo", "bar", "baz"]
	ass.Equal(t, barbaz.AsArray(), list.GetValues(2, 3).AsArray())
	ass.Equal(t, foo.AsArray(), list.GetValues(1, 1).AsArray())
	iterator = col.Iterator[string](list)               // ["foo", "bar", "baz"]
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
	list.AddValues(foobarbaz)                                          // ["foo", "bar", "baz"]
	list.SortValues()                                                  // ["bar", "baz", "foo"]
	ass.Equal(t, barbazfoo.AsArray(), list.AsArray())                  // ["bar", "baz", "foo"]
	list.RemoveAll()                                                   // [ ]
	list.AddValue("foo")                                               // ["foo"]
	list.SortValues()                                                  // ["foo"]
	ass.Equal(t, 1, list.GetSize())                                    // ["foo"]
	ass.Equal(t, "foo", string(list.GetValue(1)))                      // ["foo"]
	list.AddValue("bar")                                               // ["foo", "bar"]
	list.SortValues()                                                  // ["bar", "foo"]
	ass.Equal(t, 2, list.GetSize())                                    // ["bar", "foo"]
	ass.Equal(t, "bar", string(list.GetValue(1)))                      // ["bar", "foo"]
}

func TestListsWithConcatenate(t *tes.T) {
	var list1 = col.List[int]()
	var onetwothree = col.ListFromArray([]int{1, 2, 3})
	var fourfivesix = col.ListFromArray([]int{4, 5, 6})
	var onethrusix = col.ListFromArray([]int{1, 2, 3, 4, 5, 6})
	list1.AddValues(onetwothree)
	var list2 = col.List[int]()
	list2.AddValues(fourfivesix)
	var list3 = col.Concatenate(list1, list2)
	var list4 = col.List[int]()
	list4.AddValues(onethrusix)
	ass.True(t, col.CompareValues(list3, list4))
}

func TestListsWithEmptyLists(t *tes.T) {
	var list1 = col.List[int]()
	var list2 = col.List[int]()
	var list3 = col.Concatenate(list1, list2)
	ass.True(t, col.CompareValues(list1, list2))
	ass.True(t, col.CompareValues(list2, list3))
	ass.True(t, col.CompareValues(list3, list1))
}
