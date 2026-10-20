package interpreter

import "testing"

func TestNewFieldNode(t *testing.T) {
}

func TestNewExpressionNode(t *testing.T) {
}

func TestNewOperationNode(t *testing.T) {
	tests := []struct{
		name							string
		tokens						[]Token
		expectedResponse	string
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
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := &Parser{
				tokens: test.tokens,
				r: 0,
			}

			node, err := newOperatorNode(parser)
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

func TestNewInsertNode(t *testing.T) {
}

func TestNewSelectNode(t *testing.T) {
	tests := []struct{
		name							string
		tokens						[]Token
		expectedResponse  string
		expectedColumns		[]string
		expectedTable			string
		expectedError     error
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
			expectedTable: "users",
			expectedError: nil,
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
			expectedTable: "users",
			expectedError: nil,
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
			expectedTable: "users",
			expectedError: nil,
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
				"name",
			},
			expectedTable: "users",
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := &Parser{
				tokens: test.tokens,
				r: 0,
			}

			response, err := newSelectNode(parser)
			if test.expectedError != err {
				t.Fatalf("unexpected err, %v got %v", test.expectedError, err)
			}

			if test.expectedError != nil {
				return
			}

			if str := response.String(); test.expectedResponse != str {
				t.Fatalf("unexpected response, %s got %s", test.expectedResponse, str)
			}

			if str := response.Table().String(); test.expectedTable != str {
				t.Fatalf("unexpected table name, %s got %s", test.expectedTable, str)
			}

			cols := response.Columns()
			for i := range cols {
				if str := cols[i].String(); test.expectedColumns[i] != str {
					t.Fatalf("unexpected col, %s got %s", test.expectedColumns[i], str)
				}
			}
		})
	}
}

func TestNewUpdateNode(t *testing.T) {
}

func TestNewDeleteNode(t *testing.T) {
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
