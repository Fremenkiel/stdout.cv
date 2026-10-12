package schema

import (
	"testing"

	"github.com/fremenkiel/stdout.cv/internal/column"
	"github.com/fremenkiel/stdout.cv/internal/session"
	"github.com/fremenkiel/stdout.cv/internal/testutil"
)

func TestRepositoryGetSchema(t *testing.T) {
	tests := []struct{
		name,
		tableName					string
		hasSession				bool
		expectedResponse	[]*column.Column
		expectedError			error
	} {
		{
			name: "get_tables",
			tableName: "educations",
			hasSession: true,
			expectedResponse: []*column.Column{
				{
					Name: "id",
					Type: "INTEGER",
					IsNullable: false,
					DefaultValue: "",
					IsPrimaryKey: true,
					IsForeignKey: false,
				},
				{
					Name: "user_id",
					Type: "INTEGER",
					IsNullable: false,
					DefaultValue: "",
					IsPrimaryKey: false,
					IsForeignKey: true,
				},
				{
					Name: "description",
					Type: "TEXT",
					IsNullable: false,
					DefaultValue: "",
					IsPrimaryKey: false,
					IsForeignKey: false,
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

	for ti := range tests {
		test := tests[ti]

		t.Run(test.name, func(t *testing.T) {
			env := testutil.LoadTestEnv(t)
			cache := session.NewCache(env.DatabasePath, env.DatabaseFileTemplate)

			ctx := t.Context()
			if test.hasSession {
				sessionService := session.NewService(cache)
				id, err := sessionService.CreateSession()
				if err != nil {
					t.Fatal(err)
				}

				ctx = session.NewContext(ctx, id)
				defer testutil.CleanupSession(ctx, env.DatabasePath)
			}

			repository := NewRepository(cache)

			response, err := repository.GetSchema(ctx, test.tableName)
			if test.expectedError != err {
				t.Fatalf("unexpected err, %v got %v", test.expectedError, err)
			}

			if test.expectedError != nil && err != nil {
				return
			}

			if len(test.expectedResponse) != len(response) {
				t.Fatalf("unexpected response len, %d got %d", len(test.expectedResponse), len(response))
			}


			for i := range response {
				col := response[i]
				expectedResponse := test.expectedResponse[i]

				if expectedResponse.Name != col.Name {
					t.Fatalf("unexpected column name, %s got %s", expectedResponse.Name, col.Name)
				}

				if expectedResponse.Type != col.Type {
					t.Fatalf("unexpected type for %s, %s got %s", col.Name, expectedResponse.Type, col.Type)
				}

				if expectedResponse.IsNullable != col.IsNullable {
					t.Fatalf("unexpected nullable state for %s, %v got %v", col.Name, expectedResponse.IsNullable, col.IsNullable)
				}

				if expectedResponse.DefaultValue != col.DefaultValue {
					t.Fatalf("unexpected default value for %s, '%s' got '%s'", col.Name, expectedResponse.DefaultValue, col.DefaultValue)
				}

				if expectedResponse.IsPrimaryKey != col.IsPrimaryKey {
					t.Fatalf("unexpected primary key state for %s, %v got %v", col.Name, expectedResponse.IsPrimaryKey, col.IsPrimaryKey)
				}

				if expectedResponse.IsForeignKey != col.IsForeignKey {
					t.Fatalf("unexpected foreign key state for %s, %v got %v", col.Name, expectedResponse.IsForeignKey, col.IsForeignKey)
				}
			}
		})
	}
}
