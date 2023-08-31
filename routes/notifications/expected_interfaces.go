package notifications

import (
	"github.com/desmos-labs/caerus/types"
)

type Database interface {
	GetApp(appID string) (*types.Application, bool, error)

	GetAppNotificationsRateLimit(appID string) (uint64, error)
	GetAppNotificationsCount(appID string) (uint64, error)
}

type Firebase interface {
	SendNotificationToUsers(application *types.Application, usersAddresses []string, notification *types.Notification) error
}
