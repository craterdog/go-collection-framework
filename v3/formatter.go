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

// CLASS ACCESS

// Reference

var formatterClass = &formatterClass_{
	defaultMaximum_: 8,
}

// Function

/*
Formatter defines an implementation of a formatter-like class that uses Crater
Dog Collection Notation™ (CDCN) for formatting collections.  This is required by
the Go `Stringer` interface when the `String()` method is called on a
collection.  If not for the requirement to support the Go `Stringer` interface
this class would be located in the `cdcn` package with the rest of the CDCN
classes.  Instead, the `cdcn.FormatterClass` must delegate its implementation to
this class to avoid circular dependencies.
*/
func Formatter() FormatterClassLike {
	return formatterClass
}

// CLASS METHODS

// Target

type formatterClass_ struct {
	defaultMaximum_ int
}

// Constants

func (c *formatterClass_) DefaultMaximum() int {
	return c.defaultMaximum_
}

// Constructors

func (c *formatterClass_) Make() FormatterLike {
	return &formatter_{
		maximum_: c.defaultMaximum_,
	}
}

func (c *formatterClass_) MakeWithMaximum(maximum int) FormatterLike {
	if maximum < 0 {
		maximum = c.defaultMaximum_
	}
	return &formatter_{
		maximum_: maximum,
	}
}

// INSTANCE METHODS

// Target

type formatter_ struct {
	depth_   int
	maximum_ int
	result_  sts.Builder
}

// Attributes

func (v *formatter_) GetDepth() int {
	return v.depth_
}

func (v *formatter_) GetMaximum() int {
	return v.maximum_
}

// Public

func (v *formatter_) FormatCollection(collection Collection) string {
	var reflected = ref.ValueOf(collection)
	v.formatCollection(reflected)
	return v.getResult()
}

// Private

func (v *formatter_) appendString(s string) {
	v.result_.WriteString(s)
}

func (v *formatter_) appendNewline() {
	var separator = "\n"
	for level := 0; level < v.depth_; level++ {
		separator += "    "
	}
	v.result_.WriteString(separator)
}

/*
NOTE:
Because the Go language doesn't handle generic types very well in type switches,
we use reflection to handle all generic types.
*/
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

func (v *formatter_) formatArray(array ref.Value) {
	var size = array.Len()
	switch {
	case v.depth_ == v.maximum_:
		// Truncate the recursion.
		v.appendString("...")
	case size == 0:
		v.appendString(" ")
	case size == 1:
		var value = array.Index(0)
		v.formatValue(value.Interface())
	default:
		v.depth_++
		for i := 0; i < size; i++ {
			v.appendNewline()
			var value = array.Index(i)
			v.formatValue(value.Interface())
		}
		v.depth_--
		v.appendNewline()
	}
}

func (v *formatter_) formatMap(map_ ref.Value) {
	var keys = map_.MapKeys()
	var size = len(keys)
	switch {
	case v.depth_ == v.maximum_:
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
		v.depth_++
		for i := 0; i < size; i++ {
			v.appendNewline()
			var key = keys[i]
			var value = map_.MapIndex(key)
			v.formatValue(key.Interface())
			v.appendString(": ")
			v.formatValue(value.Interface())
		}
		v.depth_--
		v.appendNewline()
	}
}

func (v *formatter_) formatNil(value any) {
	v.appendString("<nil>")
}

func (v *formatter_) formatBoolean(boolean bool) {
	v.appendString(stc.FormatBool(boolean))
}

func (v *formatter_) formatInteger(integer int64) {
	v.appendString(stc.FormatInt(integer, 10))
}

func (v *formatter_) formatUnsigned(unsigned uint64) {
	v.appendString("0x" + stc.FormatUint(unsigned, 16))
}

func (v *formatter_) formatFloat(float float64) {
	var str = stc.FormatFloat(float, 'G', -1, 64)
	if !sts.Contains(str, ".") && !sts.Contains(str, "E") {
		str += ".0"
	}
	v.appendString(str)
}

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

func (v *formatter_) formatRune(rune_ rune) {
	v.appendString(stc.QuoteRune(rune_))
}

func (v *formatter_) formatString(string_ string) {
	v.appendString(stc.Quote(string_))
}

func (v *formatter_) formatAssociation(association ref.Value) {
	var key = association.MethodByName("GetKey").Call([]ref.Value{})[0]
	v.formatValue(key.Interface())
	v.appendString(": ")
	var value = association.MethodByName("GetValue").Call([]ref.Value{})[0]
	v.formatValue(value.Interface())
}

func (v *formatter_) formatSequence(sequence ref.Value) {
	var iterator = sequence.MethodByName("GetIterator").Call([]ref.Value{})[0]
	var size = sequence.MethodByName("GetSize").Call([]ref.Value{})[0].Interface()
	switch {
	case v.depth_ == v.maximum_:
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
		v.depth_++
		for iterator.MethodByName("HasNext").Call([]ref.Value{})[0].Interface().(bool) {
			v.appendNewline()
			var value = iterator.MethodByName("GetNext").Call([]ref.Value{})[0]
			v.formatValue(value.Interface())
		}
		v.depth_--
		v.appendNewline()
	}
}

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

func (v *formatter_) getResult() string {
	var result = v.result_.String()
	v.result_.Reset()
	return result
}

/*
NOTE:
This hack is necessary since Go does not handle type switches with generics very
well.
*/
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
