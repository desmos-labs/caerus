package applications

type RegisterAppDeviceTokenRequest struct {
	AppID       string
	DeviceToken string
}

func NewRegisterAppDeviceTokenRequest(appID string, deviceToken string) *RegisterAppDeviceTokenRequest {
	return &RegisterAppDeviceTokenRequest{
		AppID:       appID,
		DeviceToken: deviceToken,
	}
}

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
