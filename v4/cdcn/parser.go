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

package cdcn

import (
	fmt "fmt"
	col "github.com/craterdog/go-collection-framework/v4/collection"
	stc "strconv"
	sts "strings"
	utf "unicode/utf8"
)

// CLASS ACCESS

// Reference

var parserClass = &parserClass_{
	queueSize_: 16,
	stackSize_: 4,
}

// Function

func Parser() ParserClassLike {
	return parserClass
}

// CLASS METHODS

// Target

type parserClass_ struct {
	queueSize_ int
	stackSize_ int
}

// Constructors

func (c *parserClass_) Make() ParserLike {
	return &parser_{
		class_: c,
	}
}

// INSTANCE METHODS

// Target

type parser_ struct {
	class_    ParserClassLike
	notation_ col.NotationLike
	source_   string                   // The original source code.
	tokens_   col.QueueLike[TokenLike] // A queue of unread tokens from the scanner.
	next_     col.StackLike[TokenLike] // A stack of read, but unprocessed tokens.
}

// Attributes

func (v *parser_) GetClass() ParserClassLike {
	return v.class_
}

// Public

func (v *parser_) ParseSource(source string) col.Collection {
	v.source_ = source
	v.notation_ = Notation().Make()
	v.tokens_ = col.Queue[TokenLike](v.notation_).MakeWithCapacity(parserClass.queueSize_)
	v.next_ = col.Stack[TokenLike](v.notation_).MakeWithCapacity(parserClass.stackSize_)

	// The scanner runs in a separate Go routine.
	Scanner().Make(v.source_, v.tokens_)

	// Attempt to parse a collection.
	var collection, token, ok = v.parseCollection()
	if !ok {
		var message = v.formatError(token)
		message += v.generateSyntax("Collection",
			"Cdcn",
			"Collection",
		)
		panic(message)
	}

	// Attempt to parse optional end-of-line characters.
	for ok {
		_, _, ok = v.parseToken(EOLToken, "")
	}

	// Attempt to parse the end-of-file marker.
	_, token, ok = v.parseToken(EOFToken, "")
	if !ok {
		var message = v.formatError(token)
		message += v.generateSyntax("EOF",
			"Cdcn",
			"Collection",
		)
		panic(message)
	}

	// Found a collection.
	return collection
}

// Private

func (v *parser_) formatError(token TokenLike) string {
	// Format the error message.
	var message = fmt.Sprintf(
		"An unexpected token was received by the parser: %v\n",
		Scanner().FormatToken(token),
	)
	var line = token.GetLine()
	var lines = sts.Split(v.source_, "\n")

	// Append the source line with the error in it.
	message += "\033[36m"
	if line > 1 {
		message += fmt.Sprintf("%04d: ", line-1) + string(lines[line-2]) + "\n"
	}
	message += fmt.Sprintf("%04d: ", line) + string(lines[line-1]) + "\n"

	// Append an arrow pointing to the error.
	message += " \033[32m>>>─"
	var count = 0
	for count < token.GetPosition() {
		message += "─"
		count++
	}
	message += "⌃\033[36m\n"

	// Append the following source line for context.
	if line < len(lines) {
		message += fmt.Sprintf("%04d: ", line+1) + string(lines[line]) + "\n"
	}
	message += "\033[0m\n"

	return message
}

func (v *parser_) generateSyntax(expected string, names ...string) string {
	var message = "Was expecting '" + expected + "' from:\n"
	for _, name := range names {
		message += fmt.Sprintf(
			"  \033[32m%v: \033[33m%v\033[0m\n\n",
			name,
			syntax[name],
		)
	}
	return message
}

func (v *parser_) getNextToken() TokenLike {
	// Check for any read, but unprocessed tokens.
	if !v.next_.IsEmpty() {
		return v.next_.RemoveTop()
	}

	// Read a new token from the token stream.
	var token, ok = v.tokens_.RemoveHead() // This will wait for a token.
	if !ok {
		panic("The token channel terminated without an EOF token.")
	}

	// Check for an error token.
	if token.GetType() == ErrorToken {
		var message = v.formatError(token)
		panic(message)
	}

	return token
}

