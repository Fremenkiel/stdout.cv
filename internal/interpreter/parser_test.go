package interpreter

import "testing"

// func TestParseQuery(t *testing.T) {
// 	tests := []struct{
// 		name						string
// 		query						string
// 		expectedTokens	[]Token
// 	} {
// 		{
// 			name: "insert",
// 			query: "INSERT INTO users (name, email) VALUES ('kevin', 'test@gmail.dk');",
// 			expectedTokens: []Token{
// 				{ Type: INSERT, Text: "INSERT" },
// 				{ Type: INTO, Text: "INTO" },
// 				{ Type: FIELD, Text: "users" },
// 				{ Type: L_PARAN, Text: "(" },
// 				{ Type: FIELD, Text: "name" },
// 				{ Type: FIELD, Text: "email" },
// 				{ Type: R_PARAN, Text: ")" },
// 				{ Type: VALUES, Text: "VALUES" },
// 				{ Type: L_PARAN, Text: "(" },
// 				{ Type: STRING, Text: "kevin" },
// 				{ Type: STRING, Text: "test@gmail.dk" },
// 				{ Type: R_PARAN, Text: ")" },
// 				{ Type: SEMICOL, Text: ";" },
// 				{ Type: EOF, Text: "" },
// 			},
// 		},
// 		{
// 			name: "select",
// 			query: "SELECT name, age, email FROM users WHERE id = 1 AND name LIKE '%kevin%';",
// 			expectedTokens: []Token{
// 				{ Type: SELECT, Text: "SELECT" },
// 				{ Type: FIELD, Text: "name" },
// 				{ Type: FIELD, Text: "age" },
// 				{ Type: FIELD, Text: "email" },
// 				{ Type: FROM, Text: "FROM" },
// 				{ Type: FIELD, Text: "users" },
// 				{ Type: WHERE, Text: "WHERE" },
// 				{ Type: FIELD, Text: "id" },
// 				{ Type: EQUAL, Text: "=" },
// 				{ Type: NUMBER, Text: "1" },
// 				{ Type: AND, Text: "AND" },
// 				{ Type: FIELD, Text: "name" },
// 				{ Type: LIKE, Text: "LIKE" },
// 				{ Type: STRING, Text: "%kevin%" },
// 				{ Type: SEMICOL, Text: ";" },
// 				{ Type: EOF, Text: "" },
// 			},
// 		},
// 		{
// 			name: "update",
// 			query: "UPDATE users SET name = 'john' WHERE id = 1;",
// 			expectedTokens: []Token{
// 				{ Type: UPDATE, Text: "UPDATE" },
// 				{ Type: FIELD, Text: "users" },
// 				{ Type: SET, Text: "SET" },
// 				{ Type: FIELD, Text: "name" },
// 				{ Type: EQUAL, Text: "=" },
// 				{ Type: STRING, Text: "john" },
// 				{ Type: WHERE, Text: "WHERE" },
// 				{ Type: FIELD, Text: "id" },
// 				{ Type: EQUAL, Text: "=" },
// 				{ Type: NUMBER, Text: "1" },
// 				{ Type: SEMICOL, Text: ";" },
// 				{ Type: EOF, Text: "" },
// 			},
// 		},
// 		{
// 			name: "delete",
// 			query: "DELETE FROM users WHERE id = 1 AND name LIKE '%kevin%';",
// 			expectedTokens: []Token{
// 				{ Type: DELETE, Text: "DELETE" },
// 				{ Type: FROM, Text: "FROM" },
// 				{ Type: FIELD, Text: "users" },
// 				{ Type: WHERE, Text: "WHERE" },
// 				{ Type: FIELD, Text: "id" },
// 				{ Type: EQUAL, Text: "=" },
// 				{ Type: NUMBER, Text: "1" },
// 				{ Type: AND, Text: "AND" },
// 				{ Type: FIELD, Text: "name" },
// 				{ Type: LIKE, Text: "LIKE" },
// 				{ Type: STRING, Text: "%kevin%" },
// 				{ Type: SEMICOL, Text: ";" },
// 				{ Type: EOF, Text: "" },
// 			},
// 		},
// 		{
// 			name: "func",
// 			query: "SELECT now();",
// 			expectedTokens: []Token{
// 				{ Type: SELECT, Text: "SELECT" },
// 				{ Type: FUNC, Text: "now" },
// 				{ Type: L_PARAN, Text: "(" },
// 				{ Type: R_PARAN, Text: ")" },
// 				{ Type: SEMICOL, Text: ";" },
// 				{ Type: EOF, Text: "" },
// 			},
// 		},
// 		{
// 			name: "is_not_null",
// 			query: "SELECT name FROM users WHERE email IS NOT NULL;",
// 			expectedTokens: []Token{
// 				{ Type: SELECT, Text: "SELECT" },
// 				{ Type: FIELD, Text: "name" },
// 				{ Type: FROM, Text: "FROM" },
// 				{ Type: FIELD, Text: "users" },
// 				{ Type: WHERE, Text: "WHERE" },
// 				{ Type: FIELD, Text: "email" },
// 				{ Type: IS, Text: "IS" },
// 				{ Type: NOT, Text: "NOT" },
// 				{ Type: NULL, Text: "NULL" },
// 				{ Type: SEMICOL, Text: ";" },
// 				{ Type: EOF, Text: "" },
// 			},
// 		},
// 		{
// 			name: "l_than",
// 			query: "SELECT name FROM users WHERE id < 2;",
// 			expectedTokens: []Token{
// 				{ Type: SELECT, Text: "SELECT" },
// 				{ Type: FIELD, Text: "name" },
// 				{ Type: FROM, Text: "FROM" },
// 				{ Type: FIELD, Text: "users" },
// 				{ Type: WHERE, Text: "WHERE" },
// 				{ Type: FIELD, Text: "id" },
// 				{ Type: L_THAN, Text: "<" },
// 				{ Type: NUMBER, Text: "2" },
// 				{ Type: SEMICOL, Text: ";" },
// 				{ Type: EOF, Text: "" },
// 			},
// 		},
// 		{
// 			name: "g_than",
// 			query: "SELECT name FROM users WHERE id > 0;",
// 			expectedTokens: []Token{
// 				{ Type: SELECT, Text: "SELECT" },
// 				{ Type: FIELD, Text: "name" },
// 				{ Type: FROM, Text: "FROM" },
// 				{ Type: FIELD, Text: "users" },
// 				{ Type: WHERE, Text: "WHERE" },
// 				{ Type: FIELD, Text: "id" },
// 				{ Type: G_THAN, Text: ">" },
// 				{ Type: NUMBER, Text: "0" },
// 				{ Type: SEMICOL, Text: ";" },
// 				{ Type: EOF, Text: "" },
// 			},
// 		},
// 		{
// 			name: "lt_equal",
// 			query: "SELECT name FROM users WHERE id <= 1;",
// 			expectedTokens: []Token{
// 				{ Type: SELECT, Text: "SELECT" },
// 				{ Type: FIELD, Text: "name" },
// 				{ Type: FROM, Text: "FROM" },
// 				{ Type: FIELD, Text: "users" },
// 				{ Type: WHERE, Text: "WHERE" },
// 				{ Type: FIELD, Text: "id" },
// 				{ Type: LT_EQUAL, Text: "<=" },
// 				{ Type: NUMBER, Text: "1" },
// 				{ Type: SEMICOL, Text: ";" },
// 				{ Type: EOF, Text: "" },
// 			},
// 		},
// 		{
// 			name: "gt_equal",
// 			query: "SELECT name FROM users WHERE id >= 1;",
// 			expectedTokens: []Token{
// 				{ Type: SELECT, Text: "SELECT" },
// 				{ Type: FIELD, Text: "name" },
// 				{ Type: FROM, Text: "FROM" },
// 				{ Type: FIELD, Text: "users" },
// 				{ Type: WHERE, Text: "WHERE" },
// 				{ Type: FIELD, Text: "id" },
// 				{ Type: GT_EQUAL, Text: ">=" },
// 				{ Type: NUMBER, Text: "1" },
// 				{ Type: SEMICOL, Text: ";" },
// 				{ Type: EOF, Text: "" },
// 			},
// 		},
// 		{
// 			name: "n_equal",
// 			query: "SELECT name FROM users WHERE id != 1;",
// 			expectedTokens: []Token{
// 				{ Type: SELECT, Text: "SELECT" },
// 				{ Type: FIELD, Text: "name" },
// 				{ Type: FROM, Text: "FROM" },
// 				{ Type: FIELD, Text: "users" },
// 				{ Type: WHERE, Text: "WHERE" },
// 				{ Type: FIELD, Text: "id" },
// 				{ Type: N_EQUAL, Text: "!=" },
// 				{ Type: NUMBER, Text: "1" },
// 				{ Type: SEMICOL, Text: ";" },
// 				{ Type: EOF, Text: "" },
// 			},
// 		},
// 		{
// 			name: "and_or",
// 			query: "SELECT name FROM users WHERE id = 1 AND (name = 'kevin' OR email = 'test@gmail.dk');",
// 			expectedTokens: []Token{
// 				{ Type: SELECT, Text: "SELECT" },
// 				{ Type: FIELD, Text: "name" },
// 				{ Type: FROM, Text: "FROM" },
// 				{ Type: FIELD, Text: "users" },
// 				{ Type: WHERE, Text: "WHERE" },
// 				{ Type: FIELD, Text: "id" },
// 				{ Type: EQUAL, Text: "=" },
// 				{ Type: NUMBER, Text: "1" },
// 				{ Type: AND, Text: "AND" },
// 				{ Type: L_PARAN, Text: "(" },
// 				{ Type: FIELD, Text: "name" },
// 				{ Type: EQUAL, Text: "=" },
// 				{ Type: STRING, Text: "kevin" },
// 				{ Type: OR, Text: "OR" },
// 				{ Type: FIELD, Text: "email" },
// 				{ Type: EQUAL, Text: "=" },
// 				{ Type: STRING, Text: "test@gmail.dk" },
// 				{ Type: R_PARAN, Text: ")" },
// 				{ Type: SEMICOL, Text: ";" },
// 				{ Type: EOF, Text: "" },
// 			},
// 		},
// 		{
// 			name: "as",
// 			query: "SELECT name AS result FROM users WHERE id = 1;",
// 			expectedTokens: []Token{
// 				{ Type: SELECT, Text: "SELECT" },
// 				{ Type: FIELD, Text: "name" },
// 				{ Type: AS, Text: "AS" },
// 				{ Type: FIELD, Text: "result" },
// 				{ Type: FROM, Text: "FROM" },
// 				{ Type: FIELD, Text: "users" },
// 				{ Type: WHERE, Text: "WHERE" },
// 				{ Type: FIELD, Text: "id" },
// 				{ Type: EQUAL, Text: "=" },
// 				{ Type: NUMBER, Text: "1" },
// 				{ Type: SEMICOL, Text: ";" },
// 				{ Type: EOF, Text: "" },
// 			},
// 		},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			parser := NewParser()
//
// 			if err := parser.ParseQuery(test.query); err != nil {
// 				t.Fatal(err)
// 			}
//
// 			if len(test.expectedTokens) != len(parser.tokens) {
// 				t.Fatalf("unexpected len, %v got %v", len(test.expectedTokens), len(parser.tokens))
// 			}
//
// 			for i := range parser.tokens {
// 				if test.expectedTokens[i].Type != parser.tokens[i].Type {
// 					t.Fatalf("unexpected token type at %v, %v (%s) got %v (%s)", i, test.expectedTokens[i].Type, test.expectedTokens[i].Text, parser.tokens[i].Type, parser.tokens[i].Text)
// 				}
//
// 				if test.expectedTokens[i].Text != parser.tokens[i].Text {
// 					t.Fatalf("unexpected token text at %v, %s got %s", i, test.expectedTokens[i].Text, parser.tokens[i].Text)
// 				}
// 			}
// 		})
// 	}
// }
//
// func TestParseQueryTokenStats(t *testing.T) {
// 	tests := []struct{
// 		name						string
// 		query						string
// 		expectedTokens	[]Token
// 	} {
// 		{
// 			name: "insert",
// 			query: "INSERT INTO users (name, email) VALUES ('kevin', 'test@gmail.dk');",
// 			expectedTokens: []Token{
// 				{ Type: INSERT, Text: "INSERT", LineNumber: 1, Start: 0, End: 6 },
// 				{ Type: INTO, Text: "INTO", LineNumber: 1, Start: 7, End: 11 },
// 				{ Type: FIELD, Text: "users", LineNumber: 1, Start: 12, End: 17 },
// 				{ Type: L_PARAN, Text: "(", LineNumber: 1, Start: 18, End: 19 },
// 				{ Type: FIELD, Text: "name", LineNumber: 1, Start: 19, End: 23 },
// 				{ Type: FIELD, Text: "email", LineNumber: 1, Start: 25, End: 30 },
// 				{ Type: R_PARAN, Text: ")", LineNumber: 1, Start: 30, End: 31 },
// 				{ Type: VALUES, Text: "VALUES", LineNumber: 1, Start: 32, End: 38 },
// 				{ Type: L_PARAN, Text: "(", LineNumber: 1, Start: 39, End: 40 },
// 				{ Type: STRING, Text: "kevin", LineNumber: 1, Start: 41, End: 46 },
// 				{ Type: STRING, Text: "test@gmail.dk", LineNumber: 1, Start: 50, End: 63 },
// 				{ Type: R_PARAN, Text: ")", LineNumber: 1, Start: 64, End: 65 },
// 				{ Type: SEMICOL, Text: ";", LineNumber: 1, Start: 65, End: 66 },
// 				{ Type: EOF, Text: "", LineNumber: 1, Start: 66, End: 66 },
// 			},
// 		},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			parser := NewParser()
//
// 			if err := parser.ParseQuery(test.query); err != nil {
// 				t.Fatal(err)
// 			}
//
// 			if len(test.expectedTokens) != len(parser.tokens) {
// 				t.Fatalf("unexpected len, %v got %v", len(test.expectedTokens), len(parser.tokens))
// 			}
//
// 			for i := range parser.tokens {
// 				if test.expectedTokens[i].Type != parser.tokens[i].Type {
// 					t.Fatalf("unexpected token type at %v, %v (%s) got %v (%s)", i, test.expectedTokens[i].Type, test.expectedTokens[i].Text, parser.tokens[i].Type, parser.tokens[i].Text)
// 				}
//
// 				if test.expectedTokens[i].Text != parser.tokens[i].Text {
// 					t.Fatalf("unexpected token text at %v, %s got %s", i, test.expectedTokens[i].Text, parser.tokens[i].Text)
// 				}
//
// 				if test.expectedTokens[i].LineNumber != parser.tokens[i].LineNumber {
// 					t.Fatalf("unexpected line number at %v, %d got %v", i, test.expectedTokens[i].LineNumber, parser.tokens[i].LineNumber)
// 				}
//
// 				if test.expectedTokens[i].Start != parser.tokens[i].Start {
// 					t.Fatalf("unexpected start position at %v, %d got %v", i, test.expectedTokens[i].Start, parser.tokens[i].Start)
// 				}
//
// 				if test.expectedTokens[i].End != parser.tokens[i].End {
// 					t.Fatalf("unexpected end position at %v, %d got %v", i, test.expectedTokens[i].End, parser.tokens[i].End)
// 				}
// 			}
// 		})
// 	}
// }

