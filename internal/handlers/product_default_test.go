package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go-web-bootcamp/internal"
	"go-web-bootcamp/internal/repository"
	"go-web-bootcamp/internal/service"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func NewRequest(method, url string, body io.Reader, urlParams map[string]string, urlQuery map[string]string) *http.Request {
	// old request
	req := httptest.NewRequest(method, url, body)

	// url params
	// - new request with a new context with key chi.RouteCtxKey and value chiCtx -> "id":"1"
	if urlParams != nil {
		chiCtx := chi.NewRouteContext()
		for k, v := range urlParams {
			chiCtx.URLParams.Add(k, v)
		}
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	}

	// url query
	// r.URL.RawQuery = "name=task 1"
	if urlQuery != nil {
		query := req.URL.Query()
		for k, v := range urlQuery {
			query.Add(k, v)
		}
		req.URL.RawQuery = query.Encode() // raw string
	}

	return req
}

func TestTaskDefault_GetProducts(t *testing.T) {
	t.Run("success 01 - should return all the products", func(t *testing.T) {
		// arrange
		// - repository
		db := map[int]internal.Product{
			1: {
				Id:          1,
				Name:        "test_1",
				Quantity:    123,
				CodeValue:   "test_code_1",
				IsPublished: false,
				Expiration:  "02/02/2012",
				Price:       123,
			},
			2: {
				Id:          2,
				Name:        "test_2",
				Quantity:    123,
				CodeValue:   "test_code_2",
				IsPublished: false,
				Expiration:  "02/02/2012",
				Price:       123,
			},
			3: {
				Id:          3,
				Name:        "test_3",
				Quantity:    123,
				CodeValue:   "test_code_3",
				IsPublished: false,
				Expiration:  "02/02/2012",
				Price:       123,
			},
		}
		// repository
		rp, _ := repository.NewProductMap(db)
		// service
		sv := service.NewProductDefault(rp)
		// - handler
		hd := NewDefaultProducts(sv)
		hdFunc := hd.GetProducts()

		// act
		req := NewRequest("GET", "/products", nil, nil, nil)
		res := httptest.NewRecorder()
		hdFunc(res, req)

		// assert
		expectedCode := http.StatusOK
		expectedBody := `[{"id":1,"name":"test_1","quantity":123,"code_value":"test_code_1","is_published":false,"expiration":"02/02/2012","price":123},{"id":2,"name":"test_2","quantity":123,"code_value":"test_code_2","is_published":false,"expiration":"02/02/2012","price":123},{"id":3,"name":"test_3","quantity":123,"code_value":"test_code_3","is_published":false,"expiration":"02/02/2012","price":123}]`
		expectedHeader := http.Header{"Content-Type": []string{"application/json"}}
		require.Equal(t, expectedCode, res.Code)
		require.JSONEq(t, expectedBody, res.Body.String())
		require.Equal(t, expectedHeader, res.Header())
	})
}
