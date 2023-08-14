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
	AppID string
}

func NewDeleteAppRequest(appID string) *DeleteAppRequest {
	return &DeleteAppRequest{
		AppId: appID,
	}
}
