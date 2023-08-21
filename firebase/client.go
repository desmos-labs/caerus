package firebase

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"

	caerus "github.com/desmos-labs/caerus/types"
)

type Client struct {
	cfg               *Config
	firebaseMessaging *messaging.Client
}

// NewClient returns a new Client instance
func NewClient(config *Config) (*Client, error) {
	if config.CredentialsFilePath == "" {
		log.Info().Msg("No Firebase credentials file found")
		return nil, nil
	}

	var firebaseMessagingService *messaging.Client

	if config.CredentialsFilePath != "" {
		ctx := context.Background()
		options := option.WithCredentialsFile(config.CredentialsFilePath)

		// Create the Firebase app
		firebaseApp, err := firebase.NewApp(ctx, nil, options)
		if err != nil {
			return nil, err
		}

		// Create Firebase Messaging client
		messagingService, err := firebaseApp.Messaging(ctx)
		if err != nil {
			return nil, err
		}
		firebaseMessagingService = messagingService
	}

	return &Client{
		firebaseMessaging: firebaseMessagingService,
	}, nil
}

// NewClientFromEnvVariables returns a new Client instance by using the environment variables
func NewClientFromEnvVariables() (*Client, error) {
	config, err := ReadConfigFromEnvVariables()
	if err != nil {
		return nil, err
	}

	return NewClient(config)
}

// --------------------------------------------------------------------------------------------------------------------

// SendNotifications sends the given notification to the devices that have registered using the given tokens
func (c *Client) SendNotifications(app *caerus.Application, deviceTokens []string, notification *caerus.Notification) error {
	// Do nothing if the client is not configured
	if c.firebaseMessaging == nil {
		return nil
	}

	// Do nothing if there are no tokens
	if len(deviceTokens) == 0 {
		return nil
	}

	// Put the application name and id inside the data
	if app != nil {
		notification.Data[ApplicationIDKey] = app.ID
		notification.Data[ApplicationNameKey] = app.Name
	}

	// Put the title and body of the notification inside the data if the config specify this
	var notificationData = notification.Notification
	if notification.MergeNotificationWithData {
		notification.Data[NotificationTitleKey] = notificationData.Title
		notification.Data[NotificationBodyKey] = notificationData.Body
		notification.Data[NotificationImageKey] = notificationData.ImageURL
		notificationData = nil
	}

	// Send the notification
	_, err := c.firebaseMessaging.SendEachForMulticast(context.Background(), &messaging.MulticastMessage{
		Tokens:       deviceTokens,
		Notification: notificationData,
		Data:         notification.Data,
		Android:      notification.Android,
		Webpush:      notification.WebPush,
		APNS:         notification.APNS,
	})
	return err
}
