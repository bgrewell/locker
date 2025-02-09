package main

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	locker "locker/api/go"
	"locker/internal/config"
	"locker/pkg/server"
	"net"
)

// NewGRPCServer creates a gRPC server, registers the LockerService, and sets up lifecycle hooks.
func NewGRPCServer(lc fx.Lifecycle, cfg *config.Configuration, log *zap.Logger, svc locker.LockerServiceServer) *grpc.Server {
	server := grpc.NewServer()
	// Register your service implementation with the gRPC server.
	locker.RegisterLockerServiceServer(server, svc)

	// Open a TCP listener on the address specified in your configuration.
	lis, err := net.Listen("tcp", cfg.GRPCAddress)
	if err != nil {
		log.Fatal("failed to listen", zap.Error(err))
	}

	// Add lifecycle hooks to start and stop the gRPC server.
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting gRPC server", zap.String("address", cfg.GRPCAddress))
			go func() {
				if err := server.Serve(lis); err != nil {
					log.Error("gRPC server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping gRPC server")
			server.GracefulStop()
			return nil
		},
	})
	return server
}

// NewLogger returns a new instance of the zap.Logger
func NewLogger() (*zap.Logger, error) {
	zapCFG := zap.NewDevelopmentConfig()
	zapCFG.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapCFG.DisableStacktrace = true
	return zapCFG.Build()
}

// LoadConfig returns the current configuration
func LoadConfig() (*config.Configuration, error) {
	return config.ReadConfiguration()
}

func main() {

	fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Provide(
			NewLogger,
			LoadConfig,
			server.NewLockerServiceServer,
			NewGRPCServer,
		),
		fx.Invoke(func(s *grpc.Server) {}),
	).Run()

}
