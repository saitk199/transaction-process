package main

import (
	"log/slog"
	"os"
	"transaction-process/internal/config"
)

func main() {

	cfg := config.MustLoad()

	log := setUpLogger()

	log.Info("starting rest service :", slog.String("env", cfg.Env))
	log.Debug("debug logs enable")

}

func setUpLogger() *slog.Logger {
	return slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug}))
}
