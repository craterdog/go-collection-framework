/*
................................................................................
.    Copyright (c) 2009-2024 Crater Dog Technologies.  All Rights Reserved.    .
................................................................................
.  DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               .
.                                                                              .
.  This code is free software; you can redistribute it and/or modify it under  .
.  the terms of The MIT License (MIT), as published by the Open Source         .
.  Initiative. (See http://opensource.org/licenses/MIT)                        .
................................................................................
*/

package cdcn

import (
	fmt "fmt"
	col "github.com/craterdog/go-collection-framework/v3"
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
		tokens_: col.Queue[TokenLike]().MakeWithCapacity(c.queueSize_),
		next_:   col.Stack[TokenLike]().MakeWithCapacity(c.stackSize_),
	}
}

// INSTANCE METHODS

// Target

type parser_ struct {
	source_ string                   // The original source code.
	tokens_ col.QueueLike[TokenLike] // A queue of unread tokens from the scanner.
	next_   col.StackLike[TokenLike] // A stack of read, but unprocessed tokens.
}

// Public

func (v *parser_) ParseSource(source string) col.Collection {
	// The scanner runs in a separate Go routine.
	v.source_ = source
	Scanner().Make(v.source_, v.tokens_)

	// Attempt to parse a collection.
	var collection, token, ok = v.parseCollection()
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("collection",
			"source",
			"collection",
		)
		panic(message)
	}

	// Attempt to parse the end-of-file marker.
	_, token, ok = v.parseToken(EOFToken, "")
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("EOF",
			"source",
			"collection",
		)
		panic(message)
	}

	// Found a collection.
	return collection
}

// Private

/*
This private instance method returns an error message containing the context for
a parsing error.
*/
func (v *parser_) formatError(token TokenLike) string {
	// Format the error message.
	var message = fmt.Sprintf(
		"An unexpected token was received by the parser: %v\n",
		token,
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

/*
This private instance method is useful when creating scanner and parser error
messages that include the required grammatical rules.
*/
func (v *parser_) generateGrammar(expected string, names ...string) string {
	var message = "Was expecting '" + expected + "' from:\n"
	for _, name := range names {
		message += fmt.Sprintf(
			"  \033[32m%v: \033[33m%v\033[0m\n\n",
			name,
			grammar[name],
		)
	}
	return message
}

/*
This private instance method attempts to read the next token from the token
stream and return it.
*/
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
	association col.AssociationLike[col.Key, col.Value],
	token TokenLike,
	ok bool,
) {
	// Attempt to parse a primitive key.
	var key col.Key
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
	var value col.Value
	value, token, ok = v.parseValue()
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("value",
			"association",
			"key",
			"value")
		panic(message)
	}

	// Found an association.
	association = col.Association[col.Key, col.Value]().MakeWithAttributes(key, value)
	return association, token, true
}

