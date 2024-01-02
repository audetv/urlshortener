package sqlite

import (
	"database/sql"
	"fmt"
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
	CREATE TABLE IF NOT EXIST url(
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
