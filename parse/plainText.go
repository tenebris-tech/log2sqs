//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package parse

import (
	"log2sqs/config"
	"strings"
	"time"
)

// PlainText turns the string into a GELF message
func PlainText(s string, g GELFMessage) error {

	// Create GELF message
	g["version"] = "1.1"
	g["host"] = config.Hostname
	g["short_message"] = strings.TrimSuffix(s, "\n")
	g["timestamp"] = time.Now().Unix()
	g["_original_format"] = "PlainText"
	return nil
}
