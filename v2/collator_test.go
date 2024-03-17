/*
................................................................................
.    Copyright (c) 2009-2024 Crater Dog Technologies.  All Rights Reserved.    .
................................................................................
.  DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               .
.                                                                              .
.  This code is free software; you can redistribute it and/or modify it under  .
.  the terms of The MIT License (MIT), as published by the Open Source         .
.  Initiative. (See http://opensource.org/licenses/MIT)                        .
................................................................................
*/

package collections_test

import (
	col "github.com/craterdog/go-collection-framework/v2"
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

// Tilda Types
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
	// Nil
	var ShouldBeNil any

	ass.True(t, col.CompareValues(nil, nil))
	ass.True(t, col.CompareValues(nil, ShouldBeNil))
	ass.True(t, col.CompareValues(ShouldBeNil, ShouldBeNil))
	ass.True(t, col.CompareValues(ShouldBeNil, nil))

	// Boolean
	var False = false
	var True = true
	var ShouldBeFalse bool

	ass.True(t, col.CompareValues(ShouldBeFalse, False))
	ass.False(t, col.CompareValues(True, ShouldBeFalse))

	ass.False(t, col.CompareValues(False, True))
	ass.True(t, col.CompareValues(False, False))
	ass.False(t, col.CompareValues(True, False))
	ass.True(t, col.CompareValues(True, True))

	// Byte
	var Zero byte = 0x00
	var One byte = 0x01
	var ShouldBeZero byte

	ass.True(t, col.CompareValues(ShouldBeZero, Zero))
	ass.False(t, col.CompareValues(One, ShouldBeZero))

	ass.False(t, col.CompareValues(Zero, One))
	ass.True(t, col.CompareValues(Zero, Zero))
	ass.False(t, col.CompareValues(One, Zero))
	ass.True(t, col.CompareValues(One, One))

	// Integer
	var Zilch = 0
	var Two = 2
	var Three = 3
	var ShouldBeZilch int

	ass.True(t, col.CompareValues(ShouldBeZilch, Zilch))
	ass.False(t, col.CompareValues(Two, ShouldBeZilch))

	ass.False(t, col.CompareValues(Two, Three))
	ass.True(t, col.CompareValues(Two, Two))
	ass.False(t, col.CompareValues(Three, Two))
	ass.True(t, col.CompareValues(Three, Three))

	// Float
	var Negligible = 0.0
	var Fourth = 0.25
	var Half = 0.5
	var ShouldBeNegligible float64

	ass.True(t, col.CompareValues(ShouldBeNegligible, Negligible))
	ass.False(t, col.CompareValues(Half, ShouldBeNegligible))

	ass.False(t, col.CompareValues(Fourth, Half))
	ass.True(t, col.CompareValues(Fourth, Fourth))
	ass.False(t, col.CompareValues(Half, Fourth))
	ass.True(t, col.CompareValues(Half, Half))

	// Complex
	var Origin = 0 + 0i
	var PiOver4 = 1 + 1i
	var PiOver2 = 1 + 0i
	var ShouldBeOrigin complex128

	ass.True(t, col.CompareValues(ShouldBeOrigin, Origin))
	ass.False(t, col.CompareValues(PiOver4, ShouldBeOrigin))

	ass.False(t, col.CompareValues(PiOver4, PiOver2))
	ass.True(t, col.CompareValues(PiOver4, PiOver4))
	ass.False(t, col.CompareValues(PiOver2, PiOver4))
	ass.True(t, col.CompareValues(PiOver2, PiOver2))

	// Rune
	var Null = rune(0)
	var Sad = '☹'
	var Happy = '☺'
	var ShouldBeNull rune

	ass.True(t, col.CompareValues(ShouldBeNull, Null))
	ass.False(t, col.CompareValues(Sad, ShouldBeNull))

	ass.False(t, col.CompareValues(Happy, Sad))
	ass.True(t, col.CompareValues(Happy, Happy))
	ass.False(t, col.CompareValues(Sad, Happy))
	ass.True(t, col.CompareValues(Sad, Sad))

	// String
	var Empty = ""
	var Hello = "Hello"
	var World = "World"
	var ShouldBeEmpty string

	ass.True(t, col.CompareValues(ShouldBeEmpty, Empty))
	ass.False(t, col.CompareValues(Hello, ShouldBeEmpty))

	ass.False(t, col.CompareValues(World, Hello))
	ass.True(t, col.CompareValues(World, World))
	ass.False(t, col.CompareValues(Hello, World))
	ass.True(t, col.CompareValues(Hello, Hello))

	// Array
	var Universe = "Universe"
	var a0 = []any{}
	var a1 = []any{Hello, World}
	var a2 = []any{Hello, Universe}
	var aNil []any

	ass.True(t, col.CompareValues(aNil, aNil))
	ass.False(t, col.CompareValues(aNil, a0))
	ass.False(t, col.CompareValues(a0, aNil))
	ass.True(t, col.CompareValues(a0, a0))

	ass.False(t, col.CompareValues(a1, a2))
	ass.True(t, col.CompareValues(a1, a1))
	ass.False(t, col.CompareValues(a2, a1))
	ass.True(t, col.CompareValues(a2, a2))

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

	ass.True(t, col.CompareValues(mNil, mNil))
	ass.False(t, col.CompareValues(mNil, m0))
	ass.False(t, col.CompareValues(m0, mNil))
	ass.True(t, col.CompareValues(m0, m0))

	ass.False(t, col.CompareValues(m1, m2))
	ass.True(t, col.CompareValues(m1, m1))
	ass.False(t, col.CompareValues(m2, m1))
	ass.True(t, col.CompareValues(m2, m2))
	ass.False(t, col.CompareValues(m2, m3))

	// Struct
	var f0 Foolish
	var f1 = FooBar(1, "one", nil)
	var f2 = FooBar(1, "one", nil)
	var f3 = FooBar(2, "two", nil)
	var f4 = Fuz{"two"}
	var f5 = Fuz{"two"}
	var f6 = Fuz{"three"}
	ass.True(t, col.CompareValues(f0, f0))
	ass.False(t, col.CompareValues(f0, f1))
	ass.True(t, col.CompareValues(f1, f1))
	ass.True(t, col.CompareValues(f1, f2))
	ass.False(t, col.CompareValues(f2, f3))
	ass.True(t, col.CompareValues(f4, f4))
	ass.True(t, col.CompareValues(f4, f5))
	ass.False(t, col.CompareValues(f5, f6))
	ass.True(t, col.CompareValues(&f4, &f4))
	ass.True(t, col.CompareValues(&f4, &f5))
	ass.False(t, col.CompareValues(&f5, &f6))
}

