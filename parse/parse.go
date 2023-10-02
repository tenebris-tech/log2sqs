//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package parse

import (
	"errors"
	"fmt"
	"regexp"

	"log2sqs/config"
)

// Enum for parser types
const (
	RegexParserType = iota
	GelfParserType
	PlainTextParserType
)

// CheckFormat checks if the format string is valid
func CheckFormat(format string) bool {
	parsersMX.RLock()
	defer parsersMX.RUnlock()

	_, ok := parsers[format]
	return ok
}

// New returns a new parser instance configured for the specified format
func New(format string) (*Parser, error) {
	parsersMX.RLock()
	defer parsersMX.RUnlock()

	// Get the parser for the format if it exists
	p, ok := parsers[format]
	if !ok {
		return &Parser{}, errors.New(fmt.Sprintf("unknown format %s", format))
	}

	// Create a deep copy for thread safety
	var parser = Parser{
		format:        p.format,
		parserType:    p.parserType,
		requireFields: p.requireFields,
		regexFields:   make(config.RegexFields),
		regex:         nil,
	}

	if parser.parserType == RegexParserType {
		// Copy the RegexFields map
		for k, v := range p.regexFields {
			parser.regexFields[k] = v
		}

		// Iterate over fields to build the Regex
		tmp := ""
		for i := 1; i <= len(parser.regexFields); i++ {
			tmp = tmp + parser.regexFields[i].Regex
		}

		r, err := regexp.Compile(tmp)
		if err != nil {
			return &Parser{}, errors.New(fmt.Sprintf("Regex failed to compile: %s", err.Error()))
		}

		// Save the pointer to the compiled Regex for future use
		parser.regex = r
	}

	return &parser, nil
}

// Parse parses the line into a GELF message based on the format
func (p *Parser) Parse(line string) (GELFMessage, error) {

	// Invoke the correct parser
	switch p.parserType {
	case RegexParserType:
		return p.regexParser(line)
	case GelfParserType:
		return p.gelfParser(line)
	case PlainTextParserType:
		return p.plainTextParser(line)
	default:
		return GELFMessage{}, errors.New("unknown parser type")
	}
}
