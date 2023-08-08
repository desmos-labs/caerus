package runner

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"

	"github.com/desmos-labs/caerus/server/chain"
	"github.com/desmos-labs/caerus/server/database"
	"github.com/desmos-labs/caerus/server/firebase"
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

// RoutesRegister represents a function that allows to register a series of routes into the given engince
type RoutesRegister func(engine *gin.Engine, context Context)
