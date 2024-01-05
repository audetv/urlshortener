package main

import (
	"github.com/audetv/urlshortener/internal/config"
	"github.com/audetv/urlshortener/internal/http-server/handlers/redirect"
	"github.com/audetv/urlshortener/internal/http-server/handlers/url/delete"
	"github.com/audetv/urlshortener/internal/http-server/handlers/url/save"
	mwLogger "github.com/audetv/urlshortener/internal/http-server/middleware/logger"
	"github.com/audetv/urlshortener/internal/http-server/server"
	"github.com/audetv/urlshortener/internal/lib/logger/handlers/slogpretty"
	"github.com/audetv/urlshortener/internal/lib/logger/sl"
	"github.com/audetv/urlshortener/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log = log.With(slog.String("env", cfg.Env)) // к каждому сообщению будет добавляться поле с информацией о текущем окружении

	log.Info("initializing server", slog.String("address", cfg.Address))
	log.Debug("logger debug mode enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
		os.Exit(1)
	}

	// создадим объект роутера и подключим к нему необходимый middleware
	router := chi.NewRouter()

	router.Use(middleware.RequestID) // Добавляет request_id в каждый запрос, для трейсинга
	//router.Use(middleware.Logger)    // Логирование всех запросов
	router.Use(mwLogger.New(log))    // Логирование всех запросов
	router.Use(middleware.Recoverer) // Если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
	router.Use(middleware.URLFormat) // Парсер URLов поступающих запросов

	// Все пути этого роутера будут начинаться с префикса `/url`
	router.Route("/url", func(r chi.Router) {
		// Подключаем middleware BasicAuth авторизацию
		r.Use(middleware.BasicAuth("urlshortener", map[string]string{
			// Передаём в BasicAuth логин и пароль из конфига
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		r.Post("/", save.New(log, storage))
		r.Delete("/{alias}", delete.New(log, storage))
	})

	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	// Создаем объект сервера
	srv := server.New(cfg, router)

	// Запускаем сервер
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
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
