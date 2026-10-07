package table

type Table struct {
	Name			string
	RowCount	int16
	Columns		[]*Column
}

type Column struct {
	Cid						int16
	Name					string
	Type					string
	IsNullable		bool
	DefaultValue	string
	IsKey					bool
}
