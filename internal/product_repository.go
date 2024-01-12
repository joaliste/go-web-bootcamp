package internal

import "errors"

var (
	ErrProductCodeAlreadyExists = errors.New("product code already exists")
	ErrProductNotFound          = errors.New("product not found")
)

type ProductRepository interface {
	GetById(id int) (Product, error)
	GetProducts() []Product
	GetProductsByPrice(price float64) []Product
	Save(product *Product) (err error)
}
