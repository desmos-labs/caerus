package grants_test

import (
	"context"
	"path"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/desmos-labs/cosmos-go-wallet/client"
	wallettypes "github.com/desmos-labs/cosmos-go-wallet/types"
	"github.com/desmos-labs/cosmos-go-wallet/wallet"
	desmosapp "github.com/desmos-labs/desmos/v6/app"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"

	"github.com/desmos-labs/caerus/authentication"
	"github.com/desmos-labs/caerus/chain"
	"github.com/desmos-labs/caerus/database"
	"github.com/desmos-labs/caerus/routes/grants"
	"github.com/desmos-labs/caerus/testutils"
	"github.com/desmos-labs/caerus/types"
)

func TestGrantsServerTestSuite(t *testing.T) {
	suite.Run(t, new(GrantsServerTestSuite))
}

type GrantsServerTestSuite struct {
	suite.Suite

	cdc     codec.Codec
	db      *database.Database
	handler *grants.Handler

	server *grpc.Server
	client grants.GrantsServiceClient

	serverChainClient *chain.Client
	appWallet         *wallet.Wallet
}

func (suite *GrantsServerTestSuite) SetupSuite() {
	// Setup the Cosmos SDK config
	desmosapp.SetupConfig(sdk.GetConfig())

	// Build the database
	db, err := testutils.CreateDatabase(path.Join("..", "..", "database", "schema"))
	suite.Require().NoError(err)
	suite.db = db

	// Build chain-related stuff
	encodingConfig := desmosapp.MakeEncodingConfig()
	txConfig, cdc := encodingConfig.TxConfig, encodingConfig.Codec
	suite.cdc = cdc

	chainCfg := &wallettypes.ChainConfig{
		Bech32Prefix: "desmos",
		RPCAddr:      "http://localhost:26657",
		GRPCAddr:     "http://localhost:9090",
		GasPrice:     "0.01stake",
	}

	// Build the server chain client
	serverClient, err := chain.NewClient(&chain.Config{
		Account: &wallettypes.AccountConfig{
			Mnemonic: "hour harbor fame unaware bunker junk garment decrease federal vicious island smile warrior fame right suit portion skate analyst multiply magnet medal fresh sweet",
			HDPath:   "m/44'/852'/0'/0/0",
		},
		Chain: chainCfg,
	}, txConfig, cdc)
	suite.Require().NoError(err)
	suite.serverChainClient = serverClient

	// Build the app wallet
	walletClient, err := client.NewClient(chainCfg, cdc)
	suite.Require().NoError(err)

	appWallet, err := wallet.NewWallet(&wallettypes.AccountConfig{
		Mnemonic: "goose kite term stamp swallow agree fruit culture cart decorate fame start decrease piece random wheel chuckle course oyster enforce jeans attract inquiry beyond",
		HDPath:   "m/44'/852'/0'/0/0",
	}, walletClient, txConfig)
	suite.Require().NoError(err)
	suite.appWallet = appWallet

	// Create the handler
	suite.handler = grants.NewHandler(suite.serverChainClient, suite.cdc, suite.db)

	// Create the server
	suite.server = testutils.CreateServer(suite.db)

	// Register the service
	service := grants.NewServer(suite.handler)
	grants.RegisterGrantsServiceServer(suite.server, service)

	// Setup the client
	conn, err := testutils.StartServerAndConnect(suite.server)
	suite.Require().NoError(err)
	suite.client = grants.NewGrantsServiceClient(conn)
}

func (suite *GrantsServerTestSuite) TearDownSuite() {
	suite.server.Stop()
}

func (suite *GrantsServerTestSuite) SetupTest() {
	// Truncate all the database data to make sure we have a clean database state
	err := testutils.TruncateDatabase(suite.db)
	suite.Require().NoError(err)

	// Remove existing grants to make sure we have a clean chain state
	suite.deleteGrants()
}

