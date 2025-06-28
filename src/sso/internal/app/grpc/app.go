package grpcapp

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	authgrpc "github.com/HAHLIK/AuthService/sso/internal/grpc/auth"
	"github.com/HAHLIK/AuthService/sso/internal/services/auth"
	"github.com/HAHLIK/AuthService/sso/internal/storage/sqlite"
	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int, tokenTTL time.Duration) *App {
	storage := sqlite.New()
	auth := auth.New(log, storage, storage, storage, tokenTTL)

	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, auth)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	err := a.run()
	if err != nil {
		panic(err)
	}
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(
		"op", op,
		"port", a.port,
	).Info("gRPC server is stopping")

	a.gRPCServer.GracefulStop()
}

func (a *App) run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		"op", op,
		"port", a.port,
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is running", slog.String("addr", lis.Addr().String()))

	if err := a.gRPCServer.Serve(lis); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
