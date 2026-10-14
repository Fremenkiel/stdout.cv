package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type middlewareIndex struct{
	id		string
	index	uint8
}

func TestChainMiddleware(t *testing.T) {
	middlewareIndexes := []middlewareIndex{
		{ id: "middleware1", index: 0 },
		{ id: "middleware2", index: 1 },
		{ id: "middleware3", index: 2 },
		{ id: "middleware4", index: 3 },
	}

	middlewares := make([]func(http.Handler) http.Handler, len(middlewareIndexes))
	for i, mi:= range middlewareIndexes {
		middlewares[i] = mockMiddleware(mi.id)
	}

	var index uint8 = 0
	ctx := context.WithValue(t.Context(), "index", index)

	rr := httptest.NewRecorder()
	req := httptest.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"/",
		nil,
		)

	handler := ChainMiddleware(middlewares...)(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			for _, mi := range middlewareIndexes {
				ctxValue := r.Context().Value(mi.id).(uint8)
				if mi.index != ctxValue {
					t.Fatalf("unexpected context index for %s, %d got %d", mi.id, mi.index, ctxValue)
				}
			}
			w.WriteHeader(http.StatusOK)
		}))

	handler.(http.HandlerFunc)(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("unexpected status code, got %v", rr.Code)
	}
}

func mockMiddleware(id string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			index := ctx.Value("index").(uint8)

			r = r.WithContext(
				context.WithValue(
					context.WithValue(ctx, "index", index + 1),
					id, index),
				)

			next.ServeHTTP(w, r)
		})
	}
}
