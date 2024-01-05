package delete

import (
	resp "github.com/audetv/urlshortener/internal/lib/api/response"
	"github.com/audetv/urlshortener/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Response struct {
	resp.Response
	Alias string `json:"alias"`
}

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Роутер Chi позволяет делать так -
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		err := urlDeleter.DeleteURL(alias)
		if err != nil {
			log.Info("url not found", sl.Err(err))

			render.JSON(w, r, resp.Error("not found"))

			return
		}

		log.Info("URL deleted", slog.String("alias", alias))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
