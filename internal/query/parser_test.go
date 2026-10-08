package query

import (
	"testing"

	"github.com/fremenkiel/stdout.cv/internal/column"
	"github.com/fremenkiel/stdout.cv/internal/testutil"
)

func TestParse(t *testing.T) {
	tests := []struct{
		name							string
		queryString				string
		expectedResponse	*Query
		expectedError			error
	} {
		{
			name: "select_with_fields",
			queryString: "select name, age, email from users where id = 1;",
			expectedResponse: &Query{
				Fields: []string{
					"name",
					"age",
					"email",
				},
				WildcardIndex: -1,
				TableName: "users",
				Type: SELECT,
				Modifiers: []string{
					FROM,
					WHERE,
				},
			},
			expectedError: nil,
		},
		{
			name: "select_with_wildcard",
			queryString: "select * from users where id = 1;",
			expectedResponse: &Query{
				Fields: []string{
					"*",
				},
				WildcardIndex: 0,
				TableName: "users",
				Type: SELECT,
				Modifiers: []string{
					FROM,
					WHERE,
				},
			},
			expectedError: nil,
		},
		{
			name: "select_with_wildcard_and_field",
			queryString: "select name, email, * from users where id = 1;",
			expectedResponse: &Query{
				Fields: []string{
					"name",
					"email",
					"*",
				},
				WildcardIndex: 2,
				TableName: "users",
				Type: SELECT,
				Modifiers: []string{
					FROM,
					WHERE,
				},
			},
			expectedError: nil,
		},
		{
			name: "select_without_space",
			queryString: "select name,email,age from users where id = 1;",
			expectedResponse: &Query{
				Fields: []string{
					"name",
					"email",
					"age",
				},
				WildcardIndex: -1,
				TableName: "users",
				Type: SELECT,
				Modifiers: []string{
					FROM,
					WHERE,
				},
			},
			expectedError: nil,
		},
		{
			name: "select_mixed_case",
			queryString: "sElEcT name, email, age fRoM users wHeRe id = 1;",
			expectedResponse: &Query{
				Fields: []string{
					"name",
					"email",
					"age",
				},
				WildcardIndex: -1,
				TableName: "users",
				Type: SELECT,
				Modifiers: []string{
					FROM,
					WHERE,
				},
			},
			expectedError: nil,
		},
		{
			name: "select_two_table_names",
			queryString: "select name, email, age from files users where id = 1;",
			expectedResponse: &Query{
				Fields: []string{
					"name",
					"email",
					"age",
				},
				WildcardIndex: -1,
				TableName: "files",
				Type: SELECT,
				Modifiers: []string{
					FROM,
					WHERE,
				},
			},
			expectedError: nil,
		},
		{
			name: "update",
			queryString: "update users set name = 'Kim' where id = 1;",
			expectedResponse: &Query{
				Fields: []string{},
				WildcardIndex: -1,
				TableName: "users",
				Type: UPDATE,
				Modifiers: []string{
					SET,
					WHERE,
				},
			},
			expectedError: nil,
		},
		{
			name: "insert",
			queryString: "insert into users (name, age, email) VALUES ('Kim', 25, 'kim@saarikivi.dk');",
			expectedResponse: &Query{
				Fields: []string{},
				WildcardIndex: -1,
				TableName: "users",
				Type: INSERT,
				Modifiers: []string{
					INTO,
					VALUES,
				},
			},
			expectedError: nil,
		},
		{
			name: "delete",
			queryString: "delete from users where id = 1;",
			expectedResponse: &Query{
				Fields: []string{},
				WildcardIndex: -1,
				TableName: "users",
				Type: DELETE,
				Modifiers: []string{
					FROM,
					WHERE,
				},
			},
			expectedError: nil,
		},
		{
			name: "multiple_select",
			queryString: "select * from files where owner_id IN (select id from users where id = 1);",
			expectedResponse: &Query{
				Fields: []string{
					"*",
				},
				WildcardIndex: 0,
				TableName: "files",
				Type: SELECT,
				Modifiers: []string{
					FROM,
					WHERE,
					FROM,
					WHERE,
				},
			},
			expectedError: nil,
		},
		{
			name: "delete_and_select",
			queryString: "delete from files where id = 1; select name from users where id = 1;",
			expectedResponse: nil,
			expectedError: ErrQueryNotAllowed,
		},
		{
			name: "semicolon",
			queryString: "select name; from users where id = 1;",
			expectedResponse: nil,
			expectedError: ErrQueryNotAllowed,
		},
		{
			name: "semicolon_at_end",
			queryString: "select * from users;",
			expectedResponse: &Query{
				Fields: []string{
					"*",
				},
				WildcardIndex: 0,
				TableName: "users",
				Type: SELECT,
				Modifiers: []string{
					FROM,
				},
			},
			expectedError: nil,
		},
		{
			name: "comment",
			queryString: "-- do not process this part\nselect * from users;",
			expectedResponse: &Query{
				Fields: []string{
					"*",
				},
				WildcardIndex: 0,
				TableName: "users",
				Type: SELECT,
				Modifiers: []string{
					FROM,
				},
			},
			expectedError: nil,
		},
		// Will be catched by sqlite
		{
			name: "comment_inside",
			queryString: "select --name from users where id = 1;",
			expectedResponse: &Query{
				Fields: []string{},
				WildcardIndex: -1,
				TableName: "",
				Type: SELECT,
				Modifiers: []string{},
			},
			expectedError: nil,
		},
		{
			name: "comment_space_at_line_break",
			queryString: "-- do not process this part \n select * from users;",
			expectedResponse: &Query{
				Fields: []string{
					"*",
				},
				WildcardIndex: 0,
				TableName: "users",
				Type: SELECT,
				Modifiers: []string{
					FROM,
				},
			},
			expectedError: nil,
		},
		{
			name: "comment_first_space_at_line_break",
			queryString: "-- do not process this part \nselect * from users;",
			expectedResponse: &Query{
				Fields: []string{
					"*",
				},
				WildcardIndex: 0,
				TableName: "users",
				Type: SELECT,
				Modifiers: []string{
					FROM,
				},
			},
			expectedError: nil,
		},
		{
			name: "comment_last_space_at_line_break",
			queryString: "-- do not process this part\n select * from users;",
			expectedResponse: &Query{
				Fields: []string{
					"*",
				},
				WildcardIndex: 0,
				TableName: "users",
				Type: SELECT,
				Modifiers: []string{
					FROM,
				},
			},
			expectedError: nil,
		},
		{
			name: "comment_end",
			queryString: "select * from users;\n-- do not process this part",
			expectedResponse: nil,
			expectedError: ErrQueryNotAllowed,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := Parse(test.queryString)

			if test.expectedError != err {
				t.Fatalf("unexpected error: %v, got %v", test.expectedError, err)
			}

			if test.expectedResponse != nil && response != nil {
				if !testutil.SliceEql(test.expectedResponse.Fields, response.Fields) {
					t.Fatalf("unexpected fields slice: %v, got %v", test.expectedResponse.Fields, response.Fields)
				}

				if test.expectedResponse.WildcardIndex != response.WildcardIndex {
					t.Fatalf("unexpected wildcard index: %d, got %d", test.expectedResponse.WildcardIndex, response.WildcardIndex)
				}

				if test.expectedResponse.TableName != response.TableName {
					t.Fatalf("unexpected table name: %s, got %s", test.expectedResponse.TableName, response.TableName)
				}

				if test.expectedResponse.Type != response.Type {
					t.Fatalf("unexpected type: %s, got %s", test.expectedResponse.Type, response.Type)
				}

				if !testutil.SliceEql(test.expectedResponse.Modifiers, response.Modifiers) {
					t.Fatalf("unexpected modifiers: %v, got %v", test.expectedResponse.Modifiers, response.Modifiers)
				}
			}

			if (test.expectedResponse == nil) != (response == nil) {
					t.Fatalf("unexpected response: %v, got %v", test.expectedResponse, response)
			}
		})
	}
}

