package api

import (
	"context"

	"github.com/MihaiBlebea/url-shortener/shortener"
	"github.com/go-kit/kit/endpoint"
)

// Endpoints is a struct for all the endpoints exposed by the application
type Endpoints struct {
	FindRedirect  endpoint.Endpoint
	StoreRedirect endpoint.Endpoint
	Report        endpoint.Endpoint
}

// MakeEndpoints returns the Endpoints struct
func MakeEndpoints(service shortener.RedirectService) Endpoints {
	return Endpoints{
		FindRedirect:  makeFindRedirectEndpoint(service),
		StoreRedirect: makeStoreRedirectEndpoint(service),
	}
}

// makeFindRedirectEndpoint creates a FindRedirect endpoint
func makeFindRedirectEndpoint(service shortener.RedirectService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FindRedirectRequest)
		redirect, err := service.Find(req.Code)

		return FindRedirectResponse{URL: redirect.URL}, err
	}
}

// makeStoreRedirectEndpoint creates a StoreRedirect endpoint
func makeStoreRedirectEndpoint(service shortener.RedirectService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(StoreRedirectRequest)
		code, err := service.Store(shortener.Redirect{URL: req.URL})

		return StoreRedirectResponse{Code: code}, err
	}
}
