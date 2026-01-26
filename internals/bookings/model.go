package bookings

import (
	"time"
)

type BookingRequest struct {
	SeatID int32 `json:"seat_id"`
}

type Booking struct {
	ID        int32     `json:"id"`
	SeatID    int32     `json:"seat_id"`
	UserID    int32     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
