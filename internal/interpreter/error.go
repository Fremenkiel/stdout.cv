package interpreter

import (
	"errors"
	"fmt"
)

var ErrQueryNotAllowed = errors.New("query: The received query is not allowed")
var ErrMultipleWildcard = errors.New("query: More than one wildcard char was given")
var ErrMultipleQueries = errors.New("query: More than one query was given")
var ErrQueryEndMarkerMissing = errors.New("query: End marker for query is missing, add a semi-colon")
var ErrNoQuery = errors.New("query: No query received")

// Lexer
var ErrInvalidString = errors.New("String is not correctly closed, add end to string")
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

		InvalidString: "String is not correctly closed, add termination to",
		InvalidQuery: "Query is invalid formatted",
		InvalidInsertStatement: "Invalid INSERT statement, add",
		InvalidField: "Unexpected token type when trying to parse field, got",

		UnexpectedEOF: "Unexpected end of file after",
		UnknownBaseNodeType: "Unknown base node type recived, this might be a resulte of a falty query before the current",
	}
)

type ErrorCode int

const (
	NoQuery ErrorCode = iota

	InvalidString
	InvalidQuery
	InvalidInsertStatement
	InvalidField

	UnexpectedEOF
	UnknownBaseNodeType
)

func parseErrorCode(code ErrorCode) string {
	return ErrorCodeString[code]
}

