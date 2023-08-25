//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package syslog

import (
	"strings"
	"time"

	"log2sqs/config"
	"log2sqs/global"
	"log2sqs/parse"
)

// plainText handles log events that can not otherwise be parsed
func plainText(buf []byte, srcIP string, g parse.GELFMessage) error {

	g["version"] = "1.1"
	g["_via_hostname"] = config.Config.Hostname
	g["_via_proto"] = "syslog_udp"
	g["host"] = srcIP
	g["short_message"] = strings.TrimSuffix(string(buf), "\n")
	g["_original_format"] = "unknown"
	g["timestamp"] = time.Now().Unix()

	if config.Config.SyslogOverrideSourceIP != "" {
		g["_event_source_ip"] = config.Config.SyslogOverrideSourceIP
	} else {
		if config.Config.SyslogReplaceLocalhost && srcIP == "127.0.0.1" {
			g["_event_source_ip"] = global.GetOutboundIP()
			if g["_event_source_ip"] == "" {
				g["_event_source_ip"] = srcIP
			}
		} else {
			g["_event_source_ip"] = srcIP
		}
	}

	if config.Config.SyslogFullMessage {
		g["full_message"] = strings.TrimSuffix(string(buf), "\n")
	}

	return nil
}
