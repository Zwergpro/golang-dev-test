package handlers

import (
	"fmt"
	"github.com/pkg/errors"
	"homework-1/internal/storage"
	"strconv"
	"strings"
)

func updateCmdHandler(cmdArgs string) string {
	args := strings.Split(cmdArgs, " ")
	if len(args) != 4 {
		return errors.Wrapf(BadArguments, "Invalid arguments count: %d", len(args)).Error()
	}

	oldProduct, err := getProductByStringId(args[0])
	if err != nil {
		return err.Error()
	}

	// Pass by value to escape partially update when error occurred
	newProduct, err := updateProduct(*oldProduct, args[1:])
	if err != nil {
		return err.Error()
	}

	if err = storage.Update(newProduct); err != nil {
		return err.Error()
	}
	return fmt.Sprintf("Product updated: %s", newProduct.String())
}

func getProductByStringId(id string) (*storage.Product, error) {
	productId, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, errors.Wrapf(BadArguments, "Can't parse id: %s", id)
	}

	product, err := storage.Get(uint(productId))
	if err != nil {
		return nil, err
	}

	return product, nil
}

func updateProduct(product storage.Product, params []string) (*storage.Product, error) {
	if err := product.SetName(params[0]); err != nil {
		return nil, err
	}

	price, err := strconv.ParseUint(params[1], 10, 64)
	if err != nil {
		return nil, errors.Wrapf(BadArguments, "Can't parse price: %s", params[1])
	}
	if err = product.SetPrice(uint(price)); err != nil {
		return nil, err
	}

	quantity, err := strconv.ParseUint(params[2], 10, 64)
	if err != nil {
		return nil, errors.Wrapf(BadArguments, "Can't parse quantity: %s", params[2])
	}
	if err = product.SetQuantity(uint(quantity)); err != nil {
		return nil, err
	}

	return &product, nil
}
