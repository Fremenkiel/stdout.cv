package table

import (
	"context"

	"github.com/fremenkiel/stdout.cv/internal/column"
	"github.com/fremenkiel/stdout.cv/internal/database"
)

type repository interface {
	GetTables(ctx context.Context) ([]*Table, error)
}

type databaseService interface {
	GetStats(ctx context.Context) (*database.Stats, error)
}

type schemaService interface {
	GetSchema(ctx context.Context, name string) ([]*column.Column, error)
}

type Service struct {
	repository			repository
	databaseService	databaseService
	schemaService		schemaService
}

var _ service = (*Service)(nil)

func NewService(r *Repository, ds databaseService, ss schemaService) *Service {
	return &Service{
		repository: r,
		databaseService: ds,
		schemaService: ss,
	}
}

func (s *Service) GetTables(ctx context.Context) ([]*Table, error) {
	tables, err := s.repository.GetTables(ctx)
	if err != nil {
		return nil, err
	}

	for i := range tables {
		tables[i].Columns, err = s.schemaService.GetSchema(ctx, tables[i].Name)
		if err != nil {
			return nil, err
		}
	}
	return tables, nil
}

func (s *Service) GetStats(ctx context.Context) (*database.Stats, error) {
	stats, err := s.databaseService.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
