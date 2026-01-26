package users

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

func (r *Repository) CreateUser(ctx context.Context, email, passwordHash string) error {
	_, err := r.db.Pool.Exec(ctx, "INSERT INTO users (email, password_hash) VALUES ($1, $2)", email, passwordHash)
	return err
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	query := "SELECT id, email, password_hash, created_at FROM users WHERE email = $1"
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &u, nil
}
