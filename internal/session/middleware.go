package session

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/fremenkiel/stdout.cv/internal/platform"
)

type Middleware struct {
	service *Service
}

func NewMiddleware(s *Service) *Middleware {
	return &Middleware{service: s}
}

func (m *Middleware) Handle(w http.ResponseWriter, r *http.Request, next func(http.ResponseWriter, *http.Request)) {
	var sessionId string
	sessionCookie, err := r.Cookie(platform.SessionKey)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			sessionId, err = m.service.CreateSession()
			if err != nil {
				log.Printf("middlware: error thrown while creating session, %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			log.Printf("middlware: error thrown while getting cookie, %v", err)
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
	} else {
		if err := m.service.UpdateSession(sessionCookie.Value); err != nil {
			log.Printf("middlware: error thrown while updating session, %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	setCookie(w, sessionId)

	ctx := context.WithValue(r.Context(), platform.SessionKey, sessionId)
	next(w, r.WithContext(ctx))
}

func setCookie(w http.ResponseWriter, id string) {
	http.SetCookie(w, &http.Cookie{
		Name:     platform.SessionKey,
		Value:    id,
		Path:     "/",
		MaxAge:   int(5*time.Minute),
		HttpOnly: true,
	})
}

