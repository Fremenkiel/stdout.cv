package interpreter

import "testing"


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
