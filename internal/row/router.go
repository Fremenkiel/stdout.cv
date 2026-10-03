package row

import (
	"github.com/fremenkiel/stdout.cv/internal/platform/middleware"
)

func NewRouter(mux *middleware.MiddlewareMux, handler *Handler) {
	mux.HandleFunc("GET /rows", handler.GetRows)
}

