package app

import (
	grpcapp "auth_service/internal/app/grpc"
	authService "auth_service/internal/services/auth"
	mapstorage "auth_service/internal/storage/map_storage"
	"log/slog"
	"time"

)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	port int,
	storagePath string,
	tokenTTL time.Duration,
) *App {

	storage := mapstorage.New()

	authService := authService.New(
		log,
		storage,
		storage,
		storage,
		tokenTTL,
	)

	grpcApp := grpcapp.New(log, port, authService)

	return &App{
		GRPCSrv: grpcApp,
	}
}
