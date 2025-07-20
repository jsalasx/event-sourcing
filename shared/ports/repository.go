package ports

// AccountRead representa una cuenta en el Read Model.
type AccountRead struct {
	AccountID string  `bson:"account_id"`
	Owner     string  `bson:"owner"`
	Balance   float64 `bson:"balance"`
}

// AccountReadRepository define operaciones para consultar y actualizar la proyección.
type AccountReadRepository interface {
	Upsert(account AccountRead) error
	UpdateBalance(accountID string, delta float64) error
	FindByID(accountID string) (*AccountRead, error)
}
