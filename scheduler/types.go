package scheduler

import (
	"github.com/go-co-op/gocron"
)

// Context contains all the data used in order to build and run a Scheduler instance
type Context struct {
	ChainClient    ChainClient
	FirebaseClient Firebase
	Database       Database
}

type OperationRegistrar func(context Context, scheduler *gocron.Scheduler)
