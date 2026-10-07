package viewmodels

type ColumnList []Column

type Column struct {
	Name				string
	Type				string
	PrivateKey,
	ForeignKey	bool
}
