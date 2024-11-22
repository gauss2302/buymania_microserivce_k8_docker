package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gauss2302/microtest/product-service/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock usecase
type MockProductUsecase struct {
	mock.Mock
}

func (m *MockProductUsecase) CreateProduct(req *entity.CreateProductRequest) (*entity.Product, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *MockProductUsecase) GetProduct(id int) (*entity.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *MockProductUsecase) ListProducts(limit, offset int) ([]*entity.Product, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Product), args.Error(1)
}

func (m *MockProductUsecase) UpdateProduct(id int, req *entity.UpdateProductRequest) (*entity.Product, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *MockProductUsecase) DeleteProduct(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestProductHandler_CreateProduct(t *testing.T) {
	mockUsecase := new(MockProductUsecase)
	handler := NewProductHandler(mockUsecase)

	t.Run("successful creation", func(t *testing.T) {
		req := entity.CreateProductRequest{
			Name:        "Test Product",
			Description: "Test Description",
			Price:       99.99,
		}
		reqBody, _ := json.Marshal(req)

		expected := &entity.Product{
			ID:          1,
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockUsecase.On("CreateProduct", &req).Return(expected, nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(reqBody))
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response entity.Product
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expected.Name, response.Name)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("invalid request body", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer([]byte("invalid json")))
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("usecase error", func(t *testing.T) {
		req := entity.CreateProductRequest{
			Name:        "Test Product",
			Description: "Test Description",
			Price:       99.99,
		}
		reqBody, _ := json.Marshal(req)

		mockUsecase.On("CreateProduct", &req).Return(nil, errors.New("usecase error")).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(reqBody))
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUsecase.AssertExpectations(t)
	})
}

func TestProductHandler_GetProduct(t *testing.T) {
	mockUsecase := new(MockProductUsecase)
	handler := NewProductHandler(mockUsecase)

	t.Run("successful retrieval", func(t *testing.T) {
		productID := 1
		expected := &entity.Product{
			ID:          productID,
			Name:        "Test Product",
			Description: "Test Description",
			Price:       99.99,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockUsecase.On("GetProduct", productID).Return(expected, nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/products/1", nil)
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var response entity.Product
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expected.ID, response.ID)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/products/invalid", nil)
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		productID := 999
		mockUsecase.On("GetProduct", productID).Return(nil, errors.New("not found")).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/products/999", nil)
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockUsecase.AssertExpectations(t)
	})
}

func TestProductHandler_ListProducts(t *testing.T) {
	mockUsecase := new(MockProductUsecase)
	handler := NewProductHandler(mockUsecase)

	t.Run("successful listing", func(t *testing.T) {
		expected := []*entity.Product{
			{ID: 1, Name: "Product 1"},
			{ID: 2, Name: "Product 2"},
		}

		mockUsecase.On("ListProducts", 10, 0).Return(expected, nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/products", nil)
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []*entity.Product
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("usecase error", func(t *testing.T) {
		mockUsecase.On("ListProducts", 10, 0).Return(nil, errors.New("usecase error")).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/products", nil)
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUsecase.AssertExpectations(t)
	})
}

func TestProductHandler_MethodNotAllowed(t *testing.T) {
	mockUsecase := new(MockProductUsecase)
	handler := NewProductHandler(mockUsecase)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/products", nil)
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
}
