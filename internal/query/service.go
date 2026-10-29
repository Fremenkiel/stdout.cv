package query

import (
	"context"
	"log"

	"github.com/fremenkiel/stdout.cv/internal/database"
	"github.com/fremenkiel/stdout.cv/internal/interpreter"
	"github.com/fremenkiel/stdout.cv/internal/schema"
)

type repository interface {
	GetAll(ctx context.Context, query string, columns []string) ([]*Row, error)
	Execute(ctx context.Context, query string) error
}

type schemaService interface {
	GetSchema(ctx context.Context) (schema.Schema, error)
	AppendWildcardSchema(ctx context.Context, tableName string, fields []string) ([]string, error)
}

type databaseService interface {
	GetStats(ctx context.Context) (*database.Stats, error)
}

type Service struct {
	repo			repository
	schema		schemaService
	database	databaseService
}

var _ service = (*Service)(nil)

func NewService(r repository, ss schemaService, ds databaseService) *Service {
	return &Service{
		repo: r,
		schema: ss,
		database: ds,
	}
}

func (s *Service) Run(ctx context.Context, sql string) (*Query, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var response Query

	parser := interpreter.NewParser()
	if err := parser.ParseSql(sql); err != nil {
		return nil, err
	}

	executeQueries := parser.GetExecureQueries()
	for i := range executeQueries {
		if err := s.repo.Execute(ctx, executeQueries[i].String()); err != nil {
			return nil, err
		}
	}

	returnQuery := parser.GetReturnQuery()
	if returnQuery != nil {
		columns, err := s.schema.AppendWildcardSchema(
			ctx,
			returnQuery.Table().String(),
			returnQuery.Columns(),
			)
		if err != nil {
			return nil, err
		}

		var result Result
		result.Columns = columns

		rows, err := s.repo.GetAll(ctx, returnQuery.String(), result.Columns)
		if err != nil {
			return nil, err
		}

		result.Rows = rows

		response.Result = &result
	}

	log.Print(parser.AltersDatabase())
	if parser.AltersDatabase() {
		schema, err := s.GetSchema(ctx)
		if err != nil {
			return nil, err
		}

		response.Schema = schema

		stats, err := s.GetStats(ctx)
		if err != nil {
			return nil, err
		}

		response.Stats = stats
	}

	return &response, nil
}

func (s *Service) GetSchema(ctx context.Context) (schema.Schema, error) {
	schema, err := s.schema.GetSchema(ctx)
	if err != nil {
		return nil, err
	}

	return schema, nil
}

func (s *Service) GetStats(ctx context.Context) (*database.Stats, error) {
	stats, err := s.database.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	return stats, nil
}
