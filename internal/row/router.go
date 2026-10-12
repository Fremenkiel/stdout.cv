package row

import (
	"net/http"
)

type muxHandler interface {
	HandleFunc(string, func(w http.ResponseWriter, r *http.Request))
}

func NewRouter(mux muxHandler, handler *Handler) {
	mux.HandleFunc("GET /rows", handler.GetRows)
}

