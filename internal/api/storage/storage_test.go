package storage

import (
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"homework-1/internal/models/products"
	pb "homework-1/pkg/api/storage/v1"
	"testing"
)

func TestProductList(t *testing.T) {
	t.Run("success getting product list", func(t *testing.T) {
		// arrange
		f := SetUp(t)

		stream := makeProductListResponseStreamMock()

		f.productRepo.EXPECT().GetAllProducts(gomock.Any(), uint64(0), uint64(0)).Return([]*products.Product{
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
		}, nil)

		// act
		err := f.service.ProductList(&pb.ProductListRequest{}, stream)

		res := stream.GetAll()

		// assert
		require.NoError(t, err)
		assert.Equal(t, res, []*pb.ProductListResponse{
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

	t.Run("fail with internal error", func(t *testing.T) {
		// arrange
		f := SetUp(t)

		stream := makeProductListResponseStreamMock()
		defer stream.Close()

		f.productRepo.EXPECT().GetAllProducts(gomock.Any(), uint64(0), uint64(0)).Return(nil, errors.New("internal error"))

		// act
		err := f.service.ProductList(&pb.ProductListRequest{}, stream)

		// assert
		assert.EqualError(t, err, "rpc error: code = Internal desc = internal error")
	})
}
