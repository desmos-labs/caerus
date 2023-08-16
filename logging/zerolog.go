package logging

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

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

// ZeroLogInterceptorLogger adapts zerolog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func ZeroLogInterceptorLogger(l zerolog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		logger := l.With().Fields(fields).Logger()

		switch lvl {
		case logging.LevelDebug:
			logger.Debug().Msg(msg)
		case logging.LevelInfo:
			logger.Info().Msg(msg)
		case logging.LevelWarn:
			logger.Warn().Msg(msg)
		case logging.LevelError:
			logger.Error().Msg(msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
	})
}

// NewLoggingInterceptor returns a list of gprc.ServerOptions that can be used to enable logging
// on the gRPC server
func NewLoggingInterceptor() []grpc.ServerOption {
	logFunc := ZeroLogInterceptorLogger(log.Logger)
	return []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(logging.UnaryServerInterceptor(logFunc)),
		grpc.ChainStreamInterceptor(logging.StreamServerInterceptor(logFunc)),
	}
}
