package query

type Query struct {
	Rows			[]*Row
	Columns		[]string
	Duration	int64
}

type Row []*string
