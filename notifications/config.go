package notifications

import (
	"github.com/desmos-labs/caerus/utils"
)

type Config struct {
	CredentialsFilePath string `json:"credentials_file_path" yaml:"credentials_file_path" toml:"credentials_file_path"`
}

func ReadConfigFromEnvVariables() (*Config, error) {
	return &Config{
		CredentialsFilePath: utils.GetEnvOr(EnvCredentialsFilePath, ""),
	}, nil
}
