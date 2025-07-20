package domain

import (
	"errors"
	"time"

	"shared/utils"

	"github.com/google/uuid"
)

type Event interface {
	EventType() string
	AggregateID() string
	Timestamp() time.Time
	GetVersion() uint64
	SetVersion(uint64)
}

// Base de todos los eventos
type BaseEvent struct {
	ID          string
	AggregateId string
	Type        string
	OccurredAt  time.Time
	Version     uint64
}

func (e BaseEvent) EventType() string    { return e.Type }
func (e BaseEvent) AggregateID() string  { return e.AggregateId }
func (e BaseEvent) Timestamp() time.Time { return e.OccurredAt }
func (e *BaseEvent) GetVersion() uint64  { return e.Version }
func (e *BaseEvent) SetVersion(v uint64) {
	e.Version = v

}

// Eventos concretos
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

// Agregado Cuenta Bancaria
type BankAccount struct {
	ID      string
	Owner   string
	Balance float64
	Version uint64
	Events  []Event
}

func NewAccount(owner string, initialBalance float64) *BankAccount {
	id := uuid.NewString()
	acc := &BankAccount{ID: id, Version: 0}
	acc.apply(&AccountCreated{
		BaseEvent: BaseEvent{
			ID:          uuid.NewString(),
			AggregateId: id,
			Type:        "AccountCreated",
			OccurredAt:  time.Now(),
			Version:     1,
		},
		OwnerName:      owner,
		InitialBalance: initialBalance,
	})
	return acc
}

func (b *BankAccount) Deposit(amount float64) {
	b.apply(&MoneyDeposited{
		BaseEvent: BaseEvent{
			ID:          uuid.NewString(),
			AggregateId: b.ID,
			Type:        "MoneyDeposited",
			OccurredAt:  time.Now(),
		},
		Amount: amount,
	})
}

func (b *BankAccount) Withdraw(amount float64) error {
	if amount > b.Balance {
		return errors.New("insufficient funds")
	}
	b.apply(&MoneyWithdrawn{
		BaseEvent: BaseEvent{
			ID:          uuid.NewString(),
			AggregateId: b.ID,
			Type:        "MoneyWithdrawn",
			OccurredAt:  time.Now(),
		},
		Amount: amount,
	})
	return nil
}

// Aplica el evento y muta el estado
func (b *BankAccount) apply(e Event) {
	b.Version++
	e.SetVersion(b.Version)
	utils.Info.Printf("** version apply to BankAccount %d, new version: %d", b.Version, e.GetVersion())
	switch ev := e.(type) {
	case *AccountCreated:
		b.Owner = ev.OwnerName
		b.Balance = ev.InitialBalance
	case *MoneyDeposited:
		b.Balance += ev.Amount
	case *MoneyWithdrawn:
		b.Balance -= ev.Amount
	}

	b.Events = append(b.Events, e)
}

func (b *BankAccount) Apply(e Event) {
	// Igual que apply pero sin registrar en Events (para reconstruir estado)
	utils.Warning.Printf("Version Apply , Bank Account v [%d] event version [%d]", b.Version, e.GetVersion())
	b.Version = e.GetVersion()
	switch ev := e.(type) {
	case *AccountCreated:
		b.Owner = ev.OwnerName
		b.Balance = ev.InitialBalance
	case *MoneyDeposited:
		b.Balance += ev.Amount
	case *MoneyWithdrawn:
		b.Balance -= ev.Amount
	}

}
