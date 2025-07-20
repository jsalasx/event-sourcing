package infrastructure

import (
	"log"
	"strings"
	"time"

	"shared/application"

	"github.com/nats-io/nats.go"
)

type JetStreamListener struct {
	nc         *nats.Conn
	js         nats.JetStreamContext
	projection application.Projection
}

func NewJetStreamListener(natsURL string, projection application.Projection) (*JetStreamListener, error) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}
	return &JetStreamListener{nc: nc, js: js, projection: projection}, nil
}

func (l *JetStreamListener) Start() error {
	streamName := "BANK_EVENTS"
	subject := "bankaccount.*"
	queueGroup := "account-query-workers"

	// Asegurar que el Stream existe
	_, err := l.js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{subject},
		Storage:  nats.FileStorage,
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		return err
	}

	// Suscribirse con ACK manual
	_, err = l.js.QueueSubscribe(subject, queueGroup, func(msg *nats.Msg) {
		eventType := strings.TrimPrefix(msg.Subject, "bankaccount.")
		log.Printf("[JetStream] Recibido evento: %s", eventType)

		// Procesar evento
		if err := l.projection.HandleEvent(eventType, msg.Data); err != nil {
			log.Printf("[JetStream] Error procesando: %v", err)
			// No hacemos Ack → el mensaje será redeliver
			return
		}

		// Confirmar que fue procesado correctamente
		if err := msg.Ack(); err != nil {
			log.Printf("[JetStream] Error en ACK: %v", err)
		}
	}, nats.ManualAck(), nats.AckWait(30*time.Second))
	if err != nil {
		return err
	}

	log.Println("[JetStream] Suscrito en", subject, "con ACKs y grupo", queueGroup)
	return nil
}
