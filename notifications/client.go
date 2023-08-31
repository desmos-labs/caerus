package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"

	caerus "github.com/desmos-labs/caerus/types"
)

type Client struct {
	firebaseMessaging *messaging.Client
	db                Database
}

// NewClient returns a new Client instance
func NewClient(config *Config, db Database) (*Client, error) {
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
		db:                db,
	}, nil
}

// NewClientFromEnvVariables returns a new Client instance by using the environment variables
func NewClientFromEnvVariables(db Database) (*Client, error) {
	config, err := ReadConfigFromEnvVariables()
	if err != nil {
		return nil, err
	}

	return NewClient(config, db)
}

// --------------------------------------------------------------------------------------------------------------------

// SendNotificationToUsers sends the given notification to the devices of the users having the given addresses
func (c *Client) SendNotificationToUsers(app *caerus.Application, usersAddresses []string, notification *caerus.Notification) error {
	// Do nothing if the client is not configured
	if c.firebaseMessaging == nil {
		return nil
	}

	// Get the notification tokens
	var deviceTokens []string
	for _, address := range usersAddresses {
		tokens, err := c.db.GetUserNotificationTokens(address)
		if err != nil {
			return err
		}

		deviceTokens = append(deviceTokens, tokens...)
	}

	// Send the notification only if there is at least one device token registered
	if len(deviceTokens) > 0 {
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
		if err != nil {
			return err
		}
	}

	// Store the sent notification
	return c.db.SaveSentNotification(caerus.NewSentNotification(app.ID, usersAddresses, notification))
}

// SendNotificationToApp allows to send the given notification to the application having the given id
func (c *Client) SendNotificationToApp(appID string, notification *caerus.Notification) error {
	// Get the app with the given id
	app, found, err := c.db.GetApp(appID)
	if err != nil {
		return err
	}

	// Do nothing if the app was not found
	if !found {
		return nil
	}

	// Do nothing if the app does not provide a notifications webhook URL
	if app.NotificationsWebhookURL == "" {
		return nil
	}

	// Crete a POST request to the notification webhook URL authenticating the request using the Authorization header
	bodyBz, err := json.Marshal(&notification)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, app.NotificationsWebhookURL, bytes.NewBuffer(bodyBz))
	if err != nil {
		return err
	}

	// Set the authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer: %s", app.SecretKey))
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client and send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Handle errors
	if resp.StatusCode != http.StatusOK {
		respBodyBz, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("error while sending the request: %s", string(respBodyBz))
	}

	return nil
}
