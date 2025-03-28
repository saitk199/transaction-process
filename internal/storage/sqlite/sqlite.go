package sqlite

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: $w", op, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS vallet(
		adress TEXT PRIMARY KEY,
		balance REAL NOT NULL
		);
	CREATE TABLE IF NOT EXISTS transaction(
		id TEXT PRIMARY KEY,
		from TEXT NOT NULL,
		FOREIGN KEY(from) REFERENCES vallet(adress),
		to TEXT NOT NULL,
		FOREIGN KEY(to) REFERENCES vallet(adress),
		date INTEGER
	CREATE INDEX IF NOT EXISTS idx_adress ON vallet(adress);
	)
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: $w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: $w", op, err)
	}

	return &Storage{db: db}, nil
}
