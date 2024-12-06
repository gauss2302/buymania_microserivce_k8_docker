package main

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
