package ports

import "account-cmd/internal/domain"

type EventPublisher interface {
	Publish(e domain.Event) error
}
