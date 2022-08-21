package storage

import (
	"context"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
	mock_repository "homework-1/internal/repository/mock"
	pb "homework-1/pkg/api/storage/v1"
	"testing"
)

type storageFixture struct {
	Ctx         context.Context
	service     *implementation
	productRepo *mock_repository.MockProduct
}

func SetUp(t *testing.T) *storageFixture {
	f := storageFixture{Ctx: context.Background()}
	f.productRepo = mock_repository.NewMockProduct(gomock.NewController(t))
	f.service = New(Deps{ProductRepository: f.productRepo})
	return &f
}

func makeProductListResponseStreamMock() *ProductListResponseStreamMock {
	return &ProductListResponseStreamMock{
		queue: make(chan *pb.ProductListResponse, 10),
	}
}

type ProductListResponseStreamMock struct {
	grpc.ServerStream
	queue chan *pb.ProductListResponse
}

func (m *ProductListResponseStreamMock) Close() {
	close(m.queue)
}

func (m *ProductListResponseStreamMock) Send(resp *pb.ProductListResponse) error {
	m.queue <- resp
	return nil
}

func (m *ProductListResponseStreamMock) GetAll() []*pb.ProductListResponse {
	m.Close()

	var resp []*pb.ProductListResponse
	for product := range m.queue {
		resp = append(resp, product)
	}
	return resp
}
