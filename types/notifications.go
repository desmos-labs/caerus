package types

import (
	"firebase.google.com/go/v4/messaging"
)

// AppNotificationDeviceToken represents a notification device token that is associated to an application
type AppNotificationDeviceToken struct {
	AppID       string
	DeviceToken string
}

func NewAppNotificationDeviceToken(appID, deviceToken string) *AppNotificationDeviceToken {
	return &AppNotificationDeviceToken{
		AppID:       appID,
		DeviceToken: deviceToken,
	}
}

// --------------------------------------------------------------------------------------------------------------------

// UserNotificationDeviceToken represents a notification device token that is associated to a user
type UserNotificationDeviceToken struct {
	UserAddress string
	DeviceToken string
}

func NewUserNotificationDeviceToken(userAddress, deviceToken string) *UserNotificationDeviceToken {
	return &UserNotificationDeviceToken{
		UserAddress: userAddress,
		DeviceToken: deviceToken,
	}
}

// --------------------------------------------------------------------------------------------------------------------

type Notification struct {
	// Notification is the notification to be sent
	Notification *messaging.Notification `json:"notification"`

	// Data is any additional data that should be sent along with the notification
	Data map[string]string `json:"data"`

	Android *messaging.AndroidConfig `json:"android_config"`
	WebPush *messaging.WebpushConfig `json:"web_push_config"`
	APNS    *messaging.APNSConfig    `json:"apns_config"`
}

// BackgroundAndroidConfig returns the AndroidConfig to be used when sending a notification in background.
// This makes sure that the notification can be properly received by the device.
func BackgroundAndroidConfig() *messaging.AndroidConfig {
	return &messaging.AndroidConfig{
		Priority: "high",
	}
}

// BackgroundAPNSConfig returns the APNSConfig to be used when sending a notification in background to iOS devices.
// This ensures that the notification can be properly received by the device.
func BackgroundAPNSConfig() *messaging.APNSConfig {
	return &messaging.APNSConfig{
		Payload: &messaging.APNSPayload{
			Aps: &messaging.Aps{
				Sound:            "default",
				ContentAvailable: true,
				MutableContent:   true,
			},
		},
		Headers: map[string]string{
			"apns-priority": "5",
		},
	}
}
