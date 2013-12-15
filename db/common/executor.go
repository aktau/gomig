package common

import (
	"database/sql"
	"io"
)

const (
	CapCopy = iota
)

type Executor interface {
	io.Closer

	Transaction(name string, statements []string) error
	Multiple(name string, statements []string) []error
	Single(name string, statement string) error

	/* for e.g. PostgreSQL COPY support */
	HasCapability(capability int) bool

	/* will return the underlying sql.DB if it exists, otherwise nil */
	GetDb() *sql.DB
}
