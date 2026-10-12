package table

import (
	"context"

	"github.com/fremenkiel/stdout.cv/internal/column"
)

type TableRepository interface {
	GetTables(ctx context.Context) ([]*Table, error)
}

type SchemaService interface {
	GetSchema(ctx context.Context, name string) ([]*column.Column, error)
}

type Service struct {
	repository		TableRepository
	schemaService	SchemaService
}

var _ TableService = (*Service)(nil)

func NewService(r *Repository, ss SchemaService) *Service {
	return &Service{repository: r, schemaService: ss}
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
