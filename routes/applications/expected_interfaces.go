package applications

import (
	"github.com/desmos-labs/caerus/types"
)

type Database interface {
	SaveAppNotificationDeviceToken(token *types.AppNotificationDeviceToken) error

	CanDeleteApp(desmosAddress string, appID string) (bool, error)
	DeleteApp(appID string) error
}
