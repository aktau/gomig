package common

type Table struct {
	Name    string
	Columns []*Column
}

type Column struct {
	TableName  string
	Name       string
	Type       string
	Length     int
	Null       bool
	PrimaryKey bool
	AutoIncr   bool
	Default    interface{}

	/* how to select the column */
	Select string
}
