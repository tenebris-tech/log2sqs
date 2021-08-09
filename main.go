//
// Copyright (c) 2021 Tenebris Technologies Inc.
//

package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"log2sqs/config"
	"log2sqs/queue"

	"github.com/tenebris-tech/tail"
)

const ProductName = "log2sqs"
const ProductVersion = "0.0.2"

func main() {

	// Default config file name
	var configFile = "log2sqs.conf"

	// Check for path to config file as argument
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	// Setup signal catching
	signals := make(chan os.Signal, 1)

	// Catch signals
	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// method invoked upon seeing signal
	go func() {
		for {
			s := <-signals
			appCleanup(s)
		}
	}()

	// Load configuration information
	err := config.Load(configFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Set up logging
	if config.LogFile != "" {
		f, err := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		// If unable to open log file, report error, but continue writing logs to stderr
		if err != nil {
			log.Printf("Error opening log file: %s", err.Error())
		} else {
			defer func(f *os.File) {
				_ = f.Close()
			}(f)
			log.SetOutput(f)
		}
	}

	log.Printf("Starting %s %s", ProductName, ProductVersion)

	// Open the SQS queue
	openQueue()

	// Create buffered channel to allow queue failure notification and restart request
	// Capacity is set high to avoid blocking.
	ch := make(chan int, 1024)

	// Iterate over list of files to monitor
	// and launch a goroutine to handle for each
	for i, inputFile := range config.InputFiles {
		go tailFile(i, inputFile, ch)
	}

	// Wait for signal or process to be killed and
	// handle any messages on channel
	for {
		// Read message from channel
		// This is a blocking function, but we're waiting anyway
		msg := <-ch
		log.Printf("Received queue restart request from worker %d [%s]", msg, config.InputFiles[msg])

		// Open the SQS queue again
		openQueue()

		// Drain the channel to avoid duplicates
		for len(ch) > 0 {
			_ = <-ch
		}
	}
}

// Open the SQS queue
func openQueue() {

	// Initialize queue - retry indefinitely until success
	for {
		err := queue.Open()
		if err != nil {
			log.Printf("Error opening queue: %s", err.Error())
			log.Printf("Sleeping for 60 seconds...")
			time.Sleep(60 * time.Second)
		} else {
			log.Printf("Queue %s opened", config.AWSQueueName)
			return
		}
	}
}

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

				// Loop until line in sent to allow retries in the event of a failure
				sent := false
				for sent == false {
					err := queue.Send(s)
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

// Graceful exit
func appCleanup(sig os.Signal) {
	log.Printf("Exiting on signal: %v", sig)
	os.Exit(0)
}
