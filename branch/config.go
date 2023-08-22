package branch

import (
	"fmt"

	"github.com/desmos-labs/caerus/utils"
)

// Config contains the data used to build a new Client instance
type Config struct {
	ApiKey string
}

// NewConfig allows to build a new Config instance
func NewConfig(apiKey string) *Config {
	return &Config{
		ApiKey: apiKey,
	}
}

// NewConfigFromEnvVariables allows to build a new Config instance from the environment variables
func NewConfigFromEnvVariables() (*Config, error) {
	apiKey := utils.GetEnvOr(EnvBranchKey, "")
	if apiKey == "" {
		return nil, fmt.Errorf("missing %s en variable", EnvBranchKey)
	}

	return NewConfig(
		apiKey,
	), nil
}
