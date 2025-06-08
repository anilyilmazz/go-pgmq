package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func startConsumer(pool *pgxpool.Pool) {
	fmt.Println("Starting queue consumer...")

	for {
		var msgId int
		var msgBody string

		err := pool.QueryRow(context.Background(),
			"SELECT msg_id, message FROM pgmq.read_with_poll($1, $2, $3, $4, $5, $6) LIMIT 1",
			"my_queue",
			30,
			1000,
			1,
			0,
			"{}").Scan(&msgId, &msgBody)

		if err != nil {
			if err.Error() == "no rows in result set" {
				time.Sleep(1 * time.Second)
				continue
			}
			fmt.Printf("Error reading message: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Printf("Received message ID: %d, content: %s\n", msgId, msgBody)

		go handleMessage(msgId, msgBody, pool)
	}
}

func handleMessage(msgId int, msgContent string, pool *pgxpool.Pool) {
	fmt.Printf("Processing message ID: %d with content: %s\n", msgId, msgContent)

	_, err := pool.Exec(context.Background(), "SELECT pgmq.delete($1, $2::bigint)", "my_queue", msgId)
	if err != nil {
		fmt.Printf("Failed to delete message ID %d: %v\n", msgId, err)
		return
	}

	fmt.Printf("Message ID %d processed and deleted successfully.\n", msgId)
}
