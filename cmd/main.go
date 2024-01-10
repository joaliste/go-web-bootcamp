package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go-web-bootcamp/internal/handlers"
	"io"
	"log"
	"net/http"
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

	err = handlers.LoadProducts(byteValue)

	jsonFile.Close()

	// Server
	rt := chi.NewRouter()

	// Endpoints
	rt.Get("/ping", handlers.Ping())

	rt.Get("/products", handlers.GetProducts())
	rt.Get("/products/{id}", handlers.GetProductById())
	rt.Get("/products/search", handlers.GetProductByPrice())

	fmt.Println("Server is running...")
	if err := http.ListenAndServe(":8080", rt); err != nil {
		panic(err)
	}
}
