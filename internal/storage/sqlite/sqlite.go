package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/audetv/urlshortener/internal/storage"
	"github.com/mattn/go-sqlite3"
)

// Storage структура для объекта Storage
type Storage struct {
	db *sql.DB
}

// New конструктор объекта Storage
func New(storagePath string) (*Storage, error) {
	// Зачем здесь константа op? Я стараюсь всегда добавлять имя текущей функции в возвращаемые ошибки и в логгер,
	// чтобы потом было проще «искать хвосты» в логах.
	// Ведь разные функции часто возвращают одинаковые ошибки и пишут одинаковые логи,
	// а нам обычно нужно понимать, где именно произошло событие.
	const op = "storage.sqlite.NewStorage" // Имя текущей функции для логов и ошибок

	db, err := sql.Open("sqlite3", storagePath) // Подключаемся к БД
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Создаем таблицу если её ещё нет
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
	    id INTEGER PRIMARY KEY,
	    alias TEXT NOT NULL UNIQUE,
	    url TEXT NOT NULL
	);
	CREATE INDEX IF NOT EXIST idx_alias ON url(alias);
`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	// Подготавливаем запрос
	stmt, err := s.db.Prepare("INSERT INTO url(url,alias) values (?,?)")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	// Выполняем запрос
	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}

		return 0, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	// Получаем ID созданной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	// Возвращаем ID
	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var resURL string

	err = stmt.QueryRow(alias).Scan(&resURL)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resURL, nil
}
