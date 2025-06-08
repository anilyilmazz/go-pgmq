package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

func main() {
	fmt.Println("Starting the service...")

	config, err := pgxpool.ParseConfig("postgres://postgres:postgres@localhost:5432/postgres?search_path=pgmq,public&sslmode=disable")
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

	err = createQueue(context.Background(), dbPool, "my_queue")
	if err != nil {
		fmt.Printf("Error creating queue: %v\n", err)
		return
	}
	fmt.Println("Queue created successfully.")

	go startConsumer(dbPool)

	http.HandleFunc("/send", sendMessageHandler)
	fmt.Println("HTTP server listening on :8080")
	http.ListenAndServe(":8080", nil)
}

func createQueue(ctx context.Context, pool *pgxpool.Pool, queueName string) error {
	_, err := pool.Exec(ctx, "SELECT pgmq.create($1)", queueName)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") || strings.Contains(err.Error(), "duplicate") {
			return nil
		}
		return fmt.Errorf("failed to create queue: %w", err)
	}
	return nil
}
