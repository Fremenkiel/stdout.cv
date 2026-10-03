package table

type Table struct {
	Name	string
}

type TableListViewData struct {
	Tables	[]TableViewData
	Rows		*[]any
}

type TableViewData struct {
	Name string
}
