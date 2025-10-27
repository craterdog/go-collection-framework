/*
................................................................................
.    Copyright (c) 2009-2025 Crater Dog Technologies.  All Rights Reserved.    .
................................................................................
.  DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               .
.                                                                              .
.  This code is free software; you can redistribute it and/or modify it under  .
.  the terms of The MIT License (MIT), as published by the Open Source         .
.  Initiative. (See https://opensource.org/license/MIT)                        .
................................................................................
*/

package ranges

import (
	fmt "fmt"
	age "github.com/craterdog/go-collection-framework/v8/agents"
	col "github.com/craterdog/go-collection-framework/v8/collections"
	uti "github.com/craterdog/go-missing-utilities/v8"
	ref "reflect"
	sts "strings"
	syn "sync"
)

// CLASS INTERFACE

// Access Function

func IntervalClass[V Discrete]() IntervalClassLike[V] {
	return intervalClass[V]()
}

// Constructor Methods

func (c *intervalClass_[V]) Interval(
	left Bracket,
	minimum V,
	maximum V,
	right Bracket,
) IntervalLike[V] {
	var instance = &interval_[V]{
		// Initialize the instance attributes.
		left_:    left,
		minimum_: minimum,
		maximum_: maximum,
		right_:   right,
	}
	instance.validateInterval()
	return instance
}

// Constant Methods

// Function Methods

// INSTANCE INTERFACE

// Principal Methods

func (v *interval_[V]) GetClass() IntervalClassLike[V] {
	return intervalClass[V]()
}

// Accessible[V] Methods

func (v *interval_[V]) GetValue(
	index int,
) V {
	var size = v.effectiveSize()
	var offset = v.effectiveMinimum() + uti.RelativeToCardinal(index, size)
	return v.valueOf(offset)
}

func (v *interval_[V]) GetValues(
	first int,
	last int,
) col.Sequential[V] {
	var firstValue = v.GetValue(first)
	var lastValue = v.GetValue(last)
	return v.GetClass().Interval(
		Inclusive,
		firstValue,
		lastValue,
		Inclusive,
	)
}

func (v *interval_[V]) GetIndex(
	value V,
) int {
	var index = 1
	var iterator = v.GetIterator()
	for iterator.HasNext() {
		var candidate = iterator.GetNext()
		if candidate.AsInteger() == value.AsInteger() {
			return index
		}
		index++
	}
	return 0
}

// Bounded[V] Methods

func (v *interval_[V]) GetLeft() Bracket {
	return v.left_
}

func (v *interval_[V]) SetLeft(left Bracket) {
	v.left_ = left
	v.validateInterval()
}

func (v *interval_[V]) GetMinimum() V {
	return v.minimum_
}

func (v *interval_[V]) SetMinimum(minimum V) {
	v.minimum_ = minimum
	v.validateInterval()
}

func (v *interval_[V]) GetMaximum() V {
	return v.maximum_
}

func (v *interval_[V]) SetMaximum(maximum V) {
	v.maximum_ = maximum
	v.validateInterval()
}

func (v *interval_[V]) GetRight() Bracket {
	return v.right_
}

func (v *interval_[V]) SetRight(right Bracket) {
	v.right_ = right
	v.validateInterval()
}

// Searchable[V] Methods

func (v *interval_[V]) ContainsValue(
	value V,
) bool {
	if v.minimum_.IsDefined() && value.AsInteger() < v.minimum_.AsInteger() {
		return false
	}
	if v.maximum_.IsDefined() && value.AsInteger() > v.maximum_.AsInteger() {
		return false
	}
	return true
}

func (v *interval_[V]) ContainsAny(
	values col.Sequential[V],
) bool {
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		if v.ContainsValue(value) {
			// This interval contains at least one of the values.
			return true
		}
	}
	// This interval does not contain any of the values.
	return false
}

func (v *interval_[V]) ContainsAll(
	values col.Sequential[V],
) bool {
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		if !v.ContainsValue(value) {
			// This interval is missing at least one of the values.
			return false
		}
	}
	// This interval does contains all of the values.
	return true
}

// col.Sequential[V] Methods

func (v *interval_[V]) IsEmpty() bool {
	return false
}

func (v *interval_[V]) GetSize() uint {
	return v.effectiveSize()
}

func (v *interval_[V]) AsArray() []V {
	var size = int(v.effectiveSize())
	if size > 256 {
		// Limit the size to something reasonable.
		size = 256
	}
	var array = make([]V, size)
	for index := 0; index < size; index++ {
		var offset = v.effectiveMinimum() + index
		var value = v.valueOf(offset)
		array[index] = value
	}
	return array
}

func (v *interval_[V]) GetIterator() uti.IteratorLike[V] {
	var iterator = &iterator_[V]{
		interval_: v,
	}
	return iterator
}

// PROTECTED INTERFACE

func (v *interval_[V]) String() string {
	var source string
	switch v.left_ {
	case Inclusive:
		source += "["
	case Exclusive:
		source += "("
	}
	if v.minimum_.IsDefined() {
		source += v.minimum_.AsSource()
	}
	source += ".."
	if v.maximum_.IsDefined() {
		source += v.maximum_.AsSource()
	}
	switch v.right_ {
	case Inclusive:
		source += "]"
	case Exclusive:
		source += ")"
	}
	return source
}

