package row

import (
	"context"

	"github.com/fremenkiel/stdout.cv/internal/query"
	"github.com/fremenkiel/stdout.cv/internal/schema"
)


type SchemaService interface {
	GetSchema(ctx context.Context, name string) ([]*schema.Schema, error)
}

type Service struct {
	repository		*Repository
	schemaService	SchemaService
}

func NewService(r *Repository, ss SchemaService) *Service {
	return &Service{
		repository: r,
		schemaService: ss,
	}
}

func (s *Service) GetRows(ctx context.Context, queryString string) ([]*Row, error) {
	q, err := s.parseQuery(ctx, queryString)
	if err != nil {
		return nil, err
	}

	return s.repository.GetRows(ctx, queryString, q)
}

func (s *Service) parseQuery(ctx context.Context, queryString string) (*query.Query, error) {
	q, err := query.Parse(queryString)
	if err != nil {
		return nil, err
	}

	if q.WildcardIndex != -1 {
		tableSchema, err := s.schemaService.GetSchema(ctx, q.TableName)
		if err != nil {
			return nil, err
		}

		q.Fields, err = query.AppendWildcardSchema(q.WildcardIndex, q.Fields, tableSchema)
		if err != nil {
			return nil, err
		}
	}

	return q, nil
}
