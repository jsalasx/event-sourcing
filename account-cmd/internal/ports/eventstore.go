package ports

import "account-cmd/internal/domain"

// EventStore para persistencia de eventos
type EventStore interface {
	Save(aggregateID string, events []domain.Event) error
	Load(aggregateID string) ([]domain.Event, error)
	LoadWithVersion(aggregateID string, version uint64) ([]domain.Event, error)
}
