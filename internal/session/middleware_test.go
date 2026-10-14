package session

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/synctest"
	"time"

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
		expiredSession	bool
	} {
		{
			name: "initial_request",
			method: "GET",
			url: "/",
			sessionId: "1a100bf-c3be-7d9c-acb8-785e8722f892",
			cookies: []*http.Cookie{},
			existingSession: false,
			expiredSession: false,
		},
		{
			name: "post_request",
			method: "POST",
			url: "/",
			sessionId: "1a100bf-c3be-7d9c-acb8-785e8722f892",
			cookies: []*http.Cookie{},
			existingSession: false,
			expiredSession: false,
		},
		{
			name: "put_request",
			method: "PUT",
			url: "/",
			sessionId: "1a100bf-c3be-7d9c-acb8-785e8722f892",
			cookies: []*http.Cookie{},
			existingSession: false,
			expiredSession: false,
		},
		{
			name: "delete_request",
			method: "DELETE",
			url: "/",
			sessionId: "1a100bf-c3be-7d9c-acb8-785e8722f892",
			cookies: []*http.Cookie{},
			existingSession: false,
			expiredSession: false,
		},
		{
			name: "random_request",
			method: "GET",
			url: "/this-is-a-non-existing-endpoint",
			sessionId: "1a100bf-c3be-7d9c-acb8-785e8722f892",
			cookies: []*http.Cookie{},
			existingSession: false,
			expiredSession: false,
		},
		{
			name: "existing_session",
			method: "GET",
			url: "/",
			sessionId: "1a100bf-c3be-7d9c-acb8-785e8722f892",
			cookies: []*http.Cookie{},
			existingSession: true,
			expiredSession: false,
		},
		{
			name: "existing_session",
			method: "GET",
			url: "/",
			sessionId: "1a100bf-c3be-7d9c-acb8-785e8722f892",
			cookies: []*http.Cookie{},
			existingSession: true,
			expiredSession: true,
		},
		{
			name: "empty_session_id",
			method: "GET",
			url: "/",
			sessionId: "",
			cookies: []*http.Cookie{},
			existingSession: true,
			expiredSession: false,
		},
	}

	for ti := range tests {
		test := tests[ti]

		t.Run(test.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				env := testutil.LoadTestEnv(t)
				cache := NewCache(env.DatabasePath, env.DatabaseFileTemplate)
				service := NewService(cache)

				handler := func(w http.ResponseWriter, r *http.Request) {
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

					w.WriteHeader(http.StatusOK)
				}

				middleware := Middleware(service)(http.HandlerFunc(handler))

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
					defer testutil.CleanupSession(ctx, env.DatabasePath)
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

				rr := httptest.NewRecorder()
				request := testutil.NewRequest(
					ctx,
					test.method,
					test.url,
					test.cookies,
					)


				f, ok := middleware.(http.HandlerFunc)
				if !ok || f == nil {
					t.Fatal("middleware not of type handler func")
				}

				f(rr, request)
				if rr.Code != http.StatusOK {
					t.Fatalf("unexpected recorder status code, got %d", rr.Code)
				}
			})
		})
	}
}

