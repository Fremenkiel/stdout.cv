package user

import "context"

type Service struct {
	repository	*Repository
}

func NewService(r *Repository) *Service {
	return &Service{repository: r}
}

func (s *Service) GetUser(ctx context.Context) (*User, error) {
	return s.repository.GetUser(ctx)
}
