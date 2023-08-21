package types

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// linkConfigJson represents the type that will be used to properly marshal the LinkConfig struct
type linkConfigJson struct {
	*OpenGraphConfig
	*TwitterConfig
	*RedirectionsConfig
	*DeepLinkConfig
}

func (c *LinkConfig) MarshalJSON() ([]byte, error) {
	data := linkConfigJson{
		OpenGraphConfig:    c.OpenGraph,
		TwitterConfig:      c.Twitter,
		RedirectionsConfig: c.Redirections,
		DeepLinkConfig:     c.DeepLinking,
	}

	// Serialize the link data to make sure everything is at the same level
	dataWithoutCustomBz, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Read the link configuration data as a map to allow merging with the custom data
	var linkData map[string]interface{}
	err = json.Unmarshal(dataWithoutCustomBz, &linkData)
	if err != nil {
		return nil, err
	}

	// Read the custom data as a map
	var customData map[string]interface{}
	err = json.Unmarshal(c.CustomData, &customData)
	if err != nil {
		return nil, err
	}

	// Merge the custom data and link data together
	for key, value := range customData {
		linkData[key] = value
	}

	return json.Marshal(linkData)
}