func (v *parser_) parseAssociations() (
	associations col.CatalogLike[col.Key, col.Value],
	token TokenLike,
	ok bool,
) {
	var association col.AssociationLike[col.Key, col.Value]
	associations = col.Catalog[col.Key, col.Value]().Make()

	// Check for an empty sequence.
	_, token, ok = v.parseToken(DelimiterToken, ":")
	if ok {
		// Found an empty sequence.
		return associations, token, true
	}

	// Attempt to parse a multi-line sequence of associations.
	_, token, ok = v.parseToken(EOLToken, "")
	if ok {
		// Attempt to parse one or more associations.
		association, _, ok = v.parseAssociation()
		if !ok {
			// This must be a sequence of values instead.
			v.putBack(token)
			return associations, token, false
		}
		for ok {
			var key = association.GetKey()
			var value = association.GetValue()
			associations.SetValue(key, value)
			_, token, ok = v.parseToken(EOLToken, "")
			if !ok {
				break
			}
			association, _, ok = v.parseAssociation()
			if !ok {
				v.putBack(token)
				break
			}
		}

		// Attempt to parse an optional end-of-line character.
		_, token, _ = v.parseToken(EOLToken, "")

		// Found a multi-line sequence of associations.
		return associations, token, true
	}

	// Attempt to parse one or more associations in an in-line sequence.
	association, token, ok = v.parseAssociation()
	if !ok {
		return associations, token, false
	}
	for ok {
		var value = association.GetValue()
		var key = association.GetKey()
		associations.SetValue(key, value)
		_, token, ok = v.parseToken(DelimiterToken, ",")
		if ok {
			// Attempt to parse a association.
			association, token, ok = v.parseAssociation()
			if !ok {
				var message = v.formatError(token)
				message += v.generateGrammar("association",
					"associations",
					"association",
				)
				panic(message)
			}
		}
	}

	// Found an in-line sequence of associations.
	return associations, token, true
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
			message += v.generateGrammar("associations",
				"collection",
				"associations",
				"values",
			)
			panic(message)
		}
	}

	// Attempt to parse the closing bracket of the collection.
	_, token, ok = v.parseToken(DelimiterToken, "]")
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("]",
			"collection",
			"associations",
			"values",
		)
		panic(message)
	}

	// Attempt to parse the opening bracket of the context.
	_, token, ok = v.parseToken(DelimiterToken, "(")
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("(",
			"collection",
			"associations",
			"values",
		)
		panic(message)
	}

	// Attempt to parse the context for the collection.
	var context string
	context, token, ok = v.parseToken(ContextToken, "")
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("CONTEXT",
			"collection",
			"associations",
			"values",
		)
		panic(message)
	}

	// Attempt to parse the opening bracket of the context.
	_, token, ok = v.parseToken(DelimiterToken, ")")
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar(")",
			"collection",
			"associations",
			"values",
		)
		panic(message)
	}

	// Found a collection of a specific type.
	switch sequence := collection.(type) {
	case col.Sequential[col.Value]:
		switch context {
		case "array":
			collection = sequence.AsArray()
		case "Array":
			collection = col.Array[col.Value]().MakeFromArray(sequence.AsArray())
		case "List":
			collection = col.List[col.Value]().MakeFromSequence(sequence)
		case "Queue":
			collection = col.Queue[col.Value]().MakeFromSequence(sequence)
		case "Set":
			collection = col.Set[col.Value]().MakeFromSequence(sequence)
		case "Stack":
			collection = col.Stack[col.Value]().MakeFromSequence(sequence)
		default:
			var message = fmt.Sprintf("Found an unknown collection type: %q", context)
			panic(message)
		}
	case col.Sequential[col.AssociationLike[col.Key, col.Value]]:
		switch context {
		case "map":
			var map_ = map[col.Key]col.Value{}
			var iterator = sequence.GetIterator()
			for iterator.HasNext() {
				var association = iterator.GetNext()
				var key = association.GetKey()
				var value = association.GetValue()
				map_[key] = value
			}
			collection = map_
		case "Map":
			collection = col.Map[col.Key, col.Value]().MakeFromArray(sequence.AsArray())
		case "Catalog":
			collection = col.Catalog[col.Key, col.Value]().MakeFromSequence(sequence)
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

func (v *parser_) parseKey() (
	key col.Key,
	token TokenLike,
	ok bool,
) {
	// Attempt to parse a primitive.
	var primitive any
	primitive, token, ok = v.parsePrimitive()
	if !ok {
		return key, token, false
	}

	// Found a primitive key.
	key = col.Key(primitive)
	return key, token, true
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
		if len(expectedValue) == 0 || expectedValue == value {
			// Found the expected token.
			return value, token, true
		}
	}

	// This is not the right token.
	v.putBack(token)
	return value, token, false
}

func (v *parser_) parseValue() (
	value col.Value,
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
	values col.ListLike[col.Value],
	token TokenLike,
	ok bool,
) {
	var value col.Value
	values = col.List[col.Value]().Make()

	// Check for an empty sequence.
	_, token, ok = v.parseToken(DelimiterToken, "]")
	if ok {
		v.putBack(token)
		return values, token, true
	}

	// Attempt to parse a multi-line sequence of values.
	_, _, ok = v.parseToken(EOLToken, "")
	if ok {
		// Attempt to parse the first value.
		value, token, ok = v.parseValue()
		if !ok {
			var message = v.formatError(token)
			message += v.generateGrammar("value",
				"values",
				"value",
			)
			panic(message)
		}

		// Parse any additional values.
		for ok {
			values.AppendValue(value)
			_, token, ok = v.parseToken(EOLToken, "")
			if !ok {
				break
			}
			value, _, ok = v.parseValue()
			if !ok {
				v.putBack(token)
				break
			}
		}

		// Attempt to parse an optional end-of-line character.
		_, token, _ = v.parseToken(EOLToken, "")

		// Found a multi-line sequence of values.
		return values, token, true
	}

	// Attempt to parse the first value in an in-line sequence.
	value, token, ok = v.parseValue()
	if !ok {
		return values, token, false
	}

	// Parse any additional values.
	for ok {
		values.AppendValue(value)
		_, token, ok = v.parseToken(DelimiterToken, ",")
		if ok {
			// Attempt to parse a value.
			value, token, ok = v.parseValue()
			if !ok {
				var message = v.formatError(token)
				message += v.generateGrammar("value",
					"values",
					"value",
				)
				panic(message)
			}
		}
	}

	// Found an in-line sequence of values.
	return values, token, true
}

func (v *parser_) putBack(token TokenLike) {
	v.next_.AddValue(token)
}

/*
This Go map captures the syntax rules for collections of Go primitives.
*/
var grammar = map[string]string{
	"association": `key ":" value`,
	"associations": `
    association ("," association)*
    (EOL association)+ EOL?
    ":"  ! No associations.`,
	"collection": `"[" (associations | values) "]" "(" Context ")"`,
	"key":        `primitive`,
	"primitive":  `Boolean | Complex | Float | Hexadecimal | Integer | Nil | Rune | String`,
	"source":     `collection EOF  ! Terminated with an end-of-file marker.`,
	"value":      `collection | primitive`,
	"values": `
    value ("," value)*
    (EOL value)+ EOL?
    " "  ! No values.`,
}
