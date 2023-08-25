//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package parse

import (
	"log2sqs/config"
	"regexp"
)

// GELFMessage type can hold GELF fields of various types
type GELFMessage map[string]interface{}

// Parser a struct that indicates hot to parse a log format
type Parser struct {
	format        string                                     // name of the format
	parserFunc    func(string, *Parser) (GELFMessage, error) // pointer to the parser function
	regexFields   config.RegexFields                         // map of regexes for each Field to put them in the correct order
	requireFields int                                        // number of fields required to be present
	regex         *regexp.Regexp                             // pointer to the compiled regex
}

// Parse formats
const (
	formatGelf                       = "gelf"
	formatText                       = "text"
	formatApacheError                = "error"
	formatApacheCombined             = "combined"
	formatApacheCombinedPlus         = "combinedplus"
	formatApacheCombinedPlusVhost    = "combinedplusvhost"
	formatApacheCombinedLoadBalancer = "combinedloadbalancer"
)

var Parsers = map[string]Parser{
	formatGelf:                       {format: formatGelf, parserFunc: gelfParser},
	formatText:                       {format: formatText, parserFunc: plainTextParser},
	formatApacheError:                {format: formatApacheError, parserFunc: regexParser, regexFields: apacheErrorRegex, requireFields: 5},
	formatApacheCombined:             {format: formatApacheCombined, parserFunc: regexParser, regexFields: apacheCombinedRegex, requireFields: 9},
	formatApacheCombinedPlus:         {format: formatApacheCombinedPlus, parserFunc: regexParser, regexFields: apacheCombinedPlusRegex, requireFields: 13},
	formatApacheCombinedPlusVhost:    {format: formatApacheCombinedPlusVhost, parserFunc: regexParser, regexFields: apacheCombinedPlusVhostRegex, requireFields: 15},
	formatApacheCombinedLoadBalancer: {format: formatApacheCombinedLoadBalancer, parserFunc: regexParser, regexFields: apacheCombinedLoadBalancerRegex, requireFields: 17},
}

var apacheErrorRegex = config.RegexFields{
	1: {Regex: `^\[([^]]+)\]\s`, Field: "timestamp", FType: "date", DateFormat: "Mon Jan 02 15:04:05.000000 2006 -0700", AddTZ: true},
	2: {Regex: `(\S+):`, Field: "_apache_module", FType: "string"},
	3: {Regex: `(\S+)\s`, Field: "_apache_level", FType: "string"},
	4: {Regex: `\[([^]]+)\]\s`, Field: "_apache_pid", FType: "string"},
	5: {Regex: `(.*?)$`, Field: "short_message", FType: "string"},
}

// LogFormat "%h %l %u %t \"%r\" %>s %O \"%{Referer}i\" \"%{User-Agent}i\"" combined
var apacheCombinedRegex = config.RegexFields{
	1: {Regex: `^(\S+)\s`, Field: "_src_ip", FType: "string"},
	2: {Regex: `(\S+)\s`, Field: "_http_ident", FType: "string"},
	3: {Regex: `(\S+)\s`, Field: "_user", FType: "string"},
	4: {Regex: `\[([^]]+)\]\s`, Field: "timestamp", FType: "date", DateFormat: "02/Jan/2006:15:04:05 -0700"},
	5: {Regex: `"(.*?)"\s`, Field: "_http_request", FType: "string", ShortMessage: true},
	6: {Regex: `(\S+)\s`, Field: "_http_status", FType: "int"},
	7: {Regex: `(\S+)\s`, Field: "_http_response_size", FType: "int"},
	8: {Regex: `"((?:[^"]*(?:\\")?)*)"\s`, Field: "_http_referer", FType: "string"},
	9: {Regex: `"(.*?)"$`, Field: "_user_agent", FType: "string"},
}

// LogFormat "%h %l %u %t \"%r\" %>s %O \"%{Referer}i\" \"%{User-Agent}i\" %D \"%m\" \"%U\" \"%q\"" combinedplus
var apacheCombinedPlusRegex = config.RegexFields{
	1:  {Regex: `^(\S+)\s`, Field: "_src_ip", FType: "string"},
	2:  {Regex: `(\S+)\s`, Field: "_http_ident", FType: "string"},
	3:  {Regex: `(\S+)\s`, Field: "_user", FType: "string"},
	4:  {Regex: `\[([^]]+)\]\s`, Field: "timestamp", FType: "date", DateFormat: "02/Jan/2006:15:04:05 -0700"},
	5:  {Regex: `"(.*?)"\s`, Field: "_http_request", FType: "string", ShortMessage: true},
	6:  {Regex: `(\S+)\s`, Field: "_http_status", FType: "int"},
	7:  {Regex: `(\S+)\s`, Field: "_http_response_size", FType: "int"},
	8:  {Regex: `"((?:[^"]*(?:\\")?)*)"\s`, Field: "_http_referer", FType: "string"},
	9:  {Regex: `"(.*?)"\s`, Field: "_user_agent", FType: "string"},
	10: {Regex: `(\S+)\s`, Field: "_duration_usec", FType: "int"},
	11: {Regex: `"(.*?)"\s`, Field: "_http_request_method", FType: "string"},
	12: {Regex: `"(.*?)"\s`, Field: "_http_request_path", FType: "string"},
	13: {Regex: `"(.*?)"$`, Field: "_http_request_query", FType: "string"},
}

