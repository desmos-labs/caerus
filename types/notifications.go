package types

import (
	"fmt"
	"time"

	"firebase.google.com/go/v4/messaging"
	"github.com/google/uuid"
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

const (
	NotificationTypeKey    = "type"
	NotificationMessageKey = "message"
)

type Notification struct {
	// Notification is the notification to be sent
	Notification *messaging.Notification `json:"notification,omitempty"`

	// Data is any additional data that should be sent along with the notification
	Data map[string]string `json:"data,omitempty"`

	// Whether the data inside Notification should be merged inside the Data field.
	// If set to true, the Notification will be set to nil when sending the notification to clients
	MergeNotificationWithData bool `json:"merge_notification_with_data,omitempty"`

	Android *messaging.AndroidConfig `json:"android_config,omitempty"`
	WebPush *messaging.WebpushConfig `json:"web_push_config,omitempty"`
	APNS    *messaging.APNSConfig    `json:"apns_config,omitempty"`
}

func (n *Notification) Validate() error {
	if n.Notification == nil && len(n.Data) == 0 {
		return fmt.Errorf("either notitication or data should be specified")
	}

	return nil
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

// --------------------------------------------------------------------------------------------------------------------

type SentNotification struct {
	ID            string
	AppID         string
	UserAddresses []string
	Notification  *Notification
	SendTime      time.Time
}

func NewSentNotification(appID string, userAddresses []string, notification *Notification) *SentNotification {
	return &SentNotification{
		ID:            uuid.NewString(),
		AppID:         appID,
		UserAddresses: userAddresses,
		Notification:  notification,
		SendTime:      time.Now(),
	}
}
