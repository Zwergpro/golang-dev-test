package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"homework-1/internal/models/products"
	"regexp"
	"testing"
)

func TestGetProductByID(t *testing.T) {
	t.Run("success getting product by id", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		mockResponse := pgxmock.NewRows([]string{"id", "name", "price", "quantity"}).AddRow(uint64(1), "product1", uint64(1), uint64(1))
		f.mockPool.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, price, quantity FROM products WHERE id = $1`)).
			WithArgs(uint64(1)).
			WillReturnRows(mockResponse)

		// act
		res, err := f.productRepo.GetProductById(context.Background(), 1)

		// assert
		require.NoError(t, err)
		assert.Equal(t, res, &products.Product{
			Id:       uint64(1),
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		})
	})

	t.Run("product by id not found", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		f.mockPool.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, price, quantity FROM products WHERE id = $1`)).
			WithArgs(uint64(1)).
			WillReturnError(pgx.ErrNoRows)

		// act
		_, err := f.productRepo.GetProductById(context.Background(), 1)

		// assert
		assert.EqualError(t, err, "1: product does not exist")
	})

	t.Run("getting with internal error", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		f.mockPool.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, price, quantity FROM products WHERE id = $1`)).
			WithArgs(uint64(1)).
			WillReturnError(errors.New("internal error"))

		// act
		_, err := f.productRepo.GetProductById(context.Background(), 1)

		// assert
		assert.EqualError(t, err, "Repository.GetProductById: select: scany: query one result row: internal error")
	})
}

func TestCreateProduct(t *testing.T) {
	t.Run("success creating product", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		f.mockPool.ExpectQuery(regexp.QuoteMeta(`INSERT INTO products (name, price, quantity) VALUES ($1,$2,$3) RETURNING id`)).
			WithArgs("product1", uint64(1), uint64(1)).
			WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(uint64(1)))

		// act
		res, err := f.productRepo.CreateProduct(context.Background(), products.Product{
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		})

		// assert
		require.NoError(t, err)
		assert.Equal(t, res, &products.Product{
			Id:       uint64(1),
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		})
	})

	t.Run("creating with internal error", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		f.mockPool.ExpectQuery(regexp.QuoteMeta(`INSERT INTO products (name, price, quantity) VALUES ($1,$2,$3) RETURNING id`)).
			WithArgs("product1", uint64(1), uint64(1)).
			WillReturnError(errors.New("internal error"))

		// act
		_, err := f.productRepo.CreateProduct(context.Background(), products.Product{
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		})

		// assert
		assert.EqualError(t, err, "Repository.CreateProduct: insert: internal error")
	})
}

func TestDeleteProduct(t *testing.T) {
	t.Run("success deleting product", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		f.mockPool.ExpectExec(regexp.QuoteMeta(`DELETE FROM products WHERE id = $1`)).
			WithArgs(uint64(1)).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		// act
		err := f.productRepo.DeleteProduct(context.Background(), 1)

		// assert
		require.NoError(t, err)
	})

	t.Run("deleting with internal error", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		f.mockPool.ExpectExec(regexp.QuoteMeta(`DELETE FROM products WHERE id = $1`)).
			WithArgs(uint64(1)).
			WillReturnError(errors.New("internal error"))

		// act
		err := f.productRepo.DeleteProduct(context.Background(), 1)

		// assert
		assert.EqualError(t, err, "Repository.DeleteProduct: to delete: internal error")
	})
}

func TestUpdateProduct(t *testing.T) {
	t.Run("success updating product", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		f.mockPool.ExpectExec(regexp.QuoteMeta(`UPDATE products SET name = $1, price = $2, quantity = $3 WHERE id = $4`)).
			WithArgs("product1", uint64(1), uint64(1), uint64(1)).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		// act
		res, err := f.productRepo.UpdateProduct(context.Background(), products.Product{
			Id:       uint64(1),
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		})

		// assert
		require.NoError(t, err)
		assert.Equal(t, res, &products.Product{
			Id:       uint64(1),
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		})
	})

	t.Run("updating with internal error", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		f.mockPool.ExpectExec(regexp.QuoteMeta(`UPDATE products SET name = $1, price = $2, quantity = $3 WHERE id = $4`)).
			WithArgs("product1", uint64(1), uint64(1), uint64(1)).
			WillReturnError(errors.New("internal error"))

		// act
		_, err := f.productRepo.UpdateProduct(context.Background(), products.Product{
			Id:       uint64(1),
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		})

		// assert
		assert.EqualError(t, err, "Repository.UpdateProduct: to update: internal error")
	})
}

func TestGetAllProducts(t *testing.T) {
	t.Run("success getting all products", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		f.mockPool.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, price, quantity FROM products ORDER BY id LIMIT 2 OFFSET 0`)).
			WillReturnRows(pgxmock.NewRows([]string{"id", "name", "price", "quantity"}).
				AddRow(uint64(1), "product1", uint64(1), uint64(1)).
				AddRow(uint64(2), "product2", uint64(2), uint64(2)))

		// act
		res, err := f.productRepo.GetAllProducts(context.Background(), uint64(1), uint64(2))

		// assert
		require.NoError(t, err)
		assert.Equal(t, res, []*products.Product{
			{
				Id:       uint64(1),
				Name:     "product1",
				Price:    uint64(1),
				Quantity: uint64(1),
			},
			{
				Id:       uint64(2),
				Name:     "product2",
				Price:    uint64(2),
				Quantity: uint64(2),
			},
		})
	})

	t.Run("getting with internal error", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		f.mockPool.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, price, quantity FROM products ORDER BY id LIMIT 2 OFFSET 0`)).
			WillReturnError(errors.New("internal error"))

		// act
		_, err := f.productRepo.GetAllProducts(context.Background(), uint64(1), uint64(2))

		// assert
		assert.EqualError(t, err, "Repository.GetAllProducts: select: scany: query multiple result rows: internal error")
	})
}

func TestGetPaginationLimitAndOffset(t *testing.T) {
	t.Run("success getting pagination limit and offset", func(t *testing.T) {
		// act
		limit, offset := (&Repository{}).getPaginationLimitAndOffset(uint64(1), uint64(2))

		// assert
		assert.Equal(t, limit, uint64(2))
		assert.Equal(t, offset, uint64(0))
	})

	t.Run("success getting pagination limit and offset with offset greater than 0", func(t *testing.T) {
		// act
		limit, offset := (&Repository{}).getPaginationLimitAndOffset(uint64(3), uint64(3))

		// assert
		assert.Equal(t, limit, uint64(3))
		assert.Equal(t, offset, uint64(6))
	})

	t.Run("success getting pagination limit and offset with 0 values", func(t *testing.T) {
		// act
		limit, offset := (&Repository{}).getPaginationLimitAndOffset(uint64(0), uint64(0))

		// assert
		assert.Equal(t, limit, defaultProductsPageSize)
		assert.Equal(t, offset, uint64(0))
	})
}
