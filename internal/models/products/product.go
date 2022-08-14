package products

import "fmt"

type Product struct {
	Id       uint64 `db:"id"`
	Name     string `db:"name"`
	Price    uint64 `db:"price"`
	Quantity uint64 `db:"quantity"`
}

func (p *Product) GetId() uint64 {
	return p.Id
}

func (p *Product) GetName() string {
	return p.Name
}

func (p *Product) SetName(name string) error {
	if err := ValidateName(name); err != nil {
		return err
	}
	p.Name = name
	return nil
}

func (p *Product) GetPrice() uint64 {
	return p.Price
}

func (p *Product) SetPrice(price uint64) error {
	if err := ValidatePrice(price); err != nil {
		return err
	}
	p.Price = price
	return nil
}

func (p *Product) GetQuantity() uint64 {
	return p.Quantity
}

func (p *Product) SetQuantity(quantity uint64) error {
	if err := ValidateQuantity(quantity); err != nil {
		return err
	}
	p.Quantity = quantity
	return nil
}

func (p *Product) String() string {
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
