package interpreter

import (
	"errors"
	"testing"
)

func TestNewInsertNode(t *testing.T) {
	tests := []struct{
		name							string
		tokens						[]Token
		expectedResponse  string
		expectedColumns		[]string
		expectedTable			Node
		expectedErrorCode	ErrorCode
	} {
		{
			name: "insert",
			tokens: []Token{
				{ Type: INSERT, Text: "INSERT" },
				{ Type: INTO, Text: "INTO" },
				{ Type: FIELD, Text: "users" },
				{ Type: L_PARAN, Text: "(" },
				{ Type: FIELD, Text: "name" },
				{ Type: FIELD, Text: "age" },
				{ Type: R_PARAN, Text: ")" },
				{ Type: VALUES, Text: "VALUES" },
				{ Type: L_PARAN, Text: "(" },
				{ Type: STRING, Text: "kevin" },
				{ Type: NUMBER, Text: "25" },
				{ Type: R_PARAN, Text: ")" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "INSERT INTO users (name, age) VALUES ('kevin', 25);",
			expectedColumns: []string{
				"name",
				"age",
			},
			expectedTable: &FieldNode{ baseNode{ Token{ Type: FIELD, Text: "users" } } },
			expectedErrorCode: -1,
		},
		{
			name: "no_into",
			tokens: []Token{
				{ Type: INSERT, Text: "INSERT" },
				{ Type: FIELD, Text: "users" },
				{ Type: L_PARAN, Text: "(" },
				{ Type: FIELD, Text: "name" },
				{ Type: FIELD, Text: "age" },
				{ Type: R_PARAN, Text: ")" },
				{ Type: VALUES, Text: "VALUES" },
				{ Type: L_PARAN, Text: "(" },
				{ Type: STRING, Text: "kevin" },
				{ Type: NUMBER, Text: "25" },
				{ Type: R_PARAN, Text: ")" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "",
			expectedColumns: nil,
			expectedTable: nil,
			expectedErrorCode: INVALID_INSERT_STATEMENT,
		},
		{
			name: "no_semicol",
			tokens: []Token{
				{ Type: INSERT, Text: "INSERT" },
				{ Type: INTO, Text: "INTO" },
				{ Type: FIELD, Text: "users" },
				{ Type: L_PARAN, Text: "(" },
				{ Type: FIELD, Text: "name" },
				{ Type: FIELD, Text: "age" },
				{ Type: R_PARAN, Text: ")" },
				{ Type: VALUES, Text: "VALUES" },
				{ Type: L_PARAN, Text: "(" },
				{ Type: STRING, Text: "kevin" },
				{ Type: NUMBER, Text: "25" },
				{ Type: R_PARAN, Text: ")" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "",
			expectedColumns: nil,
			expectedTable: nil,
			expectedErrorCode: INVALID_END_OF_STATEMENT,
		},
		{
			name: "select",
			tokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: FIELD, Text: "name" },
				{ Type: FROM, Text: "FROM" },
				{ Type: FIELD, Text: "users" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "",
			expectedColumns: nil,
			expectedTable: nil,
			expectedErrorCode: INVALID_INSERT_STATEMENT,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := &Parser{
				tokens: test.tokens,
				r: 0,
			}

			response, err := newInsetNode(parser)
			typedErr, ok := errors.AsType[*Error](err)
			if err != nil && !ok {
				t.Fatalf("unexpected error, %v", err)
			}

			if test.expectedErrorCode != -1 {
				if !ok {
					t.Fatalf("unexpected error type, %v", err)
				}

				if test.expectedErrorCode != typedErr.code {
					t.Fatalf("unexpected error, %v got %v (%s)", test.expectedErrorCode, typedErr.code, typedErr.msg)
				}
				return
			}

			if str := response.String(); test.expectedResponse != str {
				t.Fatalf("unexpected response, %s got %s", test.expectedResponse, str)
			}

			cols := response.Columns()
			for i := range cols {
				if str := cols[i]; test.expectedColumns[i] != str {
					t.Fatalf("unexpected col, %s got %s", test.expectedColumns[i], str)
				}
			}

			if response.Table() == nil && test.expectedTable == nil {
				return
			}

			if str := response.Table().String(); test.expectedTable.String() != str {
				t.Fatalf("unexpected table name, %s got %s", test.expectedTable.String(), str)
			}

		})
	}
}

func TestNewSelectNode(t *testing.T) {
	tests := []struct{
		name							string
		tokens						[]Token
		expectedResponse  string
		expectedColumns		[]string
		expectedTable			Node
		expectedErrorCode ErrorCode
	} {
		{
			name: "select",
			tokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: FIELD, Text: "name" },
				{ Type: FROM, Text: "FROM" },
				{ Type: FIELD, Text: "users" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "SELECT name FROM users;",
			expectedColumns: []string{
				"name",
			},
			expectedTable: &FieldNode{ baseNode{ Token{ Type: FIELD, Text: "users" } } },
			expectedErrorCode: -1,
		},
		{
			name: "where",
			tokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: FIELD, Text: "name" },
				{ Type: FROM, Text: "FROM" },
				{ Type: FIELD, Text: "users" },
				{ Type: WHERE, Text: "WHERE" },
				{ Type: FIELD, Text: "id" },
				{ Type: EQUAL, Text: "=" },
				{ Type: NUMBER, Text: "1" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "SELECT name FROM users WHERE id = 1;",
			expectedColumns: []string{
				"name",
			},
			expectedTable: &FieldNode{ baseNode{ Token{ Type: FIELD, Text: "users" } } },
			expectedErrorCode: -1,
		},
		{
			name: "wildcard",
			tokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: FIELD, Text: "name" },
				{ Type: FIELD, Text: "*" },
				{ Type: FROM, Text: "FROM" },
				{ Type: FIELD, Text: "users" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "SELECT name, * FROM users;",
			expectedColumns: []string{
				"name",
				"*",
			},
			expectedTable: &FieldNode{ baseNode{ Token{ Type: FIELD, Text: "users" } } },
			expectedErrorCode: -1,
		},
		{
			name: "as",
			tokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: FIELD, Text: "name" },
				{ Type: AS, Text: "AS" },
				{ Type: FIELD, Text: "first_name" },
				{ Type: FROM, Text: "FROM" },
				{ Type: FIELD, Text: "users" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "SELECT name AS first_name FROM users;",
			expectedColumns: []string{
				"first_name",
			},
			expectedTable: &FieldNode{ baseNode{ Token{ Type: FIELD, Text: "users" } } },
			expectedErrorCode: -1,
		},
		{
			name: "func",
			tokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: FUNC, Text: "now" },
				{ Type: L_PARAN, Text: "(" },
				{ Type: R_PARAN, Text: ")" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "SELECT now();",
			expectedColumns: []string{
				"now()",
			},
			expectedTable: nil,
			expectedErrorCode: -1,
		},
		{
			name: "func_as",
			tokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: FUNC, Text: "now" },
				{ Type: L_PARAN, Text: "(" },
				{ Type: R_PARAN, Text: ")" },
				{ Type: AS, Text: "AS" },
				{ Type: FIELD, Text: "today" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "SELECT now() AS today;",
			expectedColumns: []string{
				"today",
			},
			expectedTable: nil,
			expectedErrorCode: -1,
		},
		{
			name: "table_as",
			tokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: FIELD, Text: "name" },
				{ Type: FROM, Text: "FROM" },
				{ Type: FIELD, Text: "users" },
				{ Type: AS, Text: "AS" },
				{ Type: FIELD, Text: "u" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "SELECT name FROM users AS u;",
			expectedColumns: []string{
				"name",
			},
			expectedTable: &FieldNode{ baseNode{ Token{ Type: FIELD, Text: "users" } } },
			expectedErrorCode: -1,
		},
		{
			name: "update",
			tokens: []Token{
				{ Type: UPDATE, Text: "UPDATE" },
				{ Type: FIELD, Text: "users" },
				{ Type: SET, Text: "SET" },
				{ Type: FIELD, Text: "name" },
				{ Type: EQUAL, Text: "=" },
				{ Type: FIELD, Text: "kevin" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "",
			expectedColumns: nil,
			expectedTable: nil,
			expectedErrorCode: INVALID_SELECT_STATEMENT,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := &Parser{
				tokens: test.tokens,
				r: 0,
			}

			response, err := newSelectNode(parser)
			typedErr, ok := errors.AsType[*Error](err)
			if err != nil && !ok {
				t.Fatalf("unexpected error, %v", err)
			}

			if test.expectedErrorCode != -1 {
				if !ok {
					t.Fatalf("unexpected error type, %v", err)
				}

				if test.expectedErrorCode != typedErr.code {
					t.Fatalf("unexpected error, %v got %v (%s)", test.expectedErrorCode, typedErr.code, typedErr.msg)
				}
				return
			}

			if str := response.String(); test.expectedResponse != str {
				t.Fatalf("unexpected response, %s got %s", test.expectedResponse, str)
			}

			cols := response.Columns()
			for i := range cols {
				if str := cols[i]; test.expectedColumns[i] != str {
					t.Fatalf("unexpected col, %s got %s", test.expectedColumns[i], str)
				}
			}

			if response.Table() == nil && test.expectedTable == nil {
				return
			}

			if str := response.Table().String(); test.expectedTable.String() != str {
				t.Fatalf("unexpected table name, %s got %s", test.expectedTable.String(), str)
			}

		})
	}
}

func TestNewUpdateNode(t *testing.T) {
	tests := []struct{
		name							string
		tokens						[]Token
		expectedResponse  string
		expectedColumns		[]string
		expectedTable			Node
		expectedErrorCode	ErrorCode
	} {
		{
			name: "update",
			tokens: []Token{
				{ Type: UPDATE, Text: "UPDATE" },
				{ Type: FIELD, Text: "users" },
				{ Type: SET, Text: "SET" },
				{ Type: FIELD, Text: "name" },
				{ Type: EQUAL, Text: "=" },
				{ Type: STRING, Text: "kevin" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "UPDATE users SET name = 'kevin';",
			expectedColumns: []string{
				"name",
			},
			expectedTable: &FieldNode{ baseNode{ Token{ Type: FIELD, Text: "users" } } },
			expectedErrorCode: -1,
		},
		{
			name: "where",
			tokens: []Token{
				{ Type: UPDATE, Text: "UPDATE" },
				{ Type: FIELD, Text: "users" },
				{ Type: SET, Text: "SET" },
				{ Type: FIELD, Text: "name" },
				{ Type: EQUAL, Text: "=" },
				{ Type: STRING, Text: "kevin" },
				{ Type: WHERE, Text: "WHERE" },
				{ Type: FIELD, Text: "id" },
				{ Type: EQUAL, Text: "=" },
				{ Type: NUMBER, Text: "1" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "UPDATE users SET name = 'kevin' WHERE id = 1;",
			expectedColumns: []string{
				"name",
			},
			expectedTable: &FieldNode{ baseNode{ Token{ Type: FIELD, Text: "users" } } },
			expectedErrorCode: -1,
		},
		{
			name: "no_set",
			tokens: []Token{
				{ Type: UPDATE, Text: "UPDATE" },
				{ Type: FIELD, Text: "users" },
				{ Type: FIELD, Text: "name" },
				{ Type: EQUAL, Text: "=" },
				{ Type: STRING, Text: "kevin" },
				{ Type: WHERE, Text: "WHERE" },
				{ Type: FIELD, Text: "id" },
				{ Type: EQUAL, Text: "=" },
				{ Type: NUMBER, Text: "1" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "",
			expectedColumns: nil,
			expectedTable: nil,
			expectedErrorCode: INVALID_UPDATE_STATEMENT,
		},
		{
			name: "select",
			tokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: FIELD, Text: "name" },
				{ Type: FROM, Text: "FROM" },
				{ Type: FIELD, Text: "users" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "",
			expectedColumns: nil,
			expectedTable: nil,
			expectedErrorCode: INVALID_UPDATE_STATEMENT,
		},
		{
			name: "multiple_set",
			tokens: []Token{
				{ Type: UPDATE, Text: "UPDATE" },
				{ Type: FIELD, Text: "users" },
				{ Type: SET, Text: "SET" },
				{ Type: FIELD, Text: "name" },
				{ Type: EQUAL, Text: "=" },
				{ Type: STRING, Text: "kevin" },
				{ Type: FIELD, Text: "email" },
				{ Type: EQUAL, Text: "=" },
				{ Type: STRING, Text: "test@gmail.dk" },
				{ Type: FIELD, Text: "age" },
				{ Type: EQUAL, Text: "=" },
				{ Type: NUMBER, Text: "25" },
				{ Type: WHERE, Text: "WHERE" },
				{ Type: FIELD, Text: "id" },
				{ Type: EQUAL, Text: "=" },
				{ Type: NUMBER, Text: "1" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "UPDATE users SET name = 'kevin', email = 'test@gmail.dk', age = 25 WHERE id = 1;",
			expectedColumns: []string{
				"name",
			},
			expectedTable: &FieldNode{ baseNode{ Token{ Type: FIELD, Text: "users" } } },
			expectedErrorCode: -1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := &Parser{
				tokens: test.tokens,
				r: 0,
			}

			response, err := newUpdateNode(parser)
			typedErr, ok := errors.AsType[*Error](err)
			if err != nil && !ok {
				t.Fatalf("unexpected error, %v", err)
			}

			if test.expectedErrorCode != -1 {
				if !ok {
					t.Fatalf("unexpected error type, %v", err)
				}

				if test.expectedErrorCode != typedErr.code {
					t.Fatalf("unexpected error, %v got %v (%s)", test.expectedErrorCode, typedErr.code, typedErr.msg)
				}
				return
			}

			if str := response.String(); test.expectedResponse != str {
				t.Fatalf("unexpected response, %s got %s", test.expectedResponse, str)
			}

			cols := response.Columns()
			for i := range cols {
				if str := cols[i]; test.expectedColumns[i] != str {
					t.Fatalf("unexpected col, %s got %s", test.expectedColumns[i], str)
				}
			}

			if response.Table() == nil && test.expectedTable == nil {
				return
			}

			if str := response.Table().String(); test.expectedTable.String() != str {
				t.Fatalf("unexpected table name, %s got %s", test.expectedTable.String(), str)
			}

		})
	}
}

func TestNewDeleteNode(t *testing.T) {
	tests := []struct{
		name							string
		tokens						[]Token
		expectedResponse  string
		expectedColumns		[]string
		expectedTable			Node
		expectedErrorCode	ErrorCode
	} {
		{
			name: "delete",
			tokens: []Token{
				{ Type: DELETE, Text: "DELETE" },
				{ Type: FROM, Text: "FROM" },
				{ Type: FIELD, Text: "users" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "DELETE FROM users;",
			expectedColumns: []string{},
			expectedTable: &FieldNode{ baseNode{ Token{ Type: FIELD, Text: "users" } } },
			expectedErrorCode: -1,
		},
		{
			name: "where",
			tokens: []Token{
				{ Type: DELETE, Text: "DELETE" },
				{ Type: FROM, Text: "FROM" },
				{ Type: FIELD, Text: "users" },
				{ Type: WHERE, Text: "WHERE" },
				{ Type: FIELD, Text: "id" },
				{ Type: EQUAL, Text: "=" },
				{ Type: NUMBER, Text: "1" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "DELETE FROM users WHERE id = 1;",
			expectedColumns: []string{},
			expectedTable: &FieldNode{ baseNode{ Token{ Type: FIELD, Text: "users" } } },
			expectedErrorCode: -1,
		},
		{
			name: "no_from",
			tokens: []Token{
				{ Type: DELETE, Text: "DELETE" },
				{ Type: FIELD, Text: "users" },
				{ Type: WHERE, Text: "WHERE" },
				{ Type: FIELD, Text: "id" },
				{ Type: EQUAL, Text: "=" },
				{ Type: NUMBER, Text: "1" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "",
			expectedColumns: nil,
			expectedTable: nil,
			expectedErrorCode: INVALID_DELETE_STATEMENT,
		},
		{
			name: "select",
			tokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: FIELD, Text: "name" },
				{ Type: FROM, Text: "FROM" },
				{ Type: FIELD, Text: "users" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "",
			expectedColumns: nil,
			expectedTable: nil,
			expectedErrorCode: INVALID_DELETE_STATEMENT,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := &Parser{
				tokens: test.tokens,
				r: 0,
			}

			response, err := newDeleteNode(parser)
			typedErr, ok := errors.AsType[*Error](err)
			if err != nil && !ok {
				t.Fatalf("unexpected error, %v", err)
			}

			if test.expectedErrorCode != -1 {
				if !ok {
					t.Fatalf("unexpected error type, %v", err)
				}

				if test.expectedErrorCode != typedErr.code {
					t.Fatalf("unexpected error, %v got %v (%s)", test.expectedErrorCode, typedErr.code, typedErr.msg)
				}
				return
			}

			if str := response.String(); test.expectedResponse != str {
				t.Fatalf("unexpected response, %s got %s", test.expectedResponse, str)
			}

			cols := response.Columns()
			for i := range cols {
				if str := cols[i]; test.expectedColumns[i] != str {
					t.Fatalf("unexpected col, %s got %s", test.expectedColumns[i], str)
				}
			}

			if response.Table() == nil && test.expectedTable == nil {
				return
			}

			if str := response.Table().String(); test.expectedTable.String() != str {
				t.Fatalf("unexpected table name, %s got %s", test.expectedTable.String(), str)
			}

		})
	}
}

func TestNewFieldNode(t *testing.T) {
	tests := []struct{
		name							string
		tokens						[]Token
		expectedResponse  string
		expectedColumns		[]string
		expectedTable			Node
		expectedErrorCode	ErrorCode
	} {
		{
			name: "field",
			tokens: []Token{
				{ Type: FIELD, Text: "name" },
				{ Type: FROM, Text: "FROM" },
				{ Type: FIELD, Text: "users" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "name",
			expectedColumns: []string{
				"name",
			},
			expectedTable: &FieldNode{ baseNode{ Token{ Type: FIELD, Text: "name" } } },
			expectedErrorCode: -1,
		},
		{
			name: "field_as",
			tokens: []Token{
				{ Type: FIELD, Text: "name" },
				{ Type: AS, Text: "AS" },
				{ Type: FIELD, Text: "first_name" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "name AS first_name",
			expectedColumns: []string{
				"first_name",
			},
			expectedTable: &FieldNode{ baseNode{ Token{ Type: FIELD, Text: "name" } } },
			expectedErrorCode: -1,
		},
		{
			name: "string",
			tokens: []Token{
				{ Type: STRING, Text: "kevin" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "'kevin'",
			expectedColumns: []string{
				"'kevin'",
			},
			expectedTable: &FieldNode{ baseNode{ Token{ Type: FIELD, Text: "'kevin'" } } },
			expectedErrorCode: -1,
		},
		{
			name: "number",
			tokens: []Token{
				{ Type: NUMBER, Text: "1" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "1",
			expectedColumns: []string{
				"1",
			},
			expectedTable: &FieldNode{ baseNode{ Token{ Type: FIELD, Text: "1" } } },
			expectedErrorCode: -1,
		},
		{
			name: "number",
			tokens: []Token{
				{ Type: NUMBER, Text: "1" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "1",
			expectedColumns: []string{
				"1",
			},
			expectedTable: &FieldNode{ baseNode{ Token{ Type: FIELD, Text: "1" } } },
			expectedErrorCode: -1,
		},
		{
			name: "func",
			tokens: []Token{
				{ Type: FUNC, Text: "now" },
				{ Type: L_PARAN, Text: "(" },
				{ Type: R_PARAN, Text: ")" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "now()",
			expectedColumns: []string{
				"now()",
			},
			expectedTable: nil,
			expectedErrorCode: -1,
		},
		{
			name: "func_as",
			tokens: []Token{
				{ Type: FUNC, Text: "now" },
				{ Type: L_PARAN, Text: "(" },
				{ Type: R_PARAN, Text: ")" },
				{ Type: AS, Text: "AS" },
				{ Type: FIELD, Text: "time" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "now() AS time",
			expectedColumns: []string{
				"time",
			},
			expectedTable: &FieldNode{ baseNode{ Token{ Type: FIELD, Text: "now()" } } },
			expectedErrorCode: -1,
		},
		{
			name: "SELECT",
			tokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "",
			expectedColumns: nil,
			expectedTable: nil,
			expectedErrorCode: INVALID_FIELD,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := &Parser{
				tokens: test.tokens,
				r: 0,
			}

			response, err := newFieldNode(parser)

			typedErr, ok := errors.AsType[*Error](err)
			if err != nil && !ok {
				t.Fatalf("unexpected error, %v", err)
			}

			if test.expectedErrorCode != -1 {
				if !ok {
					t.Fatalf("unexpected error type, %v", err)
				}

				if test.expectedErrorCode != typedErr.code {
					t.Fatalf("unexpected error, %v got %v (%s)", test.expectedErrorCode, typedErr.code, typedErr.msg)
				}
				return
			}

			if str := response.String(); test.expectedResponse != str {
				t.Fatalf("unexpected response, %s got %s", test.expectedResponse, str)
			}

			cols := response.Columns()
			for i := range cols {
				if str := cols[i]; test.expectedColumns[i] != str {
					t.Fatalf("unexpected col, %s got %s", test.expectedColumns[i], str)
				}
			}

			if response.Table() == nil && test.expectedTable == nil {
				return
			}

			if str := response.Table().String(); test.expectedTable.String() != str {
				t.Fatalf("unexpected table name, %s got %s", test.expectedTable.String(), str)
			}
		})
	}
}

func TestNewExpressionNode(t *testing.T) {
	tests := []struct{
		name							string
		tokens						[]Token
		expectedResponse  string
		expetedPrecedense	int
		expectedColumns		[]string
		expectedTable			Node
		expectedErrorCode	ErrorCode
	} {
		{
			name: "expression",
			tokens: []Token{
				{ Type: FIELD, Text: "id" },
				{ Type: EQUAL, Text: "=" },
				{ Type: NUMBER, Text: "1" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "id",
			expetedPrecedense: 0,
			expectedColumns: nil,
			expectedTable: nil,
			expectedErrorCode: -1,
		},
		{
			name: "and",
			tokens: []Token{
				{ Type: AND, Text: "AND" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "AND",
			expetedPrecedense: 2,
			expectedColumns: nil,
			expectedTable: nil,
			expectedErrorCode: -1,
		},
		{
			name: "semicol",
			tokens: []Token{
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: ";",
			expetedPrecedense: -1,
			expectedColumns: nil,
			expectedTable: nil,
			expectedErrorCode: -1,
		},
		{
			name: "invalid",
			tokens: []Token{
				{ Type: TokenType("invalid"), Text: "invalid" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "",
			expetedPrecedense: 0,
			expectedColumns: nil,
			expectedTable: nil,
			expectedErrorCode: INVALID_TYPE_FOR_PRECEDENCE,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := &Parser{
				tokens: test.tokens,
				r: 0,
			}

			response, err := newExpressionNode(parser)
			typedErr, ok := errors.AsType[*Error](err)
			if err != nil && !ok {
				t.Fatalf("unexpected error, %v", err)
			}

			if test.expectedErrorCode != -1 {
				if !ok {
					t.Fatalf("unexpected error type, %v", err)
				}

				if test.expectedErrorCode != typedErr.code {
					t.Fatalf("unexpected error, %v got %v (%s)", test.expectedErrorCode, typedErr.code, typedErr.msg)
				}
				return
			}

			if str := response.String(); test.expectedResponse != str {
				t.Fatalf("unexpected response, %s got %s", test.expectedResponse, str)
			}

			cols := response.Columns()
			for i := range cols {
				if str := cols[i]; test.expectedColumns[i] != str {
					t.Fatalf("unexpected col, %s got %s", test.expectedColumns[i], str)
				}
			}

			if response.Table() == nil && test.expectedTable == nil {
				return
			}

			if str := response.Table().String(); test.expectedTable.String() != str {
				t.Fatalf("unexpected table name, %s got %s", test.expectedTable.String(), str)
			}
		})
	}
}

func TestNewOperationNode(t *testing.T) {
	tests := []struct{
		name							string
		tokens						[]Token
		expectedResponse	string
		expectedErrorCode	ErrorCode
	} {
		{
			name: "is_not_null",
			tokens: []Token{
				{ Type: FIELD, Text: "age" },
				{ Type: IS, Text: "IS" },
				{ Type: NOT, Text: "NOT" },
				{ Type: NULL, Text: "NULL" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "age IS NOT NULL",
			expectedErrorCode: -1,
		},
		{
			name: "equal",
			tokens: []Token{
				{ Type: FIELD, Text: "age" },
				{ Type: EQUAL, Text: "=" },
				{ Type: NUMBER, Text: "5" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "age = 5",
			expectedErrorCode: -1,
		},
		{
			name: "as",
			tokens: []Token{
				{ Type: FIELD, Text: "age" },
				{ Type: AS, Text: "AS" },
				{ Type: FIELD, Text: "date" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "age AS date",
			expectedErrorCode: -1,
		},
		{
			name: "as_n_equal",
			tokens: []Token{
				{ Type: FIELD, Text: "age" },
				{ Type: AS, Text: "AS" },
				{ Type: FIELD, Text: "date" },
				{ Type: N_EQUAL, Text: "!=" },
				{ Type: STRING, Text: "kevin" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "age AS date != 'kevin'",
			expectedErrorCode: -1,
		},
		{
			name: "no_prefix",
			tokens: []Token{
				{ Type: N_EQUAL, Text: "!=" },
				{ Type: STRING, Text: "kevin" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "!=",
			expectedErrorCode: -1,
		},
		{
			name: "no_infix",
			tokens: []Token{
				{ Type: FIELD, Text: "name" },
				{ Type: STRING, Text: "kevin" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "name 'kevin'",
			expectedErrorCode: -1,
		},
		{
			name: "no_suffix",
			tokens: []Token{
				{ Type: FIELD, Text: "name" },
				{ Type: EQUAL, Text: "=" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "name =",
			expectedErrorCode: -1,
		},
		{
			name: "invalid",
			tokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: FIELD, Text: "name" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "",
			expectedErrorCode: INVALID_OPERATOR_TOKEN,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := &Parser{
				tokens: test.tokens,
				r: 0,
			}

			node, err := newOperatorNode(parser)
			typedErr, ok := errors.AsType[*Error](err)
			if err != nil && !ok {
				t.Fatalf("unexpected error, %v", err)
			}

			if test.expectedErrorCode != -1 {
				if !ok {
					t.Fatalf("unexpected error type, %v", err)
				}

				if test.expectedErrorCode != typedErr.code {
					t.Fatalf("unexpected error, %v got %v (%s)", test.expectedErrorCode, typedErr.code, typedErr.msg)
				}
				return
			}

			response := node.String()
			if test.expectedResponse != response {
				t.Fatalf("unexpected response, '%s' got '%s'", test.expectedResponse, response)
			}

			if node.Table() != nil {
				t.Fatalf("unexpected table, nil got %v", node.Table())
			}

			if node.Columns() != nil {
				t.Fatalf("unexpected columns, nil got %v", node.Columns())
			}
		})
	}
}

func TestNewCastNode(t *testing.T) {
	tests := []struct{
		name							string
		tokens						[]Token
		expectedResponse  string
		expectedColumns		[]string
		expectedTable			Node
		expectedErrorCode	ErrorCode
	} {
		{
			name: "cast",
			tokens: []Token{
				{ Type: FIELD, Text: "id" },
				{ Type: AS, Text: "AS" },
				{ Type: FIELD, Text: "identifier" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "id AS identifier",
			expectedColumns: []string{
				"identifier",
			},
			expectedTable: &FieldNode{ baseNode{ Token{ Type: FIELD, Text: "id" } } },
			expectedErrorCode: -1,
		},
		{
			name: "no_as",
			tokens: []Token{
				{ Type: FIELD, Text: "id" },
				{ Type: FIELD, Text: "identifier" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "",
			expectedColumns: nil,
			expectedTable: nil,
			expectedErrorCode: INVALID_BASE_NODE_TYPE,
		},
		{
			name: "invalid",
			tokens: []Token{
				{ Type: FIELD, Text: "id" },
				{ Type: AS, Text: "AS" },
				{ Type: STRING, Text: "kevin" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "",
			expectedColumns: nil,
			expectedTable: nil,
			expectedErrorCode: INVALID_BASE_NODE_TYPE,
		},
		{
			name: "non_field",
			tokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: AS, Text: "AS" },
				{ Type: STRING, Text: "kevin" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "",
			expectedColumns: nil,
			expectedTable: nil,
			expectedErrorCode: INVALID_BASE_NODE_TYPE,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := &Parser{
				tokens: test.tokens,
				r: 0,
			}

			fieldNode := FieldNode{}
			fieldNode.Token = parser.currentToken()

			if err := parser.nextToken(); err != nil {
				t.Fatalf("unexpected err while creating field, %v", err)
			}

			response, err := newCastNode(parser, &fieldNode)
			typedErr, ok := errors.AsType[*Error](err)
			if err != nil && !ok {
				t.Fatalf("unexpected error, %v", err)
			}

			if test.expectedErrorCode != -1 {
				if !ok {
					t.Fatalf("unexpected error type, %v", err)
				}

				if test.expectedErrorCode != typedErr.code {
					t.Fatalf("unexpected error, %v got %v (%s)", test.expectedErrorCode, typedErr.code, typedErr.msg)
				}
				return
			}

			if str := response.String(); test.expectedResponse != str {
				t.Fatalf("unexpected response, %s got %s", test.expectedResponse, str)
			}

			cols := response.Columns()
			for i := range cols {
				if str := cols[i]; test.expectedColumns[i] != str {
					t.Fatalf("unexpected col, %s got %s", test.expectedColumns[i], str)
				}
			}

			if response.Table() == nil && test.expectedTable == nil {
				return
			}

			if str := response.Table().String(); test.expectedTable.String() != str {
				t.Fatalf("unexpected table name, %s got %s", test.expectedTable.String(), str)
			}
		})
	}
}

func TestNewFuncNode(t *testing.T) {
	tests := []struct{
		name							string
		tokens						[]Token
		expectedResponse	string
		expectedColumns		[]string
		expectedTable			Node
		expectedErrorCode	ErrorCode
	} {
		{
			name: "func",
			tokens: []Token{
				{ Type: FUNC, Text: "now" },
				{ Type: L_PARAN, Text: "(" },
				{ Type: R_PARAN, Text: ")" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "now()",
			expectedColumns: []string{
				"now()",
			},
			expectedTable: nil,
			expectedErrorCode: -1,
		},
		{
			name: "func_infix",
			tokens: []Token{
				{ Type: FUNC, Text: "CONCAT" },
				{ Type: L_PARAN, Text: "(" },
				{ Type: FIELD, Text: "name" },
				{ Type: STRING, Text: "kevin" },
				{ Type: NUMBER, Text: "12" },
				{ Type: R_PARAN, Text: ")" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "CONCAT(name, 'kevin', 12)",
			expectedColumns: []string{
				"CONCAT(name, 'kevin', 12)",
			},
			expectedTable: nil,
			expectedErrorCode: -1,
		},
		{
			name: "invalid",
			tokens: []Token{
				{ Type: FUNC, Text: "now" },
				{ Type: L_PARAN, Text: "(" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "",
			expectedColumns: nil,
			expectedTable: nil,
			expectedErrorCode: INVALID_FUNC_CHAR,
		},
		{
			name: "select",
			tokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: L_PARAN, Text: "(" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "",
			expectedColumns: nil,
			expectedTable: nil,
			expectedErrorCode: INVALID_BASE_NODE_TYPE,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := &Parser{
				tokens: test.tokens,
				r: 0,
			}

			node, err := newFuncNode(parser)
			typedErr, ok := errors.AsType[*Error](err)
			if err != nil && !ok {
				t.Fatalf("unexpected error, %v", err)
			}

			if test.expectedErrorCode != -1 {
				if !ok {
					t.Fatalf("unexpected error type, %v", err)
				}

				if test.expectedErrorCode != typedErr.code {
					t.Fatalf("unexpected error, %v got %v (%s)", test.expectedErrorCode, typedErr.code, typedErr.msg)
				}
				return
			}

			response := node.String()
			if test.expectedResponse != response {
				t.Fatalf("unexpected response, '%s' got '%s'", test.expectedResponse, response)
			}
		})
	}
}

func TestNewListNode(t *testing.T) {
	tests := []struct{
		name							string
		tokens						[]Token
		expectedResponse	string
		expectedColumns		[]string
		expectedTable			Node
		expectedErrorCode	ErrorCode
	} {
		{
			name: "list",
			tokens: []Token{
				{ Type: L_PARAN, Text: "(" },
				{ Type: FIELD, Text: "name" },
				{ Type: FIELD, Text: "age" },
				{ Type: R_PARAN, Text: ")" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "(name, age)",
			expectedColumns: []string{
				"name",
				"age",
			},
			expectedTable: nil,
			expectedErrorCode: -1,
		},
		{
			name: "string",
			tokens: []Token{
				{ Type: L_PARAN, Text: "(" },
				{ Type: FIELD, Text: "name" },
				{ Type: STRING, Text: "kevin" },
				{ Type: R_PARAN, Text: ")" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "(name, 'kevin')",
			expectedColumns: []string{
				"name",
				"'kevin'",
			},
			expectedTable: nil,
			expectedErrorCode: -1,
		},
		{
			name: "numer",
			tokens: []Token{
				{ Type: L_PARAN, Text: "(" },
				{ Type: FIELD, Text: "name" },
				{ Type: NUMBER, Text: "1" },
				{ Type: R_PARAN, Text: ")" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "(name, 1)",
			expectedColumns: []string{
				"name",
				"1",
			},
			expectedTable: nil,
			expectedErrorCode: -1,
		},
		{
			name: "func",
			tokens: []Token{
				{ Type: L_PARAN, Text: "(" },
				{ Type: FUNC, Text: "now" },
				{ Type: L_PARAN, Text: "(" },
				{ Type: R_PARAN, Text: ")" },
				{ Type: FIELD, Text: "name" },
				{ Type: R_PARAN, Text: ")" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "(now(), name)",
			expectedColumns: []string{
				"now()",
				"name",
			},
			expectedTable: nil,
			expectedErrorCode: -1,
		},
		{
			name: "no_start",
			tokens: []Token{
				{ Type: FIELD, Text: "name" },
				{ Type: R_PARAN, Text: ")" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "",
			expectedColumns: nil,
			expectedTable: nil,
			expectedErrorCode: INVALID_LIST_FORMAT,
		},
		{
			name: "no_end",
			tokens: []Token{
				{ Type: L_PARAN, Text: "(" },
				{ Type: FIELD, Text: "name" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "",
			expectedColumns: nil,
			expectedTable: nil,
			expectedErrorCode: INVALID_LIST_FORMAT,
		},
		{
			name: "select",
			tokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "",
			expectedColumns: nil,
			expectedTable: nil,
			expectedErrorCode: INVALID_LIST_FORMAT,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := &Parser{
				tokens: test.tokens,
				r: 0,
			}

			node, err := newListNode(parser)
			typedErr, ok := errors.AsType[*Error](err)
			if err != nil && !ok {
				t.Fatalf("unexpected error, %v", err)
			}

			if test.expectedErrorCode != -1 {
				if !ok {
					t.Fatalf("unexpected error type, %v", err)
				}

				if test.expectedErrorCode != typedErr.code {
					t.Fatalf("unexpected error, %v got %v (%s)", test.expectedErrorCode, typedErr.code, typedErr.msg)
				}
				return
			}

			response := node.String()
			if test.expectedResponse != response {
				t.Fatalf("unexpected response, '%s' got '%s'", test.expectedResponse, response)
			}
		})
	}
}

func TestNewWhereNode(t *testing.T) {
	tests := []struct{
		name							string
		tokens						[]Token
		expectedResponse	string
	} {
		{
			name: "where",
			tokens: []Token{
				{ Type: WHERE, Text: "WHERE" },
				{ Type: FIELD, Text: "id" },
				{ Type: EQUAL, Text: "=" },
				{ Type: NUMBER, Text: "1" },
				{ Type: AND, Text: "AND" },
				{ Type: FIELD, Text: "name" },
				{ Type: LIKE, Text: "LIKE" },
				{ Type: STRING, Text: "%kev%" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "WHERE id = 1 AND name LIKE '%kev%'",
		},
		{
			name: "where_not",
			tokens: []Token{
				{ Type: WHERE, Text: "WHERE" },
				{ Type: NOT, Text: "NOT" },
				{ Type: FIELD, Text: "id" },
				{ Type: EQUAL, Text: "=" },
				{ Type: NUMBER, Text: "1" },
				{ Type: AND, Text: "AND" },
				{ Type: FIELD, Text: "name" },
				{ Type: LIKE, Text: "LIKE" },
				{ Type: STRING, Text: "%kev%" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "WHERE NOT id = 1 AND name LIKE '%kev%'",
		},
		{
			name: "is_not_null",
			tokens: []Token{
				{ Type: WHERE, Text: "WHERE" },
				{ Type: FIELD, Text: "id" },
				{ Type: IS, Text: "IS" },
				{ Type: NOT, Text: "NOT" },
				{ Type: NULL, Text: "NULL" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "WHERE id IS NOT NULL",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := &Parser{
				tokens: test.tokens,
				r: 0,
			}

			node, err := newWhereNode(parser)
			if err != nil {
				t.Fatal(err)
			}

			response := node.String()
			if test.expectedResponse != response {
				t.Fatalf("unexpected response, %s got %s", test.expectedResponse, response)
			}
		})
	}
}
