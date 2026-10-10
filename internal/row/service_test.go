package row

import (
	"context"
	"fmt"
	"testing"
	"testing/synctest"
	"time"

	"github.com/fremenkiel/stdout.cv/internal/column"
	"github.com/fremenkiel/stdout.cv/internal/platform/database"
	"github.com/fremenkiel/stdout.cv/internal/query"
	"github.com/fremenkiel/stdout.cv/internal/schema"
	"github.com/fremenkiel/stdout.cv/internal/session"
	"github.com/fremenkiel/stdout.cv/internal/testutil"
)

func TestServiceGetRows(t *testing.T) {
	tests := []struct{
		name,
		tableName					string
		hasSession				bool
		expectedResponse	*Result
		expectedError			error
	} {
		{
			name: "get_rows",
			tableName: "users",
			hasSession: true,
			expectedResponse: &Result{
				Rows: []*Row{
					{
						stringPtr("1"),
						stringPtr("Kevin Wagner"),
					},
				},	
				Columns: []string{
					"id",
					"name",
				},
			},
			expectedError: nil,
		},
		{
			name: "no_rows",
			tableName: "educations",
			hasSession: true,
			expectedResponse: &Result{
				Rows: []*Row{},
				Columns: []string{
					"id",
					"user_id",
					"description",
				},
			},
			expectedError: nil,
		},
		{
			name: "no_session",
			tableName: "users",
			hasSession: false,
			expectedResponse: nil,
			expectedError: session.ErrNoSessionInContext,
		},
	}

	for ti := range tests {
		test := tests[ti]

		t.Run(test.name, func(t *testing.T) {
			cache := session.NewCache(database.TestDatabaseFileNameTemplate)

			ctx := t.Context()

			if test.hasSession {
				sessionService := session.NewService(cache)

				sessionId, err := sessionService.CreateSession()
				if err != nil {
					t.Fatal(err)
				}

				ctx = session.NewContext(t.Context(), sessionId)
				defer testutil.CleanupSession(ctx)
			}

			service := NewService(
				NewRepository(cache),
				schema.NewService(schema.NewRepository(cache)),
				)

			response, err := service.GetRows(ctx, fmt.Sprintf("SELECT * FROM %s;", test.tableName))
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

func TestParseQuery(t *testing.T) {
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

	for ti := range tests {
		test := tests[ti]

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
	tableSchema	[]*column.Column
}

func NewTestSchemaService(tableName string, tableSchemaSlice []string) *testSchemaService {
	tableSchema := make([]*column.Column, len(tableSchemaSlice))

	for i, col := range tableSchemaSlice {
		tableSchema[i] = &column.Column{
			Name: col,
		}
	}

	return &testSchemaService{
		tableName: tableName,
		tableSchema: tableSchema,
	}
}

func (s *testSchemaService) GetSchema(ctx context.Context, name string) ([]*column.Column, error) {
	if name != s.tableName {
		return nil, fmt.Errorf("Unknown table, %s got %s", s.tableName, name)
	}

	return s.tableSchema, nil
}

func stringPtr(s string) *string {
	return &s
}

func TestServiceGetRowDuration(t *testing.T) {
	tests := []struct{
		name,
		query							string
		expectedResponse	time.Duration
		expectedError			error
	} {
		{
			name: "get_rows",
			expectedResponse: time.Duration(10 * time.Microsecond),
			expectedError: nil,
		},
	}

	for ti := range tests {
		test := tests[ti]

		t.Run(test.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
			ctx := t.Context()
			service := NewService(
				&testRepository{},
				nil,
				)

			response, err := service.GetRows(ctx, test.query)
			if test.expectedError != err {
				t.Fatalf("unexpected error, %v got %v", test.expectedError, err)
			}

			if test.expectedError != nil && err != nil {
				return
			}

			if test.expectedResponse.Microseconds() != response.Duration {
				t.Fatalf("unexpected duration, %d got %d", test.expectedResponse.Microseconds(), response.Duration)
			}
			})
		})
	}
}

type testRepository struct {}

func (r *testRepository) GetRows(ctx context.Context, queryString string, query *query.Query) ([]*Row, error) {
	time.Sleep(10 * time.Microsecond)

	return nil, nil
}
