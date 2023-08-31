/*******************************************************************************
 *   Copyright (c) 2009-2023 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package collections_test

import (
	fmt "fmt"
	bal "github.com/craterdog/go-collection-framework/v2"
	osx "os"
	tes "testing"
)

const filename = "./collections.cdsn"

func TestGenerateGrammar(t *tes.T) {
	var err = osx.WriteFile(filename, []byte(bal.FormatGrammar()), 0644)
	if err != nil {
		var message = fmt.Sprintf("Could not create the bwsn file: %v.", err)
		panic(message)
	}
}
