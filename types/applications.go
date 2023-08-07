package types

import (
	"time"
)

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

type Application struct {
	ID             string
	Name           string
	WalletAddress  string
	SubscriptionID uint64
	CreationTime   time.Time
}

func NewApplication(id string, name string, walletAddress string, subscriptionID uint64, creationTime time.Time) *Application {
	return &Application{
		ID:             id,
		Name:           name,
		WalletAddress:  walletAddress,
		SubscriptionID: subscriptionID,
		CreationTime:   creationTime,
	}
}
