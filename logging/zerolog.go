package logging

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/caerus/types"
	"github.com/desmos-labs/caerus/utils"
)

func init() {
	logLevelStr := utils.GetEnvOr(types.EnvLoggingLevel, zerolog.TraceLevel.String())
	logLevel, err := zerolog.ParseLevel(logLevelStr)
	if err != nil {
		panic(err)
	}

	// Setup logging
	log.Logger = zerolog.
		New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(logLevel).
		With().Timestamp().
		Logger()
}

// ZeroLog returns a Gin Handler function that logs endpoint calls using ZeroLog
func ZeroLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Log request
		log.Trace().Str("path", c.Request.URL.Path).Msg("received request")

		// Log errors
		for _, err := range c.Errors {
			log.Error().
				Int("status", c.Writer.Status()).
				Str("path", c.Request.URL.Path).
				Msg(err.Error())
		}
	}
}
