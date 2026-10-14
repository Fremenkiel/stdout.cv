package table

import (
	"net/http"

	"github.com/fremenkiel/stdout.cv/internal/platform/middleware"
)

func NewRouter(mux *http.ServeMux, handler *Handler, middlewares ...func(http.Handler) http.Handler) {
	mux.Handle("GET /{$}", middleware.ChainMiddleware(middlewares...)(http.HandlerFunc(handler.Index)))
}

