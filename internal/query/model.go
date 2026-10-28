package query

import (
	"github.com/fremenkiel/stdout.cv/internal/database"
	"github.com/fremenkiel/stdout.cv/internal/schema"
)

type Query struct {
	Schema		[]*schema.Table
	Stats			*database.Stats
	Result		*Result
}

type Result struct {
	Rows			[]*Row
	Columns		[]string
	Duration	int64
}

type Row []*string
