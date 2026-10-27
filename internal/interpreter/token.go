package interpreter 

type Token struct {
	Type				TokenType
	Text				string
	LineNumber,
	Start,
	End					int
}

func (t Token) isField() bool {
	_, ok := fieldMap[t.Type]

	return ok
}

var fieldMap = map[TokenType]struct{}{
	FIELD: {},
	STRING: {},
	NUMBER: {},
	NULL: {},
	WILDCARD: {},
}

func (t Token) isOperator() bool {
	_, ok := operatorMap[t.Type]

	return ok
}

var operatorMap = map[TokenType]struct{}{
	L_THAN: {},
	G_THAN: {},
	EQUAL: {},
	N_EQUAL: {},
	LT_EQUAL: {},
	GT_EQUAL: {},
	LIKE: {},
	IS: {},
	IN: {},
	NOT: {},
	IF: {},
	EXISTS: {},
	ADD: {},
	RENAME: {},
	ALTER: {},
	DROP: {},
	COLUMN: {},
	CONSTRAINT: {},
	CHECK: {},
	CREATE: {},
	TO: {},
	PRIMARY: {},
	KEY: {},
}

func (t Token) isTablePrefix() bool {
	_, ok := tablePrefixMap[t.Type]

	return ok
}

var tablePrefixMap = map[TokenType]struct{}{
	INTO: {},
	TABLE: {},
	FROM: {},
	UPDATE: {},
	DROP: {},
}

type TokenType string

const (
	EOF					TokenType = "EOF"										// End of file
	STRING			TokenType = "STRING"								// Strings marked by ''
	NUMBER			TokenType = "NUMBER"
	FIELD				TokenType = "FIELD"									// column or table name

	// Operator
	AND					TokenType = "AND"										// AND operator
	OR					TokenType = "OR"										// OR operator
	L_THAN			TokenType = "L_THAN"								// < operator
	G_THAN			TokenType = "G_THAN"								// > operator
	EQUAL				TokenType = "EQUAL"									// = operator
	N_EQUAL			TokenType = "N_EQUAL"								// != operator
	LT_EQUAL		TokenType = "LT_EQUAL"							// <= operator
	GT_EQUAL		TokenType = "GT_EQUAL"							// >= operator
	IS					TokenType = "IS"										// IS operator
	NOT					TokenType = "NOT"										// NOT operator
	LIKE				TokenType = "LIKE"									// LIKE operator
	AS					TokenType = "AS"										// AS operator
	IN					TokenType = "IN"										// IN operator
	IF					TokenType = "IF"										// IF operator
	EXISTS			TokenType = "EXISTS"								// EXISTS operator

	// 	Delimitor
	SEMICOL			TokenType = "SEMICOL"								// ;
	COMMA				TokenType = "COMMA"									// ,
	L_PARAN			TokenType = "L_PARAN"								// (
	R_PARAN			TokenType = "R_PARAN"								// )

	// Keywords
	SELECT			TokenType = "SELECT"
	UPDATE			TokenType = "UPDATE"
	INSERT			TokenType = "INSERT"
	DELETE			TokenType = "DELETE"
	DROP				TokenType = "DROP"
	ALTER				TokenType = "ALTER"
	CREATE			TokenType = "CREATE"

	// Expressions
	INTO 				TokenType = "INTO"
	FROM 				TokenType = "FROM"
	WHERE 			TokenType = "WHERE"
	VALUES 			TokenType = "VALUES"
	SET 				TokenType = "SET"
	TABLE 			TokenType = "TABLE"
	ADD					TokenType = "ADD"
	RENAME			TokenType = "RENAME"
	COLUMN			TokenType = "COLUMN"
	TO					TokenType = "TO"
	PRIMARY			TokenType = "PRIMARY"
	KEY					TokenType = "KEY"
	CONSTRAINT	TokenType = "CONSTRAINT"
	CHECK				TokenType = "CHECK"

	// VALUES
	NULL			TokenType = "NULL"
	WILDCARD	TokenType = "WILDCARD"						// *
)

// todo: change to switch with default -1 (Kevin)
var precedenceMap = map[TokenType]int{
	EOF: -1,
	SEMICOL: -1,
	COMMA: -1,
	WILDCARD: -1,

	INSERT: -1,
	SELECT: -1,
	UPDATE: -1,
	DELETE: -1,

	INTO: -1,
	FROM: -1,
	WHERE: -1,
	VALUES: -1,
	SET: -1,
	TABLE: -1,
	COLUMN: -1,
	CONSTRAINT: -1,
	CHECK: -1,

	FIELD: 0,
	STRING: 0,
	NUMBER: 0,
	NULL: 0,

	DROP: 1,
	ADD: 1,
	ALTER: 1,
	RENAME: 1,
	CREATE: 1,

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
	IF: 1,
	EXISTS: 1,
	TO: 1,
	PRIMARY: 1,
	KEY: 1,

	AND: 2,

	OR: 3,

	L_PARAN: 4,
	R_PARAN: 4,
}
