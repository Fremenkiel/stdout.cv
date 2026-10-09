package session

import (
	"net/http"
	"testing"
	"testing/synctest"
	"time"

	"github.com/fremenkiel/stdout.cv/internal/platform/database"
	"github.com/fremenkiel/stdout.cv/internal/testutil"
	"github.com/google/uuid"
)

func TestHandle(t *testing.T) {
	tests := []struct{
		name,
		method,
		url,
		sessionId				string
		cookies					[]*http.Cookie
		existingSession,
		expiredSession,
		expectNext			bool
	} {
		{
			name: "initial_request",
			method: "GET",
			url: "/",
			sessionId: "1a100bf-c3be-7d9c-acb8-785e8722f892",
			cookies: []*http.Cookie{},
			existingSession: false,
			expiredSession: false,
			expectNext: true,
		},
		{
			name: "post_request",
			method: "POST",
			url: "/",
			sessionId: "1a100bf-c3be-7d9c-acb8-785e8722f892",
			cookies: []*http.Cookie{},
			existingSession: false,
			expiredSession: false,
			expectNext: true,
		},
		{
			name: "put_request",
			method: "PUT",
			url: "/",
			sessionId: "1a100bf-c3be-7d9c-acb8-785e8722f892",
			cookies: []*http.Cookie{},
			existingSession: false,
			expiredSession: false,
			expectNext: true,
		},
		{
			name: "delete_request",
			method: "DELETE",
			url: "/",
			sessionId: "1a100bf-c3be-7d9c-acb8-785e8722f892",
			cookies: []*http.Cookie{},
			existingSession: false,
			expiredSession: false,
			expectNext: true,
		},
		{
			name: "random_request",
			method: "GET",
			url: "/this-is-a-non-existing-endpoint",
			sessionId: "1a100bf-c3be-7d9c-acb8-785e8722f892",
			cookies: []*http.Cookie{},
			existingSession: false,
			expiredSession: false,
			expectNext: true,
		},
		{
			name: "existing_session",
			method: "GET",
			url: "/",
			sessionId: "1a100bf-c3be-7d9c-acb8-785e8722f892",
			cookies: []*http.Cookie{},
			existingSession: true,
			expiredSession: false,
			expectNext: true,
		},
		{
			name: "existing_session",
			method: "GET",
			url: "/",
			sessionId: "1a100bf-c3be-7d9c-acb8-785e8722f892",
			cookies: []*http.Cookie{},
			existingSession: true,
			expiredSession: true,
			expectNext: true,
		},
		{
			name: "empty_session_id",
			method: "GET",
			url: "/",
			sessionId: "",
			cookies: []*http.Cookie{},
			existingSession: true,
			expiredSession: false,
			expectNext: true,
		},
	}

	for ti := range tests {
		test := tests[ti]

		t.Run(test.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				cache := NewCache(database.TestDatabaseFileNameTemplate)
				service := NewService(cache)
				middleware := NewMiddleware(service)

				ctx := t.Context()

				if test.existingSession {
					if err := service.UpdateSession(test.sessionId); err != nil {
						t.Fatal(err)
					}

					test.cookies = append(test.cookies, &http.Cookie{
						Name:     cookieKey,
						Value:    test.sessionId,
						Path:     "/",
						MaxAge:   int(5*time.Minute),
						HttpOnly: true,
					})
					ctx := NewContext(t.Context(), test.sessionId)
					defer testutil.CleanupSession(ctx)
				}

				if test.expiredSession {
					if len(test.sessionId) == 0 {
						sessionUuid, err := uuid.NewV7()
						if err != nil {
							t.Fatal(err)
						}

						test.sessionId = sessionUuid.String()
					}

					cache.sessions.Swap(test.sessionId, time.Now().Add(-(10*time.Minute)))
				}

				writer := testutil.NewResponseWriter()
				request := testutil.NewRequest(
					ctx,
					test.method,
					test.url,
					test.cookies,
					)

				ranNext := false
				middleware.Handle(writer, request, func(w http.ResponseWriter, r *http.Request) {
					contextSessionId, err := FromContext(r.Context())
					if err != nil {
						t.Fatal(err)
					}

					if len(contextSessionId) == 0 {
						t.Fatal("no session id was set")
					}

					if test.existingSession && len(test.sessionId) > 0 && !test.existingSession {
						if test.sessionId != contextSessionId {
							t.Fatalf("unexpected context session id, %s got %s", test.sessionId, contextSessionId)
						}
					}

					requestTime, ok := cache.LoadSession(contextSessionId)

					if !ok {
						t.Fatal("request not present in cache")
					}

					if time.Now() != requestTime {
						t.Fatalf("unexpected request time, %s got %s", time.Now().String(), requestTime.String())
					}
					ranNext = true
				})

				if test.expectNext != ranNext {
					t.Fatalf("unexpected next run state, %v got %v", test.expectNext, ranNext)
				}
			})
		})
	}
}

