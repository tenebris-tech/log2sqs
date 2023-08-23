//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package parse

import "errors"

// AddRegexParser adds a new regex parser to the list of available parsers
func AddRegexParser(name string, fields RegexFields) error {
	if name == "" {
		return errors.New("parser name cannot be empty")
	}

	if len(fields) < 1 {
		return errors.New("parser fields cannot be empty")
	}

	Parsers[name] = Parser{format: name, parserFunc: regexParser, regexFields: fields, requireFields: len(fields)}
	return nil
}
