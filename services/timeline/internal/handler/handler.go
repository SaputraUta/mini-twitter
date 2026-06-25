package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/SaputraUta/mini-twitter/services/timeline/internal/service"
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

func (h *Handler) Timeline(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("user"), 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	tweets, err := h.svc.Timeline(userID)
	if err != nil {
		log.Printf("get timeline %d: %v", userID, err)
		http.Error(w, "could not get timeline", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"user_id": userID, "tweets": tweets})
}

func (h *Handler) TimelineDB(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("user"), 10, 64)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}
	tweets, err := h.svc.TimelineFromDB(userID)
	if err != nil {
		log.Printf("get timeline-db %d: %v", userID, err)
		http.Error(w, "could not get timeline", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"user_id": userID, "tweets": tweets})
}