func (v *parser_) parseAssociation() (
	association col.AssociationLike[any, any],
	token TokenLike,
	ok bool,
) {
	// Attempt to parse a primitive key.
	var key any
	key, token, ok = v.parseKey()
	if !ok {
		return association, token, false
	}
	_, _, ok = v.parseToken(DelimiterToken, ":")
	if !ok {
		// The primitive token is not a key.
		v.putBack(token)
		return association, token, false
	}

	// Attempt to parse a value.
	var value any
	value, token, ok = v.parseValue()
	if !ok {
		var message = v.formatError(token)
		message += v.generateSyntax("Value",
			"Association",
			"Key",
			"Value")
		panic(message)
	}

	// Found an association.
	association = col.Association[any, any]().MakeWithAttributes(key, value)
	return association, token, true
}

func (v *parser_) parseAssociations() (
	associations col.CatalogLike[any, any],
	token TokenLike,
	ok bool,
) {
	// Check for an empty sequence of associations.
	_, token, ok = v.parseToken(DelimiterToken, ":")
	if ok {
		associations = col.Catalog[any, any](v.notation_).Make()
		return associations, token, true
	}

	// Attempt to parse an inline sequence of associations.
	associations, token, ok = v.parseInlineAssociations()
	if ok {
		return associations, token, true
	}

	// Attempt to parse an multi-line sequence of associations.
	associations, token, ok = v.parseMultilineAssociations()
	if ok {
		return associations, token, true
	}

	// This is not a sequence of associations.
	return associations, token, false
}

func (v *parser_) parseCollection() (
	collection col.Collection,
	token TokenLike,
	ok bool,
) {
	// Attempt to parse the opening bracket of the collection.
	_, token, ok = v.parseToken(DelimiterToken, "[")
	if !ok {
		return collection, token, false
	}

	// Attempt to parse a sequence of associations.
	collection, _, ok = v.parseAssociations()
	if !ok {
		// Attempt to parse a sequence of values. The values must be
		// attempted second since it may start with a component which
		// cannot be put back as a single token.
		collection, _, ok = v.parseValues()
		if !ok {
			var message = v.formatError(token)
			message += v.generateSyntax("Associations",
				"Collection",
				"Associations",
				"Values",
			)
			panic(message)
		}
	}

	// Attempt to parse the closing bracket of the collection.
	_, token, ok = v.parseToken(DelimiterToken, "]")
	if !ok {
		var message = v.formatError(token)
		message += v.generateSyntax("]",
			"Collection",
			"Associations",
			"Values",
		)
		panic(message)
	}

	// Attempt to parse the opening bracket of the context.
	_, token, ok = v.parseToken(DelimiterToken, "(")
	if !ok {
		var message = v.formatError(token)
		message += v.generateSyntax("(",
			"Collection",
			"Associations",
			"Values",
		)
		panic(message)
	}

	// Attempt to parse the context for the collection.
	var context string
	context, token, ok = v.parseToken(ContextToken, "")
	if !ok {
		var message = v.formatError(token)
		message += v.generateSyntax("Context",
			"Collection",
			"Associations",
			"Values",
		)
		panic(message)
	}

	// Attempt to parse the closing bracket of the context.
	_, token, ok = v.parseToken(DelimiterToken, ")")
	if !ok {
		var message = v.formatError(token)
		message += v.generateSyntax(")",
			"Collection",
			"Associations",
			"Values",
		)
		panic(message)
	}

	// Found a collection of a specific type.
	switch sequence := collection.(type) {
	case col.Sequential[any]:
		switch context {
		case "array":
			collection = sequence.AsArray()
		case "Array":
			collection = col.Array[any](v.notation_).MakeFromSequence(sequence)
		case "List":
			collection = col.List[any](v.notation_).MakeFromSequence(sequence)
		case "Queue":
			collection = col.Queue[any](v.notation_).MakeFromSequence(sequence)
		case "Set":
			collection = col.Set[any](v.notation_).MakeFromSequence(sequence)
		case "Stack":
			collection = col.Stack[any](v.notation_).MakeFromSequence(sequence)
		default:
			var message = fmt.Sprintf("Found an unknown collection type: %q", context)
			panic(message)
		}
	case col.Sequential[col.AssociationLike[any, any]]:
		switch context {
		case "map":
			var map_ = map[any]any{}
			var iterator = sequence.GetIterator()
			for iterator.HasNext() {
				var association = iterator.GetNext()
				var key = association.GetKey()
				var value = association.GetValue()
				map_[key] = value
			}
			collection = map_
		case "Map":
			collection = col.Map[any, any](v.notation_).MakeFromSequence(sequence)
		case "Catalog":
			collection = col.Catalog[any, any](v.notation_).MakeFromSequence(sequence)
		default:
			var message = fmt.Sprintf("Found an unknown collection type: %q", context)
			panic(message)
		}
	default:
		var message = fmt.Sprintf("Found an unknown sequence type: %T", sequence)
		panic(message)
	}
	return collection, token, true
}

