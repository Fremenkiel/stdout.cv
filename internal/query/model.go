package query

const (
	SELECT	string = "SELECT"
	DELETE	string = "DELETE"
	INSERT	string = "INSERT"
	UPDATE	string = "UPDATE"

	// Modifiers
	INTO		string = "INTO"
	FROM		string = "FROM"
	WHERE		string = "WHERE"
	SET			string = "SET"
	VALUES	string = "VALUES"
)

var queryTypes []string = []string{
	SELECT,
	DELETE,
	INSERT,
	UPDATE,
}

var queryModifiers []string = []string{
	INTO,
	FROM,
	WHERE,
	SET,
	VALUES,
}

type Query struct {
	Fields					[]string
	WildcardIndex		int16
	TableName				string
	Type						string
	Modifiers				[]string
}
