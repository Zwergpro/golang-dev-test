package repository

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"homework-1/internal/models"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func (r *Repository) GetProductById(ctx context.Context, id uint64) (*models.Product, error) {
	query, args, err := psql.Select("id, name, price, quantity").
		From("products").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("Repository.GetProductById: to sql: %w", err)
	}

	var product models.Product
	if err := pgxscan.Get(ctx, r.pool, &product, query, args...); err != nil {
		return nil, fmt.Errorf("Repository.GetProductById: select: %w", err)
	}

	return &product, nil
}

func (r *Repository) CreateProduct(ctx context.Context, product models.Product) (*models.Product, error) {
	query, args, err := psql.Insert("products").
		Columns("name, price, quantity").
		Values(product.Name, product.Price, product.Quantity).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("Repository.CreateProduct: to sql: %w", err)
	}

	row := r.pool.QueryRow(ctx, query, args...)
	if err := row.Scan(&product.Id); err != nil {
		return nil, fmt.Errorf("Repository.CreateProduct: insert: %w", err)
	}

	return &product, nil
}

func (r *Repository) DeleteProduct(ctx context.Context, id uint64) error {
	query, args, err := psql.Delete("products").Where(squirrel.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("Repository.DeleteProduct: to sql: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("Repository.DeleteProduct: to delete: %w", err)
	}
	return nil
}

func (r *Repository) UpdateProduct(ctx context.Context, product models.Product) (*models.Product, error) {
	query, args, err := psql.Update("products").
		Set("name", product.Name).
		Set("price", product.Price).
		Set("quantity", product.Quantity).
		Where(squirrel.Eq{"id": product.Id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("Repository.UpdateProduct: to sql: %w", err)
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("Repository.UpdateProduct: to update: %w", err)
	}

	return &product, nil
}

func (r *Repository) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	query, args, err := psql.Select("id, name, price, quantity").
		From("products").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("Repository.GetAllProducts: to sql: %w", err)
	}

	var products []*models.Product
	if err := pgxscan.Select(ctx, r.pool, &products, query, args...); err != nil {
		return nil, fmt.Errorf("Repository.GetAllProducts: select: %w", err)
	}

	return products, nil
}
