//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

// Load the configuration
func Load(configFile string) error {

	// Set defaults
	SetDefaults()

	// Check for a legacy .conf file
	if len(configFile) > 5 && configFile[len(configFile)-5:] == ".conf" {
		return LoadLegacy(configFile)
	}

	// Load the config file
	content, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(content, &Config)
}
