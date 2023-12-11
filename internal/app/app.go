package app

import (
	grpcapp "auth_service/internal/app/grpc"
	"auth_service/internal/config"
	authService "auth_service/internal/services/auth"
	"auth_service/internal/services/notifier"
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
	kafkaCfg config.KafkaConfig,
) *App {

	storage := mapstorage.New()

	notifier, err := notifier.New(kafkaCfg)
	if err != nil {
		panic("notifier start failed")
	}

	authService := authService.New(
		log,
		storage,
		storage,
		storage,
		tokenTTL,
		notifier,
	)

	grpcApp := grpcapp.New(log, port, authService)

	return &App{
		GRPCSrv: grpcApp,
	}
}
