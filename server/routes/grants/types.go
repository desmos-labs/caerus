package grants

type RequestFeeGrantRequest struct {
	// Address of the user requesting the fee grant
	DesmosAddress string

	// Address of the user that should grant the fee grant
	GranterAddress string `json:"granter_address"`
}
