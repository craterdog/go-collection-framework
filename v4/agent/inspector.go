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

package agent

import (
	ref "reflect"
)

// CLASS ACCESS

// Reference

var inspectorClass = &inspectorClass_{
	// This class has no private constants to initialize.
}

// Function

func Inspector() InspectorClassLike {
	return inspectorClass
}

// CLASS METHODS

// Target

type inspectorClass_ struct {
	// This class has no private constants.
}

// Constructors

func (c *inspectorClass_) Make() InspectorLike {
	return &inspector_{}
}

// INSTANCE METHODS

// Target

type inspector_ struct {
	class_ InspectorClassLike
}

// Attributes

func (v *inspector_) GetClass() InspectorClassLike {
	return v.class_
}

// Public

func (v *inspector_) ImplementsAspect(
	value any,
	aspect any,
) bool {
	return v.implementsAspect(value, aspect)
}

func (v *inspector_) IsDefined(value any) bool {
	return v.isDefined(value)
}

// Private

func (v *inspector_) implementsAspect(
	value any,
	aspect any,
) bool {
	if v.isDefined(value) {
		var reflectedType = ref.TypeOf(value)
		var aspectType = ref.TypeOf(aspect).Elem()
		return reflectedType.Implements(aspectType)
	}
	return false
}

func (v *inspector_) isDefined(value any) bool {
	// This method addresses the inconsistencies in the Go language with respect
	// to whether or not a value is defined or not.  Go handles interfaces,
	// pointers and various primitive types differently.  This makes consistent
	// checking across different types problematic.  We handle it here in one
	// place (hopefully correctly).
	switch actual := value.(type) {
	case string:
		return len(actual) > 0
	default:
		var meta = ref.ValueOf(actual)
		var isPointer = meta.Kind() == ref.Ptr ||
			meta.Kind() == ref.Interface ||
			meta.Kind() == ref.Slice ||
			meta.Kind() == ref.Map ||
			meta.Kind() == ref.Chan ||
			meta.Kind() == ref.Func
		var isNil = isPointer && meta.IsNil()
		return !isNil && meta.IsValid()
	}
}
