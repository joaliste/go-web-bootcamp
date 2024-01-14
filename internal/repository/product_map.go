package repository

import (
	"encoding/json"
	"go-web-bootcamp/internal"
	"io"
	"log"
	"os"
)

func NewProductMap(db map[int]internal.Product) (*ProductMap, error) {
	// default config / values
	jsonFile, err := os.Open("data/products.json")

	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)

	if err != nil {
		log.Fatal(err)
	}

	var products []internal.Product
	err = json.Unmarshal(byteValue, &products)

	for _, value := range products {
		db[value.Id] = value
	}

	if err != nil {
		return nil, err
	}

	if err != nil {
		log.Fatal(err)
	}

	jsonFile.Close()

	return &ProductMap{
		db:     db,
		lastId: len(db),
	}, nil
}

type ProductMap struct {
	db map[int]internal.Product
	// lastId is the last id of the product
	lastId int
}

func (pm *ProductMap) GetById(id int) (internal.Product, error) {
	product, ok := pm.db[id]
	if !ok {
		err := internal.ErrProductNotFound
		return product, err
	}

	return product, nil
}

func (pm *ProductMap) GetProducts() []internal.Product {
	var products []internal.Product

	for _, value := range pm.db {
		products = append(products, value)
	}

	return products
}

func (pm *ProductMap) GetProductsByPrice(price float64) []internal.Product {
	var products []internal.Product

	for _, value := range pm.db {
		if value.Price >= price {
			products = append(products, value)
		}
	}

	return products
}

func (pm *ProductMap) Save(product *internal.Product) (err error) {
	// validate product (consistency)
	// - code_value: unique
	if err = pm.ValidateProductCode((*product).CodeValue); err != nil {
		return
	}

	// autoincrement
	// - increment id
	(*pm).lastId++
	// - set id
	(*product).Id = (*pm).lastId

	// store product
	(*pm).db[(*product).Id] = *product

	return
}

func (pm *ProductMap) ValidateProductCode(code string) (err error) {
	// validate product (consistency)
	// - code_value: unique
	for _, v := range (*pm).db {
		if v.CodeValue == code {
			return internal.ErrProductCodeAlreadyExists
		}
	}

	return
}

func (pm *ProductMap) Update(product *internal.Product) (err error) {

	_, ok := pm.db[(*product).Id]
	if !ok {
		err = internal.ErrProductNotFound
		return
	}

	pm.db[(*product).Id] = *product
	return
}
