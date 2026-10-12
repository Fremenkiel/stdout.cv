package table

import (
	"net/http"
)

func NewRouter(mux *http.ServeMux, middleware func(http.Handler) http.Handler, handler *Handler) {
	mux.Handle("GET /{$}", middleware(
		http.HandlerFunc(
			handler.Index,
			)))
}

