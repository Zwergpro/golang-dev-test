package handlers

import (
	"homework-1/internal/storage"
	"strings"
)

func listCmdHandler(_ string) string {
	products, err := storage.List()
	if err != nil {
		return err.Error()
	}

	if len(products) == 0 {
		return "Warehouse is empty"
	}

	res := make([]string, len(products))

	for _, p := range products {
		res = append(res, p.String())
	}

	return strings.Join(res, "\n")
}
