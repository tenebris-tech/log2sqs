//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package global

import "time"

//goland:noinspection GoUnusedConst
const (
	EMERG = iota
	ALERT
	CRIT
	ERR
	WARN
	NOTICE
	INFO
	DEBUG
)

//goland:noinspection GoUnusedConst
const (
	KERN = iota
	USER
	MAIL
	DAEMON
	AUTH
	SYSLOG
	LPR
	NEWS
	UUCP
	CRON
	AUTHPRIV
	FTP
	NTP
	SECURITY
	CONSOLE
	SOLARISCRON
	LOCAL0
	LOCAL1
	LOCAL2
	LOCAL3
	LOCAL4
	LOCAL5
	LOCAL6
	LOCAL7
)

var levels = map[int]string{
	0: "emerg",
	1: "alert",
	2: "crit",
	3: "err",
	4: "warning",
	5: "notice",
	6: "info",
	7: "debug",
}

var facilities = map[int]string{
	0:  "kern",
	1:  "user",
	2:  "mail",
	3:  "daemon",
	4:  "auth",
	5:  "syslog",
	6:  "lpr",
	7:  "news",
	8:  "uucp",
	9:  "cron",
	10: "authpriv",
	11: "ftp",
	12: "ntp",
	13: "security",
	14: "console",
	15: "solaris-cron",
	16: "local0",
	17: "local1",
	18: "local2",
	19: "local3",
	20: "local4",
	21: "local5",
	22: "local6",
	23: "local7",
}

// GetLevel returns "unknown" or the appropriate string
//
//goland:noinspection GoUnusedExportedFunction
func GetLevel(level int) string {
	s, ok := levels[level]
	if ok {
		return s
	}
	return "unknown"
}

// GetFacility returns "unknown" or the appropriate string
//
//goland:noinspection GoUnusedExportedFunction
func GetFacility(fac int) string {
	s, ok := facilities[fac]
	if ok {
		return s
	}
	return "unknown"
}

func TimeStamp() float64 {
	return float64(time.Now().UnixMilli()) / 1000.0
}
