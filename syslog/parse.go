//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package syslog

import (
	"fmt"
	"log"
	"log2sqs/event"
	"log2sqs/global"
	"math"
	"strings"
	"time"

	"github.com/jeromer/syslogparser"
	"github.com/jeromer/syslogparser/rfc3164"
	"github.com/jeromer/syslogparser/rfc5424"

	"log2sqs/config"
	"log2sqs/parse"
)

func parseSyslog(buf []byte, srcIP string, g parse.GELFMessage) error {

	// Try to identify the syslog format
	rfc, err := syslogparser.DetectRFC(buf)
	if err != nil {
		log.Printf("unable to determine syslog format: %s", err.Error())
		return plainText(buf, srcIP, g)
	}

	switch rfc {
	case syslogparser.RFC_UNKNOWN:
		return plainText(buf, srcIP, g)

	case syslogparser.RFC_3164:
		p := rfc3164.NewParser(buf)
		err := p.Parse()
		if err != nil {
			event.Log(fmt.Sprintf("error parsing RFC3164 message: %s", err), string(buf), global.WARN)
			return plainText(buf, srcIP, g)
		}

		// Dump() returns a map[string]interface{}
		eventMap := p.Dump()
		g["version"] = "1.1"
		g["_via_hostname"] = config.Hostname
		g["_via_proto"] = "syslog_udp"
		g["host"] = eventMap["hostname"]
		g["level"] = eventMap["severity"]
		g["_facility"] = global.GetFacility(eventMap["facility"].(int))
		g["_app_name"] = eventMap["tag"]
		g["short_message"] = strings.TrimSuffix(fmt.Sprint(eventMap["content"]), "\n")
		g["_original_format"] = "RFC3164"

		if config.SyslogOverrideSourceIP != "" {
			g["_source_ip"] = config.SyslogOverrideSourceIP
		} else {
			g["_source_ip"] = srcIP
		}

		if config.SyslogOverrideTime {
			g["timestamp"] = time.Now().Unix()
		} else {
			g["timestamp"] = eventMap["timestamp"].(time.Time).Unix()

			// Sanity check timestamp
			// This safeguards against systems that log in local time instead of UTC with no time zone
			// or have clocks that are out of whack
			if math.Abs(float64(g["timestamp"].(int64))-float64(time.Now().Unix())) > 240 {
				g["timestamp"] = time.Now().Unix()
			}
		}

		if config.SyslogFullMessage {
			g["full_message"] = strings.TrimSuffix(string(buf), "\n")
		}

	case syslogparser.RFC_5424:
		p := rfc5424.NewParser(buf)
		err := p.Parse()
		if err != nil {
			event.Log(fmt.Sprintf("error parsing RFC5242 message: %s", err), string(buf), global.WARN)
			return plainText(buf, srcIP, g)
		}

		// Dump() returns a map[string]interface{}
		eventMap := p.Dump()
		g["version"] = "1.1"
		g["_via_hostname"] = config.Hostname
		g["_via_proto"] = "syslog_udp"
		g["host"] = eventMap["hostname"]
		g["level"] = eventMap["severity"]
		g["_facility"] = global.GetFacility(eventMap["facility"].(int))
		g["_app_name"] = eventMap["app_name"]
		g["_proc_id"] = eventMap["proc_id"]
		g["short_message"] = strings.TrimSuffix(fmt.Sprint(eventMap["message"]), "\n")
		g["_structured_data"] = eventMap["structured_data"]
		g["_original_format"] = "RFC5424"

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

		if config.SyslogOverrideTime {
			g["timestamp"] = time.Now().Unix()
		} else {
			g["timestamp"] = eventMap["timestamp"].(time.Time).Unix()

			// Sanity check timestamp
			// This safeguards against systems that log in local time instead of UTC with no time zone
			// or have clocks that are out of whack
			if math.Abs(float64(g["timestamp"].(int64))-float64(time.Now().Unix())) > 240 {
				g["timestamp"] = time.Now().Unix()
			}
		}

		if config.SyslogFullMessage {
			g["full_message"] = strings.TrimSuffix(string(buf), "\n")
		}
	}
	return nil
}

// plainText handles log events that can not otherwise be parsed
func plainText(buf []byte, srcIP string, g parse.GELFMessage) error {
	g["version"] = "1.1"
	g["_via_hostname"] = config.Hostname
	g["_via_proto"] = "syslog_udp"
	g["host"] = srcIP
	g["short_message"] = strings.TrimSuffix(string(buf), "\n")
	g["_original_format"] = "unknown"
	g["timestamp"] = time.Now().Unix()

	if config.SyslogOverrideSourceIP != "" {
		g["_source_ip"] = config.SyslogOverrideSourceIP
	} else {
		g["_source_ip"] = srcIP
	}

	if config.SyslogFullMessage {
		g["full_message"] = strings.TrimSuffix(string(buf), "\n")
	}

	return nil
}
