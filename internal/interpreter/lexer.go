package interpreter

import (
	"io"
	"strconv"
	"strings"
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
		currentChar: -1,
		currentLine: -1,
	}
}

// write query to buffer
func (l *Lexer) Write(str string) {
	l.buf = []byte(str)
	l.len = len(l.buf)
	l.currentLine = 1
	l.currentChar = 0
}

func (l *Lexer) Next() (Token, error) {
	if l.len == -1 {
		return Token{}, NewError(
			NO_QUERY, 
			-1,
			-1, 
			"",
			)
	}

	var c byte
	var err error
	for {
		c, err = l.readByte()
		if err != nil {
			if err == io.EOF {
				return l.newToken(EOF, ""), nil
			}
			return Token{}, err
		}

		switch {
		case c == '-' && l.peakByte() == '-':
			err := l.skipLine()
			if err == io.EOF {
				return l.newToken(EOF, ""), nil
			}
			continue
		case isLineBreak(c): 
			l.addLineCount()
			continue
		case isWhiteSpace(c):
			continue
		case isInt(c):
			return l.nextNumber(c)
		case isStringMarker(c):
			return l.nextString()
		case c == ',':
			return l.newToken(COMMA, string(c)), nil
		case c == ';':
			return l.newToken(SEMICOL, string(c)), nil
		case c == '(':
			return l.newToken(L_PARAN, string(c)), nil
		case c == ')':
			return l.newToken(R_PARAN, string(c)), nil
		case c == '>':
			if l.peakByte() == '=' {
				b, _ := l.readByte()
				return l.newToken(GT_EQUAL, string([]byte{c, b})), nil
			}
			return l.newToken(G_THAN, string(c)), nil
		case c == '<':
			if l.peakByte() == '=' {
				b, _ := l.readByte()
				return l.newToken(LT_EQUAL, string([]byte{c, b})), nil
			}
			return l.newToken(L_THAN, string(c)), nil
		case c == '!':
			if l.peakByte() == '=' {
				b, _ := l.readByte()
				return l.newToken(N_EQUAL, string([]byte{c, b})), nil
			}
			return Token{}, NewError(
				INVALID_QUERY,
				l.currentLine,
				l.currentChar,
				string(c),
				)
		case c == '=':
			return l.newToken(EQUAL, string(c)), nil
		case c == '*':
			if !isSeperatorChar(l.peakByte()) {
			return Token{}, NewError(
				INVALID_QUERY,
				l.currentLine,
				l.currentChar,
				string(c),
				)
			}
			return l.newToken(WILDCARD, string(c)), nil
		default:
			return l.nextKeyword(c)
		}

	}
}

func (l *Lexer) nextNumber(char byte) (Token, error) {
	var num = []byte{char}

	var isFloat bool

	outerloop:
	for {
		peak := l.peakByte()
		switch {
		case isSeperatorChar(peak):
			break outerloop
		case peak == ',':
			break outerloop
		case peak == ';':
			break outerloop
		case peak == ')':
			break outerloop
		case peak == byte(0):
		case !isNumber(peak):
			return Token{}, NewError(
				INVALID_QUERY_LEX_NUMBER,
				l.currentLine,
				l.currentChar + 1,
				string(append(num, peak)),
				)
		}

		b, err := l.readByte()
		if err != nil {
			break
		}

		if b == '.' {
			isFloat = true
		}

		num = append(num, b)
	}

	str := string(num)

	switch isFloat {
	case false:
		if _, err := strconv.ParseInt(str, 0, 0); err != nil {
			return Token{}, NewError(
				INVALID_QUERY_LEX_NUMBER,
				l.currentLine,
				l.currentChar,
				str,
				)
		}
	case true:
		if _, err := strconv.ParseFloat(str, 64); err != nil {
			return Token{}, NewError(
				INVALID_QUERY_LEX_NUMBER,
				l.currentLine,
				l.currentChar,
				str,
				)
		}
	}

	return l.newToken(NUMBER, str), nil
}

func (l *Lexer) nextString() (Token, error) {
	str, err := l.readUntilEndOfString()
	if err != nil {
		if err == io.EOF {
			return l.newToken(EOF, ""), nil
		}
		return Token{}, err
	}

	// consume end of string char
	_, err = l.readByte()
	if err != nil {
		return Token{}, err
	}

	return l.newToken(STRING, string(str)), nil
}

