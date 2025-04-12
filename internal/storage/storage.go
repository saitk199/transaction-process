package storage

import "transaction-process/internal/storage/domain"

type Service interface {
	Send(payment domain.Payment) (*domain.Payment, error)
	GetBalance(address string) (*domain.Vallet, error)
	GetLast(count int) ([]domain.Payment, error)
}
