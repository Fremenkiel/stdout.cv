package query

import (
	"strings"

	"github.com/fremenkiel/stdout.cv/internal/schema"
)

// Parses query string to Query struct. 
// Fields are only populated on queries that would return fields.
// Multiple queries are not supported.
func Parse(queryString string) (*Query, error) {
	queryParts := strings.Split(queryString, " ")

	query := &Query{
		WildcardIndex: -1,
	}

	lastIndex := len(queryParts) - 1

	var current string
	Outerloop:
	for i, part := range queryParts {
		if strings.Contains(part, "--") {
			return nil, ErrQueryNotAllowed
		}
		if strings.Contains(part, ";") && i != lastIndex {
			return nil, ErrQueryNotAllowed
		}

		for _, queryMod := range queryModifiers {
			if strings.EqualFold(part, queryMod) {
				current = queryMod
				continue Outerloop
			}
		}

		for _, queryType := range queryTypes {
			if strings.EqualFold(part, queryType) {
				if query.Type != "" {
					return nil, ErrQueryNotAllowed
				}
				query.Type = queryType
				current = queryType
				continue Outerloop
			}
		}

		switch current {
		case SELECT:
			query.Fields, query.WildcardIndex = parseSelectField(query.Fields, part)
			continue
		case FROM, INTO, UPDATE:
			if query.TableName == "" {
				query.TableName = part
			}
			continue
		}
	}

	return query, nil
}

func AppendWildcardSchema(wildcardIndex int16, fields []string, schema []*schema.Schema) ([]string, error) {
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

// parseSelectField parses a comma seperated select field string and appends them to the fields slice.
// Returns the populated fields slice, alongside a windcard index whice is -1 if none found.
func parseSelectField(fields []string, str string) ([]string, int16) {
	var windcardIndex int16 = -1

	selectFields := strings.Split(str, ",")
	for i, field := range selectFields {
		if len(field) == 0 {
			continue
		}

		if field == "*" {
			windcardIndex = int16(len(fields))
		}

		fields = append(fields, selectFields[i])
	}

	return fields, windcardIndex
}
