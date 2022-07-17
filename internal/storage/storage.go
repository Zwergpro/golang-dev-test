package storage

import (
	"github.com/pkg/errors"
	"log"
	"strconv"
)

var warehouse map[uint]*Product

var ProductAlreadyExists = errors.New("product already exists")
var ProductNotExists = errors.New("product does not exist")

func init() {
	warehouse = make(map[uint]*Product)
	product, _ := NewProduct("pillow", 500, 10)
	if err := Add(product); err != nil {
		log.Panic(err)
	}
}

func Add(p *Product) error {
	if _, ok := warehouse[p.GetId()]; ok {
		return errors.Wrap(ProductAlreadyExists, strconv.FormatUint(uint64(p.GetId()), 10))
	}
	warehouse[p.GetId()] = p
	return nil
}

func Delete(id uint) error {
	if _, ok := warehouse[id]; !ok {
		return errors.Wrap(ProductNotExists, strconv.FormatUint(uint64(id), 10))
	}
	delete(warehouse, id)
	return nil
}

func Update(p *Product) error {
	if _, ok := warehouse[p.GetId()]; !ok {
		return errors.Wrap(ProductNotExists, strconv.FormatUint(uint64(p.GetId()), 10))
	}
	warehouse[p.GetId()] = p
	return nil
}

func List() []*Product {
	products := make([]*Product, 0, len(warehouse))
	for _, v := range warehouse {
		products = append(products, v)
	}
	return products
}
