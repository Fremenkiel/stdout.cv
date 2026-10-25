package query

import (
	"context"

	"github.com/fremenkiel/stdout.cv/internal/interpreter"
)

type repository interface {
	GetRows(ctx context.Context, query string, columns []string) ([]*Row, error)
	Execute(ctx context.Context, query string) error
}

type SchemaService interface {
	AppendWildcardSchema(ctx context.Context, tableName string, fields []string) ([]string, error)
}

type Service struct {
	repository		repository
	schemaService	SchemaService
}

var _ service = (*Service)(nil)

func NewService(r repository, ss SchemaService) *Service {
	return &Service{
		repository: r,
		schemaService: ss,
	}
}

func (s *Service) Run(ctx context.Context, sql string) (*Query, error) {
	var response Query

	parser := interpreter.NewParser()
	if err := parser.ParseSql(sql); err != nil {
		return nil, err
	}

	executeQueries := parser.GetExecureQueries()
	for i := range executeQueries {
		if err := s.repository.Execute(ctx, executeQueries[i].String()); err != nil {
			return nil, err
		}
	}

	returnQuery := parser.GetReturnQuery()
	if returnQuery == nil {
		return &response, nil
	}

	columns, err := s.schemaService.AppendWildcardSchema(
		ctx,
		returnQuery.Table().String(),
		returnQuery.Columns(),
		)
	if err != nil {
		return nil, err
	}

	response.Columns = columns

	rows, err := s.repository.GetRows(ctx, returnQuery.String(), response.Columns)
	if err != nil {
		return nil, err
	}

	response.Rows = rows

	return &response, nil
}
