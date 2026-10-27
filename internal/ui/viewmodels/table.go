package viewmodels

import "github.com/fremenkiel/stdout.cv/internal/platform/lib"

type TableViewData struct {
	Error					*lib.Error
	Tables				[]Table
	Stats					DatabaseStats
	Result				*Result
}

type Table struct {
	Name 			string
	RowCount	int16
	Columns		*ColumnList
}

type DatabaseStats struct {
	Version	string
	Size		float64
}
