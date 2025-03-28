package main

import (
	"log/slog"
	"os"
	"transaction-process/internal/config"
	"transaction-process/internal/storage/sqlite"
)

func main() {

	cfg := config.MustLoad()

	log := setUpLogger()

	log.Info("starting rest service :", slog.String("env", cfg.Env))
	log.Debug("debug logs enable")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Ошибка инициализации БД", Err(err))
		os.Exit(1)
	}

	_ = storage

}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func setUpLogger() *slog.Logger {
	return slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug}))
}
