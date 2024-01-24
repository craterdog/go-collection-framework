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

// CLASS NAMESPACE

// Private Class Namespace Type

type parserClass_ struct {
	channelSize int
	stackSize   int
}

// Private Class Namespace Reference

var parserClass = &parserClass_{
	channelSize: 16,
	stackSize:   4,
}

// Public Class Namespace Access

func ParserClass() col.ParserClassLike {
	return parserClass
}

// Public Class Constructors

func (c *parserClass_) Make() col.ParserLike {
	var parser = &parser_{
		next: col.StackClass[col.TokenLike]().MakeWithCapacity(c.stackSize),
	}
	return parser
}

// CLASS INSTANCES

// Private Class Type Definition

type parser_ struct {
	source string
	next   col.StackLike[col.TokenLike] // A stack of unprocessed retrieved tokens.
	tokens chan col.TokenLike           // A queue of unread tokens from the scanner.
}

// Public Interface

func (v *parser_) ParseSource(source string) col.Collection {
	// The scanner runs in a separate Go routine.
	v.source = source
	v.tokens = make(chan col.TokenLike, parserClass.channelSize)
	ScannerClass().Make(v.source, v.tokens)

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
	_, token, ok = v.parseToken(TypeEOF, "")
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
func (v *parser_) formatError(token col.TokenLike) string {
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

// This private class method is useful when creating scanner and parser error
// messages that include the required grammatical rules.
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

// This private class method attempts to read the next token from the token
// stream and return it.
func (v *parser_) getNextToken() col.TokenLike {
	var next col.TokenLike
	if v.next.IsEmpty() {
		var token, ok = <-v.tokens
		if !ok {
			panic("The token channel terminated without an EOF token.")
		}
		next = token
		if next.GetType() == TypeError {
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
	token col.TokenLike,
	ok bool,
) {
	var key col.Key
	var value col.Value
	key, token, ok = v.parseKey()
	if !ok {
		return association, token, false
	}
	_, _, ok = v.parseToken(TypeDelimiter, ":")
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
	association = col.AssociationClass[col.Key, col.Value]().Make(key, value)
	return association, token, true
}

func (v *parser_) parseAssociations() (
	associations col.CatalogLike[col.Key, col.Value],
	token col.TokenLike,
	ok bool,
) {
	var association col.AssociationLike[col.Key, col.Value]
	associations = col.CatalogClass[col.Key, col.Value]().Make()

	// Handle the empty case.
	_, token, ok = v.parseToken(TypeDelimiter, ":")
	if ok {
		return associations, token, true
	}

	// Handle the multi-line case.
	var eolToken col.TokenLike
	_, eolToken, ok = v.parseToken(TypeEOL, "")
	if ok {
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
			_, token, ok = v.parseToken(TypeEOL, "")
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
		_, token, ok = v.parseToken(TypeDelimiter, ",")
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
	token col.TokenLike,
	ok bool,
) {
	var context string
	_, token, ok = v.parseToken(TypeDelimiter, "[")
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
	_, token, ok = v.parseToken(TypeDelimiter, "]")
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("]",
			"$collection",
			"$associations",
			"$values",
		)
		panic(message)
	}
	_, token, ok = v.parseToken(TypeDelimiter, "(")
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("(",
			"$collection",
			"$associations",
			"$values",
		)
		panic(message)
	}
	context, token, ok = v.parseToken(TypeContext, "")
	if !ok {
		var message = v.formatError(token)
		message += v.generateGrammar("CONTEXT",
			"$collection",
			"$associations",
			"$values",
		)
		panic(message)
	}
	_, token, ok = v.parseToken(TypeDelimiter, ")")
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
			collection = col.ArrayClass[col.Value]().MakeFromArray(sequence.AsArray())
		case "List":
			collection = col.ListClass[col.Value]().MakeFromSequence(sequence)
		case "Queue":
			collection = col.QueueClass[col.Value]().MakeFromSequence(sequence)
		case "Set":
			collection = col.SetClass[col.Value]().MakeFromSequence(sequence)
		case "Stack":
			collection = col.StackClass[col.Value]().MakeFromSequence(sequence)
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
			collection = col.MapClass[col.Key, col.Value]().MakeFromArray(sequence.AsArray())
		case "Catalog":
			collection = col.CatalogClass[col.Key, col.Value]().MakeFromSequence(sequence)
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
	token col.TokenLike,
	ok bool,
) {
	var primitive col.Primitive
	primitive, token, ok = v.parsePrimitive()
	key = col.Key(primitive)
	return key, token, ok
}

func (v *parser_) parsePrimitive() (
	primitive col.Primitive,
	token col.TokenLike,
	ok bool,
) {
	_, token, ok = v.parseToken(TypeBoolean, "")
	if ok {
		primitive, _ = stc.ParseBool(token.GetValue())
		return primitive, token, true
	}
	_, token, ok = v.parseToken(TypeComplex, "")
	if ok {
		primitive, _ = stc.ParseComplex(token.GetValue(), 128)
		return primitive, token, true
	}
	_, token, ok = v.parseToken(TypeFloat, "")
	if ok {
		primitive, _ = stc.ParseFloat(token.GetValue(), 64)
		return primitive, token, true
	}
	_, token, ok = v.parseToken(TypeHexadecimal, "")
	if ok {
		primitive, _ = stc.ParseUint(token.GetValue()[2:], 16, 64)
		return primitive, token, true
	}
	_, token, ok = v.parseToken(TypeInteger, "")
	if ok {
		primitive, _ = stc.ParseInt(token.GetValue(), 10, 64)
		return primitive, token, true
	}
	_, token, ok = v.parseToken(TypeNil, "")
	if ok {
		primitive = nil
		return primitive, token, true
	}
	_, token, ok = v.parseToken(TypeRune, "")
	if ok {
		var matches = ScannerClass().MatchToken(TypeRune, token.GetValue())
		var match, _ = stc.Unquote(matches[0])
		primitive, _ = utf.DecodeRuneInString(match)
		return primitive, token, true
	}
	_, token, ok = v.parseToken(TypeString, "")
	if ok {
		var matches = ScannerClass().MatchToken(TypeString, token.GetValue())
		primitive, _ = stc.Unquote(matches[0])
		return primitive, token, true
	}
	return primitive, token, ok
}

func (v *parser_) parseToken(
	expectedType col.TokenType,
	expectedValue string,
) (value string, token col.TokenLike, ok bool) {
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
	token col.TokenLike,
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
	token col.TokenLike,
	ok bool,
) {
	var value col.Value
	values = col.ListClass[col.Value]().Make()

	// Handle the empty case.
	_, token, ok = v.parseToken(TypeDelimiter, "]")
	if ok {
		v.putBack(token)
		return values, token, true
	}

	// Handle the multi-line case.
	_, _, ok = v.parseToken(TypeEOL, "")
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
		for ok {
			values.AppendValue(value)
			_, token, ok = v.parseToken(TypeEOL, "")
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
		_, token, ok = v.parseToken(TypeDelimiter, ",")
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

func (v *parser_) putBack(token col.TokenLike) {
	v.next.AddValue(token)
}

// This Go map captures the syntax rules for collections of Go primitives.
var grammar = map[string]string{
	"$source":     `collection EOF  ! Terminated with an end-of-file marker.`,
	"$collection": `"[" (associations | values) "]" "(" CONTEXT ")"`,
	"$associations": `
    association ("," association)*
    EOL (association EOL)+
    ":"  ! No associations.`,
	"$association": `key ":" value`,
	"$key":         `primitive`,
	"$values": `
    value ("," value)*
    EOL (value EOL)+
    " "  ! No values.`,
	"$value":     `collection | primitive`,
	"$primitive": `BOOLEAN | COMPLEX | FLOAT | INTEGER | NIL | RUNE | STRING | UNSIGNED`,
}
