package interpreter

func isSeperatorChar(c byte) bool {
	return isWhiteSpace(c) || isEndChar(c) || isLineBreak(c)
}

func isWhiteSpace(c byte) bool {
	switch c {
	case ' ':
		return true
	case '\t':
		return true
	}

	return false
}

func isEndChar(c byte) bool {
	return c == ';'
}

func isLineBreak(c byte) bool {
	return  c == '\n'
}

func isInt(c byte) bool {
	if c == '-' {
		return true
	}

	return isNumber(c)
}

func isNumber(c byte) bool {
	if c == '.' {
		return true
	}

	c -= '0'
	if c > 9 {
		return false
	}

	return true
}

func isStringMarker(c byte) bool {
	return c == '\''
}

func isIdentifierChar(c byte, isFirstChar bool) bool {
	switch {
	case c >= 'a' && c <= 'z':
		return true
	case c >= 'A' && c <= 'Z':
		return true
	case c == '_':
		return true
	case c >= 0x80: // SQLite treats any byte >= 0x80 as a letter
		return true
	case !isFirstChar && c >= '0' && c <= '9':
		return true
	}

	return false
}

func isQuotedIdentifierChar(c byte) bool {
	switch c {
	case '"':
		return true
	case '[':
		return true
	case ']':
		return true
	case '`':
		return true
	}

	return false
}

func getEndQuotedIdentifiersChar(c byte) byte {
	switch c {
	case '"':
		return '"'
	case '[':
		return ']'
	case '`':
		return '`'
	}

	return byte(0)
}
