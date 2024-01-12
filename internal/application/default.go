package application

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go-web-bootcamp/internal"
	"go-web-bootcamp/internal/handlers"
	"go-web-bootcamp/internal/repository"
	"go-web-bootcamp/internal/service"
	"log"
	"net/http"
)

func NewDefaultHTTP(address string) *DefaultHTTP {
	// default config / values
	defaultAddress := ":8080"
	if address != "" {
		defaultAddress = address
	}

	return &DefaultHTTP{
		address: defaultAddress,
	}
}

// DefaultHTTP is the server using chi
type DefaultHTTP struct {
	// address is the address to listen on
	address string
}

// Run runs the server
func (s *DefaultHTTP) Run() error {
	// - repository
	rp, err := repository.NewProductMap(make(map[int]internal.Product))
	if err != nil {
		log.Fatal(err)
	}

	// - service
	sv := service.NewProductDefault(rp)
	// Handler
	hd := handlers.NewDefaultProducts(sv)

	// - router
	rt := chi.NewRouter()

	// Endpoints
	rt.Get("/products/{id}", hd.GetByID())
	rt.Get("/products", hd.GetProducts())
	rt.Get("/ping", hd.Ping())
	rt.Post("/products", hd.Create())
	rt.Get("/products/search", hd.GetProductByPrice())

	fmt.Println("Server is running...")
	return http.ListenAndServe(s.address, rt)
}
