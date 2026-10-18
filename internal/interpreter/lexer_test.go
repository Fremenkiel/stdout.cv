package interpreter

import (
	"errors"
	"testing"
)

func TestNextNoQuery(t *testing.T) {
	tests := []struct{
		name					string
		expectedResponse	any
		expectedError			ErrorCode
	} {
		{
			name: "no_query",
			expectedResponse: nil,
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
