package interpreter

import (
	"errors"
	"testing"
)

func TestLexer(t *testing.T) {
	tests := []struct{
		name							string
		query							string
		expectedResponse  []Token
	} {
		{
			name: "multi_line_with_comment",
			query: "-- Write something here\nSELECT * FROM users WHERE id = 1;",
			expectedResponse: []Token{
				{ Type: SELECT, Text: "SELECT",	LineNumber: 2, Start: 0, End: 6 },
				{ Type: FIELD, Text: "*",				LineNumber: 2, Start: 7, End: 8 },
				{ Type: FROM, Text: "FROM",			LineNumber: 2, Start: 9, End: 13 },
				{ Type: FIELD, Text: "users",		LineNumber: 2, Start: 14, End: 19 },
				{ Type: WHERE, Text: "WHERE",		LineNumber: 2, Start: 20, End: 25 },
				{ Type: FIELD, Text: "id",			LineNumber: 2, Start: 26, End: 28 },
				{ Type: EQUAL, Text: "=",				LineNumber: 2, Start: 29, End: 30 },
				{ Type: NUMBER, Text: "1",			LineNumber: 2, Start: 31, End: 32 },
				{ Type: SEMICOL, Text: ";",			LineNumber: 2, Start: 32, End: 33 },
				{ Type: EOF, Text: "",					LineNumber: 2, Start: 33, End: 33 },
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var tokens []Token

			lexer := NewLexer()
			lexer.Write(test.query)

			for {
				token, err := lexer.Next()
				if err != nil {
					t.Fatal(err)
				}

				tokens = append(tokens, token)

				if token.Type == EOF {
					break
				}
			}

			for i := range tokens {
				if test.expectedResponse[i].Type != tokens[i].Type {
					t.Fatalf("unexpected token type, %s got %s (%s)", test.expectedResponse[i].Type, tokens[i].Type, tokens[i].Text)
				}

				if test.expectedResponse[i].Text != tokens[i].Text {
					t.Fatalf("unexpected token text, %s got %s", test.expectedResponse[i].Text, tokens[i].Text)
				}

				if test.expectedResponse[i].LineNumber != tokens[i].LineNumber {
					t.Fatalf("unexpected token line number, %v got %v (%s)", test.expectedResponse[i].LineNumber, tokens[i].LineNumber, tokens[i].Text)
				}

				if test.expectedResponse[i].Start != tokens[i].Start {
					t.Fatalf("unexpected token start, %v got %v (%s)", test.expectedResponse[i].Start, tokens[i].Start, tokens[i].Text)
				}

				if test.expectedResponse[i].End != tokens[i].End {
					t.Fatalf("unexpected token end, %v got %v (%s)", test.expectedResponse[i].End, tokens[i].End, tokens[i].Text)
				}
			}
		})
	}
}
  

