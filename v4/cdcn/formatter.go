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

package cdcn

import (
	fmt "fmt"
	age "github.com/craterdog/go-collection-framework/v4/agent"
	col "github.com/craterdog/go-collection-framework/v4/collection"
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

func Formatter() FormatterClassLike {
	return formatterClass
}

// CLASS METHODS

// Target

type formatterClass_ struct {
	// Define the class constants.
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
		class_:   c,
		maximum_: maximum,
	}
}

// INSTANCE METHODS

// Target

type formatter_ struct {
	class_   FormatterClassLike
	depth_   int
	maximum_ int
	result_  sts.Builder
}

// Attributes

func (v *formatter_) GetClass() FormatterClassLike {
	return v.class_
}

func (v *formatter_) GetDepth() int {
	return v.depth_
}

func (v *formatter_) GetMaximum() int {
	return v.maximum_
}

// Public

func (v *formatter_) FormatValue(value any) (source string) {
	v.formatValue(value)
	v.appendNewline()
	source = v.getResult()
	return source
}

// Private

func (v *formatter_) appendNewline() {
	var separator = "\n"
	var indentation = "    "
	for level := 0; level < v.depth_; level++ {
		separator += indentation
	}
	v.appendString(separator)
}

func (v *formatter_) appendString(s string) {
	v.result_.WriteString(s)
}

func (v *formatter_) formatArray(array any) {
	var reflected = ref.ValueOf(array)
	var size = reflected.Len()
	switch {
	case v.depth_ == v.maximum_:
		// Truncate the recursion.
		v.appendString("...")
	case size == 0:
		// This is an empty sequence of values.
		v.appendString(" ")
	case size == 1:
		var value = reflected.Index(0).Interface()
		v.formatValue(value)
	default:
		// This is a multiline sequence of values.
		v.depth_++
		for i := 0; i < size; i++ {
			v.appendNewline()
			var value = reflected.Index(i).Interface()
			v.formatValue(value)
		}
		v.depth_--
		v.appendNewline()
	}
}

func (v *formatter_) formatAssociation(key any, value any) {
	v.formatIntrinsic(key)
	v.appendString(": ")
	v.formatValue(value)
}

func (v *formatter_) formatAssociations(associations any) {
	var sequence = ref.ValueOf(associations)
	var iterator = sequence.MethodByName("GetIterator").Call([]ref.Value{})[0]
	var size = sequence.MethodByName("GetSize").Call([]ref.Value{})[0].Interface()
	switch {
	case v.depth_ == v.maximum_:
		// Truncate the recursion.
		v.appendString("...")
	case size == 0:
		// This is an empty sequence of associations.
		v.appendString(":")
	case size == 1:
		var value = iterator.MethodByName("GetNext").Call([]ref.Value{})[0].Interface()
		v.formatValue(value)
	default:
		// This is a multiline sequence of associations.
		v.depth_++
		for iterator.MethodByName("HasNext").Call([]ref.Value{})[0].Interface().(bool) {
			v.appendNewline()
			var value = iterator.MethodByName("GetNext").Call([]ref.Value{})[0].Interface()
			v.formatValue(value)
		}
		v.depth_--
		v.appendNewline()
	}
}

func (v *formatter_) formatBoolean(boolean bool) {
	v.appendString(stc.FormatBool(boolean))
}

