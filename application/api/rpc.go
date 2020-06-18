package api

import (
	"github.com/MihaiBlebea/url-shortener/shortener"
)

// RPC is an struct to hold the rpc methods
type RPC struct {
	redirectService shortener.RedirectService
}

// NewRPC is a cnstructor function for the RPC struct
func NewRPC(redirectService shortener.RedirectService) *RPC {
	return &RPC{redirectService}
}

// ShortenURL takes a url string and returns a url code that can be used to redirect
func (r *RPC) ShortenURL(url string, response *string) error {
	code, err := r.redirectService.Store(shortener.Redirect{URL: url})
	if err != nil {
		return err
	}

	*response = code

	return nil
}
