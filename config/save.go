//
// Copyright (c) 2021-2023 Tenebris Technologies Inc.
//

package config

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

// Save the current configuration
func Save(configFile string) error {
	saveData, err := yaml.Marshal(Config)
	if err != nil {
		return errors.New(fmt.Sprintf("error marshaling configuration: %s", err.Error()))
	}

	err = os.WriteFile(configFile, saveData, 0644)
	if err != nil {
		return errors.New(fmt.Sprintf("error writing configuration to %s: %s", configFile, err.Error()))
	}

	return nil
}
