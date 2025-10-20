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

package agents

import (
	fmt "fmt"
	uti "github.com/craterdog/go-missing-utilities/v8"
)

// CLASS INTERFACE

// Access Function

func ControllerClass() ControllerClassLike {
	return controllerClass()
}

// Constructor Methods

func (c *controllerClass_) Controller(
	events []Event,
	transitions map[State]Transitions,
	initialState State,
) ControllerLike {
	// Validate the constructor arguments.
	if uti.IsUndefined(events) {
		panic("The \"events\" attribute is required by this class.")
	}
	if uti.IsUndefined(transitions) {
		panic("The \"transitions\" attribute is required by this class.")
	}
	if uti.IsUndefined(initialState) {
		panic("The \"initialState\" attribute is required by this class.")
	}
	var height = len(transitions)
	if height < 2 {
		var message = fmt.Sprintf(
			"The state table must have at least two possible transitions: %v",
			height,
		)
		panic(message)
	}
	var width = len(events)
	if width < 1 {
		var message = fmt.Sprintf(
			"The state table must have at least one possible event: %v",
			width,
		)
		panic(message)
	}
	for _, row := range transitions {
		if len(row) != width {
			var message = fmt.Sprintf(
				"Each row in the state table must be the same width: %v",
				width,
			)
			panic(message)
		}
	}
	var invalidState = true
	for candidate := range transitions {
		if candidate == initialState && initialState != c.invalid_ {
			invalidState = false
			break
		}
	}
	if invalidState {
		var message = fmt.Sprintf(
			"The initial state is invalid: %q",
			initialState,
		)
		panic(message)
	}

	// Create a new instance.
	var instance = &controller_{
		// Initialize the instance attributes.
		state_:       initialState,
		events_:      events,
		transitions_: transitions,
	}
	return instance
}

// Constant Methods

func (c *controllerClass_) Invalid() State {
	return c.invalid_
}

// Function Methods

// INSTANCE INTERFACE

// Principal Methods

func (v *controller_) GetClass() ControllerClassLike {
	return controllerClass()
}

func (v *controller_) ProcessEvent(
	event Event,
) State {
	var index = v.eventIndex(event)
	if index < 0 {
		var message = fmt.Sprintf(
			"Attempted to process an invalid event %q.",
			event,
		)
		panic(message)
	}
	var next = v.transitions_[v.state_][index]
	if !v.hasState(next) {
		var message = fmt.Sprintf(
			"Attempted to transition from state %q to an invalid state on event %q.",
			v.state_,
			event,
		)
		panic(message)
	}
	v.state_ = next
	return next
}

// Attribute Methods

func (v *controller_) GetState() State {
	return v.state_
}

func (v *controller_) SetState(
	state State,
) {
	if uti.IsUndefined(state) || !v.hasState(state) {
		var message = fmt.Sprintf(
			"A valid \"state\" argument is required: %v",
			state,
		)
		panic(message)
	}
	v.state_ = state
}

func (v *controller_) GetEvents() []Event {
	return v.events_
}

func (v *controller_) GetTransitions() map[State]Transitions {
	return v.transitions_
}

// PROTECTED INTERFACE

// Private Methods

func (v *controller_) eventIndex(
	event Event,
) int {
	for index, candidate := range v.events_ {
		if candidate == event {
			return index
		}
	}
	return -1
}

func (v *controller_) hasState(
	state State,
) bool {
	for candidate := range v.transitions_ {
		if candidate == state {
			return true
		}
	}
	return false
}

// Instance Structure

type controller_ struct {
	// Declare the instance attributes.
	state_       State
	events_      []Event
	transitions_ map[State]Transitions
}

// Class Structure

type controllerClass_ struct {
	// Declare the class constants.
	invalid_ State
}

// Class Reference

func controllerClass() *controllerClass_ {
	return controllerClassReference_
}

var controllerClassReference_ = &controllerClass_{
	// Initialize the class constants.
	invalid_: "$Invalid",
}
