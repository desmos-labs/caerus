package links

import (
	"github.com/desmos-labs/caerus/types"
)

type GenerateDeepLinkRequest struct {
	// ID of the app that is trying to generate the link
	AppID string

	// Configuration used to generate the link
	LinkConfig *types.LinkConfig

	// (optional) API key used to generate the link
	ApiKey string
}

func NewGenerateDeepLinkRequest(appID string, config *types.LinkConfig, apiKey string) *GenerateDeepLinkRequest {
	return &GenerateDeepLinkRequest{
		AppID:      appID,
		LinkConfig: config,
		ApiKey:     apiKey,
	}
}
