package lib

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
		INVALID_TABLE_ALTERATION: "Invalid table alteration (INVALID_TABLE_ALTERATION)",
		INVALID_TABLE_NAME: "Invalid table name (INVALID_TABLE_NAME)",

		INVALID_ALTER_NODE: "Invalid alter node (INVALID_ALTER_NODE)",
	}
)

type ErrorCode int

const (
	INVALID_TABLE_ALTERATION = iota
	INVALID_TABLE_NAME

	INVALID_ALTER_NODE
)


