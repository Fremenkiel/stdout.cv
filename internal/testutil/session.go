package testutil

import (
	"os"

	"github.com/fremenkiel/stdout.cv/db"
)

type Cache struct {
	files			map[string]string
}

const TestSessionId = "0000-0000-0000-0001"

type Service struct {
	cache	*Cache
}

func NewService() *Service {
	return &Service{cache: &Cache{
		files: map[string]string{
			TestSessionId: "./db/test.db",
		},
	}}
}

// Creates a session id, a database file and saves it to the cache.
// Uses UpdateSession under the hood
func (s *Service) CreateSession() (string, error) {
	if err := s.UpdateSession(TestSessionId); err != nil {
		return "", err
	}

	return TestSessionId, nil
}

// Updates a session id check-in time.
// If no session exists then is't created and saved.
func (s *Service) UpdateSession(id string) error {
	if err := s.createSessionDatabase(id); err != nil {
		return err
	}

	return nil
}

func (s *Service) createSessionDatabase(id string) error {
	databaseName := s.cache.files[TestSessionId]
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
