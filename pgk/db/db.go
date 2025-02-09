package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

func InitDb(path string) (db *sql.DB, error error) {
	var migration bool

	// Проверяем, существует ли файл базы данных
	if _, err := os.Stat(path); os.IsNotExist(err) {
		migration = true
	}

	// Открываем базу данных

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия базы данных: %v", err)
	}

	// Проверяем, что база данных доступна
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %v", err)
	}

	if migration {
		// Создаём таблицу, если её нет
		_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS scheduler (
				id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
				date TEXT NOT NULL,
				title TEXT NOT NULL,
				comment TEXT,
				repeat TEXT CHECK (LENGTH(repeat) <= 128)
			);
		`)
		if err != nil {
			return nil, fmt.Errorf("ошибка создания таблицы: %v", err)
		}

		// Создаём индекс для ускорения запросов
		_, err = db.Exec(`
			CREATE INDEX IF NOT EXISTS scheduler_date ON scheduler(date);
		`)
		if err != nil {
			return nil, fmt.Errorf("ошибка создания индекса: %v", err)
		}
	}

	return db, nil
}
