package table

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fremenkiel/stdout.cv/internal/testutil"
)

func TestNewRouter(t *testing.T) {
	mux := http.NewServeMux()

	handler := &Handler{
		testutil.NewRenderer(),
		&mockService{},
	}

	NewRouter(mux, handler)

	req := httptest.NewRequestWithContext(
		context.WithValue(t.Context(), "session_key", "test"),
		http.MethodGet, "/", nil)

	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if http.StatusOK != rr.Code {
		t.Fatalf("unexpected status code, %d got %d", http.StatusOK, rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "page:index") {
		t.Fatalf("unexpected response body, %s", rr.Body.String())
	}
}
