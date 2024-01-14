package storage

import (
	"encoding/json"
	"go-web-bootcamp/internal"
	"io"
	"log"
	"os"
)

func ReadJson(db map[int]internal.Product, path string) (map[int]internal.Product, error) {
	jsonFile, err := os.Open(path)

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

	jsonFile.Close()

	return db, nil
}

func WriteJson(db map[int]internal.Product, path string) error {
	var products []internal.Product

	for _, value := range db {
		products = append(products, value)
	}

	bytes, err := json.Marshal(products)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
