package row

import (
	"net/http"

	"github.com/fremenkiel/stdout.cv/internal/platform/middleware"
)

func NewRouter(mux *http.ServeMux, handler *Handler, middlewares ...func(http.Handler) http.Handler) {
	mux.Handle("GET /rows", middleware.ChainMiddleware(middlewares...)(http.HandlerFunc(handler.GetRows)))
}

