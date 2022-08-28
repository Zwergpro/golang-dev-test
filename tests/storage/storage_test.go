//go:build integration
// +build integration

package storage

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"homework-1/internal/models/products"
	pbStorage "homework-1/pkg/api/storage/v1"
	"io"
	"testing"
)

func TestTestProductListStorage(t *testing.T) {
	t.Run("success getting empty result", func(t *testing.T) {
		//arrange
		DB.SetUp(t)
		defer DB.TearDown()

		pageNum := uint64(1)
		pageSize := uint64(10)

		//act
		stream, err := StorageClient.ProductList(context.Background(), &pbStorage.ProductListRequest{Page: &pageNum, Size: &pageSize})

		_, err = stream.Recv()

		// assert
		assert.Equal(t, err, io.EOF)
	})

	t.Run("success getting result", func(t *testing.T) {
		//arrange
		DB.SetUp(t)
		defer DB.TearDown()

		pageNum := uint64(1)
		pageSize := uint64(10)

		newProduct, err := StorageClient.ProductCreate(context.Background(), &pbStorage.ProductCreateRequest{
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		})
		if err != nil {
			t.FailNow()
		}

		//act
		stream, err := StorageClient.ProductList(context.Background(), &pbStorage.ProductListRequest{Page: &pageNum, Size: &pageSize})

		var result []*products.Product
		for {
			product, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.FailNow()
			}
			result = append(result, &products.Product{
				Id:       product.GetId(),
				Name:     product.GetName(),
				Price:    product.GetPrice(),
				Quantity: product.GetQuantity(),
			})
		}

		//assert
		require.NoError(t, err)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, result[0], &products.Product{
			Id:       newProduct.Id,
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		})
	})
}

func TestProductGetProxyApi(t *testing.T) {
	t.Run("success getting result", func(t *testing.T) {
		//arrange
		DB.SetUp(t)
		defer DB.TearDown()

		newProduct, err := StorageClient.ProductCreate(context.Background(), &pbStorage.ProductCreateRequest{
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		})
		if err != nil {
			t.FailNow()
		}

		//act
		resp, err := StorageClient.ProductGet(context.Background(), &pbStorage.ProductGetRequest{Id: newProduct.Id})

		//assert
		require.NoError(t, err)
		assert.Equal(t, resp.Id, newProduct.Id)
		assert.Equal(t, resp.Name, "product1")
		assert.Equal(t, resp.Price, uint64(1))
		assert.Equal(t, resp.Quantity, uint64(1))
	})

	t.Run("product does not exists", func(t *testing.T) {
		//arrange
		DB.SetUp(t)
		defer DB.TearDown()

		//act
		_, err := StorageClient.ProductGet(context.Background(), &pbStorage.ProductGetRequest{Id: 10})

		//assert
		assert.EqualError(t, err, "rpc error: code = NotFound desc = 10: product does not exist")
	})
}

func TestProductCreateProxyApi(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		//arrange
		DB.SetUp(t)
		defer DB.TearDown()

		//act
		resp, err := StorageClient.ProductCreate(context.Background(), &pbStorage.ProductCreateRequest{
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		})

		//assert
		require.NoError(t, err)
		assert.Equal(t, resp.Name, "product1")
		assert.Equal(t, resp.Price, uint64(1))
		assert.Equal(t, resp.Quantity, uint64(1))
	})
}

func TestProductUpdateProxyApi(t *testing.T) {
	t.Run("success updating", func(t *testing.T) {
		//arrange
		DB.SetUp(t)
		defer DB.TearDown()

		newProduct, err := StorageClient.ProductCreate(context.Background(), &pbStorage.ProductCreateRequest{
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		})
		if err != nil {
			t.FailNow()
		}

		//act
		resp, err := StorageClient.ProductUpdate(context.Background(), &pbStorage.ProductUpdateRequest{
			Id:       newProduct.Id,
			Name:     "product2",
			Price:    uint64(2),
			Quantity: uint64(2),
		})

		//assert
		require.NoError(t, err)
		assert.Equal(t, resp.Id, newProduct.Id)
		assert.Equal(t, resp.Name, "product2")
		assert.Equal(t, resp.Price, uint64(2))
		assert.Equal(t, resp.Quantity, uint64(2))
	})

	t.Run("product does not exist", func(t *testing.T) {
		//arrange
		DB.SetUp(t)
		defer DB.TearDown()

		//act
		_, err := StorageClient.ProductUpdate(context.Background(), &pbStorage.ProductUpdateRequest{
			Id:       10,
			Name:     "product2",
			Price:    uint64(2),
			Quantity: uint64(2),
		})

		//assert
		assert.EqualError(t, err, "rpc error: code = NotFound desc = 10: product does not exist")
	})
}

func TestProductDeleteProxyApi(t *testing.T) {
	t.Run("success deleting", func(t *testing.T) {
		//arrange
		DB.SetUp(t)
		defer DB.TearDown()

		newProduct, err := StorageClient.ProductCreate(context.Background(), &pbStorage.ProductCreateRequest{
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		})
		if err != nil {
			t.FailNow()
		}

		//act
		_, err = StorageClient.ProductDelete(context.Background(), &pbStorage.ProductDeleteRequest{
			Id: newProduct.Id,
		})

		//assert
		require.NoError(t, err)
	})
}
