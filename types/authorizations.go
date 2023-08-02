package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// FeeGrantRequest contains the data of a fee grant allowance request that the user has made
type FeeGrantRequest struct {
	// DesmosAddress is the Desmos address of the user requesting the fee grant allowance.
	DesmosAddress string

	// RequestTime is the time at which the user requested the fee grant allowance.
	RequestTime time.Time

	// GrantTime is the time at which the user wants the fee grant allowance to be granted.
	// If nil, the fee grant allowance has not been granted yet.
	GrantTime *time.Time
}

func NewFeeGrantRequest(desmosAddress string, requestTime time.Time, grantTime *time.Time) FeeGrantRequest {
	return FeeGrantRequest{
		DesmosAddress: desmosAddress,
		RequestTime:   requestTime,
		GrantTime:     grantTime,
	}
}

// FeeGrantAllowance contains the data of a fee grant allowance that the user has been granted
type FeeGrantAllowance struct {
	// ExpirationTime is the time at which the fee grant allowance expires.
	// If nil, the fee grant allowance has no expiration time.
	ExpirationTime *time.Time

	// SpendLimit is the maximum amount of coins that can be spent by the user.
	// If nil, the fee grant allowance has no spend limit.
	SpendLimit sdk.Coins

	// AllowedMessages are the messages that the user is allowed to send.
	// If nil, the user will be able to send any message.
	AllowedMessages []string
}

func NewAuthorization(expirationTime *time.Time, spendLimit sdk.Coins, allowedMessages []string) *FeeGrantAllowance {
	return &FeeGrantAllowance{
		ExpirationTime:  expirationTime,
		AllowedMessages: allowedMessages,
		SpendLimit:      spendLimit,
	}
}
