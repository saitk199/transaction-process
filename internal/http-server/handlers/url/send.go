package url

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"transaction-process/internal/lib/api"
	"transaction-process/internal/storage"
	"transaction-process/internal/storage/domain"
)

type Response struct {
	api.Response
	Data interface{} `json:"data,omitempty"`
}

func Send(log *slog.Logger, service storage.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.getBalance.Send"
		w.Header().Set("Content-Type", "application/json")

		log.With(
			slog.String("operation", op),
			slog.String("requestId", middleware.GetReqID(r.Context())))

		var payment domain.Payment

		err := render.DecodeJSON(r.Body, &payment)

		if errors.Is(err, io.EOF) {
			// Такую ошибку встретим, если получили запрос с пустым телом.
			// Обработаем её отдельно
			log.Error("request body is empty: " + err.Error())

			render.JSON(w, r, api.Error("empty request"))

			return
		}

		if err != nil {
			log.Error("Failed to decode request body: " + err.Error())

			render.JSON(w, r, api.Error("Failed to decode request body"))

			return
		}

		log.Info("request body decoded", slog.Any("request", payment))

		result, err := service.Send(payment)

		if err != nil {
			log.Error("Failed to send payment: " + err.Error())

			render.JSON(w, r, api.Error("Failed to send payment"))

			return
		}

		log.Info("Successfully sent payment", slog.Any("payment", result))

		render.JSON(w, r, Response{
			Response: api.OK(),
			Data:     result,
		})
	}
}

func GetBalance(log *slog.Logger, service storage.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.GetBalance"
		w.Header().Set("Content-Type", "application/json")

		log.With(
			slog.String("operation", op),
			slog.String("requestId", middleware.GetReqID(r.Context())))

		address := chi.URLParam(r, "address")

		result, err := service.GetBalance(address)
		if err != nil {
			log.Error("Failed to get balance: " + err.Error())

			render.JSON(w, r, api.Error("Failed to get balance"))

			return
		}

		log.Info("Successfully get balance", slog.Any("vallet", result))

		render.JSON(w, r, Response{
			Response: api.OK(),
			Data:     result,
		})
	}
}

func GetLast(log *slog.Logger, service storage.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.GetLast"
		w.Header().Set("Content-Type", "application/json")

		log.With(
			slog.String("operation", op),
			slog.String("requestId", middleware.GetReqID(r.Context())))

		countStr := r.URL.Query().Get("count")
		count := 10 // значение по умолчанию
		if countStr != "" {
			if parsed, err := strconv.Atoi(countStr); err == nil {
				count = parsed
			} else {
				log.Error("Failed to parsed count: " + err.Error())
				render.JSON(w, r, api.Error("Failed to parsed count"))
				return
			}
		}
		result, err := service.GetLast(count)
		if err != nil {
			log.Error("Failed to get last payment count: " + err.Error())
			render.JSON(w, r, api.Error("Failed to get last payment count"))
			return
		}

		log.Info("Successfully get last", slog.Any("payment", result))

		render.JSON(w, r, Response{
			Response: api.OK(),
			Data:     result,
		})
	}
}
