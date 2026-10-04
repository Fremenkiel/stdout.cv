package main

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

const (
	SELECT string = "SELECT"
	DELETE string = "DELETE"
	INSERT string = "INSERT"
	UPDATE string = "UPDATE"
)

var queryTypeMap map[string]struct{} = map[string]struct{}{
	SELECT: {},
	DELETE: {},
	INSERT: {},
	UPDATE: {},
}

var queryTypeSlice []string = []string{
	SELECT,
	DELETE,
	INSERT,
	UPDATE,
}

var regexMatch *regexp.Regexp = regexp.MustCompile("(?i)SELECT|UPDATE|DELETE|INSERT")

const query string = "select name, age, email, address, zip_code, city, country, phone_number, password_hash, FROM users WHERE country LIKE 'mark' AND age >= 25 AND age <= 30 AND email LIKE 'gmail' AND first_name = '%s';"

func BenchmarkPerformanceEql1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		q := fmt.Sprintf(query, rand.Text())

		qParts := strings.Split(q, " ")

		var typeFound bool
		for pi := range qParts {
			for si := range queryTypeSlice {
				if strings.EqualFold(qParts[pi], queryTypeSlice[si]) {
					typeFound = true
				}
			}
		}
		if !typeFound {
			b.Fatal("No type found")
		}
	}
}

func BenchmarkPerformanceEql2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		q := fmt.Sprintf(query, rand.Text())

		qParts := strings.Split(q, " ")

		var typeFound bool
		for _, part := range qParts {
			for si := range queryTypeSlice {
				if strings.EqualFold(part, queryTypeSlice[si]) {
					typeFound = true
				}
			}
		}
		if !typeFound {
			b.Fatal("No type found")
		}
	}
}

func BenchmarkPerformanceEql3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		q := fmt.Sprintf(query, rand.Text())

		qParts := strings.Split(q, " ")

		var typeFound bool
		for pi := range qParts {
			for _, qt := range queryTypeSlice {
				if strings.EqualFold(qParts[pi], qt) {
					typeFound = true
				}
			}
		}
		if !typeFound {
			b.Fatal("No type found")
		}
	}
}

func BenchmarkPerformanceEql4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		q := fmt.Sprintf(query, rand.Text())

		qParts := strings.Split(q, " ")

		var typeFound bool
		for _, part := range qParts {
			for _, qt := range queryTypeSlice {
				if strings.EqualFold(part, qt) {
					typeFound = true
				}
			}
		}
		if !typeFound {
			b.Fatal("No type found")
		}
	}
}

func BenchmarkPerformanceToUpper(b *testing.B) {
	for i := 0; i < b.N; i++ {
		q := fmt.Sprintf(query, rand.Text())

		qParts := strings.Split(q, " ")

		var typeFound bool
		for pi := range qParts {
			uStr := strings.ToUpper(qParts[pi])
			if _, ok := queryTypeMap[uStr]; ok {
				typeFound = true
			}
		}
		if !typeFound {
			b.Fatal("No type found")
		}
	}
}

func BenchmarkPerformanceRegex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		q := fmt.Sprintf(query, rand.Text())

		qParts := strings.Split(q, " ")

		var typeFound bool
		for pi := range qParts {
			if ok := regexMatch.MatchString(qParts[pi]); ok {
				typeFound = true
			}
		}
		if !typeFound {
			b.Fatal("No type found")
		}
	}
}
