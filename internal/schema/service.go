package schema

import (
	"context"

	"github.com/fremenkiel/stdout.cv/internal/row"
	"github.com/fremenkiel/stdout.cv/internal/table"
)

type Service struct {
	repository	*Repository
}

var _ table.SchemaService = (*Service)(nil)
var _ row.SchemaService = (*Service)(nil)

func NewService(r *Repository) *Service {
	return &Service{repository: r}
}

func (s *Service) GetSchema(ctx context.Context, name string) ([]*table.Column, error) {
	return s.repository.GetSchema(ctx, name)
}
