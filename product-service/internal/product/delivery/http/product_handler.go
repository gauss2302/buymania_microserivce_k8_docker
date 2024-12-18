package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gauss2302/microtest/product-service/internal/entity"
	"github.com/gauss2302/microtest/product-service/internal/product/usecase"
)

type ProductHandler struct {
	usecase usecase.ProductUsecase
}

func NewProductHandler(usecase usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{usecase: usecase}
}

func (h *ProductHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight request
	if r.Method == http.MethodOptions {
		return
	}

	log.Printf("Product service received: %s %s", r.Method, r.URL.Path)

	// Check the method
	switch r.Method {
	case http.MethodPost:
		if strings.HasSuffix(r.URL.Path, "/products/") || strings.HasSuffix(r.URL.Path, "/products") {
			log.Printf("Handling POST request for path: %s", r.URL.Path)
			h.CreateProduct(w, r)
			return
		}

	case http.MethodGet:
		// For GET requests, check the path
		if strings.HasSuffix(r.URL.Path, "/products/") || strings.HasSuffix(r.URL.Path, "/products") {
			h.GetProducts(w, r)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/products/") {
			h.GetProduct(w, r)
			return
		}

	default:
		log.Printf("Method %s not allowed for path: %s", r.Method, r.URL.Path)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.usecase.ListProducts(10, 0) // Временно хардкодим limit и offset
	if err != nil {
		log.Printf("Error getting products: %v", err)
		http.Error(w, "Failed to get products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(products); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	log.Printf("Starting CreateProduct handler")

	var req entity.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Decoded request: %+v", req)

	product, err := h.usecase.CreateProduct(&req)
	if err != nil {
		log.Printf("Error creating product: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Product created successfully: %+v", product)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid ID: %s", idStr)
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	log.Printf("Getting product with ID: %d", id)
	product, err := h.usecase.GetProduct(id)
	if err != nil {
		log.Printf("Error getting product: %v", err)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(product); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/products/"))
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var req entity.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	product, err := h.usecase.UpdateProduct(id, &req)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/products/"))
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	if err := h.usecase.DeleteProduct(id); err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
