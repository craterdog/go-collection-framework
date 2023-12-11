/*******************************************************************************
 *   Copyright (c) 2009-2023 Crater Dog Technologies™.  All Rights Reserved.   *
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

// This private type defines the namespace structure associated with the
// constants, constructors and functions for the formatter class namespace.
type formatterClass_ struct {
	defaultDepth int
}

// This private constant defines the singleton reference to the formatter
// class namespace.  It also initializes any class constants as needed.
var formatterClassSingleton = &formatterClass_{
	defaultDepth: 8,
}

// This public function returns the singleton reference to the formatter
// class namespace.
func Formatter() *formatterClass_ {
	return formatterClassSingleton
}

// CLASS CONSTANTS

// This public class constant represents the default depth of a collection
// at which the formatter gives up and inserts "...".  This handles cycles
// in a sensible and efficient manner.
func (c *formatterClass_) GetDefaultDepth() int {
	return c.defaultDepth
}

// CLASS CONSTRUCTORS

// This public class constructor creates a new formatter with the default
// maximum traversal depth.
func (c *formatterClass_) WithDefaultDepth() FormatterLike {
	var formatter = &formatter_{
		maximumDepth: c.defaultDepth,
	}
	return formatter
}

// This public class constructor creates a new formatter with the specified
// maximum traversal depth.
func (c *formatterClass_) WithMaximumDepth(maximumDepth int) FormatterLike {
	if maximumDepth < 1 || maximumDepth > c.defaultDepth {
		maximumDepth = c.defaultDepth
	}
	var formatter = &formatter_{
		maximumDepth: maximumDepth,
	}
	return formatter
}

// CLASS FUNCTIONS

// This public class function returns a string containing the canonical format
// for the specified association.
func (c *formatterClass_) FormatAssociation(association Value) string {
	var formatter = c.WithDefaultDepth()
	var string_ = formatter.FormatAssociation(association)
	return string_
}

// This public class function returns a string containing the canonical format
// for the specified collection.
func (c *formatterClass_) FormatCollection(collection Collection) string {
	var formatter = c.WithDefaultDepth()
	var string_ = formatter.FormatCollection(collection)
	return string_
}

// This public class function returns the bytes containing the canonical format
// for the specified collection including the POSIX standard EOF marker.
func (c *formatterClass_) FormatDocument(collection Collection) []byte {
	var formatter = c.WithDefaultDepth()
	var string_ = formatter.FormatCollection(collection) + EOF
	return []byte(string_)
}

// CLASS TYPE

// Encapsulated Type

// This private class type encapsulates a Go structure containing private
// attributes that can only be accessed and manipulated using methods that
// implement the formatter-like abstract type.
type formatter_ struct {
	depth        int
	maximumDepth int
	indentation  int
	result       sts.Builder
}

// Canonical Interface

// This public class method returns a string containing the canonical format
// for the specified association.
func (v *formatter_) FormatAssociation(association Value) string {
	v.formatAssociation(ref.ValueOf(association))
	return v.getResult()
}

// This public class method returns a string containing the canonical format
// for the specified collection.
func (v *formatter_) FormatCollection(collection Collection) string {
	v.formatValue(ref.ValueOf(collection))
	return v.getResult()
}

// Private Interface

// This private class method returns the canonically formatted string result.
func (v *formatter_) getResult() string {
	var result = v.result.String()
	v.result.Reset()
	return result
}

// This private class method appends the specified string to the result.
func (v *formatter_) appendString(s string) {
	v.result.WriteString(s)
}

// This private class method appends a properly indented newline to the result.
func (v *formatter_) appendNewline() {
	var separator = "\n"
	var levels = v.depth + v.indentation
	for level := 0; level < levels; level++ {
		separator += "    "
	}
	v.result.WriteString(separator)
}

// This private class method determines the actual type of the specified value
// and calls the corresponding format function for that type. NOTE: Because the
// Go language doesn't really support polymorphism, the selection of the actual
// function called must be done explicitly using reflection and a type switch.
func (v *formatter_) formatValue(value ref.Value) {
	switch value.Kind() {

	// Handle all primitive types.
	case ref.Bool:
		v.formatBoolean(value)
	case ref.Uint, ref.Uint8, ref.Uint16, ref.Uint32, ref.Uint64, ref.Uintptr:
		v.formatUnsigned(value)
	case ref.Int, ref.Int8, ref.Int16, ref.Int64:
		v.formatInteger(value)
	case ref.Float32, ref.Float64:
		v.formatFloat(value)
	case ref.Complex64, ref.Complex128:
		v.formatComplex(value)
	case ref.Int32:
		v.formatRune(value)
	case ref.String:
		v.formatString(value)

	// Handle all primitive collections.
	case ref.Array, ref.Slice:
		v.formatArray(value, "array")
	case ref.Map:
		v.formatMap(value)

	// Handle all sequential collections.
	case ref.Interface, ref.Pointer:
		if value.MethodByName("AsArray").IsValid() {
			// The value is a collection.
			v.formatCollection(value)
		} else if value.MethodByName("GetKey").IsValid() {
			// The value is an association.
			v.formatAssociation(value)
		} else {
			// The value is a pointer to the value to be formatted.
			value = value.Elem()
			v.formatValue(value)
		}

	default:
		if !value.IsValid() {
			v.formatNil(value)
		} else {
			panic(fmt.Sprintf(
				"Attempted to format:\n    value: %v\n    type: %v\n    kind: %v\n",
				value.Interface(),
				value.Type(),
				value.Kind()))
		}
	}
}

// This private class method appends the nil string for the specified value to
// the result.
func (v *formatter_) formatNil(r ref.Value) {
	v.appendString("<nil>")
}

// This private class method appends the name of the specified boolean value to
// the result.
func (v *formatter_) formatBoolean(r ref.Value) {
	var b = r.Bool()
	v.appendString(stc.FormatBool(b))
}

// This private class method appends the base 10 string for the specified
// integer value to the result.
func (v *formatter_) formatInteger(r ref.Value) {
	var i = r.Int()
	v.appendString(stc.FormatInt(int64(i), 10))
}

// This private class method appends the base 16 string for the specified
// unsigned integer value to the result.
func (v *formatter_) formatUnsigned(r ref.Value) {
	var u = r.Uint()
	v.appendString("0x" + stc.FormatUint(uint64(u), 16))
}

// This private class method appends the base 10 string for the specified
// floating point value to the result using scientific notation if necessary.
func (v *formatter_) formatFloat(r ref.Value) {
	var flt = r.Float()
	var str = stc.FormatFloat(flt, 'G', -1, 64)
	if !sts.Contains(str, ".") && !sts.Contains(str, "E") {
		str += ".0"
	}
	v.appendString(str)
}

// This private class method appends the base 10 string for the specified
// complex number value to the result using scientific notation if necessary.
func (v *formatter_) formatComplex(r ref.Value) {
	var complex_ = r.Complex()
	var real_ = ref.ValueOf(real(complex_))
	var imag_ = ref.ValueOf(imag(complex_))
	v.appendString("(")
	v.formatFloat(real_)
	if imag_.Float() >= 0.0 {
		v.appendString("+")
	}
	v.formatFloat(imag_)
	v.appendString("i)")
}

// This private class method appends the string for the specified rune value to
// the result.
func (v *formatter_) formatRune(r ref.Value) {
	var rune_ = r.Int()
	v.appendString(stc.QuoteRune(int32(rune_)))
}

// This private class method appends the string for the specified string value
// to the result.
func (v *formatter_) formatString(r ref.Value) {
	var string_ = r.String()
	v.appendString(stc.Quote(string_))
}

// This private class method appends the string for the specified array of
// values to the result.
func (v *formatter_) formatArray(r ref.Value, typ string) {
	var size = r.Len()
	v.appendString("[")
	if size > 0 {
		if v.depth+1 > v.maximumDepth {
			// Truncate the recursion.
			v.appendString("...")
		} else {
			for i := 0; i < size; i++ {
				var value ref.Value
				v.depth++
				v.appendNewline()
				if typ == "stack" {
					// Format values in reverse order.
					value = r.Index(size - i - 1)
				} else {
					// Format values in actual order.
					value = r.Index(i)
				}
				v.formatValue(value)
				v.depth--
			}
			v.appendNewline()
		}
	} else {
		if typ == "catalog" {
			v.appendString(":") // The array of associations is empty: [:]
		} else {
			v.appendString(" ") // The array of values is empty: [ ]
		}
	}
	v.appendString("](" + typ + ")")
}

// This private class method appends the string for the specified map of
// key-value pairs to the result.
func (v *formatter_) formatMap(r ref.Value) {
	var keys = r.MapKeys()
	var size = len(keys)
	v.appendString("[")
	if size > 0 {
		if v.depth+1 > v.maximumDepth {
			// Truncate the recursion.
			v.appendString("...")
		} else {
			for i := 0; i < size; i++ {
				v.depth++
				v.appendNewline()
				var key = keys[i]
				var value = r.MapIndex(key)
				v.formatValue(key)
				v.appendString(": ")
				v.formatValue(value)
				v.depth--
			}
			v.appendNewline()
		}
	} else {
		v.appendString(":") // The map is empty: [:]
	}
	v.appendString("](map)")
}

// This private class method appends the string for the specified catalog of
// key-value pairs to the result. It uses recursion to format each pair.
func (v *formatter_) formatAssociation(r ref.Value) {
	var key = r.MethodByName("GetKey").Call([]ref.Value{})[0]
	v.formatValue(key)
	v.appendString(": ")
	var value = r.MethodByName("GetValue").Call([]ref.Value{})[0]
	v.formatValue(value)
}

// This private class method appends the string for the specified collection of
// values to the result. It uses recursion to format each value.
func (v *formatter_) formatCollection(r ref.Value) {
	var array = r.MethodByName("AsArray").Call([]ref.Value{})[0]
	var type_ = r.Type().String()
	switch {
	case sts.HasPrefix(type_, "[]"):
		type_ = "array"
	case sts.HasPrefix(type_, "map["):
		type_ = "map"
	case sts.HasPrefix(type_, "*collections.set"):
		type_ = "set"
	case sts.HasPrefix(type_, "collections.SetLike"):
		type_ = "set"
	case sts.HasPrefix(type_, "*collections.queue"):
		type_ = "queue"
	case sts.HasPrefix(type_, "collections.QueueLike"):
		type_ = "queue"
	case sts.HasPrefix(type_, "*collections.stack"):
		type_ = "stack"
	case sts.HasPrefix(type_, "collections.StackLike"):
		type_ = "stack"
	case sts.HasPrefix(type_, "*collections.list"):
		type_ = "list"
	case sts.HasPrefix(type_, "collections.ListLike"):
		type_ = "list"
	case sts.HasPrefix(type_, "*collections.catalog"):
		type_ = "catalog"
	case sts.HasPrefix(type_, "collections.CatalogLike"):
		type_ = "catalog"
	default:
		fmt.Printf("UNKNOWN: %v\n", type_)
		type_ = "unknown"
	}
	v.formatArray(array, type_)
}
