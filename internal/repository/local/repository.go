package repository

type Repository struct {
	warehouse *Warehouse
}

func NewRepository() *Repository {
	return &Repository{
		warehouse: NewWarehouse(),
	}
}
