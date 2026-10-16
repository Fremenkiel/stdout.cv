package lexer

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

type TokenType int

const (
	EOF TokenType = iota	// End of file
	STRING								// Strings marked by ''
	FIELD									// Column or table name
	FUNC									// Build in func call, e.g. now()

	// Operator
	AND										// AND operator
	OR										// OR operator
	L_THAN								// > operator
	S_THAN								// < operator
	EQUAL									// = operator
	N_EQUAL								// != operator
	LT_EQUAL							// >= operator
	ST_EQUAL							// <= operator
	IS										// IS operator
	NOT										// NOT operator
	LIKE									// LIKE operator

	// Delimitor
	SEMICOL								// ;
	L_PARAN								// (
	R_PARAN								// )

	// Keywords
	SELECT
	UPDATE
	INSERT
	DELETE

	// Expressions
	INTO
	FROM
	WHERE
	VALUES
	SET
)

type Token struct {
	Type	TokenType
	Text	string
}
