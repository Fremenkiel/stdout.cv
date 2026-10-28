package session

import (
	"context"
	"errors"
)

var ErrNoSessionInContext = errors.New("No session id found in context")

type Id = string

type key = string

const idKey key = "session_key"

func NewContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, idKey, id)
}

func FromContext(ctx context.Context) (Id, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	id, ok := ctx.Value(idKey).(string)
	if !ok {
		return "", ErrNoSessionInContext
	}

	return id, nil
}
