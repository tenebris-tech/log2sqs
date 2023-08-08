//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package main

import (
	"encoding/json"
	"io"
	"log"
	"strings"
	"time"

	"github.com/tenebris-tech/tail"

	"log2sqs/config"
	"log2sqs/event"
	"log2sqs/parse"
)

// Tail the file and write to the queue
func tailFile(f config.InputFileDef) {
	var send bool
	var g parse.GELFMessage

	// Infinite loop to facilitate restart on error
	for {
		// Set up to tail the file, starting at the current end
		t, err := tail.TailFile(
			f.Name,
			tail.Config{Follow: true, ReOpen: true, Location: &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd}})

		if err != nil {
			log.Printf("Error tailing file: %s [%d %s %s]", err.Error(), f.Index, f.Name, f.Type)
			log.Printf("Sleeping for 60 seconds...")
			time.Sleep(60 * time.Second)
		}

		// Loop and read
		for line := range t.Lines {

			// Create an empty message
			g = parse.GELFMessage{}

			// Trim leading and trailing whitespace
			s := strings.TrimSpace(line.Text)

			// Assume no data to send to queue
			send = false

			// Handle different file types here
			switch strings.ToLower(f.Type) {
			case "gelf":
				// Unmarshal to verify JSON and allow adding fields
				if json.Unmarshal([]byte(s), &g) == nil {
					send = true
				}

			case "text":
				err := parse.PlainText(s, g)
				if err != nil {
					log.Printf("Error parsing text log format: %s", err.Error())
				} else {
					send = true
				}

			case "combined":
				// Apache/NGINX combined log format
				err := parse.ApacheCombined(s, g)
				if err != nil {
					log.Printf("Error parsing combined log format: %s", err.Error())
				} else {
					send = true
				}

			case "combinedplus":
				// Apache combined log format with additional fields
				err := parse.ApacheCombinedPlus(s, g)
				if err != nil {
					log.Printf("Error parsing combinedplus log format: %s", err.Error())
				} else {
					send = true
				}

			case "combinedplusvhost":
				// Apache combined log format with additional fields
				err := parse.ApacheCombinedPlusVhost(s, g)
				if err != nil {
					log.Printf("Error parsing combinedplusvhost log format: %s", err.Error())
				} else {
					send = true
				}

			case "error":
				// Apache error log format
				err := parse.ApacheError(s, g)
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
				g["_log_file"] = f.Name
				g["_log_source"] = config.Hostname

				// Do we have addFields to add?
				for key, value := range config.AddFields {
					g[key] = value
				}

				// Marshal JSON for queue
				gBytes, err := json.Marshal(g)
				if err != nil {
					log.Printf("Failed to marshal JSON %s [%d %s %s]", err.Error(), f.Index, f.Name, f.Type)
					// Drop this log event
					continue
				}

				// Loop until line is sent to allow retries in the event of a failure
				// Since these are log files, there is no need to buffer them in memory
				sent := false
				for sent == false {
					err := event.Send(gBytes)
					if err != nil {
						// Log error
						log.Printf("Error sending to queue: %s [%d %s %s]", err.Error(), f.Index, f.Name, f.Type)
						log.Printf("Sleeping for 30 seconds...")
						time.Sleep(30 * time.Second)
					} else {
						sent = true
					}
				} // send to queue loop
			} else {
				log.Printf("Validation failed, ignoring: %s [%d %s %s]", s, f.Index, f.Name, f.Type)
			}
		}

		// For loop fell through. If there is an error, wait and restart the tail.
		err = t.Wait()
		if err != nil {
			log.Printf("Wait error: %s [%d %s %s]", err.Error(), f.Index, f.Name, f.Type)
			log.Printf("Sleeping for 60 seconds...")
			time.Sleep(60 * time.Second)
		}
	}
}
