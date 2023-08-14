package applications

type RegisterAppDeviceTokenRequest struct {
	AppID       string
	DeviceToken string `json:"token"`
}

type DeleteApplicationRequest struct {
	AppID string
}
