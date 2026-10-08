package schema

import (
	"context"

	"github.com/fremenkiel/stdout.cv/internal/column"
)

type Service struct {
	repository	*Repository
}

func NewService(r *Repository) *Service {
	return &Service{repository: r}
}

func (s *Service) GetSchema(ctx context.Context, name string) ([]*column.Column, error) {
	return s.repository.GetSchema(ctx, name)
}
