package viewmodels

type Error struct {
	Message					string
	Code						string
	Helpertext			string
	Duration				int64
	LineNumber			int
	CharStartNumber int
	CharEndNumber 	int
}
