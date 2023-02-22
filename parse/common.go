//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package parse

import (
	"strconv"
	"strings"
)

type GELFMessage map[string]interface{}

// String2Int returns the integer contained in string or 0
func String2Int(s string) int {
	ret, err := strconv.Atoi(s)
	if err != nil {
		ret = 0
	}
	return ret
}

// EmptyString returns a cleaned string or "-" if empty
func EmptyString(s string) string {
	r := strings.TrimSpace(s)
	if r == "" {
		r = "-"
	}
	return r
}
