package testutil

import (
	"context"
	"errors"
	"fmt"
	"os"
)

func CleanupSession(ctx context.Context, filePath string) {
	sessionId := ctx.Value("session_key").(string)

	if err := os.Remove(fmt.Sprintf(filePath, sessionId)); err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}
}
