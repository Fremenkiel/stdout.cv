package user

import "net/http"

func NewRouter(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("GET /user", handler.Index)
}

