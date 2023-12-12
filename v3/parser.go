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
	stc "strconv"
	sts "strings"
	utf "unicode/utf8"
)

// PARSER INTERFACE

// This function parses the specified document source retrieved from a POSIX
// compliant file and returns the corresponding collection that was used to
// generate the document using the collection formatting capabilities.
// A POSIX compliant file must end with an EOF marker.
func ParseDocument(source []byte) Collection {
	var ok bool
	var token TokenLike
	var collection Collection
	var tokens = make(chan TokenLike, 256)
	Scanner().FromSource(source, tokens) // Starts scanning in a separate go routine.
	var Stack = Stack[TokenLike]()
	var p = &parser{
		source: source,
		next:   Stack.WithCapacity(4),
		tokens: tokens,
	}
	collection, token, ok = p.parseCollection()
	if !ok {
		var message = p.formatError(token)
		message += generateGrammar("collection",
			"$source",
			"$collection",
			"$values",
			"$associations",
			"$context")
		panic(message)
	}
	_, token, ok = p.parseEOF()
	if !ok {
		var message = p.formatError(token)
		message += generateGrammar("EOF",
			"$source",
			"$collection",
			"$values",
			"$associations",
			"$context")
		panic(message)
	}
	return collection
}

// This function parses a source string rather than the bytes from a document
// file. It is useful when parsing strings within source code.
func ParseCollection(source string) Collection {
	var document = []byte(source + EOF) // Append the POSIX compliant EOF marker.
	return ParseDocument(document)
}

// PARSER IMPLEMENTATION

// This type defines the structure and methods for the parser agent.
type parser struct {
	source         []byte
	next           StackLike[TokenLike] // The stack of the retrieved tokens that have been put back.
	tokens         chan TokenLike       // The queue of unread tokens coming from the scanner.
	p1, p2, p3, p4 TokenLike            // The previous four tokens that have been retrieved.
}

// This method attempts to read the next token from the token stream and return
// it.
func (v *parser) nextToken() TokenLike {
	var Token = Token()
	var next TokenLike
	if v.next.IsEmpty() {
		var token, ok = <-v.tokens
		if !ok {
			panic("The token channel terminated without an EOF or error token.")
		}
		next = token
		if next.GetType() == Token.TypeError() {
			var message = v.formatError(next)
			panic(message)
		}
	} else {
		next = v.next.RemoveTop()
	}
	v.p4, v.p3, v.p2, v.p1 = v.p3, v.p2, v.p1, next
	return next
}

// This method puts back the current token onto the token stream so that it can
// be retrieved by another parsing method.
func (v *parser) backupOne() {
	v.next.AddValue(v.p1)
	v.p1, v.p2, v.p3, v.p4 = v.p2, v.p3, v.p4, nil
}

// This method returns an error message containing the context for a parsing
// error.
func (v *parser) formatError(token TokenLike) string {
	var message = fmt.Sprintf("An unexpected token was received by the parser: %v\n", token)
	var line = token.GetLine()
	var lines = sts.Split(string(v.source), EOL)

	message += "\033[36m"
	if line > 1 {
		message += fmt.Sprintf("%04d: ", line-1) + string(lines[line-2]) + EOL
	}
	message += fmt.Sprintf("%04d: ", line) + string(lines[line-1]) + EOL

	message += " \033[32m>>>─"
	var count = 0
	for count < token.GetPosition() {
		message += "─"
		count++
	}
	message += "⌃\033[36m\n"

	if line < len(lines) {
		message += fmt.Sprintf("%04d: ", line+1) + string(lines[line]) + EOL
	}
	message += "\033[0m\n"

	return message
}

// This method attempts to parse an association between a key and value. It
// returns the association and whether or not the association was successfully
// parsed.
func (v *parser) parseAssociation() (AssociationLike[Key, Value], TokenLike, bool) {
	var ok bool
	var token TokenLike
	var key Key
	var value Value
	var association AssociationLike[Key, Value]
	key, token, ok = v.parsePrimitive()
	if !ok {
		// This is not an association.
		return association, token, false
	}
	_, token, ok = v.parseDelimiter(":")
	if !ok {
		// This is not an association.
		v.backupOne() // Put back the primitive key token.
		return association, token, false
	}
	// This must be an association.
	value, token, ok = v.parseValue()
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("value",
			"$association",
			"$key",
			"$value")
		panic(message)
	}
	var Association = Association[Key, Value]()
	association = Association.FromPair(key, value)
	return association, token, true
}

