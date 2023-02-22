//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package main

import (
	"fmt"
	"log"
	"log2sqs/global"
	"os"
	"os/signal"
	"syscall"

	"log2sqs/config"
	"log2sqs/event"
	"log2sqs/syslog"
)

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

	// Iterate over list of files to monitor
	// and launch a goroutine to handle each
	for _, inputFile := range config.InputFiles {
		go tailFile(inputFile)
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
