package application

import (
	"account-cmd/internal/domain"
	"account-cmd/internal/ports"
	"fmt"
	"time"
)

type AccountService struct {
	store     ports.EventStore
	publisher ports.EventPublisher
	snapshot  ports.SnapshotStore
}

func NewAccountService(s ports.EventStore, p ports.EventPublisher, ss ports.SnapshotStore) *AccountService {
	return &AccountService{store: s, publisher: p, snapshot: ss}
}

func (svc *AccountService) CreateAccount(owner string, initial float64) (*domain.BankAccount, error) {
	acc := domain.NewAccount(owner, initial)
	if err := svc.store.Save(acc.ID, acc.Events); err != nil {
		return nil, err
	}
	for _, e := range acc.Events {
		_ = svc.publisher.Publish(e)
	}
	return acc, nil
}

func (svc *AccountService) GetEvents(id string, acc *domain.BankAccount) ([]domain.Event, error) {

	lastSnapshot, err := svc.snapshot.LoadSnapshot(id)
	if err != nil {
		return nil, err
	}
	var events []domain.Event
	if lastSnapshot == nil {
		acc.ID = id
		acc.Version = 0
		events, err = svc.store.Load(id)
		if err != nil {
			return nil, err
		}
	} else {
		acc.ID = id
		acc.Version = lastSnapshot.Version
		acc.Owner = lastSnapshot.State.Owner
		acc.Balance = lastSnapshot.State.Balance
		events, err = svc.store.LoadWithVersion(id, lastSnapshot.Version)
		if err != nil {
			return nil, err
		}
		if events == nil {
			return []domain.Event{}, nil
		}
	}
	return events, nil
}

func (svc *AccountService) SaveSnapshot(id string, acc *domain.BankAccount, version int) error {

	if version%3 == 0 {
		return svc.snapshot.UpsertSnapshot(domain.BankAccountSnapshot{
			AggregateID: id,
			Version:     acc.Version,
			State: domain.BankAccountState{
				ID:      acc.ID,
				Owner:   acc.Owner,
				Balance: acc.Balance,
			},
			CreatedAt: time.Now(),
		})
	}
	return nil
}

func (svc *AccountService) Deposit(id string, amount float64) (*domain.BankAccount, error) {
	// Cargar eventos del store
	acc := &domain.BankAccount{ID: id, Version: 0}
	events, err := svc.GetEvents(id, acc)
	if err != nil {
		return nil, err
	}

	if len(events) == 0 && acc.Version == 0 {
		return nil, fmt.Errorf("account with ID %s not found", id)
	}

	for _, e := range events {
		acc.Apply(e)

	}

	acc.Deposit(amount)
	if err := svc.store.Save(id, acc.Events); err != nil {
		return nil, err
	}

	err = svc.SaveSnapshot(id, acc, int(acc.Version))
	if err != nil {
		return nil, err
	}

	for _, e := range acc.Events {
		_ = svc.publisher.Publish(e)
	}
	return acc, nil
}

func (svc *AccountService) Withdraw(id string, amount float64) (*domain.BankAccount, error) {
	acc := &domain.BankAccount{ID: id, Version: 0}
	events, err := svc.GetEvents(id, acc)
	if err != nil {
		return nil, err
	}

	if len(events) == 0 && acc.Version == 0 {
		return nil, fmt.Errorf("account with ID %s not found", id)
	}

	for _, e := range events {
		acc.Apply(e)
	}
	if err := acc.Withdraw(amount); err != nil {
		return nil, err
	}

	err = svc.SaveSnapshot(id, acc, int(acc.Version))
	if err != nil {
		return nil, err
	}

	if err := svc.store.Save(id, acc.Events); err != nil {
		return nil, err
	}
	for _, e := range acc.Events {
		_ = svc.publisher.Publish(e)
	}
	return acc, nil
}
