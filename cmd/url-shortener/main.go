package main

import (
	"fmt"
	"log/slog"
	"os"
	"url-shortener/internals/config"
	"url-shortener/internals/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	log := SetupLogger(cfg.Env)
	log.Info("Logger succesfully setuped", "Env", cfg.Env)
	log.Debug("Debug messages are enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to init Storage", "err", err)
		os.Exit(1)
	}

	_ = storage //to be deleted

	router := chi.NewRouter()

	router.Use(middleware.RequestID)

	//run server
}

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New( //LevelInfo, don't want to see debug messages on server
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
