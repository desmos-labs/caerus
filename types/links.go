package types

import (
	"time"

	"github.com/google/uuid"
)

const (
	DeepLinkChainTypeKey = "chain_type"
	DeepLinkActionKey    = "action"
	DeepLinkAddressKey   = "address"
	DeepLinkAmountKey    = "amount"

	DeepLinkActionSendTokens  = "send_tokens"
	DeepLinkActionViewProfile = "view_profile"
)

// CreatedDeepLink contains the data of a deep link created by an application
type CreatedDeepLink struct {
	ID           string
	AppID        string
	URL          string
	Config       *LinkConfig
	CreationTime time.Time
}

func NewCreatedDeepLink(appID string, linkURL string, config *LinkConfig) *CreatedDeepLink {
	return &CreatedDeepLink{
		ID:           uuid.NewString(),
		AppID:        appID,
		URL:          linkURL,
		Config:       config,
		CreationTime: time.Now(),
	}
}

// --------------------------------------------------------------------------------------------------------------------

// LinkDetails contains the details of a link
type LinkDetails struct {
	Data map[string]interface{}
}
