package proxyApi

import (
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
	mock_storage "homework-1/internal/api/proxyApi/mock"
	pbStorage "homework-1/pkg/api/storage/v1"
	"io"
	"testing"
)

type proxyApiFixture struct {
	service       *implementation
	ctrl          *gomock.Controller
	storageClient *mock_storage.MockStorageServiceClient
}

func SetUp(t *testing.T) *proxyApiFixture {
	f := proxyApiFixture{ctrl: gomock.NewController(t)}
	f.storageClient = mock_storage.NewMockStorageServiceClient(f.ctrl)
	f.service = New(Deps{StorageClient: f.storageClient})
	return &f
}

func (p *proxyApiFixture) TearDown() {
	p.ctrl.Finish()
}

func makeProductListResponseStreamMock() *ProductListResponseStreamMock {
	return &ProductListResponseStreamMock{
		queue: make(chan *pbStorage.ProductListResponse, 10),
	}
}

type ProductListResponseStreamMock struct {
	grpc.ClientStream
	queue chan *pbStorage.ProductListResponse
}

func (m *ProductListResponseStreamMock) Close() {
	close(m.queue)
}

func (m *ProductListResponseStreamMock) Send(resp *pbStorage.ProductListResponse) {
	m.queue <- resp
}

func (m *ProductListResponseStreamMock) Recv() (*pbStorage.ProductListResponse, error) {
	if len(m.queue) == 0 {
		return nil, io.EOF
	}
	return <-m.queue, nil
}
