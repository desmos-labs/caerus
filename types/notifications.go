package types

import (
	"firebase.google.com/go/v4/messaging"
)

type NotificationDeviceToken struct {
	UserAddress string
	DeviceToken string
}

func NewNotificationDeviceToken(userAddress, deviceToken string) *NotificationDeviceToken {
	return &NotificationDeviceToken{
		UserAddress: userAddress,
		DeviceToken: deviceToken,
	}
}

type Notification struct {
	Data         map[string]string        `json:"data"`
	Notification *messaging.Notification  `json:"notification"`
	Android      *messaging.AndroidConfig `json:"android_config"`
	WebPush      *messaging.WebpushConfig `json:"web_push_config"`
	APNS         *messaging.APNSConfig    `json:"apns_config"`
}

func NewNotification(data map[string]string, notification *messaging.Notification, android *messaging.AndroidConfig, webPush *messaging.WebpushConfig, apns *messaging.APNSConfig) *Notification {
	return &Notification{
		Data:         data,
		Notification: notification,
		Android:      android,
		WebPush:      webPush,
		APNS:         apns,
	}
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