func (suite *GrantsServerTestSuite) deleteGrants() {
	// -------------------------------------------------------------------
	// --- Remove all the fee grants given to the user from the app
	// -------------------------------------------------------------------

	feegrantClient := feegrant.NewQueryClient(suite.appWallet.Client.GRPCConn)
	allowancesRes, err := feegrantClient.AllowancesByGranter(context.Background(), &feegrant.QueryAllowancesByGranterRequest{
		Granter: suite.appWallet.AccAddress(),
	})
	suite.Require().NoError(err)

	sdkMsgs := make([]sdk.Msg, len(allowancesRes.Allowances))
	for i, allowance := range allowancesRes.Allowances {
		granterAddr, err := sdk.AccAddressFromBech32(allowance.Granter)
		suite.Require().NoError(err)

		granteeAddr, err := sdk.AccAddressFromBech32(allowance.Grantee)
		suite.Require().NoError(err)

		revokeMsg := feegrant.NewMsgRevokeAllowance(granterAddr, granteeAddr)

		sdkMsgs[i] = &revokeMsg
	}

	if len(sdkMsgs) == 0 {
		return
	}

	// Broadcast the transaction
	res, err := suite.appWallet.BroadcastTxCommit(&wallettypes.TransactionData{
		Messages: sdkMsgs,
		GasAuto:  true,
		FeeAuto:  true,
	})
	suite.Require().NoError(err)
	suite.Require().Zero(res.Code)
	suite.Require().Zerof(res.Code, res.RawLog)

	// -------------------------------------------------------------------------------
	// --- Remove all the on-chain authorizations given to the server from the app
	// -------------------------------------------------------------------------------

	authzClient := authz.NewQueryClient(suite.appWallet.Client.GRPCConn)
	authorizationsRes, err := authzClient.GranterGrants(context.Background(), &authz.QueryGranterGrantsRequest{
		Granter: suite.appWallet.AccAddress(),
	})
	suite.Require().NoError(err)

	sdkMsgs = make([]sdk.Msg, len(authorizationsRes.Grants))
	for i, grant := range authorizationsRes.Grants {
		granterAddress, err := sdk.AccAddressFromBech32(grant.Granter)
		suite.Require().NoError(err)

		granteeAddress, err := sdk.AccAddressFromBech32(grant.Grantee)
		suite.Require().NoError(err)

		var authorization authz.Authorization
		err = suite.cdc.UnpackAny(grant.Authorization, &authorization)
		suite.Require().NoError(err)

		revokeMsg := authz.NewMsgRevoke(granterAddress, granteeAddress, authorization.MsgTypeURL())

		sdkMsgs[i] = &revokeMsg
	}

	if len(sdkMsgs) == 0 {
		return
	}

	// Broadcast the transaction
	res, err = suite.appWallet.BroadcastTxCommit(&wallettypes.TransactionData{
		Messages: sdkMsgs,
		GasAuto:  true,
		FeeAuto:  true,
	})
	suite.Require().NoError(err)
	suite.Require().Zero(res.Code)
	suite.Require().Zerof(res.Code, res.RawLog)
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *GrantsServerTestSuite) grantOnChainExecAuthorization() {
	// Create the authorization
	authorization := authz.NewGenericAuthorization(sdk.MsgTypeURL(&feegrant.MsgGrantAllowance{}))

	// Get the addresses
	granterAddress, err := sdk.AccAddressFromBech32(suite.appWallet.AccAddress())
	suite.Require().NoError(err)

	granteeAddress, err := sdk.AccAddressFromBech32(suite.serverChainClient.AccAddress())
	suite.Require().NoError(err)

	// Create the authorization message
	msgAuth, err := authz.NewMsgGrant(granterAddress, granteeAddress, authorization, nil)
	suite.Require().NoError(err)

	// Execute the transaction
	res, err := suite.appWallet.BroadcastTxCommit(&wallettypes.TransactionData{
		Messages: []sdk.Msg{msgAuth},
		GasAuto:  true,
		FeeAuto:  true,
	})
	suite.Require().NoError(err)
	suite.Require().Zero(res.Code)
	suite.Require().Zerof(res.Code, res.RawLog)
}

func (suite *GrantsServerTestSuite) sendFundsToUser(address string) {
	// Get the addresses
	senderAddress, err := sdk.AccAddressFromBech32(suite.appWallet.AccAddress())
	suite.Require().NoError(err)

	receiverAddress, err := sdk.AccAddressFromBech32(address)
	suite.Require().NoError(err)

	// Build the message
	amount := sdk.NewCoins(sdk.NewCoin(suite.appWallet.Client.GasPrice.Denom, sdk.NewInt(10)))
	sendMsg := banktypes.NewMsgSend(senderAddress, receiverAddress, amount)

	// Execute the transaction
	res, err := suite.appWallet.BroadcastTxCommit(&wallettypes.TransactionData{
		Messages: []sdk.Msg{sendMsg},
		GasAuto:  true,
		FeeAuto:  true,
	})
	suite.Require().NoError(err)
	suite.Require().Zero(res.Code)
	suite.Require().Zerof(res.Code, res.RawLog)
}

func (suite *GrantsServerTestSuite) grantFeeAllowanceToUser(address string) {
	// Get the addresses
	granterAddress, err := sdk.AccAddressFromBech32(suite.appWallet.AccAddress())
	suite.Require().NoError(err)

	granteeAddress, err := sdk.AccAddressFromBech32(address)
	suite.Require().NoError(err)

	// Build the message
	msg, err := feegrant.NewMsgGrantAllowance(&feegrant.BasicAllowance{}, granterAddress, granteeAddress)
	suite.Require().NoError(err)

	// Execute the transaction
	res, err := suite.appWallet.BroadcastTxCommit(&wallettypes.TransactionData{
		Messages: []sdk.Msg{msg},
		GasAuto:  true,
		FeeAuto:  true,
	})
	suite.Require().NoError(err)
	suite.Require().Zero(res.Code)
	suite.Require().Zerof(res.Code, res.RawLog)
}

