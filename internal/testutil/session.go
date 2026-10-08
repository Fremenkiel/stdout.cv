package testutil

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/fremenkiel/stdout.cv/internal/platform/database"
)

func CleanupSession(ctx context.Context) {
	sessionId := ctx.Value("session_key").(string)

	if err := os.Remove(fmt.Sprintf(database.TestDatabaseFileNameTemplate, sessionId)); err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}
}
