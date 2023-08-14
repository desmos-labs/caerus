package applications

import (
	"github.com/desmos-labs/caerus/routes/base"
	"github.com/desmos-labs/caerus/types"
)

type Database interface {
	base.Database

	DeleteApp(appID string) error
	SaveAppNotificationDeviceToken(token *types.AppNotificationDeviceToken) error
}
