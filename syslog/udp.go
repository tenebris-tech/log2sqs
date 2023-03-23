//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package syslog

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"log2sqs/config"
	"log2sqs/event"
	"log2sqs/global"
)

func UDP() {

	// Infinite loop to allow restarts
	for {

		// Listen for incoming udp packets
		pc, err := net.ListenPacket("udp", config.SyslogUDP)
		if err != nil {
			event.Log(fmt.Sprintf("Error starting UDP listener on %s: %s", config.SyslogUDP, err.Error()), "", global.ERR)
			time.Sleep(10 * time.Second)
			continue
		}

		event.Log(fmt.Sprintf("Listening for Syslog messages on UDP %s", config.SyslogUDP), "", global.INFO)

		// Loop and receive UDP datagrams
		for {
			// ReadFrom will respect the length of buf, so we don't need to worry about buffer
			// overflows. If the packet contains more data than len(buf) it will be truncated.
			buf := make([]byte, config.SyslogUDPMax)
			n, addr, err := pc.ReadFrom(buf)
			if err != nil {
				// Read error is unusual, restart listener
				event.Log(fmt.Sprintf("UDP read error: %s", err.Error()), "", global.ERR)
				break
			}

			srcIP := strings.Split(safeAddrString(addr), ":")[0]

			if config.Debug {
				log.Printf("Received %d bytes from %s", n, safeAddrString(addr))
				event.Dump(buf[:n])
			}

			// Process the message
			err = syslogProcess(buf[:n], srcIP)
			if err != nil {
				event.Log(err.Error(), string(buf[:n]), global.ERR)
				event.Dump(buf[:n])
			}
		}
	}
}
