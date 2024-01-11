package application

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"go-web-bootcamp/internal/handlers"
	"io"
	"log"
	"net/http"
	"os"
)

func NewServerChi(address string) *ServerChi {
	// default config / values
	defaultAddress := ":8080"
	if address != "" {
		defaultAddress = address
	}

	return &ServerChi{
		address: defaultAddress,
	}
}

// ServerChi is the server using chi
type ServerChi struct {
	// address is the address to listen on
	address string
}

// Run runs the server
func (s *ServerChi) Run() error {
	// - db
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

	if err != nil {
		log.Fatal(err)
	}

	jsonFile.Close()

	// - router
	rt := chi.NewRouter()

	// Endpoints
	rt.Get("/ping", hd.Ping())
	rt.Get("/products", hd.GetProducts())
	rt.Post("/products", hd.Create())
	rt.Get("/products/{id}", hd.GetProductById())
	rt.Get("/products/search", hd.GetProductByPrice())

	fmt.Println("Server is running...")
	return http.ListenAndServe(s.address, rt)
}
