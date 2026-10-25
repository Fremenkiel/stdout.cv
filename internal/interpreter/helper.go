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

