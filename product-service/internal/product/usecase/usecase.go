package usecase

import "github.com/gauss2302/microtest/product-service/internal/entity"

type ProductUsecase interface {
	CreateProduct(req *entity.CreateProductRequest) (*entity.Product, error)
	GetProduct(id int) (*entity.Product, error)
	ListProducts(limit, offset int) ([]*entity.Product, error)
	UpdateProduct(id int, req *entity.UpdateProductRequest) (*entity.Product, error)
	DeleteProduct(id int) error
}
