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
	col "github.com/craterdog/go-collection-framework/v3"
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

// Tilde Types
type Boolean bool
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

func TestCollatorConstants(t *tes.T) {
	var Collator = col.Collator()
	ass.Equal(t, 16, Collator.DefaultMaximum())
}

func TestCompareMaximum(t *tes.T) {
	var collator = col.Collator().MakeWithMaximum(1)
	var array = col.Array[any]().MakeFromArray([]any{"foo", []int{1, 2, 3}})
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The maximum traversal depth was exceeded: 1", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	_ = collator.CompareValues(array, array)
}

func TestRankMaximum(t *tes.T) {
	var collator = col.Collator().MakeWithMaximum(1)
	var array = col.Array[any]().MakeFromArray([]any{"foo", []int{1, 2, 3}})
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The maximum traversal depth was exceeded: 1", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	_ = collator.RankValues(array, array)
}

func TestComparison(t *tes.T) {
	var collator = col.Collator().Make()

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
	var collator = col.Collator().Make()

	// Boolean
	var False = Boolean(false)
	var True = Boolean(true)
	var ShouldBeFalse Boolean

	ass.Equal(t, 0, collator.RankValues(ShouldBeFalse, ShouldBeFalse))
	ass.Equal(t, -1, collator.RankValues(ShouldBeFalse, True))
	ass.Equal(t, 0, collator.RankValues(False, ShouldBeFalse))
	ass.Equal(t, 1, collator.RankValues(True, ShouldBeFalse))
	ass.Equal(t, 0, collator.RankValues(ShouldBeFalse, False))
	ass.Equal(t, -1, collator.RankValues(False, True))
	ass.Equal(t, 0, collator.RankValues(False, False))
	ass.Equal(t, 1, collator.RankValues(True, False))
	ass.Equal(t, 0, collator.RankValues(True, True))

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
	var collator = col.Collator().Make()
	var array = col.Array[any]().MakeFromArray([]any{0})
	array.SetValue(1, array) // Now it is recursive.
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The maximum traversal depth was exceeded: 16", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	collator.CompareValues(array, array) // This should panic.
}

func TestCompareRecursiveMaps(t *tes.T) {
	var collator = col.Collator().Make()
	var m = col.Map[string, any]().MakeFromMap(map[string]any{"first": 1})
	m.SetValue("first", m) // Now it is recursive.
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The maximum traversal depth was exceeded: 16", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	collator.CompareValues(m, m) // This should panic.
}

