package runner

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/caerus/server/analytics"
	"github.com/desmos-labs/caerus/server/logging"
	applicationsroutes "github.com/desmos-labs/caerus/server/routes/applications"
	filesroutes "github.com/desmos-labs/caerus/server/routes/files"
	grantsroutes "github.com/desmos-labs/caerus/server/routes/grants"
	notificationsroutes "github.com/desmos-labs/caerus/server/routes/notifications"
	userroutes "github.com/desmos-labs/caerus/server/routes/user"
	"github.com/desmos-labs/caerus/server/types"
	"github.com/desmos-labs/caerus/utils"
)

// Runner represents the structure that allows to run everything related to the server
type Runner struct {
	ctx Context

	routesRegisters []RoutesRegister
}

// New returns a new Runner instance
func New(ctx Context) *Runner {
	return &Runner{
		ctx: ctx,
	}

}

func (r *Runner) RegisterDefaultRoutes() {
	r.AddRouteRegister(applicationsroutes.RoutesRegistrar)
	r.AddRouteRegister(filesroutes.RoutesRegistrar)
	r.AddRouteRegister(grantsroutes.RoutesRegistrar)
	r.AddRouteRegister(userroutes.RoutesRegistrar)
	r.AddRouteRegister(notificationsroutes.RoutesRegistrar)
}

func (r *Runner) AddRouteRegister(register RoutesRegister) {
	r.routesRegisters = append(r.routesRegisters, register)
}

func (r *Runner) Run() {
	// Setup the CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")

	// Build the Gin server
	router := gin.New()
	router.Use(logging.ZeroLog(), gin.Recovery(), cors.New(corsConfig))

	// Register the routes
	for _, register := range r.routesRegisters {
		register(router, r.ctx)
	}

	// Build the HTTP server to be able to shut it down if needed
	runningAddress := utils.GetEnvOr(types.EnvAPIsAddress, "0.0.0.0")
	runningPort := utils.GetEnvOr(types.EnvAPIsPort, "3000")
	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", runningAddress, runningPort),
		Handler:           router,
		ReadHeaderTimeout: time.Minute,
		ReadTimeout:       time.Minute,
		WriteTimeout:      time.Minute,
	}

	// Listen for and trap any OS signal to gracefully shutdown and exit
	go r.trapSignal(r.ctx.Scheduler, httpServer)

	// Start the HTTP server
	// Block main process (signal capture will call WaitGroup's Done)
	log.Info().Str("address", httpServer.Addr).Msg("Starting API server")
	err := httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}

	// Start the scheduler
	log.Info().Msg("Starting scheduler")
	r.ctx.Scheduler.StartAsync()
}

// trapSignal traps the stops signals to gracefully shut down the server
func (r *Runner) trapSignal(scheduler *gocron.Scheduler, httpServer *http.Server) {
	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)

	// Kill (no param) default send syscall.SIGTERM
	// Kill -2 is syscall.SIGINT
	// Kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Debug().Msg("Shutting down API server")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("API server forces to shutdown")
	}

	log.Info().Msg("API server shutdown")

	// Stop the scheduler
	scheduler.Stop()
	log.Info().Msg("Scheduler stopped")

	// Perform the cleanup of other things
	analytics.Stop()
}
