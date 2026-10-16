package lexer

import (
	"io"
)

type Lexer struct {
	buf								[]byte
	r, p							int
	len								int
}

func NewParser(str string) *Lexer {
	return &Lexer{
		buf: []byte(str),
		r: -1,
		p: 0,
		len: len(str),
	}
}

// func (l *Lexer) parseWord(word string) error {
// 	for _, queryMod := range queryModifiers {
// 		if strings.EqualFold(word, queryMod) {
// 			p.currentContext = queryMod
// 			p.Result.Modifiers = append(p.Result.Modifiers, queryMod)
// 			return nil
// 		}
// 	}
//
// 	for _, queryType := range queryTypes {
// 		if strings.EqualFold(word, queryType) {
// 			if len(p.Result.Type) != 0 {
// 				return ErrQueryNotAllowed
// 			}
// 			p.Result.Type = queryType
// 			p.currentContext = queryType
// 			return nil
// 		}
// 	}
//
// 	switch p.currentContext {
// 	case SELECT:
// 		if err := p.parseSelectField(word); err != nil {
// 			return err
// 		}
// 	case FROM, INTO, UPDATE:
// 		if p.Result.TableName != "" {
// 				log.Printf("has tn %s, got tn %s", p.Result.TableName, word)
// 			return ErrQueryNotAllowed
// 		}
// 		p.Result.TableName = word
// 	}
//
// 	return nil
// }
//
// // parseSelectField parses a comma seperated select field string and appends them to the fields slice.
// // Returns the populated fields slice, alongside a windcard index whice is -1 if none found.
// func (l *Lexer) parseSelectField(word string) error {
// 	if word == "*" {
// 		if p.Result.WildcardIndex != -1 {
// 			return ErrMultipleWildcard 
// 		}
// 		p.Result.WildcardIndex = int16(len(p.Result.Fields))
// 	}
//
// 	p.Result.Fields = append(p.Result.Fields, word)
//
// 	return nil
// }

// returns word and whether EOF
func (l *Lexer) Next() (Token, error) {
	// Assumig each word is at mostly 10 char in len of the string based on:
	// SELECT * FROM educations WHERE user_id = 1;

	var c byte
	var err error
	for {
		c, err = l.readByte()
		if err != nil {
			return Token{ Type: EOF, Text: "" }, nil
		}

		if l.isComment(c) {
			err := l.skipUntilEOL()
			if err != nil {
				if err == io.EOF {
					return Token{ Type: EOF, Text: "" }, nil
				}
				return Token{}, err
			}

			// consume EOL char
			_, err = l.readByte()
			if err != nil {
				if err == io.EOF {
					return Token{ Type: EOF, Text: "" }, nil
				}
				return Token{}, err
			}
		}

		if l.isWhiteSpaceOrSeperatorChar(c) {
			continue
		}

		switch c {
		case ';':
			return Token{ Type: SEMICOL, Text: string(c) }, nil
		case '(':
			return Token{ Type: L_PARAN, Text: string(c) }, nil
		case ')':
			return Token{ Type: R_PARAN, Text: string(c) }, nil
		case '>':
			if l.peakByte() == '=' {
				return Token{ Type: LT_EQUAL, Text: string(c) }, nil
			}
			return Token{ Type: L_THAN, Text: string(c) }, nil
		case '<':
			if l.peakByte() == '=' {
				return Token{ Type: ST_EQUAL, Text: string(c) }, nil
			}
			return Token{ Type: S_THAN, Text: string(c) }, nil
		case '!':
			if l.peakByte() == '=' {
				return Token{ Type: N_EQUAL, Text: string(c) }, nil
			}
			return Token{}, ErrInvalidQuery
		case '=':
			return Token{ Type: EQUAL, Text: string(c) }, nil
		}

		if l.isStringStart(c) {
			str, err := l.readUntilEndOfString()
			if err != nil {
				return Token{}, err
			}

			// consume end of string char
			_, err = l.readByte()
			if err != nil {
				return Token{}, err
			}

			return Token{ Type: STRING, Text: string(str) }, nil
		}

		return l.nextKeyword()
	}
}

func (l *Lexer) nextKeyword() (Token, error) {
	var kw []byte

	var p byte
	for {
		p := l.peakByte()

		if l.isWhiteSpaceOrSeperatorChar(p) {
			break
		}

		switch p {
		case '-':
			return Token{}, ErrInvalidQuery
		case '!':
			return Token{}, ErrInvalidQuery
		case '=':
			return Token{}, ErrInvalidQuery
		case '<':
			return Token{}, ErrInvalidQuery
		case '>':
			return Token{}, ErrInvalidQuery
		case '(':
			break
		case ')':
			break
		}

		if l.isEndChar(p) {
			break
		}

		if l.isStringStart(p) {
			break
		}

		b, err := l.readByte()
		if err != nil {
			break
		}
		kw = append(kw, b)
	}

	keyword := string(kw)
	switch keyword {
	case "AND":
		return Token{ Type: AND, Text: keyword }, nil
	case "OR":
		return Token{ Type: OR, Text: keyword }, nil
	case "IS":
		return Token{ Type: IS, Text: keyword }, nil
	case "NOT":
		return Token{ Type: NOT, Text: keyword }, nil
	case "LIKE":
		return Token{ Type: LIKE, Text: keyword }, nil
	case "SELECT":
		return Token{ Type: SELECT, Text: keyword }, nil
	case "UPDATE":
		return Token{ Type: UPDATE, Text: keyword }, nil
	case "INSERT":
		return Token{ Type: INSERT, Text: keyword }, nil
	case "DELETE":
		return Token{ Type: DELETE, Text: keyword }, nil
	case "INTO":
		return Token{ Type: INTO, Text: keyword }, nil
	case "FROM":
		return Token{ Type: FROM, Text: keyword }, nil
	case "WHERE":
		return Token{ Type: WHERE, Text: keyword }, nil
	case "VALUES":
		return Token{ Type: VALUES, Text: keyword }, nil
	case "SET":
		return Token{ Type: SET, Text: keyword }, nil
	}

	if p == '(' {
		return Token{ Type: FUNC, Text: keyword }, nil
	}

	return Token{ Type: FIELD, Text: keyword }, nil
}

func (l *Lexer) isComment(c byte) bool {
	if c != '-' {
		return false
	}
	
	if l.peakByte() != '-' {
		return false
	}

	return true
}

func (l *Lexer) isWhiteSpaceOrSeperatorChar(c byte) bool {
	switch c {
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

func (l *Lexer) isEndChar(c byte) bool {
	return c == ';'
}

func (l *Lexer) isStringStart(c byte) bool {
	return l.peakByte() == '\''
}

func (l *Lexer) skipUntilEOL() error {
	for {
		b, err := l.readByte()
		if err != nil {
			return err
		}

		if b == '\n' {
			return nil
		}
	}
}

func (l *Lexer) readUntilEndOfString() ([]byte, error) {
	buf := make([]byte, 0, 20)
	for {
		b, err := l.readByte()
		if err != nil {
			return nil, ErrInvalidString
		}

		if b == '\'' {
			return buf, nil
		}

		buf = append(buf, b)
	}
}

// resets peak to read index + 1
func (l *Lexer) readByte() (byte, error) {
	l.r = l.r + 1
	l.p = l.r + 1

	if l.r >= (l.len - 1) {
		return byte(0), io.EOF
	}

	return l.buf[l.r], nil
}

func (l *Lexer) peakByte() byte {
	if l.p >= (l.len - 1) {
		return byte(0)
	}

	return l.buf[l.p]
}
