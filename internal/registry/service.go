package registry

import "context"

type repository interface {
	GetAll(ctx context.Context) (*Registry, error)
	Create(ctx context.Context, name string) (*int, error)
	Update(ctx context.Context, id int, name string) error
	Remove(ctx context.Context, id int) error
}

type Service struct {
	repo repository
}

func NewService(r repository) *Service {
	return &Service{
		repo: r,
	}
}

func (s *Service) GetAll(ctx context.Context) (*Registry, error) {
	registry, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return registry, nil
}

func (s *Service) Create(ctx context.Context, name string) (*int, error) {
	id, err := s.repo.Create(ctx, name)
	if err != nil {
		return nil, err
	}

	return id, nil
}

func (s *Service) Update(ctx context.Context, id int, name string) error {
	if err := s.repo.Update(ctx, id, name); err != nil {
		return err
	}

	return nil
}

func (s *Service) Remove(ctx context.Context, id int) error {
	if err := s.repo.Remove(ctx, id); err != nil {
		return err
	}

	return nil
}
