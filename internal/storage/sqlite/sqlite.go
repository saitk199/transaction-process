package sqlite

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
)

type Storage struct {
	db *sql.DB
}

func InitDataBase(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.InitDataBase"

	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = createTable(db)
	if err != nil {
		return nil, err
	}

	err = loadSeedData(db)
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func createTable(db *sql.DB) error {
	const op = "storage.sqlite.createTable"
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS vallet(
    	adress TEXT PRIMARY KEY,
    	balance REAL NOT NULL
	);

	CREATE TABLE IF NOT EXISTS payment(
    	id TEXT PRIMARY KEY,
    	sender TEXT NOT NULL,
    	recipient TEXT NOT NULL,
    	payment_date INTEGER,
	    amount REAL NOT NULL,
    	FOREIGN KEY(sender) REFERENCES vallet(adress),
    	FOREIGN KEY(recipient) REFERENCES vallet(adress)
	);

	CREATE INDEX IF NOT EXISTS idx_adress ON vallet(adress);
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func loadSeedData(db *sql.DB) error {
	const op = "storage.loadSeedData"

	// Определяем путь к файлу ресурсов
	seedFilePath := filepath.Join(".", "storage", "resources", "insertTables.sql")

	// Читаем файл
	data, err := os.ReadFile(seedFilePath)
	if err != nil {
		return fmt.Errorf("%s: не удалось прочитать seed-файл: %w", op, err)
	}

	// Выполняем SQL-скрипт
	_, err = db.Exec(string(data))
	if err != nil {
		return fmt.Errorf("%s: ошибка выполнения SQL-скрипта: %w", op, err)
	}

	fmt.Println("✅ Данные успешно загружены в базу!")
	return nil
}
