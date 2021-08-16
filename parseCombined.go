//
// Copyright (c) 2021 Tenebris Technologies Inc.
//

package main

import (
	"bytes"
	"errors"
	"os"
	"regexp"
	"time"
)

// Parse Apache2 "combined" log format into object
func parseCombined(s string, j arbitraryJSON) error {

	var buffer bytes.Buffer
	buffer.WriteString(`^(\S+)\s`)                 // IP (1)
	buffer.WriteString(`(\S+)\s`)                  // ident (2)
	buffer.WriteString(`(\S+)\s`)                  // user (3)
	buffer.WriteString(`\[([^]]+)\]\s`)            // date, time, and zone (4)
	buffer.WriteString(`"(.*?)"\s`)                // URL (5)
	buffer.WriteString(`(\S+)\s`)                  // status code (6)
	buffer.WriteString(`(\S+)\s`)                  // size (7)
	buffer.WriteString(`"((?:[^"]*(?:\\")?)*)"\s`) // referrer (8)
	buffer.WriteString(`"(.*)"$`)                  // user agent (9)

	re1, err := regexp.Compile(buffer.String())
	if err != nil {
		return err
	}
	result := re1.FindStringSubmatch(s)

	if len(result) < 10 {
		// Too few items found
		return errors.New("error parsing Combined Log Format")
	}

	// Parse time
	layout := "02/Jan/2006:15:04:05 -0700"
	t, err := time.Parse(layout, result[4])
	if err != nil {
		return err
	}

	// Create GELF message
	j["version"] = "1.1"
	j["host"], _ = os.Hostname()
	j["short_message"] = result[5]
	j["timestamp"] = t.Unix()
	j["_src_ip"] = result[1]
	j["_user"] = result[3]
	j["_http_request"] = result[5]
	j["_http_status"] = result[6]
	j["_http_response_size"] = result[7]
	j["_http_referer"] = result[8]
	j["_user_agent"] = result[9]
	return nil
}
