//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package parse

import (
	"errors"
	"fmt"
	"log2sqs/config"
)

// Parse parses the line into a GELF message based on the format
func Parse(line string, format string) (GELFMessage, error) {

	// Get the right parser for the format
	parser, err := getParser(format)
	if err != nil {
		return GELFMessage{}, err
	}

	// Call the parser function for the format
	gelf, err := parser.parserFunc(line, &parser)
	if err != nil {
		// If an error occurs, augment it
		return GELFMessage{}, errors.New(fmt.Sprintf("error parsing log format %s: %s", format, err.Error()))
	}

	return gelf, nil
}

// CheckFormat returns true if the format is in the parsers map
func CheckFormat(format string) bool {
	parsersMX.RLock()
	defer parsersMX.RUnlock()
	_, ok := parsers[format]
	return ok
}

// getParser returns the correct parser for the format or an error if the format is unknown
func getParser(format string) (Parser, error) {
	parsersMX.RLock()
	defer parsersMX.RUnlock()

	parser, ok := parsers[format]

	if !ok {
		return Parser{}, errors.New(fmt.Sprintf("unknown format %s", format))
	}

	// Create a deep copy for thread safety
	deepCopy := Parser{
		format:        parser.format,
		parserFunc:    parser.parserFunc,
		regex:         parser.regex,
		regexFields:   make(config.RegexFields),
		requireFields: parser.requireFields,
	}

	// Copy the RegexFields map
	for k, v := range parser.regexFields {
		deepCopy.regexFields[k] = v
	}

	return deepCopy, nil
}