func (v *parser_) parseInlineAssociations() (
	associations col.CatalogLike[any, any],
	token TokenLike,
	ok bool,
) {
	// Attempt to parse one or more inline associations.
	var association col.AssociationLike[any, any]
	association, token, ok = v.parseAssociation()
	if !ok {
		// This is not an inline sequence of associations.
		return associations, token, false
	}
	associations = col.Catalog[any, any](v.notation_).Make()
	for ok {
		var value = association.GetValue()
		var key = association.GetKey()
		associations.SetValue(key, value)
		_, token, ok = v.parseToken(DelimiterToken, ",")
		if ok {
			association, token, ok = v.parseAssociation()
			if !ok {
				var message = v.formatError(token)
				message += v.generateSyntax("Association",
					"Associations",
					"Association",
				)
				panic(message)
			}
		}
	}

	// Found an inline sequence of associations.
	return associations, token, true
}

func (v *parser_) parseInlineValues() (
	values col.ListLike[any],
	token TokenLike,
	ok bool,
) {
	// Attempt to parse one or more inline values.
	var value any
	value, token, ok = v.parseValue()
	if !ok {
		return values, token, false
	}
	values = col.List[any](v.notation_).Make()
	for ok {
		values.AppendValue(value)
		_, token, ok = v.parseToken(DelimiterToken, ",")
		if ok {
			value, token, ok = v.parseValue()
			if !ok {
				var message = v.formatError(token)
				message += v.generateSyntax("Value",
					"Values",
					"Value",
				)
				panic(message)
			}
		}
	}

	// Found an inline sequence of values.
	return values, token, true
}

func (v *parser_) parseKey() (
	key any,
	token TokenLike,
	ok bool,
) {
	// Attempt to parse a primitive.
	key, token, ok = v.parsePrimitive()
	if !ok {
		return key, token, false
	}

	// Found a primitive key.
	return key, token, true
}

func (v *parser_) parseMultilineAssociations() (
	associations col.CatalogLike[any, any],
	token TokenLike,
	ok bool,
) {
	// Attempt to parse one or more multi-line associations.
	var eolToken TokenLike
	_, eolToken, ok = v.parseToken(EOLToken, "")
	if !ok {
		// This is not a multi-line sequence of associations.
		return associations, eolToken, false
	}
	var association col.AssociationLike[any, any]
	association, token, ok = v.parseAssociation()
	if !ok {
		// This must be a sequence of values instead.
		v.putBack(eolToken)
		return associations, token, false
	}
	associations = col.Catalog[any, any](v.notation_).Make()
	for ok {
		var key = association.GetKey()
		var value = association.GetValue()
		associations.SetValue(key, value)
		_, eolToken, ok = v.parseToken(EOLToken, "")
		if !ok {
			var message = v.formatError(eolToken)
			message += v.generateSyntax("EOL",
				"Associations",
				"Association",
			)
			panic(message)
		}
		association, token, ok = v.parseAssociation()
		if !ok {
			break
		}
	}

	// Found a multi-line sequence of associations.
	return associations, token, true
}

func (v *parser_) parseMultilineValues() (
	values col.ListLike[any],
	token TokenLike,
	ok bool,
) {
	// Attempt to parse one or more multi-line values.
	var eolToken TokenLike
	_, eolToken, ok = v.parseToken(EOLToken, "")
	if !ok {
		// This is not a multi-line sequence of values.
		return values, eolToken, false
	}
	var value any
	value, token, ok = v.parseValue()
	if !ok {
		var message = v.formatError(token)
		message += v.generateSyntax("Value",
			"Values",
			"Value",
		)
		panic(message)
	}
	values = col.List[any](v.notation_).Make()
	for ok {
		values.AppendValue(value)
		_, eolToken, ok = v.parseToken(EOLToken, "")
		if !ok {
			var message = v.formatError(eolToken)
			message += v.generateSyntax("EOL",
				"Values",
				"Value",
			)
			panic(message)
		}
		value, token, ok = v.parseValue()
		if !ok {
			break
		}
	}

	// Found a multi-line sequence of values.
	return values, token, true
}

