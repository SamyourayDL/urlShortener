package main

import (
	"fmt"
	"log/slog"
	"url-shortener/internals/config"
)

const ()

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	//init logger

	//init storage

	//init router

	//run server
}

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case "local":

	}

	return log
}
