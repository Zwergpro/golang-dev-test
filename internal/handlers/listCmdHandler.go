package handlers

import (
	"homework-1/internal/storage"
	"strings"
)

func listCmdHandler(_ string) string {
	products := storage.List()
	res := make([]string, len(products))

	for _, p := range products {
		res = append(res, p.String())
	}

	return strings.Join(res, "\n")
}
