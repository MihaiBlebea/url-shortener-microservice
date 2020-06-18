package api

import (
	"context"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// NewHTTPServer returns a new HTTP server
func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware)

	r.Methods("GET").Path("/r/{code}").Handler(httptransport.NewServer(
		endpoints.FindRedirect,
		decodeFindRedirectReq,
		redirectResponse,
	))

	r.Methods("GET").Path("/report").Handler(httptransport.NewServer(
		endpoints.Report,
		nil,
		encodeResponse,
	))

	r.Methods("POST").Path("/").Handler(httptransport.NewServer(
		endpoints.StoreRedirect,
		decodeStoreRedirectReq,
		encodeResponse,
	))

	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
