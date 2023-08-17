package grants

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
)

type RequestFeeGrantRequest struct {
	// AppID represents the ID of the application that is requesting the fee grant
	AppID string

	// DesmosAddress is the address of the user requesting the fee grant
	DesmosAddress string

	// Allowance represents the allowance to be granted to the user
	Allowance *codectypes.Any
}

func NewRequestFeeGrantRequest(appID string, desmosAddress string, allowance *codectypes.Any) *RequestFeeGrantRequest {
	return &RequestFeeGrantRequest{
		AppID:         appID,
		DesmosAddress: desmosAddress,
		Allowance:     allowance,
	}
}