func TestNext(t *testing.T) {
	tests := []struct{
		name													string
		query													string
		expectedResponse  						Token
		expectedErrorCode 						ErrorCode
		expectedErrorLineNumber				int
		expectedErrorCharStartNumber	int
		expectedErrorCharEndNumber		int
	} {
		{
			name: "keyword",
			query: "SeLeCt",
			expectedResponse: Token{ Type: SELECT, Text: "SeLeCt" },
			expectedErrorCode: -1,
			expectedErrorLineNumber: -1,
			expectedErrorCharStartNumber: -1,
			expectedErrorCharEndNumber: -1,
		},
		{
			name: "string",
			query: "'test'",
			expectedResponse: Token{ Type: STRING, Text: "test" },
			expectedErrorCode: -1,
			expectedErrorLineNumber: -1,
			expectedErrorCharStartNumber: -1,
			expectedErrorCharEndNumber: -1,
		},
		{
			name: "number",
			query: "-12",
			expectedResponse: Token{ Type: NUMBER, Text: "-12" },
			expectedErrorCode: -1,
			expectedErrorLineNumber: -1,
			expectedErrorCharStartNumber: -1,
			expectedErrorCharEndNumber: -1,
		},
		{
			name: "field",
			query: "users",
			expectedResponse: Token{ Type: FIELD, Text: "users" },
			expectedErrorCode: -1,
			expectedErrorLineNumber: -1,
			expectedErrorCharStartNumber: -1,
			expectedErrorCharEndNumber: -1,
		},
		{
			name: "invalid_field",
			query: "user--test",
			expectedResponse: Token{},
			expectedErrorCode: INVALID_QUERY_LEX_KEYWORD,
			expectedErrorLineNumber: 1,
			expectedErrorCharStartNumber: 0,
			expectedErrorCharEndNumber: 5,
		},
		{
			name: "left_paran",
			query: "(",
			expectedResponse: Token{ Type: L_PARAN, Text: "(" },
			expectedErrorCode: -1,
			expectedErrorLineNumber: -1,
			expectedErrorCharStartNumber: -1,
			expectedErrorCharEndNumber: -1,
		},
		{
			name: "less_than",
			query: "< =",
			expectedResponse: Token{ Type: L_THAN, Text: "<" },
			expectedErrorCode: -1,
			expectedErrorLineNumber: -1,
			expectedErrorCharStartNumber: -1,
			expectedErrorCharEndNumber: -1,
		},
		{
			name: "less_than_or_equal",
			query: "<=",
			expectedResponse: Token{ Type: LT_EQUAL, Text: "<=" },
			expectedErrorCode: -1,
			expectedErrorLineNumber: -1,
			expectedErrorCharStartNumber: -1,
			expectedErrorCharEndNumber: -1,
		},
		{
			name: "line_break",
			query: "\n",
			expectedResponse: Token{ Type: EOF, Text: "" },
			expectedErrorCode: -1,
			expectedErrorLineNumber: -1,
			expectedErrorCharStartNumber: -1,
			expectedErrorCharEndNumber: -1,
		},
		{
			name: "comment",
			query: "-- test",
			expectedResponse: Token{ Type: EOF, Text: "" },
			expectedErrorCode: -1,
			expectedErrorLineNumber: -1,
			expectedErrorCharStartNumber: -1,
			expectedErrorCharEndNumber: -1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lexer := NewLexer()
			lexer.Write(test.query)

			response, err := lexer.Next()

			if typedErr, ok := errors.AsType[*Error](err); ok {
				if test.expectedErrorCode != typedErr.code {
					t.Fatalf("unexpected error, %v got %v (%s)", test.expectedErrorCode, typedErr.code, typedErr.msg)
				}

				errLn, errStartCn, errEndCn := typedErr.Pos()
				if test.expectedErrorLineNumber != errLn {
					t.Fatalf("unexpected error line number, %v got %v", test.expectedErrorLineNumber, errLn)
				}

				if test.expectedErrorCharStartNumber != errStartCn {
					t.Fatalf("unexpected error char start number, %v got %v", test.expectedErrorCharStartNumber, errStartCn)
				}

				if test.expectedErrorCharEndNumber != errEndCn {
					t.Fatalf("unexpected error char end number, %v got %v", test.expectedErrorCharEndNumber, errEndCn)
				}
			}

			if err != nil && test.expectedErrorCode == -1 {
				t.Fatal(err)
			}

			if test.expectedErrorCode != -1 {
				return
			}
    	
			if test.expectedResponse.Type != response.Type {
				t.Fatalf("unexpected response type, %s got %s (%s)", test.expectedResponse.Type, response.Type, response.Text)
			}
		
			if test.expectedResponse.Text != response.Text {
				t.Fatalf("unexpected response text, %s got %s", test.expectedResponse.Text, response.Text)
			}
		})
	}
}
  

func TestNextNoQuery(t *testing.T) {
	tests := []struct{
		name					string
		expectedErrorCode			ErrorCode
	} {
		{
			name: "no_query",
			expectedErrorCode: NO_QUERY,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lexer := NewLexer()
			
			_, err := lexer.Next()
			lexErr, ok := errors.AsType[*Error](err)
			if !ok {
				t.Fatalf("unexpected error type, %v", err)
			}

			if test.expectedErrorCode != lexErr.code {
				t.Fatalf("unexpected error code, %v got %v", test.expectedErrorCode, lexErr.code)
			}
		})
	}
}

