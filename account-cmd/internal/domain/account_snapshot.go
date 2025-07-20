package domain

import "time"

type BankAccountSnapshot struct {
	AggregateID string           `bson:"aggregate_id"`
	Version     uint64           `bson:"version"`
	State       BankAccountState `bson:"state"`
	CreatedAt   time.Time        `bson:"created_at"`
}

type BankAccountState struct {
	ID      string  `bson:"id"`
	Owner   string  `bson:"owner"`
	Balance float64 `bson:"balance"`
}
