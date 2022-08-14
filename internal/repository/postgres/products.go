package repository

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/pkg/errors"
	"homework-1/internal/models/products"
	"homework-1/internal/repository"
	"strconv"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

var defaultProductsPageSize = uint64(20)

func (r *Repository) GetProductById(ctx context.Context, id uint64) (*products.Product, error) {
	query, args, err := psql.Select("id, name, price, quantity").
		From("products").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("Repository.GetProductById: to sql: %w", err)
	}

	var product products.Product
	if err = pgxscan.Get(ctx, r.pool, &product, query, args...); err != nil {
		if pgxscan.NotFound(err) {
			return nil, errors.Wrap(repository.ProductNotExists, strconv.FormatUint(id, 10))
		}
		return nil, fmt.Errorf("Repository.GetProductById: select: %w", err)
	}

	return &product, nil
}

func (r *Repository) CreateProduct(ctx context.Context, product products.Product) (*products.Product, error) {
	query, args, err := psql.Insert("products").
		Columns("name, price, quantity").
		Values(product.Name, product.Price, product.Quantity).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("Repository.CreateProduct: to sql: %w", err)
	}

	row := r.pool.QueryRow(ctx, query, args...)
	if err = row.Scan(&product.Id); err != nil {
		return nil, fmt.Errorf("Repository.CreateProduct: insert: %w", err)
	}

	return &product, nil
}

func (r *Repository) DeleteProduct(ctx context.Context, id uint64) error {
	query, args, err := psql.Delete("products").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("Repository.DeleteProduct: to sql: %w", err)
	}

	if _, err = r.pool.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("Repository.DeleteProduct: to delete: %w", err)
	}
	return nil
}

func (r *Repository) UpdateProduct(ctx context.Context, product products.Product) (*products.Product, error) {
	query, args, err := psql.Update("products").
		Set("name", product.Name).
		Set("price", product.Price).
		Set("quantity", product.Quantity).
		Where(squirrel.Eq{"id": product.Id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("Repository.UpdateProduct: to sql: %w", err)
	}

	if _, err = r.pool.Exec(ctx, query, args...); err != nil {
		return nil, fmt.Errorf("Repository.UpdateProduct: to update: %w", err)
	}

	return &product, nil
}

func (r *Repository) GetAllProducts(ctx context.Context, page uint64, size uint64) ([]*products.Product, error) {
	limit, offset := r.getPaginationLimitAndOffset(page, size)

	query, args, err := psql.Select("id, name, price, quantity").
		From("products").
		OrderBy("id").
		Limit(limit).
		Offset(offset).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("Repository.GetAllProducts: to sql: %w", err)
	}

	var allProducts []*products.Product
	if err = pgxscan.Select(ctx, r.pool, &allProducts, query, args...); err != nil {
		return nil, fmt.Errorf("Repository.GetAllProducts: select: %w", err)
	}

	return allProducts, nil
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
