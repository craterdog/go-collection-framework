/*******************************************************************************
 *   Copyright (c) 2009-2024 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package collections

import (
	fmt "fmt"
	ref "reflect"
	stc "strconv"
	sts "strings"
)

// CLASS NAMESPACE

// Private Class Namespace Type

type formatterClass_ struct {
	defaultDepth int
}

// Private Class Namespace Reference

var formatterClass = &formatterClass_{
	defaultDepth: 8,
}

// Public Class Namespace Access

/*
FormatterClass defines an implementation of a formatter-like class that uses
Crater Dog Collection Notation™ (CDCN) for formatting collections.  This is
required by the Go `Stringer` interface when the `String()` method is called on
a collection.  If not for the requirement to support the Go `Stringer` interface
this class would be located in the `cdcn` package with the rest of the CDCN
classes.  Instead, the `cdcn.FormatterClass` must delegate its implementation to
this class to avoid circular dependencies.
*/
func FormatterClass() FormatterClassLike {
	return formatterClass
}

// Public Class Constants

func (c *formatterClass_) DefaultDepth() int {
	return c.defaultDepth
}

// Public Class Constructors

func (c *formatterClass_) Make() FormatterLike {
	return c.MakeWithDepth(c.defaultDepth)
}

func (c *formatterClass_) MakeWithDepth(depth int) FormatterLike {
	var formatter = &formatter_{
		maximum: depth,
	}
	return formatter
}

// CLASS INSTANCES

// Private Class Type Definition

type formatter_ struct {
	depth   int
	maximum int
	result  sts.Builder
}

// Public Interface

func (v *formatter_) FormatCollection(collection Collection) string {
	var reflected = ref.ValueOf(collection)
	v.formatCollection(reflected)
	return v.getResult()
}

// Private Interface

// This private class method appends the specified string to the result.
func (v *formatter_) appendString(s string) {
	v.result.WriteString(s)
}

// This private class method appends a properly indented newline to the result.
func (v *formatter_) appendNewline() {
	var separator = "\n"
	for level := 0; level < v.depth; level++ {
		separator += "    "
	}
	v.result.WriteString(separator)
}

// This private class method determines the actual type of the specified value
// and calls the corresponding format function for that type.  NOTE: Because the
// Go language doesn't handle generic types very well in type switches, we use
// reflection to handle all generic types.
func (v *formatter_) formatValue(value any) {
	switch actual := value.(type) {
	// Handle primitive types.
	case nil:
		v.formatNil(actual)
	case bool:
		v.formatBoolean(actual)
	case uint:
		v.formatUnsigned(uint64(actual))
	case uint8:
		v.formatUnsigned(uint64(actual))
	case uint16:
		v.formatUnsigned(uint64(actual))
	case uint32:
		v.formatUnsigned(uint64(actual))
	case uint64:
		v.formatUnsigned(uint64(actual))
	case uintptr:
		v.formatUnsigned(uint64(actual))
	case int:
		v.formatInteger(int64(actual))
	case int8:
		v.formatInteger(int64(actual))
	case int16:
		v.formatInteger(int64(actual))
	case int64:
		v.formatInteger(int64(actual))
	case float32:
		v.formatFloat(float64(actual))
	case float64:
		v.formatFloat(float64(actual))
	case complex64:
		v.formatComplex(complex128(actual))
	case complex128:
		v.formatComplex(complex128(actual))
	case rune:
		v.formatRune(rune(actual))
	case string:
		v.formatString(actual)

	// Handle generic types.
	default:
		var reflected = ref.ValueOf(value)
		if reflected.MethodByName("GetKey").IsValid() {
			v.formatAssociation(reflected)
		} else if reflected.MethodByName("AsArray").IsValid() {
			v.formatCollection(reflected)
		} else {
			switch reflected.Kind() {
			case ref.Array, ref.Slice, ref.Map:
				v.formatCollection(reflected)
			default:
				panic(fmt.Sprintf(
					"Attempted to format:\n    value: %v\n    type: %v\n    kind: %v\n",
					reflected.Interface(),
					reflected.Type(),
					reflected.Kind(),
				))
			}
		}
	}
}

