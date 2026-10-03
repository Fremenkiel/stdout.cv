package row

import (
	"context"
	"errors"
	"strings"

	"github.com/fremenkiel/stdout.cv/internal/schema"
)

var ErrQueryNotAllowed = errors.New("row: The received query is not allowed")

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

func (s *Service) GetRows(ctx context.Context, query string) ([]*Row, error) {
	queryParts := strings.Split(query, " ")

	if !strings.EqualFold(queryParts[0], "SELECT") {
		return nil, ErrQueryNotAllowed
	}

	var fields []string
	var table string
	var windcardIndex int16

	current := ""
	for i := range queryParts {
		if strings.EqualFold(queryParts[i], "SELECT") {
			current = "SELECT"
			continue
		}
		if strings.EqualFold(queryParts[i], "FROM") {
			current = "FROM"
			continue
		}
		if strings.EqualFold(queryParts[i], "WHERE") {
			current = "WHERE"
			continue
		}

		switch current {
		case "SELECT":
			fields, windcardIndex = s.parseSelectField(fields, queryParts[i])
		break
		case "FROM":
			table = queryParts[i]
		break
		case "WHERE":
		break
		}
	}

	if windcardIndex != -1 {
			tableSchema, err := s.schemaService.GetSchema(ctx, table)
			if err != nil {
				return nil, err
			}
			endFields := fields[windcardIndex + 1:]
			fields = fields[:windcardIndex]

			for ti := range tableSchema {
				fields = append(fields, tableSchema[ti].Name)
			}

			for ei := range endFields {
				fields = append(fields, endFields[ei])
			}
	}

	return s.repository.GetRows(ctx, query, fields, table)
}

func (s *Service) parseSelectField(fields []string, str string) ([]string, int16) {
	var windcardIndex int16 = -1

	selectFields := strings.Split(str, ",")
	for i := range selectFields {
		if selectFields[i] == "*" {
			windcardIndex = int16(len(fields) + i)
		}
		if len(selectFields[i]) > 0 {
			fields = append(fields, selectFields[i])
		}
	}

	return fields, windcardIndex
}
