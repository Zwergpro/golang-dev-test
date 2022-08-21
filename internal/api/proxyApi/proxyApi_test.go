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
				Name:     "product1",
				Price:    uint64(1),
				Quantity: uint64(1),
			}, nil)

		// act
		res, err := f.service.ProductGet(context.Background(), &pbApi.ProductGetRequest{Id: uint64(1)})

		// assert
		require.NoError(t, err)
		assert.Equal(t, res, &pbApi.ProductGetResponse{
			Id:       uint64(1),
			Name:     "product1",
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

func TestProductCreate(t *testing.T) {
	t.Run("success creating", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		f.storageClient.EXPECT().ProductCreate(gomock.Any(), &pbStorage.ProductCreateRequest{
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		}).Return(&pbStorage.ProductCreateResponse{
			Id:       uint64(1),
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		}, nil)

		// act
		res, err := f.service.ProductCreate(context.Background(), &pbApi.ProductCreateRequest{
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		})

		// assert
		require.NoError(t, err)
		assert.Equal(t, res, &pbApi.ProductCreateResponse{
			Id:       uint64(1),
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		})
	})

	t.Run("StorageClient fail", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		f.storageClient.EXPECT().ProductCreate(gomock.Any(), &pbStorage.ProductCreateRequest{
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		}).Return(nil, status.Error(codes.Internal, "internal error"))

		// act
		_, err := f.service.ProductCreate(context.Background(), &pbApi.ProductCreateRequest{
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(1),
		})

		// assert
		assert.EqualError(t, err, "rpc error: code = Internal desc = internal error")
	})

	t.Run("fail with wrong name", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		// act
		_, err := f.service.ProductCreate(context.Background(), &pbApi.ProductCreateRequest{
			Name:     "",
			Price:    uint64(1),
			Quantity: uint64(1),
		})

		// assert
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = name length must be greater than 0")
	})

	t.Run("fail with wrong price", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		// act
		_, err := f.service.ProductCreate(context.Background(), &pbApi.ProductCreateRequest{
			Name:     "product1",
			Price:    uint64(0),
			Quantity: uint64(1),
		})

		// assert
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = price must be greater than 0")
	})

	t.Run("fail with wrong quantity", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		// act
		_, err := f.service.ProductCreate(context.Background(), &pbApi.ProductCreateRequest{
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(0),
		})

		// assert
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = quantity must be greater than 0")
	})

	t.Run("fail with wrong args", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		// act
		_, err := f.service.ProductCreate(context.Background(), &pbApi.ProductCreateRequest{
			Name:     "",
			Price:    uint64(0),
			Quantity: uint64(0),
		})

		// assert
		expectedErr := "rpc error: code = InvalidArgument desc = name length must be greater than 0; price must be greater than 0; quantity must be greater than 0"
		assert.EqualError(t, err, expectedErr)
	})
}

