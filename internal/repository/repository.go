package repository

import (
	"context"
	"homework-1/internal/models"
)

type Product interface {
	GetProductById(ctx context.Context, id uint64) (*models.Product, error)
	GetAllProducts(ctx context.Context) ([]*models.Product, error)
	CreateProduct(ctx context.Context, product models.Product) (*models.Product, error)
	UpdateProduct(ctx context.Context, product models.Product) (*models.Product, error)
	DeleteProduct(ctx context.Context, id uint64) error
}
