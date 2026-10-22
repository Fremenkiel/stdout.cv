package interpreter 

// const (
// 	SELECT	string = "SELECT"
// 	DELETE	string = "DELETE"
// 	INSERT	string = "INSERT"
// 	UPDATE	string = "UPDATE"
//
// 	// Modifiers
// 	INTO		string = "INTO"
// 	FROM		string = "FROM"
// 	WHERE		string = "WHERE"
// 	SET			string = "SET"
// 	VALUES	string = "VALUES"
//
// 	COMMENT	string = "COMMENT"
// )
//
// var queryTypes []string = []string{
// 	SELECT,
// 	DELETE,
// 	INSERT,
// 	UPDATE,
// }
//
// var queryModifiers []string = []string{
// 	INTO,
// 	FROM,
// 	WHERE,
// 	SET,
// 	VALUES,
// }
//
// type Query struct {
// 	Fields						[]string
// 	WildcardIndex			int16
// 	TableName					string
// 	Type							string
// 	Modifiers					[]string
// }

type Token struct {
	Type				TokenType
	Text				string
	LineNumber,
	Start,
	End					int
}

func NewToken(t TokenType, text string, lexer *Lexer) Token {
	endPos := lexer.currentChar

	// removing the string marker as it's ignored in the string
	if t == STRING {
		endPos = endPos - 1
	}

	startPos := endPos - len(text)

	return Token{ 
		Type: t,
		Text: text,
		LineNumber: lexer.currentLine,
		Start: startPos,
		End: endPos,
	}
}

func (t Token) isField() bool {
	switch t.Type {
	case FIELD:
		return true
	case STRING:
		return true
	case NUMBER:
		return true
	case FUNC:
		return true
	case NULL:
		return true
	}

	return false
}

func (t Token) isOperator() bool {
	switch t.Type {
	case L_THAN:
		return true
	case G_THAN:
		return true
	case EQUAL:
		return true
	case N_EQUAL:
		return true
	case LT_EQUAL:
		return true
	case GT_EQUAL:
		return true
	case LIKE:
		return true
	case IS:
		return true
	case IN:
		return true
	case NOT:
		return true
	}

	return false
}

type TokenType string

const (
	EOF				TokenType = "EOF"										// End of file
	STRING		TokenType = "STRING"								// Strings marked by ''
	NUMBER		TokenType = "NUMBER"
	FIELD			TokenType = "FIELD"									// column or table name
	FUNC			TokenType = "FUNC"									// Build in func call, e.g. now()

	// Operator
	AND				TokenType = "AND"										// AND operator
	OR				TokenType = "OR"										// OR operator
	L_THAN		TokenType = "L_THAN"								// < operator
	G_THAN		TokenType = "G_THAN"								// > operator
	EQUAL			TokenType = "EQUAL"									// = operator
	N_EQUAL		TokenType = "N_EQUAL"								// != operator
	LT_EQUAL	TokenType = "LT_EQUAL"							// <= operator
	GT_EQUAL	TokenType = "GT_EQUAL"							// >= operator
	IS				TokenType = "IS"										// IS operator
	NOT				TokenType = "NOT"										// NOT operator
	LIKE			TokenType = "LIKE"									// LIKE operator
	AS				TokenType = "AS"										// AS operator
	IN				TokenType = "IN"										// IN operator

	// Delimitor
	SEMICOL		TokenType = "SEMICOL"								// ;
	L_PARAN		TokenType = "L_PARAN"								// (
	R_PARAN		TokenType = "R_PARAN"								// )

	// Keywords
	SELECT		TokenType = "SELECT"
	UPDATE		TokenType = "UPDATE"
	INSERT		TokenType = "INSERT"
	DELETE		TokenType = "DELETE"

	// Expressions
	INTO 			TokenType = "INTO"
	FROM 			TokenType = "FROM"
	WHERE 		TokenType = "WHERE"
	VALUES 		TokenType = "VALUES"
	SET 			TokenType = "SET"

	// VALUES
	NULL			TokenType = "NULL"
)

var precedenceMap = map[TokenType]int{
	EOF: -1,
	SEMICOL: -1,

	INSERT: -1,
	SELECT: -1,
	UPDATE: -1,
	DELETE: -1,

	INTO: -1,
	FROM: -1,
	WHERE: -1,
	VALUES: -1,
	SET: -1,

	FIELD: 0,
	STRING: 0,
	NUMBER: 0,
	NULL: 0,

	EQUAL: 1,
	L_THAN: 1,
	G_THAN: 1,
	N_EQUAL: 1,
	LT_EQUAL: 1,
	GT_EQUAL: 1,
	IS: 1,
	NOT: 1,
	LIKE: 1,
	IN: 1,

	AND: 2,

	OR: 3,

	L_PARAN: 4,
	R_PARAN: 4,
}
