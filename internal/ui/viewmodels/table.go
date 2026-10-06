package viewmodels

import "github.com/fremenkiel/stdout.cv/internal/platform/lib"

type TableViewData struct {
	Error					*lib.Error
	Tables				[]Table
	Result				*Result
}

type Table struct {
	Name 		string
	Schema	ColumnList
}