// This private class method adds the canonical format for the specified Go
// array of values to the state of the formatter.
func (v *formatter_) formatArray(array ref.Value) {
	var size = array.Len()
	switch {
	case v.depth == v.maximum:
		// Truncate the recursion.
		v.appendString("...")
	case size == 0:
		v.appendString(" ")
	case size == 1:
		var value = array.Index(0)
		v.formatValue(value.Interface())
	default:
		v.depth++
		for i := 0; i < size; i++ {
			v.appendNewline()
			var value = array.Index(i)
			v.formatValue(value.Interface())
		}
		v.depth--
		v.appendNewline()
	}
}

// This private class method adds the canonical format for the specified Go map
// of key-value pairs to the state of the formatter.
func (v *formatter_) formatMap(map_ ref.Value) {
	var keys = map_.MapKeys()
	var size = len(keys)
	switch {
	case v.depth == v.maximum:
		// Truncate the recursion.
		v.appendString("...")
	case size == 0:
		v.appendString(":")
	case size == 1:
		var key = keys[0]
		var value = map_.MapIndex(key)
		v.formatValue(key.Interface())
		v.appendString(": ")
		v.formatValue(value.Interface())
	default:
		v.depth++
		for i := 0; i < size; i++ {
			v.appendNewline()
			var key = keys[i]
			var value = map_.MapIndex(key)
			v.formatValue(key.Interface())
			v.appendString(": ")
			v.formatValue(value.Interface())
		}
		v.depth--
		v.appendNewline()
	}
}

// This private class method appends the nil string for the specified value to
// the result.
func (v *formatter_) formatNil(value any) {
	v.appendString("<nil>")
}

// This private class method appends the name of the specified boolean value to
// the result.
func (v *formatter_) formatBoolean(boolean bool) {
	v.appendString(stc.FormatBool(boolean))
}

// This private class method appends the base 10 string for the specified
// integer value to the result.
func (v *formatter_) formatInteger(integer int64) {
	v.appendString(stc.FormatInt(integer, 10))
}

// This private class method appends the base 16 string for the specified
// unsigned integer value to the result.
func (v *formatter_) formatUnsigned(unsigned uint64) {
	v.appendString("0x" + stc.FormatUint(unsigned, 16))
}

// This private class method appends the base 10 string for the specified
// floating point value to the result using scientific notation if necessary.
func (v *formatter_) formatFloat(float float64) {
	var str = stc.FormatFloat(float, 'G', -1, 64)
	if !sts.Contains(str, ".") && !sts.Contains(str, "E") {
		str += ".0"
	}
	v.appendString(str)
}

// This private class method appends the base 10 string for the specified
// complex number value to the result using scientific notation if necessary.
func (v *formatter_) formatComplex(complex_ complex128) {
	var real_ = real(complex_)
	var imag_ = imag(complex_)
	v.appendString("(")
	v.formatFloat(real_)
	if imag_ >= 0.0 {
		v.appendString("+")
	}
	v.formatFloat(imag_)
	v.appendString("i)")
}

// This private class method appends the string for the specified rune value to
// the result.
func (v *formatter_) formatRune(rune_ rune) {
	v.appendString(stc.QuoteRune(rune_))
}

// This private class method appends the string for the specified string value
// to the result.
func (v *formatter_) formatString(string_ string) {
	v.appendString(stc.Quote(string_))
}

// This private class method appends the string for the specified association to
// the result.
func (v *formatter_) formatAssociation(association ref.Value) {
	var key = association.MethodByName("GetKey").Call([]ref.Value{})[0]
	v.formatValue(key.Interface())
	v.appendString(": ")
	var value = association.MethodByName("GetValue").Call([]ref.Value{})[0]
	v.formatValue(value.Interface())
}

