package users

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo   *Repository
	jwtKey string
}

func NewService(repo *Repository, jwt string) *Service {
	return &Service{repo: repo, jwtKey: jwt}
}

// Register handles hashing and saving
func (s *Service) Register(ctx context.Context, req RegisterRequest) error {
	// 1. Hash Password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 2. Save to Repo
	return s.repo.CreateUser(ctx, req.Email, string(hashed))
}

// Login handles verification and token generation
func (s *Service) Login(ctx context.Context, req LoginRequest) (string, error) {
	// 1. Find User
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// 2. Verify Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// 3. Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	// Note: In production, move "MY_SUPER_SECRET_KEY" to .env
	return token.SignedString([]byte(s.jwtKey))
}
