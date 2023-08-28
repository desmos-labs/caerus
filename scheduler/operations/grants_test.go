package operations_test

import (
	"path"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	desmosapp "github.com/desmos-labs/desmos/v6/app"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/desmos-labs/caerus/database"
	"github.com/desmos-labs/caerus/scheduler"
	"github.com/desmos-labs/caerus/scheduler/operations"
	schedulertestutils "github.com/desmos-labs/caerus/scheduler/testutils"
	"github.com/desmos-labs/caerus/testutils"
	"github.com/desmos-labs/caerus/types"
)

func TestGrantsOperationsTestSuite(t *testing.T) {
	suite.Run(t, new(GrantsOperationsTestSuite))
}

type GrantsOperationsTestSuite struct {
	suite.Suite

	chainClient *schedulertestutils.MockChainClient
	firebase    *schedulertestutils.MockFirebase
	db          *database.Database

	ctx scheduler.Context
}

func (suite *GrantsOperationsTestSuite) SetupSuite() {
	// Setup the Cosmos SDK config
	desmosapp.SetupConfig(sdk.GetConfig())

	// Build the database
	db, err := testutils.CreateDatabase(path.Join("..", "..", "database", "schema"))
	suite.Require().NoError(err)
	suite.db = db
}

func (suite *GrantsOperationsTestSuite) SetupTest() {
	// Setup the mocks
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	// Build the mocks
	suite.chainClient = schedulertestutils.NewMockChainClient(ctrl)
	suite.firebase = schedulertestutils.NewMockFirebase(ctrl)

	// Setup the context
	suite.ctx = scheduler.Context{
		ChainClient:    suite.chainClient,
		FirebaseClient: suite.firebase,
		Database:       suite.db,
	}

	// Truncate all the database data to make sure we have a clean database state
	err := testutils.TruncateDatabase(suite.db)
	suite.Require().NoError(err)
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *GrantsOperationsTestSuite) TestGrantAuthorizations() {
	testCases := []struct {
		name      string
		setup     func()
		shouldErr bool
		check     func()
	}{
		{
			name:      "not found requests return no error",
			shouldErr: false,
			check: func() {
				// Make sure the clients have not been used
				suite.chainClient.EXPECT().
					BroadcastTxCommit(gomock.Any()).
					Times(0)

				suite.firebase.EXPECT().
					SendNotifications(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0)
			},
		},
		{
			name: "application without on-chain authorization is notified",
			setup: func() {
				// ----------------------------------
				// --- Save the app
				// ----------------------------------
				err := suite.db.SaveAppSubscription(types.ApplicationSubscription{
					ID:                     1,
					Name:                   "Test App Subscription",
					FeeGrantLimit:          10,
					NotificationsRateLimit: 10,
				})
				suite.Require().NoError(err)

				err = suite.db.SaveApp(types.Application{
					ID:             "1",
					Name:           "Test Application",
					WalletAddress:  "desmos1ncxvkdpq3fl7rlj625s433a7xaraskv0pvmlwp",
					SubscriptionID: 1,
					Admins: []string{
						"desmos1ncxvkdpq3fl7rlj625s433a7xaraskv0pvmlwp",
					},
					CreationTime: time.Now(),
				})
				suite.Require().NoError(err)

				err = suite.db.SaveAppNotificationDeviceToken(&types.AppNotificationDeviceToken{
					AppID:       "1",
					DeviceToken: "DeviceToken",
				})
				suite.Require().NoError(err)

				// ----------------------------------
				// --- Save the requests
				// ----------------------------------
				err = suite.db.SaveFeeGrantRequest(types.FeeGrantRequest{
					ID:            "1",
					AppID:         "1",
					DesmosAddress: "desmos14p3u4zvqx7rupxlc9ja6qzcwucdvxv0cwk8pv5",
					Allowance:     &feegrant.BasicAllowance{},
					RequestTime:   time.Now().Add(-5 * time.Minute),
					GrantTime:     nil,
				})
				suite.Require().NoError(err)

				err = suite.db.SaveFeeGrantRequest(types.FeeGrantRequest{
					ID:            "2",
					AppID:         "1",
					DesmosAddress: "desmos1jdxfc2wsrdu8lyj9w4nfjarsna3faeghwpenxv",
					Allowance:     &feegrant.BasicAllowance{},
					RequestTime:   time.Now().Add(-5 * time.Minute),
					GrantTime:     nil,
				})
				suite.Require().NoError(err)

				// ----------------------------------
				// --- Setup the mocks
				// ----------------------------------
				suite.chainClient.EXPECT().
					HasGrantedMsgGrantAllowanceAuthorization(gomock.Any()).
					Return(false, nil)

				suite.firebase.EXPECT().
					SendNotifications(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			shouldErr: false,
			check: func() {
				// Make sure the clients have been used
				suite.chainClient.EXPECT().
					HasGrantedMsgGrantAllowanceAuthorization(gomock.Any()).
					Times(1)

				suite.firebase.EXPECT().
					SendNotifications(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				suite.chainClient.EXPECT().
					BroadcastTxCommit(gomock.Any()).
					Times(0)

				// Make sure the grant requests are not set as granted
				var count int
				err := suite.db.SQL.Get(&count, `SELECT COUNT(*) FROM fee_grant_requests WHERE grant_time IS NULL`)
				suite.Require().NoError(err)
				suite.Require().Equal(2, count)
			},
		},
		{
			name: "valid requests are handled properly",
			setup: func() {
				// ----------------------------------
				// --- Save the app
				// ----------------------------------
				err := suite.db.SaveAppSubscription(types.ApplicationSubscription{
					ID:                     1,
					Name:                   "Test App Subscription",
					FeeGrantLimit:          10,
					NotificationsRateLimit: 10,
				})
				suite.Require().NoError(err)

				err = suite.db.SaveApp(types.Application{
					ID:             "1",
					Name:           "Test Application",
					WalletAddress:  "desmos1ncxvkdpq3fl7rlj625s433a7xaraskv0pvmlwp",
					SubscriptionID: 1,
					Admins: []string{
						"desmos1ncxvkdpq3fl7rlj625s433a7xaraskv0pvmlwp",
					},
					CreationTime: time.Now(),
				})
				suite.Require().NoError(err)

				err = suite.db.SaveAppNotificationDeviceToken(&types.AppNotificationDeviceToken{
					AppID:       "1",
					DeviceToken: "DeviceToken",
				})
				suite.Require().NoError(err)

				// ----------------------------------
				// --- Save the requests
				// ----------------------------------
				err = suite.db.SaveFeeGrantRequest(types.FeeGrantRequest{
					ID:            "1",
					AppID:         "1",
					DesmosAddress: "desmos14p3u4zvqx7rupxlc9ja6qzcwucdvxv0cwk8pv5",
					Allowance:     &feegrant.BasicAllowance{},
					RequestTime:   time.Now().Add(-5 * time.Minute),
					GrantTime:     nil,
				})
				suite.Require().NoError(err)

				err = suite.db.SaveFeeGrantRequest(types.FeeGrantRequest{
					ID:            "2",
					AppID:         "1",
					DesmosAddress: "desmos1jdxfc2wsrdu8lyj9w4nfjarsna3faeghwpenxv",
					Allowance:     &feegrant.BasicAllowance{},
					RequestTime:   time.Now().Add(-5 * time.Minute),
					GrantTime:     nil,
				})
				suite.Require().NoError(err)

				// ----------------------------------
				// --- Setup the mocks
				// ----------------------------------
				suite.chainClient.EXPECT().
					AccAddress().
					Return("desmos1sfklnhd5fu5jtgtmxpdm3dsg2l895rl95j8zvn")

				suite.chainClient.EXPECT().
					HasGrantedMsgGrantAllowanceAuthorization(gomock.Any()).
					Return(true, nil)

				suite.chainClient.EXPECT().
					BroadcastTxCommit(gomock.Any()).AnyTimes().
					Return(&sdk.TxResponse{Code: 0}, nil)

				suite.firebase.EXPECT().
					SendNotifications(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
			shouldErr: false,
			check: func() {
				// Make sure the clients have been used
				suite.chainClient.EXPECT().
					HasGrantedMsgGrantAllowanceAuthorization(gomock.Any()).
					Times(1)

				suite.chainClient.EXPECT().
					BroadcastTxCommit(gomock.Any()).
					Times(1)

				suite.firebase.EXPECT().
					SendNotifications(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)

				// Make sure the grant requests are set as granted
				var count int
				err := suite.db.SQL.Get(&count, `SELECT COUNT(*) FROM fee_grant_requests WHERE grant_time IS NOT NULL`)
				suite.Require().NoError(err)
				suite.Require().Equal(2, count)
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

			err := operations.GrantAuthorizations(suite.ctx)

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
