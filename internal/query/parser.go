package query

import (
	"strings"

	"github.com/fremenkiel/stdout.cv/internal/column"
)

type Query struct {
	Fields						[]string
	WildcardIndex			int16
	TableName					string
	Type							string
	Modifiers					[]string
	currentContext		string
	endMarkerReached	bool
}

func NewQuery() *Query {
	return &Query{
		WildcardIndex: -1,
		currentContext: "",
		endMarkerReached: false,
	}
}

// Parses query string to Query struct.
// Fields are only populated on queries that would return fields.
// Multiple queries are not supported.
func Parse(queryString string) (*Query, error) {
	lines := strings.Split(queryString, "\n")

	query := NewQuery()

	for _, line := range lines {
		if err := query.parseLine(line); err != nil {
			return nil, err
		}
	}

	if !query.endMarkerReached {
		return nil, ErrQueryEndMarkerMissing
	}

	return query, nil
}

func AppendWildcardSchema(wildcardIndex int16, fields []string, schema []*column.Column) ([]string, error) {
	if fields[wildcardIndex] != "*" {
		return nil, ErrFieldNotWildcard
	}

	endFields := fields[wildcardIndex + 1:]
	fields = fields[:wildcardIndex]

	for ti := range schema {
		fields = append(fields, schema[ti].Name)
	}

	for ei := range endFields {
		fields = append(fields, endFields[ei])
	}

	return fields, nil
}

func (q *Query) parseLine(line string) error {
	queryParts := strings.Split(line, " ")

	Outerloop:
	for _, part := range queryParts {
		if strings.Contains(part, "--") {
			return nil
		}

		if q.endMarkerReached {
			return ErrMultipleQueries
		}

		if strings.Contains(part, ";") {
			q.endMarkerReached = true
			part = strings.Trim(part, ";")
		}

		for _, queryMod := range queryModifiers {
			if strings.EqualFold(part, queryMod) {
				q.currentContext = queryMod
				q.Modifiers = append(q.Modifiers, queryMod)
				continue Outerloop
			}
		}

		for _, queryType := range queryTypes {
			if strings.EqualFold(part, queryType) {
				if q.Type != "" {
					return ErrQueryNotAllowed
				}
				q.Type = queryType
				q.currentContext = queryType
				continue Outerloop
			}
		}

		switch q.currentContext {
		case SELECT:
			if err := q.parseSelectField(part); err != nil {
				return err
			}
			continue
		case FROM, INTO, UPDATE:
			if q.TableName == "" {
				q.TableName = part
			}
			continue
		}
	}

	return nil
}


// parseSelectField parses a comma seperated select field string and appends them to the fields slice.
// Returns the populated fields slice, alongside a windcard index whice is -1 if none found.
func (q *Query) parseSelectField(str string) error {
	selectFields := strings.Split(str, ",")
	for i, field := range selectFields {
		if len(field) == 0 {
			continue
		}

		if field == "*" {
			if q.WildcardIndex != -1 {
				return ErrDuplicateWildcard
			}
			q.WildcardIndex = int16(len(q.Fields))
		}

		q.Fields = append(q.Fields, selectFields[i])
	}

	return nil
}
