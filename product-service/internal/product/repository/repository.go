package repository

import "github.com/gauss2302/microtest/product-service/internal/entity"

type ProductRepository interface {
	Create(product *entity.CreateProductRequest) (*entity.Product, error)
	GetByID(id int) (*entity.Product, error)
	GetAll(limit, offset int) ([]*entity.Product, error)
	Update(id int, product *entity.UpdateProductRequest) (*entity.Product, error)
	Delete(id int) error
}