func (v *formatter_) formatCollection(collection any) {
	v.formatSequence(collection)
	v.formatContext(collection)
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

func (v *formatter_) formatContext(collection any) {
	var type_ string
	var reflectedString = ref.TypeOf(collection).String()
	var inspector = age.Inspector().Make()
	switch {
	// First try the implemented interface.  NOTE: These will not match any
	// generic types other than of type "any" as defined above.
	case inspector.ImplementsAspect(collection, (*col.CatalogLike[any, any])(nil)):
		type_ = "Catalog"
	case inspector.ImplementsAspect(collection, (*col.ListLike[any])(nil)):
		type_ = "List"
	case inspector.ImplementsAspect(collection, (*col.QueueLike[any])(nil)):
		type_ = "Queue"
	case inspector.ImplementsAspect(collection, (*col.SetLike[any])(nil)):
		type_ = "Set"
	case inspector.ImplementsAspect(collection, (*col.StackLike[any])(nil)):
		type_ = "Stack"

	// These two must come last since they implement a subset of the list
	// and catalog interfaces respectively.
	case inspector.ImplementsAspect(collection, (*col.ArrayLike[any])(nil)):
		type_ = "Array"
	case inspector.ImplementsAspect(collection, (*col.MapLike[any, any])(nil)):
		type_ = "Map"

	// Next try the actual value type.  This will catch generic types other
	// than "any" but is specific to these class implementations.
	case sts.HasPrefix(reflectedString, "collection.array_"):
		type_ = "Array"
	case sts.HasPrefix(reflectedString, "*collection.catalog_"):
		type_ = "Catalog"
	case sts.HasPrefix(reflectedString, "*collection.list_"):
		type_ = "List"
	case sts.HasPrefix(reflectedString, "collection.map_"):
		type_ = "Map"
	case sts.HasPrefix(reflectedString, "*collection.queue_"):
		type_ = "Queue"
	case sts.HasPrefix(reflectedString, "*collection.set_"):
		type_ = "Set"
	case sts.HasPrefix(reflectedString, "*collection.stack_"):
		type_ = "Stack"

	// And finally look for intrinsic arrays and maps.
	case sts.HasPrefix(reflectedString, "[]"):
		type_ = "Array"
	case sts.HasPrefix(reflectedString, "map["):
		type_ = "Map"
	default:
		type_ = reflectedString
	}
	v.appendString("(" + type_ + ")")
}

func (v *formatter_) formatFloat(float float64) {
	var str = stc.FormatFloat(float, 'G', -1, 64)
	if !sts.Contains(str, ".") && !sts.Contains(str, "E") {
		str += ".0"
	}
	v.appendString(str)
}

func (v *formatter_) formatInteger(integer int64) {
	v.appendString(stc.FormatInt(integer, 10))
}

func (v *formatter_) formatItems(items any) {
	var reflected = ref.ValueOf(items)
	var kind = reflected.Kind()
	switch kind {
	case ref.Map:
		v.formatMap(items)
	case ref.Array, ref.Slice:
		v.formatArray(items)
	case ref.Interface, ref.Pointer:
		switch {
		case reflected.MethodByName("GetKeys").IsValid():
			v.formatAssociations(items)
		default:
			v.formatValues(items)
		}
	default:
		var message = fmt.Sprintf(
			"Attempted to format:\n    value: %v\n    type: %v\n    kind: %v\n",
			reflected.Interface(),
			reflected.Type(),
			reflected.Kind(),
		)
		panic(message)
	}
}

func (v *formatter_) formatMap(map_ any) {
	var reflected = ref.ValueOf(map_)
	var size = reflected.Len()
	var keys = reflected.MapKeys()
	switch {
	case v.depth_ == v.maximum_:
		// Truncate the recursion.
		v.appendString("...")
	case size == 0:
		// This is an empty map of associations.
		v.appendString(":")
	case size == 1:
		var key = keys[0].Interface()
		var value = reflected.MapIndex(keys[0]).Interface()
		v.formatAssociation(key, value)
	default:
		// This is a multiline map of associations.
		v.depth_++
		for i := 0; i < size; i++ {
			v.appendNewline()
			var key = keys[i].Interface()
			var value = reflected.MapIndex(keys[i]).Interface()
			v.formatAssociation(key, value)
		}
		v.depth_--
		v.appendNewline()
	}
}

func (v *formatter_) formatNil(value any) {
	v.appendString("nil")
}

func (v *formatter_) formatIntrinsic(intrinsic any) {
	switch actual := intrinsic.(type) {
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
	default:
		var message = fmt.Sprintf(
			"Attempted to format unknown intrinsic: %T",
			intrinsic,
		)
		panic(message)
	}
}

func (v *formatter_) formatRune(rune_ rune) {
	v.appendString(stc.QuoteRune(rune_))
}

func (v *formatter_) formatSequence(sequence any) {
	v.appendString("[")
	v.formatItems(sequence)
	v.appendString("]")
}

func (v *formatter_) formatString(string_ string) {
	v.appendString(stc.Quote(string_))
}

func (v *formatter_) formatUnsigned(unsigned uint64) {
	v.appendString("0x" + stc.FormatUint(unsigned, 16))
}

func (v *formatter_) formatValue(value any) {
	var reflected = ref.ValueOf(value)
	var kind = reflected.Kind()
	switch kind {
	case ref.Array, ref.Slice, ref.Map:
		v.formatCollection(value)
	case ref.Interface, ref.Pointer:
		if reflected.MethodByName("GetKey").IsValid() {
			var key = reflected.MethodByName("GetKey").Call([]ref.Value{})[0].Interface()
			var value = reflected.MethodByName("GetValue").Call([]ref.Value{})[0].Interface()
			v.formatAssociation(key, value)
		} else {
			v.formatCollection(value)
		}
	default:
		v.formatIntrinsic(value)
	}
}

func (v *formatter_) formatValues(values any) {
	var sequence = ref.ValueOf(values)
	var iterator = sequence.MethodByName("GetIterator").Call([]ref.Value{})[0]
	var size = sequence.MethodByName("GetSize").Call([]ref.Value{})[0].Interface()
	switch {
	case v.depth_ == v.maximum_:
		// Truncate the recursion.
		v.appendString("...")
	case size == 0:
		// This is an empty sequence of values.
		v.appendString(" ")
	case size == 1:
		var value = iterator.MethodByName("GetNext").Call([]ref.Value{})[0].Interface()
		v.formatValue(value)
	default:
		// This is a multiline sequence of values.
		v.depth_++
		for iterator.MethodByName("HasNext").Call([]ref.Value{})[0].Interface().(bool) {
			v.appendNewline()
			var value = iterator.MethodByName("GetNext").Call([]ref.Value{})[0].Interface()
			v.formatValue(value)
		}
		v.depth_--
		v.appendNewline()
	}
}

func (v *formatter_) getResult() string {
	var result = v.result_.String()
	v.result_.Reset()
	return result
}
