package storage

import (
	"github.com/pkg/errors"
	"log"
	"strconv"
)

var warehouse map[uint64]*Product

var ProductAlreadyExists = errors.New("product already exists")
var ProductNotExists = errors.New("product does not exist")

type Interface interface {
	Add(p *Product) error
	Get(id uint64) (*Product, error)
	Update(p *Product) error
	Delete(id uint64) error
	List() []*Product
}

func init() {
	warehouse = make(map[uint64]*Product)
	product, _ := NewProduct("pillow", 500, 10)
	if err := Add(product); err != nil {
		log.Fatal(err)
	}
}

func Get(id uint64) (*Product, error) {
	product, ok := warehouse[id]
	if ok {
		return product, nil
	}
	return nil, errors.Wrap(ProductNotExists, strconv.FormatUint(uint64(id), 10))
}

func Add(p *Product) error {
	if _, ok := warehouse[p.GetId()]; ok {
		return errors.Wrap(ProductAlreadyExists, strconv.FormatUint(uint64(p.GetId()), 10))
	}
	warehouse[p.GetId()] = p
	return nil
}

func Delete(id uint64) error {
	if _, ok := warehouse[id]; !ok {
		return errors.Wrap(ProductNotExists, strconv.FormatUint(id, 10))
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
