package service

import (
	"errors"
	"fmt"
	"go-web-bootcamp/internal"
	"time"
)

func NewProductDefault(rp internal.ProductRepository) *ProductDefault {
	return &ProductDefault{
		rp: rp,
	}
}

type ProductDefault struct {
	rp internal.ProductRepository
}

func (p *ProductDefault) GetById(id int) (internal.Product, error) {
	product, err := p.rp.GetById(id)

	if err != nil {
		switch {
		case errors.Is(err, internal.ErrProductNotFound):
			err = fmt.Errorf("%w: id", internal.ErrProductNotFound)
		}
		return product, err
	}

	return product, nil
}

func (p *ProductDefault) GetProducts() []internal.ProductJSON {
	products := p.rp.GetProducts()

	return products
}

func (p *ProductDefault) GetProductsByPrice(price float64) []internal.Product {
	products := p.rp.GetProductsByPrice(price)

	return products
}

func (p *ProductDefault) Save(product *internal.Product) (err error) {
	// external services
	// ...

	// business logic
	// - validate required fields
	if err = ValidateBusinessRule(product); err != nil {
		return
	}

	// save product
	err = p.rp.Save(product)
	if err != nil {
		switch err {
		case internal.ErrProductCodeAlreadyExists:
			err = fmt.Errorf("%w: title", internal.ErrProductAlreadyExists)
		}
		return
	}

	return
}

func ValidateBusinessRule(p *internal.Product) error {
	// validate product (business rules)
	// - required
	if p.Name == "" {
		return fmt.Errorf("%w: name", internal.ErrFieldRequired)
	}
	if p.Quantity == 0 {
		return fmt.Errorf("%w: quantity", internal.ErrFieldRequired)
	}
	if p.CodeValue == "" {
		return fmt.Errorf("%w: code_value", internal.ErrFieldRequired)
	}
	if p.Expiration == "" {
		return fmt.Errorf("%w: expiration", internal.ErrFieldRequired)
	}
	if p.Price == 0 {
		return fmt.Errorf("%w: price", internal.ErrFieldRequired)
	}
	_, err := time.Parse("02/01/2006", p.Expiration)
	if err != nil {

		return fmt.Errorf("%w: must be in dd/mm/Y format", internal.ErrFieldFormat)
	}

	return nil
}

func (p *ProductDefault) Update(product *internal.Product) (err error) {
	// validate
	if err = ValidateBusinessRule(product); err != nil {
		return
	}

	// update product
	err = p.rp.Update(product)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrProductNotFound):
			err = fmt.Errorf("%w: id", internal.ErrProductNotFound)
		}
		return
	}
	return
}

func (p *ProductDefault) Delete(id int) (err error) {
	// delete product
	err = p.rp.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrProductNotFound):
			err = fmt.Errorf("%w: id", internal.ErrProductNotFound)
		}
		return
	}
	return
}
