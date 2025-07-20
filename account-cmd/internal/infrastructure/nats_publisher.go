package infrastructure

import (
	"encoding/json"
	"log"

	"account-cmd/internal/domain"
	"account-cmd/internal/ports"

	"github.com/nats-io/nats.go"
)

type JetStreamPublisher struct {
	js      nats.JetStreamContext
	stream  string
	subject string
}

// NewJetStreamPublisher inicializa el publisher y asegura que exista el stream.
func NewJetStreamPublisher(natsURL, streamName, subject string) (ports.EventPublisher, error) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	// Asegurar que el stream existe
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{subject},
		Storage:  nats.FileStorage,
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		return nil, err
	}

	return &JetStreamPublisher{
		js:      js,
		stream:  streamName,
		subject: subject,
	}, nil
}

// Publish envía un evento a JetStream y espera confirmación.
func (p *JetStreamPublisher) Publish(e domain.Event) error {
	data, _ := json.Marshal(e)
	subject := p.subjectPrefix(e.EventType())

	ack, err := p.js.Publish(subject, data)
	if err != nil {
		log.Printf("[JetStream] Error publicando evento %s: %v", e.EventType(), err)
		return err
	}

	log.Printf("[JetStream] Evento %s publicado (stream: %s, seq: %d)",
		e.EventType(), ack.Stream, ack.Sequence)
	return nil
}

// subjectPrefix genera el subject final (ej: bankaccount.AccountCreated).
func (p *JetStreamPublisher) subjectPrefix(eventType string) string {
	return "bankaccount." + eventType
}
