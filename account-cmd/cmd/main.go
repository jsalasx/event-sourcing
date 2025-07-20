package main

import (
	"account-cmd/internal/application"
	"account-cmd/internal/domain"
	"account-cmd/internal/infrastructure"
	"log"
	"os"
)

func main() {

	natsUrl := os.Getenv("NATS_URL")
	mongoUrl := os.Getenv("MONGO_URL")

	if natsUrl == "" {
		panic("NATS_URL environment variable is required")
	}

	log.Println("NATS_URL:", natsUrl)

	if mongoUrl == "" {
		panic("MONGO_URL environment variable is required")
	}

	publisher, err := infrastructure.NewJetStreamPublisher(
		natsUrl,
		"BANK_EVENTS",   // nombre del stream
		"bankaccount.*", // patrón de subjects
	)
	if err != nil {
		log.Fatal("Error creando publisher:", err)
	}
	store, _ := infrastructure.NewMongoEventStore(mongoUrl, "bank", "events")
	snapshotStore, _ := infrastructure.NewMongoSnapshotStore(mongoUrl, "bank", "snapshots")

	service := application.NewAccountService(store, publisher, snapshotStore)
	server := infrastructure.NewHTTPServer(service)

	domain.RegisterEvent("AccountCreated", func() domain.Event { return &domain.AccountCreated{} })
	domain.RegisterEvent("MoneyDeposited", func() domain.Event { return &domain.MoneyDeposited{} })
	domain.RegisterEvent("MoneyWithdrawn", func() domain.Event { return &domain.MoneyWithdrawn{} })

	log.Println("Servidor escuchando en http://localhost:8080")
	if err := server.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
