package interpreter

import (
	"errors"
	"fmt"
)

// Lexer
var ErrInvalidQuery = errors.New("Query is invalid formatted")


type Error struct {
	msg  string
	code ErrorCode
}

var _ error = (*Error)(nil)

func (e *Error) Error() string { return e.msg }
func (e *Error) Code() int { return int(e.code) }

func NewError(code ErrorCode, line, pos int, str string) *Error {
	return &Error{
		msg: fmt.Sprintf("line %v,%v | %s (%s)",
			line, pos, parseErrorCode(code), str,
			),
		code: code,
	}
}

var (
	ErrorCodeString = map[ErrorCode]string{
		NoQuery: "No query given",

		InvalidInsertStatement: "Invalid INSERT statement, add",
		InvalidSelectStatement: "Invalid SELECT statement, add",
		InvalidUpdateStatement: "Invalid UPDATE statement, add",
		InvalidDeleteStatement: "Invalid DELETE statement, add",

		InvalidString: "String is not correctly closed, add termination to",
		InvalidListFormat: "Invalid list format, add",
		InvalidQuery: "Query is invalid formatted",
		InvalidField: "Unexpected token type when trying to parse field, got",
		InvalidTypeForPrecedence: "Invalid token type for expression, precedence not set for",
		InvalidEndOfStatement: "Invalid end of statement, add",
		InvalidBaseNodeType: "Invalid base node type",
		InvalidFuncChar: "Invalid func char, add",
		InvalidOperatiorToken: "Invalid operator token",

		UnexpectedEOF: "Unexpected end of file after",
		UnknownBaseNodeType: "Unknown base node type recived, this might be a resulte of a falty query before the current",
	}
)

type ErrorCode int

const (
	NoQuery ErrorCode = iota

	InvalidInsertStatement
	InvalidSelectStatement
	InvalidUpdateStatement
	InvalidDeleteStatement

	InvalidString
	InvalidListFormat
	InvalidQuery
	InvalidField
	InvalidTypeForPrecedence
	InvalidEndOfStatement
	InvalidBaseNodeType
	InvalidFuncChar
	InvalidOperatiorToken

	UnexpectedEOF
	UnknownBaseNodeType
)

func parseErrorCode(code ErrorCode) string {
	return ErrorCodeString[code]
}

