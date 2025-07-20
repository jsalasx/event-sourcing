package domain

import "time"

// Interfaz base
type Event interface {
	EventType() string
	AggregateID() string
	Timestamp() time.Time
}

// Estructura base para todos los eventos
type BaseEvent struct {
	ID          string
	AggregateId string
	Type        string
	OccurredAt  time.Time
}

func (e BaseEvent) EventType() string    { return e.Type }
func (e BaseEvent) AggregateID() string  { return e.AggregateId }
func (e BaseEvent) Timestamp() time.Time { return e.OccurredAt }

// Eventos del dominio
type AccountCreated struct {
	BaseEvent
	OwnerName      string
	InitialBalance float64
}

type MoneyDeposited struct {
	BaseEvent
	Amount float64
}

type MoneyWithdrawn struct {
	BaseEvent
	Amount float64
}
