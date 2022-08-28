package repository

import (
	"homework-1/internal/repository"
	"testing"
)

type productRepoFixture struct {
	productRepo repository.Product
	warehouse   *Warehouse
}

func SetUp(_ *testing.T) *productRepoFixture {
	var fixture productRepoFixture

	fixture.warehouse = NewWarehouse()
	fixture.productRepo = NewRepository(fixture.warehouse)

	return &fixture
}
