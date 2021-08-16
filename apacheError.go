//
// Copyright (c) 2021 Tenebris Technologies Inc.
//

package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

// Parse Apache2 error log format into object
func apacheError(s string, j arbitraryJSON) error {
	fmt.Println(s)
	var buffer bytes.Buffer
	buffer.WriteString(`\[([^]]+)\]\s`) // date, time(1)
	buffer.WriteString(`\[([^]]+)\]\s`) // module:level (2)
	buffer.WriteString(`\[([^]]+)\]\s`) // pid (3)
	buffer.WriteString(`(.*?)$`)        // message (4)

	re1, err := regexp.Compile(buffer.String())
	if err != nil {
		return err
	}
	result := re1.FindStringSubmatch(s)

	fmt.Println(result[1])
	fmt.Println(result[2])
	fmt.Println(result[3])
	fmt.Println(result[4])

	if len(result) < 4 {
		return errors.New("too few fields found")
	}

	// Parse time
	// Apache annoyingly omits the timezone from error log messages,
	// so add UTC
	layout := "Mon Jan 02 15:04:05.000000 2006 -0700"
	t, err := time.Parse(layout, result[1]+" +0000")
	if err != nil {
		return err
	}

	// Parse module and error level
	module := "apache"
	level := "unknown"
	split := strings.Split(result[2], ":")
	if len(split) > 0 {
		module = split[0]
	}
	if len(split) > 1 {
		level = split[1]
	}

	// Create GELF message
	j["version"] = "1.1"
	j["host"], _ = os.Hostname()
	j["short_message"] = result[4]
	j["timestamp"] = t.Unix()
	j["_apache_pid"] = result[3]
	j["_apache_module"] = module
	j["_apache_level"] = level
	return nil
}
