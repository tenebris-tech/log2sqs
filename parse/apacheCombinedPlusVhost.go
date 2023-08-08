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

// ApacheCombinedPlusVhost parses Apache/NGINX combined log format into a GELF message with additional fields
// LogFormat "%v:%p %h %l %u %t \"%r\" %>s %O \"%{Referer}i\" \"%{User-Agent}i\" %D \"%m\" \"%U\" \"%q\"" combinedplusvhost
func ApacheCombinedPlusVhost(s string, g GELFMessage) error {

	var buffer bytes.Buffer
	buffer.WriteString(`^(\S+):`)                  // Vhost (1)
	buffer.WriteString(`(\S+)\s`)                  // Vhost Port (2)
	buffer.WriteString(`(\S+)\s`)                  // IP (3)
	buffer.WriteString(`(\S+)\s`)                  // ident (4)
	buffer.WriteString(`(\S+)\s`)                  // user (5)
	buffer.WriteString(`\[([^]]+)\]\s`)            // date, time, and zone (6)
	buffer.WriteString(`"(.*?)"\s`)                // URL (7)
	buffer.WriteString(`(\S+)\s`)                  // status code (8)
	buffer.WriteString(`(\S+)\s`)                  // size (9)
	buffer.WriteString(`"((?:[^"]*(?:\\")?)*)"\s`) // referrer (10)
	buffer.WriteString(`"(.*?)"\s`)                // user agent (11)
	buffer.WriteString(`(\S+)\s`)                  // processing time (12)
	buffer.WriteString(`"(.*?)"\s`)                // request method (13)
	buffer.WriteString(`"(.*?)"\s`)                // request path (14)
	buffer.WriteString(`"(.*?)"$`)                 // request query (15)

	re1, err := regexp.Compile(buffer.String())
	if err != nil {
		return err
	}
	result := re1.FindStringSubmatch(s)

	if len(result) < 14 {
		return errors.New("too few fields found")
	}

	// Parse time
	layout := "02/Jan/2006:15:04:05 -0700"
	t, err := time.Parse(layout, result[6])
	if err != nil {
		return err
	}

	// Create GELF message
	g["version"] = "1.1"
	g["host"] = config.Hostname
	g["short_message"] = result[7]
	g["timestamp"] = t.Unix()
	g["_src_ip"] = result[3]
	g["_user"] = result[5]
	g["_http_request"] = result[7]
	g["_http_status"] = String2Int(result[8])
	g["_http_response_size"] = String2Int(result[9])
	g["_http_referer"] = result[10]
	g["_user_agent"] = result[11]
	g["_duration_usec"] = String2Int(result[12])
	g["_http_request_method"] = EmptyString(result[13])
	g["_http_request_path"] = EmptyString(result[14])
	g["_http_request_query"] = EmptyString(result[15])
	g["_vhost"] = EmptyString(result[1])
	g["_vhost_port"] = String2Int(result[2])
	g["_original_format"] = "ApacheCombinedPlusVhost"
	return nil
}
