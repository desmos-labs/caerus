package links

import (
	"github.com/desmos-labs/caerus/types"
	"github.com/desmos-labs/caerus/utils"
)

type Handler struct {
	deepLinksClient DeepLinksClient
	db              Database
}

// NewHandler allows to build a new Handler instance
func NewHandler(client DeepLinksClient, db Database) *Handler {
	return &Handler{
		deepLinksClient: client,
		db:              db,
	}
}

// NewHandlerFromEnvVariables returns a new Handler instance reading the config from the env variables
func NewHandlerFromEnvVariables(client DeepLinksClient, db Database) *Handler {
	return NewHandler(client, db)
}

// HandleGenerateDeepLinkRequest handles the given GenerateDeepLinkRequest
func (h *Handler) HandleGenerateDeepLinkRequest(request *GenerateDeepLinkRequest) (*CreateLinkResponse, error) {
	// Make sure the application has not exceeded the rate limit
	deepLinksRateLimit, err := h.db.GetAppDeepLinksRateLimit(request.AppID)
	if err != nil {
		return nil, err
	}

	deepLinksCount, err := h.db.GetAppDeepLinksCount(request.AppID)
	if err != nil {
		return nil, err
	}

	if deepLinksCount > deepLinksRateLimit {
		return nil, utils.NewTooManyRequestsError("deep links rate limit exceeded")
	}

	// Generate the deep link
	deepLink, err := h.deepLinksClient.CreateDynamicLink(request.ApiKey, request.LinkConfig)
	if err != nil {
		return nil, err
	}

	// Store the created deep link
	err = h.db.SaveCreatedDeepLink(types.NewCreatedDeepLink(request.AppID, deepLink))
	if err != nil {
		return nil, err
	}

	// Return the response data
	return &CreateLinkResponse{Url: deepLink}, nil
}
