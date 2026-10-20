package interpreter

import (
	"errors"
	"testing"
)

func TestNext(t *testing.T) {
	tests := []struct{
		name							string
		query							string
		expectedResponse  Token
		expectedError     error
	} {
		{
			name: "keyword",
			query: "SeLeCt",
			expectedResponse: Token{ Type: SELECT, Text: "SeLeCt" },
			expectedError: nil,
		},
		{
			name: "string",
			query: "'test'",
			expectedResponse: Token{ Type: STRING, Text: "test" },
			expectedError: nil,
		},
		{
			name: "number",
			query: "-12",
			expectedResponse: Token{ Type: NUMBER, Text: "-12" },
			expectedError: nil,
		},
		{
			name: "field",
			query: "users",
			expectedResponse: Token{ Type: FIELD, Text: "users" },
			expectedError: nil,
		},
		{
			name: "invalid_field",
			query: "user--test",
			expectedResponse: Token{},
			expectedError: ErrInvalidQuery,
		},
		{
			name: "left_paran",
			query: "(",
			expectedResponse: Token{ Type: L_PARAN, Text: "(" },
			expectedError: nil,
		},
		{
			name: "less_than",
			query: "< =",
			expectedResponse: Token{ Type: L_THAN, Text: "<" },
			expectedError: nil,
		},
		{
			name: "less_than_or_equal",
			query: "<=",
			expectedResponse: Token{ Type: LT_EQUAL, Text: "<=" },
			expectedError: nil,
		},
		{
			name: "line_break",
			query: "\n",
			expectedResponse: Token{ Type: EOF, Text: "" },
			expectedError: nil,
		},
		{
			name: "comment",
			query: "-- test",
			expectedResponse: Token{ Type: EOF, Text: "" },
			expectedError: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lexer := NewLexer()
			lexer.Write(test.query)

			response, err := lexer.Next()

			if test.expectedError != err {
				t.Fatalf("unexpected error, %v got %v", test.expectedError, err)
			}

			if test.expectedError != nil {
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
		expectedError			ErrorCode
	} {
		{
			name: "no_query",
			expectedError: NoQuery,
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

			if test.expectedError != lexErr.code {
				t.Fatalf("unexpected error code, %v got %v", test.expectedError, lexErr.code)
			}
		})
	}
}

func TestNextNumber(t *testing.T) {
	tests := []struct{
		name							string
		query							string
		expectedResponse  Token
		expectedError     error
	} {
		{
			name: "number",
			query: "12",
			expectedResponse: Token{ Type: NUMBER, Text: "12" },
			expectedError: nil,
		},
		{
			name: "negative_number",
			query: "-12",
			expectedResponse: Token{ Type: NUMBER, Text: "-12" },
			expectedError: nil,
		},
		{
			name: "invalid_negative_number",
			query: "1-2",
			expectedResponse: Token{},
			expectedError: ErrInvalidQuery,
		},
		{
			name: "float",
			query: "14.23",
			expectedResponse: Token{ Type: NUMBER, Text: "14.23" },
			expectedError: nil,
		},
		{
			name: "invalid_float",
			query: "14.2.3",
			expectedResponse: Token{},
			expectedError: ErrInvalidQuery,
		},
		{
			name: "line_break",
			query: "14\n3",
			expectedResponse: Token{ Type: NUMBER, Text: "14", },
			expectedError: nil,
		},
		{
			name: "comma",
			query: "14,3",
			expectedResponse: Token{ Type: NUMBER, Text: "14", },
			expectedError: nil,
		},
		{
			name: "semi_col",
			query: "14;3",
			expectedResponse: Token{ Type: NUMBER, Text: "14", },
			expectedError: nil,
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

			if test.expectedError != err {
				t.Fatalf("unexpected error, %v got %v", test.expectedError, err)
			}

			if test.expectedError != nil {
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
		name							string
		query							string
		expectedResponse  Token
		expectedError     ErrorCode
	} {
		{
			name: "normal",
			query: "'test'",
			expectedResponse: Token{ Type: STRING, Text: "test" },
			expectedError: 0,
		},
		{
			name: "no_end_char",
			query: "'test",
			expectedResponse: Token{},
			expectedError: InvalidString,
		},
		{
			name: "rand_char",
			query: "'4&€g_*^2ı}!$?`#'",
			expectedResponse: Token{ Type: STRING, Text: "4&€g_*^2ı}!$?`#" },
			expectedError: 0,
		},
		{
			name: "line_break",
			query: "'test\ntest'",
			expectedResponse: Token{ Type: STRING, Text: "test\ntest" },
			expectedError: 0,
		},
		{
			name: "semi_col",
			query: "'test;test'",
			expectedResponse: Token{ Type: STRING, Text: "test;test" },
			expectedError: 0,
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

			if castError, ok := errors.AsType[*Error](err); ok && test.expectedError != castError.code {
				t.Fatalf("unexpected error, %v got %v", test.expectedError, err)
			}

			if test.expectedError != 0 {
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
		name							string
		query							string
		expectedResponse  Token
		expectedError     error
	} {
		{
			name: "select",
			query: "select ",
			expectedResponse: Token{ Type: SELECT, Text: "select" },
			expectedError: nil,
		},
		{
			name: "test_field",
			query: "test ",
			expectedResponse: Token{ Type: FIELD, Text: "test" },
			expectedError: nil,
		},
		{
			name: "func",
			query: "now() ",
			expectedResponse: Token{ Type: FUNC, Text: "now" },
			expectedError: nil,
		},
		{
			name: "comment_in_field",
			query: "users--test ",
			expectedResponse: Token{},
			expectedError: ErrInvalidQuery,
		},
		{
			name: "string_term",
			query: "users'test ",
			expectedResponse: Token{},
			expectedError: ErrInvalidQuery,
		},
		{
			name: "line_break",
			query: "users\ntest",
			expectedResponse: Token{ Type: FIELD, Text: "users" },
			expectedError: nil,
		},
		{
			name: "comma",
			query: "users,test",
			expectedResponse: Token{ Type: FIELD, Text: "users" },
			expectedError: nil,
		},
		{
			name: "semi_col",
			query: "users;test",
			expectedResponse: Token{ Type: FIELD, Text: "users" },
			expectedError: nil,
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

			if test.expectedError != err {
				t.Fatalf("unexpected error, %v got %v", test.expectedError, err)
			}

			if test.expectedError != nil {
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
  
