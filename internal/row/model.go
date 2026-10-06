package row

type Result struct {
	Rows			[]*Row
	Columns		[]string
	Duration	int64
}

type Row []*string