// Private Methods

func (v *interval_[V]) effectiveMaximum() int {
	var maximum = v.maximum_.AsInteger()
	maximum -= int(v.right_)
	return maximum
}

func (v *interval_[V]) effectiveMinimum() int {
	var minimum = v.minimum_.AsInteger()
	minimum += int(v.left_)
	return minimum
}

func (v *interval_[V]) effectiveSize() uint {
	var size = uint(v.effectiveMaximum() - v.effectiveMinimum() + 1)
	return size
}

// This method ensures that the endpoints are valid.
func (v *interval_[V]) validateInterval() {
	// Validate the left bracket.
	switch v.left_ {
	case Inclusive:
	case Exclusive:
	default:
		var message = fmt.Sprintf(
			"Received an invalid left bracket for an interval: %v",
			v.left_,
		)
		panic(message)
	}

	// Validate the right bracket.
	switch v.right_ {
	case Inclusive:
	case Exclusive:
	default:
		var message = fmt.Sprintf(
			"Received an invalid right bracket for an interval: %v",
			v.right_,
		)
		panic(message)
	}

	// Validate the endpoints.
	if v.minimum_.IsDefined() && v.maximum_.IsDefined() {
		var collator = age.CollatorClass[V]().Collator()
		if collator.RankValues(v.minimum_, v.maximum_) != age.LesserRank {
			var message = fmt.Sprintf(
				"The minimum %v in an interval must be less than the maximum %v.",
				v.minimum_,
				v.maximum_,
			)
			panic(message)
		}
		var size = v.effectiveSize()
		if size <= 0 {
			var message = fmt.Sprintf(
				"The effective size of an interval must be greater than zero: %v.",
				size,
			)
			panic(message)
		}
	}
}

func (v *interval_[V]) valueOf(offset int) V {
	// We cannot get access to the constructor we need without reflection.
	var offsetRef = ref.ValueOf(offset)
	var resultRef ref.Value
	var valueRef = ref.ValueOf(v.minimum_)
	var method = valueRef.MethodByName("GetClass")
	var classRef = method.Call([]ref.Value{})[0]
	var typeRef = classRef.Type()
	var count = typeRef.NumMethod()
	for index := 0; index < count; index++ {
		var name = typeRef.Method(index).Name
		if sts.HasSuffix(name, "FromInteger") {
			resultRef = classRef.MethodByName(name).Call([]ref.Value{offsetRef})[0]
			break
		}
	}
	return resultRef.Interface().(V)
}

// Instance Structure

type interval_[V Discrete] struct {
	// Declare the instance attributes.
	left_    Bracket
	minimum_ V
	maximum_ V
	right_   Bracket
}

// Class Structure

type intervalClass_[V Discrete] struct {
	// Declare the class constants.
}

// Class Reference

var intervalMap_ = map[string]any{}
var intervalMutex_ syn.Mutex

func intervalClass[V Discrete]() *intervalClass_[V] {
	// Generate the name of the bound class type.
	var class *intervalClass_[V]
	var name = fmt.Sprintf("%T", class)

	// Check for an existing bound class type.
	intervalMutex_.Lock()
	var value = intervalMap_[name]
	switch actual := value.(type) {
	case *intervalClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Add a new bound class type.
		class = &intervalClass_[V]{
			// Initialize the class constants.
		}
		intervalMap_[name] = class
	}
	intervalMutex_.Unlock()

	// Return a reference to the bound class type.
	return class
}

/*
NOTE:
The following is a private implementation of an Iterator class that operates
on an interval directly without requiring the overhead of creating an array.
This allows iteration over portions of very large intervals with no memory
overhead.
*/

type iterator_[V Discrete] struct {
	slot_     uint
	interval_ IntervalLike[V]
}

func (v *iterator_[V]) IsEmpty() bool {
	return v.interval_.GetSize() == 0
}

func (v *iterator_[V]) ToStart() {
	v.slot_ = 0
}

func (v *iterator_[V]) ToEnd() {
	v.slot_ = v.interval_.GetSize()
}

func (v *iterator_[V]) HasPrevious() bool {
	return v.slot_ > 0
}

func (v *iterator_[V]) GetPrevious() V {
	var result_ V
	if v.slot_ > 0 {
		result_ = v.interval_.GetValue(int(v.slot_))
		v.slot_--
	}
	return result_
}

func (v *iterator_[V]) HasNext() bool {
	return v.slot_ < v.interval_.GetSize()
}

func (v *iterator_[V]) GetNext() V {
	var result_ V
	if v.slot_ < v.interval_.GetSize() {
		v.slot_++
		result_ = v.interval_.GetValue(int(v.slot_))
	}
	return result_
}

func (v *iterator_[V]) GetSize() uint {
	return v.interval_.GetSize()
}

func (v *iterator_[V]) GetSlot() uint {
	return v.slot_
}

func (v *iterator_[V]) SetSlot(
	slot uint,
) {
	var size = v.interval_.GetSize()
	if slot > size {
		slot = size
	}
	v.slot_ = slot
}