// This method attempts to parse a collection of values. It returns the
// collection and whether or not the collection was successfully parsed.
func (v *parser) parseCollection() (Collection, TokenLike, bool) {
	var ok bool
	var token TokenLike
	var collection Collection
	var context string
	collection, token, ok = v.parseAssociations()
	if !ok {
		// The sequence of values must be attempted last since it starts
		// with a value which cannot be put back as a single token.
		collection, token, ok = v.parseValues()
	}
	if ok {
		_, token, ok = v.parseDelimiter("(")
		if !ok {
			var message = v.formatError(token)
			message += generateGrammar("(",
				"$collection",
				"$context")
			panic(message)
		}
		context, token, ok = v.parseContext()
		if !ok {
			var message = v.formatError(token)
			message += generateGrammar("context",
				"$collection",
				"$context")
			panic(message)
		}
		switch sequence := collection.(type) {
		case Sequential[Value]:
			switch context {
			case "array":
				var Array = Array[Value]()
				collection = Array.FromArray(sequence.AsArray())
			case "list":
				var List = List[Value]()
				collection = List.FromSequence(sequence)
			case "queue":
				var Queue = Queue[Value]()
				collection = Queue.FromSequence(sequence)
			case "set":
				var Set = Set[Value]()
				collection = Set.FromSequence(sequence)
			case "stack":
				var Stack = Stack[Value]()
				collection = Stack.FromSequence(sequence)
			default:
			}
		case Sequential[Binding[Key, Value]]:
			switch context {
			case "catalog":
				var Catalog = Catalog[Key, Value]()
				collection = Catalog.FromSequence(sequence)
			case "map":
				var Map = Map[Key, Value]()
				collection = Map.FromArray(sequence.AsArray())
			default:
			}
		default:
		}
		_, token, ok = v.parseDelimiter(")")
		if !ok {
			var message = v.formatError(token)
			message += generateGrammar(")",
				"$collection",
				"$context")
			panic(message)
		}
	}
	return collection, token, ok
}

// This method attempts to parse the specified delimiter. It returns
// the token and whether or not the delimiter was found.
func (v *parser) parseDelimiter(delimiter string) (string, TokenLike, bool) {
	var Token = Token()
	var token = v.nextToken()
	if token.GetType() == Token.TypeEOF() || token.GetValue() != delimiter {
		v.backupOne()
		return delimiter, token, false
	}
	return delimiter, token, true
}

// This method attempts to parse the end-of-file (EOF) marker. It returns
// the token and whether or not an EOL token was found.
func (v *parser) parseEOF() (TokenLike, TokenLike, bool) {
	var Token = Token()
	var token = v.nextToken()
	if token.GetType() != Token.TypeEOF() {
		v.backupOne()
		return token, token, false
	}
	return token, token, true
}

// This method attempts to parse the end-of-line (EOL) marker. It returns
// the token and whether or not an EOF token was found.
func (v *parser) parseEOL() (TokenLike, TokenLike, bool) {
	var Token = Token()
	var token = v.nextToken()
	if token.GetType() != Token.TypeEOL() {
		v.backupOne()
		return token, token, false
	}
	return token, token, true
}

// It returns a sequence of associations and whether or not the sequence of
// associations was successfully parsed.
func (v *parser) parseInlineAssociations() (Sequential[Binding[Key, Value]], TokenLike, bool) {
	var ok bool
	var token TokenLike
	var association Binding[Key, Value]
	var List = List[Binding[Key, Value]]()
	var associations = List.Empty()
	_, token, ok = v.parseDelimiter(":")
	if ok {
		// This is an empty sequence of associations.
		return associations, token, true
	}
	_, token, ok = v.parseDelimiter("]")
	if ok {
		// This is an empty sequence of values.
		v.backupOne() // Put back the ']' character.
		return associations, token, false
	}
	association, token, ok = v.parseAssociation()
	if !ok {
		// A non-empty sequence must have at least one association.
		return associations, token, false
	}
	for {
		associations.AppendValue(association)
		// Every subsequent association must be preceded by a ','.
		_, token, ok = v.parseDelimiter(",")
		if !ok {
			// No more associations.
			break
		}
		association, token, ok = v.parseAssociation()
		if !ok {
			var message = v.formatError(token)
			message += generateGrammar("association",
				"$associations",
				"$association",
				"$key",
				"$value")
			panic(message)
		}
	}
	return associations, token, true
}