func TestRanking(t *tes.T) {
	var collator = col.Collator().Make()

	// Nil
	var ShouldBeNil any

	ass.Equal(t, 0, collator.RankValues(nil, nil))
	ass.Equal(t, 0, collator.RankValues(nil, ShouldBeNil))
	ass.Equal(t, 0, collator.RankValues(ShouldBeNil, ShouldBeNil))
	ass.Equal(t, 0, collator.RankValues(ShouldBeNil, nil))

	// Boolean
	var False = false
	var True = true
	var ShouldBeFalse bool

	ass.Equal(t, 0, collator.RankValues(ShouldBeFalse, ShouldBeFalse))
	ass.Equal(t, -1, collator.RankValues(ShouldBeFalse, True))
	ass.Equal(t, 0, collator.RankValues(False, ShouldBeFalse))
	ass.Equal(t, 1, collator.RankValues(True, ShouldBeFalse))
	ass.Equal(t, 0, collator.RankValues(ShouldBeFalse, False))
	ass.Equal(t, -1, collator.RankValues(False, True))
	ass.Equal(t, 0, collator.RankValues(False, False))
	ass.Equal(t, 1, collator.RankValues(True, False))
	ass.Equal(t, 0, collator.RankValues(True, True))

	// Byte
	var Zero byte = 0x00
	var One byte = 0x01
	var ShouldBeZero byte

	ass.Equal(t, 0, collator.RankValues(ShouldBeZero, ShouldBeZero))
	ass.Equal(t, -1, collator.RankValues(ShouldBeZero, One))
	ass.Equal(t, 0, collator.RankValues(Zero, ShouldBeZero))
	ass.Equal(t, 1, collator.RankValues(One, ShouldBeZero))
	ass.Equal(t, 0, collator.RankValues(ShouldBeZero, Zero))
	ass.Equal(t, -1, collator.RankValues(Zero, One))
	ass.Equal(t, 0, collator.RankValues(Zero, Zero))
	ass.Equal(t, 1, collator.RankValues(One, Zero))
	ass.Equal(t, 0, collator.RankValues(One, One))

	// Integer
	var Zilch = 0
	var Two = 2
	var Three = 3
	var ShouldBeZilch int

	ass.Equal(t, 0, collator.RankValues(ShouldBeZilch, ShouldBeZilch))
	ass.Equal(t, -1, collator.RankValues(ShouldBeZilch, Two))
	ass.Equal(t, 0, collator.RankValues(Zilch, ShouldBeZilch))
	ass.Equal(t, 1, collator.RankValues(Two, ShouldBeZilch))
	ass.Equal(t, 0, collator.RankValues(ShouldBeZilch, Zilch))
	ass.Equal(t, -1, collator.RankValues(Two, Three))
	ass.Equal(t, 0, collator.RankValues(Two, Two))
	ass.Equal(t, 1, collator.RankValues(Three, Two))
	ass.Equal(t, 0, collator.RankValues(Three, Three))

	// Float
	var Negligible = 0.0
	var Fourth = 0.25
	var Half = 0.5
	var ShouldBeNegligible float64

	ass.Equal(t, 0, collator.RankValues(ShouldBeNegligible, ShouldBeNegligible))
	ass.Equal(t, -1, collator.RankValues(ShouldBeNegligible, Half))
	ass.Equal(t, 0, collator.RankValues(Negligible, ShouldBeNegligible))
	ass.Equal(t, 1, collator.RankValues(Half, ShouldBeNegligible))
	ass.Equal(t, 0, collator.RankValues(ShouldBeNegligible, Negligible))
	ass.Equal(t, -1, collator.RankValues(Fourth, Half))
	ass.Equal(t, 0, collator.RankValues(Fourth, Fourth))
	ass.Equal(t, 1, collator.RankValues(Half, Fourth))
	ass.Equal(t, 0, collator.RankValues(Half, Half))

	// Complex
	var Origin = 0 + 0i
	var PiOver4 = 1 + 1i
	var PiOver2 = 1 + 0i
	var ShouldBeOrigin complex128

	ass.Equal(t, 0, collator.RankValues(ShouldBeOrigin, ShouldBeOrigin))
	ass.Equal(t, -1, collator.RankValues(ShouldBeOrigin, PiOver4))
	ass.Equal(t, 0, collator.RankValues(Origin, ShouldBeOrigin))
	ass.Equal(t, 1, collator.RankValues(PiOver4, ShouldBeOrigin))
	ass.Equal(t, 0, collator.RankValues(ShouldBeOrigin, Origin))
	ass.Equal(t, -1, collator.RankValues(PiOver2, PiOver4))
	ass.Equal(t, 0, collator.RankValues(PiOver2, PiOver2))
	ass.Equal(t, 1, collator.RankValues(PiOver4, PiOver2))
	ass.Equal(t, 0, collator.RankValues(PiOver4, PiOver4))

	// Rune
	var Null = rune(0)
	var Sad = '☹'
	var Happy = '☺'
	var ShouldBeNull rune

	ass.Equal(t, 0, collator.RankValues(ShouldBeNull, ShouldBeNull))
	ass.Equal(t, -1, collator.RankValues(ShouldBeNull, Sad))
	ass.Equal(t, 0, collator.RankValues(Null, ShouldBeNull))
	ass.Equal(t, 1, collator.RankValues(Sad, ShouldBeNull))
	ass.Equal(t, 0, collator.RankValues(ShouldBeNull, Null))
	ass.Equal(t, -1, collator.RankValues(Sad, Happy))
	ass.Equal(t, 0, collator.RankValues(Sad, Sad))
	ass.Equal(t, 1, collator.RankValues(Happy, Sad))
	ass.Equal(t, 0, collator.RankValues(Happy, Happy))

	// String
	var Empty = ""
	var Hello = "Hello"
	var World = "World"
	var ShouldBeEmpty string

	ass.Equal(t, 0, collator.RankValues(ShouldBeEmpty, ShouldBeEmpty))
	ass.Equal(t, -1, collator.RankValues(ShouldBeEmpty, Hello))
	ass.Equal(t, 0, collator.RankValues(Empty, ShouldBeEmpty))
	ass.Equal(t, 1, collator.RankValues(Hello, ShouldBeEmpty))
	ass.Equal(t, 0, collator.RankValues(ShouldBeEmpty, Empty))
	ass.Equal(t, -1, collator.RankValues(Hello, World))
	ass.Equal(t, 0, collator.RankValues(Hello, Hello))
	ass.Equal(t, 1, collator.RankValues(World, Hello))
	ass.Equal(t, 0, collator.RankValues(World, World))

	// Array
	var Universe = "Universe"
	var a0 = []any{}
	var a1 = []any{Hello, World}
	var a2 = []any{Hello, Universe}
	var a3 = []any{Hello, World, Universe}
	var a4 = []any{Hello, Universe, World}
	var aNil []any

	ass.Equal(t, 0, collator.RankValues(aNil, aNil))
	ass.Equal(t, -1, collator.RankValues(aNil, a0))
	ass.Equal(t, 1, collator.RankValues(a0, aNil))
	ass.Equal(t, 0, collator.RankValues(a0, a0))
	ass.Equal(t, 1, collator.RankValues(a1, aNil))
	ass.Equal(t, -1, collator.RankValues(a2, a1))
	ass.Equal(t, 0, collator.RankValues(a2, a2))
	ass.Equal(t, 1, collator.RankValues(a1, a2))
	ass.Equal(t, 0, collator.RankValues(a1, a1))
	ass.Equal(t, -1, collator.RankValues(a2, a3))
	ass.Equal(t, 0, collator.RankValues(a2, a2))
	ass.Equal(t, 1, collator.RankValues(a3, a2))
	ass.Equal(t, 0, collator.RankValues(a3, a3))
	ass.Equal(t, -1, collator.RankValues(a4, a1))
	ass.Equal(t, 0, collator.RankValues(a4, a4))
	ass.Equal(t, 1, collator.RankValues(a1, a4))
	ass.Equal(t, 0, collator.RankValues(a1, a1))

	// Map
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

	ass.Equal(t, 0, collator.RankValues(mNil, mNil))
	ass.Equal(t, -1, collator.RankValues(mNil, m0))
	ass.Equal(t, 1, collator.RankValues(m0, mNil))
	ass.Equal(t, 0, collator.RankValues(m0, m0))
	ass.Equal(t, -1, collator.RankValues(m2, m1))
	ass.Equal(t, 0, collator.RankValues(m2, m2))
	ass.Equal(t, 1, collator.RankValues(m1, m2))
	ass.Equal(t, 0, collator.RankValues(m1, m1))
	ass.Equal(t, -1, collator.RankValues(m2, m3))
	ass.Equal(t, 0, collator.RankValues(m2, m2))
	ass.Equal(t, 1, collator.RankValues(m3, m2))
	ass.Equal(t, 0, collator.RankValues(m3, m3))
	ass.Equal(t, -1, collator.RankValues(m4, m1))
	ass.Equal(t, 0, collator.RankValues(m4, m4))
	ass.Equal(t, 1, collator.RankValues(m1, m4))
	ass.Equal(t, 0, collator.RankValues(m1, m1))

	// Struct
	var f1 = FooBar(1, "one", nil)
	var f2 = FooBar(1, "two", nil)
	var f3 = FooBar(2, "two", nil)
	var f4 = Fuz{"two"}
	var f5 = Fuz{"two"}
	var f6 = Fuz{"three"}
	ass.Equal(t, 0, collator.RankValues(f1, f1))
	ass.Equal(t, -1, collator.RankValues(f1, f2))
	ass.Equal(t, -1, collator.RankValues(f2, f3))
	ass.Equal(t, 1, collator.RankValues(f3, f1))
	ass.Equal(t, 1, collator.RankValues(f3, f2))
	ass.Equal(t, 0, collator.RankValues(f4, f4))
	ass.Equal(t, 0, collator.RankValues(f4, f5))
	ass.Equal(t, 1, collator.RankValues(f5, f6))
	ass.Equal(t, 1, collator.RankValues(f3, &f4))
	ass.Equal(t, 0, collator.RankValues(&f4, &f4))
	ass.Equal(t, 0, collator.RankValues(&f4, &f5))
	ass.Equal(t, 1, collator.RankValues(&f5, &f6))
}

func TestTildeArrays(t *tes.T) {
	var collator = col.Collator().Make()
	var ranker = collator.RankValues
	var sorter = col.Sorter[String]().MakeWithRanker(ranker)
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
	var collator = col.Collator().Make()
	var array = col.Array[any]().MakeFromArray([]any{0})
	array.SetValue(1, array) // Now it is recursive.
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The maximum traversal depth was exceeded: 16", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	collator.RankValues(array, array) // This should panic.
}

func TestRankRecursiveMaps(t *tes.T) {
	var collator = col.Collator().Make()
	var m = col.Map[string, any]().MakeFromMap(map[string]any{"first": 1})
	m.SetValue("first", m) // Now it is recursive.
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The maximum traversal depth was exceeded: 16", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	collator.RankValues(m, m) // This should panic.
}
