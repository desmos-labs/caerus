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

// --------------------------------------------------------------------------------------------------------------------

type NotificationApplication struct {
	ID   string
	Name string
}

type NotificationSender struct {
	Token       string
	Application *NotificationApplication
}

// --------------------------------------------------------------------------------------------------------------------

type Notification struct {
	// Application contains the details of the application that is sending the notification.
	// It can be nil if the notification is sent by the server itself, or if there is no
	// authentication set for sending notifications.
	Application *NotificationApplication

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
