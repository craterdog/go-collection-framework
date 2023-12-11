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

func TestComparison(t *tes.T) {
	var Collator = col.Collator()

	// Nil
	var ShouldBeNil any

	ass.True(t, Collator.CompareValues(nil, nil))
	ass.True(t, Collator.CompareValues(nil, ShouldBeNil))
	ass.True(t, Collator.CompareValues(ShouldBeNil, ShouldBeNil))
	ass.True(t, Collator.CompareValues(ShouldBeNil, nil))

	// Boolean
	var False = false
	var True = true
	var ShouldBeFalse bool

	ass.True(t, Collator.CompareValues(ShouldBeFalse, False))
	ass.False(t, Collator.CompareValues(True, ShouldBeFalse))

	ass.False(t, Collator.CompareValues(False, True))
	ass.True(t, Collator.CompareValues(False, False))
	ass.False(t, Collator.CompareValues(True, False))
	ass.True(t, Collator.CompareValues(True, True))

	// Byte
	var Zero byte = 0x00
	var One byte = 0x01
	var ShouldBeZero byte

	ass.True(t, Collator.CompareValues(ShouldBeZero, Zero))
	ass.False(t, Collator.CompareValues(One, ShouldBeZero))

	ass.False(t, Collator.CompareValues(Zero, One))
	ass.True(t, Collator.CompareValues(Zero, Zero))
	ass.False(t, Collator.CompareValues(One, Zero))
	ass.True(t, Collator.CompareValues(One, One))

	// Integer
	var Zilch = 0
	var Two = 2
	var Three = 3
	var ShouldBeZilch int

	ass.True(t, Collator.CompareValues(ShouldBeZilch, Zilch))
	ass.False(t, Collator.CompareValues(Two, ShouldBeZilch))

	ass.False(t, Collator.CompareValues(Two, Three))
	ass.True(t, Collator.CompareValues(Two, Two))
	ass.False(t, Collator.CompareValues(Three, Two))
	ass.True(t, Collator.CompareValues(Three, Three))

	// Float
	var Negligible = 0.0
	var Fourth = 0.25
	var Half = 0.5
	var ShouldBeNegligible float64

	ass.True(t, Collator.CompareValues(ShouldBeNegligible, Negligible))
	ass.False(t, Collator.CompareValues(Half, ShouldBeNegligible))

	ass.False(t, Collator.CompareValues(Fourth, Half))
	ass.True(t, Collator.CompareValues(Fourth, Fourth))
	ass.False(t, Collator.CompareValues(Half, Fourth))
	ass.True(t, Collator.CompareValues(Half, Half))

	// Complex
	var Origin = 0 + 0i
	var PiOver4 = 1 + 1i
	var PiOver2 = 1 + 0i
	var ShouldBeOrigin complex128

	ass.True(t, Collator.CompareValues(ShouldBeOrigin, Origin))
	ass.False(t, Collator.CompareValues(PiOver4, ShouldBeOrigin))

	ass.False(t, Collator.CompareValues(PiOver4, PiOver2))
	ass.True(t, Collator.CompareValues(PiOver4, PiOver4))
	ass.False(t, Collator.CompareValues(PiOver2, PiOver4))
	ass.True(t, Collator.CompareValues(PiOver2, PiOver2))

	// Rune
	var Null = rune(0)
	var Sad = '☹'
	var Happy = '☺'
	var ShouldBeNull rune

	ass.True(t, Collator.CompareValues(ShouldBeNull, Null))
	ass.False(t, Collator.CompareValues(Sad, ShouldBeNull))

	ass.False(t, Collator.CompareValues(Happy, Sad))
	ass.True(t, Collator.CompareValues(Happy, Happy))
	ass.False(t, Collator.CompareValues(Sad, Happy))
	ass.True(t, Collator.CompareValues(Sad, Sad))

	// String
	var Empty = ""
	var Hello = "Hello"
	var World = "World"
	var ShouldBeEmpty string

	ass.True(t, Collator.CompareValues(ShouldBeEmpty, Empty))
	ass.False(t, Collator.CompareValues(Hello, ShouldBeEmpty))

	ass.False(t, Collator.CompareValues(World, Hello))
	ass.True(t, Collator.CompareValues(World, World))
	ass.False(t, Collator.CompareValues(Hello, World))
	ass.True(t, Collator.CompareValues(Hello, Hello))

	// Array
	var Universe = "Universe"
	var a0 = []any{}
	var a1 = []any{Hello, World}
	var a2 = []any{Hello, Universe}
	var aNil []any

	ass.True(t, Collator.CompareValues(aNil, aNil))
	ass.False(t, Collator.CompareValues(aNil, a0))
	ass.False(t, Collator.CompareValues(a0, aNil))
	ass.True(t, Collator.CompareValues(a0, a0))

	ass.False(t, Collator.CompareValues(a1, a2))
	ass.True(t, Collator.CompareValues(a1, a1))
	ass.False(t, Collator.CompareValues(a2, a1))
	ass.True(t, Collator.CompareValues(a2, a2))

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

	ass.True(t, Collator.CompareValues(mNil, mNil))
	ass.False(t, Collator.CompareValues(mNil, m0))
	ass.False(t, Collator.CompareValues(m0, mNil))
	ass.True(t, Collator.CompareValues(m0, m0))

	ass.False(t, Collator.CompareValues(m1, m2))
	ass.True(t, Collator.CompareValues(m1, m1))
	ass.False(t, Collator.CompareValues(m2, m1))
	ass.True(t, Collator.CompareValues(m2, m2))
	ass.False(t, Collator.CompareValues(m2, m3))

	// Struct
	var f0 Foolish
	var f1 = FooBar(1, "one", nil)
	var f2 = FooBar(1, "one", nil)
	var f3 = FooBar(2, "two", nil)
	var f4 = Fuz{"two"}
	var f5 = Fuz{"two"}
	var f6 = Fuz{"three"}
	ass.True(t, Collator.CompareValues(f0, f0))
	ass.False(t, Collator.CompareValues(f0, f1))
	ass.True(t, Collator.CompareValues(f1, f1))
	ass.True(t, Collator.CompareValues(f1, f2))
	ass.False(t, Collator.CompareValues(f2, f3))
	ass.True(t, Collator.CompareValues(f4, f4))
	ass.True(t, Collator.CompareValues(f4, f5))
	ass.False(t, Collator.CompareValues(f5, f6))
	ass.True(t, Collator.CompareValues(&f4, &f4))
	ass.True(t, Collator.CompareValues(&f4, &f5))
	ass.False(t, Collator.CompareValues(&f5, &f6))
}