func (v *parser_) parsePrimitive() (
	primitive any,
	token TokenLike,
	ok bool,
) {
	_, token, ok = v.parseToken(BooleanToken, "")
	if ok {
		primitive, _ = stc.ParseBool(token.GetValue())
		return primitive, token, true
	}
	_, token, ok = v.parseToken(ComplexToken, "")
	if ok {
		primitive, _ = stc.ParseComplex(token.GetValue(), 128)
		return primitive, token, true
	}
	_, token, ok = v.parseToken(FloatToken, "")
	if ok {
		primitive, _ = stc.ParseFloat(token.GetValue(), 64)
		return primitive, token, true
	}
	_, token, ok = v.parseToken(HexadecimalToken, "")
	if ok {
		primitive, _ = stc.ParseUint(token.GetValue()[2:], 16, 64)
		return primitive, token, true
	}
	_, token, ok = v.parseToken(IntegerToken, "")
	if ok {
		primitive, _ = stc.ParseInt(token.GetValue(), 10, 64)
		return primitive, token, true
	}
	_, token, ok = v.parseToken(NilToken, "")
	if ok {
		primitive = nil
		return primitive, token, true
	}
	_, token, ok = v.parseToken(RuneToken, "")
	if ok {
		var matches = Scanner().MatchToken(RuneToken, token.GetValue())
		var match, _ = stc.Unquote(matches.GetValue(1))
		primitive, _ = utf.DecodeRuneInString(match)
		return primitive, token, true
	}
	_, token, ok = v.parseToken(StringToken, "")
	if ok {
		var matches = Scanner().MatchToken(StringToken, token.GetValue())
		primitive, _ = stc.Unquote(matches.GetValue(1))
		return primitive, token, true
	}
	return primitive, token, ok
}

func (v *parser_) parseToken(expectedType TokenType, expectedValue string) (
	value string,
	token TokenLike,
	ok bool,
) {
	// Attempt to parse a specific token.
	token = v.getNextToken()
	if token.GetType() == expectedType {
		value = token.GetValue()
		var notConstrained = len(expectedValue) == 0
		if notConstrained || value == expectedValue {
			// Found the right token.
			return value, token, true
		}
	}

	// This is not the right token.
	v.putBack(token)
	return value, token, false
}

func (v *parser_) parseValue() (
	value any,
	token TokenLike,
	ok bool,
) {
	// Attempt to parse a primitive.
	value, token, ok = v.parsePrimitive()
	if ok {
		// Found a primitive value.
		return value, token, true
	}

	// Attempt to parse a collection.
	value, token, ok = v.parseCollection()
	if ok {
		// Found a collection value.
		return value, token, true
	}

	// This is not a value.
	return value, token, false
}

func (v *parser_) parseValues() (
	values col.ListLike[any],
	token TokenLike,
	ok bool,
) {
	// Check for an empty sequence.
	_, token, ok = v.parseToken(DelimiterToken, "]")
	if ok {
		v.putBack(token)
		values = col.List[any](v.notation_).Make()
		return values, token, true
	}

	// Attempt to parse an inline sequence of values.
	values, token, ok = v.parseInlineValues()
	if ok {
		return values, token, true
	}

	// Attempt to parse an multi-line sequence of values.
	values, token, ok = v.parseMultilineValues()
	if ok {
		return values, token, true
	}

	// This is not a sequence of values.
	return values, token, false
}

func (v *parser_) putBack(token TokenLike) {
	v.next_.AddValue(token)
}

var syntax = map[string]string{
	"AST":        `Collection EOL* EOF  ! Terminated with an end-of-file marker.`,
	"Collection": `"[" (Associations | Values) "]" "(" context ")"`,
	"Associations": `
    Association ("," Association)*
    (EOL Association)+ EOL
    ":"  ! No associations.`,
	"Association": `Key ":" Value`,
	"Key":         `Primitive`,
	"Value":       `Primitive | Collection`,
	"Primitive": `
    boolean
    complex
    float
    hexadecimal
    integer
    nil
    rune
    string`,
	"Values": `
    Value ("," Value)*
    (EOL Value)+ EOL
    " "  ! No values.`,
}
