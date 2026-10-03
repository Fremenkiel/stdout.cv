package middleware

import (
	"container/list"
	"net/http"
)

type MiddlewareMux struct {
	http.ServeMux
	middlewares list.List
}

func (mux *MiddlewareMux) AppendMiddleware(middleware Middleware) {
	mux.middlewares.PushBack(middleware)
}

func (mux *MiddlewareMux) PrependMiddleware(middleware Middleware) {
	mux.middlewares.PushFront(middleware)
}

func (mux *MiddlewareMux) nextMiddleware(el *list.Element) func(w http.ResponseWriter, req *http.Request) {
	if el != nil {
		return func(w http.ResponseWriter, req *http.Request) {
			el.Value.(Middleware)(w, req, mux.nextMiddleware(el.Next()))
		}
	}
	return mux.ServeMux.ServeHTTP
}

func (mux *MiddlewareMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mux.nextMiddleware(mux.middlewares.Front())(w, req)
}

