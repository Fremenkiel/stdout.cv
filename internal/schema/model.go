package schema

type Schema struct {
	Cid						int16
	Name					string
	Type					string
	IsNullable		bool
	DefaultValue	string
	IsKey					bool
}
