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

	ListTables() []Table
	CreateView(name string, body string) error
	DropView(name string) error
}

type ReadCloser interface {
	io.Closer
	Reader
}
