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

package cdcn_test

import (
	fmt "fmt"
	not "github.com/craterdog/go-collection-framework/v4/cdcn"
	ass "github.com/stretchr/testify/assert"
	osx "os"
	sts "strings"
	tes "testing"
)

const collectionTests = "../test/input/"

func TestCollectionRoundtrips(t *tes.T) {
	var notation = not.Notation().Make()
	var files, err = osx.ReadDir(collectionTests)
	if err != nil {
		var message = fmt.Sprintf("Could not find the %s directory.", collectionTests)
		panic(message)
	}
	for _, file := range files {
		var filename = collectionTests + file.Name()
		if sts.HasSuffix(filename, ".cdcn") {
			fmt.Println(filename)
			var bytes, err = osx.ReadFile(filename)
			if err != nil {
				panic(err)
			}
			var expected = string(bytes)
			var collection = notation.ParseSource(expected)
			var actual = notation.FormatCollection(collection)
			if !sts.HasPrefix(file.Name(), "map") {
				// Skip maps since they are non-deterministic.
				ass.Equal(t, expected, actual)
				bytes = []byte(actual)
				err = osx.WriteFile(filename, bytes, 0644)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
