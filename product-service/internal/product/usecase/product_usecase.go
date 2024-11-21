package usecase

import (
	"github.com/gauss2302/microtest/product-service/internal/entity"
	"github.com/gauss2302/microtest/product-service/internal/product/repository"
)

type productUsecase struct {
	repo repository.ProductRepository
}

func NewProductUsecase(repo repository.ProductRepository) ProductUsecase {
	return &productUsecase{repo: repo}
}

func (u *productUsecase) CreateProduct(req *entity.CreateProductRequest) (*entity.Product, error) {
	return u.repo.Create(req)
}

func (u *productUsecase) GetProduct(id int) (*entity.Product, error) {
	return u.repo.GetByID(id)
}

func (u *productUsecase) ListProducts(limit, offset int) ([]*entity.Product, error) {
	if limit <= 0 {
		limit = 10 // default limit
	}
	return u.repo.GetAll(limit, offset)
}

func (u *productUsecase) UpdateProduct(id int, req *entity.UpdateProductRequest) (*entity.Product, error) {
	return u.repo.Update(id, req)
}

func (u *productUsecase) DeleteProduct(id int) error {
	return u.repo.Delete(id)
}
