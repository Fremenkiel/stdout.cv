package table

import (
	"context"

)

type SchemaService interface {
	GetSchemas(ctx context.Context, tables []*Table) ([]*Table, error)
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
	return tables, nil

	// return s.schemaService.GetSchemas(ctx, tables)
}
