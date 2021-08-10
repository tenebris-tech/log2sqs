//
// Copyright (c) 2021 Tenebris Technologies Inc.
//

package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var AWSID = ""
var AWSKey = ""
var AWSRegion = ""
var AWSQueueName = ""
var AddEC2Tags = false
var LogFile = ""
var InputFiles []string

func Load(filename string) error {
	var item []string
	var name string
	var value string

	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	//noinspection GoUnhandledErrorResult
	defer file.Close()

	lineCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		// Read line and increment line counter
		line := scanner.Text()
		lineCount++

		// Ignore empty lines and comments
		if len(line) < 1 {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, "/") {
			continue
		}

		// Split into name value pair
		item = strings.Split(line, "=")
		if len(item) < 2 {
			continue
		}

		name = strings.TrimSpace(strings.ToLower(item[0]))
		value = strings.TrimSpace(item[1])

		switch name {
		case "logfile":
			LogFile = value
		case "awsid":
			AWSID = value
		case "awskey":
			AWSKey = value
		case "awsregion":
			AWSRegion = value
		case "awsqueuename":
			AWSQueueName = value
		case "addec2tags":
			AddEC2Tags = string2bool(value)
		case "inputfile":
			// Append to list (slice)
			InputFiles = append(InputFiles, value)
		default:
			tmp := fmt.Sprintf("error parsing config file: %s", line)
			return errors.New(tmp)
		}
	}
	return nil
}

// Return true if string is yes or true (case insensitive)
func string2bool(s string) bool {

	if strings.ToLower(s) == "yes" {
		return true
	}

	if strings.ToLower(s) == "true" {
		return true
	}

	return false
}
