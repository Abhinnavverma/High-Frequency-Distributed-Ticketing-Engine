package seats

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	repo *Repository
}

// NewHandler creates a new Seat Handler
func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

// GetSeats handles GET /seats
func (h *Handler) GetSeats(w http.ResponseWriter, r *http.Request) {
	// 1. Context is crucial for timeouts/cancellation
	ctx := r.Context()

	// 2. Ask the Repository for data
	seats, err := h.repo.GetAll(ctx)
	if err != nil {
		http.Error(w, "Failed to fetch seats", http.StatusInternalServerError)
		return
	}

	// 3. Send the JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(seats); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
