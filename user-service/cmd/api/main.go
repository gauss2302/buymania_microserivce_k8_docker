package main

import (
	userHttp "github.com/gauss2302/microtest/user-service/internal/user/delivery/http"
	"github.com/gauss2302/microtest/user-service/internal/user/repository/postgres"
	"github.com/gauss2302/microtest/user-service/internal/user/usecase"
	"github.com/gauss2302/microtest/user-service/pkg/config"
	"github.com/gauss2302/microtest/user-service/pkg/db"

	"log"
	"net/http"
	"time"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Run database migrations
	if err := db.RunMigrations(
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize database connection
	dbConn, err := db.NewPostgresDB(
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	//Initialize application layers
	userRepo := postgres.NewUserRepository(dbConn)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := userHttp.NewUserHandler(userUsecase)

	//userHandler := userHttp.NewUserHandler(userUsecase)

	// Set up HTTP server with middleware
	mux := http.NewServeMux()

	// Health check handler with middleware
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	// User handler
	mux.Handle("/", userHandler)

	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("User Service running on port %s", cfg.ServerPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

//package main
//
//import (
//	"fmt"
//	"log"
//	"net/http"
//)
//
//func main() {
//	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
//		switch r.Method {
//		case http.MethodGet:
//			// Handle GET request
//			fmt.Fprintf(w, "Get Users")
//		case http.MethodPost:
//			// Handle POST request
//			fmt.Fprintf(w, "Create User")
//		default:
//			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//		}
//	})
//
//	port := ":8083"
//	fmt.Printf("User service is running on port %s\n", port)
//	if err := http.ListenAndServe(port, nil); err != nil {
//		log.Fatalf("Failed to start server: %v", err)
//	}
//}
