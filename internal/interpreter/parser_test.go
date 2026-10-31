package interpreter

import (
	"errors"
	"testing"
)

func TestParseSql(t *testing.T) {
	tests := []struct{
		name							string
		sql								string
		expectedTokens		[]Token
		expectedErrorCode ErrorCode
	} {
		{
			name: "insert",
			sql: "INSERT INTO users (name, email) VALUES ('kevin', 'test@gmail.dk');",
			expectedTokens: []Token{
				{ Type: INSERT, Text: "INSERT" },
				{ Type: INTO, Text: "INTO" },
				{ Type: FIELD, Text: "users" },
				{ Type: L_PARAN, Text: "(" },
				{ Type: FIELD, Text: "name" },
				{ Type: COMMA, Text: "," },
				{ Type: FIELD, Text: "email" },
				{ Type: R_PARAN, Text: ")" },
				{ Type: VALUES, Text: "VALUES" },
				{ Type: L_PARAN, Text: "(" },
				{ Type: STRING, Text: "kevin" },
				{ Type: COMMA, Text: "," },
				{ Type: STRING, Text: "test@gmail.dk" },
				{ Type: R_PARAN, Text: ")" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedErrorCode: -1,
		},
		{
			name: "select",
			sql: "SELECT name, age, email FROM users WHERE id = 1 AND name LIKE '%kevin%';",
			expectedTokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: FIELD, Text: "name" },
				{ Type: COMMA, Text: "," },
				{ Type: FIELD, Text: "age" },
				{ Type: COMMA, Text: "," },
				{ Type: FIELD, Text: "email" },
				{ Type: FROM, Text: "FROM" },
				{ Type: FIELD, Text: "users" },
				{ Type: WHERE, Text: "WHERE" },
				{ Type: FIELD, Text: "id" },
				{ Type: EQUAL, Text: "=" },
				{ Type: NUMBER, Text: "1" },
				{ Type: AND, Text: "AND" },
				{ Type: FIELD, Text: "name" },
				{ Type: LIKE, Text: "LIKE" },
				{ Type: STRING, Text: "%kevin%" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedErrorCode: -1,
		},
		{
			name: "update",
			sql: "UPDATE users SET name = 'john' WHERE id = 1;",
			expectedTokens: []Token{
				{ Type: UPDATE, Text: "UPDATE" },
				{ Type: FIELD, Text: "users" },
				{ Type: SET, Text: "SET" },
				{ Type: FIELD, Text: "name" },
				{ Type: EQUAL, Text: "=" },
				{ Type: STRING, Text: "john" },
				{ Type: WHERE, Text: "WHERE" },
				{ Type: FIELD, Text: "id" },
				{ Type: EQUAL, Text: "=" },
				{ Type: NUMBER, Text: "1" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedErrorCode: -1,
		},
		{
			name: "delete",
			sql: "DELETE FROM users WHERE id = 1 AND name LIKE '%kevin%';",
			expectedTokens: []Token{
				{ Type: DELETE, Text: "DELETE" },
				{ Type: FROM, Text: "FROM" },
				{ Type: FIELD, Text: "users" },
				{ Type: WHERE, Text: "WHERE" },
				{ Type: FIELD, Text: "id" },
				{ Type: EQUAL, Text: "=" },
				{ Type: NUMBER, Text: "1" },
				{ Type: AND, Text: "AND" },
				{ Type: FIELD, Text: "name" },
				{ Type: LIKE, Text: "LIKE" },
				{ Type: STRING, Text: "%kevin%" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedErrorCode: -1,
		},
		{
			name: "drop",
			sql: "DROP TABLE IF EXISTS users;",
			expectedTokens: []Token{
				{ Type: DROP, Text: "DROP" },
				{ Type: TABLE, Text: "TABLE" },
				{ Type: IF, Text: "IF" },
				{ Type: EXISTS, Text: "EXISTS" },
				{ Type: FIELD, Text: "users" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedErrorCode: -1,
		},
		{
			name: "alter_add",
			sql: "ALTER TABLE users ADD COLUMN country varchar(128);",
			expectedTokens: []Token{
				{ Type: ALTER, Text: "ALTER" },
				{ Type: TABLE, Text: "TABLE" },
				{ Type: FIELD, Text: "users" },
				{ Type: ADD, Text: "ADD" },
				{ Type: COLUMN, Text: "COLUMN" },
				{ Type: FIELD, Text: "country" },
				{ Type: FIELD, Text: "varchar" },
				{ Type: L_PARAN, Text: "(" },
				{ Type: NUMBER, Text: "128" },
				{ Type: R_PARAN, Text: ")" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedErrorCode: -1,
		},
		{
			name: "alter_drop",
			sql: "ALTER TABLE users DROP COLUMN country;",
			expectedTokens: []Token{
				{ Type: ALTER, Text: "ALTER" },
				{ Type: TABLE, Text: "TABLE" },
				{ Type: FIELD, Text: "users" },
				{ Type: DROP, Text: "DROP" },
				{ Type: COLUMN, Text: "COLUMN" },
				{ Type: FIELD, Text: "country" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedErrorCode: -1,
		},
		{
			name: "alter_alter",
			sql: "ALTER TABLE users ALTER COLUMN country varchar(256);",
			expectedTokens: []Token{
				{ Type: ALTER, Text: "ALTER" },
				{ Type: TABLE, Text: "TABLE" },
				{ Type: FIELD, Text: "users" },
				{ Type: ALTER, Text: "ALTER" },
				{ Type: COLUMN, Text: "COLUMN" },
				{ Type: FIELD, Text: "country" },
				{ Type: FIELD, Text: "varchar" },
				{ Type: L_PARAN, Text: "(" },
				{ Type: NUMBER, Text: "256" },
				{ Type: R_PARAN, Text: ")" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedErrorCode: -1,
		},
		{
			name: "create",
			sql: "CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT NOT NULL);",
			expectedTokens: []Token{
				{ Type: CREATE, Text: "CREATE" },
				{ Type: TABLE, Text: "TABLE" },
				{ Type: FIELD, Text: "users" },
				{ Type: L_PARAN, Text: "(" },
				{ Type: FIELD, Text: "id" },
				{ Type: FIELD, Text: "INTEGER" },
				{ Type: PRIMARY, Text: "PRIMARY" },
				{ Type: KEY, Text: "KEY" },
				{ Type: COMMA, Text: "," },
				{ Type: FIELD, Text: "name" },
				{ Type: FIELD, Text: "TEXT" },
				{ Type: NOT, Text: "NOT" },
				{ Type: NULL, Text: "NULL" },
				{ Type: R_PARAN, Text: ")" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedErrorCode: -1,
		},
		{
			name: "func",
			sql: "SELECT now();",
			expectedTokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: FIELD, Text: "now" },
				{ Type: L_PARAN, Text: "(" },
				{ Type: R_PARAN, Text: ")" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedErrorCode: -1,
		},
		{
			name: "is_not_null",
			sql: "SELECT name FROM users WHERE email IS NOT NULL;",
			expectedTokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: FIELD, Text: "name" },
				{ Type: FROM, Text: "FROM" },
				{ Type: FIELD, Text: "users" },
				{ Type: WHERE, Text: "WHERE" },
				{ Type: FIELD, Text: "email" },
				{ Type: IS, Text: "IS" },
				{ Type: NOT, Text: "NOT" },
				{ Type: NULL, Text: "NULL" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedErrorCode: -1,
		},
		{
			name: "l_than",
			sql: "SELECT name FROM users WHERE id < 2;",
			expectedTokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: FIELD, Text: "name" },
				{ Type: FROM, Text: "FROM" },
				{ Type: FIELD, Text: "users" },
				{ Type: WHERE, Text: "WHERE" },
				{ Type: FIELD, Text: "id" },
				{ Type: L_THAN, Text: "<" },
				{ Type: NUMBER, Text: "2" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedErrorCode: -1,
		},
		{
			name: "g_than",
			sql: "SELECT name FROM users WHERE id > 0;",
			expectedTokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: FIELD, Text: "name" },
				{ Type: FROM, Text: "FROM" },
				{ Type: FIELD, Text: "users" },
				{ Type: WHERE, Text: "WHERE" },
				{ Type: FIELD, Text: "id" },
				{ Type: G_THAN, Text: ">" },
				{ Type: NUMBER, Text: "0" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedErrorCode: -1,
		},
		{
			name: "lt_equal",
			sql: "SELECT name FROM users WHERE id <= 1;",
			expectedTokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: FIELD, Text: "name" },
				{ Type: FROM, Text: "FROM" },
				{ Type: FIELD, Text: "users" },
				{ Type: WHERE, Text: "WHERE" },
				{ Type: FIELD, Text: "id" },
				{ Type: LT_EQUAL, Text: "<=" },
				{ Type: NUMBER, Text: "1" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedErrorCode: -1,
		},
		{
			name: "gt_equal",
			sql: "SELECT name FROM users WHERE id >= 1;",
			expectedTokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: FIELD, Text: "name" },
				{ Type: FROM, Text: "FROM" },
				{ Type: FIELD, Text: "users" },
				{ Type: WHERE, Text: "WHERE" },
				{ Type: FIELD, Text: "id" },
				{ Type: GT_EQUAL, Text: ">=" },
				{ Type: NUMBER, Text: "1" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedErrorCode: -1,
		},
		{
			name: "n_equal",
			sql: "SELECT name FROM users WHERE id != 1;",
			expectedTokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: FIELD, Text: "name" },
				{ Type: FROM, Text: "FROM" },
				{ Type: FIELD, Text: "users" },
				{ Type: WHERE, Text: "WHERE" },
				{ Type: FIELD, Text: "id" },
				{ Type: N_EQUAL, Text: "!=" },
				{ Type: NUMBER, Text: "1" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedErrorCode: -1,
		},
		{
			name: "and_or",
			sql: "SELECT name FROM users WHERE id = 1 AND (name = 'kevin' OR email = 'test@gmail.dk');",
			expectedTokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: FIELD, Text: "name" },
				{ Type: FROM, Text: "FROM" },
				{ Type: FIELD, Text: "users" },
				{ Type: WHERE, Text: "WHERE" },
				{ Type: FIELD, Text: "id" },
				{ Type: EQUAL, Text: "=" },
				{ Type: NUMBER, Text: "1" },
				{ Type: AND, Text: "AND" },
				{ Type: L_PARAN, Text: "(" },
				{ Type: FIELD, Text: "name" },
				{ Type: EQUAL, Text: "=" },
				{ Type: STRING, Text: "kevin" },
				{ Type: OR, Text: "OR" },
				{ Type: FIELD, Text: "email" },
				{ Type: EQUAL, Text: "=" },
				{ Type: STRING, Text: "test@gmail.dk" },
				{ Type: R_PARAN, Text: ")" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedErrorCode: -1,
		},
		{
			name: "wildcard",
			sql: "SELECT name, * FROM users;",
			expectedTokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: FIELD, Text: "name" },
				{ Type: COMMA, Text: "," },
				{ Type: WILDCARD, Text: "*" },
				{ Type: FROM, Text: "FROM" },
				{ Type: FIELD, Text: "users" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedErrorCode: -1,
		},
		{
			name: "as",
			sql: "SELECT name AS result FROM users WHERE id = 1;",
			expectedTokens: []Token{
				{ Type: SELECT, Text: "SELECT" },
				{ Type: FIELD, Text: "name" },
				{ Type: AS, Text: "AS" },
				{ Type: FIELD, Text: "result" },
				{ Type: FROM, Text: "FROM" },
				{ Type: FIELD, Text: "users" },
				{ Type: WHERE, Text: "WHERE" },
				{ Type: FIELD, Text: "id" },
				{ Type: EQUAL, Text: "=" },
				{ Type: NUMBER, Text: "1" },
				{ Type: SEMICOL, Text: ";" },
				{ Type: EOF, Text: "" },
			},
			expectedErrorCode: -1,
		},
		{
			name: "pragma",
			sql: "PRAGMA schema.table_info(users);",
			expectedTokens: 	nil,
			expectedErrorCode: UNSUPPORTED_QUERY,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := NewParser()

			_, err := parser.ParseSql(test.sql); 

			typedErr, ok := errors.AsType[*Error](err)
			if err != nil || test.expectedErrorCode != -1 {
				if !ok {
					t.Fatalf("unexpected error type, %v", err)
				}

				if test.expectedErrorCode != typedErr.code {
					t.Fatalf("unexpected error, %v got %v (%s)", test.expectedErrorCode, typedErr.code, typedErr.msg)
				}
				return
			}

			if len(test.expectedTokens) != len(parser.tokens) {
				t.Fatalf("unexpected len, %v got %v", len(test.expectedTokens), len(parser.tokens))
			}

			for i := range parser.tokens {
				if test.expectedTokens[i].Type != parser.tokens[i].Type {
					t.Fatalf("unexpected token type at %v, %v (%s) got %v (%s)", i, test.expectedTokens[i].Type, test.expectedTokens[i].Text, parser.tokens[i].Type, parser.tokens[i].Text)
				}

				if test.expectedTokens[i].Text != parser.tokens[i].Text {
					t.Fatalf("unexpected token text at %v, %s got %s", i, test.expectedTokens[i].Text, parser.tokens[i].Text)
				}
			}
		})
	}
}

func TestParseQueryTokenStats(t *testing.T) {
	tests := []struct{
		name						string
		sql							string
		expectedTokens	[]Token
	} {
		{
			name: "insert",
			sql: "INSERT INTO users (name, email) VALUES ('kevin', 'test@gmail.dk');",
			expectedTokens: []Token{
				{ Type: INSERT, Text: "INSERT", LineNumber: 1, Start: 0, End: 6 },
				{ Type: INTO, Text: "INTO", LineNumber: 1, Start: 7, End: 11 },
				{ Type: FIELD, Text: "users", LineNumber: 1, Start: 12, End: 17 },
				{ Type: L_PARAN, Text: "(", LineNumber: 1, Start: 18, End: 19 },
				{ Type: FIELD, Text: "name", LineNumber: 1, Start: 19, End: 23 },
				{ Type: COMMA, Text: ",", LineNumber: 1, Start: 23, End: 24 },
				{ Type: FIELD, Text: "email", LineNumber: 1, Start: 25, End: 30 },
				{ Type: R_PARAN, Text: ")", LineNumber: 1, Start: 30, End: 31 },
				{ Type: VALUES, Text: "VALUES", LineNumber: 1, Start: 32, End: 38 },
				{ Type: L_PARAN, Text: "(", LineNumber: 1, Start: 39, End: 40 },
				{ Type: STRING, Text: "kevin", LineNumber: 1, Start: 41, End: 46 },
				{ Type: COMMA, Text: ",", LineNumber: 1, Start: 47, End: 48 },
				{ Type: STRING, Text: "test@gmail.dk", LineNumber: 1, Start: 50, End: 63 },
				{ Type: R_PARAN, Text: ")", LineNumber: 1, Start: 64, End: 65 },
				{ Type: SEMICOL, Text: ";", LineNumber: 1, Start: 65, End: 66 },
				{ Type: EOF, Text: "", LineNumber: 1, Start: 66, End: 66 },
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := NewParser()

			if _, err := parser.ParseSql(test.sql); err != nil {
				t.Fatal(err)
			}

			if len(test.expectedTokens) != len(parser.tokens) {
				t.Fatalf("unexpected len, %v got %v", len(test.expectedTokens), len(parser.tokens))
			}

			for i := range parser.tokens {
				if test.expectedTokens[i].Type != parser.tokens[i].Type {
					t.Fatalf("unexpected token type at %v, %v (%s) got %v (%s)", i, test.expectedTokens[i].Type, test.expectedTokens[i].Text, parser.tokens[i].Type, parser.tokens[i].Text)
				}

				if test.expectedTokens[i].Text != parser.tokens[i].Text {
					t.Fatalf("unexpected token text at %v, %s got %s", i, test.expectedTokens[i].Text, parser.tokens[i].Text)
				}

				if test.expectedTokens[i].LineNumber != parser.tokens[i].LineNumber {
					t.Fatalf("unexpected line number at %v, %d got %v", i, test.expectedTokens[i].LineNumber, parser.tokens[i].LineNumber)
				}

				if test.expectedTokens[i].Start != parser.tokens[i].Start {
					t.Fatalf("unexpected start position at %v, %d got %v", i, test.expectedTokens[i].Start, parser.tokens[i].Start)
				}

				if test.expectedTokens[i].End != parser.tokens[i].End {
					t.Fatalf("unexpected end position at %v, %d got %v", i, test.expectedTokens[i].End, parser.tokens[i].End)
				}
			}
		})
	}
}

func TestParseQueryTable(t *testing.T) {
	tests := []struct{
		name							string
		sql								string
		expectedResponse	string
	} {
		{
			name: "insert",
			sql: "INSERT INTO users (name, age, email) VALUES ('kevin', 25, 'test@gmail.dk');",
			expectedResponse: "users",
		},
		{
			name: "select",
			sql: "SELECT name, age FROM users WHERE id = 1;",
			expectedResponse: "users",
		},
		{
			name: "update",
			sql: "UPDATE users SET name = 'kevin', age = 25 WHERE id = 1;",
			expectedResponse: "users",
		},
		{
			name: "delete",
			sql: "DELETE FROM users WHERE id = 1;",
			expectedResponse: "users",
		},
		{
			name: "create",
			sql: "CREATE TABLE users(id INTEGER PRIMARY KEY);",
			expectedResponse: "users",
		},
		{
			name: "drop",
			sql: "DROP TABLE users;",
			expectedResponse: "users",
		},
		{
			name: "alter_add",
			sql: "ALTER TABLE users ADD COLUMN age INTEGER NOT NULL;",
			expectedResponse: "users",
		},
		{
			name: "alter_alter",
			sql: "ALTER TABLE users ALTER COLUMN age INTEGER;",
			expectedResponse: "users",
		},
		{
			name: "alter_drop",
			sql: "ALTER TABLE users DROP COLUMN age INTEGER;",
			expectedResponse: "users",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := NewParser()

			trees, err := parser.ParseSql(test.sql)
			if err != nil {
				t.Fatalf("unexpected error while parsing, %v", err)
			}

			if treeLen := len(trees); treeLen != 1 {
				t.Fatalf("unexpected tree len, 1 got %v", treeLen)
			}

			response := trees[0].Table()
			if test.expectedResponse != response.String() {
				t.Fatalf("unexpected response, '%s' got '%s'", test.expectedResponse, response.String())
			}
		})
	}
}

func TestParseQueryColumns(t *testing.T) {
	tests := []struct{
		name							string
		sql								string
		expectedResponse	[]string
	} {
		{
			name: "insert",
			sql: "INSERT INTO users (name, age, email) VALUES ('kevin', 25, 'test@gmail.dk');",
			expectedResponse: []string{
				"name",
				"age",
				"email",
			},
		},
		{
			name: "select",
			sql: "SELECT name, age FROM users WHERE id = 1;",
			expectedResponse: []string{
				"name",
				"age",
			},
		},
		{
			name: "update",
			sql: "UPDATE users SET name = 'kevin', age = 25 WHERE id = 1;",
			expectedResponse: []string{},
		},
		{
			name: "delete",
			sql: "DELETE FROM users WHERE id = 1;",
			expectedResponse: []string{},
		},
		{
			name: "create",
			sql: "CREATE TABLE users(id INTEGER PRIMARY KEY);",
			expectedResponse: []string{"id"},
		},
		{
			name: "drop",
			sql: "DROP TABLE users;",
			expectedResponse: []string{},
		},
		{
			name: "alter_add",
			sql: "ALTER TABLE users ADD COLUMN age INTEGER NOT NULL;",
			expectedResponse: []string{"age"},
		},
		{
			name: "alter_alter",
			sql: "ALTER TABLE users ALTER COLUMN age INTEGER;",
			expectedResponse: []string{"age"},
		},
		{
			name: "alter_drop",
			sql: "ALTER TABLE users DROP COLUMN age INTEGER;",
			expectedResponse: []string{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := NewParser()

			trees, err := parser.ParseSql(test.sql)
			if err != nil {
				t.Fatalf("unexpected error while parsing, %v", err)
			}

			if treeLen := len(trees); treeLen != 1 {
				t.Fatalf("unexpected tree len, 1 got %v (%v)", treeLen, trees)
			}

			response := trees[0].Columns()
			for i := range response {
				if test.expectedResponse[i] != response[i] {
					t.Fatalf("unexpected response, '%s' got '%s'", test.expectedResponse[i], response[i])
				}
			}
		})
	}
}

func TestParseQueryToString(t *testing.T) {
	tests := []struct{
		name							string
		sql								string
	} {
		{
			name: "insert",
			sql: "INSERT INTO users (name, age, email) VALUES ('kevin', 25, 'test@gmail.dk');",
		},
		{
			name: "select",
			sql: "SELECT name, age FROM users WHERE name LIKE '%kev%' AND age < 30;",
		},
		{
			name: "update",
			sql: "UPDATE users SET age = 26 WHERE id = 1;",
		},
		{
			name: "delete",
			sql: "DELETE FROM users WHERE id = 1;",
		},
		{
			name: "create",
			sql: "CREATE TABLE users(id INTEGER PRIMARY KEY);",
		},
		{
			name: "drop",
			sql: "DROP TABLE users;",
		},
		{
			name: "alter_add",
			sql: "ALTER TABLE users ADD COLUMN age INTEGER NOT NULL;",
		},
		{
			name: "alter_alter",
			sql: "ALTER TABLE users ALTER COLUMN age INTEGER;",
		},
		{
			name: "alter_drop",
			sql: "ALTER TABLE users DROP COLUMN age INTEGER;",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := NewParser()

			trees, err := parser.ParseSql(test.sql)
			if err != nil {
				t.Fatalf("unexpected error while parsing, %v", err)
			}

			if treeLen := len(trees); treeLen != 1 {
				t.Fatalf("unexpected tree len, 1 got %v", treeLen)
			}

			response := trees[0].String()
			if test.sql != response {
				t.Fatalf("unexpected response, '%s' got '%s'", test.sql, response)
			}
		})
	}
}

func TestExpressionParse(t *testing.T) {
	tests := []struct{
		name							string
		tokens						[]Token
		expectedResponse	string
	} {
		{
			name: "and",
			tokens: []Token{
				{ Type: FIELD, Text: "name" },
				{ Type: EQUAL, Text: "=" },
				{ Type: STRING, Text: "Kevin" },
				{ Type: AND, Text: "AND" },
				{ Type: FIELD, Text: "age" },
				{ Type: L_THAN, Text: "<" },
				{ Type: NUMBER, Text: "30" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "name = 'Kevin' AND age < 30",
		},
		{
			name: "and_or",
			tokens: []Token{
				{ Type: FIELD, Text: "name" },
				{ Type: EQUAL, Text: "=" },
				{ Type: STRING, Text: "Kevin" },
				{ Type: AND, Text: "AND" },
				{ Type: FIELD, Text: "age" },
				{ Type: L_THAN, Text: "<" },
				{ Type: NUMBER, Text: "30" },
				{ Type: OR, Text: "OR" },
				{ Type: FIELD, Text: "id" },
				{ Type: L_THAN, Text: "=" },
				{ Type: NUMBER, Text: "1" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "name = 'Kevin' AND age < 30 OR id = 1",
		},
		{
			name: "like",
			tokens: []Token{
				{ Type: FIELD, Text: "name" },
				{ Type: LIKE, Text: "LIKE" },
				{ Type: STRING, Text: "%Kev%" },
				{ Type: AND, Text: "AND" },
				{ Type: FIELD, Text: "age" },
				{ Type: L_THAN, Text: "<" },
				{ Type: NUMBER, Text: "30" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "name LIKE '%Kev%' AND age < 30",
		},
		{
			name: "is",
			tokens: []Token{
				{ Type: FIELD, Text: "age" },
				{ Type: IS, Text: "IS" },
				{ Type: NOT, Text: "NOT" },
				{ Type: NULL, Text: "NULL" },
				{ Type: EOF, Text: "" },
			},
			expectedResponse: "age IS NOT NULL",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parser := &Parser{
				tokens: test.tokens,
				r: 0,
			}

			node, err := parser.expressionParse(0)
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

