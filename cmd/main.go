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

	// Handler
	hd, err := handlers.NewDefaultProducts(byteValue)

	jsonFile.Close()

	// Server
	rt := chi.NewRouter()

	// Endpoints
	if err != nil {
		log.Fatal(err)
	}

	rt.Get("/ping", hd.Ping())
	rt.Get("/products", hd.GetProducts())
	rt.Get("/products/{id}", hd.GetProductById())
	rt.Get("/products/search", hd.GetProductByPrice())

	fmt.Println("Server is running...")
	if err := http.ListenAndServe(":8080", rt); err != nil {
		log.Fatal(err)
	}
}
