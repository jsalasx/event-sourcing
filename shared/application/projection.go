package application

// Projection define cómo un proyector aplica cambios en el Read Model.
type Projection interface {
	HandleEvent(eventType string, data []byte) error
}
