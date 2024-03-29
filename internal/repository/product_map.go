package repository

import (
	"go-web-bootcamp/internal"
	"go-web-bootcamp/internal/storage"
	"log"
)

const FilePath = "/Users/joaliste/Desktop/Meli/Bootcamp GO/Repositories/go-web-bootcamp/data/products_updated.json"

func NewProductMap(db map[int]internal.Product) (*ProductMap, error) {

	if len(db) == 0 {
		_, err := storage.ReadJson(db, FilePath)
		if err != nil {
			log.Fatal("could not read the data")
		}
	}

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

func (pm *ProductMap) GetProducts() []internal.ProductJSON {
	var products []internal.ProductJSON

	for _, value := range pm.db {
		products = append(products, internal.ProductJSON(value))
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

	err = storage.WriteJson((*pm).db, FilePath)

	if err != nil {
		return
	}

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

	err = storage.WriteJson((*pm).db, FilePath)

	if err != nil {
		return
	}

	return
}

func (pm *ProductMap) Delete(id int) (err error) {
	_, ok := pm.db[id]
	if !ok {
		err = internal.ErrProductNotFound
		return
	}
	delete(pm.db, id)

	err = storage.WriteJson(pm.db, FilePath)

	if err != nil {
		return
	}

	return
}
