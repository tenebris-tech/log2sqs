//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package parse

import "regexp"

// GELFMessage type can hold GELF fields of various types
type GELFMessage map[string]interface{}

// Parser a struct that indicates hot to parse a log format
type Parser struct {
	format        string                                     // name of the format
	parserFunc    func(string, *Parser) (GELFMessage, error) // pointer to the parser function
	regexFields   map[int]regexField                         // map of regexes for each field to put them in the correct order
	requireFields int                                        // number of fields required to be present
	regex         *regexp.Regexp                             // pointer to the compiled regex
}

// RegexField describes how to parse each field and what to map it to
type regexField struct {
	regex        string // regex to match the field
	field        string // name of the field
	fType        string // type of the field
	shortMessage bool   // if true, the field will be used as the short_message in addition to the named field
	dateFormat   string // if the field is a date, this is the format to use for parsing
	addTZ        bool   // if true, add the +0000 timezone to the timestamp to deal with annoying Apache logs
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

var apacheErrorRegex = map[int]regexField{
	1: {regex: `^\[([^]]+)\]\s`, field: "timestamp", fType: "date", dateFormat: "Mon Jan 02 15:04:05.000000 2006 -0700", addTZ: true},
	2: {regex: `(\S+):`, field: "_apache_module", fType: "string"},
	3: {regex: `(\S+)\s`, field: "_apache_level", fType: "string"},
	4: {regex: `\[([^]]+)\]\s`, field: "_apache_pid", fType: "string"},
	5: {regex: `(.*?)$`, field: "short_message", fType: "string"},
}

// LogFormat "%h %l %u %t \"%r\" %>s %O \"%{Referer}i\" \"%{User-Agent}i\"" combined
var apacheCombinedRegex = map[int]regexField{
	1: {regex: `^(\S+)\s`, field: "_src_ip", fType: "string"},
	2: {regex: `(\S+)\s`, field: "_http_ident", fType: "string"},
	3: {regex: `(\S+)\s`, field: "_user", fType: "string"},
	4: {regex: `\[([^]]+)\]\s`, field: "timestamp", fType: "date", dateFormat: "02/Jan/2006:15:04:05 -0700"},
	5: {regex: `"(.*?)"\s`, field: "_http_request", fType: "string", shortMessage: true},
	6: {regex: `(\S+)\s`, field: "_http_status", fType: "int"},
	7: {regex: `(\S+)\s`, field: "_http_response_size", fType: "int"},
	8: {regex: `"((?:[^"]*(?:\\")?)*)"\s`, field: "_http_referer", fType: "string"},
	9: {regex: `"(.*?)"$`, field: "_user_agent", fType: "string"},
}

// LogFormat "%h %l %u %t \"%r\" %>s %O \"%{Referer}i\" \"%{User-Agent}i\" %D \"%m\" \"%U\" \"%q\"" combinedplus
var apacheCombinedPlusRegex = map[int]regexField{
	1:  {regex: `^(\S+)\s`, field: "_src_ip", fType: "string"},
	2:  {regex: `(\S+)\s`, field: "_http_ident", fType: "string"},
	3:  {regex: `(\S+)\s`, field: "_user", fType: "string"},
	4:  {regex: `\[([^]]+)\]\s`, field: "timestamp", fType: "date", dateFormat: "02/Jan/2006:15:04:05 -0700"},
	5:  {regex: `"(.*?)"\s`, field: "_http_request", fType: "string", shortMessage: true},
	6:  {regex: `(\S+)\s`, field: "_http_status", fType: "int"},
	7:  {regex: `(\S+)\s`, field: "_http_response_size", fType: "int"},
	8:  {regex: `"((?:[^"]*(?:\\")?)*)"\s`, field: "_http_referer", fType: "string"},
	9:  {regex: `"(.*?)"\s`, field: "_user_agent", fType: "string"},
	10: {regex: `(\S+)\s`, field: "_duration_usec", fType: "int"},
	11: {regex: `"(.*?)"\s`, field: "_http_request_method", fType: "string"},
	12: {regex: `"(.*?)"\s`, field: "_http_request_path", fType: "string"},
	13: {regex: `"(.*?)"$`, field: "_http_request_query", fType: "string"},
}

// LogFormat "%v:%p %h %l %u %t \"%r\" %>s %O \"%{Referer}i\" \"%{User-Agent}i\" %D \"%m\" \"%U\" \"%q\"" combinedplusvhost
var apacheCombinedPlusVhostRegex = map[int]regexField{
	1:  {regex: `$(\S+):`, field: "_vhost", fType: "string"},
	2:  {regex: `(\S+)\s`, field: "_vhost_port", fType: "int"},
	3:  {regex: `(\S+)\s`, field: "_src_ip", fType: "string"},
	4:  {regex: `(\S+)\s`, field: "_http_ident", fType: "string"},
	5:  {regex: `(\S+)\s`, field: "_user", fType: "string"},
	6:  {regex: `\[([^]]+)\]\s`, field: "timestamp", fType: "date", dateFormat: "02/Jan/2006:15:04:05 -0700"},
	7:  {regex: `"(.*?)"\s`, field: "_http_request", fType: "string", shortMessage: true},
	8:  {regex: `(\S+)\s`, field: "_http_status", fType: "int"},
	9:  {regex: `(\S+)\s`, field: "_http_response_size", fType: "int"},
	10: {regex: `"((?:[^"]*(?:\\")?)*)"\s`, field: "_http_referer", fType: "string"},
	11: {regex: `"(.*?)"\s`, field: "_user_agent", fType: "string"},
	12: {regex: `(\S+)\s`, field: "_duration_usec", fType: "int"},
	13: {regex: `"(.*?)"\s`, field: "_http_request_method", fType: "string"},
	14: {regex: `"(.*?)"\s`, field: "_http_request_path", fType: "string"},
	15: {regex: `"(.*?)"$`, field: "_http_request_query", fType: "string"},
}

// LogFormat "%{X-Forwarded-Proto}i %{Host}i:%{X-Forwarded-Port}i %v:%p %{X-Forwarded-For}i %h %t \"%r\" %>s %O \"%{Referer}i\" \"%{User-Agent}i\" %D \"%m\" \"%U\" \"%q\"" combinedloadbalancer
var apacheCombinedLoadBalancerRegex = map[int]regexField{
	1:  {regex: `^(\S+)\s`, field: "_x-forwarded-proto", fType: "string"},
	2:  {regex: `(\S+):`, field: "_http_host", fType: "string"},
	3:  {regex: `(\S+)\s`, field: "_x-forwarded-port", fType: "string"},
	4:  {regex: `(\S+):`, field: "_vhost", fType: "string"},
	5:  {regex: `(\S+)\s`, field: "_vhost_port", fType: "int"},
	6:  {regex: `(\S+)\s`, field: "_x-forwarded-for", fType: "string"},
	7:  {regex: `(\S+)\s`, field: "_src_ip", fType: "string"},
	8:  {regex: `\[([^]]+)\]\s`, field: "timestamp", fType: "date", dateFormat: "02/Jan/2006:15:04:05 -0700"},
	9:  {regex: `"(.*?)"\s`, field: "_http_request", fType: "string", shortMessage: true},
	10: {regex: `(\S+)\s`, field: "_http_status", fType: "int"},
	11: {regex: `(\S+)\s`, field: "_http_response_size", fType: "int"},
	12: {regex: `"((?:[^"]*(?:\\")?)*)"\s`, field: "_http_referer", fType: "string"},
	13: {regex: `"(.*?)"\s`, field: "_user_agent", fType: "string"},
	14: {regex: `(\S+)\s`, field: "_duration_usec", fType: "int"},
	15: {regex: `"(.*?)"\s`, field: "_http_request_method", fType: "string"},
	16: {regex: `"(.*?)"\s`, field: "_http_request_path", fType: "string"},
	17: {regex: `"(.*?)"$`, field: "_http_request_query", fType: "string"},
}
