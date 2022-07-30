package repository

import (
	"context"
	"github.com/pkg/errors"
	"homework-1/internal/models"
	"homework-1/internal/repository"
	"strconv"
)

var ErrProductIdAlreadySet = errors.New("Product id already set")

func (r *Repository) GetProductById(ctx context.Context, id uint64) (*models.Product, error) {
	if err := r.warehouse.RLockWithContext(ctx); err != nil {
		return nil, err
	}
	defer r.warehouse.RUnlock()

	if product, ok := r.warehouse.storage[id]; ok {
		return product.Copy(), nil
	}
	return nil, errors.Wrap(repository.ProductNotExists, strconv.FormatUint(id, 10))
}

func (r *Repository) CreateProduct(ctx context.Context, product models.Product) (*models.Product, error) {
	if product.Id > 0 {
		return nil, errors.Wrap(ErrProductIdAlreadySet, "Can't create new product")
	}

	product.Id = r.warehouse.GetNextId()

	if err := r.warehouse.LockWithContext(ctx); err != nil {
		return nil, err
	}
	defer r.warehouse.Unlock()

	if _, ok := r.warehouse.storage[product.GetId()]; ok {
		return nil, errors.Wrap(repository.ProductAlreadyExists, strconv.FormatUint(product.GetId(), 10))
	}
	r.warehouse.storage[product.GetId()] = &product
	return product.Copy(), nil
}

func (r *Repository) DeleteProduct(ctx context.Context, id uint64) error {
	if err := r.warehouse.LockWithContext(ctx); err != nil {
		return err
	}
	defer r.warehouse.Unlock()

	if _, ok := r.warehouse.storage[id]; !ok {
		return errors.Wrap(repository.ProductNotExists, strconv.FormatUint(id, 10))
	}
	delete(r.warehouse.storage, id)
	return nil
}

func (r *Repository) UpdateProduct(ctx context.Context, product models.Product) (*models.Product, error) {
	if err := r.warehouse.LockWithContext(ctx); err != nil {
		return nil, err
	}
	defer r.warehouse.mu.Unlock()

	if _, ok := r.warehouse.storage[product.GetId()]; !ok {
		return nil, errors.Wrap(repository.ProductNotExists, strconv.FormatUint(product.GetId(), 10))
	}
	r.warehouse.storage[product.GetId()] = &product
	return product.Copy(), nil
}

func (r *Repository) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	if err := r.warehouse.RLockWithContext(ctx); err != nil {
		return nil, err
	}
	defer r.warehouse.RUnlock()

	products := make([]*models.Product, 0, len(r.warehouse.storage))
	for _, v := range r.warehouse.storage {
		products = append(products, v.Copy())
	}
	return products, nil
}
