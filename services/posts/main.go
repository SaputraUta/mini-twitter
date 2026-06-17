package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Tweet struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Text   string `json:"text"`
}

func main() {
	port := os.Getenv("POSTS_SERVICE_PORT")
	if port == "" {
		port = "8001"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	mux.HandleFunc("POST /tweet", func(w http.ResponseWriter, r *http.Request) {
		var t Tweet
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		if t.Text == "" || t.UserID == 0 {
			http.Error(w, "text and user_id required", http.StatusBadRequest)
			return
		}

		t.ID = 1
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(t)
	})

	addr := ":" + port
	log.Printf("posts-service listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
