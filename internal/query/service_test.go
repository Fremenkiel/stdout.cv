package query

import (
	"context"
	"fmt"
	"testing"

	"github.com/fremenkiel/stdout.cv/internal/testutil"
)

func TestServiceRun(t *testing.T) {
	tests := []struct{
		name,
		query							string
		expectedResponse	*Query
		expectedError			error
	} {
		{
			name: "get_rows",
			query: "users",
			expectedResponse: &Query{
				Rows: []*Row{
					{
						stringPtr("1"),
						stringPtr("Kevin"),
					},
				},	
				Columns: []string{
					"id",
					"name",
					"age",
				},
			},
			expectedError: nil,
		},
	}

	for ti := range tests {
		test := tests[ti]

		t.Run(test.name, func(t *testing.T) {
			ctx := t.Context()

			service := NewService(
				&mockRepository{},
				NewMockSchemaService(),
				)

			response, err := service.Run(ctx, fmt.Sprintf("SELECT * FROM %s;", test.query))
			if test.expectedError != err {
				t.Fatalf("unexpected error, %v got %v", test.expectedError, err)
			}

			if test.expectedError != nil && err != nil {
				return
			}

			if len(test.expectedResponse.Rows) != len(response.Rows) {
				t.Fatalf("unexpected row len, %d got %d", len(test.expectedResponse.Rows), len(response.Rows))
			}

			if !testutil.SliceEql(test.expectedResponse.Columns, response.Columns) {
				t.Fatalf("unexpected column slice, %v got %v", test.expectedResponse.Columns, response.Columns)
			}

			for i := range response.Rows {
				row := response.Rows[i]
				expectedRow := test.expectedResponse.Rows[i]

				if len(*expectedRow) != len(*row) {
					t.Fatalf("unexpected row len, %d got %d", len(*expectedRow), len(*row))
				}

				for ri := range *row {
					field := (*row)[ri]
					expectedField := (*expectedRow)[ri]

					if *expectedField != *field {
						t.Fatalf("unexpected field, %s got %s", *expectedField, *field)
					}
				}
			}
		})
	}
}

type mockRepository struct {}

var _ repository = (*mockRepository)(nil)

func (r *mockRepository) GetRows(ctx context.Context, query string, columns []string) ([]*Row, error) {
	return []*Row{
		{
			stringPtr("1"),
			stringPtr("Kevin"),
		},
	}, nil
}

func (r *mockRepository) Execute(ctx context.Context, query string) error {
	return nil
}

type mockSchemaService struct {
	tableName		string
	tableSchema	[]string
}

var _ SchemaService = (*mockSchemaService)(nil)

func NewMockSchemaService() *mockSchemaService {
	return &mockSchemaService{
		tableName: "users",
		tableSchema: []string{
			"id",
			"name",
			"age",
		},
	}
}

func (s *mockSchemaService) AppendWildcardSchema(ctx context.Context, tableName string, fields []string) ([]string, error) {
	if tableName != s.tableName {
		return nil, fmt.Errorf("Unknown table, %s got %s", s.tableName, tableName)
	}

	return s.tableSchema, nil
}

func stringPtr(s string) *string {
	return &s
}

