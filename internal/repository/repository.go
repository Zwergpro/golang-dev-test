//go:generate mockgen -source ./repository.go -destination=./mock/repository.go -package=mock_repository

package repository

import (
	"context"
	"homework-1/internal/models/products"
)

type Product interface {
	GetProductById(ctx context.Context, id uint64) (*products.Product, error)
	GetAllProducts(ctx context.Context, page uint64, size uint64) ([]*products.Product, error)
	CreateProduct(ctx context.Context, product products.Product) (*products.Product, error)
	UpdateProduct(ctx context.Context, product products.Product) (*products.Product, error)
	DeleteProduct(ctx context.Context, id uint64) error
}
