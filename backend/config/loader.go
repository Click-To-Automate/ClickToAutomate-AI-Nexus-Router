package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type ProviderDef struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	BaseURL  string   `json:"base_url"`
	EnvKey   string   `json:"env_key"`
	Prefixes []string `json:"prefixes"`
}

type ProviderConfig struct {
	Providers []ProviderDef `json:"providers"`
}

var GlobalConfig ProviderConfig

// LoadProviders parses the providers.json file into memory
func LoadProviders(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read providers config file: %v", err)
	}

	if err := json.Unmarshal(data, &GlobalConfig); err != nil {
		return fmt.Errorf("failed to parse providers config file: %v", err)
	}

	return nil
}
