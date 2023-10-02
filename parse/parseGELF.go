//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package parse

import "encoding/json"

// gelfParser parses a GELF message
func (p *Parser) gelfParser(s string) (GELFMessage, error) {
	g := GELFMessage{}
	err := json.Unmarshal([]byte(s), &g)
	if err != nil {
		// If error, return empty GELFMessage
		return GELFMessage{}, err
	}
	return g, nil
}
