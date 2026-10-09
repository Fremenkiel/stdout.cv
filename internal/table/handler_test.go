package table

import (
	"net/http"
	"testing"

	"github.com/fremenkiel/stdout.cv/internal/platform/database"
	"github.com/fremenkiel/stdout.cv/internal/schema"
	"github.com/fremenkiel/stdout.cv/internal/session"
	"github.com/fremenkiel/stdout.cv/internal/testutil"
	"github.com/fremenkiel/stdout.cv/internal/ui"
)

func TestHandlerIndex(t *testing.T) {
	tests := []struct{
		name,
		method,
		url								string
		cookies						[]*http.Cookie
		hasSession				bool
		expectedResponse	int
	} {
		{
			name: "index",
			method: http.MethodGet,
			url: "/",
			cookies: []*http.Cookie{},
			hasSession: true,
			expectedResponse: http.StatusOK,
		},
		{
			name: "no_session",
			method: http.MethodGet,
			url: "/",
			cookies: []*http.Cookie{},
			hasSession: false,
			expectedResponse: http.StatusInternalServerError,
		},
	}

	for ti := range tests {
		test := tests[ti]

		t.Run(test.name, func(t *testing.T) {
			cache := session.NewCache(database.TestDatabaseFileNameTemplate)

			ctx := t.Context()

			if test.hasSession {
				sessionService := session.NewService(cache)

				sessionId, err := sessionService.CreateSession()
				if err != nil {
					t.Fatal(err)
				}

				ctx = session.NewContext(t.Context(), sessionId)
				defer testutil.CleanupSession(ctx)
			}

			handler := NewHandler(
				ui.NewRenderer(),
				NewService(
					NewRepository(cache),
					schema.NewService(schema.NewRepository(cache)),
					),
				)

			writer := testutil.NewResponseWriter()

			handler.Index(
				writer,
				testutil.NewRequest(ctx, test.method, test.url, test.cookies),
				)

			if test.expectedResponse != writer.Status() {
				t.Fatalf("unexpected status code, %d got %d", test.expectedResponse, writer.Status())
			}
		})
	}
}
