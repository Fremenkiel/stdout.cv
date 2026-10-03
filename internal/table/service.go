package table

import "context"

type Service struct {
	repository	*Repository
}

func NewService(r *Repository) *Service {
	return &Service{repository: r}
}

func (s *Service) GetTables(ctx context.Context) ([]*Table, error) {
	return s.repository.GetTables(ctx)
}
