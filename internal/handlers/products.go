package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go-web-bootcamp/internal"
	"net/http"
	"strconv"
	"time"
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

func (d *DefaultProducts) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body internal.ProductRequestProductJSON

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid body"))
			return
		}

		product := internal.Product{
			Name:        body.Name,
			Quantity:    body.Quantity,
			CodeValue:   body.CodeValue,
			IsPublished: body.IsPublished,
			Expiration:  body.Expiration,
			Price:       body.Price,
		}

		if err := ValidateBusinessRule(&product); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprint(err)))
			return
		}

		for _, p := range d.products {
			if p.CodeValue == product.CodeValue {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("product already exists"))
				return
			}
		}

		d.lastID++
		product.Id = d.lastID

		// - store product
		d.products = append(d.products, product)

		// response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "product created",
			"data":    product,
		})
	}
}

func ValidateBusinessRule(p *internal.Product) error {
	// validate product (business rules)
	// - required
	if p.Name == "" {
		return errors.New("name is required")
	}
	if p.Quantity == 0 {
		return errors.New("quantity is required")
	}
	if p.CodeValue == "" {
		return errors.New("code_value is required")
	}
	if p.Expiration == "" {
		return errors.New("expiration is required")
	}
	if p.Price == 0 {
		return errors.New("price is required")
	}

	_, err := time.Parse("02/01/2006", p.Expiration)
	if err != nil {
		return errors.New("expiration date must be in dd/mm/Y format")
	}

	return nil
}
