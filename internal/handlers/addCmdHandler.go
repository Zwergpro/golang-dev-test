package handlers

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"homework-1/internal/models"
	"homework-1/internal/repository"
	"strconv"
	"strings"
)

func addCmdHandler(repository repository.Product, cmdArgs string) string {
	params := strings.Split(cmdArgs, " ")
	if len(params) != 3 {
		return errors.Wrapf(BadArguments, "Invalid arguments count: %d", len(params)).Error()
	}

	price, err := strconv.ParseUint(params[1], 10, 64)
	if err != nil {
		return errors.Wrapf(BadArguments, "Can't parse price: %s", params[1]).Error()
	}

	quantity, err := strconv.ParseUint(params[2], 10, 64)
	if err != nil {
		return errors.Wrapf(BadArguments, "Can't parse quantity: %s", params[2]).Error()
	}

	product, err := models.BuildProduct(params[0], price, quantity)
	if err != nil {
		return err.Error()
	}

	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	product, err = repository.CreateProduct(ctx, *product)
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("Product added: %s", product.String())
}