func TestTildaTypes(t *tes.T) {
	// Boolean
	var False = Boolean(false)
	var True = Boolean(true)
	var ShouldBeFalse Boolean

	ass.Equal(t, 0, col.RankValues(ShouldBeFalse, ShouldBeFalse))
	ass.Equal(t, -1, col.RankValues(ShouldBeFalse, True))
	ass.Equal(t, 0, col.RankValues(False, ShouldBeFalse))
	ass.Equal(t, 1, col.RankValues(True, ShouldBeFalse))
	ass.Equal(t, 0, col.RankValues(ShouldBeFalse, False))
	ass.Equal(t, -1, col.RankValues(False, True))
	ass.Equal(t, 0, col.RankValues(False, False))
	ass.Equal(t, 1, col.RankValues(True, False))
	ass.Equal(t, 0, col.RankValues(True, True))

	// Integer
	var Zilch = Integer(0)
	var Two = Integer(2)
	var Three = Integer(3)
	var ShouldBeZilch Integer

	ass.True(t, col.CompareValues(ShouldBeZilch, Zilch))
	ass.False(t, col.CompareValues(Two, ShouldBeZilch))

	ass.False(t, col.CompareValues(Two, Three))
	ass.True(t, col.CompareValues(Two, Two))
	ass.False(t, col.CompareValues(Three, Two))
	ass.True(t, col.CompareValues(Three, Three))

	// String
	var Empty = String("")
	var Hello = String("Hello")
	var World = String("World")
	var ShouldBeEmpty String

	ass.True(t, col.CompareValues(ShouldBeEmpty, Empty))
	ass.False(t, col.CompareValues(Hello, ShouldBeEmpty))

	ass.False(t, col.CompareValues(World, Hello))
	ass.True(t, col.CompareValues(World, World))
	ass.False(t, col.CompareValues(Hello, World))
	ass.True(t, col.CompareValues(Hello, Hello))
}

