package links

import (
	"github.com/desmos-labs/caerus/types"
)

type Database interface {
	GetAppDeepLinksRateLimit(appID string) (uint64, error)
	GetAppDeepLinksCount(appID string) (uint64, error)

	SaveCreatedDeepLink(link *types.CreatedDeepLink) error
}

type DeepLinksClient interface {
	CreateDynamicLink(apiKey string, linkConfig *types.LinkConfig) (string, error)
}
