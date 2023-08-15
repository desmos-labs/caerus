package notifications

import (
	"github.com/desmos-labs/caerus/types"
)

type SendAppNotificationRequest struct {
	AppID        string
	DeviceTokens []string
	Notification *types.Notification
}

func NewSendAppNotificationRequest(appID string, deviceTokens []string, notification *types.Notification) *SendAppNotificationRequest {
	return &SendAppNotificationRequest{
		AppID:        appID,
		DeviceTokens: deviceTokens,
		Notification: notification,
	}
}
