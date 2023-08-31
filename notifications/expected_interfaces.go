package notifications

import (
	"github.com/desmos-labs/caerus/types"
)

type Database interface {
	GetUserNotificationTokens(address string) ([]string, error)

	GetApp(appID string) (*types.Application, bool, error)

	SaveSentNotification(notification *types.SentNotification) error
}
