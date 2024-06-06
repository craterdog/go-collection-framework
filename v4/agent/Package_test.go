/*
................................................................................
.    Copyright (c) 2009-2024 Crater Dog Technologies.  All Rights Reserved.    .
................................................................................
.  DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               .
.                                                                              .
.  This code is free software; you can redistribute it and/or modify it under  .
.  the terms of The MIT License (MIT), as published by the Open Source         .
.  Initiative. (See https://opensource.org/license/MIT)                        .
................................................................................
*/

package agent_test

import (
	age "github.com/craterdog/go-collection-framework/v4/agent"
	not "github.com/craterdog/go-collection-framework/v4/cdcn"
	col "github.com/craterdog/go-collection-framework/v4/collection"
	ass "github.com/stretchr/testify/assert"
	tes "testing"
)

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

func TestCollatorConstants(t *tes.T) {
	var Collator = age.Collator[any]()
	ass.Equal(t, 16, Collator.DefaultMaximum())
}

func TestCompareMaximum(t *tes.T) {
	var notation = not.Notation().Make()
	var collator = age.Collator[any]().MakeWithMaximum(1)
	var array = col.Array[any](notation).MakeFromArray([]any{"foo", []int{1, 2, 3}})
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The maximum traversal depth was exceeded: 1", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	collator.CompareValues(array, array)
}

func TestRankMaximum(t *tes.T) {
	var notation = not.Notation().Make()
	var collator = age.Collator[any]().MakeWithMaximum(1)
	var array = col.Array[any](notation).MakeFromArray([]any{"foo", []int{1, 2, 3}})
	defer func() {
		if e := recover(); e != nil {
			ass.Equal(t, "The maximum traversal depth was exceeded: 1", e)
		} else {
			ass.Fail(t, "Test should result in recovered panic.")
		}
	}()
	collator.RankValues(array, array)
}

