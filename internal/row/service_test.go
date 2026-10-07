package row

import (
	"context"
	"fmt"
	"testing"

	"github.com/fremenkiel/stdout.cv/internal/table"
	"github.com/fremenkiel/stdout.cv/internal/testutil"
)

func TestGetRow(t *testing.T) {
	tests := []struct{
		name,
		query,
		queryType,
		tableName			string
		tableSchema		[]string
		wildcardIndex	int16
	} {
		{
			name: "select_with_wildcard",
			query: "SELECT * FROM users;",
			queryType: "SELECT",
			tableName: "users",
			tableSchema: []string{
				"id",
				"name",
				"age",
			},
			wildcardIndex: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := NewService(nil, NewTestSchemaService(test.tableName, test.tableSchema))

			result, err := service.parseQuery(t.Context(), test.query)
			if err != nil {
				t.Fatal(err)
			}

			if !testutil.SliceEql(test.tableSchema, result.Fields) {
				t.Fatalf("unexpected response field slice, %v", result.Fields)
			}

			if test.tableName != result.TableName {
				t.Fatalf("unexpected table name, %s got %s", test.tableName, result.TableName)
			}

			if test.queryType != result.Type {
				t.Fatalf("unexpected query type, %s got %s", test.queryType, result.TableName)
			}

			if test.wildcardIndex != result.WildcardIndex {
				t.Fatalf("unexpected wildcard index, %d got %d", test.wildcardIndex, result.WildcardIndex)
			}
		})
	}
}

type testSchemaService struct {
	tableName		string
	tableSchema	[]*table.Column
}

func NewTestSchemaService(tableName string, tableSchemaSlice []string) *testSchemaService {
	tableSchema := make([]*table.Column, len(tableSchemaSlice))

	for i, column := range tableSchemaSlice {
		tableSchema[i] = &table.Column{
			Name: column,
		}
	}

	return &testSchemaService{
		tableName: tableName,
		tableSchema: tableSchema,
	}
}

func (s *testSchemaService) GetSchema(ctx context.Context, name string) ([]*table.Column, error) {
	if name != s.tableName {
		return nil, fmt.Errorf("Unknown table, %s got %s", s.tableName, name)
	}

	return s.tableSchema, nil
}

