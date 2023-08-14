package applications

import (
	"github.com/desmos-labs/caerus/types"
)

type Database interface {
	DeleteApp(appID string) error
	SaveAppNotificationDeviceToken(token *types.AppNotificationDeviceToken) error
}
