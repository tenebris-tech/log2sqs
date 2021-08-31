//
// Copyright (c) 2021 Tenebris Technologies Inc.
//

package main

import (
	"encoding/json"
	"io"
	"log"
	"log2sqs/config"
	"strings"
	"time"

	"log2sqs/queue"

	"github.com/tenebris-tech/tail"
)

type arbitraryJSON map[string]interface{}

// Tail the file and write to the queue
func tailFile(f config.InputFileDef, ch chan int) {
	var send bool
	var j arbitraryJSON

	// Infinite loop to allow retry on error
	for {
		// Set up to tail the file, starting at the current end
		t, err := tail.TailFile(
			f.Name,
			tail.Config{Follow: true, ReOpen: true, Location: &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}})
		//tail.Config{Follow: true, ReOpen: true})
		if err != nil {
			log.Printf("Error tailing file: %s [%d %s %s]", err.Error(), f.Index, f.Name, f.Type)
			log.Printf("Sleeping for 60 seconds...")
			time.Sleep(60 * time.Second)
		}

		// Loop and read
		for line := range t.Lines {

			// Trim leading and trailing whitespace
			s := strings.TrimSpace(line.Text)

			// Assume no data to send to queue
			send = false

			// Handle different file types here
			switch strings.ToLower(f.Type) {
			case "gelf":
				// Unmarshal to verify JSON and allow adding fields
				if json.Unmarshal([]byte(s), &j) == nil {
					send = true
				}
			case "combined":
				// Apache/NGINX combined log format
				j = arbitraryJSON{}
				err := apacheCombined(s, j)
				if err != nil {
					log.Printf("Error parsing combined log format: %s", err.Error())
				} else {
					send = true
				}
			case "combinedplus":
				// Apache combined log format with additional fields
				j = arbitraryJSON{}
				err := apacheCombinedPlus(s, j)
				if err != nil {
					log.Printf("Error parsing combinedplus log format: %s", err.Error())
				} else {
					send = true
				}
			case "error":
				// Apache error log format
				j = arbitraryJSON{}
				err := apacheError(s, j)
				if err != nil {
					log.Printf("Error parsing Apache error log format: %s", err.Error())
				} else {
					send = true
				}
			default:
				// Should never get there, but just in case...
				log.Printf("Unknown log file type [%d %s %s]", f.Index, f.Name, f.Type)
			}

			if send {
				// Add filename
				j["_log_file"] = f.Name

				// Do we have addFields to add?
				for key, value := range addFields {
					j[key] = value
				}

				// Marshal JSON for queue
				jsonObj, err := json.Marshal(j)
				if err != nil {
					log.Printf("Failed to marshal JSON %s [%d %s %s]", err.Error(), f.Index, f.Name, f.Type)
					// Drop this log event
					continue
				}

				// Loop until line in sent to allow retries in the event of a failure
				sent := false
				for sent == false {
					err := queue.Send(string(jsonObj))
					if err != nil {
						// Log error
						log.Printf("Error sending to queue: %s [%d %s %s]", err.Error(), f.Index, f.Name, f.Type)
						log.Printf("Sending queue restart request for [%d %s %s]", f.Index, f.Name, f.Type)

						// Write our index to channel to request an SQS queue restart
						ch <- f.Index

						// Wait 60 seconds before trying again
						log.Printf("Sleeping for 60 seconds...")
						time.Sleep(60 * time.Second)
					} else {
						sent = true
					}
				} // send to queue loop
			} else {
				log.Printf("Validation failed, ignoring: %s [%d %s %s]", s, f.Index, f.Name, f.Type)
			}
		}

		err = t.Wait()
		if err != nil {
			log.Printf("Wait error: %s [%d %s %s]", err.Error(), f.Index, f.Name, f.Type)
			log.Printf("Sleeping for 60 seconds...")
			time.Sleep(60 * time.Second)
		}
	}
}
