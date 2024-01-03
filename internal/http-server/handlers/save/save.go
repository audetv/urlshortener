package save

import (
	"errors"
	"github.com/audetv/urlshortener/internal/lib/logger/sl"
	"github.com/audetv/urlshortener/internal/lib/random"
	"github.com/audetv/urlshortener/internal/storage"
	"github.com/go-playground/validator/v10"
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

// TODO: move to config if needed
const aliasLength = 6

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

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		// Запишем ещё один лог. Лучше больше логов
		log.Info("request body decoded", slog.Any("req", req))

		// Создаем объект валидатора
		// и передаём в него структура, которую нужно провалидировать
		if err := validator.New().Struct(req); err != nil {
			// Приводим ошибку к типу ошибки валидации
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		// Alias проверяем вручную. Если он пустой — генерируем случайный:
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		// сохранить URL и Alias, а после — вернуть ответ с сообщением об успехе.
		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			// Отдельно обрабатываем ситуацию,
			// когда запись с таким Alias уже существует
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", id))

		responseOK(w, r, alias)
	}
}