func TestComparison(t *tes.T) {
	var collator = age.Collator[any]().Make()

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
	var collator = age.Collator[any]().Make()

	// Boolean
	var False = Boolean(false)
	var True = Boolean(true)
	var ShouldBeFalse Boolean

	ass.Equal(t, age.EqualRank, collator.RankValues(ShouldBeFalse, ShouldBeFalse))
	ass.Equal(t, age.LesserRank, collator.RankValues(ShouldBeFalse, True))
	ass.Equal(t, age.EqualRank, collator.RankValues(False, ShouldBeFalse))
	ass.Equal(t, age.GreaterRank, collator.RankValues(True, ShouldBeFalse))
	ass.Equal(t, age.EqualRank, collator.RankValues(ShouldBeFalse, False))
	ass.Equal(t, age.LesserRank, collator.RankValues(False, True))
	ass.Equal(t, age.EqualRank, collator.RankValues(False, False))
	ass.Equal(t, age.GreaterRank, collator.RankValues(True, False))
	ass.Equal(t, age.EqualRank, collator.RankValues(True, True))

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
	var notation = not.Notation().Make()
	var collator = age.Collator[any]().Make()
	var array = col.Array[any](notation).MakeFromArray([]any{0})
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
	var notation = not.Notation().Make()
	var collator = age.Collator[any]().Make()
	var m = col.Map[string, any](notation).MakeFromMap(map[string]any{"first": 1})
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
	var collator = age.Collator[any]().Make()

	// Nil
	var ShouldBeNil any

	ass.Equal(t, age.EqualRank, collator.RankValues(nil, nil))
	ass.Equal(t, age.EqualRank, collator.RankValues(nil, ShouldBeNil))
	ass.Equal(t, age.EqualRank, collator.RankValues(ShouldBeNil, ShouldBeNil))
	ass.Equal(t, age.EqualRank, collator.RankValues(ShouldBeNil, nil))

	// Boolean
	var False = false
	var True = true
	var ShouldBeFalse bool

	ass.Equal(t, age.EqualRank, collator.RankValues(ShouldBeFalse, ShouldBeFalse))
	ass.Equal(t, age.LesserRank, collator.RankValues(ShouldBeFalse, True))
	ass.Equal(t, age.EqualRank, collator.RankValues(False, ShouldBeFalse))
	ass.Equal(t, age.GreaterRank, collator.RankValues(True, ShouldBeFalse))
	ass.Equal(t, age.EqualRank, collator.RankValues(ShouldBeFalse, False))
	ass.Equal(t, age.LesserRank, collator.RankValues(False, True))
	ass.Equal(t, age.EqualRank, collator.RankValues(False, False))
	ass.Equal(t, age.GreaterRank, collator.RankValues(True, False))
	ass.Equal(t, age.EqualRank, collator.RankValues(True, True))

	// Byte
	var Zero byte = 0x00
	var One byte = 0x01
	var ShouldBeZero byte

	ass.Equal(t, age.EqualRank, collator.RankValues(ShouldBeZero, ShouldBeZero))
	ass.Equal(t, age.LesserRank, collator.RankValues(ShouldBeZero, One))
	ass.Equal(t, age.EqualRank, collator.RankValues(Zero, ShouldBeZero))
	ass.Equal(t, age.GreaterRank, collator.RankValues(One, ShouldBeZero))
	ass.Equal(t, age.EqualRank, collator.RankValues(ShouldBeZero, Zero))
	ass.Equal(t, age.LesserRank, collator.RankValues(Zero, One))
	ass.Equal(t, age.EqualRank, collator.RankValues(Zero, Zero))
	ass.Equal(t, age.GreaterRank, collator.RankValues(One, Zero))
	ass.Equal(t, age.EqualRank, collator.RankValues(One, One))

	// Integer
	var Zilch = 0
	var Two = 2
	var Three = 3
	var ShouldBeZilch int

	ass.Equal(t, age.EqualRank, collator.RankValues(ShouldBeZilch, ShouldBeZilch))
	ass.Equal(t, age.LesserRank, collator.RankValues(ShouldBeZilch, Two))
	ass.Equal(t, age.EqualRank, collator.RankValues(Zilch, ShouldBeZilch))
	ass.Equal(t, age.GreaterRank, collator.RankValues(Two, ShouldBeZilch))
	ass.Equal(t, age.EqualRank, collator.RankValues(ShouldBeZilch, Zilch))
	ass.Equal(t, age.LesserRank, collator.RankValues(Two, Three))
	ass.Equal(t, age.EqualRank, collator.RankValues(Two, Two))
	ass.Equal(t, age.GreaterRank, collator.RankValues(Three, Two))
	ass.Equal(t, age.EqualRank, collator.RankValues(Three, Three))

	// Float
	var Negligible = 0.0
	var Fourth = 0.25
	var Half = 0.5
	var ShouldBeNegligible float64

	ass.Equal(t, age.EqualRank, collator.RankValues(ShouldBeNegligible, ShouldBeNegligible))
	ass.Equal(t, age.LesserRank, collator.RankValues(ShouldBeNegligible, Half))
	ass.Equal(t, age.EqualRank, collator.RankValues(Negligible, ShouldBeNegligible))
	ass.Equal(t, age.GreaterRank, collator.RankValues(Half, ShouldBeNegligible))
	ass.Equal(t, age.EqualRank, collator.RankValues(ShouldBeNegligible, Negligible))
	ass.Equal(t, age.LesserRank, collator.RankValues(Fourth, Half))
	ass.Equal(t, age.EqualRank, collator.RankValues(Fourth, Fourth))
	ass.Equal(t, age.GreaterRank, collator.RankValues(Half, Fourth))
	ass.Equal(t, age.EqualRank, collator.RankValues(Half, Half))

	// Complex
	var Origin = 0 + 0i
	var PiOver4 = 1 + 1i
	var PiOver2 = 1 + 0i
	var ShouldBeOrigin complex128

	ass.Equal(t, age.EqualRank, collator.RankValues(ShouldBeOrigin, ShouldBeOrigin))
	ass.Equal(t, age.LesserRank, collator.RankValues(ShouldBeOrigin, PiOver4))
	ass.Equal(t, age.EqualRank, collator.RankValues(Origin, ShouldBeOrigin))
	ass.Equal(t, age.GreaterRank, collator.RankValues(PiOver4, ShouldBeOrigin))
	ass.Equal(t, age.EqualRank, collator.RankValues(ShouldBeOrigin, Origin))
	ass.Equal(t, age.LesserRank, collator.RankValues(PiOver2, PiOver4))
	ass.Equal(t, age.EqualRank, collator.RankValues(PiOver2, PiOver2))
	ass.Equal(t, age.GreaterRank, collator.RankValues(PiOver4, PiOver2))
	ass.Equal(t, age.EqualRank, collator.RankValues(PiOver4, PiOver4))

	// Rune
	var Null = rune(0)
	var Sad = '☹'
	var Happy = '☺'
	var ShouldBeNull rune

	ass.Equal(t, age.EqualRank, collator.RankValues(ShouldBeNull, ShouldBeNull))
	ass.Equal(t, age.LesserRank, collator.RankValues(ShouldBeNull, Sad))
	ass.Equal(t, age.EqualRank, collator.RankValues(Null, ShouldBeNull))
	ass.Equal(t, age.GreaterRank, collator.RankValues(Sad, ShouldBeNull))
	ass.Equal(t, age.EqualRank, collator.RankValues(ShouldBeNull, Null))
	ass.Equal(t, age.LesserRank, collator.RankValues(Sad, Happy))
	ass.Equal(t, age.EqualRank, collator.RankValues(Sad, Sad))
	ass.Equal(t, age.GreaterRank, collator.RankValues(Happy, Sad))
	ass.Equal(t, age.EqualRank, collator.RankValues(Happy, Happy))

	// String
	var Empty = ""
	var Hello = "Hello"
	var World = "World"
	var ShouldBeEmpty string

	ass.Equal(t, age.EqualRank, collator.RankValues(ShouldBeEmpty, ShouldBeEmpty))
	ass.Equal(t, age.LesserRank, collator.RankValues(ShouldBeEmpty, Hello))
	ass.Equal(t, age.EqualRank, collator.RankValues(Empty, ShouldBeEmpty))
	ass.Equal(t, age.GreaterRank, collator.RankValues(Hello, ShouldBeEmpty))
	ass.Equal(t, age.EqualRank, collator.RankValues(ShouldBeEmpty, Empty))
	ass.Equal(t, age.LesserRank, collator.RankValues(Hello, World))
	ass.Equal(t, age.EqualRank, collator.RankValues(Hello, Hello))
	ass.Equal(t, age.GreaterRank, collator.RankValues(World, Hello))
	ass.Equal(t, age.EqualRank, collator.RankValues(World, World))

	// Array
	var Universe = "Universe"
	var a0 = []any{}
	var a1 = []any{Hello, World}
	var a2 = []any{Hello, Universe}
	var a3 = []any{Hello, World, Universe}
	var a4 = []any{Hello, Universe, World}
	var aNil []any

	ass.Equal(t, age.EqualRank, collator.RankValues(aNil, aNil))
	ass.Equal(t, age.LesserRank, collator.RankValues(aNil, a0))
	ass.Equal(t, age.GreaterRank, collator.RankValues(a0, aNil))
	ass.Equal(t, age.EqualRank, collator.RankValues(a0, a0))
	ass.Equal(t, age.GreaterRank, collator.RankValues(a1, aNil))
	ass.Equal(t, age.LesserRank, collator.RankValues(a2, a1))
	ass.Equal(t, age.EqualRank, collator.RankValues(a2, a2))
	ass.Equal(t, age.GreaterRank, collator.RankValues(a1, a2))
	ass.Equal(t, age.EqualRank, collator.RankValues(a1, a1))
	ass.Equal(t, age.LesserRank, collator.RankValues(a2, a3))
	ass.Equal(t, age.EqualRank, collator.RankValues(a2, a2))
	ass.Equal(t, age.GreaterRank, collator.RankValues(a3, a2))
	ass.Equal(t, age.EqualRank, collator.RankValues(a3, a3))
	ass.Equal(t, age.LesserRank, collator.RankValues(a4, a1))
	ass.Equal(t, age.EqualRank, collator.RankValues(a4, a4))
	ass.Equal(t, age.GreaterRank, collator.RankValues(a1, a4))
	ass.Equal(t, age.EqualRank, collator.RankValues(a1, a1))

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

	ass.Equal(t, age.EqualRank, collator.RankValues(mNil, mNil))
	ass.Equal(t, age.LesserRank, collator.RankValues(mNil, m0))
	ass.Equal(t, age.GreaterRank, collator.RankValues(m0, mNil))
	ass.Equal(t, age.EqualRank, collator.RankValues(m0, m0))
	ass.Equal(t, age.LesserRank, collator.RankValues(m2, m1))
	ass.Equal(t, age.EqualRank, collator.RankValues(m2, m2))
	ass.Equal(t, age.GreaterRank, collator.RankValues(m1, m2))
	ass.Equal(t, age.EqualRank, collator.RankValues(m1, m1))
	ass.Equal(t, age.LesserRank, collator.RankValues(m2, m3))
	ass.Equal(t, age.EqualRank, collator.RankValues(m2, m2))
	ass.Equal(t, age.GreaterRank, collator.RankValues(m3, m2))
	ass.Equal(t, age.EqualRank, collator.RankValues(m3, m3))
	ass.Equal(t, age.LesserRank, collator.RankValues(m4, m1))
	ass.Equal(t, age.EqualRank, collator.RankValues(m4, m4))
	ass.Equal(t, age.GreaterRank, collator.RankValues(m1, m4))
	ass.Equal(t, age.EqualRank, collator.RankValues(m1, m1))

	// Struct
	var f1 = FooBar(1, "one", nil)
	var f2 = FooBar(1, "two", nil)
	var f3 = FooBar(2, "two", nil)
	var f4 = Fuz{"two"}
	var f5 = Fuz{"two"}
	var f6 = Fuz{"three"}
	ass.Equal(t, age.EqualRank, collator.RankValues(f1, f1))
	ass.Equal(t, age.LesserRank, collator.RankValues(f1, f2))
	ass.Equal(t, age.LesserRank, collator.RankValues(f2, f3))
	ass.Equal(t, age.GreaterRank, collator.RankValues(f3, f1))
	ass.Equal(t, age.GreaterRank, collator.RankValues(f3, f2))
	ass.Equal(t, age.EqualRank, collator.RankValues(f4, f4))
	ass.Equal(t, age.EqualRank, collator.RankValues(f4, f5))
	ass.Equal(t, age.GreaterRank, collator.RankValues(f5, f6))
	ass.Equal(t, age.GreaterRank, collator.RankValues(f3, &f4))
	ass.Equal(t, age.EqualRank, collator.RankValues(&f4, &f4))
	ass.Equal(t, age.EqualRank, collator.RankValues(&f4, &f5))
	ass.Equal(t, age.GreaterRank, collator.RankValues(&f5, &f6))
}

