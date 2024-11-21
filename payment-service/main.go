package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/payment", paymentHandler)

	fmt.Println("Payment service is running on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func paymentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Process payment here
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Payment processed successfully"))
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}
}
