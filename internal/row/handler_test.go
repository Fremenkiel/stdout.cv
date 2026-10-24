package row

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fremenkiel/stdout.cv/internal/session"
	"github.com/fremenkiel/stdout.cv/internal/testutil"
)

func TestHandlerGetRows(t *testing.T) {
	tests := []struct{
		name,
		method,
		url								string
		hasSession				bool
		expectedResponseContent	[]string
		expectedResponse	int
	} {
		{
			name: "get_row",
			method: http.MethodGet,
			url: "/rows?query=select%20name%20from%20users%3B%0A",
			hasSession: true,
			expectedResponseContent: []string {
				"fragment:result",
				"test_row",
				"test_col",
			},
			expectedResponse: http.StatusOK,
		},
		{
			name: "get_row_no_query",
			method: http.MethodGet,
			url: "/rows",
			hasSession: true,
			expectedResponseContent: nil,
			expectedResponse: http.StatusNoContent,
		},
		{
			name: "no_session",
			method: http.MethodGet,
			url: "/rows?query=select%20name%20from%20users%3B%0A",
			hasSession: false,
			expectedResponseContent: nil,
			expectedResponse: http.StatusInternalServerError,
		},
	}

	for ti := range tests {
		test := tests[ti]

		t.Run(test.name, func(t *testing.T) {
			ctx := t.Context()
			if test.hasSession {
				ctx = context.WithValue(ctx, "session_key", "test")
			}

			handler := NewHandler(
				testutil.NewRenderer(),
				&mockService{},
				)

			rr := httptest.NewRecorder()

			handler.GetRows(
				rr,
				httptest.NewRequestWithContext(ctx, test.method, test.url, nil),
				)

			if test.expectedResponse != rr.Code {
				t.Fatalf("unexpected status code, %d got %d", test.expectedResponse, rr.Code)
			}

			for _, content := range test.expectedResponseContent {
				if !strings.Contains(rr.Body.String(), content) {
					t.Fatalf("unexpected response body, %s got %s", content, rr.Body.String())
				}
			}
		})
	}
}

type mockService struct {}

func (m *mockService) GetRows(ctx context.Context, query string) (*Result, error) {
	_, ok := ctx.Value("session_key").(string)
	if !ok {
		return nil, session.ErrNoSessionInContext
	}

	return &Result{
		Rows: []*Row{
			{
				stringPtr("test_row"),
			},
		},
		Columns: []string{
			"test_col",
		},
		Duration: 765,
	}, nil
}