func TestAppendWildcardSchema(t *testing.T) {
	tests := []struct{
		name							string
		fields						[]string
		wildcardIndex			int16
		schema						[]*column.Column
		expectedResponse	[]string
		expectedError			error
	} {
		{
			name: "only_wildcard",
			fields: []string{
				"*",
			},
			wildcardIndex: 0,
			schema: []*column.Column{
				&column.Column{ Name: "name" },
				&column.Column{ Name: "age" },
				&column.Column{ Name: "email" },
			},
			expectedResponse: []string{
				"name",
				"age",
				"email",
			},
			expectedError: nil,
		},
		{
			name: "fields_and_wildcard",
			fields: []string{
				"name",
				"*",
				"age",
			},
			wildcardIndex: 1,
			schema: []*column.Column{
				&column.Column{ Name: "name" },
				&column.Column{ Name: "age" },
				&column.Column{ Name: "email" },
			},
			expectedResponse: []string{
				"name",
				"name",
				"age",
				"email",
				"age",
			},
			expectedError: nil,
		},
		{
			name: "no_schema",
			fields: []string{
				"*",
			},
			wildcardIndex: 0,
			schema: []*column.Column{},
			expectedResponse: []string{},
			expectedError: nil,
		},
		{
			name: "wrong_index",
			fields: []string{
				"name",
				"*",
			},
			wildcardIndex: 0,
			schema: []*column.Column{},
			expectedResponse: nil,
			expectedError: ErrFieldNotWildcard,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, err := AppendWildcardSchema(test.wildcardIndex, test.fields, test.schema)

			if test.expectedError != err {
				t.Fatalf("unexpected error: %v, got %v", test.expectedError, err)
			}

			if !testutil.SliceEql(test.expectedResponse, response) {
				t.Fatalf("response do not match, expected %v, got %v", test.expectedResponse, response)
			}
		})
	}
}

