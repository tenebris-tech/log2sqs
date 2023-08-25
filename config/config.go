//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package config

import "os"

type Data struct {
	Debug                  bool              `yaml:"Debug"`
	LogFile                string            `yaml:"LogFile"`
	AWSID                  string            `yaml:"AWSID"`
	AWSKey                 string            `yaml:"AWSKey"`
	AWSRegion              string            `yaml:"AWSRegion"`
	AWSQueueName           string            `yaml:"AWSQueueName"`
	AddEC2Tags             bool              `yaml:"AddEC2Tags"`
	Hostname               string            `yaml:"Hostname"`
	SyslogUDP              string            `yaml:"SyslogUDP"`
	SyslogUDPMax           int               `yaml:"SyslogUDPMax"`
	SyslogFullMessage      bool              `yaml:"SyslogFullMessage"`
	SyslogOverrideTime     bool              `yaml:"SyslogOverrideTime"`
	SyslogOverrideSourceIP string            `yaml:"SyslogOverrideSourceIP"`
	SyslogReplaceLocalhost bool              `yaml:"SyslogReplaceLocalhost"`
	EventBuffer            int               `yaml:"EventBuffer"`
	InputFiles             []InputFileDef    `yaml:"InputFiles"`
	AddFields              map[string]string `yaml:"AddFields"`
}

type InputFileDef struct {
	Name    string `yaml:"Name"`
	Type    string `yaml:"Type"`
	ReadAll bool   `yaml:"ReadAll"`
}

var Config Data

// SetDefaults sets the default values for the configuration
func SetDefaults() {
	var err error

	// Create map
	Config.AddFields = make(map[string]string)

	// Set hostname as a default
	Config.Hostname, err = os.Hostname()
	if err != nil {
		Config.Hostname = ""
	}

	// Set default values
	Config.Debug = false
	Config.AddEC2Tags = false
	Config.SyslogUDPMax = 2048
	Config.SyslogFullMessage = false
	Config.SyslogOverrideTime = false
	Config.SyslogReplaceLocalhost = false
	Config.EventBuffer = 4096
}
