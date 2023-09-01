package notifications

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"

	"github.com/desmos-labs/caerus/types"
)

type SentNotificationBody struct {
	// Notification contains the notification to send
	Notification types.Notification `json:"notification"`

	// Signature contains the signature of the notification.
	// This is achieved by signing the Notification's JSON representation
	// using the application's secret key
	Signature string `json:"signature"`
}

// SignNotification signs the given notification using the given secret key
func SignNotification(notification types.Notification, secretKey string) (string, error) {
	// Serialize the notification to a JSON format
	notificationBz, err := json.Marshal(notification)
	if err != nil {
		return "", err
	}

	// Compute the signature
	mac := hmac.New(sha256.New, []byte(secretKey))
	_, err = mac.Write(notificationBz)
	if err != nil {
		return "", err
	}
	signatureBz := mac.Sum(nil)

	// Return the signature in a Base64 format
	return base64.StdEncoding.EncodeToString(signatureBz), nil
}

// VerifyNotificationSignature verifies that the given signature is valid for the given notification
func VerifyNotificationSignature(notification types.Notification, signature string, secretKey string) (bool, error) {
	// Sign the notification
	expectedSignature, err := SignNotification(notification, secretKey)
	if err != nil {
		return false, err
	}

	return expectedSignature == signature, nil
}
