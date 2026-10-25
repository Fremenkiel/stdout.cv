package lib

type Error struct {
	Message					string
	Code						string
	Duration				int64
	LineNumber			int
	CharStartNumber int
	CharEndNumber 	int
}
