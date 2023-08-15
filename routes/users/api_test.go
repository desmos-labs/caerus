package users_test

import (
	"context"
	"path"
	"testing"
	"time"

	cosmosclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	wallettypes "github.com/desmos-labs/cosmos-go-wallet/types"
	desmosapp "github.com/desmos-labs/desmos/v5/app"
	profilestypes "github.com/desmos-labs/desmos/v5/x/profiles/types"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	client "github.com/desmos-labs/caerus/chain"
	"github.com/desmos-labs/caerus/database"
	"github.com/desmos-labs/caerus/routes/users"
	"github.com/desmos-labs/caerus/testutils"
	"github.com/desmos-labs/caerus/types"
)

func TestLoginAPIsTestSuite(t *testing.T) {
	suite.Run(t, new(LoginAPIsTestSuite))
}

type LoginAPIsTestSuite struct {
	suite.Suite

	txConfig cosmosclient.TxConfig
	cdc      codec.Codec
	amino    *codec.LegacyAmino

	db      *database.Database
	handler *users.Handler

	server *grpc.Server
	client users.UsersServiceClient

	chainCfg  *wallettypes.ChainConfig
	apiClient *client.Client
}

func (suite *LoginAPIsTestSuite) SetupSuite() {
	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	// Setup the Cosmos SDK config
	desmosapp.SetupConfig(sdk.GetConfig())

	// Build the database
	db, err := testutils.CreateDatabase(path.Join("..", "..", "database", "schema"))
	suite.Require().NoError(err)
	suite.db = db

	// Build chain-related stuff
	encodingConfig := desmosapp.MakeEncodingConfig()
	suite.txConfig, suite.cdc, suite.amino = encodingConfig.TxConfig, encodingConfig.Codec, encodingConfig.Amino

	suite.chainCfg = &wallettypes.ChainConfig{
		Bech32Prefix: "desmos",
		RPCAddr:      "http://localhost:26657",
		GRPCAddr:     "http://localhost:9090",
		GasPrice:     "0.01stake",
	}

	// Build the API chain client
	apiWallet, err := client.NewClient(&client.Config{
		Account: &wallettypes.AccountConfig{
			Mnemonic: "hour harbor fame unaware bunker junk garment decrease federal vicious island smile warrior fame right suit portion skate analyst multiply magnet medal fresh sweet",
			HDPath:   "m/44'/852'/0'/0/0",
		},
		Chain: suite.chainCfg,
		FeeGrant: &client.FeeGrantConfig{
			MsgTypes:   []string{sdk.MsgTypeURL(&profilestypes.MsgSaveProfile{})},
			GrantLimit: sdk.NewCoins(sdk.NewCoin("stake", sdk.NewInt(1000000))),
			Expiration: 1 * time.Hour,
		},
	}, suite.txConfig, suite.cdc)
	suite.Require().NoError(err)
	suite.apiClient = apiWallet

	// Create the handler
	suite.handler = users.NewHandler(suite.cdc, suite.amino, suite.db)

	// Create the server
	suite.server = testutils.CreateServer(suite.db)

	// Register the service
	service := users.NewServer(suite.cdc, suite.amino, suite.db)
	users.RegisterUsersServiceServer(suite.server, service)

	// Setup the client
	conn, err := testutils.StartServerAndConnect(suite.server)
	suite.Require().NoError(err)
	suite.client = users.NewUsersServiceClient(conn)
}

func (suite *LoginAPIsTestSuite) TearDownSuite() {
	suite.server.Stop()
}

func (suite *LoginAPIsTestSuite) SetupTest() {
	// Truncate all the database data to make sure we have a clean database state
	err := testutils.TruncateDatabase(suite.db)
	suite.Require().NoError(err)

	// Remove existing grants to make sure we have a clean chain state
	suite.deleteGrants()
}

