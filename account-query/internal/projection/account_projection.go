package projection

import (
	"encoding/json"
	"log"

	"shared/domain"
	"shared/ports"
)

type AccountProjection struct {
	repo ports.AccountReadRepository
}

func NewAccountProjection(r ports.AccountReadRepository) *AccountProjection {
	return &AccountProjection{repo: r}
}

func (p *AccountProjection) HandleEvent(eventType string, data []byte) error {
	switch eventType {
	case "AccountCreated":
		var e domain.AccountCreated
		if err := json.Unmarshal(data, &e); err != nil {
			return err
		}
		return p.repo.Upsert(ports.AccountRead{
			AccountID: e.AggregateId,
			Owner:     e.OwnerName,
			Balance:   e.InitialBalance,
		})

	case "MoneyDeposited":
		var e domain.MoneyDeposited
		if err := json.Unmarshal(data, &e); err != nil {
			return err
		}
		return p.repo.UpdateBalance(e.AggregateId, e.Amount)

	case "MoneyWithdrawn":
		var e domain.MoneyWithdrawn
		if err := json.Unmarshal(data, &e); err != nil {
			return err
		}
		return p.repo.UpdateBalance(e.AggregateId, -e.Amount)

	default:
		log.Println("Evento no reconocido:", eventType)
		return nil
	}
}
