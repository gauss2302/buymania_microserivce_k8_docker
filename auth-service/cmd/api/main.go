package main

import (
	memcachediml "github.com/gauss2302/microtest/auth-service/pkg/memcached"
	"log"
	"net/http"
	"time"

	httpauth "github.com/gauss2302/microtest/auth-service/internal/auth/delivery/http"
	"github.com/gauss2302/microtest/auth-service/internal/auth/repository/memcached"
	"github.com/gauss2302/microtest/auth-service/internal/auth/usecase"
	"github.com/gauss2302/microtest/auth-service/internal/middleware"
	"github.com/gauss2302/microtest/auth-service/pkg/config"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize memcached client
	memcachedWrapper := memcachediml.NewClient(cfg.MemcachedHost, cfg.MemcachedPort)

	// Initialize rate limiters
	globalLimiter := middleware.NewRateLimiter(100, 10) // 100 requests max, refill 10 per second
	clientLimiter := middleware.NewClientRateLimiter(middleware.ClientRateLimiterConfig{
		Capacity:        10,
		RefillRate:      1,
		CleanupInterval: time.Hour,
	})

	// Initialize application layers
	authRepo := memcached.NewAuthRepository(memcachedWrapper.Client) // Pass the raw *memcache.Client
	authUsecase := usecase.NewAuthUsecase(
		authRepo,
		cfg.TokenExpiration,
		cfg.UserServiceURL,
	)
	authHandler := httpauth.NewAuthHandler(authUsecase)

	// Set up HTTP server with middleware
	mux := http.NewServeMux()

	// Apply middleware stack
	handler := middleware.Logging(
		middleware.CORS(
			middleware.RateLimit(globalLimiter)(
				middleware.PerClientRateLimit(clientLimiter)(
					authHandler,
				),
			),
		),
	)

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// Auth handler
	mux.Handle("/auth/", handler)

	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Auth Service running on port %s", cfg.ServerPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
