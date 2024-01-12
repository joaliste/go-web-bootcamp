package handlers

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"go-web-bootcamp/internal"
	"go-web-bootcamp/platform/web/request"
	"go-web-bootcamp/platform/web/response"
	"net/http"
	"strconv"
)

type DefaultProducts struct {
	sv internal.ProductService
}

type ProductJSON struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

func NewDefaultProducts(sv internal.ProductService) *DefaultProducts {
	return &DefaultProducts{
		sv: sv,
	}
}

func (d *DefaultProducts) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// - get id from path
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Text(w, http.StatusBadRequest, "invalid id")
			return
		}

		// process
		product, err := d.sv.GetById(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrProductNotFound):
				response.Text(w, http.StatusNotFound, "product not found")
			default:
				response.Text(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// response
		data := ProductJSON{
			Id:          product.Id,
			Name:        product.Name,
			Quantity:    product.Quantity,
			CodeValue:   product.CodeValue,
			IsPublished: product.IsPublished,
			Expiration:  product.Expiration,
			Price:       product.Price,
		}
		// - response
		response.JSON(w, http.StatusOK, map[string]any{
			"message": "product found",
			"data":    data,
		})
	}
}

func (d *DefaultProducts) GetProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request

		// process
		products := d.sv.GetProducts()

		// - response
		response.JSON(w, http.StatusOK, products)
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

func (d *DefaultProducts) GetProductByPrice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		price := r.URL.Query().Get("priceGT")

		priceF, err := strconv.ParseFloat(price, 64)

		if err != nil {
			response.Text(w, http.StatusBadRequest, "invalid price value")
			return
		}

		products := d.sv.GetProductsByPrice(priceF)

		response.JSON(w, http.StatusOK, products)
	}
}

func (d *DefaultProducts) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body internal.ProductRequestProductJSON

		if err := request.JSON(r, &body); err != nil {
			response.Text(w, http.StatusBadRequest, "invalid body")
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

		if err := d.sv.Save(&product); err != nil {
			switch {
			case errors.Is(err, internal.ErrFieldRequired), errors.Is(err, internal.ErrFieldFormat):
				response.Text(w, http.StatusBadRequest, "invalid or missing field in body")
			case errors.Is(err, internal.ErrProductAlreadyExists):
				response.Text(w, http.StatusConflict, "product already exists")
			default:
				response.Text(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		// response
		// - deserialize ProductJSON
		data := ProductJSON{
			Id:          product.Id,
			Name:        product.Name,
			Quantity:    product.Quantity,
			CodeValue:   product.CodeValue,
			IsPublished: product.IsPublished,
			Expiration:  product.Expiration,
			Price:       product.Price,
		}
		// w.Header().Set("Content-Type", "application/json")
		// w.WriteHeader(http.StatusCreated)
		// json.NewEncoder(w).Encode(map[string]any{
		// 	"message": "product created",
		// 	"data": data,
		// })
		response.JSON(w, http.StatusCreated, map[string]any{
			"message": "product created",
			"data":    data,
		})
	}
}