func (suite *GrantsServerTestSuite) TestRequestFeeAllowance() {
	testCases := []struct {
		name         string
		setup        func()
		setupContext func(ctx context.Context) context.Context
		buildRequest func() *grants.RequestFeeAllowanceRequest
		shouldErr    bool
		check        func()
	}{
		{
			name: "invalid app token returns error",
			setupContext: func(ctx context.Context) context.Context {
				return authentication.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *grants.RequestFeeAllowanceRequest {
				allowanceAny, err := codectypes.NewAnyWithValue(&feegrant.BasicAllowance{})
				suite.Require().NoError(err)

				return &grants.RequestFeeAllowanceRequest{
					UserDesmosAddress: "desmos1jvzuh3vzket6hlsdus5kvq7ne3fpju9t3tp84y",
					Allowance:         allowanceAny,
				}
			},
			shouldErr: true,
		},
		{
			name: "on-chain MsgGrantFeeAllowance authorization not found returns error",
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
					ID:                      "1",
					Name:                    "Test Application",
					WalletAddress:           suite.appWallet.AccAddress(),
					SubscriptionID:          1,
					SecretKey:               "secret",
					NotificationsWebhookURL: "https://example.com",
					Admins: []string{
						suite.appWallet.AccAddress(),
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
			buildRequest: func() *grants.RequestFeeAllowanceRequest {
				allowanceAny, err := codectypes.NewAnyWithValue(&feegrant.BasicAllowance{})
				suite.Require().NoError(err)

				return &grants.RequestFeeAllowanceRequest{
					UserDesmosAddress: "desmos1jvzuh3vzket6hlsdus5kvq7ne3fpju9t3tp84y",
					Allowance:         allowanceAny,
				}
			},
			shouldErr: true,
		},
		{
			name: "per-app fee grant requests limit reached returns error",
			setup: func() {
				// ----------------------------------
				// --- Save the app
				// ----------------------------------
				err := suite.db.SaveAppSubscription(types.NewApplicationSubscription(
					1,
					"Test App Subscription",
					1,
					10,
					10,
				))
				suite.Require().NoError(err)

				err = suite.db.SaveApp(types.Application{
					ID:                      "1",
					Name:                    "Test Application",
					WalletAddress:           suite.appWallet.AccAddress(),
					SubscriptionID:          1,
					SecretKey:               "secret",
					NotificationsWebhookURL: "https://example.com",
					Admins: []string{
						suite.appWallet.AccAddress(),
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
				// --- Give on-chain authorization
				// ----------------------------------
				suite.grantOnChainExecAuthorization()

				// ----------------------------------
				// --- Save past grants
				// ----------------------------------
				err = suite.db.SaveFeeGrantRequest(types.FeeGrantRequest{
					AppID:         "1",
					DesmosAddress: "desmos1cwfg6eknxz50efv2c0drnpj3dtghxfx905rzke",
					Allowance:     &feegrant.BasicAllowance{},
					RequestTime:   time.Now().Add(-5 * time.Minute),
					GrantTime:     nil,
				})
				suite.Require().NoError(err)
			},
			setupContext: func(ctx context.Context) context.Context {
				return authentication.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *grants.RequestFeeAllowanceRequest {
				allowanceAny, err := codectypes.NewAnyWithValue(&feegrant.BasicAllowance{})
				suite.Require().NoError(err)

				return &grants.RequestFeeAllowanceRequest{
					UserDesmosAddress: "desmos1jvzuh3vzket6hlsdus5kvq7ne3fpju9t3tp84y",
					Allowance:         allowanceAny,
				}
			},
			shouldErr: true,
		},
		{
			name: "user already granted a grant returns error",
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
					ID:                      "1",
					Name:                    "Test Application",
					WalletAddress:           suite.appWallet.AccAddress(),
					SubscriptionID:          1,
					SecretKey:               "secret",
					NotificationsWebhookURL: "https://example.com",
					Admins: []string{
						suite.appWallet.AccAddress(),
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
				// --- Give on-chain authorization
				// ----------------------------------
				suite.grantOnChainExecAuthorization()

				// ----------------------------------
				// --- Save past grants
				// ----------------------------------
				err = suite.db.SaveFeeGrantRequest(types.FeeGrantRequest{
					AppID:         "1",
					DesmosAddress: "desmos1jvzuh3vzket6hlsdus5kvq7ne3fpju9t3tp84y",
					Allowance:     &feegrant.BasicAllowance{},
					RequestTime:   time.Date(2023, 1, 1, 12, 00, 00, 000, time.UTC),
					GrantTime:     testutils.GetTimePointer(time.Date(2023, 1, 1, 12, 10, 00, 000, time.UTC)),
				})
				suite.Require().NoError(err)
			},
			setupContext: func(ctx context.Context) context.Context {
				return authentication.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *grants.RequestFeeAllowanceRequest {
				allowanceAny, err := codectypes.NewAnyWithValue(&feegrant.BasicAllowance{})
				suite.Require().NoError(err)

				return &grants.RequestFeeAllowanceRequest{
					UserDesmosAddress: "desmos1jvzuh3vzket6hlsdus5kvq7ne3fpju9t3tp84y",
					Allowance:         allowanceAny,
				}
			},
			shouldErr: true,
		},
		{
			name: "user with on-chain funds returns error",
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
					ID:                      "1",
					Name:                    "Test Application",
					WalletAddress:           suite.appWallet.AccAddress(),
					SubscriptionID:          1,
					SecretKey:               "secret",
					NotificationsWebhookURL: "https://example.com",
					Admins: []string{
						suite.appWallet.AccAddress(),
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
				// --- Give on-chain authorization
				// ----------------------------------
				suite.grantOnChainExecAuthorization()

				// ----------------------------------
				// --- Send funds
				// ----------------------------------
				suite.sendFundsToUser("desmos12emr8t2yw6y36vlh8qhsfh4m86xrur8m0ymatx")
			},
			setupContext: func(ctx context.Context) context.Context {
				return authentication.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *grants.RequestFeeAllowanceRequest {
				allowanceAny, err := codectypes.NewAnyWithValue(&feegrant.BasicAllowance{})
				suite.Require().NoError(err)

				return &grants.RequestFeeAllowanceRequest{
					UserDesmosAddress: "desmos12emr8t2yw6y36vlh8qhsfh4m86xrur8m0ymatx",
					Allowance:         allowanceAny,
				}
			},
			shouldErr: true,
		},
		{
			name: "user with on-chain grant returns error",
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
					ID:                      "1",
					Name:                    "Test Application",
					WalletAddress:           suite.appWallet.AccAddress(),
					SubscriptionID:          1,
					SecretKey:               "secret",
					NotificationsWebhookURL: "https://example.com",
					Admins: []string{
						suite.appWallet.AccAddress(),
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
				// --- Give on-chain authorization
				// ----------------------------------
				suite.grantOnChainExecAuthorization()

				// ----------------------------------
				// --- Grant allowance
				// ----------------------------------
				suite.sendFundsToUser("desmos1t0takal59h6djq8cceq8vh4x6q6398p6mtvp2n")
			},
			setupContext: func(ctx context.Context) context.Context {
				return authentication.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *grants.RequestFeeAllowanceRequest {
				allowanceAny, err := codectypes.NewAnyWithValue(&feegrant.BasicAllowance{})
				suite.Require().NoError(err)

				return &grants.RequestFeeAllowanceRequest{
					UserDesmosAddress: "desmos1t0takal59h6djq8cceq8vh4x6q6398p6mtvp2n",
					Allowance:         allowanceAny,
				}
			},
			shouldErr: true,
		},
		{
			name: "first time request works properly",
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
					ID:                      "1",
					Name:                    "Test Application",
					WalletAddress:           suite.appWallet.AccAddress(),
					SubscriptionID:          1,
					SecretKey:               "secret",
					NotificationsWebhookURL: "https://example.com",
					Admins: []string{
						suite.appWallet.AccAddress(),
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
				// --- Give on-chain authorization
				// ----------------------------------
				suite.grantOnChainExecAuthorization()
			},
			setupContext: func(ctx context.Context) context.Context {
				return authentication.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *grants.RequestFeeAllowanceRequest {
				allowanceAny, err := codectypes.NewAnyWithValue(&feegrant.BasicAllowance{})
				suite.Require().NoError(err)

				return &grants.RequestFeeAllowanceRequest{
					UserDesmosAddress: "desmos1ehva49d5ltaeszag20wepraf5j3dux2se95f8w",
					Allowance:         allowanceAny,
				}
			},
			shouldErr: false,
			check: func() {
				// Make sure the request is stored in the database
				var count int
				err := suite.db.SQL.Get(&count, `SELECT COUNT(*) FROM fee_grant_requests WHERE grantee_address = $1`, "desmos1ehva49d5ltaeszag20wepraf5j3dux2se95f8w")
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

			_, err := suite.client.RequestFeeAllowance(ctx, tc.buildRequest())

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
