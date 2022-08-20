package repository

import (
	"github.com/golang/mock/gomock"
	"github.com/pashagolub/pgxmock"
	"homework-1/internal/repository"
	"testing"
)

type productRepoFixture struct {
	productRepo repository.Product
	mockPool    pgxmock.PgxPoolIface
	ctrl        *gomock.Controller
}

func SetUp(t *testing.T) *productRepoFixture {
	var fixture productRepoFixture
	fixture.ctrl = gomock.NewController(t)

	mock, err := pgxmock.NewPool()
	if err != nil {
		panic(err)
	}

	fixture.mockPool = mock
	fixture.productRepo = NewRepository(mock)

	return &fixture
}

func (f *productRepoFixture) TearDown() {
	f.ctrl.Finish()
	f.mockPool.Close()
}
