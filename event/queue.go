//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package event

import (
	"log"
	"time"

	"log2sqs/config"
	"log2sqs/global"
)

// Buffered channel to queue events to be sent to SQS
var eventBuffer chan []byte

// initQueue creates the queue
func initQueue() {
	// Create buffer with a bit of extra space to avoid blocking
	eventBuffer = make(chan []byte, config.EventBuffer+10)
}

// runQueue reads the internal event buffer channel and writes to SQS
func runQueue() {
	bufferWarning := false

	for {
		bPercent := float64(len(eventBuffer)) / float64(config.EventBuffer)

		if bufferWarning {
			if bPercent < 0.6 {
				bufferWarning = false
				Log("Event buffer is now below 60%% full", "", global.INFO)
			}
		} else {
			if bPercent > 0.8 {
				bufferWarning = true
				Log("Event buffer is more than 80%% full", "", global.WARN)
			}
		}

		// This is blocking, which is fine
		msg := <-eventBuffer

		// Send to SQS
		err := sendSQS(msg)
		if err != nil {
			// Log error
			log.Printf("Error sending buffered syslog message to SQS: %s", err.Error())

			// Add the message back into the buffer to prevent loss
			Add(msg)

			// Wait 15 seconds before trying again
			log.Printf("Sleeping for 15 seconds...")
			time.Sleep(15 * time.Second)
		}
	}
}
