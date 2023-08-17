package chain

import (
	"context"
	"fmt"
	"strings"

	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/desmos-labs/cosmos-go-wallet/client"
	"github.com/desmos-labs/cosmos-go-wallet/wallet"
)

type Client struct {
	*wallet.Wallet
	FeeGrantConfig *FeeGrantConfig
	authzClient    authz.QueryClient
	bankClient     banktypes.QueryClient
	feegrantClient feegrant.QueryClient
}

// NewClient returns a new client instance
func NewClient(cfg *Config, txConfig cosmosclient.TxConfig, cdc codec.Codec) (*Client, error) {
	walletClient, err := client.NewClient(cfg.Chain, cdc)
	if err != nil {
		return nil, fmt.Errorf("error while creating wallet client: %s", err)
	}

	cosmosWallet, err := wallet.NewWallet(cfg.Account, walletClient, txConfig)
	if err != nil {
		return nil, fmt.Errorf("error while creating cosmos wallet: %s", err)
	}

	return &Client{
		Wallet:         cosmosWallet,
		FeeGrantConfig: cfg.FeeGrant,
		authzClient:    authz.NewQueryClient(walletClient.GRPCConn),
		bankClient:     banktypes.NewQueryClient(walletClient.GRPCConn),
		feegrantClient: feegrant.NewQueryClient(walletClient.GRPCConn),
	}, nil
}

// NewClientFromEnvVariables builds a new Client instance by reading the configuration from the environment variables
func NewClientFromEnvVariables(txConfig cosmosclient.TxConfig, cdc codec.Codec) (*Client, error) {
	config, err := ReadConfigFromEnvVariables()
	if err != nil {
		return nil, err
	}

	return NewClient(config, txConfig, cdc)
}

// --------------------------------------------------------------------------------------------------------------------

// HasGrantedMsgGrantAllowanceAuthorization checks whether the given address has granted a
// MsgGrantAllowance authorization to the wallet of this client
func (c *Client) HasGrantedMsgGrantAllowanceAuthorization(appAddress string) (bool, error) {
	res, err := c.authzClient.Grants(context.Background(), &authz.QueryGrantsRequest{
		Granter:    appAddress,
		Grantee:    c.Wallet.AccAddress(),
		MsgTypeUrl: sdk.MsgTypeURL(&feegrant.MsgGrantAllowance{}),
	})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		}
		return false, err
	}

	return res.Grants != nil && len(res.Grants) > 0, nil
}

// --------------------------------------------------------------------------------------------------------------------

// HasFeeGrant checks whether the given address already has a fee grant from the given granter
func (c *Client) HasFeeGrant(userAddress string, granterAddress string) (bool, error) {
	res, err := c.feegrantClient.Allowance(context.Background(), &feegrant.QueryAllowanceRequest{
		Granter: granterAddress,
		Grantee: userAddress,
	})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return false, nil
		}
		return false, err
	}

	return res.Allowance != nil, nil
}

// HasFunds checks whether the given address has funds or not inside their account
func (c *Client) HasFunds(address string) (bool, error) {
	res, err := c.bankClient.Balance(context.Background(), &banktypes.QueryBalanceRequest{
		Address: address,
		Denom:   c.Client.GasPrice.Denom,
	})
	if err != nil {
		return false, err
	}

	return res.Balance != nil && !res.Balance.IsZero(), nil
}
