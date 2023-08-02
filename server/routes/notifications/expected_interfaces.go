package notifications

import (
	"github.com/desmos-labs/caerus/server/routes/base"
	"github.com/desmos-labs/caerus/types"
)

type Database interface {
	base.Database
	CanSendNotifications(token string) (bool, error)
	SaveNotificationDeviceToken(token *types.NotificationDeviceToken) error
}

type Firebase interface {
	SendNotifications(deviceTokens []string, notification *types.Notification) error
}
