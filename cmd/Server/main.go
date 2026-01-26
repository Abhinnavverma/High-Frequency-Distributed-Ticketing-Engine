package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"ticketmaster/internals/bookings"
	"ticketmaster/internals/cache"
	database "ticketmaster/internals/db"
	authMiddleware "ticketmaster/internals/middleware"
	"ticketmaster/internals/notifications"
	"ticketmaster/internals/seats"
	"ticketmaster/internals/users"
	"time"

	"github.com/go-chi/chi/v5"            // Import Chi
	"github.com/go-chi/chi/v5/middleware" // Import Middleware (Bonus!)
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system env")
	}
	dsn := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	jwtKey := os.Getenv("MY_JWT_KEY")

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379" // Fallback
	}
	redisPassword := os.Getenv("REDIS_PASSWORD")
	hub := notifications.NewHub()
	go hub.Run()
	tokenMiddleware, ok := authMiddleware.NewMiddleware(jwtKey)

	if ok != nil {
		log.Fatal("Error in setting up middleware")
	}
	db, err := database.NewDatabase(dsn)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}
	defer db.Close()
	redisStore := cache.NewRedisStore(redisAddr, redisPassword)
	log.Println("✅ Connected to Redis")

	// --- Services ---
	seatRepo := seats.NewRepository(db)
	seatHandler := seats.NewHandler(seatRepo)

	bookingRepo := bookings.NewRepository(db)
	bookingHandler := bookings.NewHandler(bookingRepo, redisStore, hub)

	userRepo := users.NewRepository(db)
	userService := users.NewService(userRepo, jwtKey) // <--- The new layer
	userHandler := users.NewHandler(userService)

	// --- Chi Router ---
	r := chi.NewRouter()

	// 1. Middleware (The reason Chi wins)
	r.Use(middleware.Logger)    // Log every request automatically
	r.Use(middleware.Recoverer) // Don't crash if a handler panics

	r.Post("/register", userHandler.Register)
	r.Post("/login", userHandler.Login)

	// 2. Routes (Clean Grouping)
	r.Get("/seats", seatHandler.GetSeats)
	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.ServeWs(w, r)
	})
	r.Group(func(r chi.Router) {
		// Apply the Bouncer
		r.Use(tokenMiddleware.Auth)

		// Authenticated users only
		r.Post("/bookings", bookingHandler.CreateBooking)
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r, // Pass your Chi router here
	}
	go func() {
		log.Println("🚀 Ticketmaster API running on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit // 🛑 BLOCK HERE until signal is received
	log.Println("⚠️  Shutting down server...")

	// 7. Create a deadline to wait for active requests (e.g., 5 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 8. Tell the server to stop accepting new requests and finish current ones
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// 9. Now safely close the Database
	log.Println("🔌 Closing Database Connection")
	db.Close()

	log.Println("✅ Server exited properly")
}