func TestTildeTypes(t *tes.T) {
	var Collator = col.Collator()

	// Boolean
	var False = Boolean(false)
	var True = Boolean(true)
	var ShouldBeFalse Boolean

	ass.Equal(t, 0, Collator.RankValues(ShouldBeFalse, ShouldBeFalse))
	ass.Equal(t, -1, Collator.RankValues(ShouldBeFalse, True))
	ass.Equal(t, 0, Collator.RankValues(False, ShouldBeFalse))
	ass.Equal(t, 1, Collator.RankValues(True, ShouldBeFalse))
	ass.Equal(t, 0, Collator.RankValues(ShouldBeFalse, False))
	ass.Equal(t, -1, Collator.RankValues(False, True))
	ass.Equal(t, 0, Collator.RankValues(False, False))
	ass.Equal(t, 1, Collator.RankValues(True, False))
	ass.Equal(t, 0, Collator.RankValues(True, True))

	// Integer
	var Zilch = Integer(0)
	var Two = Integer(2)
	var Three = Integer(3)
	var ShouldBeZilch Integer

	ass.True(t, Collator.CompareValues(ShouldBeZilch, Zilch))
	ass.False(t, Collator.CompareValues(Two, ShouldBeZilch))

	ass.False(t, Collator.CompareValues(Two, Three))
	ass.True(t, Collator.CompareValues(Two, Two))
	ass.False(t, Collator.CompareValues(Three, Two))
	ass.True(t, Collator.CompareValues(Three, Three))

	// String
	var Empty = String("")
	var Hello = String("Hello")
	var World = String("World")
	var ShouldBeEmpty String

	ass.True(t, Collator.CompareValues(ShouldBeEmpty, Empty))
	ass.False(t, Collator.CompareValues(Hello, ShouldBeEmpty))

	ass.False(t, Collator.CompareValues(World, Hello))
	ass.True(t, Collator.CompareValues(World, World))
	ass.False(t, Collator.CompareValues(Hello, World))
	ass.True(t, Collator.CompareValues(Hello, Hello))
}

