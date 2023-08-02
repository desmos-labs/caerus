package notifications

import (
	"github.com/desmos-labs/caerus/types"
)

// RegisterDeviceTokenRequest represents the request sent when a user wants to register a new device token
// to receive notifications
type RegisterDeviceTokenRequest struct {
	UserAddress string
	DeviceToken string `json:"device_token"`
}

// SendNotificationRequest represents the request sent when a user wants to send a new notification
// to users who have registered their device token
type SendNotificationRequest struct {
	Token        string
	DeviceTokens []string            `json:"device_tokens"`
	Notification *types.Notification `json:"notification"`
}
