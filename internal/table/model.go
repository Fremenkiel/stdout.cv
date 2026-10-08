package table

import "github.com/fremenkiel/stdout.cv/internal/column"

type Table struct {
	Name			string
	RowCount	int16
	Columns		[]*column.Column
}
