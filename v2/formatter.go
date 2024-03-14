/*
................................................................................
.    Copyright (c) 2009-2024 Crater Dog Technologies™.  All Rights Reserved.   .
................................................................................
.  DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               .
.                                                                              .
.  This code is free software; you can redistribute it and/or modify it under  .
.  the terms of The MIT License (MIT), as published by the Open Source         .
.  Initiative. (See http://opensource.org/licenses/MIT)                        .
................................................................................
*/

package collections

import (
	fmt "fmt"
	ref "reflect"
	stc "strconv"
	sts "strings"
)

// FORMATTER INTERFACE

// This function returns a string containing the canonical format for the
// specified collection.
func FormatAssociation(association Value) string {
	var v = &formatter{}
	v.formatAssociation(ref.ValueOf(association))
	return v.getResult()
}

// This function returns a string containing the canonical format for the
// specified collection.
func FormatCollection(collection Collection) string {
	var v = &formatter{}
	v.formatValue(ref.ValueOf(collection))
	return v.getResult()
}

// This function returns the bytes containing the canonical format for the
// specified collection including the POSIX standard EOF marker.
func FormatDocument(collection Collection) []byte {
	var s = FormatCollection(collection) + EOF
	return []byte(s)
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

// This method returns the canonically formatted string result.
func (v *formatter) getResult() string {
	var result = v.result.String()
	v.result.Reset()
	return result
}

// This method appends the specified string to the result.
func (v *formatter) appendString(s string) {
	v.result.WriteString(s)
}

// This method appends a properly indented newline to the result.
func (v *formatter) appendNewline() {
	var separator = "\n"
	var levels = v.depth + v.indentation
	for level := 0; level < levels; level++ {
		separator += "    "
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

// This private method appends the nil string for the specified value to the
// result.
func (v *formatter) formatNil(r ref.Value) {
	v.appendString("<nil>")
}

// This private method appends the name of the specified boolean value to the
// result.
func (v *formatter) formatBoolean(r ref.Value) {
	var b = r.Bool()
	v.appendString(stc.FormatBool(b))
}

// This private method appends the base 10 string for the specified integer
// value to the result.
func (v *formatter) formatInteger(r ref.Value) {
	var i = r.Int()
	v.appendString(stc.FormatInt(int64(i), 10))
}

// This private method appends the base 16 string for the specified unsigned
// integer value to the result.
func (v *formatter) formatUnsigned(r ref.Value) {
	var u = r.Uint()
	v.appendString("0x" + stc.FormatUint(uint64(u), 16))
}

// This private method appends the base 10 string for the specified floating
// point value to the result using scientific notation if necessary.
func (v *formatter) formatFloat(r ref.Value) {
	var flt = r.Float()
	var str = stc.FormatFloat(flt, 'G', -1, 64)
	if !sts.Contains(str, ".") && !sts.Contains(str, "E") {
		str += ".0"
	}
	v.appendString(str)
}

// This private method appends the base 10 string for the specified complex
// number value to the result using scientific notation if necessary.
func (v *formatter) formatComplex(r ref.Value) {
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

// This private method appends the string for the specified rune value to the
// result.
func (v *formatter) formatRune(r ref.Value) {
	var rune_ = r.Int()
	v.appendString(stc.QuoteRune(int32(rune_)))
}

// This private method appends the string for the specified string value to the
// result.
func (v *formatter) formatString(r ref.Value) {
	var string_ = r.String()
	v.appendString(stc.Quote(string_))
}

// This private method appends the string for the specified array of values to
// the result.
func (v *formatter) formatArray(r ref.Value, typ string) {
	var size = r.Len()
	v.appendString("[")
	if size > 0 {
		if v.depth+1 > maximumDepth {
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

// This private method appends the string for the specified map of key-value
// pairs to the result.
func (v *formatter) formatMap(r ref.Value) {
	var keys = r.MapKeys()
	var size = len(keys)
	v.appendString("[")
	if size > 0 {
		if v.depth+1 > maximumDepth {
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

// This private method appends the string for the specified catalog of
// key-value pairs to the result. It uses recursion to format each pair.
func (v *formatter) formatAssociation(r ref.Value) {
	var key = r.MethodByName("GetKey").Call([]ref.Value{})[0]
	v.formatValue(key)
	v.appendString(": ")
	var value = r.MethodByName("GetValue").Call([]ref.Value{})[0]
	v.formatValue(value)
}

// This private method appends the string for the specified collection of
// values to the result. It uses recursion to format each value.
func (v *formatter) formatCollection(r ref.Value) {
	var array = r.MethodByName("AsArray").Call([]ref.Value{})[0]
	var typ = extractType(r)
	v.formatArray(array, typ)
}

// This private function extracts the type name string from the full reflected
// type.
func extractType(r ref.Value) string {
	var t = r.Type().String()
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