func TestNextNumber(t *testing.T) {
	tests := []struct{
		name													string
		query													string
		expectedResponse  						Token
		expectedErrorCode 						ErrorCode
		expectedErrorLineNumber				int
		expectedErrorCharStartNumber	int
		expectedErrorCharEndNumber		int
	} {
		{
			name: "number",
			query: "12",
			expectedResponse: Token{ Type: NUMBER, Text: "12" },
			expectedErrorCode: -1,
			expectedErrorLineNumber: -1,
			expectedErrorCharStartNumber: -1,
			expectedErrorCharEndNumber: -1,
		},
		{
			name: "negative_number",
			query: "-12",
			expectedResponse: Token{ Type: NUMBER, Text: "-12" },
			expectedErrorCode: -1,
			expectedErrorLineNumber: -1,
			expectedErrorCharStartNumber: -1,
			expectedErrorCharEndNumber: -1,
		},
		{
			name: "invalid_negative_number",
			query: "1-2",
			expectedResponse: Token{},
			expectedErrorCode: INVALID_QUERY_LEX_NUMBER,
			expectedErrorLineNumber: 1,
			expectedErrorCharStartNumber: 0,
			expectedErrorCharEndNumber: 2,
		},
		{
			name: "float",
			query: "14.23",
			expectedResponse: Token{ Type: NUMBER, Text: "14.23" },
			expectedErrorCode: -1,
			expectedErrorLineNumber: -1,
			expectedErrorCharStartNumber: -1,
			expectedErrorCharEndNumber: -1,
		},
		{
			name: "invalid_float",
			query: "14.2.3",
			expectedResponse: Token{},
			expectedErrorCode: INVALID_QUERY_LEX_NUMBER,
			expectedErrorLineNumber: 1,
			expectedErrorCharStartNumber: 0,
			expectedErrorCharEndNumber: 6,
		},
		{
			name: "line_break",
			query: "14\n3",
			expectedResponse: Token{ Type: NUMBER, Text: "14", },
			expectedErrorCode: -1,
			expectedErrorLineNumber: -1,
			expectedErrorCharStartNumber: -1,
			expectedErrorCharEndNumber: -1,
		},
		{
			name: "comma",
			query: "14,3",
			expectedResponse: Token{ Type: NUMBER, Text: "14", },
			expectedErrorCode: -1,
			expectedErrorLineNumber: -1,
			expectedErrorCharStartNumber: -1,
			expectedErrorCharEndNumber: -1,
		},
		{
			name: "semi_col",
			query: "14;3",
			expectedResponse: Token{ Type: NUMBER, Text: "14", },
			expectedErrorCode: -1,
			expectedErrorLineNumber: -1,
			expectedErrorCharStartNumber: -1,
			expectedErrorCharEndNumber: -1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lexer := NewLexer()
			lexer.Write(test.query)

			b, err := lexer.readByte()
			if err != nil {
				t.Fatal(err)
			}

			response, err := lexer.nextNumber(b)

			if typedErr, ok := errors.AsType[*Error](err); ok {
				if test.expectedErrorCode != typedErr.code {
					t.Fatalf("unexpected error, %v got %v (%s)", test.expectedErrorCode, typedErr.code, typedErr.msg)
				}

				errLn, errStartCn, errEndCn := typedErr.Pos()
				if test.expectedErrorLineNumber != errLn {
					t.Fatalf("unexpected error line number, %v got %v", test.expectedErrorLineNumber, errLn)
				}

				if test.expectedErrorCharStartNumber != errStartCn {
					t.Fatalf("unexpected error char start number, %v got %v", test.expectedErrorCharStartNumber, errStartCn)
				}

				if test.expectedErrorCharEndNumber != errEndCn {
					t.Fatalf("unexpected error char end number, %v got %v", test.expectedErrorCharEndNumber, errEndCn)
				}
			}

			if err != nil && test.expectedErrorCode == -1 {
				t.Fatal(err)
			}

			if test.expectedErrorCode != -1 {
				return
			}
    	
			if test.expectedResponse.Type != response.Type {
				t.Fatalf("unexpected response type, %s got %s (%s)", test.expectedResponse.Type, response.Type, response.Text)
			}
		
			if test.expectedResponse.Text != response.Text {
				t.Fatalf("unexpected response text, %s got %s", test.expectedResponse.Text, response.Text)
			}
		})
	}
}