func TestCompareRecursiveArrays(t *tes.T) {
	var Collator = col.Collator()
	var Array = col.Array[any]()
	var array = Array.FromArray([]any{0})
	array.SetValue(1, array) // Now it is recursive.
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The maximum recursion depth was exceeded: 16", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	Collator.CompareValues(array, array) // This should panic.
}

func TestCompareRecursiveMaps(t *tes.T) {
	var Collator = col.Collator()
	var Map = col.Map[string, any]()
	var m = Map.FromMap(map[string]any{"first": 1})
	m.SetValue("first", m) // Now it is recursive.
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The maximum recursion depth was exceeded: 16", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	Collator.CompareValues(m, m) // This should panic.
}

func TestRanking(t *tes.T) {
	var Collator = col.Collator()

	// Nil
	var ShouldBeNil any

	ass.Equal(t, 0, Collator.RankValues(nil, nil))
	ass.Equal(t, 0, Collator.RankValues(nil, ShouldBeNil))
	ass.Equal(t, 0, Collator.RankValues(ShouldBeNil, ShouldBeNil))
	ass.Equal(t, 0, Collator.RankValues(ShouldBeNil, nil))

	// Boolean
	var False = false
	var True = true
	var ShouldBeFalse bool

	ass.Equal(t, 0, Collator.RankValues(ShouldBeFalse, ShouldBeFalse))
	ass.Equal(t, -1, Collator.RankValues(ShouldBeFalse, True))
	ass.Equal(t, 0, Collator.RankValues(False, ShouldBeFalse))
	ass.Equal(t, 1, Collator.RankValues(True, ShouldBeFalse))
	ass.Equal(t, 0, Collator.RankValues(ShouldBeFalse, False))
	ass.Equal(t, -1, Collator.RankValues(False, True))
	ass.Equal(t, 0, Collator.RankValues(False, False))
	ass.Equal(t, 1, Collator.RankValues(True, False))
	ass.Equal(t, 0, Collator.RankValues(True, True))

	// Byte
	var Zero byte = 0x00
	var One byte = 0x01
	var ShouldBeZero byte

	ass.Equal(t, 0, Collator.RankValues(ShouldBeZero, ShouldBeZero))
	ass.Equal(t, -1, Collator.RankValues(ShouldBeZero, One))
	ass.Equal(t, 0, Collator.RankValues(Zero, ShouldBeZero))
	ass.Equal(t, 1, Collator.RankValues(One, ShouldBeZero))
	ass.Equal(t, 0, Collator.RankValues(ShouldBeZero, Zero))
	ass.Equal(t, -1, Collator.RankValues(Zero, One))
	ass.Equal(t, 0, Collator.RankValues(Zero, Zero))
	ass.Equal(t, 1, Collator.RankValues(One, Zero))
	ass.Equal(t, 0, Collator.RankValues(One, One))

	// Integer
	var Zilch = 0
	var Two = 2
	var Three = 3
	var ShouldBeZilch int

	ass.Equal(t, 0, Collator.RankValues(ShouldBeZilch, ShouldBeZilch))
	ass.Equal(t, -1, Collator.RankValues(ShouldBeZilch, Two))
	ass.Equal(t, 0, Collator.RankValues(Zilch, ShouldBeZilch))
	ass.Equal(t, 1, Collator.RankValues(Two, ShouldBeZilch))
	ass.Equal(t, 0, Collator.RankValues(ShouldBeZilch, Zilch))
	ass.Equal(t, -1, Collator.RankValues(Two, Three))
	ass.Equal(t, 0, Collator.RankValues(Two, Two))
	ass.Equal(t, 1, Collator.RankValues(Three, Two))
	ass.Equal(t, 0, Collator.RankValues(Three, Three))

	// Float
	var Negligible = 0.0
	var Fourth = 0.25
	var Half = 0.5
	var ShouldBeNegligible float64

	ass.Equal(t, 0, Collator.RankValues(ShouldBeNegligible, ShouldBeNegligible))
	ass.Equal(t, -1, Collator.RankValues(ShouldBeNegligible, Half))
	ass.Equal(t, 0, Collator.RankValues(Negligible, ShouldBeNegligible))
	ass.Equal(t, 1, Collator.RankValues(Half, ShouldBeNegligible))
	ass.Equal(t, 0, Collator.RankValues(ShouldBeNegligible, Negligible))
	ass.Equal(t, -1, Collator.RankValues(Fourth, Half))
	ass.Equal(t, 0, Collator.RankValues(Fourth, Fourth))
	ass.Equal(t, 1, Collator.RankValues(Half, Fourth))
	ass.Equal(t, 0, Collator.RankValues(Half, Half))

	// Complex
	var Origin = 0 + 0i
	var PiOver4 = 1 + 1i
	var PiOver2 = 1 + 0i
	var ShouldBeOrigin complex128

	ass.Equal(t, 0, Collator.RankValues(ShouldBeOrigin, ShouldBeOrigin))
	ass.Equal(t, -1, Collator.RankValues(ShouldBeOrigin, PiOver4))
	ass.Equal(t, 0, Collator.RankValues(Origin, ShouldBeOrigin))
	ass.Equal(t, 1, Collator.RankValues(PiOver4, ShouldBeOrigin))
	ass.Equal(t, 0, Collator.RankValues(ShouldBeOrigin, Origin))
	ass.Equal(t, -1, Collator.RankValues(PiOver2, PiOver4))
	ass.Equal(t, 0, Collator.RankValues(PiOver2, PiOver2))
	ass.Equal(t, 1, Collator.RankValues(PiOver4, PiOver2))
	ass.Equal(t, 0, Collator.RankValues(PiOver4, PiOver4))

	// Rune
	var Null = rune(0)
	var Sad = '☹'
	var Happy = '☺'
	var ShouldBeNull rune

	ass.Equal(t, 0, Collator.RankValues(ShouldBeNull, ShouldBeNull))
	ass.Equal(t, -1, Collator.RankValues(ShouldBeNull, Sad))
	ass.Equal(t, 0, Collator.RankValues(Null, ShouldBeNull))
	ass.Equal(t, 1, Collator.RankValues(Sad, ShouldBeNull))
	ass.Equal(t, 0, Collator.RankValues(ShouldBeNull, Null))
	ass.Equal(t, -1, Collator.RankValues(Sad, Happy))
	ass.Equal(t, 0, Collator.RankValues(Sad, Sad))
	ass.Equal(t, 1, Collator.RankValues(Happy, Sad))
	ass.Equal(t, 0, Collator.RankValues(Happy, Happy))

	// String
	var Empty = ""
	var Hello = "Hello"
	var World = "World"
	var ShouldBeEmpty string

	ass.Equal(t, 0, Collator.RankValues(ShouldBeEmpty, ShouldBeEmpty))
	ass.Equal(t, -1, Collator.RankValues(ShouldBeEmpty, Hello))
	ass.Equal(t, 0, Collator.RankValues(Empty, ShouldBeEmpty))
	ass.Equal(t, 1, Collator.RankValues(Hello, ShouldBeEmpty))
	ass.Equal(t, 0, Collator.RankValues(ShouldBeEmpty, Empty))
	ass.Equal(t, -1, Collator.RankValues(Hello, World))
	ass.Equal(t, 0, Collator.RankValues(Hello, Hello))
	ass.Equal(t, 1, Collator.RankValues(World, Hello))
	ass.Equal(t, 0, Collator.RankValues(World, World))

	// Array
	var Universe = "Universe"
	var a0 = []any{}
	var a1 = []any{Hello, World}
	var a2 = []any{Hello, Universe}
	var a3 = []any{Hello, World, Universe}
	var a4 = []any{Hello, Universe, World}
	var aNil []any

	ass.Equal(t, 0, Collator.RankValues(aNil, aNil))
	ass.Equal(t, -1, Collator.RankValues(aNil, a0))
	ass.Equal(t, 1, Collator.RankValues(a0, aNil))
	ass.Equal(t, 0, Collator.RankValues(a0, a0))
	ass.Equal(t, 1, Collator.RankValues(a1, aNil))
	ass.Equal(t, -1, Collator.RankValues(a2, a1))
	ass.Equal(t, 0, Collator.RankValues(a2, a2))
	ass.Equal(t, 1, Collator.RankValues(a1, a2))
	ass.Equal(t, 0, Collator.RankValues(a1, a1))
	ass.Equal(t, -1, Collator.RankValues(a2, a3))
	ass.Equal(t, 0, Collator.RankValues(a2, a2))
	ass.Equal(t, 1, Collator.RankValues(a3, a2))
	ass.Equal(t, 0, Collator.RankValues(a3, a3))
	ass.Equal(t, -1, Collator.RankValues(a4, a1))
	ass.Equal(t, 0, Collator.RankValues(a4, a4))
	ass.Equal(t, 1, Collator.RankValues(a1, a4))
	ass.Equal(t, 0, Collator.RankValues(a1, a1))

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

	ass.Equal(t, 0, Collator.RankValues(mNil, mNil))
	ass.Equal(t, -1, Collator.RankValues(mNil, m0))
	ass.Equal(t, 1, Collator.RankValues(m0, mNil))
	ass.Equal(t, 0, Collator.RankValues(m0, m0))
	ass.Equal(t, -1, Collator.RankValues(m2, m1))
	ass.Equal(t, 0, Collator.RankValues(m2, m2))
	ass.Equal(t, 1, Collator.RankValues(m1, m2))
	ass.Equal(t, 0, Collator.RankValues(m1, m1))
	ass.Equal(t, -1, Collator.RankValues(m2, m3))
	ass.Equal(t, 0, Collator.RankValues(m2, m2))
	ass.Equal(t, 1, Collator.RankValues(m3, m2))
	ass.Equal(t, 0, Collator.RankValues(m3, m3))
	ass.Equal(t, -1, Collator.RankValues(m4, m1))
	ass.Equal(t, 0, Collator.RankValues(m4, m4))
	ass.Equal(t, 1, Collator.RankValues(m1, m4))
	ass.Equal(t, 0, Collator.RankValues(m1, m1))

	// Struct
	var f1 = FooBar(1, "one", nil)
	var f2 = FooBar(1, "two", nil)
	var f3 = FooBar(2, "two", nil)
	var f4 = Fuz{"two"}
	var f5 = Fuz{"two"}
	var f6 = Fuz{"three"}
	ass.Equal(t, 0, Collator.RankValues(f1, f1))
	ass.Equal(t, -1, Collator.RankValues(f1, f2))
	ass.Equal(t, -1, Collator.RankValues(f2, f3))
	ass.Equal(t, 1, Collator.RankValues(f3, f1))
	ass.Equal(t, 1, Collator.RankValues(f3, f2))
	ass.Equal(t, 0, Collator.RankValues(f4, f4))
	ass.Equal(t, 0, Collator.RankValues(f4, f5))
	ass.Equal(t, 1, Collator.RankValues(f5, f6))
	ass.Equal(t, 1, Collator.RankValues(f3, &f4))
	ass.Equal(t, 0, Collator.RankValues(&f4, &f4))
	ass.Equal(t, 0, Collator.RankValues(&f4, &f5))
	ass.Equal(t, 1, Collator.RankValues(&f5, &f6))
}

func TestTildeArrays(t *tes.T) {
	var Sorter = col.Sorter[String]()
	var Collator = col.Collator()
	var alpha = String("alpha")
	var beta = String("beta")
	var gamma = String("gamma")
	var delta = String("delta")
	var array = []String{alpha, beta, gamma, delta}
	Sorter.SortValues(array, Collator.RankValues)
	ass.Equal(t, alpha, array[0])
	ass.Equal(t, beta, array[1])
	ass.Equal(t, delta, array[2])
	ass.Equal(t, gamma, array[3])
}

func TestRankRecursiveArrays(t *tes.T) {
	var Collator = col.Collator()
	var Array = col.Array[any]()
	var array = Array.FromArray([]any{0})
	array.SetValue(1, array) // Now it is recursive.
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The maximum recursion depth was exceeded: 16", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	Collator.RankValues(array, array) // This should panic.
}

func TestRankRecursiveMaps(t *tes.T) {
	var Collator = col.Collator()
	var Map = col.Map[string, any]()
	var m = Map.FromMap(map[string]any{"first": 1})
	m.SetValue("first", m) // Now it is recursive.
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The maximum recursion depth was exceeded: 16", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	Collator.RankValues(m, m) // This should panic.
}
