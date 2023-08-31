package notifications_test

import (
	"context"
	"encoding/json"
	"path"
	"testing"
	"time"

	"firebase.google.com/go/v4/messaging"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"

	"github.com/desmos-labs/caerus/authentication"
	"github.com/desmos-labs/caerus/database"
	"github.com/desmos-labs/caerus/routes/notifications"
	notificationstestutils "github.com/desmos-labs/caerus/routes/notifications/testutils"
	"github.com/desmos-labs/caerus/testutils"
	"github.com/desmos-labs/caerus/types"
)

func TestNotificationsServerTestSuite(t *testing.T) {
	suite.Run(t, new(NotificationsServerTestSuite))
}

type NotificationsServerTestSuite struct {
	suite.Suite

	db       *database.Database
	firebase *notificationstestutils.MockFirebase
	handler  *notifications.Handler

	server *grpc.Server
	client notifications.NotificationsServiceClient
}

func (suite *NotificationsServerTestSuite) SetupSuite() {
	// Build the database
	db, err := testutils.CreateDatabase(path.Join("..", "..", "database", "schema"))
	suite.Require().NoError(err)
	suite.db = db

	// Initialize the mocks
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	// Build the mocks
	suite.firebase = notificationstestutils.NewMockFirebase(ctrl)

	// Build the handler
	suite.handler = notifications.NewHandler(suite.firebase, suite.db)

	// Create the server
	suite.server = testutils.CreateServer(suite.db)

	// Register the service
	service := notifications.NewServer(suite.handler)
	notifications.RegisterNotificationsServiceServer(suite.server, service)

	// Setup the client
	conn, err := testutils.StartServerAndConnect(suite.server)
	suite.Require().NoError(err)
	suite.client = notifications.NewNotificationsServiceClient(conn)
}

func (suite *NotificationsServerTestSuite) TearDownSuite() {
	suite.server.Stop()
}

