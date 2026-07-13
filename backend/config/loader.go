package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed providers.json
var embeddedProvidersJSON []byte

type ProviderDef struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	BaseURL           string   `json:"base_url"`
	RequiresCustomURL bool     `json:"requires_custom_url"`
	EnvKey            string   `json:"env_key"`
	AuthType          string   `json:"auth_type"` // "bearer" | "cookie" | "oauth" | "bearer_token"
	Prefixes          []string `json:"prefixes"`
}

type ProviderConfig struct {
	Providers []ProviderDef `json:"providers"`
}

var GlobalConfig ProviderConfig

// LoadProviders loads the provider list from the JSON embedded in the binary.
func LoadProviders() error {
	if err := json.Unmarshal(embeddedProvidersJSON, &GlobalConfig); err != nil {
		return fmt.Errorf("failed to parse embedded providers config: %v", err)
	}
	return nil
}
