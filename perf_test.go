package main

import (
	"slices"
	"testing"
)

var valSlice []string = []string{
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
}

var valMap map[string]struct{} = map[string]struct{}{
	"test1": {},
	"test2": {},
	"test3": {},
	"test4": {},
	"test5": {},
	"test6": {},
	"test78": {},
	"testfgs": {},
	"testgsd": {},
	"test6433": {},
	"testrter": {},
	"testjfgwr": {},
	"test556e": {},
	"test56wt": {},
}

var toCheck []string = []string{
	"test4",
	"test5",
	"test6",
	"test78",
	"testy54y34wrg",
	"test4tw",
	"test1j4wt34",
	"test1rgfrweg",
}

var globalCheck bool

func BenchmarkPerformanceExists1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, val := range toCheck {
			globalCheck = slices.Contains(valSlice, val)
		}
	}
}

func BenchmarkPerformanceExists2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, val := range toCheck {
			_, ok := valMap[val]
			globalCheck = ok
		}
	}
}
