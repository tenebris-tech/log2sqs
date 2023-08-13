//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package parse

import (
	"errors"
	"fmt"
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
	_, ok := Parsers[format]
	return ok
}

// getParser returns the correct parser for the format or an error if the format is unknown
func getParser(format string) (Parser, error) {
	parser, ok := Parsers[format]
	if !ok {
		return Parser{}, errors.New(fmt.Sprintf("unknown format %s", format))
	}
	return parser, nil
}
