package main

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/v5/app"
	"github.com/go-co-op/gocron"

	client "github.com/desmos-labs/caerus/chain"
	"github.com/desmos-labs/caerus/database"
	"github.com/desmos-labs/caerus/firebase"
	applicationsroutes "github.com/desmos-labs/caerus/routes/applications"
	filesroutes "github.com/desmos-labs/caerus/routes/files"
	grantsroutes "github.com/desmos-labs/caerus/routes/grants"
	notificationsroutes "github.com/desmos-labs/caerus/routes/notifications"
	userroutes "github.com/desmos-labs/caerus/routes/user"
	"github.com/desmos-labs/caerus/runner"
)

func main() {
	// Setup Cosmos-related stuff
	app.SetupConfig(sdk.GetConfig())
	encodingCfg := app.MakeEncodingConfig()
	txConfig, cdc, amino := encodingCfg.TxConfig, encodingCfg.Codec, encodingCfg.Amino

	// Build the database
	db, err := database.NewDatabaseFromEnvVariables()
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

	// Build a scheduler
	scheduler := gocron.NewScheduler(time.UTC)

	// Build the runner
	serverRunner := runner.New(runner.Context{
		Codec:          cdc,
		Amino:          amino,
		ChainClient:    chainClient,
		FirebaseClient: firebaseClient,
		Scheduler:      scheduler,
		Database:       db,
	})

	// Register the default routes
	serverRunner.AddRouteRegister(applicationsroutes.RoutesRegistrar)
	serverRunner.AddRouteRegister(filesroutes.RoutesRegistrar)
	serverRunner.AddRouteRegister(grantsroutes.RoutesRegistrar)
	serverRunner.AddRouteRegister(userroutes.RoutesRegistrar)
	serverRunner.AddRouteRegister(notificationsroutes.RoutesRegistrar)

	// Run your server instance
	serverRunner.Run()
}
