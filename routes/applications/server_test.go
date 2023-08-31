package applications_test

import (
	"context"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"

	"github.com/desmos-labs/caerus/authentication"
	"github.com/desmos-labs/caerus/database"
	"github.com/desmos-labs/caerus/routes/applications"
	"github.com/desmos-labs/caerus/testutils"
	"github.com/desmos-labs/caerus/types"
)

func TestApplicationsServerTestSuite(t *testing.T) {
	suite.Run(t, new(ApplicationsServerTestSuite))
}

type ApplicationsServerTestSuite struct {
	suite.Suite

	db      *database.Database
	handler *applications.Handler

	server *grpc.Server
	client applications.ApplicationServiceClient
}

func (suite *ApplicationsServerTestSuite) SetupSuite() {
	// Setup the database
	db, err := testutils.CreateDatabase(path.Join("..", "..", "database", "schema"))
	suite.Require().NoError(err)
	suite.db = db

	// Create the handler
	suite.handler = applications.NewHandler(suite.db)

	// Create the server
	suite.server = testutils.CreateServer(suite.db)

	// Register the service
	service := applications.NewServer(suite.handler)
	applications.RegisterApplicationServiceServer(suite.server, service)

	// Setup the client
	conn, err := testutils.StartServerAndConnect(suite.server)
	suite.Require().NoError(err)
	suite.client = applications.NewApplicationServiceClient(conn)
}

func (suite *ApplicationsServerTestSuite) TearDownSuite() {
	suite.server.Stop()
}

func (suite *ApplicationsServerTestSuite) SetupTest() {
	// Truncate all the database data to make sure we have a clean database state
	err := testutils.TruncateDatabase(suite.db)
	suite.Require().NoError(err)
}

// --------------------------------------------------------------------------------------------------------------------

func (suite *ApplicationsServerTestSuite) TestDeleteApp() {
	testCases := []struct {
		name         string
		setup        func()
		setupContext func(ctx context.Context) context.Context
		buildRequest func() *applications.DeleteAppRequest
		shouldErr    bool
		check        func()
	}{
		{
			name: "invalid session returns error",
			buildRequest: func() *applications.DeleteAppRequest {
				return &applications.DeleteAppRequest{
					AppId: "1",
				}
			},
			shouldErr: true,
		},
		{
			name: "user that is not admin returns error",
			setup: func() {
				err := suite.db.SaveSession(types.NewUserSession(
					"desmos1c7ms9zhtgwmv5jy6ztj2vq0jj67zenw3gdl2gr",
					"token",
					time.Now(),
					time.Now().Add(3*time.Hour),
				))
				suite.Require().NoError(err)
			},
			setupContext: func(ctx context.Context) context.Context {
				return authentication.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *applications.DeleteAppRequest {
				return &applications.DeleteAppRequest{
					AppId: "1",
				}
			},
			shouldErr: true,
		},
		{
			name: "valid request returns no error",
			setup: func() {
				// ----------------------------------
				// --- Save the user session
				// ----------------------------------
				err := suite.db.SaveSession(types.NewUserSession(
					"desmos1c7ms9zhtgwmv5jy6ztj2vq0jj67zenw3gdl2gr",
					"token",
					time.Now(),
					time.Now().Add(3*time.Hour),
				))
				suite.Require().NoError(err)

				// ----------------------------------
				// --- Save the app
				// ----------------------------------
				err = suite.db.SaveAppSubscription(types.NewApplicationSubscription(
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
					WalletAddress:           "desmos1ca3pzxx65z7duwxearhxmt8cg93849vn97fmar",
					SubscriptionID:          1,
					SecretKey:               "secret",
					NotificationsWebhookURL: "https://example.com",
					Admins: []string{
						"desmos1c7ms9zhtgwmv5jy6ztj2vq0jj67zenw3gdl2gr",
					},
					CreationTime: time.Now(),
				})
				suite.Require().NoError(err)
			},
			setupContext: func(ctx context.Context) context.Context {
				return authentication.SetupContextWithAuthorization(ctx, "token")
			},
			buildRequest: func() *applications.DeleteAppRequest {
				return &applications.DeleteAppRequest{
					AppId: "1",
				}
			},
			shouldErr: false,
			check: func() {
				// Make sure the app has been deleted
				var count int
				err := suite.db.SQL.Get(&count, "SELECT COUNT(*) FROM applications")
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

			_, err := suite.client.DeleteApp(ctx, tc.buildRequest())

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
