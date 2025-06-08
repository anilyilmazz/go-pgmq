package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

func main() {
	fmt.Println("Starting the service...")

	config, err := pgxpool.ParseConfig("postgres://postgres:postgres@pgmq-postgres:5432/postgres?search_path=pgmq,public&sslmode=disable")

	if err != nil {
		panic(err)
	}

	config.MaxConns = 5
	config.MinConns = 2

	dbPool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		panic(err)
	}
	defer dbPool.Close()

	fmt.Println("Connected to PostgreSQL database.")

	if err := EnsurePgMqExtension(dbPool); err != nil {
		fmt.Printf("Failed to ensure pgmq extension: %v", err)
		return
	}

	if err := CreateQueue(dbPool, "my_queue"); err != nil {
		fmt.Printf("Failed to create queue: %v\n", err)
		return
	}

	fmt.Println("Queue created successfully.")

	go startConsumer(dbPool)

	http.HandleFunc("/send", sendMessageHandler)
	fmt.Println("HTTP server listening on :8080")
	http.ListenAndServe(":8080", nil)
}

func EnsurePgMqExtension(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), `CREATE EXTENSION IF NOT EXISTS pgmq`)
	if err != nil {
		return fmt.Errorf("failed to create pgmq extension: %w", err)
	}
	return nil
}

func CreateQueue(db *pgxpool.Pool, queueName string) error {
	_, err := db.Exec(context.Background(), `SELECT pgmq.create($1)`, queueName)
	if err != nil {
		return fmt.Errorf("failed to create queue: %w", err)
	}
	fmt.Printf("Queue '%s' created successfully.\n", queueName)
	return nil
}
