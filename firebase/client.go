package firebase

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/firebasedynamiclinks/v1"
	"google.golang.org/api/option"

	caerus "github.com/desmos-labs/caerus/types"
)

type Client struct {
	cfg               *Config
	firebaseLinks     *firebasedynamiclinks.Service
	firebaseMessaging *messaging.Client
}

// NewClient returns a new Client instance
func NewClient(config *Config) (*Client, error) {
	if config.CredentialsFilePath == "" {
		log.Info().Msg("No Firebase credentials file found")
		return nil, nil
	}

	var firebaseLinksService *firebasedynamiclinks.Service
	var firebaseMessagingService *messaging.Client

	if config.CredentialsFilePath != "" {
		ctx := context.Background()
		options := option.WithCredentialsFile(config.CredentialsFilePath)

		// Create the Firebase app
		firebaseApp, err := firebase.NewApp(ctx, nil, options)
		if err != nil {
			return nil, err
		}

		// Create Dynamic Links client
		if config.Links != nil {
			dynamicLinksService, err := firebasedynamiclinks.NewService(ctx, options)
			if err != nil {
				return nil, err
			}
			firebaseLinksService = dynamicLinksService
		}

		// Create Firebase Messaging client
		if config.Notifications != nil {
			messagingService, err := firebaseApp.Messaging(ctx)
			if err != nil {
				return nil, err
			}
			firebaseMessagingService = messagingService
		}
	}

	return &Client{
		firebaseLinks:     firebaseLinksService,
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

// --------------------------------------------------------------------------------------------------------------------

// GetLinkDesktopInfo returns the Desktop info to be used when generating a new link
func (c *Client) GetLinkDesktopInfo() *firebasedynamiclinks.DesktopInfo {
	return &firebasedynamiclinks.DesktopInfo{
		DesktopFallbackLink: c.cfg.Links.Desktop.FallbackLink,
	}
}

// GetLinkAndroidInfo returns the Android info to be used when generating a new link
func (c *Client) GetLinkAndroidInfo() *firebasedynamiclinks.AndroidInfo {
	return &firebasedynamiclinks.AndroidInfo{
		AndroidMinPackageVersionCode: c.cfg.Links.Android.MinPackageVersionCode,
		AndroidPackageName:           c.cfg.Links.Android.PackageName,
	}
}

// GetLinkIOSInfo returns the IOS info to be used when generating a new link
func (c *Client) GetLinkIOSInfo() *firebasedynamiclinks.IosInfo {
	return &firebasedynamiclinks.IosInfo{
		IosAppStoreId:     c.cfg.Links.Ios.AppStoreID,
		IosBundleId:       c.cfg.Links.Ios.BundleID,
		IosMinimumVersion: c.cfg.Links.Ios.MinimumVersion,
	}
}

// GenerateLink generates a new link using the given info
func (c *Client) GenerateLink(info caerus.LinkInfo) (string, error) {
	if c.firebaseLinks == nil {
		return "", nil
	}

	// Generate the link
	res, err := c.firebaseLinks.ShortLinks.Create(&firebasedynamiclinks.CreateShortDynamicLinkRequest{
		DynamicLinkInfo: &firebasedynamiclinks.DynamicLinkInfo{
			DynamicLinkDomain: c.cfg.Links.Domain,

			DesktopInfo: c.GetLinkDesktopInfo(),
			AndroidInfo: c.GetLinkAndroidInfo(),
			IosInfo:     c.GetLinkIOSInfo(),

			Link:              info.Link,
			SocialMetaTagInfo: info.SocialMetaTagInfo,
		},
		Suffix: info.Suffix,
	}).Do()
	if err != nil {
		return "", err
	}

	// Return the short link
	return res.ShortLink, nil
}
