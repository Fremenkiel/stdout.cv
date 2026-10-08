package column

type Column struct {
	Cid						int16
	Name					string
	Type					string
	IsNullable		bool
	DefaultValue	string
	IsPrimaryKey	bool
	IsForeignKey	bool
}
