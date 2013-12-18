package common

type Table struct {
	Name    string
	DbType  string /* mysql, postgres, sqlite, ... */
	Columns []*Column
}

type Column struct {
	TableName    string
	Name         string
	Type         *Type
	RawType      string
	Length       int
	Null         bool
	PrimaryKey   bool
	AutoIncr     bool
	Default      interface{}
	NeedsQuoting bool

	/* how to select the column */
	Select string
}
