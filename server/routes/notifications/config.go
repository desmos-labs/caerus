package notifications

import (
	"strconv"

	"github.com/desmos-labs/caerus/utils"
)

type Config struct {
	RequiresAuthentication bool `json:"requires_authentication" yaml:"requires_authentication" toml:"requires_authentication"`
}

// ReadConfigFromEnvVariables reads the configuration from the env variables
func ReadConfigFromEnvVariables() *Config {
	requiresAuth, _ := strconv.ParseBool(utils.GetEnvOr(EnvNotificationsCreationRequiresAuthorization, "false"))
	return &Config{
		RequiresAuthentication: requiresAuth,
	}
}
