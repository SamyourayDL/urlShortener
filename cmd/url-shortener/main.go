package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internals/config"
	"url-shortener/internals/http-server/handlers/redirect"
	deleteHandler "url-shortener/internals/http-server/handlers/url/delete"
	"url-shortener/internals/http-server/handlers/url/save"
	"url-shortener/internals/http-server/logger"
	"url-shortener/internals/lib/logger/handlers/slogpretty"
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

	_, err = storage.DeleteURL("abcdf") //to be deleted
	if err != nil {
		log.Error("abcdf wasn't delete", "err", err)
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer) //if panic in handler - stop panic
	router.Use(middleware.URLFormat) //nice urls in router.Get("/article/{id}"), cann access to params

	router.Post("/url", save.New(log, storage))
	router.Get("/{alias}", redirect.New(log, storage))
	router.Delete("/{alias}", deleteHandler.New(log, storage))

	log.Info("starting server", "address", cfg.Addr)

	srv := http.Server{
		Addr:         cfg.Addr,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", "err", err)
	}

	log.Error("server has stopped")

	//run server
}

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
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

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
