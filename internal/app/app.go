package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/iriskin77/grpc_go/internal/app/grpc"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	// Todo инициализировать хранилище (storage)

	// Todo: init auth service (auth)

	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
