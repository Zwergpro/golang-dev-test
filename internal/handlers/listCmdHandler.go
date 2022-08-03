package handlers

import (
	"context"
	"github.com/pkg/errors"
	"homework-1/internal/repository"
	"strconv"
	"strings"
)

func listCmdHandler(repository repository.Product, args string) string {
	ctx, cancel := context.WithTimeout(context.Background(), maxTimeout)
	defer cancel()

	page, size, err := extractPageAndSize(args)
	if err != nil {
		return err.Error()
	}

	products, err := repository.GetAllProducts(ctx, page, size)
	if err != nil {
		return err.Error()
	}

	if len(products) == 0 {
		return "nothing found"
	}

	res := make([]string, len(products))
	for _, p := range products {
		res = append(res, p.String())
	}

	return strings.Join(res, "\n")
}

func extractPageAndSize(args string) (uint64, uint64, error) {
	params := strings.Split(args, " ")
	if params[0] == "" {
		return 0, 0, nil
	}

	page, err := strconv.ParseUint(params[0], 10, 64)
	if err != nil {
		return 0, 0, errors.Wrapf(BadArguments, "Can't parse page number: %s", params[0])
	}

	if len(params) == 1 {
		return page, 0, nil
	}

	size, err := strconv.ParseUint(params[1], 10, 64)
	if err != nil {
		return 0, 0, errors.Wrapf(BadArguments, "Can't parse page size: %s", params[1])
	}

	return page, size, nil
}
