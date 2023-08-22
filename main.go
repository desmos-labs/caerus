package main

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/v5/app"
	"github.com/go-co-op/gocron"
	"google.golang.org/grpc"

	"github.com/desmos-labs/caerus/branch"
	client "github.com/desmos-labs/caerus/chain"
	"github.com/desmos-labs/caerus/database"
	"github.com/desmos-labs/caerus/firebase"
	"github.com/desmos-labs/caerus/routes/applications"
	"github.com/desmos-labs/caerus/routes/files"
	"github.com/desmos-labs/caerus/routes/grants"
	"github.com/desmos-labs/caerus/routes/links"
	"github.com/desmos-labs/caerus/routes/notifications"
	"github.com/desmos-labs/caerus/routes/users"
	"github.com/desmos-labs/caerus/runner"
	"github.com/desmos-labs/caerus/scheduler"
	"github.com/desmos-labs/caerus/scheduler/operations"
)

func main() {
	// Setup Cosmos-related stuff
	app.SetupConfig(sdk.GetConfig())
	encodingCfg := app.MakeEncodingConfig()
	txConfig, cdc, amino := encodingCfg.TxConfig, encodingCfg.Codec, encodingCfg.Amino

	// Build the database
	db, err := database.NewDatabaseFromEnvVariables(cdc)
	if err != nil {
		panic(err)
	}

	// Build the client
	chainClient, err := client.NewClientFromEnvVariables(txConfig, cdc)
	if err != nil {
		panic(err)
	}

	// Build the various clients
	firebaseClient, err := firebase.NewClientFromEnvVariables()
	if err != nil {
		panic(err)
	}

	branchClient, err := branch.NewClientFromEnvVariables()
	if err != nil {
		panic(err)
	}

	// Build a scheduler - This will be started by the server runner instance
	cronScheduler := scheduler.New(scheduler.Context{
		ChainClient:    chainClient,
		FirebaseClient: firebaseClient,
		Database:       db,
	})

	// Register the scheduler operations
	cronScheduler.SetOperationsRegistrar(func(context scheduler.Context, scheduler *gocron.Scheduler) {
		operations.RegisterGrantsOperations(context, scheduler)
	})

	// Build the runner
	serverRunner := runner.New(runner.Context{
		Codec:          cdc,
		Amino:          amino,
		ChainClient:    chainClient,
		FirebaseClient: firebaseClient,
		BranchClient:   branchClient,
		Scheduler:      cronScheduler,
		Database:       db,
	})

	// Register the default routes
	serverRunner.SetServiceRegistrar(func(context runner.Context, server *grpc.Server) {
		applications.RegisterApplicationServiceServer(server, applications.NewServerFromEnvVariables(context.Database))
		files.RegisterFilesServiceServer(server, files.NewServerFromEnvVariables(context.Database))
		grants.RegisterGrantsServiceServer(server, grants.NewServerFromEnvVariables(context.ChainClient, context.Codec, context.Database))
		notifications.RegisterNotificationsServiceServer(server, notifications.NewServerFromEnvVariables(context.FirebaseClient, context.Database))
		users.RegisterUsersServiceServer(server, users.NewServerFromEnvVariables(context.Codec, context.Amino, context.Database))
		links.RegisterLinksServiceServer(server, links.NewServerFromEnvVariables(context.BranchClient, context.Database))
	})

	// Run the server instance
	serverRunner.Run()
}