// This method attempts to parse a sequence containing inline values. It returns
// the sequence of values and whether or not the sequence of values was
// successfully parsed.
func (v *parser) parseInlineValues() (Sequential[Value], TokenLike, bool) {
	var ok bool
	var token TokenLike
	var value Value
	var List = List[Value]()
	var values = List.Empty()
	_, token, ok = v.parseDelimiter("]")
	if ok {
		// This is an empty sequence of values.
		v.backupOne() // Put back the ']' token.
		return values, token, true
	}
	value, token, ok = v.parseValue()
	if !ok {
		// A non-empty sequence must have at least one value value.
		var message = v.formatError(token)
		message += generateGrammar("value",
			"$values",
			"$value")
		panic(message)
	}
	for {
		values.AppendValue(value)
		// Every subsequent value must be preceded by a ','.
		_, token, ok = v.parseDelimiter(",")
		if !ok {
			// No more values.
			break
		}
		value, token, ok = v.parseValue()
		if !ok {
			var message = v.formatError(token)
			message += generateGrammar("value",
				"$values",
				"$value")
			panic(message)
		}
	}
	return values, token, true
}

// This method attempts to parse a sequence of associations. It returns the
// sequence of associations and whether or not the sequence of associations was
// successfully parsed.
func (v *parser) parseAssociations() (Sequential[Binding[Key, Value]], TokenLike, bool) {
	var ok bool
	var token TokenLike
	var associations Sequential[Binding[Key, Value]]
	_, token, ok = v.parseDelimiter("[")
	if !ok {
		return associations, token, false
	}
	_, _, ok = v.parseEOL()
	if !ok {
		associations, token, ok = v.parseInlineAssociations()
		if !ok {
			v.backupOne() // Put back the '[' character.
			return associations, token, false
		}
	} else {
		associations, token, ok = v.parseMultilineAssociations()
		if !ok {
			v.backupOne() // Put back the EOL character.
			v.backupOne() // Put back the '[' character.
			return associations, token, false
		}
	}
	_, token, ok = v.parseDelimiter("]")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("]",
			"$associations",
			"$association",
			"$key",
			"$value")
		panic(message)
	}
	return associations, token, true
}

// This method attempts to parse a sequence containing multiline associations.
// It returns the sequence of associations and whether or not the sequence of
// associations was successfully parsed.
func (v *parser) parseMultilineAssociations() (Sequential[Binding[Key, Value]], TokenLike, bool) {
	var ok bool
	var token TokenLike
	var association AssociationLike[Key, Value]
	var List = List[Binding[Key, Value]]()
	var associations = List.Empty()
	association, token, ok = v.parseAssociation()
	if !ok {
		// A non-empty sequence must have at least one association.
		return associations, token, false
	}
	for {
		associations.AppendValue(association)
		// Every association must be followed by an EOL.
		_, token, ok = v.parseEOL()
		if !ok {
			var message = v.formatError(token)
			message += generateGrammar("EOL",
				"$associations",
				"$association",
				"$key",
				"$value")
			panic(message)
		}
		association, token, ok = v.parseAssociation()
		if !ok {
			// No more associations.
			break
		}
	}
	return associations, token, true
}