func TestNextString(t *testing.T) {
	tests := []struct{
		name													string
		query													string
		expectedResponse  						Token
		expectedErrorCode 						ErrorCode
		expectedErrorLineNumber				int
		expectedErrorCharStartNumber	int
		expectedErrorCharEndNumber		int
	} {
		{
			name: "normal",
			query: "'test'",
			expectedResponse: Token{ Type: STRING, Text: "test" },
			expectedErrorCode: -1,
			expectedErrorLineNumber: -1,
			expectedErrorCharStartNumber: -1,
			expectedErrorCharEndNumber: -1,
		},
		{
			name: "no_end_char",
			query: "'test",
			expectedResponse: Token{},
			expectedErrorCode: INVALID_STRING,
			expectedErrorLineNumber: 1,
			expectedErrorCharStartNumber: 1,
			expectedErrorCharEndNumber: 5,
		},
		{
			name: "rand_char",
			query: "'4&€g_*^2ı}!$?`#'",
			expectedResponse: Token{ Type: STRING, Text: "4&€g_*^2ı}!$?`#" },
			expectedErrorCode: -1,
			expectedErrorLineNumber: -1,
			expectedErrorCharStartNumber: -1,
			expectedErrorCharEndNumber: -1,
		},
		{
			name: "line_break",
			query: "'test\ntest'",
			expectedResponse: Token{ Type: STRING, Text: "test\ntest" },
			expectedErrorCode: -1,
			expectedErrorLineNumber: -1,
			expectedErrorCharStartNumber: -1,
			expectedErrorCharEndNumber: -1,
		},
		{
			name: "semi_col",
			query: "'test;test'",
			expectedResponse: Token{ Type: STRING, Text: "test;test" },
			expectedErrorCode: -1,
			expectedErrorLineNumber: -1,
			expectedErrorCharStartNumber: -1,
			expectedErrorCharEndNumber: -1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lexer := NewLexer()
			lexer.Write(test.query)

			// consume start char
			_, err := lexer.readByte()
			if err != nil {
				t.Fatal(err)
			}

			response, err := lexer.nextString()

			if typedErr, ok := errors.AsType[*Error](err); ok {
				if test.expectedErrorCode != typedErr.code {
					t.Fatalf("unexpected error, %v got %v (%s)", test.expectedErrorCode, typedErr.code, typedErr.msg)
				}

				errLn, errStartCn, errEndCn := typedErr.Pos()
				if test.expectedErrorLineNumber != errLn {
					t.Fatalf("unexpected error line number, %v got %v", test.expectedErrorLineNumber, errLn)
				}

				if test.expectedErrorCharStartNumber != errStartCn {
					t.Fatalf("unexpected error char start number, %v got %v", test.expectedErrorCharStartNumber, errStartCn)
				}

				if test.expectedErrorCharEndNumber != errEndCn {
					t.Fatalf("unexpected error char end number, %v got %v", test.expectedErrorCharEndNumber, errEndCn)
				}
			}

			if err != nil && test.expectedErrorCode == -1 {
				t.Fatal(err)
			}

			if test.expectedErrorCode != -1 {
				return
			}
    	
			if test.expectedResponse.Type != response.Type {
				t.Fatalf("unexpected response type, %s got %s (%s)", test.expectedResponse.Type, response.Type, response.Text)
			}
		
			if test.expectedResponse.Text != response.Text {
				t.Fatalf("unexpected response text, %s got %s", test.expectedResponse.Text, response.Text)
			}
		})
	}
}

