package schema

import "context"

type Service struct {
	repository	*Repository
}

func NewService(r *Repository) *Service {
	return &Service{repository: r}
}

func (s *Service) GetSchema(ctx context.Context, name string) ([]*Schema, error) {
	return s.repository.GetSchema(ctx, name)
}
