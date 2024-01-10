package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
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

	// Server
	rt := chi.NewRouter()

	// Endpoints
	rt.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("pong"))

		if err != nil {
			panic(err)
		}
	})

	fmt.Println("Server is starting...")
	if err := http.ListenAndServe(":8080", rt); err != nil {
		panic(err)
	}
}
