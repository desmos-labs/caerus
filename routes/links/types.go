package links

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/desmos-labs/caerus/types"
)

type GenerateGenericDeepLinkRequest struct {
	// ID of the app that is trying to generate the link
	AppID string

	// Configuration used to generate the link
	LinkConfig *types.LinkConfig

	// (optional) API key used to generate the link
	ApiKey string
}

func NewGenerateGenericDeepLinkRequest(appID string, config *types.LinkConfig, apiKey string) *GenerateGenericDeepLinkRequest {
	return &GenerateGenericDeepLinkRequest{
		AppID:      appID,
		LinkConfig: config,
		ApiKey:     apiKey,
	}
}

func (r *GenerateGenericDeepLinkRequest) Validate() error {
	if r.ApiKey == "" {
		return fmt.Errorf("missing API key")
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

type GenerateDeepLinkRequest interface {
	// GetAppID returns the id of the app that is trying to generate the link
	GetAppID() string

	// GetAction returns the action associated to the link
	GetAction() string

	// GetChainType returns the chain type that should be associated to the link
	GetChainType() ChainType

	// GetCustomData returns the custom data that should be associated to the link
	GetCustomData() map[string]string

	// Validate checks if the request is valid or not
	Validate() error
}

// --------------------------------------------------------------------------------------------------------------------

// GenerateAddressLinkRequest represents a request to generate a deep link that allows the user
// to decide what operations
type GenerateAddressLinkRequest struct {
	// ID of the app that is trying to generate the link
	AppID string

	// Address that should be associated to the link
	Address string

	// Chain represents the chain type that should be associated to the link
	Chain ChainType
}

func NewGenerateAddressLinkRequest(appID string, address string, chainType ChainType) *GenerateAddressLinkRequest {
	return &GenerateAddressLinkRequest{
		AppID:   appID,
		Address: address,
		Chain:   chainType,
	}
}

func (r *GenerateAddressLinkRequest) GetAppID() string {
	return r.AppID
}

func (r *GenerateAddressLinkRequest) GetChainType() ChainType {
	return r.Chain
}

func (r *GenerateAddressLinkRequest) GetAction() string {
	return ""
}

func (r *GenerateAddressLinkRequest) GetCustomData() map[string]string {
	return map[string]string{
		types.DeepLinkAddressKey: r.Address,
	}
}

func (r *GenerateAddressLinkRequest) Validate() error {
	if r.Address == "" {
		return fmt.Errorf("missing address")
	}

	if r.Chain == ChainType_UNDEFINED {
		return fmt.Errorf("invalid chain type")
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

type GenerateViewProfileLinkRequest struct {
	// ID of the app that is trying to generate the link
	AppID string

	// Address that should be associated to the link
	Address string

	// Chain represents the chain type that should be associated to the link
	Chain ChainType
}

func NewGenerateViewProfileLinkRequest(appID string, address string, chainType ChainType) *GenerateViewProfileLinkRequest {
	return &GenerateViewProfileLinkRequest{
		AppID:   appID,
		Address: address,
		Chain:   chainType,
	}
}

func (r *GenerateViewProfileLinkRequest) GetAppID() string {
	return r.AppID
}

func (r *GenerateViewProfileLinkRequest) GetAction() string {
	return types.DeepLinkActionViewProfile
}

func (r *GenerateViewProfileLinkRequest) GetChainType() ChainType {
	return r.Chain
}

func (r *GenerateViewProfileLinkRequest) GetCustomData() map[string]string {
	return map[string]string{
		types.DeepLinkAddressKey: r.Address,
	}
}

func (r *GenerateViewProfileLinkRequest) Validate() error {
	if r.Address == "" {
		return fmt.Errorf("missing address")
	}

	if r.Chain == ChainType_UNDEFINED {
		return fmt.Errorf("invalid chain type")
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

type GenerateSendTokensLinkRequest struct {
	// ID of the app that is trying to generate the link
	AppID string

	// Address that should be associated to the link
	Address string

	// Chain represents the chain type that should be associated to the link
	Chain ChainType

	// (optional) Amount of tokens that should be sent
	Amount sdk.Coins
}

func NewGenerateSendTokensLinkRequest(appID string, address string, chainType ChainType, amount sdk.Coins) *GenerateSendTokensLinkRequest {
	return &GenerateSendTokensLinkRequest{
		AppID:   appID,
		Address: address,
		Chain:   chainType,
		Amount:  amount,
	}
}

func (r *GenerateSendTokensLinkRequest) GetAppID() string {
	return r.AppID
}

func (r *GenerateSendTokensLinkRequest) GetAction() string {
	return types.DeepLinkActionSendTokens
}

func (r *GenerateSendTokensLinkRequest) GetChainType() ChainType {
	return r.Chain
}

func (r *GenerateSendTokensLinkRequest) GetCustomData() map[string]string {
	return map[string]string{
		types.DeepLinkAddressKey: r.Address,
		types.DeepLinkAmountKey:  r.Amount.String(),
	}
}

func (r *GenerateSendTokensLinkRequest) Validate() error {
	if r.Address == "" {
		return fmt.Errorf("missing address")
	}

	if r.Chain == ChainType_UNDEFINED {
		return fmt.Errorf("invalid chain type")
	}

	err := r.Amount.Validate()
	if err != nil {
		return fmt.Errorf("invalid amount: %s", err)
	}

	return nil
}
