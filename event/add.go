//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package event

import (
	"fmt"
	"log"
	"time"

	"log2sqs/config"
	"log2sqs/global"
)

var discardTime int64 = 0

// Add log message to internal queue (buffer) for transmission to SQS
func Add(msg []byte) {

	if config.Debug {
		log.Printf("Buffer contains %d log events", len(eventBuffer))
	}

	// Check if the number of items in the buffer is at the limit
	if len(eventBuffer) >= config.EventBuffer {

		// Discard oldest message
		_ = <-eventBuffer

		// Limit logging this event to a maximum of once per minute to reduce flooding
		if (time.Now().Unix() - discardTime) > 60 {
			Log(fmt.Sprintf("Buffer full (%d items), discarding oldest log event", len(eventBuffer)), "", global.ERR)
			discardTime = time.Now().Unix()
		}
	}

	// Add to buffer
	eventBuffer <- msg
}
