package grants

type RequestFeeGrantRequest struct {
	// AppID represents the ID of the application that is requesting the fee grant
	AppID string

	// DesmosAddress is the address of the user requesting the fee grant
	DesmosAddress string
}

func NewRequestFeeGrantRequest(appID string, desmosAddress string) *RequestFeeGrantRequest {
	return &RequestFeeGrantRequest{
		AppID:         appID,
		DesmosAddress: desmosAddress,
	}
}
