package repository

import "github.com/pkg/errors"

var (
	ProductAlreadyExists = errors.New("product already exists")
	ProductNotExists     = errors.New("product does not exist")
)
