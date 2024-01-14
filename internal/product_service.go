package internal

import "errors"

var (
	ErrFieldRequired        = errors.New("field required")
	ErrFieldFormat          = errors.New("bad field format")
	ErrProductAlreadyExists = errors.New("product already exists")
)

type ProductService interface {
	GetById(id int) (Product, error)
	GetProducts() []Product
	GetProductsByPrice(price float64) []Product
	Save(product *Product) (err error)
	Update(product *Product) (err error)
}
