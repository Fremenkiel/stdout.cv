package table

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fremenkiel/stdout.cv/internal/column"
	"github.com/fremenkiel/stdout.cv/internal/session"
	"github.com/fremenkiel/stdout.cv/internal/testutil"
)

func TestHandlerIndex(t *testing.T) {
	tests := []struct{
		name,
		method,
		url											string
		hasSession							bool
		expectedResponseContent	[]string
		expectedResponseCode		int
	} {
		{
			name: "index",
			method: http.MethodGet,
			url: "/",
			hasSession: true,
			expectedResponseContent: []string {
				"page:index",
				"test_table",
				"876",
			},
			expectedResponseCode: http.StatusOK,
		},
		{
			name: "no_session",
			method: http.MethodGet,
			url: "/",
			hasSession: false,
			expectedResponseContent: nil,
			expectedResponseCode: http.StatusInternalServerError,
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

			handler.Index(
				rr,
				httptest.NewRequestWithContext(ctx, test.method, test.url, nil),
				)

			if test.expectedResponseCode != rr.Code {
				t.Fatalf("unexpected status code, %d got %d", test.expectedResponseCode, rr.Code)
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

func (m *mockService) GetTables(ctx context.Context) ([]*Table, error) {
	_, ok := ctx.Value("session_key").(string)
	if !ok {
		return nil, session.ErrNoSessionInContext
	}

	return []*Table{
		{
			Name: "test_table",
			RowCount: 876,
			Columns: []*column.Column{
				{
					Name: "test_col",
					Type: "INTEGER",
				},
			},
		},
	}, nil
}
