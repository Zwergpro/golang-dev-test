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
