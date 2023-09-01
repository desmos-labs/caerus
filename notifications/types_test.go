package notifications_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/desmos-labs/caerus/notifications"
	"github.com/desmos-labs/caerus/types"
)

func TestNotificationSignature(t *testing.T) {
	secretKey := "secretKey"
	notification := types.Notification{
		Data: map[string]string{
			types.NotificationTypeKey:    "fee_allowances_granted",
			types.NotificationMessageKey: "Your fee allowances have been granted",
			"users":                      "user1,user2",
		},
	}

	signature, err := notifications.SignNotification(notification, secretKey)
	require.NoError(t, err)

	isValid, err := notifications.VerifyNotificationSignature(notification, signature, secretKey)
	require.NoError(t, err)
	require.True(t, isValid)
}