// This method attempts to parse a sequence containing multiline values. It
// returns the sequence of values and whether or not the sequence of values was
// successfully parsed.
func (v *parser) parseMultilineValues() (Sequential[Value], TokenLike, bool) {
	var ok bool
	var token TokenLike
	var value Value
	var List = List[Value]()
	var values = List.Empty()
	value, token, ok = v.parseValue()
	if !ok {
		// A non-empty sequence must have at least one value.
		var message = v.formatError(token)
		message += generateGrammar("value",
			"$values",
			"$value")
		panic(message)
	}
	for {
		values.AppendValue(value)
		// Every value must be followed by an EOL.
		_, token, ok = v.parseEOL()
		if !ok {
			var message = v.formatError(token)
			message += generateGrammar("EOL",
				"$values",
				"$value")
			panic(message)
		}
		value, token, ok = v.parseValue()
		if !ok {
			// No more values.
			break
		}
	}
	return values, token, true
}

// This method attempts to parse a component entity. It returns the component
// entity and whether or not the component entity was successfully parsed.
func (v *parser) parseValue() (Value, TokenLike, bool) {
	var ok bool
	var token TokenLike
	var value Value
	value, token, ok = v.parsePrimitive()
	if !ok {
		value, token, ok = v.parseCollection()
	}
	return value, token, ok
}

// This method attempts to parse a primitive. It returns the primitive and
// whether or not the primitive was successfully parsed.
func (v *parser) parsePrimitive() (Primitive, TokenLike, bool) {
	var ok bool
	var token TokenLike
	var primitive Primitive
	primitive, token, ok = v.parseBoolean()
	if !ok {
		primitive, token, ok = v.parseComplex()
	}
	if !ok {
		primitive, token, ok = v.parseFloat()
	}
	if !ok {
		primitive, token, ok = v.parseInteger()
	}
	if !ok {
		primitive, token, ok = v.parseNil()
	}
	if !ok {
		primitive, token, ok = v.parseRune()
	}
	if !ok {
		primitive, token, ok = v.parseString()
	}
	if !ok {
		primitive, token, ok = v.parseUnsigned()
	}
	if !ok {
		// Override any zero values returned from failed parsing attempts.
		primitive = nil
	}
	return primitive, token, ok
}

// This method attempts to parse a sequence of values. It returns the
// sequence of values and whether or not the sequence of values was
// successfully parsed.
func (v *parser) parseValues() (Sequential[Value], TokenLike, bool) {
	var ok bool
	var token TokenLike
	var values Sequential[Value]
	_, token, ok = v.parseDelimiter("[")
	if !ok {
		return values, token, false
	}
	_, _, ok = v.parseEOL()
	if !ok {
		values, token, ok = v.parseInlineValues()
		if !ok {
			v.backupOne() // Put back the '[' character.
			return values, token, false
		}
	} else {
		values, token, ok = v.parseMultilineValues()
		if !ok {
			v.backupOne() // Put back the EOL character.
			v.backupOne() // Put back the '[' character.
			return values, token, false
		}
	}
	_, token, ok = v.parseDelimiter("]")
	if !ok {
		var message = v.formatError(token)
		message += generateGrammar("]",
			"$values",
			"$value")
		panic(message)
	}
	return values, token, ok
}

// This method attempts to parse a boolean primitive. It returns the boolean
// primitive and whether or not the boolean primitive was successfully parsed.
func (v *parser) parseBoolean() (bool, TokenLike, bool) {
	var Token = Token()
	var boolean bool
	var token = v.nextToken()
	if token.GetType() != Token.TypeBoolean() {
		v.backupOne()
		return boolean, token, false
	}
	boolean, _ = stc.ParseBool(token.GetValue())
	return boolean, token, true
}

// This method attempts to parse a complex number primitive. It returns the
// complex number primitive and whether or not the complex number primitive was
// successfully parsed.
func (v *parser) parseComplex() (complex128, TokenLike, bool) {
	var Token = Token()
	var complex_ complex128
	var token = v.nextToken()
	if token.GetType() != Token.TypeComplex() {
		v.backupOne()
		return complex_, token, false
	}
	complex_, _ = stc.ParseComplex(token.GetValue(), 128)
	return complex_, token, true
}

// This method attempts to parse the type of a collection. It returns the type
// string and whether or not the type string was successfully parsed.
func (v *parser) parseContext() (string, TokenLike, bool) {
	var Token = Token()
	var token = v.nextToken()
	if token.GetType() != Token.TypeContext() {
		v.backupOne()
		return "", token, false
	}
	return token.GetValue(), token, true
}