func (suite *LoginAPIsTestSuite) deleteGrants() {
	var sdkMsg []sdk.Msg

	// Remove all the feegrants made to the user
	feegrantClient := feegrant.NewQueryClient(suite.apiClient.Client.GRPCConn)
	allowancesRes, err := feegrantClient.AllowancesByGranter(context.Background(), &feegrant.QueryAllowancesByGranterRequest{
		Granter: suite.apiClient.AccAddress(),
	})
	suite.Require().NoError(err)

	for _, allowance := range allowancesRes.Allowances {
		granterAddr, err := sdk.AccAddressFromBech32(allowance.Granter)
		suite.Require().NoError(err)

		granteeAddr, err := sdk.AccAddressFromBech32(allowance.Grantee)
		suite.Require().NoError(err)

		revokeMsg := feegrant.NewMsgRevokeAllowance(granterAddr, granteeAddr)

		sdkMsg = append(sdkMsg, &revokeMsg)
	}

	if len(sdkMsg) == 0 {
		return
	}

	// Broadcast the transaction
	response, err := suite.apiClient.BroadcastTxCommit(&wallettypes.TransactionData{Messages: sdkMsg, GasAuto: true, FeeAuto: true})
	suite.Require().NoError(err)
	suite.Require().Zero(response.Code)
	suite.Require().Zerof(response.Code, response.RawLog)
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *LoginAPIsTestSuite) TestGetNonce() {
	testCases := []struct {
		name         string
		buildRequest func() *users.GetNonceRequest
		shouldErr    bool
		check        func(res *users.GetNonceResponse)
	}{
		{
			name: "invalid request returns error",
			buildRequest: func() *users.GetNonceRequest {
				return &users.GetNonceRequest{}
			},
			shouldErr: true,
		},
		{
			name: "valid request returns no error",
			buildRequest: func() *users.GetNonceRequest {
				return &users.GetNonceRequest{
					UserDesmosAddress: "desmos1n5xcfgnhd28uwyqwuy6ysf05x9hf04r0ydwxjt",
				}
			},
			shouldErr: false,
			check: func(res *users.GetNonceResponse) {
				// Make sure the nonce exists
				encryptedNonce, err := suite.db.GetNonce("desmos1n5xcfgnhd28uwyqwuy6ysf05x9hf04r0ydwxjt", res.Nonce)
				suite.Require().NoError(err)
				suite.Require().NotNil(encryptedNonce)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			res, err := suite.client.GetNonce(context.Background(), tc.buildRequest())

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

func (suite *LoginAPIsTestSuite) TestLogin() {
	testCases := []struct {
		name         string
		setup        func()
		buildRequest func() *types.SignedRequest
		shouldErr    bool
		check        func(res *users.LoginResponse)
	}{
		{
			name: "invalid signature returns error",
			buildRequest: func() *types.SignedRequest {
				return &types.SignedRequest{
					DesmosAddress:  "desmos1tamzg6rfj9wlmqhthgfmn9awq0d8ssgfr8fjns",
					PubKeyBytes:    "0a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912230a21033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d",
					SignedBytes:    "0a95010a8b010a1c2f636f736d6f732e62616e6b2e763162657461312e4d736753656e64126b0a2d6465736d6f73313379703266713374736c71366d6d747134363238713338787a6a37356574687a656c61397575122d6465736d6f733174616d7a673672666a39776c6d7168746867666d6e39617771306438737367667238666a6e731a0b0a0675646172696312013112056e6f6e636512560a4e0a460a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912230a21033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d12040a020801120410c09a0c1a066465736d6f73",
					SignatureBytes: "7acb63db1c0923a6e480a553ff86c972e4fe14226bc4eaece510e94b5b6d2f2716b7e5391da79f8d159ccab9747aa92d2dfdd7eefad0e936b6ff7ee26ba97168",
				}
			},
			shouldErr: true,
		},
		{
			name: "wrong nonce returns error",
			setup: func() {
				err := suite.db.SaveNonce(types.NewNonce(
					"desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu",
					"nonce",
					time.Now().Add(time.Hour*24),
				))
				suite.Require().NoError(err)
			},
			buildRequest: func() *types.SignedRequest {
				return &types.SignedRequest{
					DesmosAddress:  "desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu",
					PubKeyBytes:    "0a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912230a21033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d",
					SignedBytes:    "0a9b010a8b010a1c2f636f736d6f732e62616e6b2e763162657461312e4d736753656e64126b0a2d6465736d6f73313379703266713374736c71366d6d747134363238713338787a6a37356574687a656c61397575122d6465736d6f733174616d7a673672666a39776c6d7168746867666d6e39617771306438737367667238666a6e731a0b0a06756461726963120131120b77726f6e672d6e6f6e636512560a4e0a460a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912230a21033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d12040a020801120410c09a0c1a066465736d6f73",
					SignatureBytes: "5948936ebecc59ac6c7d5367052b1bedc6b0c2e4f4f6bd7915e07c73bef5d73e285a56126a4e294b041dc6007c24fa6c1f2fe30184deaa706da8748495410f88",
				}
			},
			shouldErr: true,
		},
		{
			name: "expired nonce returns error",
			setup: func() {
				err := suite.db.SaveNonce(types.NewNonce(
					"desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu",
					"nonce",
					time.Date(2020, 1, 1, 12, 00, 00, 000, time.UTC),
				))
				suite.Require().NoError(err)
			},
			buildRequest: func() *types.SignedRequest {
				return &types.SignedRequest{
					DesmosAddress:  "desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu",
					PubKeyBytes:    "0a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912230a21033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d",
					SignedBytes:    "0a95010a8b010a1c2f636f736d6f732e62616e6b2e763162657461312e4d736753656e64126b0a2d6465736d6f73313379703266713374736c71366d6d747134363238713338787a6a37356574687a656c61397575122d6465736d6f733174616d7a673672666a39776c6d7168746867666d6e39617771306438737367667238666a6e731a0b0a0675646172696312013112056e6f6e636512560a4e0a460a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912230a21033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d12040a020801120410c09a0c1a066465736d6f73",
					SignatureBytes: "7acb63db1c0923a6e480a553ff86c972e4fe14226bc4eaece510e94b5b6d2f2716b7e5391da79f8d159ccab9747aa92d2dfdd7eefad0e936b6ff7ee26ba97168",
				}
			},
			shouldErr: true,
			check: func(res *users.LoginResponse) {
				// Make sure the nonce has been deleted
				encryptedNonce, err := suite.db.GetNonce("desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu", "nonce")
				suite.Require().NoError(err)
				suite.Require().Nil(encryptedNonce)
			},
		},
		{
			name: "valid request returns no error - normal account",
			setup: func() {
				err := suite.db.SaveNonce(types.NewNonce(
					"desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu",
					"nonce",
					time.Now().Add(time.Hour*24),
				))
				suite.Require().NoError(err)
			},
			buildRequest: func() *types.SignedRequest {
				return &types.SignedRequest{
					DesmosAddress:  "desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu",
					PubKeyBytes:    "0a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912230a21033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d",
					SignedBytes:    "0a95010a8b010a1c2f636f736d6f732e62616e6b2e763162657461312e4d736753656e64126b0a2d6465736d6f73313379703266713374736c71366d6d747134363238713338787a6a37356574687a656c61397575122d6465736d6f733174616d7a673672666a39776c6d7168746867666d6e39617771306438737367667238666a6e731a0b0a0675646172696312013112056e6f6e636512560a4e0a460a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912230a21033024e9e0ad4f93045ef5a60bb92171e6418cd13b082e7a7bc3ed05312a0b417d12040a020801120410c09a0c1a066465736d6f73",
					SignatureBytes: "7acb63db1c0923a6e480a553ff86c972e4fe14226bc4eaece510e94b5b6d2f2716b7e5391da79f8d159ccab9747aa92d2dfdd7eefad0e936b6ff7ee26ba97168",
				}
			},
			shouldErr: false,
			check: func(res *users.LoginResponse) {
				suite.Require().NotEmpty(res.Token)

				// Make sure the nonce has been deleted
				encryptedNonce, err := suite.db.GetNonce("desmos13yp2fq3tslq6mmtq4628q38xzj75ethzela9uu", "nonce")
				suite.Require().NoError(err)
				suite.Require().Nil(encryptedNonce)

				// Make sure the session has been created
				encryptedSession, err := suite.db.GetUserSession(res.Token)
				suite.Require().NoError(err)
				suite.Require().NotNil(encryptedSession)

				// Make sure the user data has been inserted properly
				var count int
				err = suite.db.SQL.Get(&count, `SELECT COUNT(*) FROM users`)
				suite.Require().NoError(err)
				suite.Require().Equal(1, count)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			if tc.setup != nil {
				tc.setup()
			}

			res, err := suite.client.Login(context.Background(), tc.buildRequest())

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

func (suite *LoginAPIsTestSuite) TestRefreshSession() {
	testCases := []struct {
		name         string
		setup        func()
		setupContext func(ctx context.Context) context.Context
		shouldErr    bool
		check        func(res *users.RefreshSessionResponse)
	}{
		{
			name: "missing session returns an error",
			setupContext: func(ctx context.Context) context.Context {
				return testutils.SetupContextWithAuthorization(ctx, "token")
			},
			shouldErr: true,
		},
		{
			name: "expired session returns an error",
			setup: func() {
				err := suite.db.SaveSession(types.NewUserSession(
					"desmos1c7ms9zhtgwmv5jy6ztj2vq0jj67zenw3gdl2gr",
					"token",
					time.Date(2020, 1, 1, 12, 00, 00, 000, time.UTC),
					time.Date(2020, 1, 2, 12, 00, 00, 000, time.UTC),
				))
				suite.Require().NoError(err)
			},
			setupContext: func(ctx context.Context) context.Context {
				return testutils.SetupContextWithAuthorization(ctx, "token")
			},
			shouldErr: true,
			check: func(res *users.RefreshSessionResponse) {
				// Make sure the expired session is deleted
				encryptedSession, err := suite.db.GetUserSession("token")
				suite.Require().NoError(err)
				suite.Require().Nil(encryptedSession)
			},
		},
		{
			name: "valid request returns no error - normal user",
			setup: func() {
				err := suite.db.SaveSession(types.NewUserSession(
					"desmos1c7ms9zhtgwmv5jy6ztj2vq0jj67zenw3gdl2gr",
					"token",
					time.Now(),
					time.Now().Add(time.Hour),
				))
				suite.Require().NoError(err)
			},
			setupContext: func(ctx context.Context) context.Context {
				return testutils.SetupContextWithAuthorization(ctx, "token")
			},
			shouldErr: false,
			check: func(res *users.RefreshSessionResponse) {
				// Make sure the session has been refreshed and the expiry time has been updated
				suite.Require().True(res.ExpirationTime.After(time.Now()))

				// Make sure the session is refreshed
				encryptedSession, err := suite.db.GetUserSession("token")
				suite.Require().NoError(err)
				suite.Require().NotNil(encryptedSession)
				suite.Require().Greater(encryptedSession.ExpiryTime.Sub(time.Now()), time.Hour)
			},
		},
		{
			name: "valid request returns no error - deleted user",
			setup: func() {
				err := suite.db.SaveSession(types.NewUserSession(
					"desmos1c7ms9zhtgwmv5jy6ztj2vq0jj67zenw3gdl2gr",
					"token",
					time.Now(),
					time.Now().Add(time.Hour),
				))
				suite.Require().NoError(err)
			},
			setupContext: func(ctx context.Context) context.Context {
				return testutils.SetupContextWithAuthorization(ctx, "token")
			},
			shouldErr: false,
			check: func(res *users.RefreshSessionResponse) {
				// Make sure the session has been refreshed and the expiry time has been updated
				suite.Require().True(res.ExpirationTime.After(time.Now()))

				// Make sure the session is refreshed
				encryptedSession, err := suite.db.GetUserSession("token")
				suite.Require().NoError(err)
				suite.Require().NotNil(encryptedSession)
				suite.Require().Greater(encryptedSession.ExpiryTime.Sub(time.Now()), time.Hour)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		suite.Run(tc.name, func() {
			if tc.setup != nil {
				tc.setup()
			}

			ctx := context.Background()
			if tc.setupContext != nil {
				ctx = tc.setupContext(ctx)
			}

			res, err := suite.client.RefreshSession(ctx, &emptypb.Empty{})

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

func (suite *LoginAPIsTestSuite) TestLogout() {
	testCases := []struct {
		name         string
		setup        func()
		setupContext func(ctx context.Context) context.Context
		buildRequest func() *users.LogoutRequest
		shouldErr    bool
		check        func()
	}{
		{
			name: "logout from single session works properly",
			setup: func() {
				err := suite.db.SaveSession(types.NewUserSession(
					"desmos1c7ms9zhtgwmv5jy6ztj2vq0jj67zenw3gdl2gr",
					"token",
					time.Now(),
					time.Now().Add(time.Hour),
				))
				suite.Require().NoError(err)
			},
			setupContext: func(ctx context.Context) context.Context {
				return testutils.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *users.LogoutRequest {
				return &users.LogoutRequest{LogoutFromAll: false}
			},
			shouldErr: false,
			check: func() {
				// Make sure the session has been deleted
				var count int
				err := suite.db.SQL.Get(&count, `SELECT COUNT(*) FROM sessions`)
				suite.Require().NoError(err)
				suite.Require().Zero(count)
			},
		},
		{
			name: "logout from all sessions works properly",
			setup: func() {
				err := suite.db.SaveSession(types.NewUserSession(
					"desmos1c7ms9zhtgwmv5jy6ztj2vq0jj67zenw3gdl2gr",
					"token",
					time.Now(),
					time.Now().Add(time.Hour),
				))
				suite.Require().NoError(err)

				err = suite.db.SaveSession(types.NewUserSession(
					"desmos1c7ms9zhtgwmv5jy6ztj2vq0jj67zenw3gdl2gr",
					"second-token",
					time.Now().Add(-time.Minute),
					time.Now().Add(-time.Minute+time.Hour),
				))
				suite.Require().NoError(err)
			},
			setupContext: func(ctx context.Context) context.Context {
				return testutils.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *users.LogoutRequest {
				return &users.LogoutRequest{LogoutFromAll: true}
			},
			shouldErr: false,
			check: func() {
				// Make sure all sessions have been deleted
				var count int
				err := suite.db.SQL.Get(&count, `SELECT COUNT(*) FROM sessions`)
				suite.Require().NoError(err)
				suite.Require().Zero(count)
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

			_, err := suite.client.Logout(ctx, tc.buildRequest())

			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}

			if tc.check != nil {
				tc.check()
			}
		})
	}
}

func (suite *LoginAPIsTestSuite) TestDeleteAccount() {
	testCases := []struct {
		name         string
		setup        func()
		setupContext func(ctx context.Context) context.Context
		shouldErr    bool
		check        func()
	}{
		{
			name: "invalid session returns error",
			setup: func() {
				err := suite.db.SaveSession(types.NewUserSession(
					"desmos1c7ms9zhtgwmv5jy6ztj2vq0jj67zenw3gdl2gr",
					"token",
					time.Date(2020, 1, 1, 12, 00, 00, 000, time.UTC),
					time.Date(2020, 1, 2, 12, 00, 00, 000, time.UTC),
				))
				suite.Require().NoError(err)
			},
			setupContext: func(ctx context.Context) context.Context {
				return testutils.SetupContextWithAuthorization(ctx, "token")
			},
			shouldErr: true,
		},
		{
			name: "valid request works properly",
			setup: func() {
				err := suite.db.SaveSession(types.NewUserSession(
					"desmos1c7ms9zhtgwmv5jy6ztj2vq0jj67zenw3gdl2gr",
					"token",
					time.Now(),
					time.Now().Add(time.Hour),
				))
				suite.Require().NoError(err)

				err = suite.db.SaveSession(types.NewUserSession(
					"desmos1c7ms9zhtgwmv5jy6ztj2vq0jj67zenw3gdl2gr",
					"second-token",
					time.Now().Add(-time.Minute),
					time.Now().Add(-time.Minute+time.Hour),
				))
				suite.Require().NoError(err)
			},
			setupContext: func(ctx context.Context) context.Context {
				return testutils.SetupContextWithAuthorization(ctx, "token")
			},
			shouldErr: false,
			check: func() {
				// Make sure all sessions have been deleted
				var count int
				err := suite.db.SQL.Get(&count, `SELECT COUNT(*) FROM sessions`)
				suite.Require().NoError(err)
				suite.Require().Zero(count)
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

			_, err := suite.client.DeleteAccount(ctx, &emptypb.Empty{})

			if tc.shouldErr {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}

			if tc.check != nil {
				tc.check()
			}
		})
	}
}

// --------------------------------------------------------------------------------------------------------------------

//func (suite *LoginAPIsTestSuite) TestHasuraSession() {
//	testCases := []struct {
//		name         string
//		setup        func()
//		buildRequest func() (*http.Request, error)
//		shouldErr    bool
//		check        func(w *httptest.ResponseRecorder)
//	}{
//		{
//			name: "missing session returns unauthorized response",
//			buildRequest: func() (*http.Request, error) {
//				req, err := http.NewRequest("GET", "/hasura-session", nil)
//				if err != nil {
//					return nil, err
//				}
//
//				req.Header.Add("Authorization", "Bearer token")
//				return req, nil
//			},
//			shouldErr: false,
//			check: func(w *httptest.ResponseRecorder) {
//				// Check the response
//				resBz, err := io.ReadAll(w.Body)
//				suite.Require().NoError(err)
//
//				var res users.HasuraSessionResponseJSON
//				err = json.Unmarshal(resBz, &res)
//				suite.Require().NoError(err)
//
//				suite.Require().Equal(users.UnauthorizedUserRole, res.UserRole)
//				suite.Require().Equal("", res.UserAddress)
//			},
//		},
//		{
//			name: "expired session returns unauthorized response",
//			setup: func() {
//				err := suite.db.SaveSession(types.NewUserSession(
//					"desmos1c7ms9zhtgwmv5jy6ztj2vq0jj67zenw3gdl2gr",
//					"token",
//					time.Date(2020, 1, 1, 12, 00, 00, 000, time.UTC),
//					time.Date(2020, 1, 2, 12, 00, 00, 000, time.UTC),
//				))
//				suite.Require().NoError(err)
//			},
//			buildRequest: func() (*http.Request, error) {
//				req, err := http.NewRequest("GET", "/hasura-session", nil)
//				if err != nil {
//					return nil, err
//				}
//
//				req.Header.Add("Authorization", "Bearer token")
//				return req, nil
//			},
//			shouldErr: false,
//			check: func(w *httptest.ResponseRecorder) {
//				// Make sure the expired session is deleted
//				encryptedSession, err := suite.db.GetUserSession("token")
//				suite.Require().NoError(err)
//				suite.Require().Nil(encryptedSession)
//
//				// Check the response
//				resBz, err := io.ReadAll(w.Body)
//				suite.Require().NoError(err)
//
//				var res users.HasuraSessionResponseJSON
//				err = json.Unmarshal(resBz, &res)
//				suite.Require().NoError(err)
//
//				suite.Require().Equal(users.UnauthorizedUserRole, res.UserRole)
//				suite.Require().Equal("", res.UserAddress)
//			},
//		},
//		{
//			name: "authorized user returns correct data",
//			setup: func() {
//				err := suite.db.SaveSession(types.NewUserSession(
//					"desmos1c7ms9zhtgwmv5jy6ztj2vq0jj67zenw3gdl2gr",
//					"token",
//					time.Now(),
//					time.Now().Add(time.Hour),
//				))
//				suite.Require().NoError(err)
//			},
//			buildRequest: func() (*http.Request, error) {
//				req, err := http.NewRequest("GET", "/hasura-session", nil)
//				if err != nil {
//					return nil, err
//				}
//
//				req.Header.Add("Authorization", "Bearer token")
//				return req, nil
//			},
//			shouldErr: false,
//			check: func(w *httptest.ResponseRecorder) {
//				// Make sure the session is refreshed
//				encryptedSession, err := suite.db.GetUserSession("token")
//				suite.Require().NoError(err)
//				suite.Require().NotNil(encryptedSession)
//				suite.Require().Greater(encryptedSession.ExpiryTime.Sub(time.Now()), time.Hour)
//
//				// Check the response
//				resBz, err := io.ReadAll(w.Body)
//				suite.Require().NoError(err)
//
//				var res users.HasuraSessionResponseJSON
//				err = json.Unmarshal(resBz, &res)
//				suite.Require().NoError(err)
//
//				suite.Require().Equal(users.AuthenticatedUserRole, res.UserRole)
//				suite.Require().Equal("desmos1c7ms9zhtgwmv5jy6ztj2vq0jj67zenw3gdl2gr", res.UserAddress)
//			},
//		},
//	}
//
//	for _, tc := range testCases {
//		tc := tc
//		suite.Run(tc.name, func() {
//			suite.SetupTest()
//			if tc.setup != nil {
//				tc.setup()
//			}
//
//			req, err := tc.buildRequest()
//			suite.Require().NoError(err)
//
//			w := httptest.NewRecorder()
//			suite.r.ServeHTTP(w, req)
//
//			if tc.shouldErr {
//				suite.Require().NotEqual(http.StatusOK, w.Code)
//			} else {
//				suite.Require().Equal(http.StatusOK, w.Code)
//			}
//
//			if tc.check != nil {
//				tc.check(w)
//			}
//		})
//	}
//}
