package notifications

import (
	"github.com/desmos-labs/caerus/types"
)

// SendNotificationRequest represents the request sent when a user wants to send a new notification
// to users who have registered their device token
type SendNotificationRequest struct {
	// AppID represents the ID of the application that wants to send the notification
	AppID string

	// DeviceTokens represent the device tokens of the users that should receive the notification
	DeviceTokens []string `json:"device_tokens"`

	// Notification represents the notification to send
	Notification *types.Notification `json:"notification"`
}
