package save

import (
	"errors"
	"github.com/audetv/urlshortener/internal/lib/logger/sl"
	"io"

	// для краткости даем короткий алиас пакету
	resp "github.com/audetv/urlshortener/internal/lib/api/response"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias"`
}

// URLSaver хендлер будет сохранять полученные URL-строки,
// поэтому ему нужен Storage, а точнее его метод — SaveURL.
// Опишем соответствующий интерфейс:
type URLSaver interface {
	SaveURL(URL, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		// Добавляем к текущему объекту логгера поля op и request_id
		// Они могут очень упростить нам жизнь в будущем
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Создаём объект запроса и анмаршаллим в него запрос
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			// Такую ошибку встрети, если получили запрос с пустым телом
			// Обработаем её отдельно
			log.Error("request body is empty")

			render.JSON(w, r, resp.Response{
				Status: resp.StatusError,
				Error:  "empty request",
			})

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Response{
				Status: resp.StatusError,
				Error:  "failed to decode request",
			})

			return
		}

		// Запишем ещё один лог. Лучше больше логов
		log.Info("request body decoded", slog.Any("req", req))
	}
}