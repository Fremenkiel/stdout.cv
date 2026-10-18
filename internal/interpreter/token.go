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