func TestProductUpdate(t *testing.T) {
	t.Run("success updating", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		f.storageClient.EXPECT().ProductUpdate(gomock.Any(), &pbStorage.ProductUpdateRequest{
			Id:       uint64(1),
			Name:     "product2",
			Price:    uint64(2),
			Quantity: uint64(2),
		}).Return(&pbStorage.ProductUpdateResponse{
			Id:       uint64(1),
			Name:     "product2",
			Price:    uint64(2),
			Quantity: uint64(2),
		}, nil)

		// act
		res, err := f.service.ProductUpdate(context.Background(), &pbApi.ProductUpdateRequest{
			Id:       uint64(1),
			Name:     "product2",
			Price:    uint64(2),
			Quantity: uint64(2),
		})

		// assert
		require.NoError(t, err)
		assert.Equal(t, res, &pbApi.ProductUpdateResponse{
			Id:       uint64(1),
			Name:     "product2",
			Price:    uint64(2),
			Quantity: uint64(2),
		})
	})

	t.Run("product does not exist", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		f.storageClient.EXPECT().ProductUpdate(gomock.Any(), &pbStorage.ProductUpdateRequest{
			Id:       uint64(1),
			Name:     "product2",
			Price:    uint64(2),
			Quantity: uint64(2),
		}).Return(nil, status.Error(codes.NotFound, "product not found"))

		// act
		_, err := f.service.ProductUpdate(context.Background(), &pbApi.ProductUpdateRequest{
			Id:       uint64(1),
			Name:     "product2",
			Price:    uint64(2),
			Quantity: uint64(2),
		})

		// assert
		assert.EqualError(t, err, "rpc error: code = NotFound desc = product not found")
	})

	t.Run("storageClient internal error", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		f.storageClient.EXPECT().ProductUpdate(gomock.Any(), &pbStorage.ProductUpdateRequest{
			Id:       uint64(1),
			Name:     "product2",
			Price:    uint64(2),
			Quantity: uint64(2),
		}).Return(nil, status.Error(codes.Internal, "internal error"))

		// act
		_, err := f.service.ProductUpdate(context.Background(), &pbApi.ProductUpdateRequest{
			Id:       uint64(1),
			Name:     "product2",
			Price:    uint64(2),
			Quantity: uint64(2),
		})

		// assert
		assert.EqualError(t, err, "rpc error: code = Internal desc = internal error")
	})

	t.Run("fail with wrong name", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		// act
		_, err := f.service.ProductUpdate(context.Background(), &pbApi.ProductUpdateRequest{
			Id:       uint64(1),
			Name:     "",
			Price:    uint64(1),
			Quantity: uint64(1),
		})

		// assert
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = name length must be greater than 0")
	})

	t.Run("fail with wrong price", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		// act
		_, err := f.service.ProductUpdate(context.Background(), &pbApi.ProductUpdateRequest{
			Id:       uint64(1),
			Name:     "product1",
			Price:    uint64(0),
			Quantity: uint64(1),
		})

		// assert
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = price must be greater than 0")
	})

	t.Run("fail with wrong quantity", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		// act
		_, err := f.service.ProductUpdate(context.Background(), &pbApi.ProductUpdateRequest{
			Id:       uint64(1),
			Name:     "product1",
			Price:    uint64(1),
			Quantity: uint64(0),
		})

		// assert
		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = quantity must be greater than 0")
	})

	t.Run("fail with wrong args", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		// act
		_, err := f.service.ProductUpdate(context.Background(), &pbApi.ProductUpdateRequest{
			Id:       uint64(1),
			Name:     "",
			Price:    uint64(0),
			Quantity: uint64(0),
		})

		// assert
		expectedErr := "rpc error: code = InvalidArgument desc = name length must be greater than 0; price must be greater than 0; quantity must be greater than 0"
		assert.EqualError(t, err, expectedErr)
	})
}

func TestProductDelete(t *testing.T) {
	t.Run("success deleting", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		f.storageClient.EXPECT().ProductDelete(gomock.Any(), &pbStorage.ProductDeleteRequest{
			Id: uint64(1),
		}).Return(&pbStorage.ProductDeleteResponse{}, nil)

		// act
		res, err := f.service.ProductDelete(context.Background(), &pbApi.ProductDeleteRequest{
			Id: uint64(1),
		})

		// assert
		require.NoError(t, err)
		assert.Equal(t, res, &pbApi.ProductDeleteResponse{})
	})

	t.Run("product does not exist", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		f.storageClient.EXPECT().ProductDelete(gomock.Any(), &pbStorage.ProductDeleteRequest{
			Id: uint64(1),
		}).Return(nil, status.Error(codes.NotFound, "product not found"))

		// act
		_, err := f.service.ProductDelete(context.Background(), &pbApi.ProductDeleteRequest{
			Id: uint64(1),
		})

		// assert
		assert.EqualError(t, err, "rpc error: code = NotFound desc = product not found")
	})

	t.Run("storageClient internal error", func(t *testing.T) {
		// arrange
		f := SetUp(t)
		defer f.TearDown()

		f.storageClient.EXPECT().ProductDelete(gomock.Any(), &pbStorage.ProductDeleteRequest{
			Id: uint64(1),
		}).Return(nil, status.Error(codes.Internal, "internal error"))

		// act
		_, err := f.service.ProductDelete(context.Background(), &pbApi.ProductDeleteRequest{
			Id: uint64(1),
		})

		// assert
		assert.EqualError(t, err, "rpc error: code = Internal desc = internal error")
	})
}
