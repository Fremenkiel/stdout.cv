package stdoutcv

import (
	"math/rand"
	"runtime"
	"testing"
)


func BenchmarkMap(b *testing.B) {
	for b.Loop() {
		num := rand.Intn(len(slice))
		token := slice[num]
		runtime.KeepAlive(token)

		_, ok := tokenMap[token]
		if !ok {
			panic("unable to find token")
		}
		runtime.KeepAlive(ok)
	}
}

func BenchmarkSwitch(b *testing.B) {
	for b.Loop() {
		num := rand.Intn(len(slice))
		token := slice[num]
		runtime.KeepAlive(token)

		ok := exists(token)
		if !ok {
			panic("unable to find token")
		}
		runtime.KeepAlive(ok)
	}
}

type TokenType = string

var tokenMap = map[TokenType]struct{}{
	EOF: {},
	STRING: {},
	NUMBER: {},
	FIELD: {},
	FUNC: {},
	AND: {},
	OR: {},
	L_THAN: {},
	G_THAN: {},
	EQUAL: {},
	N_EQUAL: {},
	LT_EQUAL: {},
	GT_EQUAL: {},
	IS: {},
	NOT: {},
	LIKE: {},
	AS: {},
	IN: {},
	IF: {},
	EXISTS: {},
	SEMICOL: {},
	COMMA: {},
	L_PARAN: {},
	R_PARAN: {},
	SELECT: {},
	UPDATE: {},
	INSERT: {},
	DELETE: {},
	DROP: {},
	ALTER: {},
	CREATE: {},
	INTO: {},
	FROM: {},
	WHERE: {},
	VALUES: {},
	SET: {},
	TABLE: {},
	ADD: {},
	COLUMN: {},
	PRIMARY: {},
	KEY: {},
	NULL: {},
	WILDCARD: {},
}

func exists(token TokenType) bool {
	switch token {
	case EOF:
		return true
	case STRING:
		return true
	case NUMBER:
		return true
	case FIELD:
		return true
	case FUNC:
		return true
	case AND:
		return true
	case OR:
		return true
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
	case IS:
		return true
	case NOT:
		return true
	case LIKE:
		return true
	case AS:
		return true
	case IN:
		return true
	case IF:
		return true
	case EXISTS:
		return true
	case SEMICOL:
		return true
	case COMMA:
		return true
	case L_PARAN:
		return true
	case R_PARAN:
		return true
	case SELECT:
		return true
	case UPDATE:
		return true
	case INSERT:
		return true
	case DELETE:
		return true
	case DROP:
		return true
	case ALTER:
		return true
	case CREATE:
		return true
	case INTO:
		return true
	case FROM:
		return true
	case WHERE:
		return true
	case VALUES:
		return true
	case SET:
		return true
	case TABLE:
		return true
	case ADD:
		return true
	case COLUMN:
		return true
	case PRIMARY:
		return true
	case KEY:
		return true
	case NULL:
		return true
	case WILDCARD:
		return true
	}

	return false
}

var slice = []TokenType{
	EOF,
	STRING,
	NUMBER,
	FIELD,
	FUNC,
	AND,
	OR,
	L_THAN,
	G_THAN,
	EQUAL,
	N_EQUAL,
	LT_EQUAL,
	GT_EQUAL,
	IS,
	NOT,
	LIKE,
	AS,
	IN,
	IF,
	EXISTS,
	SEMICOL,
	COMMA,
	L_PARAN,
	R_PARAN,
	SELECT,
	UPDATE,
	INSERT,
	DELETE,
	DROP,
	ALTER,
	CREATE,
	INTO,
	FROM,
	WHERE,
	VALUES,
	SET,
	TABLE,
	ADD,
	COLUMN,
	PRIMARY,
	KEY,
	NULL,
	WILDCARD,
}

var (
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
	IF				TokenType = "IF"										// IF operator
	EXISTS		TokenType = "EXISTS"								// EXISTS operator

	// Delimitor
	SEMICOL		TokenType = "SEMICOL"								// ;
	COMMA			TokenType = "COMMA"									// ,
	L_PARAN		TokenType = "L_PARAN"								// (
	R_PARAN		TokenType = "R_PARAN"								// )

	// Keywords
	SELECT		TokenType = "SELECT"
	UPDATE		TokenType = "UPDATE"
	INSERT		TokenType = "INSERT"
	DELETE		TokenType = "DELETE"
	DROP			TokenType = "DROP"
	ALTER			TokenType = "ALTER"
	CREATE		TokenType = "CREATE"

	// Expressions
	INTO 			TokenType = "INTO"
	FROM 			TokenType = "FROM"
	WHERE 		TokenType = "WHERE"
	VALUES 		TokenType = "VALUES"
	SET 			TokenType = "SET"
	TABLE 		TokenType = "TABLE"
	ADD				TokenType = "ADD"
	COLUMN		TokenType = "COLUMN"
	PRIMARY		TokenType = "PRIMARY"
	KEY				TokenType = "KEY"

	// VALUES
	NULL			TokenType = "NULL"
	WILDCARD	TokenType = "WILDCARD"						// *
)
