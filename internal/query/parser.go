package query

import (
	"strings"

	"github.com/fremenkiel/stdout.cv/internal/column"
)

type Parser struct {
	Result						*Query
	buf								[]byte
	r, p							int
	len								int
	currentContext		string
}

func NewParser(str string) *Parser {
	return &Parser{
		Result: &Query{
			WildcardIndex: -1,
		},
		buf: []byte(str),
		r: -1,
		p: 0,
		len: len(str),
		currentContext: "",
	}
}

// Parses query string to Query struct.
// Fields are only populated on queries that would return fields.
// Multiple queries are not supported.
func Parse(query string) (*Query, error) {
	if len(query) <= 0 {
		return nil, ErrNoQuery
	}

	var endMarkerReached bool

	parser := NewParser(query)
	for !parser.atEnd() {
		w, eof := parser.readWord()

		if eof && endMarkerReached {
			return nil, ErrMultipleQueries
		}

		if len(w) > 0 {
			if err := parser.parseWord(string(w)); err != nil {
				return nil, err
			}
		}

		if eof {
			endMarkerReached = true
		}
	}

	if !endMarkerReached {
		return nil, ErrQueryEndMarkerMissing
	}

	return parser.Result, nil
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

func (p *Parser) parseWord(word string) error {
	for _, queryMod := range queryModifiers {
		if strings.EqualFold(word, queryMod) {
			p.currentContext = queryMod
			p.Result.Modifiers = append(p.Result.Modifiers, queryMod)
			return nil
		}
	}

	for _, queryType := range queryTypes {
		if strings.EqualFold(word, queryType) {
			if p.Result.Type != "" {
				return ErrQueryNotAllowed
			}
			p.Result.Type = queryType
			p.currentContext = queryType
			return nil
		}
	}

	switch p.currentContext {
	case SELECT:
		if err := p.parseSelectField(word); err != nil {
			return err
		}
	case FROM, INTO, UPDATE:
		if p.Result.TableName != "" {
			return ErrQueryNotAllowed
		}
		p.Result.TableName = word
	}

	return nil
}

// parseSelectField parses a comma seperated select field string and appends them to the fields slice.
// Returns the populated fields slice, alongside a windcard index whice is -1 if none found.
func (p *Parser) parseSelectField(word string) error {
	if word == "*" {
		if p.Result.WildcardIndex != -1 {
			return ErrMultipleWildcard 
		}
		p.Result.WildcardIndex = int16(len(p.Result.Fields))
	}

	p.Result.Fields = append(p.Result.Fields, word)

	return nil
}

// returns word and whether EOF
func (p *Parser) readWord() ([]byte, bool) {
	// Assumig each word is at mostly 10 char in len of the string based on:
	// SELECT * FROM educations WHERE user_id = 1;
	buf := make([]byte, 0, 10)

	for !p.atEnd() {
		if p.isComment() {
			_ = p.readUntil('\n')

			// consume the new line char
			_ = p.readByte()

			if len(buf) != 0 {
				return buf, false
			}
		}

		if p.isWhiteSpaceOrSeperatorChar() {
			// consume white space char
			_ = p.readByte()

			if len(buf) != 0 {
				return buf, false
			}
		}

		if p.isEndChar() {
			// consume end char
			_ = p.readByte()

			return buf, true
		}

		buf = append(buf, p.readByte())
	}

	return buf, false
}

func (p *Parser) isComment() bool {
	if p.peakByte() != '-' {
		return false
	}
	
	if atEnd := p.advancePeak(); atEnd {
		return false
	}

	if p.peakByte() != '-' {
		p.resetPeak()
		return false
	}

	return true
}

func (p *Parser) isWhiteSpaceOrSeperatorChar() bool {
	switch p.peakByte() {
	case '\n':
		return true
	case ' ':
		return true
	case '\t':
		return true
	case ',':
		return true
	}
	return false
}

func (p *Parser) isEndChar() bool {
	return p.peakByte() == ';'
}

func (p *Parser) readUntil(char byte) []byte {
	if p.atEnd() {
		return []byte{}
	}

	buf := make([]byte, 0, 10)

	b := p.readByte()
	for b != char && !p.atEnd() {
		buf = append(buf, b)

		b = p.readByte()
	}

	return buf
}

func (p *Parser) readByte() byte {
	p.r = p.r + 1

	p.resetPeak()

	if p.atEnd() {
		return byte(0)
	}

	return p.buf[p.r]
}

func (p *Parser) resetPeak() {
	p.p = p.r

	// assuming reset is only called when peak has been advanced beyond r + 1
	_ = p.advancePeak()
}

// returns whether or not EOF has been reached
func (p *Parser) advancePeak() bool {
	if pi := p.p + 1; pi < p.len {
		p.p = pi
		return false
	}
	return true
}

func (p *Parser) peakByte() byte {
	return p.buf[p.p]
}

func (p *Parser) atEnd() bool {
	return p.r >= (p.len - 1)
}
