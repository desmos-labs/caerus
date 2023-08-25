package branch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/desmos-labs/caerus/types"
)

// Client represents a client that allows to interact with Branch.io APIs
type Client struct {
	apiKey string
}

// NewClient allows to build a new Client instance
func NewClient(cfg *Config) (*Client, error) {
	if cfg.ApiKey == "" {
		return nil, fmt.Errorf("invalid API key: %s", cfg.ApiKey)
	}

	return &Client{
		apiKey: cfg.ApiKey,
	}, nil
}

// NewClientFromEnvVariables allows to build a new Client instance from the environment variables
func NewClientFromEnvVariables() (*Client, error) {
	cfg, err := NewConfigFromEnvVariables()
	if err != nil {
		return nil, err
	}

	return NewClient(cfg)
}

// --------------------------------------------------------------------------------------------------------------------

type createLinkRequest struct {
	BranchKey string           `json:"branch_key"`
	Config    CreateLinkConfig `json:"data"`
}

type createLinkResponse struct {
	URL string `json:"url"`
}

// CreateDynamicLink allows to create a new dynamic link using the given link and request configurations
func (c *Client) CreateDynamicLink(apiKey string, config *types.LinkConfig) (string, error) {
	// Get the data to be used
	branchKey := c.apiKey
	if apiKey != "" {
		branchKey = apiKey
	}

	// Build the request body
	requestBody := createLinkRequest{
		BranchKey: branchKey,
		Config:    NewCreateLinkConfig(config),
	}

	// Marshal the request body to JSON
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error while marshalling request body: %s", err)
	}

	// Make the POST request to the Branch.io API
	response, err := http.Post(APIEndpoint, "application/json", bytes.NewBuffer(requestJSON))
	if err != nil {
		return "", fmt.Errorf("error while making POST request: %s", err)
	}
	defer response.Body.Close()

	// Decode the response
	resBz, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error while reading response body: %s", err)
	}

	var responseBody createLinkResponse
	err = json.Unmarshal(resBz, &responseBody)
	if err != nil {
		return "", fmt.Errorf("error while decoding response: %s", err)
	}

	// Return the URL
	return responseBody.URL, nil
}
