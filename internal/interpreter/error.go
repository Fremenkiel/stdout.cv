package interpreter

import (
	"fmt"
)

type Error struct {
	msg  				string
	code 				ErrorCode
	lineNumber,
	startChar,
	endChar			int
}

var _ error = (*Error)(nil)

func (e *Error) Error() string { return e.msg }
func (e *Error) Code() ErrorCode { return e.code }
func (e *Error) Pos() (lineNumber, startChar, endChar int) { return e.lineNumber, e.startChar, e.endChar }

func NewError(code ErrorCode, line, endPos int, str string) *Error {
	return &Error{
		msg: fmt.Sprintf("line %v,%v | %s ´%s´",
			line, endPos, ErrorCodeString[code], str,
			),
		code: code,
		lineNumber: line,
		startChar: endPos - len(str),
		endChar: endPos,
	}
}

var (
	ErrorCodeString = map[ErrorCode]string{
		NO_QUERY:											"No query given (NO_QUERY)",
		UNSUPPORTED_QUERY:						"Query type is not supported (UNSUPPORTED_QUERY)",

		INVALID_INSERT_STATEMENT:			"INSERT statement is invalidly formatted (INVALID_INSERT_STATEMENT)",
		INVALID_SELECT_STATEMENT:			"SELECT statement is invalidly formatted (INVALID_SELECT_STATEMENT)",
		INVALID_UPDATE_STATEMENT:			"UPDATE statement is invalidly formatted (INVALID_UPDATE_STATEMENT)",
		INVALID_DELETE_STATEMENT:			"DELETE statement is invalidly formatted (INVALID_DELETE_STATEMENT)",

		INVALID_QUERY:								"Query is invalidly formatted (INVALID_QUERY)",
		INVALID_QUERY_LEX_NUMBER:			"Query formatted error while lexing number (INVALID_QUERY_LEX_NUMBER)",
		INVALID_QUERY_LEX_KEYWORD:		"Query formatted error while lexing keyword (INVALID_QUERY_LEX_KEYWORD)",

		INVALID_STRING:								"String is not correctly formatted, add termination (INVALID_STRING)",
		INVALID_LIST_FORMAT:					"Invalid list format (INVALID_LIST_FORMAT)",
		INVALID_FIELD:								"Unexpected token type while trying to parse field 8INVALID_FIELD)",
		INVALID_TYPE_FOR_PRECEDENCE:	"Invalid token type for expression, no precedence set (INVALID_TYPE_FOR_PRECEDENCE)",
		INVALID_END_OF_STATEMENT:			"Statement is invalidly ended (INVALID_END_OF_STATEMENT)",
		INVALID_BASE_NODE_TYPE:				"Base node type in not valid (INVALID_BASE_NODE_TYPE)",
		INVALID_FUNC_CHAR:						"Func call is invalidly formatted (INVALID_FUNC_CHAR)",
		INVALID_OPERATOR_TOKEN:				"Token is invalid as operator (INVALID_OPERATOR_TOKEN)",

		UNEXPECTED_EOF:								"Unexpected end of file (UNEXPECTED_EOF)",
		UNKNOWN_BASE_NODE_TYPE:				"Unknown base node type recived, this might be a resulte of a falty part earlier in the query (UNKNOWN_BASE_NODE_TYPE)",
	}
)

type ErrorCode int

const (
	NO_QUERY ErrorCode = iota
	UNSUPPORTED_QUERY

	INVALID_INSERT_STATEMENT
	INVALID_SELECT_STATEMENT
	INVALID_UPDATE_STATEMENT
	INVALID_DELETE_STATEMENT

	INVALID_QUERY
	INVALID_QUERY_LEX_NUMBER
	INVALID_QUERY_LEX_KEYWORD

	INVALID_STRING
	INVALID_LIST_FORMAT
	INVALID_FIELD
	INVALID_TYPE_FOR_PRECEDENCE
	INVALID_END_OF_STATEMENT
	INVALID_BASE_NODE_TYPE
	INVALID_FUNC_CHAR
	INVALID_OPERATOR_TOKEN

	UNEXPECTED_EOF
	UNKNOWN_BASE_NODE_TYPE
)

