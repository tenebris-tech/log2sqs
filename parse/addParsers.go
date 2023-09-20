//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package parse

import (
	"errors"
	"strings"

	"log2sqs/config"
)

// AddCustomParsers iterates through any custom parsers in config.Config and adds them to the list of available parsers
func AddCustomParsers() error {
	for _, p := range config.Config.CustomParsers {
		switch strings.ToLower(p.Type) {
		case "regex":
			err := AddRegexParser(p.Name, p.RegexFields)
			if err != nil {
				return err
			}
		default:
			return errors.New("unknown parser type")
		}
	}
	return nil
}

// AddRegexParser adds a new regex parser to the list of available parsers
func AddRegexParser(name string, fields config.RegexFields) error {
	if name == "" {
		return errors.New("parser name cannot be empty")
	}

	if len(fields) < 1 {
		return errors.New("parser fields cannot be empty")
	}

	parsersMX.Lock()
	defer parsersMX.Unlock()
	parsers[name] = Parser{format: name, parserFunc: regexParser, regexFields: fields, requireFields: len(fields)}
	return nil
}
