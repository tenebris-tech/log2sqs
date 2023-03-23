//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package syslog

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log2sqs/config"
	"log2sqs/global"
	"log2sqs/parse"
	"math"
	"strconv"
)

// plainText handles log events that can not otherwise be parsed
func gelf(buf []byte, srcIP string, g parse.GELFMessage) error {

	// Remove leading and trailing junk
	jStart := -1
	jEnd := -1
	for i := 0; i < len(buf); i++ {

		// Find first '{'
		if jStart == -1 && buf[i] == '{' {
			jStart = i
		}

		// Find last '}'
		if buf[i] == '}' {
			jEnd = i
		}
	}

	if jStart == -1 || jEnd == -1 {
		// Can't be a JSON message
		return errors.New("not JSON")
	}

	// Copy JSON part of buffer
	var jsonBuf []byte
	for i := jStart; i <= jEnd; i++ {
		jsonBuf = append(jsonBuf, buf[i])
	}

	if config.Debug {
		log.Printf("Attempting to unmarshal: %s", string(jsonBuf))
	}

	// Attempt to unmarshal
	var j parse.GELFMessage
	err := json.Unmarshal(jsonBuf, &j)
	if err != nil {
		return errors.New(fmt.Sprintf("unmarshal failed: %s", err.Error()))
	}

	if config.Debug {
		log.Printf("Successfully unmarshalled: %s", string(jsonBuf))
	}

	// Test for valid GELF
	if !gMatch(j, "version", "1.1") {
		return errors.New("invalid GELF, missing version")
	}

	if !gExists(j, "host") {
		return errors.New("invalid GELF, host missing")
	}

	if !gExists(j, "short_message") && !gExists(j, "long_message") {
		return errors.New("invalid GELF, no message")
	}

	// GELF seems valid, copy to our structure
	for k, v := range j {
		g[k] = v
	}

	// Sanity check timestamp
	if gExists(j, "timestamp") {
		if math.Abs(float64(gGetInt(j, "timestamp"))-global.TimeStamp()) > 240 {
			g["timestamp"] = global.TimeStamp()
		}
	} else {
		// Timestamp is missing, so just set it
		g["timestamp"] = global.TimeStamp()
	}

	// Add hostname and protocol
	g["_via_hostname"] = config.Hostname
	g["_via_proto"] = "syslog_gelf"

	// Add source IP
	if config.SyslogOverrideSourceIP != "" {
		g["_event_source_ip"] = config.SyslogOverrideSourceIP
	} else {
		if config.SyslogReplaceLocalhost && srcIP == "127.0.0.1" {
			g["_event_source_ip"] = global.GetOutboundIP()
			if g["_event_source_ip"] == "" {
				g["_event_source_ip"] = srcIP
			}
		} else {
			g["_event_source_ip"] = srcIP
		}
	}

	return nil
}

// gMatch returns true if the key exists and matches the supplied string
func gMatch(j parse.GELFMessage, s string, r string) bool {
	if gGetStr(j, s) == r {
		return true
	}
	return false
}

// gExists returns true if the key exists
func gExists(j parse.GELFMessage, k string) bool {
	_, ok := j[k]
	if ok {
		return true
	}
	return false
}

// gGetInt safely retrieves an int value or returns 0
func gGetInt(j parse.GELFMessage, k string) int64 {
	val, ok := j[k]
	if !ok {
		return 0
	}

	s := fmt.Sprintf("%v", val)

	// Convert string to int634
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return n
}

// gGetStr safely returns the string or ""
func gGetStr(j parse.GELFMessage, k string) string {
	val, ok := j[k]
	if !ok {
		return ""
	}

	switch r := val.(type) {
	case string:
		return r
	default:
		return ""
	}
}
