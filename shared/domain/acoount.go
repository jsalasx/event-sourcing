package domain

// AccountID es un alias para claridad en eventos y proyecciones.
type AccountID string

// AccountAggregate representa la información básica de la cuenta
// (usada principalmente para eventos y reconstrucción en Command Service).
type AccountAggregate struct {
	ID      AccountID `json:"id" bson:"id"`
	Owner   string    `json:"owner" bson:"owner"`
	Balance float64   `json:"balance" bson:"balance"`
}

// Snapshot es opcional: permite exportar el estado actual del agregado
// para proyecciones o depuración.
func (a *AccountAggregate) Snapshot() map[string]interface{} {
	return map[string]interface{}{
		"id":      a.ID,
		"owner":   a.Owner,
		"balance": a.Balance,
	}
}
