//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package parse

import (
	"bytes"
	"errors"
	"regexp"
	"time"

	"log2sqs/config"
)

// ApacheCombined parses Apache/NGINX combined log format into a GELF message
func ApacheCombined(s string, g GELFMessage) error {

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
		return errors.New("too few fields found")
	}

	// Parse time
	layout := "02/Jan/2006:15:04:05 -0700"
	t, err := time.Parse(layout, result[4])
	if err != nil {
		return err
	}

	// Create GELF message
	g["version"] = "1.1"
	g["host"] = config.Hostname
	g["short_message"] = result[5]
	g["timestamp"] = t.Unix()
	g["_src_ip"] = result[1]
	g["_user"] = result[3]
	g["_http_request"] = result[5]
	g["_http_status"] = String2Int(result[6])
	g["_http_response_size"] = String2Int(result[7])
	g["_http_referer"] = result[8]
	g["_user_agent"] = result[9]
	g["_original_format"] = "ApacheCombined"
	return nil
}
