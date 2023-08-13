//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"log2sqs/config"
	"log2sqs/event"
	"log2sqs/global"
	"log2sqs/parse"
	"log2sqs/syslog"
)

func main() {

	// Default configuration file
	var configFile = "log2sqs.conf"

	// File to ingest for testing
	var ingest = ""

	// Command line arguments

	// Check for path to config file as only argument for backward compatibility
	if len(os.Args) == 2 {
		configFile = os.Args[1]
	} else {
		cF := flag.String("config", "log2sqs.conf", "configuration file")
		iG := flag.String("ingest", "", "ingest entire file in the specified format")
		flag.Parse()

		if *cF != "" {
			configFile = *cF
		}

		if *iG != "" {
			ingest = *iG
		}
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

	// Add field to report application name and version
	config.AddFields["_via_app"] = global.ProductName + " " + global.ProductVersion

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

	// Retrieve EC2 addFields if necessary
	if config.AddEC2Tags {
		ec2Tags()
	}

	// Initialize and start queues
	event.Start()

	// Send log event
	event.Log(fmt.Sprintf("Starting %s %s", global.ProductName, global.ProductVersion), "", global.INFO)

	// If command line ingest file specified, add it to the list
	if ingest != "" {
		ingestFile, err := config.ParseInputFile(ingest)
		if err != nil {
			event.Log(fmt.Sprintf("Unable to ingest specified file %s: %s", ingest, err.Error()), "", global.INFO)
			os.Exit(1)
		}
		ingestFile.ReadAll = true
		config.InputFiles = append(config.InputFiles, ingestFile)
	}

	// Iterate over list of files to monitor
	for _, inputFile := range config.InputFiles {

		// Force all file types to lower case
		inputFile.Type = strings.ToLower(inputFile.Type)

		// Check for valid file type
		if parse.CheckFormat(inputFile.Type) == false {
			event.Log(fmt.Sprintf("Unknown input file type: [%d]%s %s", inputFile.Index, inputFile.Name, inputFile.Type), "", global.INFO)
		} else {
			// Launch a goroutine to handle this file
			go tailFile(inputFile)
		}
	}

	// Start Syslog UDP if configured
	if config.SyslogUDP != "" {
		go syslog.UDP()
	}

	// Wait for signal or process to be killed
	//goland:noinspection GoInfiniteFor
	select {}
}

// Graceful exit
func appCleanup(sig os.Signal) {
	event.Log(fmt.Sprintf("Exiting on signal: %v", sig), "", global.NOTICE)
	os.Exit(0)
}
