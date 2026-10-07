package viewmodels

type ColumnList []Column

type Column struct {
	Name					string
	Type					string
	IsNullable		bool
	IsPrimaryKey	bool
	IsForeignKey	bool
}
