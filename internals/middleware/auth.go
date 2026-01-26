package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Define a custom key type to avoid context collisions
type contextKey string

const UserIDKey contextKey = "user_id"

type authMiddleware struct {
	jwtKey string
}

func NewMiddleware(jwtString string) (*authMiddleware, error) {
	if jwtString == "" {
		return nil, fmt.Errorf("jwt secret cannot be empty")
	}
	return &authMiddleware{jwtKey: jwtString}, nil
}

func (a *authMiddleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Get the header: "Authorization: Bearer <token>"
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// 2. Remove "Bearer " prefix
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		// 3. Parse & Validate Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			// In production, use os.Getenv("JWT_SECRET")
			return []byte(a.jwtKey), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// 4. Extract User ID
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// JSON numbers are often float64 in Go
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "Invalid user_id in token", http.StatusUnauthorized)
			return
		}
		userID := int32(userIDFloat)

		// 5. Inject into Context (The critical part!)
		ctx := context.WithValue(r.Context(), UserIDKey, userID)

		// 6. Pass the request down the chain
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
