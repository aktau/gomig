package common

import (
	"database/sql"
	"io"
)

type Queryer interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type Reader interface {
	Queryer

	TableNames() []string

	/* FilteredTables() is more performant than Tables() if you
	 * only need a few tables */
	Tables() []*Table
	FilteredTables(incl, excl map[string]bool) []*Table

	Read(table *Table) (*sql.Rows, error)
	CreateView(name string, body string) error
	DropView(name string) error

	CreateProjection(name string, body string, engine string, pk []string, uks [][]string) error
	DropProjection(name string) error
}

type ReadCloser interface {
	io.Closer
	Reader
}
