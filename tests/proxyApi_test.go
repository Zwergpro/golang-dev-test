//go:build integration
// +build integration

package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pbApi "homework-1/pkg/api/v1"
	"testing"
)

func TestTestProductListProxyApi(t *testing.T) {
	t.Run("success getting empty result", func(t *testing.T) {
		//arrange
		DB.SetUp(t)
		defer DB.TearDown()

		pageNum := uint64(1)
		pageSize := uint64(10)

		//act
		resp, err := ProxyApiClient.ProductList(context.Background(), &pbApi.ProductListRequest{Page: &pageNum, Size: &pageSize})

		//assert
		require.NoError(t, err)
		assert.Equal(t, 0, len(resp.Products))
	})
	t.Run("success getting result", func(t *testing.T) {
		//arrange
		DB.SetUp(t)
		defer DB.TearDown()

		pageNum := uint64(1)
		pageSize := uint64(10)

		newProduct, err := ProxyApiClient.ProductCreate(context.Background(), &pbApi.ProductCreateRequest{
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		})
		if err != nil {
			t.Fail()
		}

		//act
		resp, err := ProxyApiClient.ProductList(context.Background(), &pbApi.ProductListRequest{Page: &pageNum, Size: &pageSize})

		//assert
		require.NoError(t, err)
		assert.Equal(t, 1, len(resp.Products))
		assert.Equal(t, &pbApi.ProductListResponse_Product{
			Id:       newProduct.Id,
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		}, resp.Products[0])
	})
}

func TestProductCreateProxyApi(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		//arrange
		DB.SetUp(t)
		defer DB.TearDown()

		//act
		resp, err := ProxyApiClient.ProductCreate(context.Background(), &pbApi.ProductCreateRequest{
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

	t.Run("fail with wrong name", func(t *testing.T) {
		//arrange
		DB.SetUp(t)
		defer DB.TearDown()

		//act
		_, err := ProxyApiClient.ProductCreate(context.Background(), &pbApi.ProductCreateRequest{
			Name:     "",
			Price:    uint64(1),
			Quantity: uint64(1),
		})

		//assert
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = name length must be greater than 0")
	})

	t.Run("fail with wrong price", func(t *testing.T) {
		//arrange
		DB.SetUp(t)
		defer DB.TearDown()

		//act
		_, err := ProxyApiClient.ProductCreate(context.Background(), &pbApi.ProductCreateRequest{
			Name:     "product1",
			Price:    uint64(0),
			Quantity: uint64(1),
		})

		//assert
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = price must be greater than 0")
	})

	t.Run("fail with wrong quantity", func(t *testing.T) {
		//arrange
		DB.SetUp(t)
		defer DB.TearDown()

		//act
		_, err := ProxyApiClient.ProductCreate(context.Background(), &pbApi.ProductCreateRequest{
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(0),
		})

		//assert
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = quantity must be greater than 0")
	})

	t.Run("fail with wrong args", func(t *testing.T) {
		//arrange
		DB.SetUp(t)
		defer DB.TearDown()

		//act
		_, err := ProxyApiClient.ProductCreate(context.Background(), &pbApi.ProductCreateRequest{
			Name:     "",
			Price:    uint64(0),
			Quantity: uint64(0),
		})

		//assert
		expectedErr := "rpc error: code = InvalidArgument desc = name length must be greater than 0; price must be greater than 0; quantity must be greater than 0"
		assert.EqualError(t, err, expectedErr)
	})
}
