package session

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
)

func TestCreateSessionDatabase(t *testing.T) {
	service := NewService(nil)

	if err := ensureFolder(); err != nil {
		t.Fatal(err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		t.Fatal(err)
	}

	databaseName := fmt.Sprintf("/tmp/stdout_cv_sessions/session_%s.db", id)

	if err := service.createSessionDatabase(databaseName); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(databaseName); err != nil {
		t.Fatal(err)
	}
}

func TestRemoveSessionDatabase(t *testing.T) {
	service := NewService(nil)

	if err := ensureFolder(); err != nil {
		t.Fatal(err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		t.Fatal(err)
	}

	databaseName := fmt.Sprintf("/tmp/stdout_cv_sessions/session_%s.db", id)

	dstFileHandler, err := os.OpenFile(databaseName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		t.Fatal(err)
	}
	defer dstFileHandler.Close()

	if _, err = dstFileHandler.Write([]byte{1}); err != nil {
		t.Fatal(err)
	}

	if err := service.removeSessionDatabase(databaseName); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(databaseName); err == nil {
		t.Fatal("file was not removed")
	}
}

func ensureFolder() error {
	_, err := os.Stat("/tmp/stdout_cv_sessions")
	if err != nil {
		if err := os.MkdirAll("/tmp/stdout_cv_sessions", os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

