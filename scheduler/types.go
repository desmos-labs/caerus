package scheduler

import (
	"github.com/go-co-op/gocron"

	"github.com/desmos-labs/caerus/chain"
	"github.com/desmos-labs/caerus/database"
	"github.com/desmos-labs/caerus/firebase"
)

// Context contains all the data used in order to build and run a Scheduler instance
type Context struct {
	ChainClient    *chain.Client
	FirebaseClient *firebase.Client
	Database       *database.Database
}

type OperationRegistrar func(context Context, scheduler *gocron.Scheduler)
