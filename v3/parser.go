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
	var token *Token
	var collection Collection
	var tokens = make(chan Token, 256)
	Scanner(source, tokens) // Starts scanning in a separate go routine.
	var p = &parser{
		source: source,
		next:   StackWithCapacity[*Token](4),
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
	next           StackLike[*Token] // The stack of the retrieved tokens that have been put back.
	tokens         chan Token        // The queue of unread tokens coming from the scanner.
	p1, p2, p3, p4 *Token            // The previous four tokens that have been retrieved.
}

// This method attempts to read the next token from the token stream and return
// it.
func (v *parser) nextToken() *Token {
	var next *Token
	if v.next.IsEmpty() {
		var token, ok = <-v.tokens
		if !ok {
			panic("The token channel terminated without an EOF or error token.")
		}
		next = &token
		if next.Type == TokenError {
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
func (v *parser) formatError(token *Token) string {
	var message = fmt.Sprintf("An unexpected token was received by the parser: %v\n", token)
	var line = token.Line
	var lines = sts.Split(string(v.source), EOL)

	message += "\033[36m"
	if line > 1 {
		message += fmt.Sprintf("%04d: ", line-1) + string(lines[line-2]) + EOL
	}
	message += fmt.Sprintf("%04d: ", line) + string(lines[line-1]) + EOL

	message += " \033[32m>>>─"
	var count = 0
	for count < token.Position {
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
func (v *parser) parseAssociation() (AssociationLike[Key, Value], *Token, bool) {
	var ok bool
	var token *Token
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
func (v *parser) parseCollection() (Collection, *Token, bool) {
	var ok bool
	var token *Token
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
				collection = QueueFromSequence(sequence)
			case "set":
				var Set = Set[Value]()
				collection = Set.FromSequence(sequence)
			case "stack":
				collection = StackFromSequence(sequence)
			default:
			}
		case Sequential[Binding[Key, Value]]:
			switch context {
			case "catalog":
				collection = Catalog[Key, Value]().FromSequence(sequence)
			case "map":
				collection = MapFromSequence(sequence)
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
func (v *parser) parseDelimiter(delimiter string) (string, *Token, bool) {
	var token = v.nextToken()
	if token.Type == TokenEOF || token.Value != delimiter {
		v.backupOne()
		return delimiter, token, false
	}
	return delimiter, token, true
}

// This method attempts to parse the end-of-file (EOF) marker. It returns
// the token and whether or not an EOL token was found.
func (v *parser) parseEOF() (*Token, *Token, bool) {
	var token = v.nextToken()
	if token.Type != TokenEOF {
		v.backupOne()
		return token, token, false
	}
	return token, token, true
}

// This method attempts to parse the end-of-line (EOL) marker. It returns
// the token and whether or not an EOF token was found.
func (v *parser) parseEOL() (*Token, *Token, bool) {
	var token = v.nextToken()
	if token.Type != TokenEOL {
		v.backupOne()
		return token, token, false
	}
	return token, token, true
}

// It returns a sequence of associations and whether or not the sequence of
// associations was successfully parsed.
func (v *parser) parseInlineAssociations() (Sequential[Binding[Key, Value]], *Token, bool) {
	var ok bool
	var token *Token
	var association Binding[Key, Value]
	var List = List[Binding[Key, Value]]()
	var associations = List.FromNothing()
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
func (v *parser) parseInlineValues() (Sequential[Value], *Token, bool) {
	var ok bool
	var token *Token
	var value Value
	var List = List[Value]()
	var values = List.FromNothing()
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
func (v *parser) parseAssociations() (Sequential[Binding[Key, Value]], *Token, bool) {
	var ok bool
	var token *Token
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
func (v *parser) parseMultilineAssociations() (Sequential[Binding[Key, Value]], *Token, bool) {
	var ok bool
	var token *Token
	var association AssociationLike[Key, Value]
	var List = List[Binding[Key, Value]]()
	var associations = List.FromNothing()
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
func (v *parser) parseMultilineValues() (Sequential[Value], *Token, bool) {
	var ok bool
	var token *Token
	var value Value
	var List = List[Value]()
	var values = List.FromNothing()
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
func (v *parser) parseValue() (Value, *Token, bool) {
	var ok bool
	var token *Token
	var value Value
	value, token, ok = v.parsePrimitive()
	if !ok {
		value, token, ok = v.parseCollection()
	}
	return value, token, ok
}

// This method attempts to parse a primitive. It returns the primitive and
// whether or not the primitive was successfully parsed.
func (v *parser) parsePrimitive() (Primitive, *Token, bool) {
	var ok bool
	var token *Token
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
func (v *parser) parseValues() (Sequential[Value], *Token, bool) {
	var ok bool
	var token *Token
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
func (v *parser) parseBoolean() (bool, *Token, bool) {
	var boolean bool
	var token = v.nextToken()
	if token.Type != TokenBoolean {
		v.backupOne()
		return boolean, token, false
	}
	boolean, _ = stc.ParseBool(token.Value)
	return boolean, token, true
}

// This method attempts to parse a complex number primitive. It returns the
// complex number primitive and whether or not the complex number primitive was
// successfully parsed.
func (v *parser) parseComplex() (complex128, *Token, bool) {
	var complex_ complex128
	var token = v.nextToken()
	if token.Type != TokenComplex {
		v.backupOne()
		return complex_, token, false
	}
	complex_, _ = stc.ParseComplex(token.Value, 128)
	return complex_, token, true
}

// This method attempts to parse the type of a collection. It returns the type
// string and whether or not the type string was successfully parsed.
func (v *parser) parseContext() (string, *Token, bool) {
	var token = v.nextToken()
	if token.Type != TokenContext {
		v.backupOne()
		return "", token, false
	}
	return token.Value, token, true
}

// This method attempts to parse a floating point primitive. It returns the
// floating point primitive and whether or not the floating point primitive was
// successfully parsed.
func (v *parser) parseFloat() (float64, *Token, bool) {
	var float float64
	var token = v.nextToken()
	if token.Type != TokenFloat {
		v.backupOne()
		return float, token, false
	}
	float, _ = stc.ParseFloat(token.Value, 64)
	return float, token, true
}

// This method attempts to parse a integer primitive. It returns the integer
// primitive and whether or not the integer primitive was successfully parsed.
func (v *parser) parseInteger() (int64, *Token, bool) {
	var integer int64
	var token = v.nextToken()
	if token.Type != TokenInteger {
		v.backupOne()
		return integer, token, false
	}
	integer, _ = stc.ParseInt(token.Value, 10, 64)
	return integer, token, true
}

// This method attempts to parse a nil primitive. It returns the nil primitive
// and whether or not the nil primitive was successfully parsed.
func (v *parser) parseNil() (Value, *Token, bool) {
	var token = v.nextToken()
	if token.Type != TokenNil {
		v.backupOne()
		return nil, token, false
	}
	return nil, token, true
}

// This method attempts to parse a rune. It returns the rune and whether or not
// the rune was successfully parsed.
func (v *parser) parseRune() (rune, *Token, bool) {
	var rune_ rune
	var size int
	var token = v.nextToken()
	if token.Type != TokenRune {
		v.backupOne()
		return rune_, token, false
	}
	var matches = scanRune([]byte(token.Value))
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
func (v *parser) parseString() (string, *Token, bool) {
	var string_ string
	var token = v.nextToken()
	if token.Type != TokenString {
		v.backupOne()
		return string_, token, false
	}
	var matches = scanString([]byte(token.Value))
	// We must unquote the full token string properly.
	string_, _ = stc.Unquote(matches[0])
	return string_, token, true
}

// This method attempts to parse an unsigned integer primitive. It returns the
// unsigned integer primitive and whether or not the unsigned integer primitive
// was successfully parsed.
func (v *parser) parseUnsigned() (uint64, *Token, bool) {
	var unsigned uint64
	var token = v.nextToken()
	if token.Type != TokenUnsigned {
		v.backupOne()
		return unsigned, token, false
	}
	unsigned, _ = stc.ParseUint(token.Value[2:], 16, 64)
	return unsigned, token, true
}
