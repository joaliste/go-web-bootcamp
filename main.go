package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

type Products struct {
	Products []Product `json:"products"`
}

func main() {
	jsonFile, err := os.Open("data/products.json")

	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)

	if err != nil {
		log.Fatal(err)
	}

	// var products Products
	var products []Product

	err = json.Unmarshal(byteValue, &products)

	if err != nil {
		log.Fatal(err)
	}

	jsonFile.Close()

}