func (suite *NotificationsServerTestSuite) SetupTest() {
	// Truncate all the database data to make sure we have a clean database state
	err := testutils.TruncateDatabase(suite.db)
	suite.Require().NoError(err)
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *NotificationsServerTestSuite) TestSendNotification() {
	testCases := []struct {
		name         string
		setup        func()
		setupContext func(ctx context.Context) context.Context
		buildRequest func() *notifications.SendNotificationRequest
		shouldErr    bool
		check        func()
	}{
		{
			name: "invalid session returns error",
			buildRequest: func() *notifications.SendNotificationRequest {
				notification := types.Notification{
					Notification: &messaging.Notification{},
					Data: map[string]string{
						"type": "test",
					},
				}

				notificationBz, err := json.Marshal(notification)
				suite.Require().NoError(err)

				return &notifications.SendNotificationRequest{
					UserAddresses: []string{"desmos1ncxvkdpq3fl7rlj625s433a7xaraskv0pvmlwp"},
					Notification:  notificationBz,
				}
			},
			shouldErr: true,
		},
		{
			name: "invalid notification data returns error",
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
					WalletAddress:           "desmos1ncxvkdpq3fl7rlj625s433a7xaraskv0pvmlwp",
					SubscriptionID:          1,
					SecretKey:               "secret",
					NotificationsWebhookURL: "https://example.com",
					Admins: []string{
						"desmos1ncxvkdpq3fl7rlj625s433a7xaraskv0pvmlwp",
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
			buildRequest: func() *notifications.SendNotificationRequest {
				return &notifications.SendNotificationRequest{
					UserAddresses: []string{"desmos1ncxvkdpq3fl7rlj625s433a7xaraskv0pvmlwp"},
					Notification:  []byte(`{"key":"value"}`),
				}
			},
			shouldErr: true,
		},
		{
			name: "notifications limit reached returns error",
			setup: func() {
				// ----------------------------------
				// --- Save the app
				// ----------------------------------
				err := suite.db.SaveAppSubscription(types.NewApplicationSubscription(
					1,
					"Test App Subscription",
					10,
					1,
					10,
				))
				suite.Require().NoError(err)

				err = suite.db.SaveApp(types.Application{
					ID:                      "1",
					Name:                    "Test Application",
					WalletAddress:           "desmos1ncxvkdpq3fl7rlj625s433a7xaraskv0pvmlwp",
					SubscriptionID:          1,
					SecretKey:               "secret",
					NotificationsWebhookURL: "https://example.com",
					Admins: []string{
						"desmos1ncxvkdpq3fl7rlj625s433a7xaraskv0pvmlwp",
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
				// --- Store past notifications
				// ----------------------------------
				err = suite.db.SaveSentNotification(&types.SentNotification{
					ID:    "1",
					AppID: "1",
					UserAddresses: []string{
						"desmos182yqg2n4xscw7564r9dphse54ux535300nx2kk",
					},
					Notification: &types.Notification{
						Notification: &messaging.Notification{},
						Data: map[string]string{
							"type": "test",
						},
					},
					SendTime: time.Now().Add(-5 * time.Minute),
				})
			},
			setupContext: func(ctx context.Context) context.Context {
				return authentication.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *notifications.SendNotificationRequest {
				notification := types.Notification{
					Notification: &messaging.Notification{},
					Data: map[string]string{
						"type": "test",
					},
				}

				notificationBz, err := json.Marshal(notification)
				suite.Require().NoError(err)

				return &notifications.SendNotificationRequest{
					UserAddresses: []string{"desmos182yqg2n4xscw7564r9dphse54ux535300nx2kk"},
					Notification:  notificationBz,
				}
			},
			shouldErr: true,
		},
		{
			name: "valid request returns no error - no notification token registered",
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
					WalletAddress:           "desmos1ncxvkdpq3fl7rlj625s433a7xaraskv0pvmlwp",
					SubscriptionID:          1,
					SecretKey:               "secret",
					NotificationsWebhookURL: "https://example.com",
					Admins: []string{
						"desmos1ncxvkdpq3fl7rlj625s433a7xaraskv0pvmlwp",
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
				suite.firebase.EXPECT().
					SendNotificationToUsers(gomock.Any(), gomock.Any(), gomock.Any()).
					AnyTimes().
					Return(nil)
			},
			setupContext: func(ctx context.Context) context.Context {
				return authentication.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *notifications.SendNotificationRequest {
				notification := types.Notification{
					Notification: &messaging.Notification{},
					Data: map[string]string{
						"type": "test",
					},
				}

				notificationBz, err := json.Marshal(notification)
				suite.Require().NoError(err)

				return &notifications.SendNotificationRequest{
					UserAddresses: []string{"desmos182yqg2n4xscw7564r9dphse54ux535300nx2kk"},
					Notification:  notificationBz,
				}
			},
			shouldErr: false,
			check: func() {
				// Expect firebase to not be called
				suite.firebase.EXPECT().
					SendNotificationToUsers(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
			},
		},
		{
			name: "valid request returns no error - existing notification tokens",
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
					WalletAddress:           "desmos1ncxvkdpq3fl7rlj625s433a7xaraskv0pvmlwp",
					SubscriptionID:          1,
					SecretKey:               "secret",
					NotificationsWebhookURL: "https://example.com",
					Admins: []string{
						"desmos1ncxvkdpq3fl7rlj625s433a7xaraskv0pvmlwp",
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
				// --- Register the tokens
				// ----------------------------------
				err = suite.db.SaveUserNotificationDeviceToken(&types.UserNotificationDeviceToken{
					UserAddress: "desmos182yqg2n4xscw7564r9dphse54ux535300nx2kk",
					DeviceToken: "device-token",
				})

				// ----------------------------------
				// --- Setup the mocks
				// ----------------------------------
				suite.firebase.EXPECT().
					SendNotificationToUsers(gomock.Any(), gomock.Any(), gomock.Any()).
					AnyTimes().
					Return(nil)
			},
			setupContext: func(ctx context.Context) context.Context {
				return authentication.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *notifications.SendNotificationRequest {
				notification := types.Notification{
					Notification: &messaging.Notification{},
					Data: map[string]string{
						"type": "test",
					},
				}

				notificationBz, err := json.Marshal(notification)
				suite.Require().NoError(err)

				return &notifications.SendNotificationRequest{
					UserAddresses: []string{"desmos182yqg2n4xscw7564r9dphse54ux535300nx2kk"},
					Notification:  notificationBz,
				}
			},
			shouldErr: false,
			check: func() {
				// Expect firebase to be called once
				suite.firebase.EXPECT().
					SendNotificationToUsers(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1)
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

			_, err := suite.client.SendNotification(ctx, tc.buildRequest())

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