func TestTildeArrays(t *tes.T) {
	var collator = age.Collator[String]().Make()
	var ranker = collator.RankValues
	var sorter = age.Sorter[String]().MakeWithRanker(ranker)
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
	var notation = not.Notation().Make()
	var collator = age.Collator[any]().Make()
	var array = col.Array[any](notation).MakeFromArray([]any{0})
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
	var notation = not.Notation().Make()
	var collator = age.Collator[any]().Make()
	var m = col.Map[string, any](notation).MakeFromMap(map[string]any{"first": 1})
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

func TestIteratorsWithLists(t *tes.T) {
	var notation = not.Notation().Make()
	var array = col.Array[int](notation).MakeFromArray([]int{1, 2, 3, 4, 5})
	var list = col.List[int](notation).MakeFromSequence(array)
	var iterator = list.GetIterator()
	ass.False(t, iterator.HasPrevious())
	ass.True(t, iterator.HasNext())
	ass.Equal(t, 1, iterator.GetNext())
	ass.True(t, iterator.HasPrevious())
	ass.True(t, iterator.HasNext())
	ass.Equal(t, 1, iterator.GetPrevious())
	iterator.ToSlot(2)
	ass.True(t, iterator.HasPrevious())
	ass.True(t, iterator.HasNext())
	ass.Equal(t, 3, iterator.GetNext())
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
	var collator = age.Collator[any]().Make()
	var ranker = collator.RankValues
	var sorter = age.Sorter[any]().MakeWithRanker(ranker)
	var empty = []any{}
	sorter.SortValues(empty)
}

func TestSortingIntegers(t *tes.T) {
	var collator = age.Collator[int]().Make()
	var ranker = collator.RankValues
	var sorter = age.Sorter[int]().MakeWithRanker(ranker)
	var unsorted = []int{4, 3, 1, 5, 2}
	var sorted = []int{1, 2, 3, 4, 5}
	sorter.SortValues(unsorted)
	ass.Equal(t, sorted, unsorted)
}

func TestSortingStrings(t *tes.T) {
	var collator = age.Collator[string]().Make()
	var ranker = collator.RankValues
	var sorter = age.Sorter[string]().MakeWithRanker(ranker)
	var unsorted = []string{"alpha", "beta", "gamma", "delta"}
	var sorted = []string{"alpha", "beta", "delta", "gamma"}
	sorter.SortValues(unsorted)
	ass.Equal(t, sorted, unsorted)
}
