package session

import (
	"os"

	"github.com/fremenkiel/stdout.cv/db"
	"github.com/google/uuid"
)

type Service struct {
	cache	*Cahce
}

func NewService(c *Cahce) *Service {
	return &Service{cache: c}
}

func (s *Service) CreateSession() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	sessionId := id.String()

	s.cache.AddSession(sessionId)

	databaseName, err := s.cache.GetFile(sessionId)
	if err != nil {
		return "", err
	}

	if err := s.createSessionDatabase(databaseName); err != nil {
		return "", err
	}

	return sessionId, nil
}

func (s *Service) UpdateSession(id string) error {
	loaded := s.cache.UpdateSession(id)
	if !loaded {
		databaseName, err := s.cache.GetFile(id)
		if err != nil {
			return err
		}

		if err := s.createSessionDatabase(databaseName); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) createSessionDatabase(databaseName string) error {
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

func (s *Service) removeSessionDatabase(databaseName string) error {
	return os.Remove(databaseName)
}
