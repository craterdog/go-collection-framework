/*******************************************************************************
 *   Copyright (c) 2009-2024 Crater Dog Technologies™.  All Rights Reserved.   *
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

// CLASS NAMESPACE

// This private type defines the namespace structure associated with the
// constants, constructors and functions for the Parser class namespace.
type parserClass_ struct {
	channelSize int
	stackSize   int
}

// Specific Namespace

// This private constant defines the singleton reference to the Parser
// class namespace.  It also initializes any class constants as needed.
var parserClassSingleton = &parserClass_{
	channelSize: 128,
	stackSize:   4,
}

// This public function returns the singleton reference to the Parser
// class namespace.
func Parser() *parserClass_ {
	return parserClassSingleton
}

// CLASS CONSTRUCTORS

// This public class constructor creates a new Parser from the specified
// source bytes.
func (c *parserClass_) FromSource(source []byte) ParserLike {
	var tokens = make(chan TokenLike, c.channelSize)
	Scanner().FromSource(source, tokens) // Starts scanning in a separate go routine.
	var parser = &parser_{
		next:   Stack[TokenLike]().WithCapacity(c.stackSize),
		source: source,
		tokens: tokens,
	}
	return parser
}

// CLASS FUNCTIONS

// This public class function parses the specified collection source retrieved
// from a POSIX compliant file and returns the corresponding collection that was
// used to generate the collection source using the collection formatting
// capabilities.  A POSIX compliant file must end with an EOF marker.
func (c *parserClass_) ParseCollection(source []byte) Collection {
	var parser = c.FromSource(source)
	var collection = parser.ParseCollection()
	return collection
}

// CLASS TYPE

// Encapsulated Type

// This private class type encapsulates a Go structure containing private
// attributes that can only be accessed and manipulated using methods that
// implement the parser-like abstract type.
type parser_ struct {
	next   StackLike[TokenLike] // A stack of the retrieved tokens still to be read.
	source []byte               // The source bytes to be parsed.
	tokens chan TokenLike       // A queue of unread tokens coming from the scanner.
}

// Lexical Interface

// This public class method uses the Parser to parse a collection from the
// current source bytes.
func (v *parser_) ParseCollection() Collection {
	var collection, token, ok = v.parseCollection()
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("collection",
			"$source",
			"$collection",
		)
		panic(message)
	}
	_, token, ok = v.parseEOF()
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("EOF",
			"$source",
			"$collection",
		)
		panic(message)
	}
	return collection
}

// Private Interface

// This private class method returns an error message containing the context for
// a parsing error.
func (v *parser_) formatError(token TokenLike) string {
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

// This private class method is useful when creating Scanner and Parser error
// messages.
func (v *parser_) generateGrammar(expected string, symbols ...string) string {
	var message = "Was expecting '" + expected + "' from:\n"
	for _, symbol := range symbols {
		message += fmt.Sprintf("  \033[32m%v: \033[33m%v\033[0m\n\n", symbol, grammar[symbol])
	}
	return message
}

// This private class method attempts to read the next token from the token
// stream and return it.
func (v *parser_) nextToken() TokenLike {
	var next TokenLike
	if v.next.IsEmpty() {
		var token, ok = <-v.tokens
		if !ok {
			panic("The token channel terminated without an EOF or error token.")
		}
		next = token
		if next.GetType() == Token().TypeError() {
			var message = v.formatError(next)
			panic(message)
		}
	} else {
		next = v.next.RemoveTop()
	}
	return next
}

// This private class method attempts to parse an association between a key and
// value. It returns the association and whether or not the association was
// successfully parsed.
func (v *parser_) parseAssociation() (AssociationLike[Key, Value], TokenLike, bool) {
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
	_, _, ok = v.parseDelimiter(":")
	if !ok {
		// This is not an association.
		v.putBack(token) // Put back the primitive key token.
		return association, token, false
	}
	// This must be an association.
	value, token, ok = v.parseValue()
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("value",
			"$association",
			"$key",
			"$value")
		panic(message)
	}
	association = Association[Key, Value]().FromPair(key, value)
	return association, token, true
}

// This private class method attempts to parse a sequence of associations. It
// returns the sequence of associations and whether or not the sequence of
// associations was successfully parsed.
func (v *parser_) parseAssociations() (Sequential[Binding[Key, Value]], TokenLike, bool) {
	var ok bool
	var token TokenLike
	var associations Sequential[Binding[Key, Value]]
	_, token, ok = v.parseDelimiter(":")
	if ok {
		// The associations is empty.
		associations = Catalog[Key, Value]().Empty()
		return associations, token, true
	}
	_, token, ok = v.parseEOL()
	if ok {
		associations, _, ok = v.parseMultilineAssociations()
		if !ok {
			v.putBack(token) // Put back the EOL character.
			return associations, token, false
		}
	} else {
		associations, token, ok = v.parseInlineAssociations()
		if !ok {
			return associations, token, false
		}
	}
	return associations, token, true
}

// This private class method attempts to parse a boolean primitive. It returns
// the boolean primitive and whether or not the boolean primitive was
// successfully parsed.
func (v *parser_) parseBoolean() (bool, TokenLike, bool) {
	var boolean bool
	var token = v.nextToken()
	if token.GetType() != Token().TypeBoolean() {
		v.putBack(token)
		return boolean, token, false
	}
	boolean, _ = stc.ParseBool(token.GetValue())
	return boolean, token, true
}

// This private class method attempts to parse a collection of values. It
// returns the collection and whether or not the collection was successfully
// parsed.
func (v *parser_) parseCollection() (Collection, TokenLike, bool) {
	var ok bool
	var token TokenLike
	var context string
	var collection Collection
	_, token, ok = v.parseDelimiter("[")
	if !ok {
		// This is not a collection.
		return collection, token, false
	}
	collection, _, ok = v.parseAssociations()
	if !ok {
		// The values must be attempted second since it may start with a component
		// which cannot be put back as a single token.
		collection, _, ok = v.parseValues()
		if !ok {
			// Not a collection (must be a range).
			v.putBack(token) // Put back the "[" character.
			return collection, token, false
		}
	}
	_, token, ok = v.parseDelimiter("]")
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("]",
			"$collection",
			"$values",
			"$associations",
		)
		panic(message)
	}
	_, token, ok = v.parseDelimiter("(")
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("(",
			"$collection",
			"$values",
			"$associations",
		)
		panic(message)
	}
	context, token, ok = v.parseContext()
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("CONTEXT",
			"$collection",
			"$values",
			"$associations",
		)
		panic(message)
	}
	switch sequence := collection.(type) {
	case Sequential[Value]:
		switch context {
		case "array":
			collection = sequence.AsArray()
		case "Array":
			collection = Array[Value]().FromArray(sequence.AsArray())
		case "List":
			collection = List[Value]().FromSequence(sequence)
		case "Queue":
			collection = Queue[Value]().FromSequence(sequence)
		case "Set":
			collection = Set[Value]().FromSequence(sequence)
		case "Stack":
			collection = Stack[Value]().FromSequence(sequence)
		default:
		}
	case Sequential[Binding[Key, Value]]:
		switch context {
		case "map":
			var map_ = map[Key]Value{}
			var iterator = sequence.GetIterator()
			for iterator.HasNext() {
				var association = iterator.GetNext()
				var key = association.GetKey()
				var value = association.GetValue()
				map_[key] = value
			}
			collection = map_
		case "Map":
			collection = Map[Key, Value]().FromArray(sequence.AsArray())
		case "Catalog":
			collection = Catalog[Key, Value]().FromSequence(sequence)
		}
	}
	_, token, ok = v.parseDelimiter(")")
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar(")",
			"$collection",
			"$context")
		panic(message)
	}
	return collection, token, true
}

// This private class method attempts to parse a complex number primitive. It
// returns the complex number primitive and whether or not the complex number
// primitive was successfully parsed.
func (v *parser_) parseComplex() (complex128, TokenLike, bool) {
	var complex_ complex128
	var token = v.nextToken()
	if token.GetType() != Token().TypeComplex() {
		v.putBack(token)
		return complex_, token, false
	}
	complex_, _ = stc.ParseComplex(token.GetValue(), 128)
	return complex_, token, true
}

// This private class method attempts to parse the type of a collection. It
// returns the type string and whether or not the type string was successfully
// parsed.
func (v *parser_) parseContext() (string, TokenLike, bool) {
	var token = v.nextToken()
	if token.GetType() != Token().TypeContext() {
		v.putBack(token)
		return "", token, false
	}
	return token.GetValue(), token, true
}

// This private class method attempts to parse the specified delimiter. It
// returns the token and whether or not the delimiter was found.
func (v *parser_) parseDelimiter(delimiter string) (string, TokenLike, bool) {
	var token = v.nextToken()
	if token.GetType() == Token().TypeEOF() || token.GetValue() != delimiter {
		v.putBack(token)
		return delimiter, token, false
	}
	return delimiter, token, true
}

// This private class method attempts to parse the end-of-file (EOF) marker. It
// returns the token and whether or not an EOF token was found.
func (v *parser_) parseEOF() (TokenLike, TokenLike, bool) {
	var token = v.nextToken()
	if token.GetType() != Token().TypeEOF() {
		v.putBack(token)
		return token, token, false
	}
	return token, token, true
}

// This private class method attempts to parse the end-of-line (EOL) marker. It
// returns the token and whether or not an EOL token was found.
func (v *parser_) parseEOL() (TokenLike, TokenLike, bool) {
	var token = v.nextToken()
	if token.GetType() != Token().TypeEOL() {
		v.putBack(token)
		return token, token, false
	}
	return token, token, true
}

// This private class method attempts to parse a floating point primitive. It
// returns the floating point primitive and whether or not the floating point
// primitive was successfully parsed.
func (v *parser_) parseFloat() (float64, TokenLike, bool) {
	var float float64
	var token = v.nextToken()
	if token.GetType() != Token().TypeFloat() {
		v.putBack(token)
		return float, token, false
	}
	float, _ = stc.ParseFloat(token.GetValue(), 64)
	return float, token, true
}

// This private class method attempts to parse a sequence containing inline
// associations. It returns a sequence of associations and whether or not the
// sequence of associations was successfully parsed.
func (v *parser_) parseInlineAssociations() (Sequential[Binding[Key, Value]], TokenLike, bool) {
	var ok bool
	var token TokenLike
	var association AssociationLike[Key, Value]
	var associations = Catalog[Key, Value]().Empty()
	association, token, ok = v.parseAssociation()
	if !ok {
		// This is not an inline associations.
		return associations, token, false
	}
	for {
		var key = association.GetKey()
		var value = association.GetValue()
		associations.SetValue(key, value)
		// Every subsequent association must be preceded by a ','.
		_, token, ok = v.parseDelimiter(",")
		if !ok {
			// No more associations.
			return associations, token, true
		}
		association, token, ok = v.parseAssociation()
		if !ok {
			var message = v.formatError(token)
			message += v.generateGrammar("association",
				"$collection",
				"$associations",
				"$association",
			)
			panic(message)
		}
	}
}

// This private class method attempts to parse a sequence containing inline
// values. It returns the sequence of values and whether or not the sequence of
// values was successfully parsed.
func (v *parser_) parseInlineValues() (Sequential[Value], TokenLike, bool) {
	var ok bool
	var token TokenLike
	var value Value
	var values = List[Value]().Empty()
	value, token, ok = v.parseValue()
	if !ok {
		// This is not an inline values.
		return values, token, false
	}
	for {
		values.AppendValue(value)
		// Every subsequent value must be preceded by a ','.
		_, token, ok = v.parseDelimiter(",")
		if !ok {
			// No more values.
			return values, token, true
		}
		value, token, ok = v.parseValue()
		if !ok {
			var message = v.formatError(token)
			message += v.generateGrammar("value",
				"$collection",
				"$values",
				"$value",
			)
			panic(message)
		}
	}
}

// This private class method attempts to parse a integer primitive. It returns
// the integer primitive and whether or not the integer primitive was
// successfully parsed.
func (v *parser_) parseInteger() (int64, TokenLike, bool) {
	var integer int64
	var token = v.nextToken()
	if token.GetType() != Token().TypeInteger() {
		v.putBack(token)
		return integer, token, false
	}
	integer, _ = stc.ParseInt(token.GetValue(), 10, 64)
	return integer, token, true
}

// This private class method attempts to parse a sequence containing multi-line
// associations.  It returns the sequence of associations and whether or not the
// sequence of associations was successfully parsed.
func (v *parser_) parseMultilineAssociations() (Sequential[Binding[Key, Value]], TokenLike, bool) {
	var ok bool
	var token TokenLike
	var association AssociationLike[Key, Value]
	var associations = Catalog[Key, Value]().Empty()

	association, token, ok = v.parseAssociation()
	if !ok {
		// This is not a multi-line associations.
		return associations, token, false
	}
	// Every association must be followed by an EOL.
	_, token, ok = v.parseEOL()
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("EOL",
			"$collection",
			"$associations",
			"$association",
		)
		panic(message)
	}
	for {
		var key = association.GetKey()
		var value = association.GetValue()
		associations.SetValue(key, value)
		association, token, ok = v.parseAssociation()
		if !ok {
			// No more associations.
			return associations, token, true
		}
		// Every association must be followed by an EOL.
		_, token, ok = v.parseEOL()
		if !ok {
			var message = v.formatError(token)
			message += v.generateGrammar("EOL",
				"$collection",
				"$associations",
				"$association",
			)
			panic(message)
		}
	}
}

// This private class method attempts to parse a sequence containing multi-line
// values. It returns the sequence of values and whether or not the sequence of
// values was successfully parsed.
func (v *parser_) parseMultilineValues() (Sequential[Value], TokenLike, bool) {
	var ok bool
	var token TokenLike
	var value Value
	var values = List[Value]().Empty()
	value, token, ok = v.parseValue()
	if !ok {
		// This is not a multi-line values.
		return values, token, false
	}
	// Every value must be followed by an EOL.
	_, token, ok = v.parseEOL()
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("EOL",
			"$collection",
			"$values",
			"$value",
		)
		panic(message)
	}
	for {
		values.AppendValue(value)
		value, token, ok = v.parseValue()
		if !ok {
			// No more values.
			return values, token, true
		}
		// Every value must be followed by an EOL.
		_, token, ok = v.parseEOL()
		if !ok {
			var message = v.formatError(token)
			message += v.generateGrammar("EOL",
				"$collection",
				"$values",
				"$value",
			)
			panic(message)
		}
	}
}

// This private class method attempts to parse a nil primitive. It returns the
// nil primitive and whether or not the nil primitive was successfully parsed.
func (v *parser_) parseNil() (Value, TokenLike, bool) {
	var token = v.nextToken()
	if token.GetType() != Token().TypeNil() {
		v.putBack(token)
		return nil, token, false
	}
	return nil, token, true
}

// This private class method attempts to parse a primitive. It returns the
// primitive and whether or not the primitive was successfully parsed.
func (v *parser_) parsePrimitive() (Primitive, TokenLike, bool) {
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

// This private class method attempts to parse a rune. It returns the rune and
// whether or not the rune was successfully parsed.
func (v *parser_) parseRune() (rune, TokenLike, bool) {
	var rune_ rune
	var size int
	var token = v.nextToken()
	if token.GetType() != Token().TypeRune() {
		v.putBack(token)
		return rune_, token, false
	}
	var matches = Scanner().MatchRune(token.GetValue())
	// We must unquote the full token string properly.
	var s, _ = stc.Unquote(matches[0])
	rune_, size = utf.DecodeRuneInString(s)
	if len(s) != size {
		// This is not a rune.
		v.putBack(token)
		return rune_, token, false
	}
	return rune_, token, true
}

// This private class method attempts to parse a string primitive. It returns
// the string primitive and whether or not the string primitive was successfully
// parsed.
func (v *parser_) parseString() (string, TokenLike, bool) {
	var string_ string
	var token = v.nextToken()
	if token.GetType() != Token().TypeString() {
		v.putBack(token)
		return string_, token, false
	}
	var matches = Scanner().MatchString(token.GetValue())
	// We must unquote the full token string properly.
	string_, _ = stc.Unquote(matches[0])
	return string_, token, true
}

// This private class method attempts to parse an unsigned integer primitive. It
// returns the unsigned integer primitive and whether or not the unsigned
// integer primitive was successfully parsed.
func (v *parser_) parseUnsigned() (uint64, TokenLike, bool) {
	var unsigned uint64
	var token = v.nextToken()
	if token.GetType() != Token().TypeUnsigned() {
		v.putBack(token)
		return unsigned, token, false
	}
	unsigned, _ = stc.ParseUint(token.GetValue()[2:], 16, 64)
	return unsigned, token, true
}

// This private class method attempts to parse a component entity. It returns
// the component entity and whether or not the component entity was successfully
// parsed.
func (v *parser_) parseValue() (Value, TokenLike, bool) {
	var ok bool
	var token TokenLike
	var value Value
	value, token, ok = v.parsePrimitive()
	if !ok {
		value, token, ok = v.parseCollection()
	}
	return value, token, ok
}

// This private class method attempts to parse a sequence of values. It returns
// the sequence of values and whether or not the sequence of values was
// successfully parsed.
func (v *parser_) parseValues() (Sequential[Value], TokenLike, bool) {
	var ok bool
	var token TokenLike
	var values Sequential[Value]
	_, token, ok = v.parseDelimiter("]")
	if ok {
		// The values is empty.
		v.putBack(token) // Put back the "]" character.
		values = List[Value]().Empty()
		return values, token, true
	}
	_, token, ok = v.parseEOL()
	if ok {
		values, _, ok = v.parseMultilineValues()
		if !ok {
			v.putBack(token) // Put back the EOL character.
			return values, token, false
		}
	} else {
		values, token, ok = v.parseInlineValues()
		if !ok {
			return values, token, false
		}
	}
	return values, token, true
}

// This private class method puts back the current token onto the token stream
// so that it can be retrieved by another parsing method.
func (v *parser_) putBack(token TokenLike) {
	v.next.AddValue(token)
}

// This map captures the syntax rules for collections of Go primitives.
var grammar = map[string]string{
	"$association": `key ":" value`,
	"$associations": `
      association ("," association)*
    | EOL (association EOL)+
    | ":"  ! No associations.`,
	"$collection": `"[" (values | associations) "]" "(" CONTEXT ")"`,
	"$key":        `primitive`,
	"$primitive":  `BOOLEAN | COMPLEX | FLOAT | INTEGER | NIL | RUNE | STRING`,
	"$source":     `collection EOF  ! EOF is the end-of-file marker.`,
	"$value":      `primitive | collection`,
	"$values": `
      value ("," value)*
    | EOL (value EOL)+
    | " "  ! No values.`,
}
