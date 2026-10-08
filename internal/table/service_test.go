package table

import (
	"testing"

	"github.com/fremenkiel/stdout.cv/internal/column"
	"github.com/fremenkiel/stdout.cv/internal/platform/database"
	"github.com/fremenkiel/stdout.cv/internal/schema"
	"github.com/fremenkiel/stdout.cv/internal/session"
	"github.com/fremenkiel/stdout.cv/internal/testutil"
)

type table = map[string]column.Column

func TestGetTables(t *testing.T) {
	tests := []struct{
		name							string
		hasSession				bool
		expectedResponse	map[string]table
		expectedError			error
	} {
		{
			name: "get_tables",
			hasSession: true,
			expectedResponse: map[string]table{
				"educations": map[string]column.Column{
					"id": column.Column{
						Type: "INTEGER",
						IsNullable: false,
						DefaultValue: "",
						IsPrimaryKey: true,
						IsForeignKey: false,
					},
					"user_id": column.Column{
						Type: "INTEGER",
						IsNullable: false,
						DefaultValue: "",
						IsPrimaryKey: false,
						IsForeignKey: true,
					},
					"description": column.Column{
						Type: "TEXT",
						IsNullable: false,
						DefaultValue: "",
						IsPrimaryKey: false,
						IsForeignKey: false,
					},
				},
				"users": map[string]column.Column{
					"id": column.Column{
						Type: "INTEGER",
						IsNullable: false,
						DefaultValue: "",
						IsPrimaryKey: true,
						IsForeignKey: false,
					},
					"name": column.Column{
						Type: "TEXT",
						IsNullable: true,
						DefaultValue: "",
						IsPrimaryKey: false,
						IsForeignKey: false,
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "no_session",
			hasSession: false,
			expectedResponse: nil,
			expectedError: session.ErrNoSessionInContext,
		},
	}

	for _, test := range tests {
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

			response, err := service.GetTables(ctx)
			if test.expectedError != err {
				t.Fatalf("unexpected error, %v got %v", test.expectedError, err)
			}

			if len(test.expectedResponse) != len(response) {
				t.Fatalf("unexpected response len, %d got %d", len(test.expectedResponse), len(response))
			}

			for _, table := range response {
				expectedTable, ok := test.expectedResponse[table.Name]
				if !ok {
					t.Fatalf("unexpected table, got %s", table.Name)
				}

				for _, col := range table.Columns {
					expectedColumn, ok := expectedTable[col.Name]
					if !ok {
						t.Fatalf("unexpected column, got %s", col.Name)
					}

					if expectedColumn.Type != col.Type {
						t.Fatalf("unexpected type for %s, %s got %s", col.Name, expectedColumn.Type, col.Type)
					}

					if expectedColumn.IsNullable != col.IsNullable {
						t.Fatalf("unexpected nullable state for %s, %v got %v", col.Name, expectedColumn.IsNullable, col.IsNullable)
					}

					if expectedColumn.DefaultValue != col.DefaultValue {
						t.Fatalf("unexpected default value for %s, '%s' got '%s'", col.Name, expectedColumn.DefaultValue, col.DefaultValue)
					}

					if expectedColumn.IsPrimaryKey != col.IsPrimaryKey {
						t.Fatalf("unexpected primary key state for %s, %v got %v", col.Name, expectedColumn.IsPrimaryKey, col.IsPrimaryKey)
					}

					if expectedColumn.IsForeignKey != col.IsForeignKey {
						t.Fatalf("unexpected foreign key state for %s, %v got %v", col.Name, expectedColumn.IsForeignKey, col.IsForeignKey)
					}
				}
			}
		})
	}
}
