//
// Copyright (c) 2021 Tenebris Technologies Inc.
//

package main

import (
	"strconv"
	"strings"
)

// Clean string and return string or "-" if empty
func emptyString(s string) string {
	r := strings.TrimSpace(s)
	if r == "" {
		r = "-"
	}
	return r
}

// Convert string to integer
// Return 0 if conversion fails
func string2Int(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
