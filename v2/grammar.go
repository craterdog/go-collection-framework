/*******************************************************************************
 *   Copyright (c) 2009-2023 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

package collections

import (
	fmt "fmt"
	sts "strings"
)

// This map captures the syntax rules for collections of Go primitives.
// This map is useful when creating scanner and parser error messages.
var grammar = map[string]string{
	"$BASE10":    `"0".."9"`,
	"$BASE16":    `"0".."9" | "a".."f"`,
	"$BOOLEAN":   `"false" | "true"`,
	"$COMPLEX":   `"(" FLOAT SIGN FLOAT "i)"`,
	"$CONTEXT":   `"array" | "catalog" | "list" | "map" | "queue" | "set" | "stack"`,
	"$DELIMITER": `"]" | "[" | ")" | "(" | ":" | ","`,
	"$EOL":       `"\n"`,
	"$ESCAPE":    `'\' (UNICODE | 'a' | 'b' | 'f' | 'n' | 'r' | 't' | 'v' | '"' | "'" | '\')`,
	"$EXPONENT":  `("e" | "E") SIGN ORDINAL`,
	"$FLOAT":     `[SIGN] SCALAR [EXPONENT]`,
	"$FRACTION":  `"." <BASE10>`,
	"$INTEGER":   `ZERO | [SIGN] ORDINAL`,
	"$NIL":       `"nil"`,
	"$ORDINAL":   `"1".."9" {"0".."9"}`,
	"$RUNE":      `"'" (ESCAPE | ~("'" | EOL)) "'"`,
	"$SCALAR":    `(ZERO | ORDINAL) FRACTION`,
	"$SIGN":      `"+" | "-"`,
	"$STRING":    `'"' {ESCAPE | ~('"' | EOL)} '"'`,
	"$UNICODE": `
    "u" BASE16 BASE16 BASE16 BASE16 |
    "U" BASE16 BASE16 BASE16 BASE16 BASE16 BASE16 BASE16 BASE16`,
	"$UNSIGNED":     `"0x" <BASE16>`,
	"$ZERO":         `"0"`,
	"$association":  `key ":" value`,
	"$associations": `
    association {"," association} |
    EOL <association EOL> |
    ":"  ! No associations.`,
	"$collection": `"[" (values | associations) "]" "(" CONTEXT ")"`,
	"$document":  `collection EOF  ! EOF is the end-of-file marker.`,
	"$key":       `primitive`,
	"$primitive": `BOOLEAN | COMPLEX | FLOAT | INTEGER | NIL | RUNE | STRING`,
	"$value":     `primitive | collection`,
	"$values": `
    value {"," value} |
    EOL <value EOL> |
    ! No values.`,
}

const header = `!>
    A formal definition of a document containing a collection using Bali Wirth
    Syntax Notation™ (BWSN):
        <https://github.com/bali-nebula/specifications/blob/main/bwsn.bwsn>

    The token names are identified by all CAPITAL characters and the rule names
    are identified by lowerCamelCase characters. The token and rule definitions
    have been alphabetized to make it easier to locate specific definitions.
    The starting rule is "$document".
<!

`

func FormatGrammar() string {
	var builder sts.Builder
	builder.WriteString(header)
	var unsorted = make([]string, len(grammar))
	var index = 0
	for key := range grammar {
		unsorted[index] = key
		index++
	}
	var keys = ListFromArray(unsorted)
	keys.SortValues()
	var iterator = Iterator[string](keys)
	for iterator.HasNext() {
		var key = iterator.GetNext()
		var value = grammar[key]
		builder.WriteString(fmt.Sprintf("%s: %s\n\n", key, value))
	}
	return builder.String()
}

// PRIVATE FUNCTIONS

func generateGrammar(expected string, symbols ...string) string {
	var message = "Was expecting '" + expected + "' from:\n"
	for _, symbol := range symbols {
		message += fmt.Sprintf("  \033[32m%v: \033[33m%v\033[0m\n\n", symbol, grammar[symbol])
	}
	return message
}
