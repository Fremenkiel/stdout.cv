package session

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/fremenkiel/stdout.cv/db"
	"github.com/fremenkiel/stdout.cv/internal/platform/logger"
	"github.com/google/uuid"
)

type Service struct {
	cache	*Cache
}

func NewService(c *Cache) *Service {
	return &Service{cache: c}
}

// Creates a session id, a database file and saves it to the cache.
// Uses UpdateSession under the hood
func (s *Service) CreateSession() (string, error) {
	sessionUuid, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	id := sessionUuid.String()

	if err := s.UpdateSession(id); err != nil {
		return "", err
	}

	return id, nil
}

// Updates a session id check-in time.
// If no session exists then is't created and saved.
func (s *Service) UpdateSession(id string) error {
	loaded := s.cache.UpdateSession(id)
	if !loaded {
		if err := s.createSessionDatabase(id); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) StartCleanupWorker(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		select {
		case <-ticker.C:
			if err := s.RemoveExpiredSessions(); err != nil {
				logger.LogError(
					"session/service",
					"error thrown while removing expired sessions",
					slog.Any("error", err),
					)
			}
		case <-ctx.Done():
			return
		}
	}()
}

// Removes all expired session, alongside any orphant cached file paths and database files not saved in the cache.
func (s *Service) RemoveExpiredSessions() error {
	expiredSessions := s.cache.LoadAndDeleteOldSessions()

	var returnErr error
	for _, id := range expiredSessions {
		if err := s.removeSessionDatabase(id); err != nil {
			returnErr = errors.Join(returnErr, err)
		}
	}

	orphantFiles := s.cache.LoadAndDeleteOrphantFiles()

	for _, id := range orphantFiles {
		if err := s.removeSessionDatabase(id); err != nil {
			returnErr = errors.Join(returnErr, err)
		}
	}

	dir, err := os.ReadDir(s.cache.databasePath)
	if err != nil {
		return err
	}

	for _, file := range dir {
		id := strings.Trim(
			strings.Trim(file.Name(), "session_"),
			".db",
			)

		if len(id) == 0 {
			continue
		}

		if _, ok := s.cache.LoadSession(id); !ok {
			if err := s.removeSessionDatabase(id); err != nil {
				returnErr = errors.Join(returnErr, err)
			}
		}
	}

	return returnErr
}

// Get the size of the session database.
func (s *Service) GetDatabaseSize(ctx context.Context, id string) (int64, error) {
	if ctx.Err() != nil {
		return -1, ctx.Err()
	}

	fileName, err := s.cache.GetFile(id)
	if err != nil {
		return -1, err
	}

	file, err := os.Stat(fileName)
	if err != nil {
		return -1, err
	}

	return file.Size(), nil
}

func (s *Service) createSessionDatabase(id string) error {
	databaseName := s.cache.AddFile(id)
	srcFile, err := db.Files.ReadFile("template.db")
	if err != nil {
		return err
	}

	dstFileHandler, err := os.OpenFile(databaseName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer dstFileHandler.Close()

	if _, err = dstFileHandler.Write(srcFile); err != nil {
		return err
	}

	return nil
}

func (s *Service) removeSessionDatabase(id string) error {
	if err := os.Remove(fmt.Sprintf(s.cache.databaseFileTemplate, id)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	s.cache.RemoveFile(id)

	return nil
}
