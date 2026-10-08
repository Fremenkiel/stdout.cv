package table

import (
	"context"

	"github.com/fremenkiel/stdout.cv/internal/column"
)

type SchemaService interface {
	GetSchema(ctx context.Context, name string) ([]*column.Column, error)
}

type Service struct {
	repository		*Repository
	schemaService	SchemaService
}

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
