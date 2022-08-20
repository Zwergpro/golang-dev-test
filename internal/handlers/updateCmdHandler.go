package handlers

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"homework-1/internal/models/products"
	"homework-1/internal/repository"
	"strconv"
	"strings"
)

func updateCmdHandler(repository repository.Product, cmdArgs string) string {
	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	args := strings.Split(cmdArgs, " ")
	if len(args) != 4 {
		return errors.Wrapf(BadArguments, "Invalid arguments count: %d", len(args)).Error()
	}

	product, err := getProductByStringId(ctx, repository, args[0])
	if err != nil {
		return err.Error()
	}

	product, err = updateProduct(product, args[1:])
	if err != nil {
		return err.Error()
	}

	product, err = repository.UpdateProduct(ctx, *product)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("Product updated: %s", product.String())
}

func getProductByStringId(ctx context.Context, repository repository.Product, id string) (*products.Product, error) {
	productId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(BadArguments, "Can't parse id: %s", id)
	}

	product, err := repository.GetProductById(ctx, productId)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func updateProduct(product *products.Product, params []string) (*products.Product, error) {
	if err := product.SetName(params[0]); err != nil {
		return nil, err
	}

	price, err := strconv.ParseUint(params[1], 10, 64)
	if err != nil {
		return nil, errors.Wrapf(BadArguments, "Can't parse price: %s", params[1])
	}
	if err = product.SetPrice(price); err != nil {
		return nil, err
	}

	quantity, err := strconv.ParseUint(params[2], 10, 64)
	if err != nil {
		return nil, errors.Wrapf(BadArguments, "Can't parse quantity: %s", params[2])
	}
	if err = product.SetQuantity(quantity); err != nil {
		return nil, err
	}

	return product, nil
}
