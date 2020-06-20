package api

import (
	"context"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// NewHTTPServer returns a new HTTP server
func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware)
	r.Use(logMiddleware)

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

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var logger log.Logger
		logger = log.NewLogfmtLogger(os.Stderr)
		logger.Log("path", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
