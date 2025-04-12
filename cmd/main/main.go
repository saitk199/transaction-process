package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"transaction-process/internal/config"
	"transaction-process/internal/http-server/handlers/url"
	"transaction-process/internal/storage/sqlite"
)

func main() {

	cfg := config.MustLoad()

	log := setUpLogger()

	log.Info("starting rest service :", slog.String("env", cfg.Env))
	log.Debug("debug logs enable")

	storage, err := sqlite.InitDataBase(cfg.StoragePath)
	if err != nil {
		log.Error("Ошибка инициализации БД", Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/api/send", url.Send(log, storage))
	router.Get("/api/wallet/{address}/balance", url.GetBalance(log, storage))
	router.Get("/api/transactions", url.GetLast(log, storage))

	log.Info("starting http server :", slog.String("env", cfg.Env))

	srv := &http.Server{
		Addr:              cfg.Adress,
		Handler:           router,
		ReadHeaderTimeout: cfg.HTTPServer.Timeout,
		WriteTimeout:      cfg.HTTPServer.Timeout,
		IdleTimeout:       cfg.HTTPServer.IDDLETimeout,
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("failed to start http server :", slog.String("err", err.Error()))
	}

	log.Info("finished http server :", slog.String("env", cfg.Env))
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
