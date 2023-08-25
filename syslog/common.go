//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package syslog

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"

	"log2sqs/config"
	"log2sqs/event"
	"log2sqs/parse"
)

func syslogProcess(buf []byte, srcIP string) error {

	// Parse the message
	g := parse.GELFMessage{}
	err := parseSyslog(buf, srcIP, g)
	if err != nil {
		return errors.New(fmt.Sprintf("error parsing syslog message: %s", err.Error()))
	}

	// Add any static fields
	for key, value := range config.Config.AddFields {
		g[key] = value
	}

	// Marshal JSON for queue
	gBytes, err := json.Marshal(g)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshaling to JSON: %s", err.Error()))
	}

	if config.Config.Debug {
		log.Printf("Syslog: %s", string(gBytes[:]))
	}

	// Add to memory buffer
	event.Add(gBytes)
	return nil
}

// safeAddrString returns a string or "" if addr is nil
func safeAddrString(addr net.Addr) string {
	if addr == nil {
		return ""
	}
	return addr.String()
}
