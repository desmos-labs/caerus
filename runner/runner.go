package runner

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/reflection"

	"github.com/desmos-labs/caerus/analytics"
	"github.com/desmos-labs/caerus/scheduler"
	"github.com/desmos-labs/caerus/server"
	"github.com/desmos-labs/caerus/utils"
)

// Runner represents the structure that allows to run everything related to the server
type Runner struct {
	ctx              Context
	serviceRegistrar ServiceRegistrar
}

// New returns a new Runner instance
func New(ctx Context) *Runner {
	return &Runner{
		ctx: ctx,
	}
}

// SetServiceRegistrar sets the given registrar as the service registrar
func (r *Runner) SetServiceRegistrar(registrar ServiceRegistrar) {
	r.serviceRegistrar = registrar
}

// Run runs the server blocking the current goroutine until the server is shut down
func (r *Runner) Run() {
	serverInstance := server.New(r.ctx.Database)

	if r.serviceRegistrar != nil {
		r.serviceRegistrar(r.ctx, serverInstance)
	}

	// Register the reflection service to allow listing methods provided by the server
	reflection.Register(serverInstance)

	// Build the HTTP server to be able to shut it down if needed
	runningAddress := utils.GetEnvOr(EnvServerAddress, "0.0.0.0")
	runningPort := utils.GetEnvOr(EnvServerPort, "3000")
	netListener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", runningAddress, runningPort))
	if err != nil {
		panic(err)
	}

	// Listen for and trap any OS signal to gracefully shutdown and exit
	go r.trapSignal(r.ctx.Scheduler, netListener)

	// Start the scheduler
	log.Info().Msg("Starting scheduler")
	r.ctx.Scheduler.StartAsync()

	// Start the serer (as this is blocking)
	log.Info().Str("address", runningAddress).Str("port", runningPort).Msg("Starting API server")
	serverInstance.Serve(netListener)
}

// trapSignal traps the stops signals to gracefully shut down the server
func (r *Runner) trapSignal(scheduler *scheduler.Scheduler, netListener net.Listener) {
	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)

	// Kill (no param) default send syscall.SIGTERM
	// Kill -2 is syscall.SIGINT
	// Kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Debug().Msg("Shutting down API server")

	if err := netListener.Close(); err != nil {
		log.Error().Err(err).Msg("API server forces to shutdown")
	}

	log.Info().Msg("API server shutdown")

	// Stop the scheduler
	scheduler.Stop()
	log.Info().Msg("Scheduler stopped")

	// Perform the cleanup of other things
	analytics.Stop()
}
