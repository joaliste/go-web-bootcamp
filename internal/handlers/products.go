package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"go-web-bootcamp/internal"
	"net/http"
	"strconv"
)

var products []internal.Product

func GetProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusOK
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(code)
		err := json.NewEncoder(w).Encode(products)

		if err != nil {
			panic(err)
		}
	}
}

func LoadProducts(p []byte) error {
	err := json.Unmarshal(p, &products)

	if err != nil {
		return err
	}

	return nil
}

func Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("pong"))

		if err != nil {
			panic(err)
		}
	}
}

func GetProductById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		product := products[id-1]

		code := http.StatusOK
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(code)
		err = json.NewEncoder(w).Encode(product)

		if err != nil {
			panic(err)
		}
	}
}

func GetProductByPrice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		price := r.URL.Query().Get("priceGT")

		priceF, err := strconv.ParseFloat(price, 64)

		if err != nil {
			panic(err)
		}

		var queryProducts []internal.Product

		for _, p := range products {
			if p.Price > priceF {
				queryProducts = append(queryProducts, p)
			}
		}

		code := http.StatusOK
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(code)
		err = json.NewEncoder(w).Encode(queryProducts)

		if err != nil {
			panic(err)
		}
	}
}
