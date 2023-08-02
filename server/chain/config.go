package chain

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	wallettypes "github.com/desmos-labs/cosmos-go-wallet/types"

	"github.com/desmos-labs/caerus/utils"
)

type Config struct {
	Account  *wallettypes.AccountConfig
	Chain    *wallettypes.ChainConfig
	FeeGrant *FeeGrantConfig
}

type FeeGrantConfig struct {
	// GrantLimit specifies the maximum amount of coins that can be spent
	// by any allowance and will be updated as coins are spent. If it is
	// empty, there is no spend limit and any amount of coins can be spent.
	GrantLimit sdk.Coins

	// Messages for which the grantee has the access. If it is empty,
	// the granted address has access to all the messages.
	MsgTypes []string

	// Expiration specifies an optional time limit on the grant. After this
	// time, the grantee can no longer use the grant. If it is zero, there
	// is no expiration time.
	Expiration time.Duration
}

// Validate validates the given configuration returning any error
func (c *Config) Validate() error {
	if strings.TrimSpace(c.Account.Mnemonic) == "" {
		return fmt.Errorf("missing account mnemonic")
	}

	if c.FeeGrant != nil {
		for _, msgType := range c.FeeGrant.MsgTypes {
			if strings.HasPrefix(msgType, "/") {
				return fmt.Errorf("invalid message type: %s", msgType)
			}
		}
	}

	return nil
}

// ReadConfigFromEnvVariables reads a Config instance from the env variables values
func ReadConfigFromEnvVariables() (*Config, error) {
	grantLimit, err := sdk.ParseCoinsNormalized(utils.GetEnvOr(EnvFeeGrantAmountLimit, ""))
	if err != nil {
		return nil, fmt.Errorf("invalid %s: %s", EnvFeeGrantAmountLimit, err)
	}

	grantExpiration, err := time.ParseDuration(utils.GetEnvOr(EnvFeeGrantExpiration, "0"))
	if err != nil {
		return nil, fmt.Errorf("invalid %s: %s", EnvFeeGrantExpiration, err)
	}

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
		FeeGrant: &FeeGrantConfig{
			GrantLimit: grantLimit,
			MsgTypes:   strings.Split(utils.GetEnvOr(EnvFeeGrantMessagesTypes, ""), ","),
			Expiration: grantExpiration,
		},
	}
	return cfg, cfg.Validate()
}
