package branch_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/desmos-labs/caerus/branch"
	"github.com/desmos-labs/caerus/types"
	"github.com/desmos-labs/caerus/utils"
)

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

type ClientTestSuite struct {
	suite.Suite

	client *branch.Client
}

func (suite *ClientTestSuite) SetupSuite() {
	branchKey := utils.GetEnvOr(branch.EnvBranchKey, "")
	if branchKey == "" {
		suite.T().Skip("Skipping ClientTestSuite since BRANCH_KEY is not set")
	}

	client, err := branch.NewClient(&branch.Config{
		ApiKey: branchKey,
	})
	suite.Require().NoError(err)
	suite.client = client
}

func (suite *ClientTestSuite) TestCreateLink() {
	testCases := []struct {
		name            string
		buildLinkConfig func() *types.LinkConfig
		shouldErr       bool
		check           func(url string)
	}{
		{
			name: "link with full data is created properly",
			buildLinkConfig: func() *types.LinkConfig {
				customDataBz, err := json.Marshal(map[string]interface{}{
					"custom_key_string": "custom_value",
					"custom_key_int":    123,
					"custom_key_bool":   true,
				})
				suite.Require().NoError(err)

				return &types.LinkConfig{
					OpenGraph: &types.OpenGraphConfig{
						Title:       "Custom OG Title",
						Description: "Custom OG Description",
						ImageUrl:    "Custom OG Image URL",
					},
					Twitter: &types.TwitterConfig{
						CardType:    "Custom Twitter Card Type",
						Title:       "Custom Twitter Title",
						Description: "Custom Twitter Description",
						ImageUrl:    "Custom Twitter Image URL",
					},
					CustomData: customDataBz,
					Redirections: &types.RedirectionsConfig{
						FallbackUrl:    "Custom Fallback URL",
						DesktopUrl:     "Custom Desktop URL",
						IosUrl:         "Custom iOS URL",
						AndroidUrl:     "Custom Android URL",
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
			},
			shouldErr: false,
			check: func(url string) {
				suite.Require().NotEmpty(url)
			},
		},
		{
			name: "link with partial data is created properly",
			buildLinkConfig: func() *types.LinkConfig {
				customDataBz, err := json.Marshal(map[string]interface{}{
					types.DeepLinkActionKey:    types.DeepLinkActionViewProfile,
					types.DeepLinkAddressKey:   "desmos1aph74nw42mk7330pftwwmj7lxr7j3drgmlu3zc",
					types.DeepLinkChainTypeKey: "mainnet",
				})
				suite.Require().NoError(err)

				return &types.LinkConfig{
					CustomData: customDataBz,
					DeepLinking: &types.DeepLinkConfig{
						DeepLinkPath: "/view_profile?address=desmos1aph74nw42mk7330pftwwmj7lxr7j3drgmlu3zc&chain_type=mainnet",
					},
				}
			},
			shouldErr: false,
			check: func(url string) {
				fmt.Println(url)
				suite.Require().NotEmpty(url)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			res, err := suite.client.CreateDynamicLink("", tc.buildLinkConfig())

			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}

			if tc.check != nil {
				tc.check(res)
			}
		})
	}
}
