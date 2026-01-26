package bookings

import (
	"context"
	"fmt"
	database "ticketmaster/internals/db"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

// CreateBooking attempts to book a seat inside a transaction
func (r *Repository) CreateBooking(ctx context.Context, seatID, userID int32) error {
	// 1. Start a Transaction
	// This opens a "sandbox" session. Nothing is permanent until we Commit.
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Safety Net: If anything fails or panics, Rollback changes.
	defer tx.Rollback(ctx)

	// 2. Lock the Seat (The Secret Sauce 🔒)
	// "FOR UPDATE" tells Postgres: "Lock this row. Make everyone else wait."
	var currentStatus string
	queryCheck := `SELECT status FROM seats WHERE id = $1 FOR UPDATE`

	err = tx.QueryRow(ctx, queryCheck, seatID).Scan(&currentStatus)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("seat %d does not exist", seatID)
		}
		return fmt.Errorf("failed to lock seat: %w", err)
	}

	// 3. The Logic Check
	if currentStatus != "available" {
		return fmt.Errorf("seat is already %s", currentStatus)
	}
	// 4. Update the Seat
	_, err = tx.Exec(ctx, `UPDATE seats SET status = 'booked' WHERE id = $1`, seatID)
	if err != nil {
		return fmt.Errorf("failed to update seat status: %w", err)
	}

	// 5. Create the Booking Record
	_, err = tx.Exec(ctx, `INSERT INTO bookings (seat_id, user_id) VALUES ($1, $2)`, seatID, userID)
	if err != nil {
		return fmt.Errorf("failed to insert booking: %w", err)
	}

	// 6. Commit (Make it permanent)
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
