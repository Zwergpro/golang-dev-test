package handlers

import (
	"fmt"
	"github.com/pkg/errors"
	"homework-1/internal/storage"
	"strconv"
	"strings"
)

type addParams struct {
	name     string
	price    uint
	quantity uint
}

func parseArgsToParams(args string) (*addParams, error) {
	params := strings.Split(args, " ")
	if len(params) != 3 {
		return nil, errors.Wrapf(BadArguments, "Invalid arguments count: %d", len(params))
	}

	name := params[0]

	price, err := strconv.ParseUint(params[1], 10, 64)
	if err != nil {
		return nil, errors.Wrapf(BadArguments, "Can't parse price: %s", params[1])
	}

	quantity, err := strconv.ParseUint(params[2], 10, 64)
	if err != nil {
		return nil, errors.Wrapf(BadArguments, "Can't parse quantity: %s", params[2])
	}

	return &addParams{
		name:     name,
		price:    uint(price),
		quantity: uint(quantity),
	}, nil

}

func addCmdHandler(cmdArgs string) string {
	params, err := parseArgsToParams(cmdArgs)
	if err != nil {
		return err.Error()
	}

	product, err := storage.NewProduct(params.name, params.price, params.quantity)
	if err != nil {
		return err.Error()
	}

	if err := storage.Add(product); err != nil {
		return err.Error()
	}

	return fmt.Sprintf("Product added: %s", product.String())
}
