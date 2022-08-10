package models

import "fmt"

type Product struct {
	Id       uint64 `db:"id"`
	Name     string `db:"name"`
	Price    uint64 `db:"price"`
	Quantity uint64 `db:"quantity"`
}

func (p Product) GetId() uint64 {
	return p.Id
}

func (p Product) GetName() string {
	return p.Name
}

func (p *Product) SetName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("name length must be greater than 0")
	}
	p.Name = name
	return nil
}

func (p Product) GetPrice() uint64 {
	return p.Price
}

func (p *Product) SetPrice(price uint64) error {
	if price == 0 {
		return fmt.Errorf("price must be greater than 0")
	}
	p.Price = price
	return nil
}

func (p Product) GetQuantity() uint64 {
	return p.Quantity
}

func (p *Product) SetQuantity(quantity uint64) error {
	if quantity == 0 {
		return fmt.Errorf("quantity must be greater than 0")
	}
	p.Quantity = quantity
	return nil
}

func (p Product) String() string {
	return fmt.Sprintf("[%d] name:%s price:%d quantity:%d", p.Id, p.Name, p.Price, p.Quantity)
}

func (p *Product) Copy() *Product {
	return &Product{
		Id:       p.Id,
		Name:     p.Name,
		Price:    p.Price,
		Quantity: p.Quantity,
	}
}

func BuildProduct(name string, price uint64, quantity uint64) (*Product, error) {
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

	return &p, nil
}
