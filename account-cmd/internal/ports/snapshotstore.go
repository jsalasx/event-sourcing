package ports

import "account-cmd/internal/domain"

// EventStore para persistencia de eventos
type SnapshotStore interface {
	UpsertSnapshot(snapshot domain.BankAccountSnapshot) error
	LoadSnapshot(aggregateID string) (*domain.BankAccountSnapshot, error)
}
