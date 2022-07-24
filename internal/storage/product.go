package storage

import "fmt"

var lastId = uint64(0)

type Product struct {
	id       uint64
	name     string
	price    uint64
	quantity uint64
}

func (p Product) GetId() uint64 {
	return p.id
}

func (p Product) GetName() string {
	return p.name
}

func (p *Product) SetName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("name length must be greater than 0")
	}
	p.name = name
	return nil
}

func (p Product) GetPrice() uint64 {
	return p.price
}

func (p *Product) SetPrice(price uint64) error {
	if price == 0 {
		return fmt.Errorf("price must be greater than 0")
	}
	p.price = price
	return nil
}

func (p Product) GetQuantity() uint64 {
	return p.quantity
}

func (p *Product) SetQuantity(quantity uint64) error {
	if quantity == 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}
	p.quantity = quantity
	return nil
}

func (p Product) String() string {
	return fmt.Sprintf("[%d] name:%s price:%d quantity:%d", p.id, p.name, p.price, p.quantity)
}

func NewProduct(name string, price uint64, quantity uint64) (*Product, error) {
	p := Product{}
	if err := p.SetName(name); err != nil {
		return nil, err
	}
	if err := p.SetPrice(price); err != nil {
		return nil, err
	}
	if err := p.SetQuantity(quantity); err != nil {
		return nil, err
	}

	lastId++
	p.id = lastId
	return &p, nil
}