func TestParseFields(t *testing.T) {
	tests := []struct{
		name									string
		fieldStr							string
		fieldSlice						[]string
		expectedResponse			[]string
		expectedWindcardIndex	int16
	} {
		{
			name: "without_windcard",
			fieldStr: "first_name,last_name,email,password_hash",
			fieldSlice: []string{},
			expectedResponse: []string{
				"first_name",
				"last_name",
				"email",
				"password_hash",
			},
			expectedWindcardIndex: -1,
		},
		{
			name: "with_windcard",
			fieldStr: "first_name,*,last_name,email,password_hash",
			fieldSlice: []string{},
			expectedResponse: []string{
				"first_name",
				"*",
				"last_name",
				"email",
				"password_hash",
			},
			expectedWindcardIndex: 1,
		},
		{
			name: "without_windcard_append_correct",
			fieldStr: "first_name,last_name,email,password_hash",
			fieldSlice: []string{
				"username",
			},
			expectedResponse: []string{
				"username",
				"first_name",
				"last_name",
				"email",
				"password_hash",
			},
			expectedWindcardIndex: -1,
		},
		{
			name: "with_windcard_append_correct",
			fieldStr: "first_name,*,last_name,email,password_hash",
			fieldSlice: []string{
				"username",
			},
			expectedResponse: []string{
				"username",
				"first_name",
				"*",
				"last_name",
				"email",
				"password_hash",
			},
			expectedWindcardIndex: 2,
		},
		{
			name: "with_windcard_wrong_format",
			fieldStr: "first_name,*last_name,email,password_hash",
			fieldSlice: []string{},
			expectedResponse: []string{
				"first_name",
				"*last_name",
				"email",
				"password_hash",
			},
			expectedWindcardIndex: -1,
		},
		{
			name: "no_comma",
			fieldStr: "first_name*last_nameemailpassword_hash",
			fieldSlice: []string{},
			expectedResponse: []string{
				"first_name*last_nameemailpassword_hash",
			},
			expectedWindcardIndex: -1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response, windcardIndex := parseSelectField(test.fieldSlice, test.fieldStr)

			if !testutil.SliceEql(test.expectedResponse, response) {
				t.Fatalf("response do not match, expected %v, got %v", test.expectedResponse, response)
			}

			if windcardIndex != test.expectedWindcardIndex {
				t.Fatalf("wrong windcard index, expected %d, got %d", test.expectedWindcardIndex, windcardIndex)
			}
		})
	}
}
