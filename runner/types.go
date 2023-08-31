package runner

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"google.golang.org/grpc"

	"github.com/desmos-labs/caerus/branch"
	"github.com/desmos-labs/caerus/chain"
	"github.com/desmos-labs/caerus/database"
	"github.com/desmos-labs/caerus/notifications"
	"github.com/desmos-labs/caerus/scheduler"
)

// Context contains all the data used in order to build and run a Runner instance
type Context struct {
	Codec codec.Codec
	Amino *codec.LegacyAmino

	ChainClient         *chain.Client
	NotificationsClient *notifications.Client
	BranchClient        *branch.Client
	Scheduler           *scheduler.Scheduler
	Database            *database.Database
}

type ServiceRegistrar func(context Context, server *grpc.Server)
