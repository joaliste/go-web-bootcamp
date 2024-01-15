package internal

import "errors"

var (
	ErrProductCodeAlreadyExists = errors.New("product code already exists")
	ErrProductNotFound          = errors.New("product not found")
)

type ProductRepository interface {
	GetById(id int) (Product, error)
	GetProducts() []ProductJSON
	GetProductsByPrice(price float64) []Product
	Save(product *Product) (err error)
	Update(product *Product) (err error)
	Delete(id int) (err error)
}
