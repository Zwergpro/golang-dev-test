package storage

import (
	"github.com/pkg/errors"
	"log"
	"strconv"
)

var warehouse *Warehouse

var (
	ProductAlreadyExists = errors.New("product already exists")
	ProductNotExists     = errors.New("product does not exist")
)

func init() {
	warehouse = NewWarehouse()

	product, _ := NewProduct("pillow", 500, 10)
	if err := Add(product); err != nil {
		log.Fatal(err)
	}
}

func Get(id uint64) (*Product, error) {
	warehouse.RLock()
	defer warehouse.RUnlock()

	product, ok := warehouse.storage[id]
	if ok {
		return product, nil
	}
	return nil, errors.Wrap(ProductNotExists, strconv.FormatUint(id, 10))
}

func Add(p *Product) error {
	warehouse.Lock()
	defer warehouse.Unlock()

	if _, ok := warehouse.storage[p.GetId()]; ok {
		return errors.Wrap(ProductAlreadyExists, strconv.FormatUint(p.GetId(), 10))
	}
	warehouse.storage[p.GetId()] = p
	return nil
}

func Delete(id uint64) error {
	warehouse.Lock()
	defer warehouse.Unlock()

	if _, ok := warehouse.storage[id]; !ok {
		return errors.Wrap(ProductNotExists, strconv.FormatUint(id, 10))
	}
	delete(warehouse.storage, id)
	return nil
}

func Update(p *Product) error {
	warehouse.Lock()
	defer warehouse.mu.Unlock()

	if _, ok := warehouse.storage[p.GetId()]; !ok {
		return errors.Wrap(ProductNotExists, strconv.FormatUint(p.GetId(), 10))
	}
	warehouse.storage[p.GetId()] = p
	return nil
}

func List() []*Product {
	warehouse.RLock()
	defer warehouse.RUnlock()

	products := make([]*Product, 0, len(warehouse.storage))
	for _, v := range warehouse.storage {
		products = append(products, v)
	}
	return products
}
