package seats

import (
	"context"
	database "ticketmaster/internals/db"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *database.DB
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateSeat(ctx context.Context, rowNo string, seatNo int32, price int32) error {
	query := `INSERT INTO seats (row_number, seat_number, status,price) VALUES ($1,$2,'AVAILABLE',$3)`
	_, err := r.db.Pool.Exec(ctx, query, rowNo, seatNo, price)
	if err != nil {
		return err
	}
	return nil

}

func (r *Repository) GetAll(ctx context.Context) ([]Seat, error) {
	// Query remains the same
	query := `SELECT id, row_number, seat_number, status, price FROM seats ORDER BY row_number, seat_number ASC`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Use pgx to map to the Seat struct defined in this same package
	return pgx.CollectRows(rows, pgx.RowToStructByName[Seat])
}
