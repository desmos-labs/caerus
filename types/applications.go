package types

import (
	"time"

	"github.com/google/uuid"
)

type ApplicationSubscription struct {
	ID                     uint64
	Name                   string
	FeeGrantLimit          uint64
	NotificationsRateLimit uint64
	DeepLinksRateLimit     uint64
}

func NewApplicationSubscription(id uint64, name string, feeGrantLimit uint64, notificationsRateLimit uint64, deepLinksRateLimit uint64) ApplicationSubscription {
	return ApplicationSubscription{
		ID:                     id,
		Name:                   name,
		FeeGrantLimit:          feeGrantLimit,
		NotificationsRateLimit: notificationsRateLimit,
		DeepLinksRateLimit:     deepLinksRateLimit,
	}
}

// --------------------------------------------------------------------------------------------------------------------

type Application struct {
	ID                      string
	Name                    string
	WalletAddress           string
	SubscriptionID          uint64
	SecretKey               string
	NotificationsWebhookURL string
	Admins                  []string
	CreationTime            time.Time
}

func NewApplication(
	id string,
	name string,
	walletAddress string,
	subscriptionID uint64,
	secretKey string,
	notificationsWebhookURL string,
	admins []string,
	creationTime time.Time,
) *Application {
	return &Application{
		ID:                      id,
		Name:                    name,
		WalletAddress:           walletAddress,
		SubscriptionID:          subscriptionID,
		SecretKey:               secretKey,
		NotificationsWebhookURL: notificationsWebhookURL,
		Admins:                  admins,
		CreationTime:            creationTime,
	}
}

// --------------------------------------------------------------------------------------------------------------------

type AppToken struct {
	AppID string
	Name  string
	Value string
}

func NewAppToken(appID string, tokenName string) AppToken {
	return AppToken{
		AppID: appID,
		Name:  tokenName,
		Value: uuid.NewString(),
	}
}

type EncryptedAppToken struct {
	AppID        string
	TokenName    string
	TokenValue   string
	CreationTime time.Time
}

// NewEncryptedAppToken returns a new EncryptedAppToken instance
func NewEncryptedAppToken(appID string, tokenName string, tokenValue string, creationTime time.Time) *EncryptedAppToken {
	return &EncryptedAppToken{
		AppID:        appID,
		TokenName:    tokenName,
		TokenValue:   tokenValue,
		CreationTime: creationTime,
	}
}
