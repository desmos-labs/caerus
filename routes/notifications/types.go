package notifications

import (
	"fmt"

	"github.com/desmos-labs/caerus/types"
)

type SendAppNotificationRequest struct {
	AppID         string
	UserAddresses []string
	Notification  *types.Notification
}

func NewSendAppNotificationRequest(appID string, userAddresses []string, notification *types.Notification) *SendAppNotificationRequest {
	return &SendAppNotificationRequest{
		AppID:         appID,
		UserAddresses: userAddresses,
		Notification:  notification,
	}
}

func (r SendAppNotificationRequest) Validate() error {
	if len(r.UserAddresses) == 0 {
		return fmt.Errorf("user addresses cannot be empty")
	}

	return r.Notification.Validate()
}