func TestCompareRecursiveArrays(t *tes.T) {
	var array = col.Array[any]([]any{0})
	array.SetValue(1, array) // Now it is recursive.
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The maximum recursion depth was exceeded: 100", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	col.CompareValues(array, array) // This should panic.
}

func TestCompareRecursiveMaps(t *tes.T) {
	var m = col.Map[string, any](map[string]any{"first": 1})
	m.SetValue("first", m) // Now it is recursive.
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The maximum recursion depth was exceeded: 100", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	col.CompareValues(m, m) // This should panic.
}

func TestRanking(t *tes.T) {
	// Nil
	var ShouldBeNil any

	ass.Equal(t, 0, col.RankValues(nil, nil))
	ass.Equal(t, 0, col.RankValues(nil, ShouldBeNil))
	ass.Equal(t, 0, col.RankValues(ShouldBeNil, ShouldBeNil))
	ass.Equal(t, 0, col.RankValues(ShouldBeNil, nil))

	// Boolean
	var False = false
	var True = true
	var ShouldBeFalse bool

	ass.Equal(t, 0, col.RankValues(ShouldBeFalse, ShouldBeFalse))
	ass.Equal(t, -1, col.RankValues(ShouldBeFalse, True))
	ass.Equal(t, 0, col.RankValues(False, ShouldBeFalse))
	ass.Equal(t, 1, col.RankValues(True, ShouldBeFalse))
	ass.Equal(t, 0, col.RankValues(ShouldBeFalse, False))
	ass.Equal(t, -1, col.RankValues(False, True))
	ass.Equal(t, 0, col.RankValues(False, False))
	ass.Equal(t, 1, col.RankValues(True, False))
	ass.Equal(t, 0, col.RankValues(True, True))

	// Byte
	var Zero byte = 0x00
	var One byte = 0x01
	var ShouldBeZero byte

	ass.Equal(t, 0, col.RankValues(ShouldBeZero, ShouldBeZero))
	ass.Equal(t, -1, col.RankValues(ShouldBeZero, One))
	ass.Equal(t, 0, col.RankValues(Zero, ShouldBeZero))
	ass.Equal(t, 1, col.RankValues(One, ShouldBeZero))
	ass.Equal(t, 0, col.RankValues(ShouldBeZero, Zero))
	ass.Equal(t, -1, col.RankValues(Zero, One))
	ass.Equal(t, 0, col.RankValues(Zero, Zero))
	ass.Equal(t, 1, col.RankValues(One, Zero))
	ass.Equal(t, 0, col.RankValues(One, One))

	// Integer
	var Zilch = 0
	var Two = 2
	var Three = 3
	var ShouldBeZilch int

	ass.Equal(t, 0, col.RankValues(ShouldBeZilch, ShouldBeZilch))
	ass.Equal(t, -1, col.RankValues(ShouldBeZilch, Two))
	ass.Equal(t, 0, col.RankValues(Zilch, ShouldBeZilch))
	ass.Equal(t, 1, col.RankValues(Two, ShouldBeZilch))
	ass.Equal(t, 0, col.RankValues(ShouldBeZilch, Zilch))
	ass.Equal(t, -1, col.RankValues(Two, Three))
	ass.Equal(t, 0, col.RankValues(Two, Two))
	ass.Equal(t, 1, col.RankValues(Three, Two))
	ass.Equal(t, 0, col.RankValues(Three, Three))

	// Float
	var Negligible = 0.0
	var Fourth = 0.25
	var Half = 0.5
	var ShouldBeNegligible float64

	ass.Equal(t, 0, col.RankValues(ShouldBeNegligible, ShouldBeNegligible))
	ass.Equal(t, -1, col.RankValues(ShouldBeNegligible, Half))
	ass.Equal(t, 0, col.RankValues(Negligible, ShouldBeNegligible))
	ass.Equal(t, 1, col.RankValues(Half, ShouldBeNegligible))
	ass.Equal(t, 0, col.RankValues(ShouldBeNegligible, Negligible))
	ass.Equal(t, -1, col.RankValues(Fourth, Half))
	ass.Equal(t, 0, col.RankValues(Fourth, Fourth))
	ass.Equal(t, 1, col.RankValues(Half, Fourth))
	ass.Equal(t, 0, col.RankValues(Half, Half))

	// Complex
	var Origin = 0 + 0i
	var PiOver4 = 1 + 1i
	var PiOver2 = 1 + 0i
	var ShouldBeOrigin complex128

	ass.Equal(t, 0, col.RankValues(ShouldBeOrigin, ShouldBeOrigin))
	ass.Equal(t, -1, col.RankValues(ShouldBeOrigin, PiOver4))
	ass.Equal(t, 0, col.RankValues(Origin, ShouldBeOrigin))
	ass.Equal(t, 1, col.RankValues(PiOver4, ShouldBeOrigin))
	ass.Equal(t, 0, col.RankValues(ShouldBeOrigin, Origin))
	ass.Equal(t, -1, col.RankValues(PiOver2, PiOver4))
	ass.Equal(t, 0, col.RankValues(PiOver2, PiOver2))
	ass.Equal(t, 1, col.RankValues(PiOver4, PiOver2))
	ass.Equal(t, 0, col.RankValues(PiOver4, PiOver4))

	// Rune
	var Null = rune(0)
	var Sad = '☹'
	var Happy = '☺'
	var ShouldBeNull rune

	ass.Equal(t, 0, col.RankValues(ShouldBeNull, ShouldBeNull))
	ass.Equal(t, -1, col.RankValues(ShouldBeNull, Sad))
	ass.Equal(t, 0, col.RankValues(Null, ShouldBeNull))
	ass.Equal(t, 1, col.RankValues(Sad, ShouldBeNull))
	ass.Equal(t, 0, col.RankValues(ShouldBeNull, Null))
	ass.Equal(t, -1, col.RankValues(Sad, Happy))
	ass.Equal(t, 0, col.RankValues(Sad, Sad))
	ass.Equal(t, 1, col.RankValues(Happy, Sad))
	ass.Equal(t, 0, col.RankValues(Happy, Happy))

	// String
	var Empty = ""
	var Hello = "Hello"
	var World = "World"
	var ShouldBeEmpty string

	ass.Equal(t, 0, col.RankValues(ShouldBeEmpty, ShouldBeEmpty))
	ass.Equal(t, -1, col.RankValues(ShouldBeEmpty, Hello))
	ass.Equal(t, 0, col.RankValues(Empty, ShouldBeEmpty))
	ass.Equal(t, 1, col.RankValues(Hello, ShouldBeEmpty))
	ass.Equal(t, 0, col.RankValues(ShouldBeEmpty, Empty))
	ass.Equal(t, -1, col.RankValues(Hello, World))
	ass.Equal(t, 0, col.RankValues(Hello, Hello))
	ass.Equal(t, 1, col.RankValues(World, Hello))
	ass.Equal(t, 0, col.RankValues(World, World))

	// Array
	var Universe = "Universe"
	var a0 = []any{}
	var a1 = []any{Hello, World}
	var a2 = []any{Hello, Universe}
	var a3 = []any{Hello, World, Universe}
	var a4 = []any{Hello, Universe, World}
	var aNil []any

	ass.Equal(t, 0, col.RankValues(aNil, aNil))
	ass.Equal(t, -1, col.RankValues(aNil, a0))
	ass.Equal(t, 1, col.RankValues(a0, aNil))
	ass.Equal(t, 0, col.RankValues(a0, a0))
	ass.Equal(t, 1, col.RankValues(a1, aNil))
	ass.Equal(t, -1, col.RankValues(a2, a1))
	ass.Equal(t, 0, col.RankValues(a2, a2))
	ass.Equal(t, 1, col.RankValues(a1, a2))
	ass.Equal(t, 0, col.RankValues(a1, a1))
	ass.Equal(t, -1, col.RankValues(a2, a3))
	ass.Equal(t, 0, col.RankValues(a2, a2))
	ass.Equal(t, 1, col.RankValues(a3, a2))
	ass.Equal(t, 0, col.RankValues(a3, a3))
	ass.Equal(t, -1, col.RankValues(a4, a1))
	ass.Equal(t, 0, col.RankValues(a4, a4))
	ass.Equal(t, 1, col.RankValues(a1, a4))
	ass.Equal(t, 0, col.RankValues(a1, a1))

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

	ass.Equal(t, 0, col.RankValues(mNil, mNil))
	ass.Equal(t, -1, col.RankValues(mNil, m0))
	ass.Equal(t, 1, col.RankValues(m0, mNil))
	ass.Equal(t, 0, col.RankValues(m0, m0))
	ass.Equal(t, -1, col.RankValues(m2, m1))
	ass.Equal(t, 0, col.RankValues(m2, m2))
	ass.Equal(t, 1, col.RankValues(m1, m2))
	ass.Equal(t, 0, col.RankValues(m1, m1))
	ass.Equal(t, -1, col.RankValues(m2, m3))
	ass.Equal(t, 0, col.RankValues(m2, m2))
	ass.Equal(t, 1, col.RankValues(m3, m2))
	ass.Equal(t, 0, col.RankValues(m3, m3))
	ass.Equal(t, -1, col.RankValues(m4, m1))
	ass.Equal(t, 0, col.RankValues(m4, m4))
	ass.Equal(t, 1, col.RankValues(m1, m4))
	ass.Equal(t, 0, col.RankValues(m1, m1))

	// Struct
	var f1 = FooBar(1, "one", nil)
	var f2 = FooBar(1, "two", nil)
	var f3 = FooBar(2, "two", nil)
	var f4 = Fuz{"two"}
	var f5 = Fuz{"two"}
	var f6 = Fuz{"three"}
	ass.Equal(t, 0, col.RankValues(f1, f1))
	ass.Equal(t, -1, col.RankValues(f1, f2))
	ass.Equal(t, -1, col.RankValues(f2, f3))
	ass.Equal(t, 1, col.RankValues(f3, f1))
	ass.Equal(t, 1, col.RankValues(f3, f2))
	ass.Equal(t, 0, col.RankValues(f4, f4))
	ass.Equal(t, 0, col.RankValues(f4, f5))
	ass.Equal(t, 1, col.RankValues(f5, f6))
	ass.Equal(t, 1, col.RankValues(f3, &f4))
	ass.Equal(t, 0, col.RankValues(&f4, &f4))
	ass.Equal(t, 0, col.RankValues(&f4, &f5))
	ass.Equal(t, 1, col.RankValues(&f5, &f6))
}

func TestTildaArrays(t *tes.T) {
	var alpha = String("alpha")
	var beta = String("beta")
	var gamma = String("gamma")
	var delta = String("delta")
	var array = []String{alpha, beta, gamma, delta}
	col.SortArray(array, col.RankValues)
	ass.Equal(t, alpha, array[0])
	ass.Equal(t, beta, array[1])
	ass.Equal(t, delta, array[2])
	ass.Equal(t, gamma, array[3])
}

func TestRankRecursiveArrays(t *tes.T) {
	var array = col.Array[any]([]any{0})
	array.SetValue(1, array) // Now it is recursive.
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The maximum recursion depth was exceeded: 100", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	col.RankValues(array, array) // This should panic.
}

func TestRankRecursiveMaps(t *tes.T) {
	var m = col.Map[string, any](map[string]any{"first": 1})
	m.SetValue("first", m) // Now it is recursive.
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The maximum recursion depth was exceeded: 100", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	col.RankValues(m, m) // This should panic.
}
