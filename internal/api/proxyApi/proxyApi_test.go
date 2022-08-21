package proxyApi

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pbStorage "homework-1/pkg/api/storage/v1"
	pbApi "homework-1/pkg/api/v1"
	"testing"
)

func TestProductGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		f.storageClient.EXPECT().ProductGet(gomock.Any(), &pbStorage.ProductGetRequest{Id: uint64(1)}).
			Return(&pbStorage.ProductGetResponse{
				Id:       uint64(1),
				Name:     "product 1",
				Price:    uint64(1),
				Quantity: uint64(1),
			}, nil)

		// act
		res, err := f.service.ProductGet(context.Background(), &pbApi.ProductGetRequest{Id: uint64(1)})

		// assert
		require.NoError(t, err)
		assert.Equal(t, res, &pbApi.ProductGetResponse{
			Id:       uint64(1),
			Name:     "product 1",
			Price:    uint64(1),
			Quantity: uint64(1),
		})
	})

	t.Run("not found error", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		f.storageClient.EXPECT().ProductGet(gomock.Any(), &pbStorage.ProductGetRequest{Id: uint64(1)}).
			Return(nil, status.Error(codes.NotFound, "not found"))

		// act
		_, err := f.service.ProductGet(context.Background(), &pbApi.ProductGetRequest{Id: uint64(1)})

		// assert
		assert.EqualError(t, err, "rpc error: code = NotFound desc = product not found")
	})

	t.Run("internal error", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		f.storageClient.EXPECT().ProductGet(gomock.Any(), &pbStorage.ProductGetRequest{Id: uint64(1)}).
			Return(nil, errors.New("internal error"))

		// act
		_, err := f.service.ProductGet(context.Background(), &pbApi.ProductGetRequest{Id: uint64(1)})

		// assert
		assert.EqualError(t, err, "rpc error: code = Internal desc = internal error")
	})
}