func (l *Lexer) nextKeyword(char byte) (Token, error) {
	var kw = []byte{char}

	var p byte
	outerloop:
	for {
		p = l.peakByte()

		switch {
		case p == '-':
			return Token{}, l.newKeywordError(string(append(kw, p)))
		case p == '!':
			return Token{}, l.newKeywordError(string(append(kw, p)))
		case p == '=':
			return Token{}, l.newKeywordError(string(append(kw, p)))
		case p == '<':
			return Token{}, l.newKeywordError(string(append(kw, p)))
		case p == '>':
			return Token{}, l.newKeywordError(string(append(kw, p)))
		case p == '*':
			return Token{}, l.newKeywordError(string(append(kw, p)))
		case isStringMarker(p):
			return Token{}, l.newKeywordError(string(append(kw, p)))
		case isSeperatorChar(p):
			break outerloop
		case p == '(':
			break outerloop
		case p == ')':
			break outerloop
		case p == ',':
			break outerloop
		}

		b, err := l.readByte()
		if err != nil {
			break
		}
		kw = append(kw, b)
	}

	keyword := string(kw)
	switch strings.ToUpper(keyword) {
	case string(AND):
		return l.newToken(AND, keyword), nil
	case string(OR):
		return l.newToken(OR, keyword), nil
	case string(IS):
		return l.newToken(IS, keyword), nil
	case string(NOT):
		return l.newToken(NOT, keyword), nil
	case string(LIKE):
		return l.newToken(LIKE, keyword), nil
	case string(SELECT):
		return l.newToken(SELECT, keyword), nil
	case string(UPDATE):
		return l.newToken(UPDATE, keyword), nil
	case string(INSERT):
		return l.newToken(INSERT, keyword), nil
	case string(DELETE):
		return l.newToken(DELETE, keyword), nil
	case string(DROP):
		return l.newToken(DROP, keyword), nil
	case string(CREATE):
		return l.newToken(CREATE, keyword), nil
	case string(INTO):
		return l.newToken(INTO, keyword), nil
	case string(FROM):
		return l.newToken(FROM, keyword), nil
	case string(WHERE):
		return l.newToken(WHERE, keyword), nil
	case string(VALUES):
		return l.newToken(VALUES, keyword), nil
	case string(SET):
		return l.newToken(SET, keyword), nil
	case string(NULL):
		return l.newToken(NULL, keyword), nil
	case string(AS):
		return l.newToken(AS, keyword), nil
	case string(IF):
		return l.newToken(IF, keyword), nil
	case string(EXISTS):
		return l.newToken(EXISTS, keyword), nil
	case string(TABLE):
		return l.newToken(TABLE, keyword), nil
	case string(ALTER):
		return l.newToken(ALTER, keyword), nil
	case string(COLUMN):
		return l.newToken(COLUMN, keyword), nil
	case string(RENAME):
		return l.newToken(RENAME, keyword), nil
	case string(ADD):
		return l.newToken(ADD, keyword), nil
	case string(PRIMARY):
		return l.newToken(PRIMARY, keyword), nil
	case string(KEY):
		return l.newToken(KEY, keyword), nil
	}

	return l.newToken(FIELD, keyword), nil
}

func (l *Lexer) newKeywordError(str string) error {
	return NewError(
		INVALID_QUERY_LEX_KEYWORD,
		l.currentLine,
		(l.currentChar + 1),
		str,
		)
}

func (l *Lexer) newToken(t TokenType, text string) Token {
	endPos := l.currentChar

	// removing the string marker as it's ignored in the string
	if t == STRING {
		endPos = endPos - 1
	}

	startPos := endPos - len(text)

	return Token{ 
		Type: t,
		Text: text,
		LineNumber: l.currentLine,
		Start: startPos,
		End: endPos,
	}
}

func (l *Lexer) readUntilEndOfString() ([]byte, error) {
	buf := make([]byte, 0, 20)
	for {
		p := l.peakByte()
		if isStringMarker(p) {
			return buf, nil
		}

		b, err := l.readByte()
		if err != nil {
			return nil, NewError(
				INVALID_STRING,
				l.currentLine,
				l.currentChar,
				string(buf),
				)
		}

		buf = append(buf, b)
	}
}

func (l *Lexer) skipLine() error {
	for {
		b, err := l.readByte()
		if err != nil {
			return err
		}

		if b == '\n' {
			l.addLineCount()
			return nil
		}
	}
}

func (l *Lexer) addLineCount() {
	l.currentChar = 0
	l.currentLine = l.currentLine + 1
}

// resets peak to read index + 1
func (l *Lexer) readByte() (byte, error) {
	if l.r >= (l.len - 1) {
		return byte(0), io.EOF
	}

	l.r = l.r + 1

	// set current char for error messages
	l.currentChar = l.currentChar + 1

	return l.buf[l.r], nil
}

func (l *Lexer) peakByte() byte {
	if (l.r + 1) >= l.len {
		return byte(0)
	}

	return l.buf[l.r + 1]
}
