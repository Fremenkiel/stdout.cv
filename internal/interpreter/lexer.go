package interpreter

import (
	"io"
)

type Lexer struct {
	buf					[]byte
	r,
	len,

	// used for error messages
	currentChar,
	currentLine	int
}

func NewLexer() *Lexer {
	return &Lexer{
		buf: nil,
		r: -1,
		len: -1,
		currentChar: 0,
		currentLine: 1,
	}
}

// write query to buffer
func (l *Lexer) Write(str string) {
	l.buf = []byte(str)
	l.len = len(l.buf)
}

func (l *Lexer) Next() (Token, error) {
	if l.len == -1 {
		return Token{}, NewError(NoQuery, -1, -1, "")
	}

	var c byte
	var err error
	for {
		c, err = l.readByte()
		if err != nil {
			if err == io.EOF {
				return NewToken(EOF, "", l), nil
			}
			return Token{}, err
		}

		if l.isComment(c) {
			err := l.skipUntilEOL()
			if err != nil {
				if err == io.EOF {
					return NewToken(EOF, "", l), nil
				}
				return Token{}, err
			}

			// consume EOL char
			_, err = l.readByte()
			if err != nil {
				if err == io.EOF {
					return NewToken(EOF, "", l), nil
				}
				return Token{}, err
			}
			l.addLineCount()
		}

		if l.isLineBreak(c) {
			l.addLineCount()
			continue
		}

		if l.isWhiteSpaceOrSeperatorChar(c) {
			continue
		}

		if l.isNumber(c) {
			return l.nextNumber(c)
		}

		if l.isStringMarker(c) {
			return l.nextString()
		}

		switch c {
		case ';':
			return NewToken(SEMICOL, string(c), l), nil
		case '(':
			return NewToken(L_PARAN, string(c), l), nil
		case ')':
			return NewToken(R_PARAN, string(c), l), nil
		case '>':
			if l.peakByte() == '=' {
				b, _ := l.readByte()
				return NewToken(GT_EQUAL, string([]byte{c, b}), l), nil
			}
			return NewToken(G_THAN, string(c), l), nil
		case '<':
			if l.peakByte() == '=' {
				b, _ := l.readByte()
				return NewToken(LT_EQUAL, string([]byte{c, b}), l), nil
			}
			return NewToken(L_THAN, string(c), l), nil
		case '!':
			if l.peakByte() == '=' {
				b, _ := l.readByte()
				return NewToken(N_EQUAL, string([]byte{c, b}), l), nil
			}
			return Token{}, ErrInvalidQuery
		case '=':
			return NewToken(EQUAL, string(c), l), nil
		}

		return l.nextKeyword(c)
	}
}

func (l *Lexer) addLineCount() {
	l.currentChar = 0
	l.currentLine = l.currentLine + 1
}

func (l *Lexer) nextNumber(char byte) (Token, error) {
	var num = []byte{char}

	var p byte
	for {
		p = l.peakByte()
		if !l.isInt(p) {
			break
		}

		b, err := l.readByte()
		if err != nil {
			break
		}
		num = append(num, b)
	}

	return NewToken(NUMBER, string(num), l), nil
}

func (l *Lexer) nextString() (Token, error) {
	str, err := l.readUntilEndOfString()
	if err != nil {
		if err == io.EOF {
			return NewToken(EOF, "", l), nil
		}
		return Token{}, err
	}

	// consume end of string char
	_, err = l.readByte()
	if err != nil {
		return Token{}, err
	}

	return NewToken(STRING, string(str), l), nil
}

func (l *Lexer) nextKeyword(char byte) (Token, error) {
	var kw = []byte{char}

	var p byte
	outerloop:
	for {
		p = l.peakByte()

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
			break outerloop
		case ')':
			break outerloop
		}

		if l.isEndChar(p) {
			break
		}

		if l.isStringMarker(p) {
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
	case string(AND):
		return NewToken(AND, keyword, l), nil
	case string(OR):
		return NewToken(OR, keyword, l), nil
	case string(IS):
		return NewToken(IS, keyword, l), nil
	case string(NOT):
		return NewToken(NOT, keyword, l), nil
	case string(LIKE):
		return NewToken(LIKE, keyword, l), nil
	case string(SELECT):
		return NewToken(SELECT, keyword, l), nil
	case string(UPDATE):
		return NewToken(UPDATE, keyword, l), nil
	case string(INSERT):
		return NewToken(INSERT, keyword, l), nil
	case string(DELETE):
		return NewToken(DELETE, keyword, l), nil
	case string(INTO):
		return NewToken(INTO, keyword, l), nil
	case string(FROM):
		return NewToken(FROM, keyword, l), nil
	case string(WHERE):
		return NewToken(WHERE, keyword, l), nil
	case string(VALUES):
		return NewToken(VALUES, keyword, l), nil
	case string(SET):
		return NewToken(SET, keyword, l), nil
	case string(NULL):
		return NewToken(NULL, keyword, l), nil
	case string(AS):
		return NewToken(AS, keyword, l), nil
	}

	if p == '(' {
		return NewToken(FUNC, keyword, l), nil
	}

	return NewToken(FIELD, keyword, l), nil
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
	case ' ':
		return true
	case '\t':
		return true
	case ',':
		return true
	}
	return false
}

func (l *Lexer) isInt(c byte) bool {
	if c == '-' {
		p := l.peakByte()
		return p > 0 && p <= 9
	}

	return l.isNumber(c)
}

func (l *Lexer) isNumber(c byte) bool {
	c -= '0'
	if c > 9 {
		return false
	}

	return true
}

func (l *Lexer) readUntilEndOfString() ([]byte, error) {
	buf := make([]byte, 0, 20)
	for {
		p := l.peakByte()
		if l.isStringMarker(p) {
			return buf, nil
		}

		b, err := l.readByte()
		if err != nil {
			return nil, NewError(InvalidString, l.currentLine, l.r - len(buf), string(buf))
		}

		buf = append(buf, b)
	}
}

func (l *Lexer) isLineBreak(c byte) bool {
	return  c == '\n'
}

func (l *Lexer) isEndChar(c byte) bool {
	return c == ';'
}

func (l *Lexer) isStringMarker(c byte) bool {
	return c == '\''
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

// resets peak to read index + 1
func (l *Lexer) readByte() (byte, error) {
	if l.r >= (l.len - 1) {
		return byte(0), io.EOF
	}

	l.r = l.r + 1

	// set current char for error messages
	l.currentChar = l.r + 1

	return l.buf[l.r], nil
}

func (l *Lexer) peakByte() byte {
	if (l.r + 1) >= l.len {
		return byte(0)
	}

	return l.buf[l.r + 1]
}