// This method attempts to parse a floating point primitive. It returns the
// floating point primitive and whether or not the floating point primitive was
// successfully parsed.
func (v *parser) parseFloat() (float64, TokenLike, bool) {
	var Token = Token()
	var float float64
	var token = v.nextToken()
	if token.GetType() != Token.TypeFloat() {
		v.backupOne()
		return float, token, false
	}
	float, _ = stc.ParseFloat(token.GetValue(), 64)
	return float, token, true
}

// This method attempts to parse a integer primitive. It returns the integer
// primitive and whether or not the integer primitive was successfully parsed.
func (v *parser) parseInteger() (int64, TokenLike, bool) {
	var Token = Token()
	var integer int64
	var token = v.nextToken()
	if token.GetType() != Token.TypeInteger() {
		v.backupOne()
		return integer, token, false
	}
	integer, _ = stc.ParseInt(token.GetValue(), 10, 64)
	return integer, token, true
}

// This method attempts to parse a nil primitive. It returns the nil primitive
// and whether or not the nil primitive was successfully parsed.
func (v *parser) parseNil() (Value, TokenLike, bool) {
	var Token = Token()
	var token = v.nextToken()
	if token.GetType() != Token.TypeNil() {
		v.backupOne()
		return nil, token, false
	}
	return nil, token, true
}

// This method attempts to parse a rune. It returns the rune and whether or not
// the rune was successfully parsed.
func (v *parser) parseRune() (rune, TokenLike, bool) {
	var Token = Token()
	var rune_ rune
	var size int
	var token = v.nextToken()
	if token.GetType() != Token.TypeRune() {
		v.backupOne()
		return rune_, token, false
	}
	var matches = Scanner().MatchRune(token.GetValue())
	// We must unquote the full token string properly.
	var s, _ = stc.Unquote(matches[0])
	rune_, size = utf.DecodeRuneInString(s)
	if len(s) != size {
		// This is not a rune.
		v.backupOne() // Put back the quote token.
		return rune_, token, false
	}
	return rune_, token, true
}

// This method attempts to parse a string primitive. It returns the string
// primitive and whether or not the string primitive was successfully parsed.
func (v *parser) parseString() (string, TokenLike, bool) {
	var Token = Token()
	var string_ string
	var token = v.nextToken()
	if token.GetType() != Token.TypeString() {
		v.backupOne()
		return string_, token, false
	}
	var matches = Scanner().MatchString(token.GetValue())
	// We must unquote the full token string properly.
	string_, _ = stc.Unquote(matches[0])
	return string_, token, true
}

// This method attempts to parse an unsigned integer primitive. It returns the
// unsigned integer primitive and whether or not the unsigned integer primitive
// was successfully parsed.
func (v *parser) parseUnsigned() (uint64, TokenLike, bool) {
	var Token = Token()
	var unsigned uint64
	var token = v.nextToken()
	if token.GetType() != Token.TypeUnsigned() {
		v.backupOne()
		return unsigned, token, false
	}
	unsigned, _ = stc.ParseUint(token.GetValue()[2:], 16, 64)
	return unsigned, token, true
}

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
	"$UNSIGNED":    `"0x" <BASE16>`,
	"$ZERO":        `"0"`,
	"$association": `key ":" value`,
	"$associations": `
    association {"," association} |
    EOL <association EOL> |
    ":"  ! No associations.`,
	"$collection": `"[" (values | associations) "]" "(" CONTEXT ")"`,
	"$key":        `primitive`,
	"$primitive":  `BOOLEAN | COMPLEX | FLOAT | INTEGER | NIL | RUNE | STRING`,
	"$source":     `collection EOF  ! EOF is the end-of-file marker.`,
	"$value":      `primitive | collection`,
	"$values": `
    value {"," value} |
    EOL <value EOL> |
    ! No values.`,
}

// PRIVATE FUNCTIONS

func generateGrammar(expected string, symbols ...string) string {
	var message = "Was expecting '" + expected + "' from:\n"
	for _, symbol := range symbols {
		message += fmt.Sprintf("  \033[32m%v: \033[33m%v\033[0m\n\n", symbol, grammar[symbol])
	}
	return message
}
