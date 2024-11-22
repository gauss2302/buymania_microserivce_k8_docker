package postgres

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gauss2302/microtest/product-service/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestProductRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewProductRepository(db)

	t.Run("successful creation", func(t *testing.T) {
		req := &entity.CreateProductRequest{
			Name:        "Test Product",
			Description: "Test Description",
			Price:       99.99,
		}

		rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "created_at", "updated_at"}).
			AddRow(1, req.Name, req.Description, req.Price, time.Now(), time.Now())

		mock.ExpectQuery("INSERT INTO products").
			WithArgs(req.Name, req.Description, req.Price).
			WillReturnRows(rows)

		product, err := repo.Create(req)
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, req.Name, product.Name)
	})

	t.Run("database error", func(t *testing.T) {
		req := &entity.CreateProductRequest{
			Name:        "Test Product",
			Description: "Test Description",
			Price:       99.99,
		}

		mock.ExpectQuery("INSERT INTO products").
			WithArgs(req.Name, req.Description, req.Price).
			WillReturnError(sql.ErrConnDone)

		product, err := repo.Create(req)
		assert.Error(t, err)
		assert.Nil(t, product)
	})
}

func TestProductRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewProductRepository(db)

	t.Run("found", func(t *testing.T) {
		productID := 1
		rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "created_at", "updated_at"}).
			AddRow(productID, "Test Product", "Description", 99.99, time.Now(), time.Now())

		mock.ExpectQuery("SELECT (.+) FROM products").
			WithArgs(productID).
			WillReturnRows(rows)

		product, err := repo.GetByID(productID)
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, productID, product.ID)
	})

	t.Run("not found", func(t *testing.T) {
		productID := 999

		mock.ExpectQuery("SELECT (.+) FROM products").
			WithArgs(productID).
			WillReturnError(sql.ErrNoRows)

		product, err := repo.GetByID(productID)
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Equal(t, "product not found", err.Error())
	})
}

func TestProductRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewProductRepository(db)

	t.Run("successful retrieval", func(t *testing.T) {
		limit, offset := 10, 0
		rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "created_at", "updated_at"}).
			AddRow(1, "Product 1", "Desc 1", 99.99, time.Now(), time.Now()).
			AddRow(2, "Product 2", "Desc 2", 199.99, time.Now(), time.Now())

		mock.ExpectQuery("SELECT (.+) FROM products").
			WithArgs(limit, offset).
			WillReturnRows(rows)

		products, err := repo.GetAll(limit, offset)
		assert.NoError(t, err)
		assert.Len(t, products, 2)
	})

	t.Run("database error", func(t *testing.T) {
		limit, offset := 10, 0

		mock.ExpectQuery("SELECT (.+) FROM products").
			WithArgs(limit, offset).
			WillReturnError(sql.ErrConnDone)

		products, err := repo.GetAll(limit, offset)
		assert.Error(t, err)
		assert.Nil(t, products)
	})
}

func TestProductRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewProductRepository(db)

	t.Run("successful update", func(t *testing.T) {
		productID := 1
		newName := "Updated Product"
		req := &entity.UpdateProductRequest{
			Name: &newName,
		}

		// Mock GetByID
		selectRows := sqlmock.NewRows([]string{"id", "name", "description", "price", "created_at", "updated_at"}).
			AddRow(productID, "Old Name", "Description", 99.99, time.Now(), time.Now())

		mock.ExpectQuery("SELECT (.+) FROM products").
			WithArgs(productID).
			WillReturnRows(selectRows)

		// Mock Update
		updateRows := sqlmock.NewRows([]string{"id", "name", "description", "price", "created_at", "updated_at"}).
			AddRow(productID, newName, "Description", 99.99, time.Now(), time.Now())

		mock.ExpectQuery("UPDATE products").
			WillReturnRows(updateRows)

		product, err := repo.Update(productID, req)
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, newName, product.Name)
	})

	t.Run("not found", func(t *testing.T) {
		productID := 999
		newName := "Updated Product"
		req := &entity.UpdateProductRequest{
			Name: &newName,
		}

		mock.ExpectQuery("SELECT (.+) FROM products").
			WithArgs(productID).
			WillReturnError(sql.ErrNoRows)

		product, err := repo.Update(productID, req)
		assert.Error(t, err)
		assert.Nil(t, product)
	})
}

func TestProductRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewProductRepository(db)

	t.Run("successful deletion", func(t *testing.T) {
		productID := 1

		mock.ExpectExec("DELETE FROM products").
			WithArgs(productID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(productID)
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		productID := 999

		mock.ExpectExec("DELETE FROM products").
			WithArgs(productID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete(productID)
		assert.Error(t, err)
		assert.Equal(t, "product not found", err.Error())
	})
}
