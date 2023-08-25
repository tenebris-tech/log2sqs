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

func LoadLegacy(filename string) error {
	var item []string
	var name string
	var value string
	var err error

	// Open the file
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
			Config.Debug = string2bool(value)
		case "logfile":
			Config.LogFile = value
		case "awsid":
			Config.AWSID = value
		case "awskey":
			Config.AWSKey = value
		case "awsregion":
			Config.AWSRegion = value
		case "awsqueuename":
			Config.AWSQueueName = value
		case "addec2tags":
			Config.AddEC2Tags = string2bool(value)
		case "inputfile":
			// Append to list (slice)
			n, err := ParseInputFile(value)
			if err != nil {
				tmp := fmt.Sprintf("error parsing config file: %s, %s", line, err.Error())
				return errors.New(tmp)
			} else {
				Config.InputFiles = append(Config.InputFiles, n)
			}
		case "hostname":
			Config.Hostname = value
		case "site":
			Config.AddFields["_site"] = value
		case "syslogudp":
			Config.SyslogUDP = value
		case "syslogudpmax":
			tmp := string2int(value)
			if tmp > 0 {
				Config.SyslogUDPMax = tmp
			}
		case "syslogfullmessage":
			Config.SyslogFullMessage = string2bool(value)
		case "syslogoverridetime":
			Config.SyslogOverrideTime = string2bool(value)
		case "syslogreplacelocalhost":
			Config.SyslogReplaceLocalhost = string2bool(value)
		case "syslogoverridesourceip":
			Config.SyslogOverrideSourceIP = value
		case "eventbuffer":
			tmp := string2int(value)
			if tmp > 0 {
				Config.EventBuffer = tmp
			}
		default:
			tmp := fmt.Sprintf("error parsing config file: %s", line)
			return errors.New(tmp)
		}
	}
	return nil
}

// ParseInputFile converts the string into a filename and type
func ParseInputFile(value string) (InputFileDef, error) {
	var f InputFileDef
	s := strings.Split(value, ",")
	if len(s) != 2 {
		return f, errors.New("must have two elements (path and type)")
	}
	f.Name = s[0]
	f.Type = s[1]
	f.ReadAll = false
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
