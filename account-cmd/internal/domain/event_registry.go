package domain

import (
	"encoding/json"
	"fmt"
)

var eventRegistry = map[string]func() Event{}

// Registrar un evento por tipo
func RegisterEvent(eventType string, factory func() Event) {
	eventRegistry[eventType] = factory
}

// Reconstruir evento a partir del tipo y JSON
func BuildEvent(eventType string, data []byte) (Event, error) {
	factory, ok := eventRegistry[eventType]
	if !ok {
		return nil, fmt.Errorf("event type %s not registered", eventType)
	}
	evt := factory()
	if err := json.Unmarshal(data, evt); err != nil {
		return nil, err
	}
	return evt, nil
}
