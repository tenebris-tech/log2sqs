//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type InputFileDef struct {
	Index int
	Name  string
	Type  string
}

var Debug = false
var AWSID = ""
var AWSKey = ""
var AWSRegion = ""
var AWSQueueName = ""
var AddEC2Tags = false
var LogFile = ""
var Hostname = ""
var SyslogUDP = ""
var SyslogUDPMax = 2048
var SyslogFullMessage = false
var SyslogOverrideTime = false
var EventBuffer = 4096
var InputFiles []InputFileDef
var AddFields = map[string]string{}

func Load(filename string) error {
	var item []string
	var name string
	var value string
	var index = 0
	var err error

	// Set default
	Hostname, err = os.Hostname()
	if err != nil {
		Hostname = ""
	}

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
		case "debug":
			Debug = string2bool(value)
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
			n, err := parseInputFile(value)
			if err != nil {
				tmp := fmt.Sprintf("error parsing config file: %s, %s", line, err.Error())
				return errors.New(tmp)
			} else {
				index++
				n.Index = index
				InputFiles = append(InputFiles, n)
			}
		case "hostname":
			Hostname = value
		case "site":
			AddFields["_site"] = value
		case "syslogudp":
			SyslogUDP = value
		case "syslogudpmax":
			tmp := string2int(value)
			if tmp > 0 {
				SyslogUDPMax = tmp
			}
		case "syslogfullmessage":
			SyslogFullMessage = string2bool(value)
		case "syslogoverridetime":
			SyslogOverrideTime = string2bool(value)
		case "eventbuffer":
			tmp := string2int(value)
			if tmp > 0 {
				EventBuffer = tmp
			}
		default:
			tmp := fmt.Sprintf("error parsing config file: %s", line)
			return errors.New(tmp)
		}
	}
	return nil
}

// Parse input file into filename and type
func parseInputFile(value string) (InputFileDef, error) {
	var f InputFileDef
	s := strings.Split(value, ",")
	if len(s) != 2 {
		return f, errors.New("must have two elements (path and type)")
	}
	f.Name = s[0]
	f.Type = s[1]
	return f, nil
}

// Return true if string is yes or true (case-insensitive)
func string2bool(s string) bool {
	if strings.ToLower(s) == "yes" {
		return true
	}

	if strings.ToLower(s) == "true" {
		return true
	}

	return false
}

// Return integer contained in string or 0
func string2int(s string) int {
	ret, err := strconv.Atoi(s)
	if err != nil {
		ret = 0
	}
	return ret
}
