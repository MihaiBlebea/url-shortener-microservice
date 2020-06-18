package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/MihaiBlebea/url-shortener/shortener"
	"github.com/gorilla/mux"
)

type (
	// FindRedirectRequest takes a code as param
	FindRedirectRequest struct {
		Code string `json:"code"`
	}

	// FindRedirectResponse returns ok bool as param
	FindRedirectResponse struct {
		URL      string `json:"url"`
		Response *http.Request
	}

	// StoreRedirectRequest takes an url as param
	StoreRedirectRequest struct {
		URL string `json:"url"`
	}

	// StoreRedirectResponse returns a code as param
	StoreRedirectResponse struct {
		Code string `json:"code"`
	}

	// ReportRequest takes no params
	ReportRequest struct {
	}

	// ReportResponse returns a list of redirects
	ReportResponse struct {
		Report []shortener.Redirect `json:"report"`
	}
)

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func redirectResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(FindRedirectResponse)
	http.Redirect(w, resp.Response, resp.URL, 303)

	return nil
}

func decodeStoreRedirectReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req StoreRedirectRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func decodeFindRedirectReq(ctx context.Context, r *http.Request) (interface{}, error) {
	var req FindRedirectRequest
	vars := mux.Vars(r)

	req = FindRedirectRequest{
		Code: vars["code"],
	}

	return req, nil
}
