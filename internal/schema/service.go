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

func (s *Service) AppendWildcardSchema(ctx context.Context, tableName string, columns []string) ([]string, error) {
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

	schema, err := s.GetSchema(ctx, tableName)
	if err != nil {
		return nil, err
	}

	for ti := range schema {
		columns = append(columns, schema[ti].Name)
	}

	columns = append(columns, endColumns...)

	return columns, nil
}

func (s *Service) GetSchema(ctx context.Context, name string) ([]*column.Column, error) {
	return s.repository.GetSchema(ctx, name)
}
