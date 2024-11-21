package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gauss2302/microtest/product-service/internal/entity"
	"github.com/gauss2302/microtest/product-service/internal/product/repository"
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) repository.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(req *entity.CreateProductRequest) (*entity.Product, error) {
	query := `
		 INSERT INTO products (name, description, price)
		 VALUES ($1, $2, $3)
		 RETURNING id, name, description, price, created_at, updated_at`

	product := &entity.Product{}
	err := r.db.QueryRow(
		query,
		req.Name,
		req.Description,
		req.Price,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error creating product: %w", err)
	}

	return product, nil
}

func (r *productRepository) GetByID(id int) (*entity.Product, error) {
	query := `
		 SELECT id, name, description, price, created_at, updated_at
		 FROM products
		 WHERE id = $1`

	product := &entity.Product{}
	err := r.db.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error getting product: %w", err)
	}

	return product, nil
}

func (r *productRepository) GetAll(limit, offset int) ([]*entity.Product, error) {
	query := `
		 SELECT id, name, description, price, created_at, updated_at
		 FROM products
		 ORDER BY id
		 LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting products: %w", err)
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		product := &entity.Product{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning product: %w", err)
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *productRepository) Update(id int, req *entity.UpdateProductRequest) (*entity.Product, error) {
	// First, get the current product
	current, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		current.Name = *req.Name
	}
	if req.Description != nil {
		current.Description = *req.Description
	}
	if req.Price != nil {
		current.Price = *req.Price
	}

	query := `
		 UPDATE products
		 SET name = $1, description = $2, price = $3
		 WHERE id = $4
		 RETURNING id, name, description, price, created_at, updated_at`

	product := &entity.Product{}
	err = r.db.QueryRow(
		query,
		current.Name,
		current.Description,
		current.Price,
		id,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error updating product: %w", err)
	}

	return product, nil
}

func (r *productRepository) Delete(id int) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("product not found")
	}

	return nil
}
