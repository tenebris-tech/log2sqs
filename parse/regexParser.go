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

// regexParser parses Apache/NGINX logs into a GELF message based on the supplied Regex to Field map
func regexParser(s string, p *Parser) (GELFMessage, error) {

	// Check if the parser has been initialized
	if p.regex == nil {

		// Iterate over fields in order to build the Regex
		tmp := ""
		for i := 1; i <= len(p.regexFields); i++ {
			tmp = tmp + p.regexFields[i].Regex
		}

		r, err := regexp.Compile(tmp)
		if err != nil {
			return GELFMessage{}, errors.New(fmt.Sprintf("Regex failed to compile: %s", err.Error()))
		}

		// Save the pointer to the compiled Regex for future use
		p.regex = r
	}

	// Parse the line using the Regex
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
		switch p.regexFields[i].FType {

		case "int":
			g[p.regexFields[i].Field] = String2Int(result[i])

		case "date":
			tmp := result[i]
			if p.regexFields[i].AddTZ {
				tmp = tmp + " +0000"
			}
			t, err := time.Parse(p.regexFields[i].DateFormat, tmp)
			if err != nil {
				return GELFMessage{}, errors.New(fmt.Sprintf("unable to parse date %s using format %s: %s", result[i], p.regexFields[i].DateFormat, err.Error()))
			}
			g[p.regexFields[i].Field] = t.Unix()

		case "string":
			g[p.regexFields[i].Field] = EmptyString(result[i])

		default:
			g[p.regexFields[i].Field] = EmptyString(result[i])
		}

		// Should this also be the short message Field (required)?
		if p.regexFields[i].ShortMessage {
			g["short_message"] = g[p.regexFields[i].Field]
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