// This private class method appends the string for the specified sequence of
// associations to the result.
func (v *formatter_) formatSequence(sequence ref.Value) {
	var iterator = sequence.MethodByName("GetIterator").Call([]ref.Value{})[0]
	var size = sequence.MethodByName("GetSize").Call([]ref.Value{})[0].Interface()
	switch {
	case v.depth == v.maximum:
		// Truncate the recursion.
		v.appendString("...")
	case size == 0:
		if sequence.MethodByName("GetKeys").IsValid() {
			v.appendString(":") // This is an empty sequence of associations.
		} else {
			v.appendString(" ") // This is an empty sequence of values.
		}
	case size == 1:
		var value = iterator.MethodByName("GetNext").Call([]ref.Value{})[0]
		v.formatValue(value.Interface())
	default:
		v.depth++
		for iterator.MethodByName("HasNext").Call([]ref.Value{})[0].Interface().(bool) {
			v.appendNewline()
			var value = iterator.MethodByName("GetNext").Call([]ref.Value{})[0]
			v.formatValue(value.Interface())
		}
		v.depth--
		v.appendNewline()
	}
}

// This private class method appends the string for the specified collection of
// values to the result. It uses recursion to format each value.
func (v *formatter_) formatCollection(collection ref.Value) {
	v.appendString("[")
	var type_ = v.getName(collection)
	switch collection.Kind() {
	case ref.Array, ref.Slice:
		v.formatArray(collection)
	case ref.Map:
		v.formatMap(collection)
	case ref.Interface, ref.Pointer:
		v.formatSequence(collection)
	default:
		panic(fmt.Sprintf(
			"Attempted to format:\n    value: %v\n    type: %v\n    kind: %v\n",
			collection.Interface(),
			collection.Type(),
			collection.Kind(),
		))
	}
	v.appendString("](" + type_ + ")")
}

// This private class method returns the canonically formatted string result.
func (v *formatter_) getResult() string {
	var result = v.result.String()
	v.result.Reset()
	return result
}

// This private class method extracts the type name string from the full
// reflected type.  NOTE: This hack is necessary since Go does not handle type
// switches with generics very well.
func (v *formatter_) getName(collection ref.Value) string {
	var type_ = collection.Type().String()
	switch {
	case sts.HasPrefix(type_, "[]"):
		return "array"
	case sts.HasPrefix(type_, "collections.ArrayLike"):
		return "Array"
	case sts.HasPrefix(type_, "collections.array_"):
		return "Array"
	case sts.HasPrefix(type_, "map["):
		return "map"
	case sts.HasPrefix(type_, "collections.MapLike"):
		return "Map"
	case sts.HasPrefix(type_, "collections.map_"):
		return "Map"
	case sts.HasPrefix(type_, "*collections.set_"):
		return "Set"
	case sts.HasPrefix(type_, "collections.SetLike"):
		return "Set"
	case sts.HasPrefix(type_, "*collections.queue_"):
		return "Queue"
	case sts.HasPrefix(type_, "collections.QueueLike"):
		return "Queue"
	case sts.HasPrefix(type_, "*collections.stack_"):
		return "Stack"
	case sts.HasPrefix(type_, "collections.StackLike"):
		return "Stack"
	case sts.HasPrefix(type_, "*collections.list_"):
		return "List"
	case sts.HasPrefix(type_, "collections.ListLike"):
		return "List"
	case sts.HasPrefix(type_, "*collections.catalog_"):
		return "Catalog"
	case sts.HasPrefix(type_, "collections.CatalogLike"):
		return "Catalog"
	case sts.HasPrefix(type_, "*collections.association_"):
		return "Association"
	case sts.HasPrefix(type_, "collections.AssociationLike"):
		return "Association"
	default:
		return "<unknown>"
	}
}
