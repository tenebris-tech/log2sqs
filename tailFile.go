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
	"log2sqs/global"
	"log2sqs/parse"
)

// Tail the file and write to the queue
func tailFile(f config.InputFileDef) {

	// Infinite loop to facilitate restart on error
	for {
		// Determine where we should start reading. By default, always start at the end to avoid
		// reprocessing old data. But, if ReadAll is set, start at the beginning.
		whence := io.SeekEnd
		if f.ReadAll {
			// Override for file ingestion
			whence = io.SeekStart
		}

		// Tail the file
		t, err := tail.TailFile(f.Name, tail.Config{Follow: true, ReOpen: true, Location: &tail.SeekInfo{Offset: 0, Whence: whence}})
		if err != nil {
			log.Printf("Error tailing file: %s [%d %s %s]", err.Error(), f.Index, f.Name, f.Type)
			log.Printf("Sleeping for 60 seconds...")
			time.Sleep(60 * time.Second)
		}

		// Loop and read
		for line := range t.Lines {

			// Trim leading and trailing whitespace and parse the line
			s := strings.TrimSpace(line.Text)
			g, err2 := parse.Parse(s, f.Type)
			if err2 != nil {
				log.Printf("error parsing %s: %s", s, err2.Error())
			} else {
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

				// For debugging only
				if config.Debug || dryRun {
					global.JSONPretty(gBytes)
				}

				if !dryRun {
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
					}
				}
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
