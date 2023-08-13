//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package parse

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"log2sqs/config"
)

// regexParser parses Apache/NGINX logs into a GELF message based on the supplied regex to field map
func regexParser(s string, p *Parser) (GELFMessage, error) {

	// Check if the parser has been initialized
	if p.regex == nil {

		// Iterate over fields in order to build the regex
		tmp := ""
		for i := 1; i <= len(p.regexFields); i++ {
			tmp = tmp + p.regexFields[i].regex
		}

		r, err := regexp.Compile(tmp)
		if err != nil {
			return GELFMessage{}, errors.New(fmt.Sprintf("regex failed to compile: %s", err.Error()))
		}

		// Save the pointer to the compiled regex for future use
		p.regex = r
	}

	// Parse the line using the regex
	result := p.regex.FindStringSubmatch(s)

	if len(result) < (p.requireFields + 1) {
		return GELFMessage{}, errors.New("too few fields found")
	}

	// Start the GELF message
	g := GELFMessage{}
	g["version"] = "1.1"
	g["host"] = config.Hostname
	g["_original_format"] = p.format

	// Iterate over the fields and add them to the GELF message
	for i := 1; i < len(result); i++ {
		switch p.regexFields[i].fType {

		case "int":
			g[p.regexFields[i].field] = String2Int(result[i])

		case "date":
			tmp := result[i]
			if p.regexFields[i].addTZ {
				tmp = tmp + " +0000"
			}
			t, err := time.Parse(p.regexFields[i].dateFormat, tmp)
			if err != nil {
				return GELFMessage{}, errors.New(fmt.Sprintf("unable to parse date %s using format %s: %s", result[i], p.regexFields[i].dateFormat, err.Error()))
			}
			g[p.regexFields[i].field] = t.Unix()

		case "string":
			g[p.regexFields[i].field] = EmptyString(result[i])

		default:
			g[p.regexFields[i].field] = EmptyString(result[i])
		}

		// Should this also be the short message field (required)?
		if p.regexFields[i].shortMessage {
			g["short_message"] = g[p.regexFields[i].field]
		}

	}
	return g, nil
}

// String2Int returns the integer contained in string or 0
func String2Int(s string) int {
	ret, err := strconv.Atoi(s)
	if err != nil {
		ret = 0
	}
	return ret
}

// EmptyString returns a cleaned string or "-" if empty
func EmptyString(s string) string {
	r := strings.TrimSpace(s)
	if r == "" {
		r = "-"
	}
	return r
}
