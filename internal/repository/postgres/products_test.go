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
