package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go-web-bootcamp/internal"
	"net/http"
	"strconv"
)

var products []internal.Product

type DefaultProducts struct {
	products []internal.Product
	lastID   int
}

func NewDefaultProducts(p []byte) (*DefaultProducts, error) {
	var products []internal.Product
	err := json.Unmarshal(p, &products)

	if err != nil {
		return nil, err
	}

	return &DefaultProducts{
		products: products,
		lastID:   len(products),
	}, nil
}

func (d *DefaultProducts) GetProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusOK
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		err := json.NewEncoder(w).Encode(d.products)

		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error getting the products"))
			return
		}
	}
}

func (d *DefaultProducts) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("pong"))

		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error writing pong response"))
			return
		}
	}
}

func (d *DefaultProducts) GetProductById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))

		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("id is not an integer"))
			return
		}

		if id <= 0 || id > d.lastID {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			message := fmt.Sprintf("id must be between 1 and %d", d.lastID)
			w.Write([]byte(message))
			return
		}

		product := d.products[id-1]

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(product)

		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error encoding the product"))
			return
		}
	}
}

func (d *DefaultProducts) GetProductByPrice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		price := r.URL.Query().Get("priceGT")

		priceF, err := strconv.ParseFloat(price, 64)

		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid price value"))
			return
		}

		var queryProducts []internal.Product

		for _, p := range d.products {
			if p.Price > priceF {
				queryProducts = append(queryProducts, p)
			}
		}

		code := http.StatusOK
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(code)
		err = json.NewEncoder(w).Encode(queryProducts)

		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error encoding the product"))
			return
		}
	}
}
