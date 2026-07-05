package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL non impostata nel file .env")
	}

	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatalf("Impossibile creare la connection pool: %v", err)
	}

	// Verifica che la connessione funzioni davvero
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Impossibile connettersi al database: %v", err)
	}

	Pool = pool
	log.Println("Connesso al database con successo")
}