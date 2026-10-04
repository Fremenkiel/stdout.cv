package lib

import (
	"errors"

	"modernc.org/sqlite"
	"modernc.org/sqlite/lib"
)

func ParseSqliteError(err error) (*Error, bool) {
	if sqliteError, ok := errors.AsType[*sqlite.Error](err); ok {
		switch sqliteError.Code() {
		case sqlite3.SQLITE_ERROR, 
			sqlite3.SQLITE_ABORT,
			sqlite3.SQLITE_CONSTRAINT,
			sqlite3.SQLITE_INTERRUPT,
			sqlite3.SQLITE_LOCKED,
			sqlite3.SQLITE_PERM,
			sqlite3.SQLITE_RANGE,
			sqlite3.SQLITE_READONLY,
			sqlite3.SQLITE_ROW,
			sqlite3.SQLITE_SCHEMA,
			sqlite3.SQLITE_TOOBIG,
			sqlite3.SQLITE_WARNING:
			return &Error{
				Title: "SQL query error",
				Message: sqliteError.Error(),
			}, true
		}
	}

	return nil, false
}
