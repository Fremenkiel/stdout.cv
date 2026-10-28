package viewmodels

type Result struct {
	Stats		*ResultStats
	Columns	ColumnList
	Rows		RowList
}

type ResultStats struct {
	RowCount,
	ColumnCount		uint16
	Duration			int64
}
