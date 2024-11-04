package sqlite

import (
	"ShortRestAPI/internal/storage"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	//"url-shortener/internal/storage"
)

type Storage struct {
	db *sql.DB
}

//func New(storagePath string) (*Storage, error) {
//	const op = "storage.sqlite.New"
//
//	// Открываем соединение с базой данных
//	db, err := sql.Open("sqlite3", storagePath)
//	if err != nil {
//		return nil, fmt.Errorf("%s: %w", op, err)
//	}
//
//	// Подготавливаем SQL-запрос для создания таблицы и индекса
//	stmt, err := db.Prepare(`
//		CREATE TABLE IF NOT EXISTS url(
//			id INTEGER PRIMARY KEY,
//			alias TEXT NOT NULL UNIQUE,
//			url TEXT NOT NULL);
//		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
//	`)
//	if err != nil {
//		// Закрываем базу данных, если что-то пошло не так
//		db.Close()
//		return nil, fmt.Errorf("%s: %w", op, err)
//	}
//	defer stmt.Close() // Закрываем подготовленный оператор после выполнения
//
//	// Выполняем SQL-запрос
//	_, err = stmt.Exec()
//	if err != nil {
//		db.Close()
//		return nil, fmt.Errorf("%s: %w", op, err)
//	}
//
//	// Возвращаем объект хранилища, если все прошло успешно
//	return &Storage{db: db}, nil
//}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
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

func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?,?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: faild to get last insert id %w", op, err)
	}
	return id, nil
}
func (s *Storage) GetUrl(alias string) (string, error) {
	const op = "storage.sqlite.GetUrl"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var resURL string
	err = stmt.QueryRow(alias).Scan(&resURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("%s:execute statement: %w", op, err)
	}
	return resURL, nil

}

// func (s *Storage) DeleteURL(alias string) error
