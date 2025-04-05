package url

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
	"transaction-process/internal/lib/api"
	"transaction-process/internal/storage"
	"transaction-process/internal/storage/domain"
)

type Response struct {
	api.Response
	Data interface{} `json:"data,omitempty"`
}

func New(log *slog.Logger, service storage.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.getBalance.New"

		log.With(
			slog.String("operation", op),
			slog.String("requestId", middleware.GetReqID(r.Context())))

		var payment domain.Payment

		err := render.DecodeJSON(r.Body, &payment)

		if errors.Is(err, io.EOF) {
			// Такую ошибку встретим, если получили запрос с пустым телом.
			// Обработаем её отдельно
			log.Error("request body is empty")

			render.JSON(w, r, api.Error("empty request"))

			return
		}

		if err != nil {
			log.Error("Failed to decode request body")

			render.JSON(w, r, api.Error("Failed to decode request body"))

			return
		}

		log.Info("request body decoded", slog.Any("request", payment))

		result, err := service.Send(payment)

		if err != nil {
			log.Error("Failed to send payment")

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
