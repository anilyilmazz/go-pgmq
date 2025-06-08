package main

import (
	"context"
	"fmt"
	"time"

	pgx "github.com/jackc/pgx/v5"
)

func main() {
	fmt.Println("Starting the message queue listener...")

	go startQueueListener()

	select {}
}

func startQueueListener() {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/postgres?search_path=pgmq,public&sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}

	fmt.Println("Successfully connected to the database!")

	for {
		var msgId int
		var msgBody string

		fmt.Println("Reading message from the queue...")

		err = conn.QueryRow(context.Background(),
			"SELECT msg_id, message FROM pgmq.read_with_poll($1, $2, $3, $4, $5, $6) LIMIT 1",
			"my_queue",
			30,
			1000,
			1,
			0,
			"{}").Scan(&msgId, &msgBody)

		if err != nil {
			if err == pgx.ErrNoRows {
				fmt.Println("No message available, waiting...")
				time.Sleep(1 * time.Second)
				continue
			}
			fmt.Printf("An error occurred: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Printf("Message ID: %d, Content: %s\n", msgId, msgBody)

		go handleMessage(msgId, msgBody, conn)
	}
}

func handleMessage(msgId int, msgContent string, conn *pgx.Conn) {
	fmt.Printf("Handling message ID: %d with content: %s\n", msgId, msgContent)

	_, err := conn.Exec(context.Background(), "SELECT pgmq.delete($1, $2)", "my_queue", msgId)
	if err != nil {
		fmt.Println("Error deleting message:", err)
		return
	}

	fmt.Printf("Message ID: %d successfully processed and deleted.\n", msgId)
}
