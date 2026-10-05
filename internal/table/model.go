package table

import "github.com/fremenkiel/stdout.cv/internal/platform/lib"

type Table struct {
	Name	string
}

type TableListViewData struct {
	DatabaseName	string
	Error					*lib.Error
	Tables				[]TableViewData
	Rows					*[]any
}

type TableViewData struct {
	Name string
}
