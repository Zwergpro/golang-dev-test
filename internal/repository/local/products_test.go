package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"homework-1/internal/models/products"
	"testing"
)

func TestGetProductByID(t *testing.T) {
	t.Run("success getting product by id", func(t *testing.T) {
		// arrange
		f := SetUp(t)

		f.warehouse.storage[uint64(1)] = &products.Product{
			Id:       uint64(1),
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		}

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

		// act
		_, err := f.productRepo.GetProductById(context.Background(), 1)

		// assert
		assert.EqualError(t, err, "1: product does not exist")
	})
}

func TestCreateProduct(t *testing.T) {
	t.Run("success creating product", func(t *testing.T) {
		// arrange
		f := SetUp(t)

		expectedProduct := products.Product{
			Id:       uint64(1),
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		}

		// act
		res, err := f.productRepo.CreateProduct(context.Background(), products.Product{
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		})

		// assert
		require.NoError(t, err)
		assert.Equal(t, res, &expectedProduct)
		assert.Equal(t, f.warehouse.storage[uint64(1)], &expectedProduct)
	})

	t.Run("try to create existed product", func(t *testing.T) {
		// arrange
		f := SetUp(t)

		f.warehouse.storage[uint64(1)] = &products.Product{
			Id:       uint64(1),
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		}

		// act
		_, err := f.productRepo.CreateProduct(context.Background(), products.Product{
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		})

		// assert
		assert.EqualError(t, err, "1: product already exists")
	})

	t.Run("try to create product with ID", func(t *testing.T) {
		// arrange
		f := SetUp(t)

		// act
		_, err := f.productRepo.CreateProduct(context.Background(), products.Product{
			Id:       uint64(1),
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		})

		// assert
		assert.EqualError(t, err, "Can't create new products: Product id already set")
	})
}

func TestDeleteProduct(t *testing.T) {
	t.Run("success deleting product", func(t *testing.T) {
		// arrange
		f := SetUp(t)

		f.warehouse.storage[uint64(1)] = &products.Product{
			Id:       uint64(1),
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		}

		// act
		err := f.productRepo.DeleteProduct(context.Background(), 1)

		// assert
		require.NoError(t, err)
		_, ok := f.warehouse.storage[uint64(1)]
		assert.False(t, ok)
	})

	t.Run("product does not exists", func(t *testing.T) {
		// arrange
		f := SetUp(t)

		// act
		err := f.productRepo.DeleteProduct(context.Background(), 1)

		// assert
		assert.EqualError(t, err, "1: product does not exist")
	})
}

func TestUpdateProduct(t *testing.T) {
	t.Run("success updating product", func(t *testing.T) {
		// arrange
		f := SetUp(t)

		updatedProduct := products.Product{
			Id:       uint64(1),
			Name:     "product2",
			Price:    uint64(2),
			Quantity: uint64(2),
		}

		f.warehouse.storage[uint64(1)] = &products.Product{
			Id:       uint64(1),
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		}

		// act
		res, err := f.productRepo.UpdateProduct(context.Background(), updatedProduct)

		// assert
		require.NoError(t, err)
		assert.Equal(t, res, &updatedProduct)
		assert.Equal(t, f.warehouse.storage[uint64(1)], &updatedProduct)
	})

	t.Run("updating not existed product", func(t *testing.T) {
		// arrange
		f := SetUp(t)

		// act
		_, err := f.productRepo.UpdateProduct(context.Background(), products.Product{
			Id:       uint64(1),
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		})

		// assert
		assert.EqualError(t, err, "1: product does not exist")
	})
}

func TestGetAllProducts(t *testing.T) {
	t.Run("success getting all products", func(t *testing.T) {
		// arrange
		f := SetUp(t)

		f.warehouse.storage[uint64(1)] = &products.Product{
			Id:       uint64(1),
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		}

		f.warehouse.storage[uint64(2)] = &products.Product{
			Id:       uint64(2),
			Name:     "product2",
			Price:    uint64(2),
			Quantity: uint64(2),
		}

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