// LogFormat "%v:%p %h %l %u %t \"%r\" %>s %O \"%{Referer}i\" \"%{User-Agent}i\" %D \"%m\" \"%U\" \"%q\"" combinedplusvhost
var apacheCombinedPlusVhostRegex = config.RegexFields{
	1:  {Regex: `$(\S+):`, Field: "_vhost", FType: "string"},
	2:  {Regex: `(\S+)\s`, Field: "_vhost_port", FType: "int"},
	3:  {Regex: `(\S+)\s`, Field: "_src_ip", FType: "string"},
	4:  {Regex: `(\S+)\s`, Field: "_http_ident", FType: "string"},
	5:  {Regex: `(\S+)\s`, Field: "_user", FType: "string"},
	6:  {Regex: `\[([^]]+)\]\s`, Field: "timestamp", FType: "date", DateFormat: "02/Jan/2006:15:04:05 -0700"},
	7:  {Regex: `"(.*?)"\s`, Field: "_http_request", FType: "string", ShortMessage: true},
	8:  {Regex: `(\S+)\s`, Field: "_http_status", FType: "int"},
	9:  {Regex: `(\S+)\s`, Field: "_http_response_size", FType: "int"},
	10: {Regex: `"((?:[^"]*(?:\\")?)*)"\s`, Field: "_http_referer", FType: "string"},
	11: {Regex: `"(.*?)"\s`, Field: "_user_agent", FType: "string"},
	12: {Regex: `(\S+)\s`, Field: "_duration_usec", FType: "int"},
	13: {Regex: `"(.*?)"\s`, Field: "_http_request_method", FType: "string"},
	14: {Regex: `"(.*?)"\s`, Field: "_http_request_path", FType: "string"},
	15: {Regex: `"(.*?)"$`, Field: "_http_request_query", FType: "string"},
}

// LogFormat "%{X-Forwarded-Proto}i %{Host}i:%{X-Forwarded-Port}i %v:%p %{X-Forwarded-For}i %h %t \"%r\" %>s %O \"%{Referer}i\" \"%{User-Agent}i\" %D \"%m\" \"%U\" \"%q\"" combinedloadbalancer
var apacheCombinedLoadBalancerRegex = config.RegexFields{
	1:  {Regex: `^(\S+)\s`, Field: "_x-forwarded-proto", FType: "string"},
	2:  {Regex: `(\S+):`, Field: "_http_host", FType: "string"},
	3:  {Regex: `(\S+)\s`, Field: "_x-forwarded-port", FType: "string"},
	4:  {Regex: `(\S+):`, Field: "_vhost", FType: "string"},
	5:  {Regex: `(\S+)\s`, Field: "_vhost_port", FType: "int"},
	6:  {Regex: `(\S+)\s`, Field: "_x-forwarded-for", FType: "string"},
	7:  {Regex: `(\S+)\s`, Field: "_src_ip", FType: "string"},
	8:  {Regex: `\[([^]]+)\]\s`, Field: "timestamp", FType: "date", DateFormat: "02/Jan/2006:15:04:05 -0700"},
	9:  {Regex: `"(.*?)"\s`, Field: "_http_request", FType: "string", ShortMessage: true},
	10: {Regex: `(\S+)\s`, Field: "_http_status", FType: "int"},
	11: {Regex: `(\S+)\s`, Field: "_http_response_size", FType: "int"},
	12: {Regex: `"((?:[^"]*(?:\\")?)*)"\s`, Field: "_http_referer", FType: "string"},
	13: {Regex: `"(.*?)"\s`, Field: "_user_agent", FType: "string"},
	14: {Regex: `(\S+)\s`, Field: "_duration_usec", FType: "int"},
	15: {Regex: `"(.*?)"\s`, Field: "_http_request_method", FType: "string"},
	16: {Regex: `"(.*?)"\s`, Field: "_http_request_path", FType: "string"},
	17: {Regex: `"(.*?)"$`, Field: "_http_request_query", FType: "string"},
}
