//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package event

import (
	"encoding/json"
	"log"

	"log2sqs/config"
	"log2sqs/global"
	"log2sqs/parse"
)

// Log reports significant internal event through GELF message
func Log(message string, full string, level int) {

	// Send to local logger
	log.Print(message)

	// Create GELF message
	g := parse.GELFMessage{}
	g["version"] = "1.1"
	g["host"] = config.Config.Hostname
	g["level"] = level
	g["_app_name"] = global.ProductName + " " + global.ProductVersion
	g["short_message"] = message
	g["full_message"] = full
	g["timestamp"] = global.TimeStamp()
	g["_via_hostname"] = config.Config.Hostname
	g["_via_proto"] = "gelf"

	// Add source IP
	if config.Config.SyslogOverrideSourceIP != "" {
		g["_event_source_ip"] = config.Config.SyslogOverrideSourceIP
	} else {
		g["_event_source_ip"] = global.GetOutboundIP()
	}

	// Add any static fields
	for key, value := range config.Config.AddFields {
		g[key] = value
	}

	// Marshal JSON for queue
	gBytes, err := json.Marshal(g)
	if err != nil {
		log.Printf("Internal event: error marshaling to JSON: %s", err.Error())
		return
	}

	if config.Config.Debug {
		log.Printf("Syslog: %s", string(gBytes[:]))
	}

	// Add to memory buffer
	Add(gBytes)
}
