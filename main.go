package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

type Message struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

func main() {
	fmt.Println("Starting the service...")

	// Initialize the connection pool
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

	go startQueueListener()

	http.HandleFunc("/send", sendMessageHandler)
	fmt.Println("HTTP server started at :8080")
	http.ListenAndServe(":8080", nil)
}

func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	jsonPayload, err := json.Marshal(map[string]string{
		"type":    msg.Type,
		"payload": msg.Payload,
	})
	if err != nil {
		http.Error(w, "Failed to encode payload", http.StatusInternalServerError)
		return
	}

	_, err = dbPool.Exec(context.Background(), "SELECT pgmq.send($1, $2::jsonb)", "my_queue", string(jsonPayload))
	if err != nil {
		http.Error(w, "Failed to enqueue message: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Message enqueued successfully"))
}

func startQueueListener() {
	config, err := pgxpool.ParseConfig("postgres://postgres:postgres@localhost:5432/postgres?search_path=pgmq,public&sslmode=disable")
	if err != nil {
		fmt.Println("Error parsing config:", err)
		return
	}

	config.MaxConns = 5
	config.MinConns = 2

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		fmt.Println("Error creating connection pool:", err)
		return
	}
	defer pool.Close()

	fmt.Println("Successfully connected to the database!")

	for {
		var msgId int
		var msgBody string

		fmt.Println("Reading message from the queue...")

		err = pool.QueryRow(context.Background(),
			"SELECT msg_id, message FROM pgmq.read_with_poll($1, $2, $3, $4, $5, $6) LIMIT 1",
			"my_queue",
			30,
			1000,
			1,
			0,
			"{}").Scan(&msgId, &msgBody)

		if err != nil {
			if err.Error() == "no rows in result set" {
				fmt.Println("No message available, waiting...")
				time.Sleep(1 * time.Second)
				continue
			}
			fmt.Printf("An error occurred: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Printf("Message ID: %d, Content: %s\n", msgId, msgBody)

		go handleMessage(msgId, msgBody, pool)
	}
}

func handleMessage(msgId int, msgContent string, pool *pgxpool.Pool) {
	fmt.Printf("Handling message ID: %d with content: %s\n", msgId, msgContent)

	_, err := pool.Exec(context.Background(), "SELECT pgmq.delete($1, $2::bigint)", "my_queue", msgId)
	if err != nil {
		fmt.Println("Error deleting message:", err)
		return
	}

	fmt.Printf("Message ID: %d successfully processed and deleted.\n", msgId)
}
