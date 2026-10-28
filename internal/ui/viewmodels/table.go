package viewmodels

import "github.com/fremenkiel/stdout.cv/internal/platform/lib"

type ViewData struct {
	Error					*lib.Error
	Tables				[]*Table
	Stats					*Stats
	Result				*Result
}

type Table struct {
	Name 			string
	RowCount	int16
	Columns		*ColumnList
}

type Stats struct {
	Version	string
	Size		float64
}