func TestNextKeyword(t *testing.T) {
	tests := []struct{
		name													string
		query													string
		expectedResponse  						Token
		expectedErrorCode 						ErrorCode
		expectedErrorLineNumber				int
		expectedErrorCharStartNumber	int
		expectedErrorCharEndNumber		int
	} {
		{
			name: "select",
			query: "select ",
			expectedResponse: Token{ Type: SELECT, Text: "select" },
			expectedErrorCode: -1,
			expectedErrorLineNumber: -1,
			expectedErrorCharStartNumber: -1,
			expectedErrorCharEndNumber: -1,
		},
		{
			name: "test_field",
			query: "test ",
			expectedResponse: Token{ Type: FIELD, Text: "test" },
			expectedErrorCode: -1,
			expectedErrorLineNumber: -1,
			expectedErrorCharStartNumber: -1,
			expectedErrorCharEndNumber: -1,
		},
		{
			name: "func",
			query: "now() ",
			expectedResponse: Token{ Type: FUNC, Text: "now" },
			expectedErrorCode: -1,
			expectedErrorLineNumber: -1,
			expectedErrorCharStartNumber: -1,
			expectedErrorCharEndNumber: -1,
		},
		{
			name: "comment_in_field",
			query: "users--test ",
			expectedResponse: Token{},
			expectedErrorCode: INVALID_QUERY_LEX_KEYWORD,
			expectedErrorLineNumber: 1,
			expectedErrorCharStartNumber: 0,
			expectedErrorCharEndNumber: 6,
		},
		{
			name: "string_term",
			query: "users'test ",
			expectedResponse: Token{},
			expectedErrorCode: INVALID_QUERY_LEX_KEYWORD,
			expectedErrorLineNumber: 1,
			expectedErrorCharStartNumber: 0,
			expectedErrorCharEndNumber: 6,
		},
		{
			name: "line_break",
			query: "users\ntest",
			expectedResponse: Token{ Type: FIELD, Text: "users" },
			expectedErrorCode: -1,
			expectedErrorLineNumber: -1,
			expectedErrorCharStartNumber: -1,
			expectedErrorCharEndNumber: -1,
		},
		{
			name: "comma",
			query: "users,test",
			expectedResponse: Token{ Type: FIELD, Text: "users" },
			expectedErrorCode: -1,
			expectedErrorLineNumber: -1,
			expectedErrorCharStartNumber: -1,
			expectedErrorCharEndNumber: -1,
		},
		{
			name: "semi_col",
			query: "users;test",
			expectedResponse: Token{ Type: FIELD, Text: "users" },
			expectedErrorCode: -1,
			expectedErrorLineNumber: -1,
			expectedErrorCharStartNumber: -1,
			expectedErrorCharEndNumber: -1,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lexer := NewLexer()
			lexer.Write(test.query)

			b, err := lexer.readByte()
			if err != nil {
				t.Fatal(err)
			}

			response, err := lexer.nextKeyword(b)

			if typedErr, ok := errors.AsType[*Error](err); ok {
				if test.expectedErrorCode != typedErr.code {
					t.Fatalf("unexpected error, %v got %v (%s)", test.expectedErrorCode, typedErr.code, typedErr.msg)
				}

				errLn, errStartCn, errEndCn := typedErr.Pos()
				if test.expectedErrorLineNumber != errLn {
					t.Fatalf("unexpected error line number, %v got %v", test.expectedErrorLineNumber, errLn)
				}

				if test.expectedErrorCharStartNumber != errStartCn {
					t.Fatalf("unexpected error char start number, %v got %v", test.expectedErrorCharStartNumber, errStartCn)
				}

				if test.expectedErrorCharEndNumber != errEndCn {
					t.Fatalf("unexpected error char end number, %v got %v", test.expectedErrorCharEndNumber, errEndCn)
				}
			}

			if err != nil && test.expectedErrorCode == -1 {
				t.Fatal(err)
			}

			if test.expectedErrorCode != -1 {
				return
			}
    	
			if test.expectedResponse.Type != response.Type {
				t.Fatalf("unexpected response type, %s got %s (%s)", test.expectedResponse.Type, response.Type, response.Text)
			}
		
			if test.expectedResponse.Text != response.Text {
				t.Fatalf("unexpected response text, %s got %s", test.expectedResponse.Text, response.Text)
			}
		})
	}
}
  
