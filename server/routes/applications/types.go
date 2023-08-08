package applications

type RegisterAppDeviceTokenRequest struct {
	AppID       string
	DeviceToken string `json:"device_token"`
}

type DeleteApplicationRequest struct {
	AppID string
}
