package api

import (
	"net/http"

	"github.com/fareez-ahamed/go-ledger-rest/internal/api/handler"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handler.Health)
	return mux
}
