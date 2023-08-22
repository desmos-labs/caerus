package links

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

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

// HandleGenerateGenericDeepLinkRequest handles the given GenerateDeepLinkRequest
// Note: This method does *NOT* check for deep link creation rate limit, as the application will have to provide
// its own API key to create the deep link. This means that the application will be rate limited by Branch itself.
func (h *Handler) HandleGenerateGenericDeepLinkRequest(request *GenerateGenericDeepLinkRequest) (*CreateLinkResponse, error) {
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

// checkAppLinksRateLimit checks if the given app has exceeded the rate limit for deep links creation
func (h *Handler) checkAppLinksRateLimit(appID string) error {
	// Make sure the application has not exceeded the rate limit
	deepLinksRateLimit, err := h.db.GetAppDeepLinksRateLimit(appID)
	if err != nil {
		return err
	}

	deepLinksCount, err := h.db.GetAppDeepLinksCount(appID)
	if err != nil {
		return err
	}

	if deepLinksCount >= deepLinksRateLimit {
		return utils.NewTooManyRequestsError("deep links rate limit exceeded")
	}

	return nil
}

// HandleGenerateDeepLinkRequest handles the given GenerateDeepLinkRequest
func (h *Handler) HandleGenerateDeepLinkRequest(request GenerateDeepLinkRequest) (*CreateLinkResponse, error) {
	// Make sure the application has not exceeded the rate limit
	err := h.checkAppLinksRateLimit(request.GetAppID())
	if err != nil {
		return nil, err
	}

	// Build the custom data to use in order to generate the link
	customData := request.GetCustomData()
	customData[types.DeepLinkActionKey] = request.GetAction()
	customData[types.DeepLinkChainTypeKey] = strings.ToLower(request.GetChainType().String())

	customDataBz, err := json.Marshal(customData)
	if err != nil {
		return nil, err
	}

	// Build the path to use in order to generate the link
	values := url.Values{}
	for key, value := range customData {
		if key != types.DeepLinkActionKey {
			values.Add(key, value)
		}
	}

	deepLinkPath := fmt.Sprintf("/%s?%s", request.GetAction(), values.Encode())

	// Generate the link
	deepLink, err := h.deepLinksClient.CreateDynamicLink("", &types.LinkConfig{
		CustomData: customDataBz,
		DeepLinking: &types.DeepLinkConfig{
			DeepLinkPath: deepLinkPath,
		},
	})
	if err != nil {
		return nil, err
	}

	// Store the created deep link
	err = h.db.SaveCreatedDeepLink(types.NewCreatedDeepLink(request.GetAppID(), deepLink))
	if err != nil {
		return nil, err
	}

	return &CreateLinkResponse{Url: deepLink}, nil
}
