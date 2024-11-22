package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/gauss2302/microtest/product-service/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repository
type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(req *entity.CreateProductRequest) (*entity.Product, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *MockProductRepository) GetByID(id int) (*entity.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *MockProductRepository) GetAll(limit, offset int) ([]*entity.Product, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Product), args.Error(1)
}

func (m *MockProductRepository) Update(id int, req *entity.UpdateProductRequest) (*entity.Product, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *MockProductRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestProductUsecase_CreateProduct(t *testing.T) {
	mockRepo := new(MockProductRepository)
	usecase := NewProductUsecase(mockRepo)

	t.Run("successful creation", func(t *testing.T) {
		req := &entity.CreateProductRequest{
			Name:        "Test Product",
			Description: "Test Description",
			Price:       99.99,
		}

		expected := &entity.Product{
			ID:          1,
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockRepo.On("Create", req).Return(expected, nil).Once()

		result, err := usecase.CreateProduct(req)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		req := &entity.CreateProductRequest{
			Name:        "Test Product",
			Description: "Test Description",
			Price:       99.99,
		}

		mockRepo.On("Create", req).Return(nil, errors.New("repository error")).Once()

		result, err := usecase.CreateProduct(req)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestProductUsecase_GetProduct(t *testing.T) {
	mockRepo := new(MockProductRepository)
	usecase := NewProductUsecase(mockRepo)

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

		mockRepo.On("GetByID", productID).Return(expected, nil).Once()

		result, err := usecase.GetProduct(productID)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		productID := 999
		mockRepo.On("GetByID", productID).Return(nil, errors.New("not found")).Once()

		result, err := usecase.GetProduct(productID)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestProductUsecase_ListProducts(t *testing.T) {
	mockRepo := new(MockProductRepository)
	usecase := NewProductUsecase(mockRepo)

	t.Run("with valid limit", func(t *testing.T) {
		limit, offset := 10, 0
		expected := []*entity.Product{
			{ID: 1, Name: "Product 1"},
			{ID: 2, Name: "Product 2"},
		}

		mockRepo.On("GetAll", limit, offset).Return(expected, nil).Once()

		result, err := usecase.ListProducts(limit, offset)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("with zero limit", func(t *testing.T) {
		limit, offset := 0, 0
		expected := []*entity.Product{
			{ID: 1, Name: "Product 1"},
		}

		mockRepo.On("GetAll", 10, offset).Return(expected, nil).Once() // Default limit should be 10

		result, err := usecase.ListProducts(limit, offset)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestProductUsecase_UpdateProduct(t *testing.T) {
	mockRepo := new(MockProductRepository)
	usecase := NewProductUsecase(mockRepo)

	t.Run("successful update", func(t *testing.T) {
		productID := 1
		newName := "Updated Product"
		req := &entity.UpdateProductRequest{
			Name: &newName,
		}

		expected := &entity.Product{
			ID:          productID,
			Name:        newName,
			Description: "Test Description",
			Price:       99.99,
			UpdatedAt:   time.Now(),
		}

		mockRepo.On("Update", productID, req).Return(expected, nil).Once()

		result, err := usecase.UpdateProduct(productID, req)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		productID := 999
		newName := "Updated Product"
		req := &entity.UpdateProductRequest{
			Name: &newName,
		}

		mockRepo.On("Update", productID, req).Return(nil, errors.New("not found")).Once()

		result, err := usecase.UpdateProduct(productID, req)
		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestProductUsecase_DeleteProduct(t *testing.T) {
	mockRepo := new(MockProductRepository)
	usecase := NewProductUsecase(mockRepo)

	t.Run("successful deletion", func(t *testing.T) {
		productID := 1
		mockRepo.On("Delete", productID).Return(nil).Once()

		err := usecase.DeleteProduct(productID)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		productID := 999
		mockRepo.On("Delete", productID).Return(errors.New("not found")).Once()

		err := usecase.DeleteProduct(productID)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
