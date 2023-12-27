/*******************************************************************************
 *   Copyright (c) 2009-2024 Crater Dog Technologies™.  All Rights Reserved.   *
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
	col "github.com/craterdog/go-collection-framework/v3"
	ass "github.com/stretchr/testify/assert"
	osx "os"
	sts "strings"
	tes "testing"
)

const collectionTests = "./test/"

func TestCollectionRoundtrips(t *tes.T) {
	var cdcn = col.CDCNClass().Default()
	var files, err = osx.ReadDir(collectionTests)
	if err != nil {
		var message = fmt.Sprintf("Could not find the %s directory.", collectionTests)
		panic(message)
	}
	for _, file := range files {
		var filename = collectionTests + file.Name()
		if sts.HasSuffix(filename, ".cdcn") {
			fmt.Println(filename)
			var source, _ = osx.ReadFile(filename)
			var expected = string(source[:len(source)-1])
			var collection = cdcn.ParseCollection(expected)
			var actual = cdcn.FormatCollection(collection)
			if !sts.HasPrefix(file.Name(), "map") {
				// Skip maps since they are non-deterministic.
				ass.Equal(t, expected, actual)
			}
		}
	}
}
