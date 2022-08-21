package repository

type Repository struct {
	warehouse *Warehouse
}

func NewRepository(warehouse *Warehouse) *Repository {
	return &Repository{
		warehouse: warehouse,
	}
}
