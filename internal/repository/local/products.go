package repository

import (
	"context"
	"github.com/pkg/errors"
	"homework-1/internal/math"
	"homework-1/internal/models/products"
	"homework-1/internal/repository"
	"sort"
	"strconv"
)

var ErrProductIdAlreadySet = errors.New("Product id already set")

var defaultProductsPageSize = uint64(20)

func (r *Repository) GetProductById(ctx context.Context, id uint64) (*products.Product, error) {
	if err := r.warehouse.RLockWithContext(ctx); err != nil {
		return nil, err
	}
	defer r.warehouse.RUnlock()

	if product, ok := r.warehouse.storage[id]; ok {
		return product.Copy(), nil
	}
	return nil, errors.Wrap(repository.ProductNotExists, strconv.FormatUint(id, 10))
}

func (r *Repository) CreateProduct(ctx context.Context, product products.Product) (*products.Product, error) {
	if product.Id > 0 {
		return nil, errors.Wrap(ErrProductIdAlreadySet, "Can't create new products")
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

func (r *Repository) UpdateProduct(ctx context.Context, product products.Product) (*products.Product, error) {
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

func (r *Repository) GetAllProducts(ctx context.Context, page uint64, size uint64) ([]*products.Product, error) {
	limit, offset := r.getPaginationLimitAndOffset(page, size)

	if err := r.warehouse.RLockWithContext(ctx); err != nil {
		return nil, err
	}
	defer r.warehouse.RUnlock()

	warehouseLen := uint64(len(r.warehouse.storage))

	allProducts := make([]*products.Product, 0, warehouseLen)
	for _, v := range r.warehouse.storage {
		allProducts = append(allProducts, v.Copy())
	}
	sort.SliceStable(allProducts, func(i, j int) bool {
		return allProducts[i].Id < allProducts[j].Id
	})

	start := math.MinUint64(warehouseLen, offset)
	end := math.MinUint64(warehouseLen, offset+limit)
	return allProducts[start:end], nil
}

func (r *Repository) getPaginationLimitAndOffset(page uint64, size uint64) (uint64, uint64) {
	if page <= 0 {
		page = 1 // min page number
	}

	if size <= 0 {
		size = defaultProductsPageSize
	}

	offset := (page - 1) * size // first page does not have offset
	return size, offset
}
