package main

import (
	"account-query/internal/infrastructure"
	"account-query/internal/infrastructure/repository"
	"account-query/internal/projection"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	mongoURI := os.Getenv("MONGO_URL")
	natsURL := os.Getenv("NATS_URL")
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	// Crear el repositorio de lectura (Mongo)
	repo, err := repository.NewMongoAccountRepository(mongoURI, "bank_read", "account_readmodel")
	if err != nil {
		log.Fatal("Error conectando a Mongo:", err)
	}

	// Crear el proyector (lógica de negocio para read model)
	proj := projection.NewAccountProjection(repo)

	// Iniciar el listener de NATS (adaptador de entrada)
	listener, err := infrastructure.NewJetStreamListener(natsURL, proj)
	if err != nil {
		log.Fatal("Error conectando a NATS:", err)
	}

	if err := listener.Start(); err != nil {
		log.Fatal("Error iniciando listener:", err)
	}

	// Iniciar el servidor HTTP para consultas
	httpServer := infrastructure.NewHTTPServer(repo)

	// Ejecutar NATS y HTTP en paralelo
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := httpServer.Listen(":" + port); err != nil {
			log.Fatal("Error en servidor HTTP:", err)
		}
	}()

	log.Println("Account Query Service listo. API en puerto", port)
	wg.Wait()
}