func TestParseQueryTable(t *testing.T) {
	tests := []struct{
		name							string
		query							string
		expectedResponse	string
		expectedError			error
	} {
		{
			name: "insert",
			query: "INSERT INTO users (name, age, email) VALUES ('kevin', 25, 'test@gmail.dk');",
			expectedResponse: "users",
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := NewParser()

			if err := parser.ParseQuery(test.query); err != test.expectedError {
				t.Fatalf("unexpected error while parsing, %v got %v", test.expectedError, err)
			}

			if test.expectedError != nil {
				return
			}

			if treeLen := len(parser.trees); treeLen != 1 {
				t.Fatalf("unexpected tree len, 1 got %v", treeLen)
			}

			response := parser.trees[0].Table()
			if test.expectedResponse != response.String() {
				t.Fatalf("unexpected response, '%s' got '%s'", test.expectedResponse, response.String())
			}
		})
	}
}

func TestParseQueryColumns(t *testing.T) {
	tests := []struct{
		name							string
		query							string
		expectedResponse	[]string
		expectedError			error
	} {
		{
			name: "insert",
			query: "INSERT INTO users (name, age, email) VALUES ('kevin', 25, 'test@gmail.dk');",
			expectedResponse: []string{
				"name",
				"age",
				"email",
			},
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := NewParser()

			if err := parser.ParseQuery(test.query); err != test.expectedError {
				t.Fatalf("unexpected error while parsing, %v got %v", test.expectedError, err)
			}

			if test.expectedError != nil {
				return
			}

			if treeLen := len(parser.trees); treeLen != 1 {
				t.Fatalf("unexpected tree len, 1 got %v", treeLen)
			}

			response := parser.trees[0].Columns()
			for i := range response {
				if test.expectedResponse[i] != response[i].String() {
					t.Fatalf("unexpected response, '%s' got '%s'", test.expectedResponse[i], response[i].String())
				}
			}
		})
	}
}

func TestParseQueryToString(t *testing.T) {
	tests := []struct{
		name							string
		query							string
		expectedError			error
	} {
		{
			name: "insert",
			query: "INSERT INTO users (name, age, email) VALUES ('kevin', 25, 'test@gmail.dk');",
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := NewParser()

			if err := parser.ParseQuery(test.query); err != test.expectedError {
				t.Fatalf("unexpected error while parsing, %v got %v", test.expectedError, err)
			}

			if test.expectedError != nil {
				return
			}

			if treeLen := len(parser.trees); treeLen != 1 {
				t.Fatalf("unexpected tree len, 1 got %v", treeLen)
			}

			response := parser.trees[0].String()
			if test.query != response {
				t.Fatalf("unexpected response, '%s' got '%s'", test.query, response)
			}
		})
	}
}
