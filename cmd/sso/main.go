package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/iriskin77/grpc_go/internal/app"
	"github.com/iriskin77/grpc_go/internal/config"
)

func main() {

	// TODO: инициализировать объект конфига
	cfg := config.MustLoad()

	// TODO: инициализировать логгер

	log := setupLogger(cfg.Env)

	log.Info("starting application", slog.String("env", cfg.Env), slog.Any("cfg", cfg))

	// Инициализация приложения (app)

	// Запустить gRPC сервер приложения

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPCSrv.MustRun()

	// Graceful shutdown

	stop := make(chan os.Signal, 1)
	// Notify ждет, когда придет один из перечисленных сигналов от операционной системы
	// И когда это произойдет, то Notify запишет что-то в наш канал
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	// В этой точке мы будем висеть, пока не придет сигнал от
	// ОС о завершении. И пока мы ждем, горутина go application
	// будет обрабатывать запросы, но как только что-то придет в канал stop,
	// Эта строчка будет выполнена.
	sign := <-stop

	application.GRPCSrv.Stop()

	log.Info("application stop", sign)

}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case envDev:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log

}
