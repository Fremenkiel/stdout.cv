package registry

type Registry = []*Table

type Table struct {
	Id					int
	CurrentName	string
}
