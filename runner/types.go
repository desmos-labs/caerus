package runner

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/go-co-op/gocron"
	"google.golang.org/grpc"

	"github.com/desmos-labs/caerus/chain"
	"github.com/desmos-labs/caerus/database"
	"github.com/desmos-labs/caerus/firebase"
)

// Context contains all the data used in order to build and run a Runner instance
type Context struct {
	Codec codec.Codec
	Amino *codec.LegacyAmino

	ChainClient    *chain.Client
	FirebaseClient *firebase.Client
	Scheduler      *gocron.Scheduler
	Database       *database.Database
}

type ServiceRegistrar func(context Context, server *grpc.Server)
