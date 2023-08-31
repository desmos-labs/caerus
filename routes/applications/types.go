package applications

type DeleteApplicationRequest struct {
	UserAddress string
	AppID       string
}

func NewDeleteApplicationRequest(desmosAddress string, appID string) *DeleteApplicationRequest {
	return &DeleteApplicationRequest{
		UserAddress: desmosAddress,
		AppID:       appID,
	}
}
