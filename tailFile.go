//
// Copyright (c) 2021 Tenebris Technologies Inc.
//

package main

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"log2sqs/queue"

	"github.com/tenebris-tech/tail"
)

// Tail the file and write to the queue
func tailFile(index int, filename string, ch chan int) {

	// Infinite loop to allow retry on error
	for {

		// Map for arbitrary JSON
		var v map[string]interface{}

		// Set up to tail the file
		t, err := tail.TailFile(filename, tail.Config{Follow: true, ReOpen: true})
		if err != nil {
			log.Printf("Error tailing file: %s [%d %s]", err.Error(), index, filename)
			log.Printf("Sleeping for 60 seconds...")
			time.Sleep(60 * time.Second)
		}

		// Loop and read
		for line := range t.Lines {

			// Trim leading and trailing whitespace
			s := strings.TrimSpace(line.Text)

			// Only accept lines with valid JSON
			if json.Unmarshal([]byte(s), &v) == nil {

				// Do we have addFields to add?
				for key, value := range addFields {
					v[key] = value
				}

				// Marshal JSON for queue
				jsonObj, err := json.Marshal(v)
				if err != nil {
					log.Printf("Failed to marshal JSON %s [%d %s]", err.Error(), index, filename)
					// Drop this log event
					continue
				}

				// Loop until line in sent to allow retries in the event of a failure
				sent := false
				for sent == false {
					err := queue.Send(string(jsonObj))
					if err != nil {
						// Log error
						log.Printf("Error sending to queue: %s [%d %s]", err.Error(), index, filename)
						log.Printf("Sending queue restart request for [%d %s]", index, filename)

						// Write our index to channel to request an SQS queue restart
						ch <- index

						// Wait 60 seconds before trying again
						log.Printf("Sleeping for 60 seconds...")
						time.Sleep(60 * time.Second)
					} else {
						sent = true
					}
				}
			} else {
				log.Printf("JSON validation failed, ignoring: %s [%d %s]", s, index, filename)
			}
		}

		err = t.Wait()
		if err != nil {
			log.Printf("Wait error: %s [%d %s]", err.Error(), index, filename)
			log.Printf("Sleeping for 60 seconds...")
			time.Sleep(60 * time.Second)
		}
	}
}
