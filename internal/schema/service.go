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

func AppendWildcardSchema(wildcardIndex int16, fields []string, schema []*column.Column) ([]string, error) {
	if fields[wildcardIndex] != "*" {
		return nil, ErrFieldNotWildcard
	}

	endFields := fields[wildcardIndex + 1:]
	fields = fields[:wildcardIndex]

	for ti := range schema {
		fields = append(fields, schema[ti].Name)
	}

	for ei := range endFields {
		fields = append(fields, endFields[ei])
	}

	return fields, nil
}

