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
	syn "sync"
)

// CLASS INTERFACE

// Access Function

func ContinuumClass[V Continuous]() ContinuumClassLike[V] {
	return continuumClass[V]()
}

// Constructor Methods

func (c *continuumClass_[V]) Continuum(
	left Bracket,
	minimum V,
	maximum V,
	right Bracket,
) ContinuumLike[V] {
	var instance = &continuum_[V]{
		// Initialize the instance attributes.
		left_:    left,
		minimum_: minimum,
		maximum_: maximum,
		right_:   right,
	}
	instance.validateContinuum()
	return instance
}

// Constant Methods

// Function Methods

// INSTANCE INTERFACE

// Principal Methods

func (v *continuum_[V]) GetClass() ContinuumClassLike[V] {
	return continuumClass[V]()
}

// Bounded[V] Methods

func (v *continuum_[V]) GetLeft() Bracket {
	return v.left_
}

func (v *continuum_[V]) SetLeft(left Bracket) {
	v.left_ = left
	v.validateContinuum()
}

func (v *continuum_[V]) GetMinimum() V {
	return v.minimum_
}

func (v *continuum_[V]) SetMinimum(minimum V) {
	v.minimum_ = minimum
	v.validateContinuum()
}

func (v *continuum_[V]) GetMaximum() V {
	return v.maximum_
}

func (v *continuum_[V]) SetMaximum(maximum V) {
	v.maximum_ = maximum
	v.validateContinuum()
}

func (v *continuum_[V]) GetRight() Bracket {
	return v.right_
}

func (v *continuum_[V]) SetRight(right Bracket) {
	v.right_ = right
	v.validateContinuum()
}

// Searchable[V] Methods

func (v *continuum_[V]) ContainsValue(
	value V,
) bool {
	if v.minimum_.IsDefined() && value.AsFloat() < v.minimum_.AsFloat() {
		return false
	}
	if v.maximum_.IsDefined() && value.AsFloat() > v.maximum_.AsFloat() {
		return false
	}
	return true
}

func (v *continuum_[V]) ContainsAny(
	values col.Sequential[V],
) bool {
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		if v.ContainsValue(value) {
			// This continuum contains at least one of the values.
			return true
		}
	}
	// This continuum does not contain any of the values.
	return false
}

func (v *continuum_[V]) ContainsAll(
	values col.Sequential[V],
) bool {
	var iterator = values.GetIterator()
	for iterator.HasNext() {
		var value = iterator.GetNext()
		if !v.ContainsValue(value) {
			// This continuum is missing at least one of the values.
			return false
		}
	}
	// This continuum does contains all of the values.
	return true
}

// PROTECTED INTERFACE

func (v *continuum_[V]) String() string {
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

// This method ensures that the endpoints are valid.
func (v *continuum_[V]) validateContinuum() {
	// Validate the left bracket.
	switch v.left_ {
	case Inclusive:
	case Exclusive:
	default:
		var message = fmt.Sprintf(
			"Received an invalid left bracket for a continuum: %v",
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
			"Received an invalid right bracket for a continuum: %v",
			v.right_,
		)
		panic(message)
	}

	// Validate the endpoints.
	if v.minimum_.IsDefined() && v.maximum_.IsDefined() {
		var collator = age.CollatorClass[V]().Collator()
		if collator.RankValues(v.minimum_, v.maximum_) != age.LesserRank {
			var message = fmt.Sprintf(
				"The minimum %v in a continuum must be less than the maximum %v.",
				v.minimum_,
				v.maximum_,
			)
			panic(message)
		}
		var size = v.maximum_.AsFloat() - v.minimum_.AsFloat()
		if size <= 0 {
			var message = fmt.Sprintf(
				"The size of a continuum must be greater than zero: %v.",
				size,
			)
			panic(message)
		}
	}
}

// Instance Structure

type continuum_[V Continuous] struct {
	// Declare the instance attributes.
	left_    Bracket
	minimum_ V
	maximum_ V
	right_   Bracket
}

// Class Structure

type continuumClass_[V Continuous] struct {
	// Declare the class constants.
}

// Class Reference

var continuumMap_ = map[string]any{}
var continuumMutex_ syn.Mutex

func continuumClass[V Continuous]() *continuumClass_[V] {
	// Generate the name of the bound class type.
	var class *continuumClass_[V]
	var name = fmt.Sprintf("%T", class)

	// Check for an existing bound class type.
	continuumMutex_.Lock()
	var value = continuumMap_[name]
	switch actual := value.(type) {
	case *continuumClass_[V]:
		// This bound class type already exists.
		class = actual
	default:
		// Add a new bound class type.
		class = &continuumClass_[V]{
			// Initialize the class constants.
		}
		continuumMap_[name] = class
	}
	continuumMutex_.Unlock()

	// Return a reference to the bound class type.
	return class
}
