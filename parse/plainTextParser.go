//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package parse

import (
	"strings"

	"log2sqs/config"
	"log2sqs/global"
)

// plainTextParser turns the string into a GELF message
func plainTextParser(s string, _ *Parser) (GELFMessage, error) {
	g := GELFMessage{}
	g["version"] = "1.1"
	g["host"] = config.Config.Hostname
	g["short_message"] = strings.TrimSuffix(s, "\n")
	g["timestamp"] = global.TimeStamp()
	g["_original_format"] = "PlainText"
	return g, nil
}
