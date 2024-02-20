/*******************************************************************************
 *   Copyright (c) 2009-2024 Crater Dog Technologies™.  All Rights Reserved.   *
 *******************************************************************************
 * DO NOT ALTER OR REMOVE COPYRIGHT NOTICES OR THIS FILE HEADER.               *
 *                                                                             *
 * This code is free software; you can redistribute it and/or modify it under  *
 * the terms of The MIT License (MIT), as published by the Open Source         *
 * Initiative. (See http://opensource.org/licenses/MIT)                        *
 *******************************************************************************/

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
	queueSize: 16,
	stackSize: 4,
}

// Function

func Parser() ParserClassLike {
	return parserClass
}

// CLASS METHODS

// Target

type parserClass_ struct {
	queueSize int
	stackSize int
}

// Constructors

func (c *parserClass_) Make() ParserLike {
	var parser = &parser_{
		next: col.Stack[TokenLike]().MakeWithCapacity(c.stackSize),
	}
	return parser
}

// INSTANCE METHODS

// Target

type parser_ struct {
	next   col.StackLike[TokenLike] // A stack of unprocessed retrieved tokens.
	source string                       // The original source code.
	tokens col.QueueLike[TokenLike] // A queue of unread tokens from the scanner.
}

// Public

func (v *parser_) ParseSource(source string) col.Collection {
	// The scanner runs in a separate Go routine.
	v.source = source
	v.tokens = col.Queue[TokenLike]().MakeWithCapacity(parserClass.queueSize)
	Scanner().Make(v.source, v.tokens)

	// Parse the tokens from the scanner.
	var collection, token, ok = v.parseCollection()
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("collection",
			"$source",
			"$collection",
		)
		panic(message)
	}
	_, token, ok = v.parseToken(EOFToken, "")
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

// Private

/*
This private class method returns an error message containing the context for a
parsing error.
*/
func (v *parser_) formatError(token TokenLike) string {
	var message = fmt.Sprintf(
		"An unexpected token was received by the parser: %v\n",
		token,
	)
	var line = token.GetLine()
	var lines = sts.Split(v.source, "\n")

	message += "\033[36m"
	if line > 1 {
		message += fmt.Sprintf("%04d: ", line-1) + string(lines[line-2]) + "\n"
	}
	message += fmt.Sprintf("%04d: ", line) + string(lines[line-1]) + "\n"

	message += " \033[32m>>>─"
	var count = 0
	for count < token.GetPosition() {
		message += "─"
		count++
	}
	message += "⌃\033[36m\n"

	if line < len(lines) {
		message += fmt.Sprintf("%04d: ", line+1) + string(lines[line]) + "\n"
	}
	message += "\033[0m\n"

	return message
}

/*
This private class method is useful when creating scanner and parser error
messages that include the required grammatical rules.
*/
func (v *parser_) generateGrammar(expected string, symbols ...string) string {
	var message = "Was expecting '" + expected + "' from:\n"
	for _, symbol := range symbols {
		message += fmt.Sprintf(
			"  \033[32m%v: \033[33m%v\033[0m\n\n",
			symbol,
			grammar[symbol],
		)
	}
	return message
}

/*
This private class method attempts to read the next token from the token
stream and return it.
*/
func (v *parser_) getNextToken() TokenLike {
	var next TokenLike
	if v.next.IsEmpty() {
		var token, ok = v.tokens.RemoveHead() // Will block if queue is empty.
		if !ok {
			panic("The token channel terminated without an EOF token.")
		}
		next = token
		if next.GetType() == ErrorToken {
			var message = v.formatError(next)
			panic(message)
		}
	} else {
		next = v.next.RemoveTop()
	}
	return next
}

func (v *parser_) parseAssociation() (
	association col.AssociationLike[col.Key, col.Value],
	token TokenLike,
	ok bool,
) {
	var key col.Key
	var value col.Value
	key, token, ok = v.parseKey()
	if !ok {
		return association, token, false
	}
	_, _, ok = v.parseToken(DelimiterToken, ":")
	if !ok {
		// Put back the primitive key token.
		v.putBack(token)
		return association, token, false
	}
	value, token, ok = v.parseValue()
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("value",
			"$association",
			"$key",
			"$value")
		panic(message)
	}
	association = col.Association[col.Key, col.Value]().Make(key, value)
	return association, token, true
}

func (v *parser_) parseAssociations() (
	associations col.CatalogLike[col.Key, col.Value],
	token TokenLike,
	ok bool,
) {
	var association col.AssociationLike[col.Key, col.Value]
	associations = col.Catalog[col.Key, col.Value]().Make()

	// Handle the empty case.
	_, token, ok = v.parseToken(DelimiterToken, ":")
	if ok {
		return associations, token, true
	}

	// Handle the multi-line case.
	var _, eolToken, isMultilined = v.parseToken(EOLToken, "")
	if isMultilined {
		association, token, ok = v.parseAssociation()
		if !ok {
			// This must be a collection of values instead.
			v.putBack(eolToken)
			return associations, token, false
		}
		for ok {
			var key = association.GetKey()
			var value = association.GetValue()
			associations.SetValue(key, value)
			_, token, ok = v.parseToken(EOLToken, "")
			if !ok {
				var message = v.formatError(token)
				message += v.generateGrammar("EOL",
					"$associations",
					"$association",
				)
				panic(message)
			}
			association, token, ok = v.parseAssociation()
		}
		return associations, token, true
	}

	// Handle the in-line case.
	association, token, ok = v.parseAssociation()
	if !ok {
		return associations, token, false
	}
	for ok {
		var key = association.GetKey()
		var value = association.GetValue()
		associations.SetValue(key, value)
		_, token, ok = v.parseToken(DelimiterToken, ",")
		if ok {
			association, token, ok = v.parseAssociation()
			if !ok {
				var message = v.formatError(token)
				message += v.generateGrammar("association",
					"$associations",
					"$association",
				)
				panic(message)
			}
		}
	}
	return associations, token, true
}

