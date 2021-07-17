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
const ProductVersion = "0.0.1"

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

	// Loop to allow recovery if an error occurs
	for {

		// Open the SQS queue
		openQueue()

		// Tail the file and write to the queue
		err := tailFile(config.InputFile, tail.Config{Follow: true, ReOpen: true})
		if err != nil {
			log.Printf("Error tailing file: %s", err.Error())
			log.Printf("Sleeping for 60 seconds...")
			time.Sleep(60 * time.Second)
		}
	}
}

// Open the SQS queue
func openQueue() {

	// Initialize queue - retry indefinitely
	for {
		err := queue.Open()
		if err != nil {
			log.Printf("Error opening queue: %s", err.Error())
			log.Printf("Sleeping for 60 seconds...")
			time.Sleep(60 * time.Second)
		} else {
			log.Print("Queue opened")
			return
		}
	}
}

// Tail the file and write to the queue
func tailFile(filename string, config tail.Config) error {

	// Map for arbitrary JSON
	var v map[string]interface{}

	// Set up to tail the file
	t, err := tail.TailFile(filename, config)
	if err != nil {
		return err
	}

	// Loop and read
	for line := range t.Lines {

		// Trim leading and trailing whitespace
		s := strings.TrimSpace(line.Text)

		// Only accept lines with valid JSON
		if json.Unmarshal([]byte(s), &v) == nil {
			err := queue.Send(s)
			if err != nil {
				return err
			}
		} else {
			log.Printf("JSON validation failed, ignoring: %s", s)
		}
	}

	err = t.Wait()
	if err != nil {
		return err
	}
	return nil
}

// Graceful exit
func appCleanup(sig os.Signal) {
	log.Printf("Exiting on signal: %v", sig)
	os.Exit(0)
}
