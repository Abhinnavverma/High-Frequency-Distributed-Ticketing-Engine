package bookings

import (
	"encoding/json"
	"fmt"
	"net/http"
	"ticketmaster/internals/cache"
	"ticketmaster/internals/middleware"
	"ticketmaster/internals/notifications"
)

type Handler struct {
	repo       *Repository
	hub        *notifications.Hub
	redisStore *cache.RedisStore
}

func NewHandler(repo *Repository, redisStore *cache.RedisStore, hub *notifications.Hub) *Handler {
	return &Handler{repo: repo, hub: hub, redisStore: redisStore}
}

func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var req BookingRequest

	// Decode JSON body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	userID, ok := r.Context().Value(middleware.UserIDKey).(int32)
	if !ok {
		// This should never happen if the middleware is running
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	lockKey := fmt.Sprintf("seat_lock:%d", req.SeatID)

	// Attempt to acquire the lock in Redis (Atomic Lua Script)
	// We set a 60-second expiry just in case the server crashes before DB write
	if err := h.redisStore.AtomicBook(r.Context(), lockKey, userID, 60); err != nil {
		// 🛑 STOP! Redis says the seat is taken.
		// Return 409 Conflict immediately. Do not touch Postgres.
		http.Error(w, "Seat is currently reserved or booked", http.StatusConflict)
		return
	}

	// Call the logic
	if err := h.repo.CreateBooking(r.Context(), req.SeatID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	go func() {
		msg := map[string]interface{}{
			"type":    "seat_booked",
			"seat_id": req.SeatID,
			"user_id": userID,
		}
		jsonMsg, _ := json.Marshal(msg)
		h.hub.Broadcast <- jsonMsg
	}()

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Booked!"))
}
