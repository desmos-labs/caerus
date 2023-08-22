package chain

import (
	"fmt"
	"strings"

	wallettypes "github.com/desmos-labs/cosmos-go-wallet/types"

	"github.com/desmos-labs/caerus/utils"
)

type Config struct {
	Account *wallettypes.AccountConfig
	Chain   *wallettypes.ChainConfig
}

// Validate validates the given configuration returning any error
func (c *Config) Validate() error {
	if strings.TrimSpace(c.Account.Mnemonic) == "" {
		return fmt.Errorf("missing account mnemonic")
	}

	return nil
}

// ReadConfigFromEnvVariables reads a Config instance from the env variables values
func ReadConfigFromEnvVariables() (*Config, error) {
	cfg := &Config{
		Account: &wallettypes.AccountConfig{
			Mnemonic: utils.GetEnvOr(EnvChainAccountRecoveryPhrase, ""),
			HDPath:   utils.GetEnvOr(EnvChainAccountDerivationPath, "m/44'/852'/0'/0/0"),
		},
		Chain: &wallettypes.ChainConfig{
			Bech32Prefix:  utils.GetEnvOr(EnChainBech32Prefix, "desmos"),
			RPCAddr:       utils.GetEnvOr(EnvChainRPCUrl, "https://rpc.morpheus.desmos.network:443"),
			GRPCAddr:      utils.GetEnvOr(EnvChainGRPCUrl, "https://grpc.morpheus.desmos.network:443"),
			GasPrice:      utils.GetEnvOr(EnvChainGasPrice, "0.01udaric"),
			GasAdjustment: 1.5,
		},
	}
	return cfg, cfg.Validate()
}
