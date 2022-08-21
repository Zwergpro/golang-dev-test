package proxyApi

import (
	"github.com/golang/mock/gomock"
	mock_storage "homework-1/internal/api/proxyApi/mock"
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
