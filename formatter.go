/*******************************************************************************
 *   Copyright (c) 2009-2022 Crater Dog Technologies™.  All Rights Reserved.   *
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

// FORMATTER INTERFACE

// This function returns a string containing the canonical format for the
// specified value.
func FormatValue(value Value) string {
	var v = Formatter(0)
	v.FormatValue(value)
	return v.GetResult()
}

// This function returns a string containing the canonical format for the
// specified value indented the specified number of levels.
func FormatValueWithIndentation(value Value, indentation int) string {
	var v = Formatter(indentation)
	v.FormatValue(value)
	return v.GetResult()
}

// This constructor creates a new instance of a formatter that can be used to
// format any value using the specified number of levels of indentation.
func Formatter(indentation int) FormatterLike {
	return &formatter{indentation: indentation}
}

// FORMATTER IMPLEMENTATION

// The maximum depth level handles overly complex data models and cyclical
// references in a clean and efficient way.
const maximumDepth int = 8

// This type defines the structure and methods for a canonical formatter agent.
type formatter struct {
	indentation int
	depth       int
	result      sts.Builder
}

// This method returns the number of levels that each line is indented in the
// resulting canonical string.
func (v *formatter) GetIndentation() int {
	return v.indentation
}

// This method returns the canonical string for the specified value.
func (v *formatter) FormatValue(value Value) {
	v.formatValue(ref.ValueOf(value))
}

// This method returns the canonically formatted string result.
func (v *formatter) GetResult() string {
	result := v.result.String()
	v.result.Reset()
	return result
}

// This method appends the specified string to the result.
func (v *formatter) AppendString(s string) {
	v.result.WriteString(s)
}

// This method appends a properly indented newline to the result.
func (v *formatter) AppendNewline() {
	separator := "\n"
	levels := v.depth + v.indentation
	for level := 0; level < levels; level++ {
		separator += "\t"
	}
	v.result.WriteString(separator)
}

// This private method determines the actual type of the specified value and
// calls the corresponding format function for that type. Note that because the
// Go language doesn't really support polymorphism the selection of the actual
// function called must be done explicitly using reflection and a type switch.
func (v *formatter) formatValue(value ref.Value) {
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
			// The value is an attribute.
			v.formatAttribute(value)
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
				"Attempted to format:\n\tvalue: %v\n\ttype: %v\n\tkind: %v\n",
				value.Interface(),
				value.Type(),
				value.Kind()))
		}
	}
}

// This private method appends the nil string for the specified value to the
// result.
func (v *formatter) formatNil(r ref.Value) {
	v.AppendString("<nil>")
}

// This private method appends the name of the specified boolean value to the
// result.
func (v *formatter) formatBoolean(r ref.Value) {
	b := r.Bool()
	v.AppendString(stc.FormatBool(b))
}

// This private method appends the base 10 string for the specified integer
// value to the result.
func (v *formatter) formatInteger(r ref.Value) {
	i := r.Int()
	v.AppendString(stc.FormatInt(int64(i), 10))
}

// This private method appends the base 16 string for the specified unsigned
// integer value to the result.
func (v *formatter) formatUnsigned(r ref.Value) {
	u := r.Uint()
	v.AppendString("0x" + stc.FormatUint(uint64(u), 16))
}

// This private method appends the base 10 string for the specified floating
// point value to the result using scientific notation if necessary.
func (v *formatter) formatFloat(r ref.Value) {
	flt := r.Float()
	str := stc.FormatFloat(flt, 'G', -1, 64)
	if !sts.Contains(str, ".") && !sts.Contains(str, "E") {
		str += ".0"
	}
	v.AppendString(str)
}

// This private method appends the base 10 string for the specified complex
// number value to the result using scientific notation if necessary.
func (v *formatter) formatComplex(r ref.Value) {
	c := r.Complex()
	v.AppendString(stc.FormatComplex(c, 'G', -1, 128))
}

// This private method appends the string for the specified rune value to the
// result.
func (v *formatter) formatRune(r ref.Value) {
	rn := r.Int()
	v.AppendString(stc.QuoteRune(int32(rn)))
}

// This private method appends the string for the specified string value to the
// result.
func (v *formatter) formatString(r ref.Value) {
	str := r.String()
	v.AppendString(stc.Quote(str))
}

// This private method appends the string for the specified array of values to
// the result.
func (v *formatter) formatArray(r ref.Value, typ string) {
	size := r.Len()
	v.AppendString("[")
	if size > 0 {
		if v.depth+1 > maximumDepth {
			// Truncate the recursion.
			v.AppendString("...")
		} else {
			for i := 0; i < size; i++ {
				v.depth++
				v.AppendNewline()
				item := r.Index(i)
				v.formatValue(item)
				v.depth--
			}
			v.AppendNewline()
		}
	} else {
		if typ == "catalog" {
			v.AppendString(":") // The array of attributes is empty: [:]
		} else {
			v.AppendString(" ") // The array of items is empty: [ ]
		}
	}
	v.AppendString("](" + typ + ")")
}

// This private method appends the string for the specified map of key-value
// pairs to the result.
func (v *formatter) formatMap(r ref.Value) {
	keys := r.MapKeys()
	size := len(keys)
	v.AppendString("[")
	if size > 0 {
		if v.depth+1 > maximumDepth {
			// Truncate the recursion.
			v.AppendString("...")
		} else {
			for i := 0; i < size; i++ {
				v.depth++
				v.AppendNewline()
				key := keys[i]
				value := r.MapIndex(key)
				v.formatValue(key)
				v.AppendString(": ")
				v.formatValue(value)
				v.depth--
			}
			v.AppendNewline()
		}
	} else {
		v.AppendString(":") // The map is empty: [:]
	}
	v.AppendString("](map)")
}

// This private method appends the string for the specified catalog of
// key-value pairs to the result. It uses recursion to format each pair.
func (v *formatter) formatAttribute(r ref.Value) {
	key := r.MethodByName("GetKey").Call([]ref.Value{})[0]
	v.formatValue(key)
	v.AppendString(": ")
	value := r.MethodByName("GetValue").Call([]ref.Value{})[0]
	v.formatValue(value)
}

// This private method appends the string for the specified collection of
// items to the result. It uses recursion to format each item.
func (v *formatter) formatCollection(r ref.Value) {
	array := r.MethodByName("AsArray").Call([]ref.Value{})[0]
	typ := extractType(r)
	v.formatArray(array, typ)
}

// This private function extracts the type name string from the full reflected
// type.
func extractType(r ref.Value) string {
	t := r.Type().String()
	switch {
	case sts.HasPrefix(t, "[]"):
		return "array"
	case sts.HasPrefix(t, "map["):
		return "map"
	case sts.HasPrefix(t, "*collections.set"):
		return "set"
	case sts.HasPrefix(t, "collections.SetLike"):
		return "set"
	case sts.HasPrefix(t, "*collections.queue"):
		return "queue"
	case sts.HasPrefix(t, "collections.QueueLike"):
		return "queue"
	case sts.HasPrefix(t, "*collections.stack"):
		return "stack"
	case sts.HasPrefix(t, "collections.StackLike"):
		return "stack"
	case sts.HasPrefix(t, "*collections.list"):
		return "list"
	case sts.HasPrefix(t, "collections.ListLike"):
		return "list"
	case sts.HasPrefix(t, "*collections.catalog"):
		return "catalog"
	case sts.HasPrefix(t, "collections.CatalogLike"):
		return "catalog"
	default:
		fmt.Printf("UNKNOWN: %v\n", t)
		return "unknown"
	}
}