func (v *parser_) parseCollection() (
	collection col.Collection,
	token TokenLike,
	ok bool,
) {
	var context string
	_, token, ok = v.parseToken(DelimiterToken, "[")
	if !ok {
		return collection, token, false
	}
	collection, _, ok = v.parseAssociations()
	if !ok {
		// The values must be attempted second since it may start with a component
		// which cannot be put back as a single token.
		collection, _, ok = v.parseValues()
		if !ok {
			var message = v.formatError(token)
			message += v.generateGrammar("associations",
				"$collection",
				"$associations",
				"$values",
			)
			panic(message)
		}
	}
	_, token, ok = v.parseToken(DelimiterToken, "]")
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("]",
			"$collection",
			"$associations",
			"$values",
		)
		panic(message)
	}
	_, token, ok = v.parseToken(DelimiterToken, "(")
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("(",
			"$collection",
			"$associations",
			"$values",
		)
		panic(message)
	}
	context, token, ok = v.parseToken(ContextToken, "")
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("CONTEXT",
			"$collection",
			"$associations",
			"$values",
		)
		panic(message)
	}
	_, token, ok = v.parseToken(DelimiterToken, ")")
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar(")",
			"$collection",
			"$associations",
			"$values",
		)
		panic(message)
	}
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
	var primitive col.Primitive
	primitive, token, ok = v.parsePrimitive()
	key = col.Key(primitive)
	return key, token, ok
}

func (v *parser_) parsePrimitive() (
	primitive col.Primitive,
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
	token = v.getNextToken()
	value = token.GetValue()
	if token.GetType() == expectedType {
		var constrained = len(expectedValue) > 0
		if !constrained || value == expectedValue {
			return value, token, true
		}
	}
	v.putBack(token)
	return "", token, false
}

func (v *parser_) parseValue() (
	value col.Value,
	token TokenLike,
	ok bool,
) {
	value, token, ok = v.parsePrimitive()
	if ok {
		return value, token, true
	}
	value, token, ok = v.parseCollection()
	if ok {
		return value, token, true
	}
	return value, token, false
}

func (v *parser_) parseValues() (
	values col.ListLike[col.Value],
	token TokenLike,
	ok bool,
) {
	var value col.Value
	values = col.List[col.Value]().Make()

	// Handle the empty case.
	_, token, ok = v.parseToken(DelimiterToken, "]")
	if ok {
		v.putBack(token)
		return values, token, true
	}

	// Handle the multi-line case.
	var _, _, isMultilined = v.parseToken(EOLToken, "")
	if isMultilined {
		value, token, ok = v.parseValue()
		if !ok {
			var message = v.formatError(token)
			message += v.generateGrammar("value",
				"$values",
				"$value",
			)
			panic(message)
		}
		for ok {
			values.AppendValue(value)
			_, token, ok = v.parseToken(EOLToken, "")
			if !ok {
				var message = v.formatError(token)
				message += v.generateGrammar("EOL",
					"$values",
					"$value",
				)
				panic(message)
			}
			value, token, ok = v.parseValue()
		}
		return values, token, true
	}

	// Handle the in-line case.
	value, token, ok = v.parseValue()
	if !ok {
		return values, token, false
	}
	for ok {
		values.AppendValue(value)
		_, token, ok = v.parseToken(DelimiterToken, ",")
		if ok {
			value, token, ok = v.parseValue()
			if !ok {
				var message = v.formatError(token)
				message += v.generateGrammar("value",
					"$values",
					"$value",
				)
				panic(message)
			}
		}
	}
	return values, token, true
}

func (v *parser_) putBack(token TokenLike) {
	v.next.AddValue(token)
}

/*
This Go map captures the syntax rules for collections of Go primitives.
*/
var grammar = map[string]string{
	"$association": `key ":" value`,
	"$associations": `
    association ("," association)*
    EOL (association EOL)+
    ":"  ! No associations.`,
	"$collection": `"[" (associations | values) "]" "(" CONTEXT ")"`,
	"$key":        `primitive`,
	"$primitive":  `BOOLEAN | COMPLEX | FLOAT | HEXADECIMAL | INTEGER | NIL | RUNE | STRING`,
	"$source":     `collection EOF  ! Terminated with an end-of-file marker.`,
	"$value":      `collection | primitive`,
	"$values": `
    value ("," value)*
    EOL (value EOL)+
    " "  ! No values.`,
}
