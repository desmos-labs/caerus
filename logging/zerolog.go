package logging

import (
	"context"
	"os"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/caerus/utils"
)

func init() {
	logLevelStr := utils.GetEnvOr(EnvLoggingLevel, zerolog.InfoLevel.String())
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

// ZeroLogInterceptorLogger adapts zerolog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func ZeroLogInterceptorLogger(l zerolog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		// Build the logger
		logger := l.With().Fields(fields).Logger()

		// Log the message using the trace level
		logger.Trace().Msg(msg)
	})
}
