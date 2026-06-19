package handler

import (
	"encoding/json"
	"net/http"

	"github.com/SaputraUta/mini-twitter/services/posts/internal/model"
	"github.com/SaputraUta/mini-twitter/services/posts/internal/service"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) CreateTweet(w http.ResponseWriter, r *http.Request) {
	var t model.Tweet
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if t.Text == "" || t.UserID == 0 {
		http.Error(w, "text and user_id required", http.StatusBadRequest)
		return
	}

	created, err := h.svc.CreateTweet(t)
	if err != nil {
		http.Error(w, "could not save tweet", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}
