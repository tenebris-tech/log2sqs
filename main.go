//
// Copyright (c) 2021 Tenebris Technologies Inc.
//

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log2sqs/config"
	"log2sqs/queue"
)

const ProductName = "log2sqs"
const ProductVersion = "0.1.0"

// Map to store addFields for adding to log events
var addFields = map[string]string{}

// main entry point
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

	// Retrieve EC2 addFields if necessary
	if config.AddEC2Tags {
		ec2Tags()
	}

	// Log fields to be added
	for key, value := range addFields {
		log.Printf("Adding field %s=%s", key, value)
	}

	// Open the SQS queue
	openQueue()

	// Create buffered channel to allow queue failure notification and restart request
	// Capacity is set high to avoid blocking.
	ch := make(chan int, 1024)

	// Iterate over list of files to monitor
	// and launch a goroutine to handle for each
	for _, inputFile := range config.InputFiles {
		go tailFile(inputFile, ch)
	}

	// Wait for signal or process to be killed and
	// handle any messages on channel
	for {
		// Read message from channel
		// This is a blocking function, but we're waiting anyway
		msg := <-ch
		log.Printf("Received queue restart request from worker %d [%s %s]",
			msg, config.InputFiles[msg].Name, config.InputFiles[msg].Type)

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
			log.Printf("SQS queue %s opened", config.AWSQueueName)
			return
		}
	}
}

// Graceful exit
func appCleanup(sig os.Signal) {
	log.Printf("Exiting on signal: %v", sig)
	os.Exit(0)
}
