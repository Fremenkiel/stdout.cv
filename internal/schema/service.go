package schema

import (
	"context"
)

type repository interface {
	GetSchema(ctx context.Context) (Schema, error)
	GetColumns(ctx context.Context, name string) ([]*Column, error)
}

type Service struct {
	repo	repository
}

func NewService(r repository) *Service {
	return &Service{repo: r}
}

func (s *Service) GetSchema(ctx context.Context) (Schema, error) {
	schema, err := s.repo.GetSchema(ctx)
	if err != nil {
		return nil, err
	}

	for i := range schema {
		schema[i].Columns, err = s.GetColumns(ctx, schema[i].Name)
		if err != nil {
			return nil, err
		}
	}

	return schema, nil
}

func (s *Service) AppendWildcardSchema(ctx context.Context, tableName string, columns []string) ([]string, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var wildcardIndex int = -1
	for i, column := range columns {
		if column == "*" {
			wildcardIndex = i
		}
	}

	if wildcardIndex == -1 {
		return columns, nil
	}

	endColumns := columns[wildcardIndex + 1:]
	columns = columns[:wildcardIndex]

	schema, err := s.GetColumns(ctx, tableName)
	if err != nil {
		return nil, err
	}

	for ti := range schema {
		columns = append(columns, schema[ti].Name)
	}

	columns = append(columns, endColumns...)

	return columns, nil
}

func (s *Service) GetColumns(ctx context.Context, name string) ([]*Column, error) {
	columns, err := s.repo.GetColumns(ctx, name)
	if err != nil {
		return nil, err
	}

	return columns, nil
}
