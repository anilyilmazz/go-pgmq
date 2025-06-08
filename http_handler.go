package main

import (
	"context"
	"encoding/json"
	"net/http"
)

type Message struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
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
