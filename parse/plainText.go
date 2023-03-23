//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package parse

import (
	"strings"

	"log2sqs/config"
	"log2sqs/global"
)

// PlainText turns the string into a GELF message
func PlainText(s string, g GELFMessage) error {

	// Create GELF message
	g["version"] = "1.1"
	g["host"] = config.Hostname
	g["short_message"] = strings.TrimSuffix(s, "\n")
	g["timestamp"] = global.TimeStamp()
	g["_original_format"] = "PlainText"
	return nil
}
