package viewmodels

type ResultViewData Result

type Result struct {
	Stats		*Stats
	Columns	ColumnList
	Rows		RowList
}

type Stats struct {
	RowCount,
	ColumnCount		uint16
	Duration			int64
}
