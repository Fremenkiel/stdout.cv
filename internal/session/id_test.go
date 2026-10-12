package session

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestNewContext(t *testing.T) {
	sessionUuid, err := uuid.NewV7()
	if err != nil {
		t.Fatal(err)
	}
	id := sessionUuid.String()

	ctx := NewContext(t.Context(), id)

	if ctxValue := ctx.Value(idKey).(string); id != ctxValue {
		t.Fatalf("unexpected context value, %s got %s", id, ctxValue)
	}
}

func TestFromContext(t *testing.T) {
	sessionUuid, err := uuid.NewV7()
	if err != nil {
		t.Fatal(err)
	}
	id := sessionUuid.String()

	ctx := context.WithValue(t.Context(), idKey, id)

	ctxValue, err := FromContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if id != ctxValue {
		t.Fatalf("unexpected context value, %s got %s", id, ctxValue)
	}
}
