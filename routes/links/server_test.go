package links_test

import (
	"context"
	"path"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	desmosapp "github.com/desmos-labs/desmos/v5/app"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"

	"github.com/desmos-labs/caerus/authentication"
	"github.com/desmos-labs/caerus/database"
	"github.com/desmos-labs/caerus/routes/links"
	linkstestutils "github.com/desmos-labs/caerus/routes/links/testutils"
	"github.com/desmos-labs/caerus/testutils"
	"github.com/desmos-labs/caerus/types"
)

func TestLinksServerTestSuite(t *testing.T) {
	suite.Run(t, new(LinksServerTestSuite))
}

type LinksServerTestSuite struct {
	suite.Suite

	db              *database.Database
	deepLinksClient *linkstestutils.MockDeepLinksClient
	handler         *links.Handler

	server *grpc.Server
	client links.LinksServiceClient
}

func (suite *LinksServerTestSuite) SetupSuite() {
	// Setup the Cosmos SDK
	desmosapp.SetupConfig(sdk.GetConfig())

	// Setup the database
	db, err := testutils.CreateDatabase(path.Join("..", "..", "database", "schema"))
	suite.Require().NoError(err)
	suite.db = db
}

func (suite *LinksServerTestSuite) SetupTest() {
	// Truncate all the database data to make sure we have a clean database state
	err := testutils.TruncateDatabase(suite.db)
	suite.Require().NoError(err)

	// Setup the mocks
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	// Build the mocks
	suite.deepLinksClient = linkstestutils.NewMockDeepLinksClient(ctrl)

	// Create the handler
	suite.handler = links.NewHandler(suite.deepLinksClient, suite.db)

	// Create the server
	suite.server = testutils.CreateServer(suite.db)

	// Register the service
	service := links.NewServer(suite.handler)
	links.RegisterLinksServiceServer(suite.server, service)

	// Setup the client
	conn, err := testutils.StartServerAndConnect(suite.server)
	suite.Require().NoError(err)
	suite.client = links.NewLinksServiceClient(conn)
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *LinksServerTestSuite) TestCreateAddressLink() {
	testCases := []struct {
		name         string
		setup        func()
		setupContext func(ctx context.Context) context.Context
		buildRequest func() *links.CreateAddressLinkRequest
		shouldErr    bool
		check        func(res *links.CreateLinkResponse)
	}{
		{
			name: "invalid session returns error",
			buildRequest: func() *links.CreateAddressLinkRequest {
				return &links.CreateAddressLinkRequest{
					Address: "desmos1aph74nw42mk7330pftwwmj7lxr7j3drgmlu3zc",
					Chain:   links.ChainType_MAINNET,
				}
			},
			shouldErr: true,
		},
		{
			name: "invalid address returns error",
			setup: func() {
				err := suite.db.SaveAppSubscription(types.NewApplicationSubscription(
					1,
					"Test App Subscription",
					10,
					10,
					10,
				))
				suite.Require().NoError(err)

				err = suite.db.SaveApp(types.Application{
					ID:             "1",
					Name:           "Test Application",
					WalletAddress:  "desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					SubscriptionID: 1,
					Admins: []string{
						"desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					},
					CreationTime: time.Now(),
				})
				suite.Require().NoError(err)

				err = suite.db.SaveAppToken(types.AppToken{
					AppID: "1",
					Name:  "Test token",
					Value: "token",
				})
				suite.Require().NoError(err)
			},
			setupContext: func(ctx context.Context) context.Context {
				return authentication.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *links.CreateAddressLinkRequest {
				return &links.CreateAddressLinkRequest{
					Address: "",
					Chain:   links.ChainType_MAINNET,
				}
			},
			shouldErr: true,
		},
		{
			name: "invalid chain returns error",
			setup: func() {
				err := suite.db.SaveAppSubscription(types.NewApplicationSubscription(
					1,
					"Test App Subscription",
					10,
					10,
					10,
				))
				suite.Require().NoError(err)

				err = suite.db.SaveApp(types.Application{
					ID:             "1",
					Name:           "Test Application",
					WalletAddress:  "desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					SubscriptionID: 1,
					Admins: []string{
						"desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					},
					CreationTime: time.Now(),
				})
				suite.Require().NoError(err)

				err = suite.db.SaveAppToken(types.AppToken{
					AppID: "1",
					Name:  "Test token",
					Value: "token",
				})
				suite.Require().NoError(err)
			},
			setupContext: func(ctx context.Context) context.Context {
				return authentication.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *links.CreateAddressLinkRequest {
				return &links.CreateAddressLinkRequest{
					Address: "desmos1aph74nw42mk7330pftwwmj7lxr7j3drgmlu3zc",
				}
			},
			shouldErr: true,
		},
		{
			name: "rate limit reached returns error",
			setup: func() {
				// ----------------------------------
				// --- Save the app
				// ----------------------------------
				err := suite.db.SaveAppSubscription(types.NewApplicationSubscription(
					1,
					"Test App Subscription",
					10,
					10,
					1,
				))
				suite.Require().NoError(err)

				err = suite.db.SaveApp(types.Application{
					ID:             "1",
					Name:           "Test Application",
					WalletAddress:  "desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					SubscriptionID: 1,
					Admins: []string{
						"desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					},
					CreationTime: time.Now(),
				})
				suite.Require().NoError(err)

				err = suite.db.SaveAppToken(types.AppToken{
					AppID: "1",
					Name:  "Test token",
					Value: "token",
				})
				suite.Require().NoError(err)

				// ----------------------------------
				// --- Save the existing deep links
				// ----------------------------------
				err = suite.db.SaveCreatedDeepLink(&types.CreatedDeepLink{
					ID:           "1",
					AppID:        "1",
					URL:          "https://example.com",
					CreationTime: time.Now(),
				})
				suite.Require().NoError(err)
			},
			setupContext: func(ctx context.Context) context.Context {
				return authentication.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *links.CreateAddressLinkRequest {
				return &links.CreateAddressLinkRequest{
					Address: "desmos1aph74nw42mk7330pftwwmj7lxr7j3drgmlu3zc",
					Chain:   links.ChainType_MAINNET,
				}
			},
			shouldErr: true,
		},
		{
			name: "valid request works properly",
			setup: func() {
				// ----------------------------------
				// --- Save the app
				// ----------------------------------
				err := suite.db.SaveAppSubscription(types.NewApplicationSubscription(
					1,
					"Test App Subscription",
					10,
					10,
					10,
				))
				suite.Require().NoError(err)

				err = suite.db.SaveApp(types.Application{
					ID:             "1",
					Name:           "Test Application",
					WalletAddress:  "desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					SubscriptionID: 1,
					Admins: []string{
						"desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					},
					CreationTime: time.Now(),
				})
				suite.Require().NoError(err)

				err = suite.db.SaveAppToken(types.AppToken{
					AppID: "1",
					Name:  "Test token",
					Value: "token",
				})
				suite.Require().NoError(err)

				// ----------------------------------
				// --- Setup the mocks
				// ----------------------------------
				suite.deepLinksClient.
					EXPECT().
					CreateDynamicLink(gomock.Any(), gomock.Any()).
					DoAndReturn(func(apiKey string, config *types.LinkConfig) (string, error) {
						// Verify the arguments used
						expectedPath := "/?address=desmos1aph74nw42mk7330pftwwmj7lxr7j3drgmlu3zc&chain_type=mainnet"
						suite.Require().Equal(expectedPath, config.DeepLinking.DeepLinkPath)

						return "https://example.com", nil
					})
			},
			setupContext: func(ctx context.Context) context.Context {
				return authentication.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *links.CreateAddressLinkRequest {
				return &links.CreateAddressLinkRequest{
					Address: "desmos1aph74nw42mk7330pftwwmj7lxr7j3drgmlu3zc",
					Chain:   links.ChainType_MAINNET,
				}
			},
			shouldErr: false,
			check: func(res *links.CreateLinkResponse) {
				// Make sure the returned response is correct
				suite.Require().Equal("https://example.com", res.Url)

				// Make sure the deep link was saved properly
				var count int
				err := suite.db.SQL.Get(&count, `SELECT COUNT(*) FROM deep_links`)
				suite.Require().NoError(err)
				suite.Require().Equal(1, count)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			if tc.setup != nil {
				tc.setup()
			}

			ctx := context.Background()
			if tc.setupContext != nil {
				ctx = tc.setupContext(ctx)
			}

			res, err := suite.client.CreateAddressLink(ctx, tc.buildRequest())
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

func (suite *LinksServerTestSuite) TestCreateViewProfileLink() {
	testCases := []struct {
		name         string
		setup        func()
		setupContext func(ctx context.Context) context.Context
		buildRequest func() *links.CreateViewProfileLinkRequest
		shouldErr    bool
		check        func(res *links.CreateLinkResponse)
	}{
		{
			name: "invalid session returns error",
			buildRequest: func() *links.CreateViewProfileLinkRequest {
				return &links.CreateViewProfileLinkRequest{
					Address: "desmos1aph74nw42mk7330pftwwmj7lxr7j3drgmlu3zc",
					Chain:   links.ChainType_MAINNET,
				}
			},
			shouldErr: true,
		},
		{
			name: "invalid address returns error",
			setup: func() {
				err := suite.db.SaveAppSubscription(types.NewApplicationSubscription(
					1,
					"Test App Subscription",
					10,
					10,
					10,
				))
				suite.Require().NoError(err)

				err = suite.db.SaveApp(types.Application{
					ID:             "1",
					Name:           "Test Application",
					WalletAddress:  "desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					SubscriptionID: 1,
					Admins: []string{
						"desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					},
					CreationTime: time.Now(),
				})
				suite.Require().NoError(err)

				err = suite.db.SaveAppToken(types.AppToken{
					AppID: "1",
					Name:  "Test token",
					Value: "token",
				})
				suite.Require().NoError(err)
			},
			setupContext: func(ctx context.Context) context.Context {
				return authentication.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *links.CreateViewProfileLinkRequest {
				return &links.CreateViewProfileLinkRequest{
					Address: "",
					Chain:   links.ChainType_MAINNET,
				}
			},
			shouldErr: true,
		},
		{
			name: "invalid chain returns error",
			setup: func() {
				err := suite.db.SaveAppSubscription(types.NewApplicationSubscription(
					1,
					"Test App Subscription",
					10,
					10,
					10,
				))
				suite.Require().NoError(err)

				err = suite.db.SaveApp(types.Application{
					ID:             "1",
					Name:           "Test Application",
					WalletAddress:  "desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					SubscriptionID: 1,
					Admins: []string{
						"desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					},
					CreationTime: time.Now(),
				})
				suite.Require().NoError(err)

				err = suite.db.SaveAppToken(types.AppToken{
					AppID: "1",
					Name:  "Test token",
					Value: "token",
				})
				suite.Require().NoError(err)
			},
			setupContext: func(ctx context.Context) context.Context {
				return authentication.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *links.CreateViewProfileLinkRequest {
				return &links.CreateViewProfileLinkRequest{
					Address: "desmos1aph74nw42mk7330pftwwmj7lxr7j3drgmlu3zc",
				}
			},
			shouldErr: true,
		},
		{
			name: "rate limit reached returns error",
			setup: func() {
				// ----------------------------------
				// --- Save the app
				// ----------------------------------
				err := suite.db.SaveAppSubscription(types.NewApplicationSubscription(
					1,
					"Test App Subscription",
					10,
					10,
					1,
				))
				suite.Require().NoError(err)

				err = suite.db.SaveApp(types.Application{
					ID:             "1",
					Name:           "Test Application",
					WalletAddress:  "desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					SubscriptionID: 1,
					Admins: []string{
						"desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					},
					CreationTime: time.Now(),
				})
				suite.Require().NoError(err)

				err = suite.db.SaveAppToken(types.AppToken{
					AppID: "1",
					Name:  "Test token",
					Value: "token",
				})
				suite.Require().NoError(err)

				// ----------------------------------
				// --- Save the existing deep links
				// ----------------------------------
				err = suite.db.SaveCreatedDeepLink(&types.CreatedDeepLink{
					ID:           "1",
					AppID:        "1",
					URL:          "https://example.com",
					CreationTime: time.Now(),
				})
				suite.Require().NoError(err)
			},
			setupContext: func(ctx context.Context) context.Context {
				return authentication.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *links.CreateViewProfileLinkRequest {
				return &links.CreateViewProfileLinkRequest{
					Address: "desmos1aph74nw42mk7330pftwwmj7lxr7j3drgmlu3zc",
					Chain:   links.ChainType_MAINNET,
				}
			},
			shouldErr: true,
		},
		{
			name: "valid request works properly",
			setup: func() {
				// ----------------------------------
				// --- Save the app
				// ----------------------------------
				err := suite.db.SaveAppSubscription(types.NewApplicationSubscription(
					1,
					"Test App Subscription",
					10,
					10,
					10,
				))
				suite.Require().NoError(err)

				err = suite.db.SaveApp(types.Application{
					ID:             "1",
					Name:           "Test Application",
					WalletAddress:  "desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					SubscriptionID: 1,
					Admins: []string{
						"desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					},
					CreationTime: time.Now(),
				})
				suite.Require().NoError(err)

				err = suite.db.SaveAppToken(types.AppToken{
					AppID: "1",
					Name:  "Test token",
					Value: "token",
				})
				suite.Require().NoError(err)

				// ----------------------------------
				// --- Setup the mocks
				// ----------------------------------
				suite.deepLinksClient.
					EXPECT().
					CreateDynamicLink(gomock.Any(), gomock.Any()).
					DoAndReturn(func(apiKey string, config *types.LinkConfig) (string, error) {
						// Verify the arguments used
						expectedPath := "/view_profile?address=desmos1aph74nw42mk7330pftwwmj7lxr7j3drgmlu3zc&chain_type=mainnet"
						suite.Require().Equal(expectedPath, config.DeepLinking.DeepLinkPath)

						return "https://example.com", nil
					})
			},
			setupContext: func(ctx context.Context) context.Context {
				return authentication.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *links.CreateViewProfileLinkRequest {
				return &links.CreateViewProfileLinkRequest{
					Address: "desmos1aph74nw42mk7330pftwwmj7lxr7j3drgmlu3zc",
					Chain:   links.ChainType_MAINNET,
				}
			},
			shouldErr: false,
			check: func(res *links.CreateLinkResponse) {
				// Make sure the returned response is correct
				suite.Require().Equal("https://example.com", res.Url)

				// Make sure the deep link was saved properly
				var count int
				err := suite.db.SQL.Get(&count, `SELECT COUNT(*) FROM deep_links`)
				suite.Require().NoError(err)
				suite.Require().Equal(1, count)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			if tc.setup != nil {
				tc.setup()
			}

			ctx := context.Background()
			if tc.setupContext != nil {
				ctx = tc.setupContext(ctx)
			}

			res, err := suite.client.CreateViewProfileLink(ctx, tc.buildRequest())
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

func (suite *LinksServerTestSuite) TestCreateSendLink() {
	testCases := []struct {
		name         string
		setup        func()
		setupContext func(ctx context.Context) context.Context
		buildRequest func() *links.CreateSendLinkRequest
		shouldErr    bool
		check        func(res *links.CreateLinkResponse)
	}{
		{
			name: "invalid session returns error",
			buildRequest: func() *links.CreateSendLinkRequest {
				return &links.CreateSendLinkRequest{
					Address: "desmos1aph74nw42mk7330pftwwmj7lxr7j3drgmlu3zc",
					Chain:   links.ChainType_MAINNET,
				}
			},
			shouldErr: true,
		},
		{
			name: "invalid address returns error",
			setup: func() {
				err := suite.db.SaveAppSubscription(types.NewApplicationSubscription(
					1,
					"Test App Subscription",
					10,
					10,
					10,
				))
				suite.Require().NoError(err)

				err = suite.db.SaveApp(types.Application{
					ID:             "1",
					Name:           "Test Application",
					WalletAddress:  "desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					SubscriptionID: 1,
					Admins: []string{
						"desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					},
					CreationTime: time.Now(),
				})
				suite.Require().NoError(err)

				err = suite.db.SaveAppToken(types.AppToken{
					AppID: "1",
					Name:  "Test token",
					Value: "token",
				})
				suite.Require().NoError(err)
			},
			setupContext: func(ctx context.Context) context.Context {
				return authentication.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *links.CreateSendLinkRequest {
				return &links.CreateSendLinkRequest{
					Address: "",
					Chain:   links.ChainType_MAINNET,
				}
			},
			shouldErr: true,
		},
		{
			name: "invalid chain returns error",
			setup: func() {
				err := suite.db.SaveAppSubscription(types.NewApplicationSubscription(
					1,
					"Test App Subscription",
					10,
					10,
					10,
				))
				suite.Require().NoError(err)

				err = suite.db.SaveApp(types.Application{
					ID:             "1",
					Name:           "Test Application",
					WalletAddress:  "desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					SubscriptionID: 1,
					Admins: []string{
						"desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					},
					CreationTime: time.Now(),
				})
				suite.Require().NoError(err)

				err = suite.db.SaveAppToken(types.AppToken{
					AppID: "1",
					Name:  "Test token",
					Value: "token",
				})
				suite.Require().NoError(err)
			},
			setupContext: func(ctx context.Context) context.Context {
				return authentication.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *links.CreateSendLinkRequest {
				return &links.CreateSendLinkRequest{
					Address: "desmos1aph74nw42mk7330pftwwmj7lxr7j3drgmlu3zc",
				}
			},
			shouldErr: true,
		},
		{
			name: "invalid amount returns error",
			setup: func() {
				err := suite.db.SaveAppSubscription(types.NewApplicationSubscription(
					1,
					"Test App Subscription",
					10,
					10,
					10,
				))
				suite.Require().NoError(err)

				err = suite.db.SaveApp(types.Application{
					ID:             "1",
					Name:           "Test Application",
					WalletAddress:  "desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					SubscriptionID: 1,
					Admins: []string{
						"desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					},
					CreationTime: time.Now(),
				})
				suite.Require().NoError(err)

				err = suite.db.SaveAppToken(types.AppToken{
					AppID: "1",
					Name:  "Test token",
					Value: "token",
				})
				suite.Require().NoError(err)
			},
			setupContext: func(ctx context.Context) context.Context {
				return authentication.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *links.CreateSendLinkRequest {
				return &links.CreateSendLinkRequest{
					Address: "desmos1aph74nw42mk7330pftwwmj7lxr7j3drgmlu3zc",
					Chain:   links.ChainType_MAINNET,
					Amount:  sdk.Coins{sdk.Coin{Denom: "test", Amount: sdk.NewInt(-1)}},
				}
			},
			shouldErr: true,
		},
		{
			name: "rate limit reached returns error",
			setup: func() {
				// ----------------------------------
				// --- Save the app
				// ----------------------------------
				err := suite.db.SaveAppSubscription(types.NewApplicationSubscription(
					1,
					"Test App Subscription",
					10,
					10,
					1,
				))
				suite.Require().NoError(err)

				err = suite.db.SaveApp(types.Application{
					ID:             "1",
					Name:           "Test Application",
					WalletAddress:  "desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					SubscriptionID: 1,
					Admins: []string{
						"desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					},
					CreationTime: time.Now(),
				})
				suite.Require().NoError(err)

				err = suite.db.SaveAppToken(types.AppToken{
					AppID: "1",
					Name:  "Test token",
					Value: "token",
				})
				suite.Require().NoError(err)

				// ----------------------------------
				// --- Save the existing deep links
				// ----------------------------------
				err = suite.db.SaveCreatedDeepLink(&types.CreatedDeepLink{
					ID:           "1",
					AppID:        "1",
					URL:          "https://example.com",
					CreationTime: time.Now(),
				})
				suite.Require().NoError(err)
			},
			setupContext: func(ctx context.Context) context.Context {
				return authentication.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *links.CreateSendLinkRequest {
				return &links.CreateSendLinkRequest{
					Address: "desmos1aph74nw42mk7330pftwwmj7lxr7j3drgmlu3zc",
					Chain:   links.ChainType_MAINNET,
					Amount:  sdk.NewCoins(sdk.NewCoin("test", sdk.NewInt(10))),
				}
			},
			shouldErr: true,
		},
		{
			name: "valid request works properly",
			setup: func() {
				// ----------------------------------
				// --- Save the app
				// ----------------------------------
				err := suite.db.SaveAppSubscription(types.NewApplicationSubscription(
					1,
					"Test App Subscription",
					10,
					10,
					10,
				))
				suite.Require().NoError(err)

				err = suite.db.SaveApp(types.Application{
					ID:             "1",
					Name:           "Test Application",
					WalletAddress:  "desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					SubscriptionID: 1,
					Admins: []string{
						"desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					},
					CreationTime: time.Now(),
				})
				suite.Require().NoError(err)

				err = suite.db.SaveAppToken(types.AppToken{
					AppID: "1",
					Name:  "Test token",
					Value: "token",
				})
				suite.Require().NoError(err)

				// ----------------------------------
				// --- Setup the mocks
				// ----------------------------------
				suite.deepLinksClient.
					EXPECT().
					CreateDynamicLink(gomock.Any(), gomock.Any()).
					DoAndReturn(func(apiKey string, config *types.LinkConfig) (string, error) {
						// Verify the arguments used
						expectedPath := "/send_tokens?address=desmos1aph74nw42mk7330pftwwmj7lxr7j3drgmlu3zc&amount=10test&chain_type=mainnet"
						suite.Require().Equal(expectedPath, config.DeepLinking.DeepLinkPath)

						return "https://example.com", nil
					})
			},
			setupContext: func(ctx context.Context) context.Context {
				return authentication.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *links.CreateSendLinkRequest {
				return &links.CreateSendLinkRequest{
					Address: "desmos1aph74nw42mk7330pftwwmj7lxr7j3drgmlu3zc",
					Chain:   links.ChainType_MAINNET,
					Amount:  sdk.NewCoins(sdk.NewCoin("test", sdk.NewInt(10))),
				}
			},
			shouldErr: false,
			check: func(res *links.CreateLinkResponse) {
				// Make sure the returned response is correct
				suite.Require().Equal("https://example.com", res.Url)

				// Make sure the deep link was saved properly
				var count int
				err := suite.db.SQL.Get(&count, `SELECT COUNT(*) FROM deep_links`)
				suite.Require().NoError(err)
				suite.Require().Equal(1, count)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			suite.SetupTest()
			if tc.setup != nil {
				tc.setup()
			}

			ctx := context.Background()
			if tc.setupContext != nil {
				ctx = tc.setupContext(ctx)
			}

			res, err := suite.client.CreateSendLink(ctx, tc.buildRequest())
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
