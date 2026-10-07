package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

var tableNames []string = []string{
	"test1",
	"test2",
	"test3",
	"test4",
	"test5",
	"test6",
	"test78",
	"testfgs",
	"testgsd",
	"test6433",
	"testrter",
	"testjfgwr",
	"test556e",
	"test56wt",
	"testy54y34wrg",
	"test4tw",
	"test4tw",
	"test1j4wt34",
	"test1rgfrweg",
}

var globalQuery string

func BenchmarkPerformanceConcat1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		querySlice := make([]string, len(tableNames))

		for i, name := range tableNames {
			querySlice[i] = fmt.Sprintf(`
				SELECT * FROM pragma_table_info('%s');
				`, name)
		}

		query := strings.Join(querySlice, "")

		globalQuery = query
	}
}

func BenchmarkPerformanceConcat2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var query string

		for _, name := range tableNames {
			query += fmt.Sprintf(`
				SELECT * FROM pragma_table_info('%s');
				`, name)
		}

		globalQuery = query
	}
}

func BenchmarkPerformanceConcat3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var query string

		for _, name := range tableNames {
			query = fmt.Sprintf(`
				%sSELECT * FROM pragma_table_info('%s');
				`, query, name)
		}

		globalQuery = query
	}
}

func BenchmarkPerformanceConcat4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer

		for _, name := range tableNames {
			buf.WriteString(fmt.Sprintf(`
				SELECT * FROM pragma_table_info('%s');
				`, name))
		}

		globalQuery = buf.String()
	}
}

func BenchmarkPerformanceConcat5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var builder strings.Builder

		for _, name := range tableNames {
			builder.WriteString(fmt.Sprintf(`
				SELECT * FROM pragma_table_info('%s');
				`, name))
		}

		globalQuery = builder.String()
	}
}
