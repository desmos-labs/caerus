package chain

import (
	"context"
	"fmt"
	"strings"
	"time"

	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/desmos-labs/cosmos-go-wallet/client"
	wallettypes "github.com/desmos-labs/cosmos-go-wallet/types"
	"github.com/desmos-labs/cosmos-go-wallet/wallet"

	"github.com/desmos-labs/caerus/types"
	"github.com/desmos-labs/caerus/utils"
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

// --------------------------------------------------------------------------------------------------------------------

// BuildFeeAllowance builds the fee allowance to be used when broadcasting the transaction
func (c *Client) BuildFeeAllowance() (feegrant.FeeAllowanceI, *types.FeeGrantAllowance, error) {
	// Get the expiration
	var expiration *time.Time
	if c.FeeGrantConfig.Expiration > 0 {
		expiration = utils.GetTimePointer(time.Now().Add(c.FeeGrantConfig.Expiration))
	}

	// Build the basic allowance
	var allowance feegrant.FeeAllowanceI = &feegrant.BasicAllowance{
		SpendLimit: c.FeeGrantConfig.GrantLimit,
		Expiration: expiration,
	}

	// Build the allowed message allowance
	if len(c.FeeGrantConfig.MsgTypes) != 0 {
		allowedMsgAllowance, err := feegrant.NewAllowedMsgAllowance(allowance, c.FeeGrantConfig.MsgTypes)
		if err != nil {
			return nil, nil, err
		}
		allowance = allowedMsgAllowance
	}

	return allowance, types.NewAuthorization(expiration, c.FeeGrantConfig.GrantLimit, c.FeeGrantConfig.MsgTypes), nil
}

// BroadcastFeeAllowancesTransaction broadcasts the transaction that grants the given fee allowance to the given addresses
func (c *Client) BroadcastFeeAllowancesTransaction(allowance feegrant.FeeAllowanceI, grantees []string) error {
	// Parse the addresses
	granterAddress, err := c.Client.ParseAddress(c.Wallet.AccAddress())
	if err != nil {
		return err
	}

	// Build the list of messages to be broadcast
	var msgs = make([]sdk.Msg, len(grantees))
	for i, grantee := range grantees {
		granteeAddress, err := c.Client.ParseAddress(grantee)
		if err != nil {
			return err
		}

		// Build the message
		feeGrantMsg, err := feegrant.NewMsgGrantAllowance(allowance, granterAddress, granteeAddress)
		if err != nil {
			return err
		}

		// Append the message to the list of ones that will be broadcast
		msgs[i] = feeGrantMsg
	}

	// Broadcast the transaction
	response, err := c.BroadcastTxSync(&wallettypes.TransactionData{Messages: msgs, GasAuto: true, FeeAuto: true})
	if err != nil {
		return err
	}

	// Check the response
	if response.Code != 0 {
		return fmt.Errorf("error while granting fee permission: %s", response.RawLog)
	}

	return nil
}

// GrantFeePermission grants a fee permission to the given address
func (c *Client) GrantFeePermission(address string) (*types.FeeGrantAllowance, error) {
	allowance, authorization, err := c.BuildFeeAllowance()
	if err != nil {
		return nil, err
	}

	err = c.BroadcastFeeAllowancesTransaction(allowance, []string{address})
	if err != nil {
		return nil, err
	}

	return authorization, nil
}

// GrantFeePermissions grants a fee permission to the given addresses
func (c *Client) GrantFeePermissions(addresses []string) (*types.FeeGrantAllowance, error) {
	allowance, authorization, err := c.BuildFeeAllowance()
	if err != nil {
		return nil, err
	}

	err = c.BroadcastFeeAllowancesTransaction(allowance, addresses)
	if err != nil {
		return nil, err
	}

	return authorization, nil
}
