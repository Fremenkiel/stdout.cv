package row

import (
	"context"
	"time"

	"github.com/fremenkiel/stdout.cv/internal/column"
	"github.com/fremenkiel/stdout.cv/internal/interpreter"
)

type repository interface {
	GetRows(ctx context.Context, query string, columns []string) ([]*Row, error)
	Execute(ctx context.Context, query string) error
}

type SchemaService interface {
	AppendWildcardSchema(ctx context.Context, tableName string, fields []string) ([]string, error)
	GetSchema(ctx context.Context, name string) ([]*column.Column, error)
}

type Service struct {
	repository		repository
	schemaService	SchemaService
}

var _ RowService = (*Service)(nil)

func NewService(r repository, ss SchemaService) *Service {
	return &Service{
		repository: r,
		schemaService: ss,
	}
}

func (s *Service) GetRows(ctx context.Context, query string) (*Result, error) {
	var response Result

	parser := interpreter.NewParser()
	if err := parser.ParseQuery(query); err != nil {
		return nil, err
	}

	start := time.Now()

	executeQueries := parser.GetExecureQueries()
	for i := range executeQueries {
		if err := s.repository.Execute(ctx, executeQueries[i].String()); err != nil {
			return nil, err
		}
	}

	returnQuery := parser.GetReturnQuery()
	if returnQuery == nil {
		response.Duration = time.Now().UnixMicro() - start.UnixMicro()

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
	response.Duration = time.Now().UnixMicro() - start.UnixMicro()

	return &response, nil
}
