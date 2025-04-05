package sqlite

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"transaction-process/internal/storage/domain"
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

	s := &Storage{db: db}

	err = s.createTable()
	if err != nil {
		return nil, err
	}

	err = s.loadSeedData()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Storage) Send(payment domain.Payment) (*domain.Payment, error) {
	return nil, nil
}

func (s *Storage) GetBalance(address string) (*domain.Vallet, error) {
	const op = "storage.sqlite.GetBalance"
	vallet, err := s.db.Prepare("SELECT * FROM vallet WHERE address=?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var valletDto domain.Vallet

	err = vallet.QueryRow(address).Scan(&valletDto)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: Not found vallet with address %s", op, address)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &valletDto, nil
}

func (s *Storage) createTable() error {
	const op = "storage.sqlite.createTable"
	stmt, err := s.db.Prepare(`
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

func (s *Storage) loadSeedData() error {
	const op = "storage.loadSeedData"

	// Определяем путь к файлу ресурсов
	seedFilePath := filepath.Join(".", "storage", "resources", "insertTables.sql")

	// Читаем файл
	data, err := os.ReadFile(seedFilePath)
	if err != nil {
		return fmt.Errorf("%s: не удалось прочитать seed-файл: %w", op, err)
	}

	// Выполняем SQL-скрипт
	_, err = s.db.Exec(string(data))
	if err != nil {
		return fmt.Errorf("%s: ошибка выполнения SQL-скрипта: %w", op, err)
	}

	fmt.Println("✅ Данные успешно загружены в базу!")
	return nil
}
