package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/HAHLIK/AuthService/sso/internal/app/grpc"
)

type App struct {
	GRPCApp *grpcapp.App
}

func New(
	log *slog.Logger,
	gRPCPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	gRPCApp := grpcapp.New(log, gRPCPort)

	return &App{
		GRPCApp: gRPCApp,
	}
}
