package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/desmos-labs/caerus/types"
)

func TestLinkConfig_MarshalJSON(t *testing.T) {
	config := types.LinkConfig{
		OpenGraph: &types.OpenGraphConfig{
			Title:       "Custom OG Title",
			Description: "Custom OG Description",
			ImageURL:    "Custom OG Image URL",
		},
		Twitter: &types.TwitterConfig{
			CardType:    "Custom Twitter Card Type",
			Title:       "Custom Twitter Title",
			Description: "Custom Twitter Description",
			ImageURL:    "Custom Twitter Image URL",
		},
		CustomData: map[string]interface{}{
			"custom_key_string": "custom_value",
			"custom_key_int":    123,
			"custom_key_bool":   true,
		},
		Redirections: &types.RedirectionsConfig{
			FallbackURL:    "Custom Fallback URL",
			DesktopURL:     "Custom Desktop URL",
			IosURL:         "Custom iOS URL",
			AndroidURL:     "Custom Android URL",
			WebOnly:        true,
			DesktopWebOnly: true,
			MobileWebOnly:  true,
		},
		DeepLinking: &types.DeepLinkConfig{
			DeepLinkPath:        "?address=desmos1qz26yufk0m2s3z2v&action=send&amount=1000udaric",
			AndroidDeepLinkPath: "Custom Android Deep Link Path",
			IosDeepLinkPath:     "Custom iOS Deep Link Path",
			DesktopDeepLinkPath: "Custom Desktop Deep Link Path",
		},
	}

	configBz, err := json.Marshal(config)
	require.NoError(t, err)

	// -------------------------------
	// --- Serialization check
	// -------------------------------

	var expectedData map[string]interface{}
	err = json.Unmarshal([]byte(`
{
  "$og_title": "Custom OG Title",
  "$og_description": "Custom OG Description",
  "$og_image_url": "Custom OG Image URL",
  "$twitter_card": "Custom Twitter Card Type",
  "$twitter_title": "Custom Twitter Title",
  "$twitter_description": "Custom Twitter Description",
  "$twitter_image_url": "Custom Twitter Image URL",
  "custom_key_bool": true,
  "custom_key_int": 123,
  "custom_key_string": "custom_value",
  "$fallback_url": "Custom Fallback URL",
  "$desktop_url": "Custom Desktop URL",
  "$ios_url": "Custom iOS URL",
  "$android_url": "Custom Android URL",
  "$web_only": true,
  "$desktop_web_only": true,
  "$mobile_web_only": true,
  "$deeplink_path": "?address=desmos1qz26yufk0m2s3z2v&action=send&amount=1000udaric",
  "$android_deeplink_path": "Custom Android Deep Link Path",
  "$ios_deeplink_path": "Custom iOS Deep Link Path",
  "$desktop_deeplink_path": "Custom Desktop Deep Link Path"
}`), &expectedData)
	require.NoError(t, err)

	var serializedData map[string]interface{}
	err = json.Unmarshal(configBz, &serializedData)
	require.NoError(t, err)

	require.Equal(t, expectedData, serializedData)
}
