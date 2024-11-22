package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"
)

type ServiceConfig struct {
	Name string
	URL  string
	Path string
}

func createProxy(serviceConfig ServiceConfig) http.HandlerFunc {
	serviceURL, err := url.Parse(serviceConfig.URL)
	if err != nil {
		log.Fatalf("Invalid URL for %s service: %v", serviceConfig.Name, err)
	}

	proxy := httputil.NewSingleHostReverseProxy(serviceURL)

	// Modify the director
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		log.Printf("Directing request: %s %s to %s", req.Method, req.URL.Path, serviceURL)
		originalDirector(req)
		// Preserve the original path
		if strings.HasPrefix(req.URL.Path, serviceConfig.Path) {
			req.URL.Path = strings.TrimPrefix(req.URL.Path, serviceConfig.Path)
			if req.URL.Path == "" {
				req.URL.Path = "/"
			}
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Single set of CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.Printf("API Gateway: proxying %s %s to %s", r.Method, r.URL.Path, serviceURL)
		proxy.ServeHTTP(w, r)
	}
}
func main() {
	services := []ServiceConfig{
		{
			Name: "Products",
			URL:  getEnv("PRODUCT_SERVICE_URL", "http://product-service:8082"),
			Path: "/products",
		},
		{
			Name: "Payments",
			URL:  getEnv("PAYMENT_SERVICE_URL", "http://payment-service:8081"),
			Path: "/payment/",
		},
		{
			Name: "Users",
			URL:  getEnv("USER_SERVICE_URL", "http://user-service:8083"),
			Path: "/users/",
		},
	}

	// Set up routes for each service
	for _, service := range services {
		path := service.Path
		handler := createProxy(service)

		log.Printf("Registering handler for %s at path: %s", service.Name, path)

		// Регистрируем и для корневого пути, и для путей с ID
		http.Handle(path+"/", handler)
		http.Handle(path, handler)
	}

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Root handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Welcome to the API Gateway!")
			return
		}
		http.NotFound(w, r)
	})

	port := getEnv("PORT", "8080")
	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("API Gateway starting on port %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
