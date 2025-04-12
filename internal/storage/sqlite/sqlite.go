package sqlite

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"time"
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

	const op = "storage.sqlite.Send"
	sender, err := s.GetBalance(payment.Sender)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if sender.Balance < payment.Amount {
		return nil, fmt.Errorf("%s: insufficient funds", op)
	}

	recipient, err := s.GetBalance(payment.Recipient)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//Снимаем средства у отправителя
	_, err = s.db.Exec("UPDATE vallet SET balance = balance - ? WHERE adress = ?", payment.Amount, payment.Sender)
	if err != nil {
		return nil, fmt.Errorf("%s: update sender balance: %w", op, err)
	}

	//Пополняем счет получателя
	if recipient != nil {
		_, err = s.db.Exec("UPDATE vallet SET balance = balance + ? WHERE adress = ?", payment.Amount, payment.Recipient)
	} else {
		_, err = s.db.Exec("INSERT INTO vallet (adress, balance) VALUES (?, ?)", payment.Recipient, payment.Amount)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: update/insert receiver balance: %w", op, err)
	}

	// Вставляем запись о платеже
	paymentDate := time.Now().UTC().Unix()
	paymentId := uuid.New().String()
	_, err = s.db.Exec(`INSERT INTO payment (id, sender, recipient, payment_date, amount) VALUES (?, ?, ?, ?, ?)`,
		paymentId, payment.Sender, payment.Recipient, paymentDate, payment.Amount)
	if err != nil {
		return nil, fmt.Errorf("%s: insert payment: %w", op, err)
	}

	var save domain.Payment

	err = s.db.QueryRow("SELECT id, sender, recipient, payment_date, amount FROM payment WHERE id = ?",
		paymentId).Scan(&save.Id, &save.Sender, &save.Recipient, &save.PaymentDate, &save.Amount)
	if err != nil {
		return nil, fmt.Errorf("%s: payment don't save: %w", op, err)
	}

	return &save, nil
}

func (s *Storage) GetBalance(address string) (*domain.Vallet, error) {
	const op = "storage.sqlite.GetBalance"
	vallet, err := s.db.Prepare("SELECT adress, balance FROM vallet WHERE adress=?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var valletDto domain.Vallet

	err = vallet.QueryRow(address).Scan(&valletDto.Address, &valletDto.Balance)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%s: Not found vallet with address %s", op, address)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &valletDto, nil
}

func (s *Storage) GetLast(count int) ([]domain.Payment, error) {
	const op = "storage.sqlite.GetLast"

	query := `
		SELECT id, sender, recipient, payment_date, amount
		FROM payment
		ORDER BY payment_date DESC
		LIMIT ?
	`

	rows, err := s.db.Query(query, count)
	if err != nil {
		return nil, fmt.Errorf("%s: query failed: %w", op, err)
	}
	defer rows.Close()

	var payments []domain.Payment

	for rows.Next() {
		var payment domain.Payment
		err := rows.Scan(&payment.Id, &payment.Sender, &payment.Recipient, &payment.PaymentDate, &payment.Amount)
		if err != nil {
			return nil, fmt.Errorf("%s: row scan failed: %w", op, err)
		}
		payments = append(payments, payment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	return payments, nil
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
